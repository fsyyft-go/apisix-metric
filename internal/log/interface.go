// Copyright 2024 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package log

// Level 定义日志级别
type Level int

const (
	// DebugLevel 调试级别
	DebugLevel Level = iota
	// InfoLevel 信息级别
	InfoLevel
	// WarnLevel 警告级别
	WarnLevel
	// ErrorLevel 错误级别
	ErrorLevel
	// FatalLevel 致命错误级别
	FatalLevel
)

// Logger 定义日志接口
type Logger interface {
	// Debug 输出调试日志
	Debug(args ...interface{})
	// Debugf 输出格式化的调试日志
	Debugf(format string, args ...interface{})

	// Info 输出信息日志
	Info(args ...interface{})
	// Infof 输出格式化的信息日志
	Infof(format string, args ...interface{})

	// Warn 输出警告日志
	Warn(args ...interface{})
	// Warnf 输出格式化的警告日志
	Warnf(format string, args ...interface{})

	// Error 输出错误日志
	Error(args ...interface{})
	// Errorf 输出格式化的错误日志
	Errorf(format string, args ...interface{})

	// Fatal 输出致命错误日志
	Fatal(args ...interface{})
	// Fatalf 输出格式化的致命错误日志
	Fatalf(format string, args ...interface{})

	// WithField 添加单个字段
	WithField(key string, value interface{}) Logger
	// WithFields 添加多个字段
	WithFields(fields map[string]interface{}) Logger
} 