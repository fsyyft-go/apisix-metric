// Copyright 2024 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package etcd

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/fsyyft-go/apisix-metric/internal/config"
	"github.com/fsyyft-go/apisix-metric/internal/log"
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

	// 使用前缀获取数据
	resp, err := s.client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		log.WithFields(map[string]interface{}{
			"prefix": prefix,
			"error":  err,
		}).Error("Failed to get data from etcd")
		return nil, fmt.Errorf("failed to get data from etcd: %v", err)
	}

	// 创建结果映射
	result := make(map[string]string)
	for _, kv := range resp.Kvs {
		result[string(kv.Key)] = string(kv.Value)
	}

	// 记录获取数据成功的日志
	log.WithFields(map[string]interface{}{
		"prefix": prefix,
		"count":  len(result),
	}).Info("Successfully got data from etcd")

	return result, nil
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
