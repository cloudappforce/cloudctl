package cmd

import (
	"context"

	"github.com/cloudappforce/cloudctl/internal/constants"
	"github.com/cloudappforce/cloudctl/internal/formatter"
	api "github.com/cloudappforce/cloudctl/pkg/api"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Database() *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "database",
		Short: "Create and configure databases",
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here

		},
	}
	var createCmd = &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {

			input, err := api.NewCreateDatabaseInputFromFlags(cmd.Flags())
			if err != nil {
				panic(err)
			}

			client := api.Client{
				Scheme: viper.GetString("SCHEME"),
				Host:   viper.GetString("HOST"),
				JWT:    viper.GetString("ACCESS_TOKEN"),
			}

			client.CreateDatabase(context.Background(), input)
		},
	}
	var deleteCmd = &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			input, err := api.NewDeleteInputFromFlags(cmd.Flags())
			if err != nil {
				panic(err)
			}

			client := api.Client{
				Scheme: viper.GetString("SCHEME"),
				Host:   viper.GetString("HOST"),
				JWT:    viper.GetString("ACCESS_TOKEN"),
			}

			client.DeleteDatabase(context.Background(), input)
		},
	}

	var listCmd = &cobra.Command{
		Use: "list",
		Run: func(cmd *cobra.Command, args []string) {
			input, err := api.NewListDatabaseInputFromFlags(cmd.Flags())
			if err != nil {
				panic(err)
			}
			client := api.Client{
				Scheme: viper.GetString("SCHEME"),
				Host:   viper.GetString("HOST"),
				JWT:    viper.GetString("ACCESS_TOKEN"),
			}
			resp, err := client.ListDatabases(context.Background(), input)
			if err != nil {
				panic(err)
			}
			formatter.FormatDatabaseListResponse(resp)
		},
	}

	createCmd.Flags().String(constants.ArgSize, "1Gib", "required size for the database cluster")
	createCmd.Flags().String(constants.ArgInstances, "1", "number of instances for the database cluster")
	createCmd.Flags().String(constants.ArgType, "postgres", "type of database cluster [postgres]")
	createCmd.Flags().String(constants.ArgName, "", "name of the database cluster")
	createCmd.Flags().StringToString(constants.ArgLabels, map[string]string{"provisioner": "cloudctl"}, "labels to add to the label cluster")

	deleteCmd.Flags().String(constants.ArgName, "", "name of the database cluster")

	cmd.AddCommand(createCmd)
	cmd.AddCommand(deleteCmd)
	cmd.AddCommand(listCmd)
	return cmd
}
