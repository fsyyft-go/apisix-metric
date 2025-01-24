package main

import (
	"github.com/gin-gonic/gin"

	"github.com/fsyyft-go/apisix-metric/handlers"
	"github.com/fsyyft-go/apisix-metric/internal/config"
)

// 主程序入口。
// 初始化 Gin Web 框架并启动 HTTP 服务。
func main() {
	// 加载配置文件，配置文件路径为 config/config.yaml。
	// 如果加载失败，程序将直接 panic。
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		panic(err)
	}

	// 创建默认的 Gin 实例，包含 Logger 和 Recovery 中间件。
	r := gin.Default()

	// 注册根路由处理函数。
	// 当访问 / 路径时，返回 JSON 格式的欢迎信息。
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, apisix-metric!",
		})
	})

	// 创建并注册 Prometheus 指标处理器。
	promHandler := handlers.NewPrometheusHandler(cfg.Prometheus.Path)
	promHandler.Register(r)

	// 启动 HTTP 服务，监听端口来自配置文件或环境变量。
	r.Run(":" + cfg.Server.Port)
}
