// Copyright 2024 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

// Package service 提供了 Web 服务相关的核心功能实现。
package service

import (
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/fsyyft-go/apisix-metric/internal/config"
	"github.com/fsyyft-go/apisix-metric/internal/middleware"
	"github.com/fsyyft-go/apisix-metric/internal/service/web"
)

// WebService 提供 Web 相关的服务逻辑。
// 该结构体负责：
// 1. 初始化和配置 HTTP 服务器
// 2. 注册所有路由和处理器
// 3. 集成中间件
// 4. 启动和管理 HTTP 服务
type WebService struct {
	// cfg 存储服务配置信息。
	cfg *config.Config
	// engine 存储 Gin 引擎实例。
	engine *gin.Engine
	// once 确保 SetupRouter 只被调用一次。
	once sync.Once
}

// NewWebService 创建一个新的 WebService 实例。
//
// 参数：
//   - cfg：配置信息，包含服务器端口、日志配置等。
//
// 返回：
//   - *WebService：初始化后的 WebService 实例。
func NewWebService(cfg *config.Config) *WebService {
	return &WebService{
		cfg: cfg,
	}
}

// setupEngine 初始化并配置 Gin 引擎实例。
// 该方法是 SetupRouter 的内部实现，由 sync.Once 保证只被调用一次。
//
// 返回：
//   - *gin.Engine：配置完成的 Gin 引擎实例。
func (s *WebService) setupEngine() *gin.Engine {
	// 创建新的 Gin 引擎实例。
	s.engine = gin.New()

	// 添加 Recovery 中间件，用于处理 panic。
	s.engine.Use(gin.Recovery())

	// 添加请求日志中间件，用于记录所有 HTTP 请求的详细信息。
	s.engine.Use(middleware.RequestLogger())

	// 注册根路由处理程序。
	// 使用 web 包中的 RootHandler 处理根路径请求。
	rootHandler := web.NewRootHandler()
	rootHandler.Register(s.engine)

	// 注册 Prometheus 指标路由。
	// 使用 web 包中的 PrometheusHandler 处理指标收集。
	promHandler := web.NewPrometheusHandler(s.cfg.Prometheus.Path)
	promHandler.Register(s.engine)

	// 注册代理路由处理程序。
	// 用于转发和处理 APISIX 的 Prometheus 指标数据。
	proxyHandler := web.NewProxyHandler(s.cfg)
	proxyHandler.Register(s.engine)

	return s.engine
}

// SetupRouter 设置并返回一个配置好的 Gin 实例。
// 该方法完成以下工作：
// 1. 使用 sync.Once 确保只初始化一次
// 2. 添加必要的中间件
// 3. 注册所有路由处理器
// 4. 设置必要的路由参数
//
// 返回：
//   - *gin.Engine：配置完成的 Gin 引擎实例，包含所有注册的路由。
func (s *WebService) SetupRouter() *gin.Engine {
	s.once.Do(func() {
		s.setupEngine()
	})
	return s.engine
}

// Run 启动 HTTP 服务并开始监听请求。
// 该方法会阻塞直到服务器关闭或发生错误。
//
// 返回：
//   - error：如果服务启动失败则返回错误信息。
func (s *WebService) Run() error {
	// 确保路由已经设置。
	s.SetupRouter()
	// 启动服务，监听端口来自配置文件。
	return s.engine.Run(":" + s.cfg.Server.Port)
}

// Engine 返回当前的 Gin 引擎实例。
// 这个方法主要用于测试目的，允许测试代码直接访问引擎实例。
//
// 返回：
//   - *gin.Engine：当前的 Gin 引擎实例。
func (s *WebService) Engine() *gin.Engine {
	return s.SetupRouter()
}
