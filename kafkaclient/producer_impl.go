// Package kafkaclient
// 消费者和生产者实现
package kafkaclient

import (
	"os"

	sarama "github.com/Shopify/sarama"

	"github.com/ydbt/devtool/v3/logger"
	"github.com/ydbt/devtool/v3/usetool"
)

// KafkaService kafka客户端基本信息
type KafkaClient struct {
	pConfig *sarama.Config
	pLog    logger.LogI
}

/* 生产者接口实现 */

// KafkaProducer kafka生产者客户端
type KafkaProducer struct {
	KafkaClient
	producer       sarama.AsyncProducer
	iMaxBufferSize uint
	bufferMessage  chan *KafkaProducerMessage
}

func NewKafkaProducer(cfgProducer *KafkaProducerConfig, lgi logger.LogI) *KafkaProducer {
	kfk := &KafkaProducer{}
	kfk.pLog = lgi
	kfk.pConfig = sarama.NewConfig()
	var err error
	if kfk.pConfig.Version, err = sarama.ParseKafkaVersion(cfgProducer.ServiceConfig.ApiVersion); err != nil {
		// 版本解析失败使用默认版本
		//		pKC.pConfig.Version = sarama.V2_5_0_0
		kfk.pConfig.Version = sarama.MaxVersion
		kfk.pLog.Warnf("[producer] parse kafka version failed, exception:%v", err)
	}
	kfk.pConfig.Producer.Partitioner = sarama.NewRandomPartitioner
	kfk.pConfig.Producer.Return.Successes = true
	kfk.pConfig.Producer.Return.Errors = true
	if kfk.producer, err = sarama.NewAsyncProducer(cfgProducer.ServiceConfig.BrokerList, kfk.pConfig); err != nil {
		kfk.pLog.Warnf("[producer] NewAsyncProducer failed, exception:%v, cfg:%v", err, cfgProducer)
		return nil
	}
	kfk.iMaxBufferSize = cfgProducer.CustomConfig.BufferSize
	return kfk
}

func (pKC *KafkaProducer) GoTimerPushMsg(psi usetool.ProcessSignalI) {
	go pKC.TimerPushMsg(psi)
}

func (pKC *KafkaProducer) TimerPushMsg(psi usetool.ProcessSignalI) {
	pKC.pLog.Debugf("kafka producer async push msg beg")
	defer pKC.pLog.Debugf("kafka producer async push msg end")
	//	timerSignal := time.NewTimer(time.Second
	interruptSignal := psi.TopicOSSignal(os.Interrupt)
	killSignal := psi.TopicOSSignal(os.Kill)
	var errPush error
	defer pKC.producer.AsyncClose()
	pKC.bufferMessage = make(chan *KafkaProducerMessage, pKC.iMaxBufferSize)
	// 缓存队列初始化
	for {
		select {
		case objMsg := <-pKC.bufferMessage:
			pKC.producer.Input() <- &sarama.ProducerMessage{Topic: objMsg.TopicName, Key: sarama.ByteEncoder(objMsg.Key), Value: sarama.ByteEncoder(objMsg.Data)}
		case <-interruptSignal:
			pKC.pLog.Infof("receive kill singal")
			return
		case <-killSignal:
			pKC.pLog.Infof("receive kill singal")
			return
		case <-pKC.producer.Successes():
			continue
		case errPush = <-pKC.producer.Errors():
			pKC.pLog.Errorf("[producer] exception:%v", errPush)
		}
	}
}

// PushMessage kafka生产者推送消息接口
func (pKC *KafkaProducer) PushMessage(_msg *KafkaProducerMessage) {
	pKC.bufferMessage <- _msg
}
