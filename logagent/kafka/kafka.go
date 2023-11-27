package kafka

import (
    "github.com/IBM/sarama"
    "github.com/sirupsen/logrus"
)

var (
    Client sarama.SyncProducer
)

// Init 是初始化全局的kafka Client
func Init(address []string) (err error) {
    // 1.生产者配置
    config := sarama.NewConfig()
    config.Producer.RequiredAcks = sarama.WaitForAll
    config.Producer.Partitioner = sarama.NewRandomPartitioner
    config.Producer.Return.Successes = true

    // 2.连接 kafka
    Client, err = sarama.NewSyncProducer(address, config)
    if err != nil {
        logrus.Errorf("kafka:produce closed, err:%v", err)
        return
    }
    return
}
