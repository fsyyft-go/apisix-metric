# 使用官方 Go 镜像作为构建环境。
FROM golang:1.24.0-alpine AS builder

# 设置工作目录。
WORKDIR /app

# 复制 go.mod 和 go.sum 文件。
COPY go.mod go.sum ./

# 设置 GOPROXY 并下载依赖。
RUN go env -w GOPROXY=https://goproxy.cn,direct && \
    go mod download

# 复制项目文件。
COPY . .

# 构建应用。
RUN CGO_ENABLED=0 GOOS=linux go build -o apisix-metric ./cmd/apisix-metric

# 使用精简的 alpine 镜像作为运行时环境。
FROM alpine:3.19

# 设置工作目录。
WORKDIR /app

# 从构建阶段复制可执行文件。
COPY --from=builder /app/apisix-metric .

# 复制配置文件。
COPY config/config.yaml ./config/config.yaml

# 创建日志目录。
RUN mkdir -p /app/logs && \
    chmod 755 /app/logs

# 设置默认的环境变量。
ENV FSYYFT_APISIX_METRIC_LOG_TYPE=logrus \
    FSYYFT_APISIX_METRIC_LOG_OUTPUT=/app/logs/app.log

# 暴露端口。
EXPOSE 32780

# 声明日志卷。
VOLUME ["/app/logs"]

# 启动应用。
CMD ["./apisix-metric"]