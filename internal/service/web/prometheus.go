package web

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// PrometheusHandler 结构体封装了 Prometheus 指标处理器。
type PrometheusHandler struct {
	// Path 表示 Prometheus 指标的暴露路径。
	Path string
}

// NewPrometheusHandler 创建并返回一个新的 PrometheusHandler 实例。
// 参数 path 指定 Prometheus 指标的暴露路径。
func NewPrometheusHandler(path string) *PrometheusHandler {
	return &PrometheusHandler{
		Path: path,
	}
}

// Register 将 Prometheus 处理器注册到 Gin 路由中。
func (h *PrometheusHandler) Register(r *gin.Engine) {
	// 使用 promhttp.Handler() 处理 Prometheus 指标请求。
	r.GET(h.Path, gin.WrapH(promhttp.Handler()))
}
