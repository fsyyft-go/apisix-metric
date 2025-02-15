// Copyright 2024 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package log

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggers(t *testing.T) {
	// 创建临时测试目录
	tmpDir := filepath.Join(os.TempDir(), "apisix-metric-test")
	err := os.MkdirAll(tmpDir, 0755)
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	testCases := []struct {
		name     string
		logType  LogType
		logPath  string
		testFunc func(t *testing.T, logger Logger)
	}{
		{
			name:    "Console Logger",
			logType: LogTypeConsole,
			logPath: "",
			testFunc: func(t *testing.T, logger Logger) {
				logger.Info("测试控制台日志")
				logger.WithField("test", "field").Info("测试带字段的控制台日志")
			},
		},
		{
			name:    "Std Logger File",
			logType: LogTypeStd,
			logPath: filepath.Join(tmpDir, "std.log"),
			testFunc: func(t *testing.T, logger Logger) {
				logger.Info("测试标准库日志文件")
				logger.WithFields(map[string]interface{}{
					"test1": "value1",
					"test2": "value2",
				}).Info("测试带多个字段的标准库日志")
			},
		},
		{
			name:    "Logrus Logger File",
			logType: LogTypeLogrus,
			logPath: filepath.Join(tmpDir, "logrus.log"),
			testFunc: func(t *testing.T, logger Logger) {
				logger.Debug("测试 logrus 调试日志")
				logger.Info("测试 logrus 信息日志")
				logger.Warn("测试 logrus 警告日志")
				logger.WithField("component", "test").Error("测试 logrus 错误日志")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 初始化日志
			err := InitLogger(tc.logType, tc.logPath)
			assert.NoError(t, err)

			// 执行测试
			tc.testFunc(t, GetLogger())

			// 如果有日志文件，验证文件是否创建
			if tc.logPath != "" {
				_, err := os.Stat(tc.logPath)
				assert.NoError(t, err)

				// 读取日志文件内容
				content, err := os.ReadFile(tc.logPath)
				assert.NoError(t, err)
				assert.NotEmpty(t, content)
			}
		})
	}
}

func TestLogLevels(t *testing.T) {
	// 创建临时测试目录
	tmpDir := filepath.Join(os.TempDir(), "apisix-metric-test-levels")
	err := os.MkdirAll(tmpDir, 0755)
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	logPath := filepath.Join(tmpDir, "all-levels.log")

	// 初始化 logrus 日志器
	err = InitLogger(LogTypeLogrus, logPath)
	assert.NoError(t, err)

	logger := GetLogger()

	// 测试所有日志级别
	logger.Debug("Debug message")
	logger.Debugf("Debug message with %s", "format")

	logger.Info("Info message")
	logger.Infof("Info message with %s", "format")

	logger.Warn("Warn message")
	logger.Warnf("Warn message with %s", "format")

	logger.Error("Error message")
	logger.Errorf("Error message with %s", "format")

	// 读取日志文件内容
	content, err := os.ReadFile(logPath)
	assert.NoError(t, err)
	assert.NotEmpty(t, content)
}

func TestWithFieldsAndFormat(t *testing.T) {
	// 创建临时测试目录
	tmpDir := filepath.Join(os.TempDir(), "apisix-metric-test-fields")
	err := os.MkdirAll(tmpDir, 0755)
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	logPath := filepath.Join(tmpDir, "fields.log")

	// 初始化 logrus 日志器
	err = InitLogger(LogTypeLogrus, logPath)
	assert.NoError(t, err)

	logger := GetLogger()

	// 测试字段和格式化
	logger.WithField("single", "field").Info("Single field test")

	logger.WithFields(map[string]interface{}{
		"field1": "value1",
		"field2": 123,
		"field3": true,
	}).Info("Multiple fields test")

	logger.WithField("request_id", "123").
		WithField("user_id", "456").
		Info("Chained fields test")

	// 读取日志文件内容
	content, err := os.ReadFile(logPath)
	assert.NoError(t, err)
	assert.NotEmpty(t, content)
}
