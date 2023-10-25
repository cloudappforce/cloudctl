package cmd

import (
	"context"
	"fmt"

	"github.com/cloudappforce/cloudctl/internal/constants"
	api "github.com/cloudappforce/cloudctl/pkg/api"
	"github.com/cloudappforce/cloudctl/pkg/api/models"
	"github.com/cloudappforce/cloudctl/pkg/builder"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NOTES
// https://github.com/docker/cli/issues/2533

func Deploy() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "deploy",
		Short: "Create and configure deployments",
		RunE: func(cmd *cobra.Command, args []string) error {

			dockerfileName, err := cmd.Flags().GetString(constants.ArgDockerfile)
			if err != nil {
				return err
			}

			dockerArchive, err := cmd.Flags().GetString(constants.ArgDockerArchive)
			if err != nil {
				return err
			}
			dockerTag, err := cmd.Flags().GetString(constants.ArgDockerTag)
			if err != nil {
				return err
			}

			repo, err := cmd.Flags().GetString(constants.ArgDockerRepo)
			if err != nil {
				return err
			}
			if repo == "" {
				err = fmt.Errorf("missing docker repo")
				return err
			}

			// fetch database arguments [size, instances, type, name]
			dbSize, err := cmd.Flags().GetString(constants.ArgDBSize)
			if err != nil {
				return err
			}

			dbInstances, err := cmd.Flags().GetString(constants.ArgDBInstances)
			if err != nil {
				return err
			}

			dbType, err := cmd.Flags().GetString(constants.ArgDBType)
			if err != nil {
				return err
			}
			dbName, err := cmd.Flags().GetString(constants.ArgDBName)
			if err != nil {
				return err
			}

			dbLabels, err := cmd.Flags().GetStringToString(constants.ArgDBLabels)
			if err != nil {
				return err
			}

			_, err = builder.BuildDocker(context.Background(), builder.DockerOptions{
				ArchivePath: dockerArchive,
				ArchiveOpts: &archive.TarOptions{},
				Dockerfile:  dockerfileName,
				ImageTags:   []string{dockerTag},
			})
			if err != nil {
				return err
			}

			repoUsername := viper.GetString("REPO_USERNAME")
			repoPassword := viper.GetString("REPO_PASSWORD")

			auth, err := registry.EncodeAuthConfig(registry.AuthConfig{
				Username:      repoUsername,
				Password:      repoPassword,
				ServerAddress: "registry.cloudappforce.com",
			})
			if err != nil {
				panic(err)
			}

			sourceImage := dockerTag
			targetImage := fmt.Sprintf("registry.cloudappforce.com/%s/%s", repo, dockerTag)
			err = builder.TagImage(cmd.Context(), sourceImage, targetImage)
			if err != nil {
				panic(err)
			}
			err = builder.PushDocker(cmd.Context(), targetImage, types.ImagePushOptions{
				RegistryAuth: auth,
			})

			if err != nil {
				panic(err)
			}

			// generate database
			client := api.Client{
				Host:   viper.GetString("HOST"),
				Scheme: viper.GetString("SCHEME"),
				JWT:    viper.GetString("ACCESS_TOKEN"),
			}

			if dbLabels == nil {
				panic("no labels")
			}
			fmt.Println(dbLabels)

			createDatabaseInput := &api.CreateDatabaseInput{
				Size:      dbSize,
				Instances: dbInstances,
				Type:      dbType,
				Name:      dbName,
				Labels:    dbLabels,
			}

			fmt.Println(createDatabaseInput)
			database, err := client.CreateDatabase(cmd.Context(), createDatabaseInput)
			if err != nil {
				panic(err)
			}

			// inject the db-name-rw service as an environment variable if DATABASE_URL isn't provided
			databaseURLEnv := fmt.Sprintf("%s-rw", database.Name)

			response, err := client.CreateDeployment(cmd.Context(), &models.CreateDeploymentInput{
				Name:     "docker-image-nginx",
				Replicas: 1,
				Containers: []models.DeploymentContainerSpec{
					{
						Name:   "docker-image-nginx",
						CPU:    "1",
						Image:  targetImage,
						Memory: "500m",
						Environment: []models.Environment{
							{
								Name:  "DATABASE_URL",
								Value: databaseURLEnv,
							},
						},
					},
				},
			})
			fmt.Println(response)
			return err

		},
	}

	cmd.Flags().String(constants.ArgDeploymentName, "1", "deployment name")
	cmd.Flags().String(constants.ArgDeploymentReplicas, "", "deployment replica count")
	cmd.Flags().Bool(constants.ArgDeploymentEndpoint, true, "deployment public endpoint")
	cmd.Flags().String(constants.ArgDeploymentCPU, "1", "deployment cpu count")
	cmd.Flags().String(constants.ArgDeploymentMem, "1", "deployment mem count")
	cmd.Flags().StringToString(constants.ArgDeploymentLabels, map[string]string{}, "deployment labels")

	// DATABASE ARGS
	cmd.Flags().String(constants.ArgDBInstances, "1", "db instances")
	cmd.Flags().String(constants.ArgDBName, "", "db name")
	cmd.Flags().String(constants.ArgDBSize, "1Gi", "db size")
	cmd.Flags().String(constants.ArgDBType, "postgres", "db type")
	cmd.Flags().StringToString(constants.ArgDBLabels, map[string]string{"cloudctl": "0.0.1"}, "db labels")

	cmd.Flags().String(constants.ArgDockerfile, "Dockerfile", "dockerfile name")
	cmd.Flags().String(constants.ArgDockerRepo, "", "docker repo name")
	cmd.Flags().String(constants.ArgDockerArchive, "./", "path to the archive")
	cmd.Flags().String(constants.ArgDockerTag, fmt.Sprintf("%s:latest", namesgenerator.GetRandomName(0)), "docker tag")

	return cmd
}
