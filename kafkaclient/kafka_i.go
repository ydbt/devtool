// Package KafkaClient
// kafka 客户端接口定义
package kafkaclient

// KafkaProducerI kafka生产者接口定义
type KafkaProducerI interface {
	PushMessage(_msg *KafkaProducerMessage)
}

// KafkaConsumerI kafka消费者接口定义
type KafkaConsumerI interface {
	PullMessage() []KafkaConsumerMessage
}
