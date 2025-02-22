// Copyright 2024 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package service

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/fsyyft-go/apisix-metric/internal/config"

	"github.com/stretchr/testify/assert"
)

// waitForServer 等待服务启动。
func waitForServer(url string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			resp, err := http.Get(url)
			if err == nil {
				resp.Body.Close()
				return nil
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// TestRun 测试 Run 方法是否正确启动 HTTP 服务。
func TestRun(t *testing.T) {
	// 创建一个模拟的配置对象。
	cfg := &config.Config{
		Server: config.Server{
			Port: "44480",
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

	// 创建一个用于停止服务的 context。
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动 HTTP 服务。
	go func() {
		if err := webService.Run(ctx); err != nil && err != context.Canceled {
			t.Errorf("服务启动失败：%v", err)
		}
	}()

	// 等待服务启动。
	err := waitForServer("http://localhost:44480/")
	assert.NoError(t, err, "服务启动超时")

	// 创建 HTTP 客户端。
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// 测试根路径。
	t.Run("Test Root Path", func(t *testing.T) {
		resp, err := client.Get("http://localhost:44480/")
		assert.NoError(t, err)
		if resp != nil {
			defer resp.Body.Close()
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}
	})

	// 测试 metrics 路径。
	t.Run("Test Metrics Path", func(t *testing.T) {
		resp, err := client.Get("http://localhost:44480/metrics")
		assert.NoError(t, err)
		if resp != nil {
			defer resp.Body.Close()
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}
	})

	// 停止服务。
	cancel()
}
