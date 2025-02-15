// Copyright 2024 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package web

import (
	"github.com/gin-gonic/gin"
)

// RootHandler 结构体封装了根路径处理器。
type RootHandler struct{}

// NewRootHandler 创建并返回一个新的 RootHandler 实例。
func NewRootHandler() *RootHandler {
	return &RootHandler{}
}

// Register 将根路径处理器注册到 Gin 路由中。
func (h *RootHandler) Register(r *gin.Engine) {
	// 处理根路径请求，返回欢迎信息。
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, apisix-metric!",
		})
	})
}
