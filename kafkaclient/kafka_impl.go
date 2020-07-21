// Package kafkaclient
// 消费者和生产者实现
package kafkaclient

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	logger "github.com/ydbt/devtool/v1/logger"
	"sync"
	"time"

	sarama "github.com/Shopify/sarama"
)

// KafkaService kafka客户端基本信息
type KafkaClient struct {
	pConfig *sarama.Config
	pLog    logger.LoggerI
}

/* 生产者接口实现 */

// KafkaProducer kafka生产者客户端
type KafkaProducer struct {
	KafkaClient
	producer       sarama.AsyncProducer
	iMaxBufferSize uint
	bufferMessage  chan *KafkaProducerMessage
}

// InitConfigByJson 根据json格式+日志对象，初始化kafka生产者
func (pKC *KafkaProducer) InitConfigByJson(_strJson string, _log logger.LoggerI) error {
	pKC.pLog = _log
	var err error
	cfgProducer := &KafkaProducerConfig{}
	err = json.Unmarshal([]byte(_strJson), cfgProducer)
	if err != nil {
		pKC.pLog.LogWarn("\"InitConfigByJson\":\"[producer] parse config json failed\",\"except\":\"%v\"", err)
		return err
	}
	pKC.pConfig = sarama.NewConfig()
	if pKC.pConfig.Version, err = sarama.ParseKafkaVersion(cfgProducer.ServiceConfig.ApiVersion); err != nil {
		// 版本解析失败使用默认版本
		//		pKC.pConfig.Version = sarama.V2_5_0_0
		pKC.pConfig.Version = sarama.MaxVersion
		pKC.pLog.LogWarn("\"InitConfigByJson\":\"[producer] parse kafka version failed\",\"except\":\"%v\"", err)
	}
	pKC.pConfig.Producer.Partitioner = sarama.NewRandomPartitioner
	pKC.pConfig.Producer.Return.Successes = true
	pKC.pConfig.Producer.Return.Errors = true
	if pKC.producer, err = sarama.NewAsyncProducer(cfgProducer.ServiceConfig.BrokerList, pKC.pConfig); err != nil {
		pKC.pLog.LogWarn("\"InitConfigByJson\":\"[producer] NewAsyncProducer failed\",\"except\":\"%v\"", err)
		return err
	}
	pKC.iMaxBufferSize = cfgProducer.CustomConfig.BufferSize
	signalKill := make(chan os.Signal, 1)
	signal.Notify(signalKill, os.Interrupt)
	// 启动监听事件
	go func() {
		var errPush error
		defer pKC.producer.AsyncClose()
		pKC.bufferMessage = make(chan *KafkaProducerMessage, pKC.iMaxBufferSize)
		for {
			select {
			case objMsg := <-pKC.bufferMessage:
				pKC.producer.Input() <- &sarama.ProducerMessage{Topic: objMsg.TopicName, Key: sarama.ByteEncoder(objMsg.Key), Value: sarama.ByteEncoder(objMsg.Data)}
			case <-signalKill:
				pKC.pLog.LogInfo("\"ProducerPush\":\"[producer] receive kill singal\"")
				return
			case <-pKC.producer.Successes():
				continue
			case errPush = <-pKC.producer.Errors():
				pKC.pLog.LogError("\"ProducerPush\":\"[producer] failed\",\"except\":\"%v\"", errPush)
			}
		}
	}()
	pKC.pLog.LogDebug("\"InitConfigByJson\":\"[producer] init producer success\"")
	return nil
}

// PushMessage kafka生产者推送消息接口
func (pKC *KafkaProducer) PushMessage(_msg *KafkaProducerMessage) {
	pKC.bufferMessage <- _msg
}

/*  消费者接口实现 */
type SynMessage struct{}

// KafkaConsumer kafka消费者客户端
type KafkaConsumer struct {
	KafkaClient
	consumerGroup sarama.ConsumerGroup
	listTopics    []string
	intervalRetry time.Duration
	bExit         bool
	synMsg        chan SynMessage
	mtxBufferMsg  sync.Mutex
	bufferMessage []KafkaConsumerMessage
}

