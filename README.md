# APISIX Metric

APISIX Metric 是一个用于收集和展示 APISIX 指标的服务。它提供了高性能的缓存组件、灵活的日志系统和可靠的指标收集功能。

## 主要特性

- 高性能的缓存系统
  - 支持泛型接口，提供类型安全
  - 支持 TTL（生存时间）设置
  - 提供全局缓存实例
  - 线程安全
- 灵活的日志系统
  - 支持多种日志后端
  - 结构化日志
  - 可配置的输出格式
- 可靠的指标收集
  - 支持 Prometheus 格式
  - 支持自定义指标
  - 实时数据更新

## 性能指标

### 缓存组件性能

在 Apple M3 处理器上的基准测试结果（单核心）：

- 基本操作
  - Set: ~900ns/op, 215B/op, 5 allocs/op
  - Get: ~100ns/op, 0B/op, 0 allocs/op
  - SetWithTTL: ~1000ns/op, 215B/op, 5 allocs/op
  - GetWithTTL: ~150ns/op, 0B/op, 0 allocs/op

- 类型安全操作（使用泛型）
  - Set: ~950ns/op, 215B/op, 5 allocs/op
  - Get: ~120ns/op, 0B/op, 0 allocs/op
  - SetWithTTL: ~1050ns/op, 215B/op, 5 allocs/op
  - GetWithTTL: ~170ns/op, 0B/op, 0 allocs/op

## 快速开始

### 安装

```bash
go get github.com/fsyyft-go/apisix-metric
```

### 使用缓存组件

1. 基本用法

```go
// 创建缓存实例
cache, err := cache.NewCache(cache.DefaultConfig())
if err != nil {
    panic(err)
}
defer cache.Close()

// 基本操作
cache.Set("key", "value")
value, exists := cache.Get("key")

// 带 TTL 的操作
cache.SetWithTTL("temp", "value", time.Hour)
value, exists, ttl := cache.GetWithTTL("temp")
```

2. 类型安全的缓存

```go
// 创建类型安全的缓存
strCache := cache.NewTypedCache[string](baseCache)
intCache := cache.NewTypedCache[int](baseCache)

// 类型安全的操作
strCache.Set("name", "Alice")
name, exists := strCache.Get("name")

intCache.Set("age", 25)
age, exists := intCache.Get("age")
```

3. 全局缓存

```go
// 初始化全局缓存
if err := cache.InitCache(cache.DefaultConfig()); err != nil {
    panic(err)
}
defer cache.Close()

// 使用全局缓存
cache.Set("global", "value")
value, exists := cache.Get("global")
```

更多示例请参考 [example/cache](example/cache) 目录。

## 配置说明

### 缓存配置

```go
type Config struct {
    NumCounters int64  // 计数器数量，建议为预期键数量的 10 倍
    MaxCost     int64  // 最大成本（可理解为最大条目数）
    BufferItems int64  // 写入缓冲区大小
}
```

默认配置：
- NumCounters: 1000 万
- MaxCost: 1GB
- BufferItems: 64

## 开发

### 运行测试

```bash
# 运行所有测试
make test

# 运行性能测试
make bench

# 生成覆盖率报告
make coverage
```

### 构建

```bash
# 构建 Docker 镜像
make image

# 运行容器
make run
```

## 许可证

MIT License

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
