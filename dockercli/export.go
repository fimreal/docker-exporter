package dockercli

import (
	"github.com/docker/docker/api/types"
)

// Export 导出容器 json 格式详细信息
func (d *DockerClient) ExportContainersJSON(dockerNameORID []string, showAll bool) ([]types.ContainerJSON, error) {
	// 按照传入条件查询容器，如果没有传入，则查询所有容器
	var containers []types.Container
	var err error
	if len(dockerNameORID) == 0 {
		containers, err = d.List(showAll)
		if err != nil {
			return nil, err
		}
	} else {
		containers, err = d.Find(dockerNameORID)
		if err != nil {
			return nil, err
		}
	}

	return d.Inspect(containers)

}

// 将容器信息格式化成指定格式
func ParseContainers(containersJSON []types.ContainerJSON, format string, pretty bool) interface{} {
	switch format {
	case "json":
		// string
		return Containers2JSON(containersJSON)
	case "compose", "yaml", "yml":
		// map[string]string
		return Containers2Compose(containersJSON)
	default:
		// string
		return Containers2CMD(containersJSON, pretty)
	}
}
