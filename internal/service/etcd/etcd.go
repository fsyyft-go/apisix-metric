// Copyright 2024 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/fsyyft-go/apisix-metric/internal/config"
	"github.com/fsyyft-go/apisix-metric/pkg/log"
)

// Service 结构体封装了 etcd 服务的功能。
type Service struct {
	// cfg 存储服务配置信息。
	cfg *config.Config
	// client 是 etcd 客户端实例。
	client *clientv3.Client
}

// NewService 创建并返回一个新的 Service 实例。
//
// 参数：
//   - cfg：配置信息，包含 etcd 连接信息等。
//
// 返回：
//   - *Service：初始化后的 Service 实例。
//   - error：如果初始化失败则返回错误信息。
func NewService(cfg *config.Config) (*Service, error) {
	// 创建 etcd 客户端配置
	etcdConfig := clientv3.Config{
		Endpoints:   cfg.Proxy.Etcd.Endpoints,
		DialTimeout: time.Duration(cfg.Proxy.Etcd.Timeout) * time.Second,
	}

	// 如果配置了认证信息，则添加到配置中
	if cfg.Proxy.Etcd.Auth.Username != "" && cfg.Proxy.Etcd.Auth.Password != "" {
		etcdConfig.Username = cfg.Proxy.Etcd.Auth.Username
		etcdConfig.Password = cfg.Proxy.Etcd.Auth.Password
	}

	// 创建 etcd 客户端
	client, err := clientv3.New(etcdConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd client: %v", err)
	}

	return &Service{
		cfg:    cfg,
		client: client,
	}, nil
}

// Close 关闭 etcd 客户端连接。
func (s *Service) Close() error {
	if s.client != nil {
		return s.client.Close()
	}
	return nil
}

// GetData 从 etcd 获取指定前缀的所有数据。
//
// 参数：
//   - ctx：上下文对象，用于控制请求的生命周期。
//   - prefix：要获取的键的前缀。
//
// 返回：
//   - map[string]string：键值对形式的数据。
//   - error：如果获取失败则返回错误信息。
func (s *Service) GetData(ctx context.Context, prefix string) (map[string]string, error) {
	// 记录开始获取数据的日志
	log.WithFields(map[string]interface{}{
		"prefix": prefix,
	}).Info("Getting data from etcd")

	// 构建需要查询的前缀列表
	prefixes := []string{
		prefix + "/routes/",
		prefix + "/services/",
	}

	// 创建结果映射
	result := make(map[string]string)

	// 遍历所有前缀进行查询
	for _, p := range prefixes {
		// 使用前缀获取数据
		resp, err := s.client.Get(ctx, p, clientv3.WithPrefix())
		if err != nil {
			log.WithFields(map[string]interface{}{
				"prefix": p,
				"error":  err,
			}).Error("Failed to get data from etcd")
			return nil, fmt.Errorf("failed to get data from etcd with prefix %s: %v", p, err)
		}

		// 将结果添加到结果映射中
		for _, kv := range resp.Kvs {
			result[string(kv.Key)] = string(kv.Value)
		}
	}

	// 创建新的映射来存储处理后的结果
	nameMap := make(map[string]string)

	// 遍历结果，提取 name 字段
	for key, value := range result {
		// 尝试解析 JSON
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(value), &data); err != nil {
			// 记录跳过的数据
			log.WithFields(map[string]interface{}{
				"key":   key,
				"value": value,
				"error": err.Error(),
			}).Debug("Skipping non-JSON data")
			// 如果不是有效的 JSON，跳过这条记录
			continue
		}

		// 尝试获取 name 字段
		if name, ok := data["name"].(string); ok {
			// 处理 key，去掉前缀
			for _, p := range prefixes {
				if strings.HasPrefix(key, p) {
					// 去掉前缀，获取最后一个节点
					key = strings.TrimPrefix(key, p)
					break
				}
			}
			nameMap[key] = name
		}
	}

	// 记录获取数据成功的日志
	log.WithFields(map[string]interface{}{
		"prefix": prefix,
		"count":  len(nameMap),
	}).Info("Successfully got data from etcd")

	return nameMap, nil
}

// GetDataWithPrefix 从 etcd 获取带有配置前缀的数据。
//
// 参数：
//   - ctx：上下文对象，用于控制请求的生命周期。
//
// 返回：
//   - map[string]string：键值对形式的数据。
//   - error：如果获取失败则返回错误信息。
func (s *Service) GetDataWithPrefix(ctx context.Context) (map[string]string, error) {
	return s.GetData(ctx, s.cfg.Proxy.Etcd.Prefix)
}
