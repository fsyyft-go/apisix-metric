# APISIX Metric

APISIX Metric 是一个专门用于优化和转换 APISIX Prometheus 指标展示的服务。它能够将 APISIX 中难以理解的指标标识符（如服务 ID、路由 ID）转换为更易读的名称，同时提供高性能的缓存系统和灵活的日志功能。

## 核心功能

### 指标转换

- **ID 到名称的映射**
  - 支持服务 ID 到服务名称的转换
  - 支持路由 ID 到路由名称的转换
  - 自动从 etcd 获取最新的映射关系

- **数据源支持**
  - 配置文件静态映射
  - etcd 动态数据获取
  - 支持多数据源组合

- **高效的转换处理**
  - 使用反向代理实现透明转换
  - 支持实时数据更新
  - 智能缓存机制减少 etcd 访问

### 工作原理

1. **代理转发**：接收对 `/apisix/prometheus/metrics` 的请求
2. **数据获取**：从多个来源获取 ID 到名称的映射关系
   - 配置文件中的静态映射
   - etcd 中的动态数据
3. **实时转换**：在响应返回前，将指标中的 ID 替换为对应的名称
4. **缓存优化**：使用缓存存储 etcd 数据，减少查询压力

## 快速开始

### 安装

```bash
go get github.com/fsyyft-go/apisix-metric
```

### 配置

创建 `config.yml` 文件：

```yaml
# 服务器配置。
server:
  # 服务器监听端口，默认为 32780。
  port: 32780

# 日志配置。
log:
  # 日志类型，可选值：console, std, logrus，默认为 console。
  type: logrus
  # 日志输出路径，默认输出到控制台。
  output: "logs/apisix-metric.log"
  # 日志级别，可选值：debug, info, warn, error, fatal，默认为 info。
  level: "debug"

# Prometheus 监控配置。
prometheus:
  # Prometheus 指标暴露路径，默认为 /metrics。
  path: /metrics

# 代理配置。
proxy:
  local:
    # 本地代理路径配置。
    path: /apisix/prometheus/metrics

  # 代理服务配置。
  service: {
    "548130462736843557": "apisix.example.com",
    "548144976421192485": "user-service.example.com"
  }
  
  # 代理路由配置。
  route: {
    "548130589790700325": "apisix-route.example.com",
    "548145045006451493": "user-service-route.example.com"
  }

  # 远程代理配置。
  remote:
    # 远程代理协议，默认为 http。
    scheme: http
    # 远程代理主机地址。
    host: "apisix-host:port"
    # 远程代理路径。
    path: "/apisix/prometheus/metrics"

  # etcd 配置。
  etcd:
    # etcd 服务器地址列表。
    endpoints:
      - "http://etcd-host:2379"
    # 连接超时时间（秒）。
    timeout: 30
    # 认证配置（可选）。
    auth:
      username: "root"
      password: "password"
    # 前缀配置。
    prefix: "/apisix"
```

### 运行

```bash
./apisix-metric
```

## 辅助功能

### 缓存系统

为提高性能，项目实现了高效的缓存组件：

- 支持泛型接口，提供类型安全
- 支持 TTL（生存时间）设置
- 提供全局缓存实例
- 线程安全

### 日志系统

为便于问题诊断和监控，实现了灵活的日志组件：

- 支持多种日志后端（标准输出、Logrus）
- 结构化日志记录
- 多级别日志控制
- 可配置输出目标

## API 文档

### 指标转换接口

- 路径: `/apisix/prometheus/metrics`
- 方法: `GET`
- 描述: 返回转换后的 APISIX Prometheus 指标
- 响应格式: Prometheus 文本格式

### Prometheus 原生指标

- 路径: `/metrics`
- 方法: `GET`
- 描述: 暴露服务自身的 Prometheus 指标
- 响应格式: Prometheus 文本格式

### 健康检查

- 路径: `/`
- 方法: `GET`
- 描述: 服务健康检查接口
- 响应示例:
```json
{
    "message": "Hello, apisix-metric!"
}
```

## 环境变量

支持通过环境变量覆盖配置：

- `FSYYFT_APISIX_METRIC_SERVER_PORT`: 服务端口
- `FSYYFT_APISIX_METRIC_PROMETHEUS_PATH`: Prometheus 路径
- `FSYYFT_APISIX_METRIC_PROXY_LOCAL_PATH`: 本地代理路径
- `FSYYFT_APISIX_METRIC_PROXY_REMOTE_SCHEME`: 远程协议
- `FSYYFT_APISIX_METRIC_PROXY_REMOTE_HOST`: 远程地址
- `FSYYFT_APISIX_METRIC_PROXY_REMOTE_PATH`: 远程路径
- `FSYYFT_APISIX_METRIC_LOG_TYPE`: 日志类型
- `FSYYFT_APISIX_METRIC_LOG_OUTPUT`: 日志输出路径
- `FSYYFT_APISIX_METRIC_LOG_LEVEL`: 日志级别

## 开发指南

### 环境要求

- Go 1.23.5 或更高版本
- etcd v3.5.x
- APISIX 3.x

### 构建

```bash
make build
```

### 测试

```bash
make test
```

### Docker 支持

构建镜像：
```bash
make docker-build
```

运行容器：
```bash
docker run -p 32780:32780 fsyyft-go/apisix-metric
```

使用 Docker Compose：
```bash
# 启动所有服务
docker compose up -d

# 查看服务状态
docker compose ps

# 查看服务日志
docker compose logs -f apisix-metric

# 停止所有服务
docker compose down
```

项目提供了完整的 docker-compose.yaml 配置文件，包含了 APISIX Metric、APISIX 和 etcd 服务，可以快速搭建完整的开发和测试环境。

## 维护者

- [@fsyyft-go](https://github.com/fsyyft-go)

## 许可证

MIT License