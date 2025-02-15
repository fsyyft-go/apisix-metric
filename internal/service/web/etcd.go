// Copyright 2024 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package web

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/fsyyft-go/apisix-metric/internal/config"
	"github.com/fsyyft-go/apisix-metric/internal/service/etcd"
)

// EtcdHandler 结构体封装了 etcd 处理器。
type EtcdHandler struct {
	// cfg 存储服务配置信息。
	cfg *config.Config
	// service 是 etcd 服务实例。
	service *etcd.Service
}

// NewEtcdHandler 创建并返回一个新的 EtcdHandler 实例。
//
// 参数：
//   - cfg：配置信息，包含 etcd 连接信息等。
//
// 返回：
//   - *EtcdHandler：初始化后的 EtcdHandler 实例。
//   - error：如果初始化失败则返回错误信息。
func NewEtcdHandler(cfg *config.Config) (*EtcdHandler, error) {
	// 创建 etcd 服务实例
	service, err := etcd.NewService(cfg)
	if err != nil {
		return nil, err
	}

	return &EtcdHandler{
		cfg:     cfg,
		service: service,
	}, nil
}

// Register 将 etcd 处理器注册到 Gin 路由中。
func (h *EtcdHandler) Register(r *gin.Engine) {
	// 注册获取 etcd 数据的路由
	r.GET("/etcd/data", h.GetEtcdData)
}

// GetEtcdData 处理获取 etcd 数据的请求。
//
// 参数：
//   - c：Gin 上下文对象，包含请求和响应信息。
func (h *EtcdHandler) GetEtcdData(c *gin.Context) {
	// 从 etcd 获取数据
	data, err := h.service.GetDataWithPrefix(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 返回数据
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

// Close 关闭 etcd 服务连接。
func (h *EtcdHandler) Close() error {
	if h.service != nil {
		return h.service.Close()
	}
	return nil
}
