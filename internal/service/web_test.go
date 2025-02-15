// Copyright 2024 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package service

import (
	"net/http"
	"testing"

	"github.com/fsyyft-go/apisix-metric/internal/config"

	"github.com/stretchr/testify/assert"
)

// TestRun 测试 Run 方法是否正确启动 HTTP 服务。
func TestRun(t *testing.T) {
	// 创建一个模拟的配置对象。
	cfg := &config.Config{
		Server: config.Server{
			Port: "8080",
		},
		Prometheus: config.Prometheus{
			Path: "/metrics",
		},
		Proxy: config.Proxy{
			Local: config.Local{
				Path: "/apisix/prometheus/metrics",
			},
		},
	}

	// 创建 WebService 实例。
	webService := NewWebService(cfg)

	// 启动 HTTP 服务。
	go func() {
		err := webService.Run()
		assert.NoError(t, err)
	}()

	// 创建一个测试请求。
	req, err := http.NewRequest("GET", "http://localhost:8080/", nil)
	assert.NoError(t, err)

	// 创建一个响应记录器。
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// 验证响应状态码。
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// 创建一个测试请求。
	req, err = http.NewRequest("GET", "http://localhost:8080/metrics", nil)
	assert.NoError(t, err)

	// 创建一个响应记录器。
	resp, err = client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// 验证响应状态码。
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
