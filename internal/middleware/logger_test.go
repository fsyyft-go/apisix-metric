// Copyright 2024 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

// Package middleware 提供了 HTTP 请求处理的中间件功能。
package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/fsyyft-go/apisix-metric/internal/log"
)

// TestRequestLogger 测试请求日志中间件的功能。
// 该测试用例验证以下功能：
// 1. 中间件是否正确记录请求信息
// 2. 日志文件是否正确创建和写入
// 3. 日志内容是否包含所有必要字段
// 4. 请求处理是否正常完成
func TestRequestLogger(t *testing.T) {
	// 创建临时测试目录，用于存放测试日志文件。
	// 使用系统临时目录以确保在不同环境下都能正常工作。
	tmpDir := filepath.Join(os.TempDir(), "apisix-metric-middleware-test")
	err := os.MkdirAll(tmpDir, 0755)
	assert.NoError(t, err)
	// 测试完成后清理临时目录。
	defer os.RemoveAll(tmpDir)

	// 设置测试日志文件路径。
	logPath := filepath.Join(tmpDir, "request.log")

	// 初始化日志系统，使用 logrus 作为日志记录器。
	err = log.InitLogger(log.LogTypeLogrus, logPath)
	assert.NoError(t, err)

	// 创建测试用的 gin 引擎，使用测试模式避免额外的日志输出。
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// 将请求日志中间件添加到路由中。
	r.Use(RequestLogger())

	// 添加一个测试路由，返回简单的 JSON 响应。
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "test",
		})
	})

	// 创建一个测试请求记录器。
	w := httptest.NewRecorder()
	// 创建一个 GET 请求，添加必要的请求头。
	req, _ := http.NewRequest("GET", "/test", nil)
	// 设置用户代理信息。
	req.Header.Set("User-Agent", "test-agent")
	// 设置来源页面信息。
	req.Header.Set("Referer", "test-referer")
	// 设置请求跟踪 ID。
	req.Header.Set("X-Request-ID", "test-request-id")

	// 执行请求。
	r.ServeHTTP(w, req)

	// 验证响应状态码是否为 200。
	assert.Equal(t, http.StatusOK, w.Code)

	// 读取并验证日志文件内容。
	content, err := os.ReadFile(logPath)
	assert.NoError(t, err)
	// 确保日志文件不为空。
	assert.NotEmpty(t, content)

	// 将日志内容转换为字符串，用于后续验证。
	logContent := string(content)
	// 验证日志中是否包含所有必要的信息。
	assert.Contains(t, logContent, "HTTP Request")
	assert.Contains(t, logContent, "GET")
	assert.Contains(t, logContent, "/test")
	assert.Contains(t, logContent, "test-agent")
	assert.Contains(t, logContent, "test-referer")
	assert.Contains(t, logContent, "test-request-id")
}
