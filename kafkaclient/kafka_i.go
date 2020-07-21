// Package KafkaClient
// kafka 客户端接口定义
package kafkaclient

type KafkaI interface {
	// InitConfigByJson
	// 初始kafka根据json字符串；(根据实现自己定义)
	InitConfigByJson(_strJson string) error
}

// KafkaProducerI kafka生产者接口定义
type KafkaProducerI interface {
	KafkaI
	PushMessage(_msg *KafkaProducerMessage)
}

// KafkaConsumerI kafka消费者接口定义
type KafkaConsumerI interface {
	KafkaI
	PullMessage() []KafkaConsumerMessage
}
