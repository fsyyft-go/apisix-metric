// Copyright 2024 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

// Package web 提供了 Web 服务相关的功能实现，包括反向代理、数据缓存和请求处理等。
package web

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/fsyyft-go/apisix-metric/internal/config"
	"github.com/fsyyft-go/kit/cache"
	"github.com/fsyyft-go/kit/log"
)

const (
	// ETCD_DATA_CACHE_KEY 定义了在缓存中存储 etcd 数据的键名。
	ETCD_DATA_CACHE_KEY = "ETCD_DATA"
	// CACHE_DURATION 定义了缓存的有效期为 1 小时。
	CACHE_DURATION = time.Hour
	// CACHE_REFRESH_THRESHOLD 定义了触发缓存刷新的时间阈值为 50 分钟。
	// 当缓存的剩余有效期小于此值时，将触发异步刷新。
	CACHE_REFRESH_THRESHOLD = 50 * time.Minute
)

// ProxyHandler 实现了一个反向代理处理器，用于处理 HTTP 请求并进行数据转发。
// 它支持从配置和 etcd 中获取数据，并使用缓存来提高性能。
type ProxyHandler struct {
	// cfg 存储了代理的配置信息。
	cfg *config.Config
	// cache 提供了数据缓存功能。
	cache cache.Cache
}

// NewProxyHandler 创建并返回一个新的代理处理器实例。
// 它会初始化缓存，如果缓存初始化失败，将返回一个没有缓存功能的处理器。
//
// 参数：
//   - cfg：代理的配置信息，包含了路由规则、服务定义等。
//
// 返回：
//   - *ProxyHandler：初始化后的代理处理器实例。
func NewProxyHandler(cfg *config.Config) *ProxyHandler {
	// 使用默认配置初始化缓存。
	cache, err := cache.NewCache()
	if err != nil {
		log.WithFields(map[string]interface{}{
			"error": err,
		}).Error("初始化缓存失败。")
		return &ProxyHandler{
			cfg: cfg,
		}
	}

	return &ProxyHandler{
		cfg:   cfg,
		cache: cache,
	}
}

// Register 将处理器的路由规则注册到 Gin 引擎中。
//
// 参数：
//   - r：Gin 的路由引擎实例。
func (h *ProxyHandler) Register(r *gin.Engine) {
	// 注册 GET 请求处理器。
	r.GET(h.cfg.Proxy.Local.Path, h.ReverseProxyHandler)
}

// ReverseProxyHandler 处理来自客户端的请求，将其转发到 APISIX 的 Prometheus 指标接口。
// 在转发过程中，会对响应内容进行处理，替换特定的标识符。
//
// 参数：
//   - c：Gin 的上下文对象，包含了请求和响应的信息。
func (h *ProxyHandler) ReverseProxyHandler(c *gin.Context) {
	config := h.cfg.Proxy

	// 创建反向代理实例。
	proxy := &httputil.ReverseProxy{
		// Director 函数负责修改请求信息。
		Director: func(req *http.Request) {
			// 设置目标服务器的协议。
			req.URL.Scheme = config.Remote.Scheme
			// 设置目标服务器的主机地址。
			req.URL.Host = config.Remote.Host
			// 设置目标服务器的路径。
			req.URL.Path = config.Remote.Path
			// 设置请求头中的 Host 字段。
			req.Host = config.Remote.Host

			// 记录请求转发的信息。
			log.WithFields(map[string]interface{}{
				"scheme": req.URL.Scheme,
				"host":   req.URL.Host,
				"path":   req.URL.Path,
				"method": req.Method,
			}).Debug("转发请求到目标服务器。")
		},
		// ModifyResponse 函数负责处理响应内容。
		ModifyResponse: func(resp *http.Response) error {
			// 记录响应状态码。
			log.WithFields(map[string]interface{}{
				"status_code": resp.StatusCode,
				"status":      resp.Status,
				"url":         resp.Request.URL.String(),
				"headers":     resp.Header,
			}).Debug("收到目标服务器响应。")

			// 只处理状态码为 200 的响应。
			if resp.StatusCode != http.StatusOK {
				log.WithFields(map[string]interface{}{
					"status_code": resp.StatusCode,
					"status":      resp.Status,
					"url":         resp.Request.URL.String(),
					"headers":     resp.Header,
				}).Warn("目标服务器返回非 200 状态码。")
				return nil
			}

			// 检查响应体是否为空。
			if resp.Body == nil {
				log.WithFields(map[string]interface{}{
					"url": resp.Request.URL.String(),
				}).Warn("目标服务器返回空响应体。")
				return nil
			}

			// 读取响应体内容。
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.WithFields(map[string]interface{}{
					"error": err,
					"url":   resp.Request.URL.String(),
				}).Error("读取响应体失败。")
				return err
			}
			resp.Body.Close()

			// 检查响应体内容是否为空。
			if len(body) == 0 {
				log.WithFields(map[string]interface{}{
					"url": resp.Request.URL.String(),
				}).Warn("响应体内容为空。")
			}

			// 将响应体转换为字符串。
			content := string(body)

			// 获取用于替换的数据。
			replaceData := h.getData(c.Request.Context())

			// 检查替换数据是否为空。
			if len(replaceData) == 0 {
				log.WithFields(map[string]interface{}{
					"url": resp.Request.URL.String(),
				}).Warn("没有可用的替换数据。")
			}

			// 使用获取的数据进行内容替换。
			for oldStr, newStr := range replaceData {
				content = strings.ReplaceAll(content, oldStr, newStr)
			}

			// 使用处理后的内容更新响应体。
			resp.Body = io.NopCloser(strings.NewReader(content))
			// 更新内容长度。
			contentLength := int64(len(content))
			resp.ContentLength = contentLength
			// 更新 Content-Length 响应头。
			resp.Header.Set("Content-Length", strconv.Itoa(len(content)))

			// 记录响应处理完成的信息。
			log.WithFields(map[string]interface{}{
				"url":             resp.Request.URL.String(),
				"content_length":  contentLength,
				"replace_count":   len(replaceData),
				"original_length": len(body),
			}).Debug("响应内容处理完成。")

			return nil
		},
	}

	// 处理请求。
	proxy.ServeHTTP(c.Writer, c.Request)
}

