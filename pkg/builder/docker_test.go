package builder_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/cloudappforce/cloudctl/pkg/builder"
	"github.com/docker/docker/pkg/archive"
)

func TestCloudPlatformBuild(t *testing.T) {
	archivePath, err := os.Getwd()

	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(archivePath)
	opts := builder.DockerOptions{
		ArchivePath: archivePath,
		ArchiveOpts: &archive.TarOptions{},
		Dockerfile:  "Dockerfile",
		ImageTags:   []string{"cloudctl-python"},
	}
	result, err := builder.BuildDocker(context.Background(), opts)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(result)
}
