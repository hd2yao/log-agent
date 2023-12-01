package kafka

import (
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

var (
	client  sarama.SyncProducer
	msgChan chan *sarama.ProducerMessage
)

// Init 是初始化全局的kafka client
func Init(address []string, chanSize int64) (err error) {
	// 1.生产者配置
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	// 2.连接 kafka
	client, err = sarama.NewSyncProducer(address, config)
	if err != nil {
		logrus.Errorf("kafka:produce closed, err:%v", err)
		return
	}
	// 初始化 msgChan
	msgChan = make(chan *sarama.ProducerMessage, chanSize)
	// 起一个后台的 goroutine 从 msgChan 中读数据
	go sendMsg()
	return
}

// 从 msgChan 中读取 msg，发送给 kafka
func sendMsg() {
	for {
		select {
		case msg := <-msgChan:
			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				logrus.Warning("send msg failed, err:", err)
				return
			}
			logrus.Infof("send msg to kafka success. pid:%v offset:%v", pid, offset)
		}
	}
}

// ToMsgChan 定义一个函数 向外暴露 msgChan
func ToMsgChan(msg *sarama.ProducerMessage) {
	msgChan <- msg
}
