// Copyright 2024 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package web

import (
	"io"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/fsyyft-go/apisix-metric/internal/config"
)

// ProxyHandler 结构体封装了代理处理器。
type ProxyHandler struct {
	cfg *config.Config
}

// NewProxyHandler 创建并返回一个新的 ProxyHandler 实例。
// 参数 path 指定代理路径。
func NewProxyHandler(cfg *config.Config) *ProxyHandler {
	return &ProxyHandler{
		cfg: cfg,
	}
}

// Register 将代理处理器注册到 Gin 路由中。
func (h *ProxyHandler) Register(r *gin.Engine) {
	// 默认实现返回 200 OK
	r.GET(h.cfg.Proxy.Local.Path, h.ReverseProxyHandler)
}

// ReverseProxyHandler 处理反向代理请求，转发到 APISIX 的 Prometheus metrics 接口。
// c: Gin 上下文对象，包含请求和响应信息。
func (h *ProxyHandler) ReverseProxyHandler(c *gin.Context) {
	config := h.cfg.Proxy

	// 创建反向代理实例
	proxy := &httputil.ReverseProxy{
		// Director 函数用于修改请求
		Director: func(req *http.Request) {
			req.URL.Scheme = config.Remote.Scheme // 设置协议
			req.URL.Host = config.Remote.Host     // 设置目标地址
			req.URL.Path = config.Remote.Path     // 设置路径
			req.Host = config.Remote.Host         // 设置 Host 头
		},
		// ModifyResponse 函数用于修改响应
		ModifyResponse: func(resp *http.Response) error {
			// 仅处理 200 状态码的响应
			if resp.StatusCode == http.StatusOK {
				// 读取响应体
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return err
				}
				resp.Body.Close()

				// 将响应体转换为字符串
				content := string(body)

				// 使用 ServiceList 进行内容替换
				for oldStr, newStr := range config.Service {
					content = strings.ReplaceAll(content, oldStr, newStr)
				}

				// 使用 RouteList 进行内容替换
				for oldStr, newStr := range config.Route {
					content = strings.ReplaceAll(content, oldStr, newStr)
				}

				// 更新响应体
				resp.Body = io.NopCloser(strings.NewReader(content))
				resp.ContentLength = int64(len(content))                      // 更新内容长度
				resp.Header.Set("Content-Length", strconv.Itoa(len(content))) // 更新 Content-Length 头
			}
			return nil
		},
	}

	// 处理请求
	proxy.ServeHTTP(c.Writer, c.Request)
}