// getData 获取用于内容替换的数据，包括配置中的服务、路由信息以及从 etcd 获取的数据。
// 如果从 etcd 获取数据失败，将只返回配置中的数据。
//
// 参数：
//   - ctx：上下文对象，用于控制请求的生命周期。
//
// 返回：
//   - map[string]string：包含了所有用于替换的键值对数据。
func (h *ProxyHandler) getData(ctx context.Context) map[string]string {
	// 创建用于存储结果的映射。
	result := make(map[string]string)

	// 从配置中获取服务数据。
	for k, v := range h.cfg.Proxy.Service {
		result[k] = v
	}

	// 从配置中获取路由数据。
	for k, v := range h.cfg.Proxy.Route {
		result[k] = v
	}

	// 从缓存或 etcd 获取数据。
	etcdData, err := h.getDataFromCache(ctx, h.cfg.Proxy.Etcd.Prefix)
	if err != nil {
		// 记录错误日志但继续执行。
		log.WithFields(map[string]interface{}{
			"error": err,
		}).Error("从缓存获取数据失败。")
		// 使用空映射作为默认值。
		etcdData = make(map[string]string)
	}

	// 将 etcd 数据添加到结果中。
	for k, v := range etcdData {
		result[k] = v
	}

	return result
}

// getDataFromEtcd 从 etcd 中获取指定前缀的所有数据。
// 它会连接到 etcd 服务器，获取数据并解析其中的 name 字段。
//
// 参数：
//   - ctx：上下文对象，用于控制请求的生命周期。
//   - prefix：要获取的键的前缀。
//
// 返回：
//   - map[string]string：解析后的键值对数据。
//   - error：如果获取或解析过程中发生错误。
func (h *ProxyHandler) getDataFromEtcd(ctx context.Context, prefix string) (map[string]string, error) {
	// 创建 etcd 客户端配置。
	etcdConfig := clientv3.Config{
		Endpoints:   h.cfg.Proxy.Etcd.Endpoints,
		DialTimeout: time.Duration(h.cfg.Proxy.Etcd.Timeout) * time.Second,
	}

	// 如果配置了认证信息，添加到配置中。
	if h.cfg.Proxy.Etcd.Auth.Username != "" && h.cfg.Proxy.Etcd.Auth.Password != "" {
		etcdConfig.Username = h.cfg.Proxy.Etcd.Auth.Username
		etcdConfig.Password = h.cfg.Proxy.Etcd.Auth.Password
	}

	// 创建 etcd 客户端。
	client, err := clientv3.New(etcdConfig)
	if err != nil {
		log.WithFields(map[string]interface{}{
			"error": err,
		}).Error("创建 etcd 客户端失败。")
		return nil, fmt.Errorf("创建 etcd 客户端失败：%v", err)
	}
	defer client.Close()

	// 记录开始获取数据的日志。
	log.WithFields(map[string]interface{}{
		"prefix": prefix,
	}).Info("正在从 etcd 获取数据。")

	// 构建需要查询的前缀列表。
	prefixes := []string{
		prefix + "/routes/",
		prefix + "/services/",
	}

	// 创建用于存储结果的映射。
	result := make(map[string]string)

	// 遍历所有前缀进行查询。
	for _, p := range prefixes {
		// 使用前缀获取数据。
		resp, err := client.Get(ctx, p, clientv3.WithPrefix())
		if err != nil {
			log.WithFields(map[string]interface{}{
				"prefix": p,
				"error":  err,
			}).Error("从 etcd 获取数据失败。")
			return nil, fmt.Errorf("从 etcd 获取数据失败，前缀 %s：%v", p, err)
		}

		// 将结果添加到映射中。
		for _, kv := range resp.Kvs {
			result[string(kv.Key)] = string(kv.Value)
		}
	}

	// 创建用于存储处理后结果的映射。
	nameMap := make(map[string]string)

	// 遍历结果，提取 name 字段。
	for key, value := range result {
		// 尝试解析 JSON 数据。
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(value), &data); err != nil {
			// 记录跳过的数据。
			log.WithFields(map[string]interface{}{
				"key":   key,
				"value": value,
				"error": err.Error(),
			}).Debug("跳过非 JSON 数据。")
			continue
		}

		// 尝试获取 name 字段。
		if name, ok := data["name"].(string); ok {
			// 处理键，去掉前缀。
			for _, p := range prefixes {
				if strings.HasPrefix(key, p) {
					key = strings.TrimPrefix(key, p)
					break
				}
			}
			nameMap[key] = name
		}
	}

	// 记录获取数据成功的日志。
	log.WithFields(map[string]interface{}{
		"prefix": prefix,
		"count":  len(nameMap),
	}).Info("从 etcd 获取数据成功。")

	return nameMap, nil
}

