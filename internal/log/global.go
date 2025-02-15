// Copyright 2024 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package log

import "fmt"

// LogType 定义日志类型
type LogType string

const (
	// LogTypeConsole 控制台日志
	LogTypeConsole LogType = "console"
	// LogTypeStd 标准库日志
	LogTypeStd LogType = "std"
	// LogTypeLogrus logrus日志
	LogTypeLogrus LogType = "logrus"
)

var (
	// globalLogger 全局日志实例
	globalLogger Logger
)

// InitLogger 初始化全局日志实例
func InitLogger(logType LogType, output string) error {
	var err error
	switch logType {
	case LogTypeConsole:
		globalLogger, err = NewStdLogger("")
	case LogTypeStd:
		globalLogger, err = NewStdLogger(output)
	case LogTypeLogrus:
		globalLogger, err = NewLogrusLogger(output)
	default:
		return fmt.Errorf("unsupported log type: %s", logType)
	}
	return err
}

// GetLogger 获取全局日志实例
func GetLogger() Logger {
	if globalLogger == nil {
		// 如果未初始化，使用控制台日志作为默认值
		globalLogger, _ = NewStdLogger("")
	}
	return globalLogger
}

// Debug 全局 Debug 日志
func Debug(args ...interface{}) {
	GetLogger().Debug(args...)
}

// Debugf 全局格式化 Debug 日志
func Debugf(format string, args ...interface{}) {
	GetLogger().Debugf(format, args...)
}

// Info 全局 Info 日志
func Info(args ...interface{}) {
	GetLogger().Info(args...)
}

// Infof 全局格式化 Info 日志
func Infof(format string, args ...interface{}) {
	GetLogger().Infof(format, args...)
}

// Warn 全局 Warn 日志
func Warn(args ...interface{}) {
	GetLogger().Warn(args...)
}

// Warnf 全局格式化 Warn 日志
func Warnf(format string, args ...interface{}) {
	GetLogger().Warnf(format, args...)
}

// Error 全局 Error 日志
func Error(args ...interface{}) {
	GetLogger().Error(args...)
}

// Errorf 全局格式化 Error 日志
func Errorf(format string, args ...interface{}) {
	GetLogger().Errorf(format, args...)
}

// Fatal 全局 Fatal 日志
func Fatal(args ...interface{}) {
	GetLogger().Fatal(args...)
}

// Fatalf 全局格式化 Fatal 日志
func Fatalf(format string, args ...interface{}) {
	GetLogger().Fatalf(format, args...)
}

// WithField 全局添加单个字段
func WithField(key string, value interface{}) Logger {
	return GetLogger().WithField(key, value)
}

// WithFields 全局添加多个字段
func WithFields(fields map[string]interface{}) Logger {
	return GetLogger().WithFields(fields)
} 