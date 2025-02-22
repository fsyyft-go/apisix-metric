// Copyright 2024 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package main

import (
	"github.com/fsyyft-go/apisix-metric/internal/config"
	"github.com/fsyyft-go/apisix-metric/internal/service"
	"github.com/fsyyft-go/apisix-metric/pkg/log"
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

	// 初始化日志。
	if err := log.InitLogger(log.LogType(cfg.Log.Type), cfg.Log.Output); err != nil {
		panic(err)
	}

	// 设置日志级别。
	level, err := log.ParseLevel(cfg.Log.Level)
	if err != nil {
		log.Warnf("无效的日志级别 '%s'，使用默认级别 'info'", cfg.Log.Level)
		level = log.InfoLevel
	}
	log.SetLevel(level)

	// 记录应用启动日志。
	log.Info("Starting APISIX Metric service...")
	log.WithFields(map[string]interface{}{
		"port":      cfg.Server.Port,
		"log_type":  cfg.Log.Type,
		"log_path":  cfg.Log.Output,
		"log_level": cfg.Log.Level,
	}).Info("Application configuration loaded")

	// 创建并启动 Web 服务。
	webService := service.NewWebService(cfg)
	if err := webService.Run(); err != nil {
		log.Fatal("Failed to start web service: ", err)
	}
}
