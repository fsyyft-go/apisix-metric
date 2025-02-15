// Copyright 2024 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package log

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// StdLogger 实现了 Logger 接口的标准库日志
type StdLogger struct {
	logger *log.Logger
	fields map[string]interface{}
}

// NewStdLogger 创建一个新的 StdLogger 实例
func NewStdLogger(output string) (Logger, error) {
	var writer *os.File = os.Stdout

	// 如果指定了输出目录
	if output != "" {
		// 确保目录存在
		if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
			return nil, err
		}
		// 打开日志文件
		file, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		writer = file
	}

	return &StdLogger{
		logger: log.New(writer, "", log.LstdFlags),
		fields: make(map[string]interface{}),
	}, nil
}

func (l *StdLogger) formatFields() string {
	if len(l.fields) == 0 {
		return ""
	}
	fields := "["
	for k, v := range l.fields {
		fields += fmt.Sprintf("%s=%v ", k, v)
	}
	return fields[:len(fields)-1] + "]"
}

func (l *StdLogger) log(level string, args ...interface{}) {
	fields := l.formatFields()
	if fields != "" {
		l.logger.Printf("%s %s %v", level, fields, fmt.Sprint(args...))
	} else {
		l.logger.Printf("%s %v", level, fmt.Sprint(args...))
	}
}

func (l *StdLogger) logf(level string, format string, args ...interface{}) {
	fields := l.formatFields()
	if fields != "" {
		l.logger.Printf("%s %s "+format, append([]interface{}{level, fields}, args...)...)
	} else {
		l.logger.Printf("%s "+format, append([]interface{}{level}, args...)...)
	}
}

// Debug implements Logger
func (l *StdLogger) Debug(args ...interface{}) {
	l.log("[DEBUG]", args...)
}

// Debugf implements Logger
func (l *StdLogger) Debugf(format string, args ...interface{}) {
	l.logf("[DEBUG]", format, args...)
}

// Info implements Logger
func (l *StdLogger) Info(args ...interface{}) {
	l.log("[INFO]", args...)
}

// Infof implements Logger
func (l *StdLogger) Infof(format string, args ...interface{}) {
	l.logf("[INFO]", format, args...)
}

// Warn implements Logger
func (l *StdLogger) Warn(args ...interface{}) {
	l.log("[WARN]", args...)
}

// Warnf implements Logger
func (l *StdLogger) Warnf(format string, args ...interface{}) {
	l.logf("[WARN]", format, args...)
}

// Error implements Logger
func (l *StdLogger) Error(args ...interface{}) {
	l.log("[ERROR]", args...)
}

// Errorf implements Logger
func (l *StdLogger) Errorf(format string, args ...interface{}) {
	l.logf("[ERROR]", format, args...)
}

// Fatal implements Logger
func (l *StdLogger) Fatal(args ...interface{}) {
	l.log("[FATAL]", args...)
	os.Exit(1)
}

// Fatalf implements Logger
func (l *StdLogger) Fatalf(format string, args ...interface{}) {
	l.logf("[FATAL]", format, args...)
	os.Exit(1)
}

// WithField implements Logger
func (l *StdLogger) WithField(key string, value interface{}) Logger {
	newFields := make(map[string]interface{})
	for k, v := range l.fields {
		newFields[k] = v
	}
	newFields[key] = value
	return &StdLogger{
		logger: l.logger,
		fields: newFields,
	}
}

// WithFields implements Logger
func (l *StdLogger) WithFields(fields map[string]interface{}) Logger {
	newFields := make(map[string]interface{})
	for k, v := range l.fields {
		newFields[k] = v
	}
	for k, v := range fields {
		newFields[k] = v
	}
	return &StdLogger{
		logger: l.logger,
		fields: newFields,
	}
} 