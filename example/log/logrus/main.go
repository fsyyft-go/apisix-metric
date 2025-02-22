// Copyright 2024 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

// Package main 提供了使用 Logrus 作为日志后端的高级示例，展示了文件输出和结构化日志等特性。
package main

import (
	"os"
	"path/filepath"

	"github.com/fsyyft-go/apisix-metric/pkg/log"
)

func main() {
	// 创建日志目录。
	// 使用 0755 权限确保目录可读可执行，且所有者可写。
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(err)
	}

	// 初始化 Logrus 日志系统。
	// 配置日志输出到文件，支持结构化日志记录。
	logPath := filepath.Join(logDir, "app.log")
	if err := log.InitLogger(log.LogTypeLogrus, logPath); err != nil {
		panic(err)
	}

	// 演示不同级别的日志记录。
	// 每个级别的日志都会包含时间戳、日志级别等元数据。
	log.Debug("Logrus 调试信息。")
	log.Info("Logrus 普通信息。")
	log.Warn("Logrus 警告信息。")
	log.Error("Logrus 错误信息。")

	// 使用结构化字段记录系统信息。
	// Logrus 会以 JSON 格式输出这些字段，便于日志解析和分析。
	log.WithFields(map[string]interface{}{
		"os":      "linux",
		"version": "1.0.0",
		"pid":     os.Getpid(),
	}).Info("系统启动信息。")

	// 演示错误处理和日志记录的集成。
	// 在实际应用中，应该始终记录错误的详细信息。
	err := someOperation()
	if err != nil {
		log.WithFields(map[string]interface{}{
			"error": err.Error(),
			"func":  "someOperation",
		}).Error("操作失败。")
	}
}

// someOperation 是一个示例函数，用于演示错误处理和日志记录。
// 在实际应用中，这里可能是数据库操作、网络请求等。
func someOperation() error {
	return nil
}
