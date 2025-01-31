package service

import (
	"github.com/gin-gonic/gin"

	"github.com/fsyyft-go/apisix-metric/internal/config"
	"github.com/fsyyft-go/apisix-metric/internal/service/web"
)

// WebService 提供 Web 相关的服务逻辑。
type WebService struct {
	cfg *config.Config
}

// NewWebService 创建一个新的 WebService 实例。
// 参数 cfg 为配置信息，返回初始化后的 WebService 指针。
func NewWebService(cfg *config.Config) *WebService {
	return &WebService{cfg: cfg}
}

// SetupRouter 设置并返回一个配置好的 Gin 实例。
// 返回配置完成的 Gin 引擎实例，包含所有注册的路由。
func (s *WebService) SetupRouter() *gin.Engine {
	// 创建默认的 Gin 实例，包含 Logger 和 Recovery 中间件。
	r := gin.Default()

	// 注册根路由处理程序。
	// 使用 web 包中的 RootHandler 处理根路径请求。
	rootHandler := web.NewRootHandler()
	rootHandler.Register(r)

	// 注册 Prometheus 指标路由。
	// 使用 web 包中的 PrometheusHandler 处理指标收集。
	promHandler := web.NewPrometheusHandler(s.cfg.Prometheus.Path)
	promHandler.Register(r)

	// 注册代理路由处理程序
	proxyHandler := web.NewProxyHandler(s.cfg)
	proxyHandler.Register(r)

	return r
}

// Run 启动 HTTP 服务并开始监听请求。
// 返回 error 类型，如果启动失败则返回错误信息。
func (s *WebService) Run() error {
	// 初始化路由配置。
	r := s.SetupRouter()
	// 启动服务，监听端口来自配置文件。
	return r.Run(":" + s.cfg.Server.Port)
}
