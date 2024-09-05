package dockercli

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
)

// List 列出所有 Docker 容器
func (d *DockerClient) List(showAll bool) ([]types.Container, error) {
	ctx := context.Background()

	// 获取容器查询选项
	options := container.ListOptions{All: showAll}

	// 获取容器
	return d.cli.ContainerList(ctx, options)

}

// listPrint 打印每个容器的信息
func ListPrint(containerSummary []types.Container, pretty bool) {

	// 打印表头
	if pretty {
		fmt.Printf("%-12s %-15s %-20s %-10s %-25s\n", "CONTAINER ID", "NAMES", "CREATED", "STATUS", "IMAGE")
	} else {
		fmt.Println("CONTAINER ID        NAMES         CREATED         STATUS     IMAGE")
	}

	for _, containerSummary := range containerSummary {
		id := strings.TrimPrefix(containerSummary.ID, "sha256:")
		image := containerSummary.Image                                                     // 镜像名称
		createdTime := time.Unix(containerSummary.Created, 0).Format("2006-01-02T15:04:05") // 创建时间
		status := containerSummary.State                                                    // 状态
		name := strings.TrimPrefix(containerSummary.Names[0], "/")                          // 容器名称（可能有多个）

		if pretty {
			fmt.Printf("%-12s %-15s %-20s %-10s %-25s\n", id[:12], name, createdTime, status, image)
		} else {
			fmt.Printf("%-20s %-15s %-20s %-10s %-25s\n", id, name, createdTime, status, image)
		}
	}
}

// Find 查找指定容器
func (d *DockerClient) Find(containerNameOrID []string) (filteredContainers []types.Container, err error) {
	ctx := context.Background()

	// 获取容器查询选项
	options := container.ListOptions{All: true}

	// 获取容器
	containers, err := d.cli.ContainerList(ctx, options)
	if err != nil {
		return nil, err
	}

	// 手动过滤符合条件的容器
	for _, containerSummary := range containers {
		for _, f := range containerNameOrID {
			if strings.HasPrefix(containerSummary.ID, f) || strings.Contains(containerSummary.Names[0], f) {
				filteredContainers = append(filteredContainers, containerSummary)
			}
		}
	}
	return
}

func (d *DockerClient) FindImage(imageName string) (imagesSummary []image.Summary, err error) {
	ctx := context.Background()

	// 获取镜像查询选项
	options := image.ListOptions{All: true}

	// 获取镜像
	images, err := d.cli.ImageList(ctx, options)
	if err != nil {
		return nil, err
	}

	// 手动过滤符合条件的镜像
	for _, imageSummary := range images {
		if strings.Contains(imageSummary.RepoTags[0], imageName) {
			imagesSummary = append(imagesSummary, imageSummary)
		}
	}
	return
}
