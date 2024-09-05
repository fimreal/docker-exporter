// create docker api client
package dockercli

import (
	"github.com/docker/docker/client"
)

// DockerClient 结构体封装了 Docker 客户端及相关配置
type DockerClient struct {
	cli           *client.Client
	dockerHost    string
	clientVersion string
}

// NewCli 创建新的 DockerClient 实例
func NewCli(dockerHost, clientVersion string) (*DockerClient, error) {
	cli, err := client.NewClientWithOpts(client.WithHost(dockerHost), client.WithVersion(clientVersion))
	if err != nil {
		return nil, err
	}
	return &DockerClient{
		cli:           cli,
		dockerHost:    dockerHost,
		clientVersion: clientVersion,
	}, nil
}
