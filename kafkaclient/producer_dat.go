// Package kafkaclient
// kafka client 数据定义
package kafkaclient

import (
	"encoding/json"

	"gopkg.in/yaml.v2"
)

// KafkaConfig kafka服务配置
type KafkaConfig struct {
	BrokerList []string `json:"brokers" yaml:"brokers"`
	ApiVersion string   `json:"apiversion" yaml:"apiversion"`
}

type KafkaProducerCustomConfig struct {
	BufferSize uint `json:"buffersize"`
}

// KafkaProducerConfig kafka服务配置+生产者配置
type KafkaProducerConfig struct {
	ServiceConfig KafkaConfig               `json:"service_config" yaml:"service_config"`
	CustomConfig  KafkaProducerCustomConfig `json:"custom_config" yaml:"custom_config"`
}

func newKafkaProducer() *KafkaProducerConfig {
	return &KafkaProducerConfig{
		ServiceConfig: KafkaConfig{
			BrokerList: []string{"127.0.0.1:9092", "127.0.0.1:9093", "127.0.0.1:9094"},
			ApiVersion: "2.5.0",
		},
		CustomConfig: KafkaProducerCustomConfig{
			BufferSize: 128,
		},
	}
}

func Json2KafkaProducerCfg(js string) (*KafkaProducerConfig, error) {
	cfg := newKafkaProducer()
	err := json.Unmarshal([]byte(js), &cfg)
	return cfg, err
}

func Yaml2KafkaProducerCfg(yml string) (*KafkaProducerConfig, error) {
	cfg := newKafkaProducer()
	err := yaml.Unmarshal([]byte(yml), &cfg)
	return cfg, err
}
