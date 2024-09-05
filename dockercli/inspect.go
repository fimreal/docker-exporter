package dockercli

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
)

// Inspect 查看容器详细信息
func (d *DockerClient) Inspect(containers []types.Container) (containersJSON []types.ContainerJSON, err error) {

	// 遍历并显示每个容器的详细信息
	for _, containerSummary := range containers {
		containerJSON, err := d.cli.ContainerInspect(context.Background(), containerSummary.ID)
		if err != nil {
			return nil, err
		}
		containersJSON = append(containersJSON, containerJSON)
	}

	return
}

func (d *DockerClient) InspectImageByID(imageID string) (imageJSON types.ImageInspect, err error) {
	imageJSON, _, err = d.cli.ImageInspectWithRaw(context.Background(), imageID)
	return
}

func (d *DockerClient) InspectImageByName(imageName string) (imageJSON types.ImageInspect, err error) {
	images, err := d.cli.ImageList(context.Background(), image.ListOptions{All: true})
	if err != nil {
		return
	}
	for _, image := range images {
		if image.RepoTags[0] == imageName {
			imageJSON, _, err = d.cli.ImageInspectWithRaw(context.Background(), image.ID)
			return
		}
	}
	return
}

// func (d *DockerClient) getImageEntrypoint(imageID string) (entrypoint []string, err error) {
// 	imageJSON, _, err := d.cli.ImageInspectWithRaw(context.Background(), imageID)
// 	if err != nil {
// 		return
// 	}
// 	return imageJSON.Config.Entrypoint, nil
// }
