package web

import (
	"github.com/gin-gonic/gin"
)

// ProxyHandler 结构体封装了代理处理器。
type ProxyHandler struct {
	// Path 表示代理路径。
	Path string
}

// NewProxyHandler 创建并返回一个新的 ProxyHandler 实例。
// 参数 path 指定代理路径。
func NewProxyHandler(path string) *ProxyHandler {
	return &ProxyHandler{
		Path: path,
	}
}

// Register 将代理处理器注册到 Gin 路由中。
func (h *ProxyHandler) Register(r *gin.Engine) {
	// 默认实现返回 200 OK
	r.GET(h.Path, func(c *gin.Context) {
		c.String(200, "Proxy handler response")
	})
}
