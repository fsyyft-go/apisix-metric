// Copyright 2024 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package main

import (
	"github.com/fsyyft-go/apisix-metric/pkg/log"
)

func main() {
	// 初始化日志系统，使用标准输出
	if err := log.InitLogger(log.LogTypeStd, ""); err != nil {
		panic(err)
	}

	// 基本日志记录
	log.Debug("这是一条调试日志")
	log.Info("这是一条信息日志")
	log.Warn("这是一条警告日志")
	log.Error("这是一条错误日志")

	// 使用结构化字段
	log.WithField("user", "admin").Info("用户登录")

	// 使用多个字段
	log.WithFields(map[string]interface{}{
		"user":   "admin",
		"action": "login",
		"time":   "2024-02-23 10:00:00",
	}).Info("用户操作记录")
}
