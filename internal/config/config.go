package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config 结构体定义了应用程序的配置结构。
// 使用 yaml 标签来映射配置文件中的字段。
type Config struct {
	Server struct {
		// Port 表示服务器监听的端口号。
		// 可以通过配置文件或环境变量 FSYYFT_APISIX_METRIC_SERVER_PORT 进行配置。
		Port string `yaml:"port"`
	} `yaml:"server"`

	Prometheus struct {
		// Path 表示 Prometheus 指标的暴露路径。
		// 可以通过配置文件或环境变量 FSYYFT_APISIX_METRIC_PROMETHEUS_PATH 进行配置。
		// 如果未指定，默认值为 /metrics。
		Path string `yaml:"path"`
	} `yaml:"prometheus"`
}

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

	return &config, nil
}
