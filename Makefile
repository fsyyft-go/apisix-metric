# Makefile for apisix-metric

# 镜像名称，格式为 用户名 / 项目名。
IMAGE_NAME=fsyyft/apisix-metric

# 获取当前日期，格式为 yyMMdd。
DATE=$(shell date +%y%m%d)

# 日志目录。
LOG_DIR=logs

# 运行所有测试。
# 使用 -v 标志显示详细的测试输出。
# 使用 -race 标志检测数据竞争。
test:
	mkdir -p $(LOG_DIR)
	go test -v -race ./...

# 运行性能测试。
# 使用 -bench 标志运行基准测试。
# 使用 -benchmem 标志显示内存分配统计。
# 使用 -count 标志指定运行次数。
# 使用 -cpu 标志指定 CPU 核心数。
bench:
	go test -bench=. -benchmem -count=5 -cpu=1,2,4,8 ./pkg/cache/...

# 运行测试并生成覆盖率报告。
# 生成 HTML 格式的覆盖率报告。
# 报告将保存在 coverage 目录下。
coverage:
	mkdir -p coverage
	mkdir -p $(LOG_DIR)
	go test -v -race -coverprofile=coverage/coverage.out ./...
	go tool cover -html=coverage/coverage.out -o coverage/coverage.html

# 构建 Docker 镜像并打标签。
# 标签包括日期标签和 latest 标签。
image:
	docker build -t $(IMAGE_NAME):$(DATE) -t $(IMAGE_NAME):latest .

# 运行容器。
# 将容器内的 32780 端口映射到主机的 32780 端口。
# 将容器内的 /app/logs 目录挂载到主机的 ./logs 目录。
# 支持通过环境变量配置：
# - FSYYFT_APISIX_METRIC_SERVER_PORT：服务端口，默认为 32780
# - FSYYFT_APISIX_METRIC_LOG_TYPE：日志类型，默认为 logrus
# - FSYYFT_APISIX_METRIC_LOG_OUTPUT：日志输出路径，默认为 /app/logs/app.log
run:
	mkdir -p $(LOG_DIR)
	docker run \
		-p 32780:32780 \
		-v $(PWD)/$(LOG_DIR):/app/logs \
		-e FSYYFT_APISIX_METRIC_SERVER_PORT=$(or $(FSYYFT_APISIX_METRIC_SERVER_PORT),32780) \
		-e FSYYFT_APISIX_METRIC_LOG_TYPE=$(or $(FSYYFT_APISIX_METRIC_LOG_TYPE),logrus) \
		-e FSYYFT_APISIX_METRIC_LOG_OUTPUT=$(or $(FSYYFT_APISIX_METRIC_LOG_OUTPUT),/app/logs/app.log) \
		$(IMAGE_NAME)

# 推送镜像到 Docker Hub。
# 推送所有标签的镜像。
push:
	docker push $(IMAGE_NAME)

# 清理构建缓存和日志文件。
# 删除本地构建的镜像、覆盖率报告和日志文件。
clean:
	docker rmi $(IMAGE_NAME)
	rm -rf coverage
	rm -rf $(LOG_DIR)