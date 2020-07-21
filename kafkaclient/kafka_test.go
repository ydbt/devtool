package kafkaclient_test

import (
	. "gitlab.qn.com/ydbt/kafkaclient"
	logger "gitlab.qn.com/ydbt/logger"
	"testing"
	"time"
)

var pLog *logger.LoggerFake

const gc_strTopicName = "ut-go-client-test"

func init() {
	pLog = new(logger.LoggerFake)
	pLog.Init(true)
	//	pLog.Init(false)
}

//const gc_strProducerJson = "{\"service_config\":{\"brokers\":[\"kafka001:9092\",\"kafka002:9092\",\"kafka003:9092\"],\"apiversion\":\"1.1.0\"},\"custom_config\":{\"buffersize\":3}}"
//func TestProducer(_t *testing.T) {
//	var producer KafkaProducer
//	producer.InitConfigByJson(gc_strProducerJson, pLog)
//	var objMsg = KafkaProducerMessage{
//		TopicName: gc_strTopicName,
//		Data:      []byte("hello world!")}
//	producer.PushMessage(&objMsg)
//}
//
const gc_strConsumerJson = "{\"service_config\":{\"brokers\":[\"kafka001:9092\",\"kafka002:9092\",\"kafka003:9092\"],\"apiversion\":\"1.1.0\"},\"custom_config\":{},\"groupid\":\"ut-go-active-01\",\"topic\":\"ut-go-client-test\",\"partitions\":[0],\"goto_offset\":-2}"

func TestConsumer(_t *testing.T) {
	var consumer KafkaConsumer
	consumer.InitConfigByJson(gc_strConsumerJson, pLog)
	time.Sleep(time.Second * 10)
}
