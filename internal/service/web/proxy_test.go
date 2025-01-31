package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewProxyHandler(t *testing.T) {
	path := "/apisix/prometheus/metrics"
	handler := NewProxyHandler(path)

	assert.Equal(t, path, handler.path, "handler path should match")
}

func TestProxyHandler_Register(t *testing.T) {
	path := "/apisix/prometheus/metrics"

	// 设置测试路由
	router := gin.Default()
	handler := NewProxyHandler(path)
	handler.Register(router)

	// 创建测试请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", path, nil)
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, 200, w.Code, "status code should be 200")
	assert.Equal(t, "Proxy handler response", w.Body.String(), "response body should match")
}
