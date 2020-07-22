// Package kafkaclient
// kafka client 数据定义
package kafkaclient

import (
	"encoding/json"

	"gopkg.in/yaml.v2"
)

type KafkaConsumerCustomConfig struct {
	BufferSize    int `json:"buffersize",yaml:"buffersize"`
	IntervalRetry int `json:"retryinterval",yaml:"retryinterval"`
}

// KafkaProducerConfig kafka服务配置+消费者配置
type KafkaConsumerConfig struct {
	ServiceConfig KafkaConfig               `json:"service_config",yaml:"service_config"`
	CustomConfig  KafkaConsumerCustomConfig `json:"custom_config",yaml:"custom_config"`
	GroupId       string                    `json:"groupid",yaml:"groupid"`
	Topics        []string                  `json:"topics",yaml:"topics"`
	GotoOffset    int64                     `json:"goto_offset",yaml:"goto_offset"`
}

func newKafkaConsumer() *KafkaConsumerConfig {
	return &KafkaConsumerConfig{
		ServiceConfig: KafkaConfig{
			BrokerList: []string{"127.0.0.1:9092", "127.0.0.1:9093", "127.0.0.1:9094"},
			ApiVersion: "2.5.0.0",
		},
		CustomConfig: KafkaConsumerCustomConfig{
			BufferSize:    128,
			IntervalRetry: 10,
		},
		GroupId:    "ydbt-consumer",
		Topics:     []string{"test"},
		GotoOffset: -1,
	}
}

func Json2KafkaConsumerCfg(js string) (*KafkaConsumerConfig, error) {
	cfg := newKafkaConsumer()
	err := json.Unmarshal([]byte(js), &cfg)
	return cfg, err
}

func Yaml2KafkaConsumerCfg(yml string) (*KafkaConsumerConfig, error) {
	cfg := newKafkaConsumer()
	err := yaml.Unmarshal([]byte(yml), &cfg)
	return cfg, err
}
