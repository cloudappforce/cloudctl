package builder

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

type DockerOptions struct {
	ArchivePath string
	ArchiveOpts *archive.TarOptions
	Dockerfile  string
	ImageTags   []string
}

// Build docker images from an archive
// Currently Docker client is loaded from environment
func BuildDocker(cxt context.Context, opts DockerOptions) (string, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return "", err
	}
	tar, err := archive.TarWithOptions(opts.ArchivePath, opts.ArchiveOpts)
	if err != nil {
		return "", err

	}

	buildOptions := types.ImageBuildOptions{
		Dockerfile: opts.Dockerfile,
		Tags:       opts.ImageTags,
		Remove:     true,
		Platform:   "linux/amd64",
		NoCache:    true,
	}

	res, err := cli.ImageBuild(context.Background(), tar, buildOptions)
	if err != nil {
		return "", err

	}

	defer res.Body.Close()

	return writeDockerBuildToConsole(context.Background(), res.Body)

}

func TagImage(ctx context.Context, source, target string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	return cli.ImageTag(ctx, source, target)
}

func PushDocker(ctx context.Context, image string, opts types.ImagePushOptions) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	res, err := cli.ImagePush(ctx, image, opts)
	if err != nil {
		return err
	}
	defer res.Close()
	// TODO handle response
	scanner := bufio.NewScanner(res)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	return nil
}

// writeDockerBuildToConsole writes the docker image build respponse to stdout
// caller is responsible for closing the ReadCloser.
func writeDockerBuildToConsole(ctx context.Context, body io.ReadCloser) (string, error) {
	imageID := ""
	scanner := bufio.NewScanner(body)
	for scanner.Scan() {
		m := map[string]interface{}{}
		data := scanner.Bytes()
		err := json.Unmarshal(data, &m)
		if err != nil {
			return "", err
		}
		if val, ok := m["stream"]; ok {
			fmt.Print(val)
		}
		if val, ok := m["aux"]; ok {
			// hack, decode to bytes then encode into map
			innerMap := map[string]interface{}{}
			auxData, err := json.Marshal(&val)
			if err != nil {
				return "", err
			}

			err = json.Unmarshal(auxData, &innerMap)
			if err != nil {
				return "", err
			}

			if id, hasID := innerMap["ID"]; hasID {
				// TODO: validate assertion is safe
				imageID = id.(string)
			}
		}
	}
	return imageID, nil
}
