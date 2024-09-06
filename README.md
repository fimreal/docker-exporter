# docker-exporter

`docker-exporter` 是一个命令行工具，用于导出和显示 Docker 容器的配置信息和参数，支持远程连接，导出完整 `docker run` 命令或者 `docker compose` 格式 yaml 文件。

## 下载

见 [release](https://github.com/fimreal/docker-exporter/releases) 页面，选择对应平台下载。

## 使用

#### 帮助

```bash
docker-exporter -h
```

#### 查看容器列表

```bash
docker-exporter list [-H tcp://remote-host:2375] [-V client_version]
```

#### 导出容器配置

```bash
docker-exporter export [-H tcp://remote-host:2375] [-V client_version] [container_name|container_id] [-f yaml|cmd]
```


## 开发

#### 编译

详见 Makefile

```bash
make
```