// getDataFromCache 从缓存中获取数据，如果缓存不存在或已过期则从 etcd 获取。
// 当缓存即将过期时，会异步刷新缓存数据。
//
// 参数：
//   - ctx：上下文对象，用于控制请求的生命周期。
//   - prefix：要获取的键的前缀。
//
// 返回：
//   - map[string]string：缓存的键值对数据。
//   - error：如果获取过程中发生错误。
func (h *ProxyHandler) getDataFromCache(ctx context.Context, prefix string) (map[string]string, error) {
	// 如果缓存未初始化，直接从 etcd 获取数据。
	if h.cache == nil {
		log.WithFields(map[string]interface{}{
			"prefix": prefix,
		}).Warn("缓存未初始化，将直接从 etcd 获取数据。")
		return h.getDataFromEtcd(ctx, prefix)
	}

	// 创建类型安全的缓存实例。
	typedCache := cache.AsTypedCache[map[string]string](h.cache)

	// 尝试从缓存获取数据。
	cacheData, exists, ttl := typedCache.GetWithTTL(ETCD_DATA_CACHE_KEY)
	if exists {
		// 如果缓存存在且剩余时间小于阈值，异步刷新缓存。
		if ttl > 0 && ttl < CACHE_REFRESH_THRESHOLD {
			log.WithFields(map[string]interface{}{
				"ttl": ttl.String(),
			}).Debug("缓存即将过期，启动异步刷新。")

			go func() {
				newCtx := context.Background()
				newData, err := h.getDataFromEtcd(newCtx, prefix)
				if err != nil {
					log.WithFields(map[string]interface{}{
						"error": err,
					}).Error("异步刷新缓存失败。")
					return
				}

				if ok := typedCache.SetWithTTL(ETCD_DATA_CACHE_KEY, newData, CACHE_DURATION); !ok {
					log.Error("异步更新缓存失败。")
				}
			}()
		}
		return cacheData, nil
	}

	// 如果缓存不存在，从 etcd 获取数据。
	data, err := h.getDataFromEtcd(ctx, prefix)
	if err != nil {
		return nil, err
	}

	// 将数据写入缓存。
	if ok := typedCache.SetWithTTL(ETCD_DATA_CACHE_KEY, data, CACHE_DURATION); !ok {
		log.Error("设置缓存数据失败。")
		return data, nil
	}

	log.WithFields(map[string]interface{}{
		"size": len(data),
	}).Info("成功设置缓存数据。")

	return data, nil
}

// Close 关闭处理器并释放相关资源。
// 主要用于关闭缓存连接。
//
// 返回：
//   - error：如果关闭过程中发生错误。
func (h *ProxyHandler) Close() error {
	if h.cache != nil {
		return h.cache.Close()
	}
	return nil
}
