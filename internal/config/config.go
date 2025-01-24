package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type (
	// Server 结构体定义了服务器的配置信息。
	Server struct {
		// Port 表示服务器监听的端口号。
		// 可以通过配置文件或环境变量 FSYYFT_APISIX_METRIC_SERVER_PORT 进行配置。
		Port string `yaml:"port"`
	}

	// Prometheus 结构体定义了 Prometheus 指标的配置信息。
	Prometheus struct {
		// Path 表示 Prometheus 指标的暴露路径。
		// 可以通过配置文件或环境变量 FSYYFT_APISIX_METRIC_PROMETHEUS_PATH 进行配置。
		// 如果未指定，默认值为 /metrics。
		Path string `yaml:"path"`
	}

	// Proxy 结构体定义了代理相关的配置信息。
	Proxy struct {
		// Local 结构体定义了本地代理的配置信息。
		Local struct {
			// Path 表示本地代理的路径。
			// 可以通过配置文件或环境变量 FSYYFT_APISIX_METRIC_PROXY_LOCAL_PATH 进行配置。
			Path string `yaml:"path"`
		} `yaml:"local"`
	}

	// Config 结构体定义了应用程序的配置结构。
	// 使用 yaml 标签来映射配置文件中的字段。
	Config struct {
		Server     `yaml:"server"`
		Prometheus `yaml:"prometheus"`
		Proxy      `yaml:"proxy"`
	}
)

// LoadConfig 从指定路径加载配置文件，并返回配置对象。
// 配置文件采用 YAML 格式，支持环境变量覆盖配置值。
// 环境变量前缀为 FSYYFT_APISIX_METRIC_。
func LoadConfig(path string) (*Config, error) {
	// 读取配置文件内容。
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// 解析 YAML 配置文件到 Config 结构体。
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	// 使用环境变量覆盖配置文件中的端口号。
	// 环境变量优先级高于配置文件。
	if port := os.Getenv("FSYYFT_APISIX_METRIC_SERVER_PORT"); port != "" {
		config.Server.Port = port
	}

	// 使用环境变量覆盖 Prometheus 路径配置。
	// 如果未指定，则使用默认值 /metrics。
	if path := os.Getenv("FSYYFT_APISIX_METRIC_PROMETHEUS_PATH"); path != "" {
		config.Prometheus.Path = path
	} else if config.Prometheus.Path == "" {
		config.Prometheus.Path = "/metrics"
	}

	// 使用环境变量覆盖本地代理路径配置。
	if path := os.Getenv("FSYYFT_APISIX_METRIC_PROXY_LOCAL_PATH"); path != "" {
		config.Proxy.Local.Path = path
	} else if config.Proxy.Local.Path == "" {
		config.Proxy.Local.Path = "/apisix/prometheus/metrics"
	}

	return &config, nil
}
