package dockercli

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/fimreal/goutils/ezap"
)

// Containers2JSON 将容器详细信息打印为 JSON 格式
func Containers2JSON(containers []types.ContainerJSON) string {
	output, err := json.MarshalIndent(containers, "", "  ")
	if err != nil {
		ezap.Fatalf("Error marshaling containers to JSON: %v", err)
	}
	return string(output)
}

// Containers2CMD 将容器详细信息打印为 docker run 格式
func Containers2CMD(containersJSON []types.ContainerJSON, pretty bool) string {
	var command strings.Builder
	for c, containerJSON := range containersJSON {
		if c != 0 {
			command.WriteString("\n\n")
		}
		command.WriteString(buildDockerRunCommand(containerJSON, pretty))
	}
	return command.String()
}

// buildDockerRunCommand 构建 docker run 命令字符串
func buildDockerRunCommand(containerJSON types.ContainerJSON, pretty bool) string {
	var command strings.Builder
	end := " "
	if pretty {
		end = " \\\n"
	}

	cname := strings.TrimPrefix(containerJSON.Name, "/")
	created := containerJSON.Created

	// description
	command.WriteString(fmt.Sprintf("# Container name: %s\n", cname))
	command.WriteString(fmt.Sprintf("# Created at: %s\n", created))
	command.WriteString(fmt.Sprintf("# Description: %s\n", containerJSON.Config.Labels["description"]))

	// docker run
	command.WriteString("docker run -d" + end)

	// container name
	command.WriteString("--name " + cname + end)

	// hostname
	hostname := containerJSON.Config.Hostname
	csha := containerJSON.ID[:12]
	if hostname != csha {
		command.WriteString("--hostname " + hostname + end)
	}

	// restart policy
	if containerJSON.HostConfig.RestartPolicy.Name != "" && containerJSON.HostConfig.RestartPolicy.Name != "no" {
		command.WriteString(fmt.Sprintf("--restart %s%s", containerJSON.HostConfig.RestartPolicy.Name, end))
		if containerJSON.HostConfig.RestartPolicy.MaximumRetryCount > 0 {
			command.WriteString(fmt.Sprintf("--restart-max-attempts %d%s", containerJSON.HostConfig.RestartPolicy.MaximumRetryCount, end))
		}
	}

	// user
	if containerJSON.Config.User != "" {
		command.WriteString(fmt.Sprintf("--user %s%s", containerJSON.Config.User, end))
	}

	// workdir
	if containerJSON.Config.WorkingDir != "" {
		command.WriteString("--workdir " + containerJSON.Config.WorkingDir + end)
	}

	// entrypoint
	if len(containerJSON.Config.Entrypoint) > 0 {
		command.WriteString("--entrypoint " + strings.Join(containerJSON.Config.Entrypoint, " ") + end)
	}

	// environment variables
	for _, env := range containerJSON.Config.Env {
		if env == "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin" {
			continue
		}
		command.WriteString(fmt.Sprintf("-e %s%s", env, end))
	}

	// host-to-IP mapping
	for _, hosts := range containerJSON.HostConfig.ExtraHosts {
		command.WriteString("--add-host " + hosts + end)
	}

	// privileged mode
	if containerJSON.HostConfig.Privileged {
		command.WriteString("--privileged" + end)
	}

	// capabilities
	for _, cap := range containerJSON.HostConfig.CapAdd {
		command.WriteString("--cap-add " + cap + end)
	}
	for _, cap := range containerJSON.HostConfig.CapDrop {
		command.WriteString("--cap-drop " + cap + end)
	}

	// readonly root fs
	if containerJSON.HostConfig.ReadonlyRootfs {
		command.WriteString("--read-only " + end)
	}

	// OOM score adjustment
	if containerJSON.HostConfig.OomScoreAdj != 0 {
		command.WriteString("--oom-score-adj=" + strconv.Itoa(containerJSON.HostConfig.OomScoreAdj) + end)
	}

	// user namespace mode
	if containerJSON.HostConfig.UsernsMode != "" {
		command.WriteString(fmt.Sprintf("--userns-mode=%s%s", containerJSON.HostConfig.UsernsMode, end))
	}

	// ipc mode
	if containerJSON.HostConfig.PidMode != "" {
		command.WriteString(fmt.Sprintf("--pid %s%s", containerJSON.HostConfig.PidMode, end))
	}

	// link
	for _, link := range containerJSON.HostConfig.Links {
		linkParts := strings.Split(link, ":")
		sourceContainer := strings.TrimPrefix(linkParts[0], "/")
		if len(linkParts) == 2 {
			aliasName := linkParts[1]
			command.WriteString(fmt.Sprintf("--link %s:%s%s", sourceContainer, aliasName, end))
		} else {
			command.WriteString(fmt.Sprintf("--link %s%s", sourceContainer, end))
		}
	}

	// cpu limit
	if containerJSON.HostConfig.NanoCPUs > 0 {
		// Convert NanoCPUs to a value for --cpus
		cpuLimit := float64(containerJSON.HostConfig.NanoCPUs) / 1e9 // convert to CPUs
		command.WriteString(fmt.Sprintf("--cpus=%.2f%s", cpuLimit, end))
	}

	// cpu shares
	if containerJSON.HostConfig.CPUShares > 0 {
		command.WriteString(fmt.Sprintf("--cpu-shares=%d%s", containerJSON.HostConfig.CPUShares, end))
	}

	// cpuset
	if containerJSON.HostConfig.CpusetCpus != "" {
		command.WriteString("--cpuset-cpus " + containerJSON.HostConfig.CpusetCpus + end)
	}

	// memory limit
	if containerJSON.HostConfig.Memory > 0 {
		command.WriteString("--memory=" + strconv.FormatInt(containerJSON.HostConfig.Memory, 10) + end)
	}

	// network mode
	if containerJSON.HostConfig.NetworkMode != "" && containerJSON.HostConfig.NetworkMode != "default" {
		command.WriteString(fmt.Sprintf("--network %s%s", containerJSON.HostConfig.NetworkMode, end))
	}

	// dns
	for _, dns := range containerJSON.HostConfig.DNS {
		command.WriteString("--dns " + dns + end)
	}

	// port mapping
	addedPorts := make(map[string]bool)
	for port, bindings := range containerJSON.NetworkSettings.Ports {
		for _, binding := range bindings {
			var portMapping string

			if binding.HostIP != "0.0.0.0" && binding.HostIP != "" && binding.HostIP != "::" {
				portMapping = fmt.Sprintf("%s:%s:%s", binding.HostIP, binding.HostPort, port.Port())
			} else {
				portMapping = fmt.Sprintf("%s:%s", binding.HostPort, port.Port())
			}

			// Check if the port mapping has already been added
			if _, exists := addedPorts[portMapping]; !exists {
				command.WriteString("-p " + portMapping + end)
				addedPorts[portMapping] = true // Mark as added
			}
		}
	}

	// mount
	for _, mount := range containerJSON.Mounts {
		if mount.Type == "bind" {
			// Bind mount
			command.WriteString(fmt.Sprintf("-v %s:%s%s", mount.Source, mount.Destination, end))
		} else {
			// Volume mount
			command.WriteString(fmt.Sprintf("--mount type=%s,source=%s,target=%s%s", mount.Type, mount.Source, mount.Destination, end))
		}
	}

	// devices
	for _, device := range containerJSON.HostConfig.Devices {
		command.WriteString(fmt.Sprintf("--device %s:%s%s", device.PathOnHost, device.PathInContainer, end))
	}

	// label
	for key, value := range containerJSON.Config.Labels {
		command.WriteString(fmt.Sprintf("--label %s=%s%s", key, value, end))
	}

	// log driver
	if containerJSON.HostConfig.LogConfig.Type != "" {
		if containerJSON.HostConfig.LogConfig.Type != "json-file" {
			command.WriteString(fmt.Sprintf("--log-driver %s%s", containerJSON.HostConfig.LogConfig.Type, end))
		}
	}
	if containerJSON.HostConfig.LogConfig.Config != nil {
		for key, value := range containerJSON.HostConfig.LogConfig.Config {
			command.WriteString(fmt.Sprintf("--log-opt %s=%s%s", key, value, end))
		}
	}

	// image
	command.WriteString(containerJSON.Config.Image + end)

	// command
	if len(containerJSON.Config.Cmd) > 0 {
		command.WriteString(strings.Join(containerJSON.Config.Cmd, " "))
	}

	return command.String()
}