// InitConfigByJson 根据json格式+日志对象，初始化kafka生产者
func (pKC *KafkaConsumer) InitConfigByJson(_strJson string, _log logger.LoggerI) error {
	pKC.pLog = _log
	var err error
	cfgConsumer := &KafkaConsumerConfig{}
	err = json.Unmarshal([]byte(_strJson), cfgConsumer)
	if err != nil {
		pKC.pLog.LogWarn("\"InitConfigByJson\":\"[consumer] parse config json failed\",\"except\":\"%v\"", err)
		return err
	}
	pKC.pConfig = sarama.NewConfig()
	pKC.pConfig.Consumer.Return.Errors = true
	pKC.pConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	pKC.pConfig.Consumer.Offsets.Initial = cfgConsumer.GotoOffset
	pKC.pConfig.ChannelBufferSize = cfgConsumer.CustomConfig.BufferSize
	if pKC.pConfig.Version, err = sarama.ParseKafkaVersion(cfgConsumer.ServiceConfig.ApiVersion); err != nil {
		// 版本解析失败使用最新版本
		pKC.pConfig.Version = sarama.MaxVersion
		pKC.pLog.LogWarn("\"InitConfigByJson\":\"[consumer] parse kafka version failed\",\"except\":\"%v\"", err)
	}
	if pKC.consumerGroup, err = sarama.NewConsumerGroup(cfgConsumer.ServiceConfig.BrokerList, cfgConsumer.GroupId, pKC.pConfig); err != nil {
		pKC.pLog.LogWarn("\"InitConfigByJson\":\"[consumer] init consumer group failed\",\"except\":\"%v\"", err)
	}
	defer func() { _ = pKC.consumerGroup.Close() }()
	go func() {
		for err := range pKC.consumerGroup.Errors() {
			pKC.pLog.LogWarn("\"InitConfigByJson\":\"[consumer] consumer group error\",\"except\":\"%v\"", err)
		}
	}()
	pKC.listTopics = cfgConsumer.Topics
	if cfgConsumer.CustomConfig.IntervalRetry < 0 {
		// 默认重试间隔为1秒
		pKC.intervalRetry = time.Millisecond * 1000
	} else {
		pKC.intervalRetry = time.Millisecond * time.Duration(cfgConsumer.CustomConfig.IntervalRetry)
	}
	go func() {
		signalKill := make(chan os.Signal, 1)
		signal.Notify(signalKill, os.Interrupt)
		select {
		case <-signalKill:
			pKC.bExit = true
			pKC.consumerGroup.Close()
		}
	}()
	pKC.pLog.LogDebug("\"InitConfigByJson\":\"[consumer] init success\"")
	ctx := context.Background()
	pKC.bExit = false
	pKC.synMsg = make(chan SynMessage, 0)
	for {
		if err := pKC.consumerGroup.Consume(ctx, pKC.listTopics, pKC); err != nil {
			pKC.pLog.LogError("\"KafkaConsumer\":\"create consumer failed\",\"except\":\"%v\"", err)
			time.Sleep(pKC.intervalRetry)
		}
		if pKC.bExit {
			return nil
		}
	}
	return nil
}

func (pKC KafkaConsumer) Setup(_ sarama.ConsumerGroupSession) error {
	pKC.pLog.LogDebug("\"ConsumerSession\":\"preparation start subscribe [setup]\"")
	return nil
}

func (pKC KafkaConsumer) Cleanup(_ sarama.ConsumerGroupSession) error {
	pKC.pLog.LogDebug("\"ConsumerSession\":\"terminate consumer [cleanup]\"")
	return nil
}

func (pKC KafkaConsumer) ConsumeClaim(grpSess sarama.ConsumerGroupSession, grpClaim sarama.ConsumerGroupClaim) error {
	defer pKC.mtxBufferMsg.Unlock()
	pKC.mtxBufferMsg.Lock()
	for msg := range grpClaim.Messages() {
		pKC.pLog.LogDebug("Message topic:%s partition:%d offset:%d\n,msg:%s", msg.Topic, msg.Partition, msg.Offset, string(msg.Value))
		grpSess.MarkMessage(msg, "")
	}
	pKC.synMsg <- SynMessage{}
	return nil
}

// 监控订阅 KafKa 消息
func (pKC KafkaConsumer) PullMessage() []KafkaConsumerMessage {
	_ := <-pKC.synMsg
	defer pKC.mtxBufferMsg.Unlock()
	pKC.mtxBufferMsg.Lock()
	listMsg := pKC.bufferMessage
	pKC.bufferMessage = make([]KafkaConsumerMessage, 0)
	return listMsg
}
