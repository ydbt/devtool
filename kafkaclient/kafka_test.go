package kafkaclient_test

import (
	"os"
	"testing"
	"time"

	"github.com/ydbt/devtool/v3/kafkaclient"
	. "github.com/ydbt/devtool/v3/kafkaclient"
	logger "github.com/ydbt/devtool/v3/logger"
	"github.com/ydbt/devtool/v3/usetool"
)

var pLog *logger.LoggerFake

const gc_strTopicName = "ut-go-client-test"

func init() {
	pLog = new(logger.LoggerFake)
	pLog.Init(true)
	pLog.SetLevel("trace")
	//	pLog.Init(false)
}

const gc_strProducerJson = "{\"service_config\":{\"brokers\":[\"kafka01:9092\",\"kafka02:9092\",\"kafka03:9092\"],\"apiversion\":\"2.5.0\"},\"custom_config\":{\"buffersize\":3}}"

func TestProducer(t *testing.T) {
	cfgProducer, err := kafkaclient.Json2KafkaProducerCfg(gc_strProducerJson)
	if err != nil {
		t.Error(err)
		return
	}
	producer := kafkaclient.NewKafkaProducer(cfgProducer, pLog)
	psr := usetool.NewProcessSignaler()
	producer.GoTimerPushMsg(psr)
	time.Sleep(time.Second)
	var objMsg0 = KafkaProducerMessage{
		TopicName: gc_strTopicName,
		Data:      []byte("hello world!")}
	producer.PushMessage(&objMsg0)
	var objMsg1 = KafkaProducerMessage{
		TopicName: gc_strTopicName,
		Data:      []byte("你好！")}
	producer.PushMessage(&objMsg1)
	time.Sleep(time.Second)
	psr.WaitSignalProcess(os.Interrupt)
}

const gc_strConsumerJson = "{\"service_config\":{\"brokers\":[\"kafka01:9092\",\"kafka02:9092\",\"kafka03:9092\"],\"apiversion\":\"2.5.0\"},\"custom_config\":{\"buffersize\":128,\"retryinterval\":2000},\"groupid\":\"ut-go-active-01\",\"topics\":[\"ut-go-client-test\"],\"partitions\":[0],\"goto_offset\":-2}"

func TestConsumer(t *testing.T) {
	cfg, err := kafkaclient.Json2KafkaConsumerCfg(gc_strConsumerJson)
	if err != nil {
		t.Error(err)
		return
	}
	psr := usetool.NewProcessSignaler()
	consumer := kafkaclient.NewKafkaConsumer(cfg, pLog)
	consumer.GoTimerPullMsg(psr)
	time.Sleep(time.Second)
	msg := consumer.PullMessage(gc_strTopicName)
	t.Log(msg)
	time.Sleep(time.Second)
	psr.WaitSignalProcess(os.Interrupt)
}