// Containers2Compose 将容器详细信息打印为 docker-compose 格式
func Containers2Compose(containersJSON []types.ContainerJSON) map[string]string {
	services := make(map[string]string)

	for _, containerJSON := range containersJSON {
		var compose strings.Builder
		compose.WriteString("version: '3'\n")
		compose.WriteString("services:\n")
		compose.WriteString(generateServiceConfig(containerJSON))

		cname := strings.TrimPrefix(containerJSON.Name, "/")
		services[cname] = compose.String()
	}

	return services
}

// generateServiceConfig 生成单个服务的配置
func generateServiceConfig(containerJSON types.ContainerJSON) string {
	var serviceConfig strings.Builder
	cname := strings.TrimPrefix(containerJSON.Name, "/")

	// Service name
	serviceConfig.WriteString(fmt.Sprintf("%s:\n", cname))

	// image
	serviceConfig.WriteString(fmt.Sprintf("  image: %s\n", containerJSON.Config.Image))

	// command
	if len(containerJSON.Config.Cmd) > 0 {
		serviceConfig.WriteString(fmt.Sprintf("  command: [%s]\n", strings.Join(containerJSON.Config.Cmd, ", ")))
	}

	// environment
	if len(containerJSON.Config.Env) > 0 {
		serviceConfig.WriteString("  environment:\n")
		for _, env := range containerJSON.Config.Env {
			serviceConfig.WriteString(fmt.Sprintf("    - %s\n", env))
		}
	}

	// workdir
	if containerJSON.Config.WorkingDir != "" {
		serviceConfig.WriteString(fmt.Sprintf("  working_dir: %s\n", containerJSON.Config.WorkingDir))
	}

	// entrypoint
	if len(containerJSON.Config.Entrypoint) > 0 {
		serviceConfig.WriteString(fmt.Sprintf("  entrypoint: [%s]\n", strings.Join(containerJSON.Config.Entrypoint, ", ")))
	}

	// hostname
	if containerJSON.Config.Hostname != "" {
		serviceConfig.WriteString(fmt.Sprintf("  hostname: %s\n", containerJSON.Config.Hostname))
	}

	// user
	if containerJSON.Config.User != "" {
		serviceConfig.WriteString(fmt.Sprintf("  user: %s\n", containerJSON.Config.User))
	}

	// Privileged mode
	if containerJSON.HostConfig.Privileged {
		serviceConfig.WriteString("  privileged: true\n")
	}

	// restart policy
	if containerJSON.HostConfig.RestartPolicy.Name != "" && containerJSON.HostConfig.RestartPolicy.Name != "no" {
		serviceConfig.WriteString(fmt.Sprintf("  restart: %s\n", containerJSON.HostConfig.RestartPolicy.Name))
	}

	// port mapping
	if len(containerJSON.NetworkSettings.Ports) > 0 {
		serviceConfig.WriteString("  ports:\n")
		for port, bindings := range containerJSON.NetworkSettings.Ports {
			for _, binding := range bindings {
				if binding.HostIP == "0.0.0.0" || binding.HostIP == "" || binding.HostIP == "::" {
					serviceConfig.WriteString(fmt.Sprintf("    - \"%s:%s\"\n", binding.HostPort, port.Port()))
				} else {
					serviceConfig.WriteString(fmt.Sprintf("    - \"%s:%s:%s\"\n", binding.HostIP, binding.HostPort, port.Port()))
				}
			}
		}
	}

	// volume
	if len(containerJSON.Mounts) > 0 {
		serviceConfig.WriteString("  volumes:\n")
		for _, mount := range containerJSON.Mounts {
			if mount.Type == "bind" {
				serviceConfig.WriteString(fmt.Sprintf("    - \"%s:%s\"\n", mount.Source, mount.Destination))
			} else {
				serviceConfig.WriteString(fmt.Sprintf("    - \"%s:%s\"\n", mount.Source, mount.Destination)) // Volume mounts can also be handled here
			}
		}
	}

	// capabilities
	if len(containerJSON.HostConfig.CapAdd) > 0 {
		serviceConfig.WriteString("  cap_add:\n")
		for _, cap := range containerJSON.HostConfig.CapAdd {
			serviceConfig.WriteString(fmt.Sprintf("    - %s\n", cap))
		}
	}
	if len(containerJSON.HostConfig.CapDrop) > 0 {
		serviceConfig.WriteString("  cap_drop:\n")
		for _, cap := range containerJSON.HostConfig.CapDrop {
			serviceConfig.WriteString(fmt.Sprintf("    - %s\n", cap))
		}
	}

	// OOM score adjustment
	if containerJSON.HostConfig.OomScoreAdj != 0 {
		serviceConfig.WriteString(fmt.Sprintf("  oom_score_adj: %d\n", containerJSON.HostConfig.OomScoreAdj))
	}

	// User namespace mode
	if containerJSON.HostConfig.UsernsMode != "" {
		serviceConfig.WriteString(fmt.Sprintf("  userns_mode: %s\n", containerJSON.HostConfig.UsernsMode))
	}

	// IPC mode
	if containerJSON.HostConfig.IpcMode != "" {
		serviceConfig.WriteString(fmt.Sprintf("  ipc: %s\n", containerJSON.HostConfig.IpcMode))
	}

	return serviceConfig.String()
}
