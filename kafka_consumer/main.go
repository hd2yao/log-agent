package main

import (
    "fmt"
    "sync"

    "github.com/IBM/sarama"
)

// kafka consumer(消费者)

func main() {
    // 创建新的消费者
    consumer, err := sarama.NewConsumer([]string{"127.0.0.1:9092"}, nil)
    if err != nil {
        fmt.Printf("fail to start consumer, err:%v\n", err)
        return
    }
    // 拿到指定 topic 下面的所有分区列表
    partitionList, err := consumer.Partitions("web_log") // 根据 topic 取到所有的分区
    if err != nil {
        fmt.Printf("fail to get list of partition:err%v\n", err)
        return
    }
    fmt.Println(partitionList)

    var wg sync.WaitGroup
    for partition := range partitionList { // 遍历所有的分区
        // 针对每个分区创建一个对应的分区消费者
        pc, err := consumer.ConsumePartition("web_log", int32(partition), sarama.OffsetOldest)
        if err != nil {
            fmt.Printf("failed to start consumer for partition %d, err:%v\v", partition, err)
            return
        }
        defer pc.AsyncClose()
        // 异步从每个分区消费信息
        wg.Add(1)
        go func(sarama.PartitionConsumer) {
            for message := range pc.Messages() {
                fmt.Printf("Partition:%d Offset:%d Key:%s Value:%s\n", message.Partition, message.Offset, message.Key, message.Value)
            }
        }(pc)
    }
    wg.Wait()
}
