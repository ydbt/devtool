// Package kafkaclient
// kafka client 数据定义
package kafkaclient

// KafkaConfig kafka服务配置
type KafkaConfig struct {
	BrokerList []string `json:"brokers"`
	ApiVersion string   `json:"apiversion"`
}

type KafkaProducerCustomConfig struct {
	BufferSize uint `json:"buffersize"`
}

// KafkaProducerConfig kafka服务配置+生产者配置
type KafkaProducerConfig struct {
	ServiceConfig KafkaConfig               `json:"service_config"`
	CustomConfig  KafkaProducerCustomConfig `json:"custom_config"`
}

type KafkaConsumerCustomConfig struct {
	BufferSize    int `json:"buffersize"`
	IntervalRetry int `json:"retry_interval"`
}

// KafkaProducerConfig kafka服务配置+消费者配置
type KafkaConsumerConfig struct {
	ServiceConfig KafkaConfig               `json:"service_config"`
	CustomConfig  KafkaConsumerCustomConfig `json:"custom_config"`
	GroupId       string                    `json:"groupid"`
	Topics        []string                  `json:"topics"`
	GotoOffset    int64                     `json:"goto_offset"`
}

//KafkaMessage 推送到Kafka消息信息
type KafkaProducerMessage struct {
	TopicName string
	Key       []byte
	Data      []byte
}

type KafkaConsumerMessage struct {
	Partition int32
	Offset    int64
	TopicName string
	Key       []byte
	Data      []byte
}
