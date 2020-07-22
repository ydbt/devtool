// Package kafkaclient
// 消费者和生产者实现
package kafkaclient

import (
	"context"
	"os"
	"time"

	logger "github.com/ydbt/devtool/v3/logger"
	"github.com/ydbt/devtool/v3/usetool"

	sarama "github.com/Shopify/sarama"
)

/*  消费者接口实现 */
type SynMessage struct{}

// KafkaConsumer kafka消费者客户端
type KafkaConsumer struct {
	KafkaClient
	consumerGroup sarama.ConsumerGroup
	listTopics    []string
	intervalRetry time.Duration
	chanTopicMsg  map[string]chan KafkaConsumerMessage
}

func NewKafkaConsumer(cfgConsumer *KafkaConsumerConfig, lgi logger.LogI) *KafkaConsumer {
	lgi.Tracef("new kafka consumer beg")
	defer lgi.Tracef("new kafka consumer end")
	kfk := &KafkaConsumer{}
	kfk.pLog = lgi
	kfk.pConfig = sarama.NewConfig()
	var err error
	kfk.pConfig.Consumer.Return.Errors = true
	kfk.pConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	kfk.pConfig.Consumer.Offsets.Initial = cfgConsumer.GotoOffset
	if cfgConsumer.CustomConfig.BufferSize < 16 {
		kfk.pConfig.ChannelBufferSize = 16
	} else if cfgConsumer.CustomConfig.BufferSize > 1024 {
		kfk.pConfig.ChannelBufferSize = 1024
	}
	if kfk.pConfig.Version, err = sarama.ParseKafkaVersion(cfgConsumer.ServiceConfig.ApiVersion); err != nil {
		// 版本解析失败使用最新版本
		kfk.pConfig.Version = sarama.MaxVersion
		kfk.pLog.Warnf("parse config kafka consumer version exception(%v),use default %v", err, kfk.pConfig.Version)
	}
	if kfk.consumerGroup, err = sarama.NewConsumerGroup(cfgConsumer.ServiceConfig.BrokerList, cfgConsumer.GroupId, kfk.pConfig); err != nil {
		kfk.pLog.Warnf("init kafka consumer group exception(%v)", err)
	}
	kfk.listTopics = make([]string, len(cfgConsumer.Topics))
	kfk.chanTopicMsg = make(map[string]chan KafkaConsumerMessage)
	for i, topic := range cfgConsumer.Topics {
		kfk.listTopics[i] = topic
		kfk.chanTopicMsg[topic] = make(chan KafkaConsumerMessage, kfk.pConfig.ChannelBufferSize/3)
	}
	if cfgConsumer.CustomConfig.IntervalRetry < 50 {
		cfgConsumer.CustomConfig.IntervalRetry = 100
	} else if cfgConsumer.CustomConfig.IntervalRetry > 60*1000 {
		cfgConsumer.CustomConfig.IntervalRetry = 60000
	}
	kfk.intervalRetry = time.Millisecond * time.Duration(cfgConsumer.CustomConfig.IntervalRetry)
	return kfk
}

func (pKC *KafkaConsumer) Setup(_ sarama.ConsumerGroupSession) error {
	pKC.pLog.Tracef("\"ConsumerSession\":\"preparation start subscribe [setup]\"")
	return nil
}

func (pKC *KafkaConsumer) Cleanup(_ sarama.ConsumerGroupSession) error {
	pKC.pLog.Tracef("\"ConsumerSession\":\"terminate consumer [cleanup]\"")
	return nil
}

func (pKC *KafkaConsumer) ConsumeClaim(grpSess sarama.ConsumerGroupSession, grpClaim sarama.ConsumerGroupClaim) error {
	pKC.pLog.Tracef("kafka consumer pull once beg")
	defer pKC.pLog.Tracef("kafka consumer pull once end")
	for msg := range grpClaim.Messages() {
		pKC.pLog.Tracef("Message topic:%s partition:%d offset:%d\n,msg:%s", msg.Topic, msg.Partition, msg.Offset, string(msg.Value))
		pKC.chanTopicMsg[msg.Topic] <- KafkaConsumerMessage{
			Partition: msg.Partition,
			Offset:    msg.Offset,
			TopicName: msg.Topic,
			Key:       msg.Key,
			Data:      msg.Value,
		}
		grpSess.MarkMessage(msg, "")
	}
	return nil
}

func (pKC *KafkaConsumer) GoTimerPullMsg(psi usetool.ProcessSignalI) {
	go pKC.TimerPullMsg(psi)
}

func (pKC *KafkaConsumer) TimerPullMsg(psi usetool.ProcessSignalI) {
	pKC.pLog.Tracef("kafka consumer timer pull msg beg")
	defer pKC.pLog.Tracef("kafka consumer timer pull msg end")
	defer func() { _ = pKC.consumerGroup.Close() }()
	//	timerSignal := time.NewTimer(time.Second)
	interruptSignal := psi.TopicOSSignal(os.Interrupt)
	killSignal := psi.TopicOSSignal(os.Kill)

	ctx := context.Background()
	for {
		select {
		case <-interruptSignal:
			pKC.pLog.Infof("[consumer] receive interrupt singal")
			break
		case <-killSignal:
			pKC.pLog.Infof("[consumer] receive kill singal")
			break
		default:
			if err := pKC.consumerGroup.Consume(ctx, pKC.listTopics, pKC); err != nil {
				pKC.pLog.Errorf("[consumer] consumer exception(%v)", err)
				time.Sleep(pKC.intervalRetry)
			}
		}
	}
	pKC.consumerGroup.Close()

}

// 监控订阅 KafKa 消息
func (pKC *KafkaConsumer) PullMessage(topic string) KafkaConsumerMessage {
	msg := <-pKC.chanTopicMsg[topic]
	return msg
}
