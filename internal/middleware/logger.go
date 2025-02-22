// Copyright 2024 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

// Package middleware 提供了 HTTP 请求处理的中间件功能。
package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/fsyyft-go/apisix-metric/pkg/log"
)

// RequestLogger 创建一个用于记录 HTTP 请求信息的中间件。
// 该中间件会记录请求的详细信息，包括：
// - 请求方法（GET、POST等）
// - 请求路径
// - 响应状态码
// - 客户端IP
// - 请求处理时间
// - User-Agent
// - Referer
// - 请求ID（X-Request-ID）
// - 错误信息
//
// 所有信息都会通过全局日志记录器以结构化的方式记录。
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间，用于计算请求处理时长。
		start := time.Now()

		// 调用下一个处理器，等待请求处理完成。
		c.Next()

		// 记录请求结束时间。
		end := time.Now()
		// 计算请求处理时长。
		latency := end.Sub(start)

		// 从请求头中获取用户代理信息。
		userAgent := c.Request.UserAgent()
		// 从请求头中获取来源页面信息。
		referer := c.Request.Referer()

		// 使用结构化字段记录请求信息。
		// WithFields 方法允许我们一次性添加多个字段。
		log.WithFields(map[string]interface{}{
			// 记录 HTTP 请求方法（GET、POST等）。
			"method": c.Request.Method,
			// 记录请求的完整路径。
			"path": c.Request.URL.Path,
			// 记录响应状态码。
			"status": c.Writer.Status(),
			// 记录客户端 IP 地址。
			"ip": c.ClientIP(),
			// 记录请求处理时长。
			"latency": latency.String(),
			// 记录客户端类型。
			"user_agent": userAgent,
			// 记录请求来源页面。
			"referer": referer,
			// 记录请求跟踪 ID，用于请求追踪。
			"request_id": c.GetHeader("X-Request-ID"),
			// 记录处理过程中的错误信息。
			"error": c.Errors.String(),
		}).Info("HTTP Request")
	}
}
