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
// 如果 clientVersion 为空，则使用自动版本协商
func NewCli(dockerHost, clientVersion string) (*DockerClient, error) {
	opts := []client.Opt{client.FromEnv, client.WithHost(dockerHost)}

	// 只有指定了版本时才设置版本，否则使用自动协商
	if clientVersion != "" {
		opts = append(opts, client.WithVersion(clientVersion))
	} else {
		// 使用自动版本协商
		opts = append(opts, client.WithAPIVersionNegotiation())
	}

	cli, err := client.NewClientWithOpts(opts...)
	if err != nil {
		return nil, err
	}
	return &DockerClient{
		cli:           cli,
		dockerHost:    dockerHost,
		clientVersion: clientVersion,
	}, nil
}
