// Copyright 2024 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package log

import (
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// LogrusLogger 实现了 Logger 接口的 logrus 日志
type LogrusLogger struct {
	logger *logrus.Entry
}

// NewLogrusLogger 创建一个新的 LogrusLogger 实例
func NewLogrusLogger(output string) (Logger, error) {
	log := logrus.New()
	
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
		// 同时输出到文件和控制台
		log.SetOutput(io.MultiWriter(os.Stdout, file))
	}

	// 设置日志格式
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	return &LogrusLogger{
		logger: logrus.NewEntry(log),
	}, nil
}

// Debug implements Logger
func (l *LogrusLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

// Debugf implements Logger
func (l *LogrusLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

// Info implements Logger
func (l *LogrusLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

// Infof implements Logger
func (l *LogrusLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

// Warn implements Logger
func (l *LogrusLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

// Warnf implements Logger
func (l *LogrusLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

// Error implements Logger
func (l *LogrusLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

// Errorf implements Logger
func (l *LogrusLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

// Fatal implements Logger
func (l *LogrusLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

// Fatalf implements Logger
func (l *LogrusLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

// WithField implements Logger
func (l *LogrusLogger) WithField(key string, value interface{}) Logger {
	return &LogrusLogger{
		logger: l.logger.WithField(key, value),
	}
}

// WithFields implements Logger
func (l *LogrusLogger) WithFields(fields map[string]interface{}) Logger {
	return &LogrusLogger{
		logger: l.logger.WithFields(fields),
	}
} 