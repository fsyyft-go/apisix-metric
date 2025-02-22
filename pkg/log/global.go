// Copyright 2024 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

// Package log 提供了全局日志功能，支持多种日志后端的统一管理。
package log

import "fmt"

// LogType 定义了支持的日志类型，用于在初始化时选择具体的日志实现。
type LogType string

const (
	// LogTypeConsole 表示控制台日志类型。
	// 这种类型的日志会直接输出到标准输出，适合开发调试使用。
	LogTypeConsole LogType = "console"

	// LogTypeStd 表示标准库日志类型。
	// 使用 Go 标准库的 log 包实现，提供基本的日志功能。
	LogTypeStd LogType = "std"

	// LogTypeLogrus 表示 Logrus 日志类型。
	// 使用 Logrus 库实现，提供丰富的日志功能，包括结构化日志、多种输出格式等。
	LogTypeLogrus LogType = "logrus"
)

var (
	// globalLogger 是全局日志实例，所有的全局日志方法都会使用这个实例。
	// 在使用前必须通过 InitLogger 初始化，否则会使用默认的标准输出日志器。
	globalLogger Logger
)

// InitLogger 初始化全局日志实例。
// 参数 logType 指定要使用的日志类型，output 指定日志输出路径。
// 当 output 为空字符串时，日志将输出到标准输出。
// 返回初始化过程中可能发生的错误。
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

// GetLogger 获取全局日志实例。
// 如果全局日志实例未初始化，将创建一个默认的标准输出日志器。
// 返回当前的全局日志实例。
func GetLogger() Logger {
	if globalLogger == nil {
		// 如果未初始化，使用标准输出日志器作为默认值。
		globalLogger, _ = NewStdLogger("")
	}
	return globalLogger
}

// Debug 使用全局日志实例记录调试级别的日志。
// 参数 args 支持任意类型的值，这些值会被转换为字符串并连接。
func Debug(args ...interface{}) {
	GetLogger().Debug(args...)
}

// Debugf 使用全局日志实例记录格式化的调试级别日志。
// 参数 format 是格式化字符串，args 是对应的参数。
func Debugf(format string, args ...interface{}) {
	GetLogger().Debugf(format, args...)
}

// Info 使用全局日志实例记录信息级别的日志。
// 参数 args 支持任意类型的值，这些值会被转换为字符串并连接。
func Info(args ...interface{}) {
	GetLogger().Info(args...)
}

// Infof 使用全局日志实例记录格式化的信息级别日志。
// 参数 format 是格式化字符串，args 是对应的参数。
func Infof(format string, args ...interface{}) {
	GetLogger().Infof(format, args...)
}

// Warn 使用全局日志实例记录警告级别的日志。
// 参数 args 支持任意类型的值，这些值会被转换为字符串并连接。
func Warn(args ...interface{}) {
	GetLogger().Warn(args...)
}

// Warnf 使用全局日志实例记录格式化的警告级别日志。
// 参数 format 是格式化字符串，args 是对应的参数。
func Warnf(format string, args ...interface{}) {
	GetLogger().Warnf(format, args...)
}

// Error 使用全局日志实例记录错误级别的日志。
// 参数 args 支持任意类型的值，这些值会被转换为字符串并连接。
func Error(args ...interface{}) {
	GetLogger().Error(args...)
}

// Errorf 使用全局日志实例记录格式化的错误级别日志。
// 参数 format 是格式化字符串，args 是对应的参数。
func Errorf(format string, args ...interface{}) {
	GetLogger().Errorf(format, args...)
}

// Fatal 使用全局日志实例记录致命错误级别的日志。
// 参数 args 支持任意类型的值，这些值会被转换为字符串并连接。
// 记录日志后会导致程序以状态码 1 退出。
func Fatal(args ...interface{}) {
	GetLogger().Fatal(args...)
}

// Fatalf 使用全局日志实例记录格式化的致命错误级别日志。
// 参数 format 是格式化字符串，args 是对应的参数。
// 记录日志后会导致程序以状态码 1 退出。
func Fatalf(format string, args ...interface{}) {
	GetLogger().Fatalf(format, args...)
}

// WithField 使用全局日志实例添加一个字段到日志上下文。
// 参数 key 是字段名，value 是字段值。
// 返回一个新的 Logger 实例，原实例不会被修改。
func WithField(key string, value interface{}) Logger {
	return GetLogger().WithField(key, value)
}

// WithFields 使用全局日志实例添加多个字段到日志上下文。
// 参数 fields 是要添加的字段映射。
// 返回一个新的 Logger 实例，原实例不会被修改。
func WithFields(fields map[string]interface{}) Logger {
	return GetLogger().WithFields(fields)
}
