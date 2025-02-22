# APISIX Metric

APISIX Metric 是一个专门为 Apache APISIX 设计的指标收集代理服务。它可以收集、转换和暴露 APISIX 的 Prometheus 指标数据，支持服务和路由的动态映射配置。

## 功能特性

- 提供 HTTP 服务，支持自定义监听端口（默认 32780）
- 暴露 Prometheus 指标接口
- 实现反向代理功能，转发 APISIX 的 Prometheus 指标数据
- 支持服务和路由的动态映射配置
- 支持配置文件和环境变量配置
- Docker 容器化部署支持

## 系统架构

### 核心组件

- **Web 服务**：基于 Gin 框架的 HTTP 服务器
- **Prometheus 处理器**：处理 Prometheus 指标的收集和暴露
- **代理处理器**：实现反向代理功能，转发和修改指标数据
- **配置管理**：支持 YAML 配置文件和环境变量

### 目录结构

```
.
├── cmd/                    # 主程序入口
│   └── apisix-metric/     # 主程序包
├── internal/              # 内部包
│   ├── config/           # 配置管理
│   └── service/          # 核心服务实现
├── config/               # 配置文件
├── Dockerfile            # Docker 构建文件
├── Makefile             # 项目管理工具
├── go.mod               # Go 模块定义
└── README.md            # 项目文档
```

## 快速开始

### 前置条件

- Go 1.23.5 或更高版本
- Docker（可选，用于容器化部署）

### 安装

1. 克隆项目：

```bash
git clone https://github.com/fsyyft-go/apisix-metric.git
cd apisix-metric
```

2. 安装依赖：

```bash
go mod download
```

3. 编译项目：

```bash
go build -o apisix-metric ./cmd/apisix-metric
```

### Docker 部署

1. 构建镜像：

```bash
make image
```

2. 运行容器：

```bash
make run
```

## 配置说明

### 配置文件

配置文件位于 `config/config.yaml`，支持以下配置项：

```yaml
server:
  port: 32780                      # 服务器监听端口（默认）

prometheus:
  path: /metrics                    # Prometheus 指标暴露路径

proxy:
  local:
    path: /apisix/prometheus/metrics  # 本地代理路径
  
  service:                          # 服务映射配置
    "service_id": "service_name"
  
  route:                           # 路由映射配置
    "route_id": "route_name"
  
  remote:
    scheme: http                   # 远程服务协议
    host: "remote_host:port"       # 远程服务地址
    path: "/metrics_path"          # 远程服务路径
```

### 环境变量

支持通过环境变量覆盖配置文件中的设置：

- `FSYYFT_APISIX_METRIC_SERVER_PORT`: 服务器端口（默认 32780）
- `FSYYFT_APISIX_METRIC_PROMETHEUS_PATH`: Prometheus 指标路径
- `FSYYFT_APISIX_METRIC_PROXY_LOCAL_PATH`: 本地代理路径
- `FSYYFT_APISIX_METRIC_PROXY_REMOTE_SCHEME`: 远程服务协议
- `FSYYFT_APISIX_METRIC_PROXY_REMOTE_HOST`: 远程服务地址
- `FSYYFT_APISIX_METRIC_PROXY_REMOTE_PATH`: 远程服务路径

## API 文档

### 根路径

- 路径: `/`
- 方法: `GET`
- 描述: 返回服务欢迎信息
- 响应示例:
```json
{
    "message": "Hello, apisix-metric!"
}
```

### Prometheus 指标

- 路径: `/metrics`（可配置）
- 方法: `GET`
- 描述: 暴露 Prometheus 格式的指标数据
- 响应格式: Prometheus 文本格式

### 代理端点

- 路径: `/apisix/prometheus/metrics`（可配置）
- 方法: `GET`
- 描述: 代理 APISIX 的 Prometheus 指标数据
- 响应格式: Prometheus 文本格式

## 开发指南

### 开发环境设置

1. 安装 Go 1.23.5 或更高版本
2. 安装开发工具：
```bash
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### 测试

运行单元测试：

```bash
go test ./...
```

### 代码风格

- 遵循 Go 标准代码风格
- 使用 `goimports` 格式化代码
- 使用 `golangci-lint` 进行代码检查

## 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 维护者

- [@fsyyft-go](https://github.com/fsyyft-go)

## 致谢

感谢所有为这个项目做出贡献的开发者！

## 日志组件

本项目提供了一个灵活的日志组件，支持多种日志后端，可以独立使用。

### 特性

- 支持多种日志后端（标准输出、Logrus）
- 提供统一的日志接口
- 支持结构化日志
- 支持多个日志级别
- 可配置输出目标

### 快速开始

```go
import "github.com/fsyyft-go/apisix-metric/pkg/log"

// 初始化日志系统（使用标准输出）
if err := log.InitLogger(log.LogTypeStd, ""); err != nil {
    panic(err)
}

// 记录日志
log.Info("应用启动")

// 使用结构化字段
log.WithFields(map[string]interface{}{
    "user": "admin",
    "action": "login",
}).Info("用户操作")
```

### 配置示例

```yaml
log:
  # 日志类型：std 或 logrus
  type: logrus
  # 日志输出路径（留空表示标准输出）
  output: logs/app.log
```

更多示例请参考 [example/log](example/log) 目录。
