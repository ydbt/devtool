// Package kafkaclient
// kafka client 数据定义
package kafkaclient

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
