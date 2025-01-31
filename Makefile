# Makefile for apisix-metric

# 镜像名称，格式为 用户名 / 项目名。
IMAGE_NAME=fsyyft/apisix-metric

# 获取当前日期，格式为 yyMMdd。
DATE=$(shell date +%y%m%d)

# 构建 Docker 镜像并打标签。
# 标签包括日期标签和 latest 标签。
image:
	docker build -t $(IMAGE_NAME):$(DATE) -t $(IMAGE_NAME):latest .

# 运行容器。
# 将容器内的 44444 端口映射到主机的 44444 端口。
# 支持通过环境变量 FSYYFT_APISIX_METRIC_SERVER_PORT 指定端口，默认为 44444。
run:
	docker run -p 48080:8080 -e FSYYFT_APISIX_METRIC_SERVER_PORT=$(or $(FSYYFT_APISIX_METRIC_SERVER_PORT),8080) $(IMAGE_NAME)

# 推送镜像到 Docker Hub。
# 推送所有标签的镜像。
push:
	docker push $(IMAGE_NAME)

# 清理构建缓存。
# 删除本地构建的镜像。
clean:
	docker rmi $(IMAGE_NAME)