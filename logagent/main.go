package main

// 日志收集的客户端
// 类似的开源项目还有 filebeat
// 收集指定目录下的日志问价，发送到 kafka 中

// 当前已掌握：
// 向 kafka 中发送数据
// 使用 tailf 读日志文件

import (
    "fmt"

    "github.com/sirupsen/logrus"
    "gopkg.in/ini.v1"
)

func main() {
    // 0.读配置文件
    cfg, err := ini.Load("./conf/conf.ini")
    if err != nil {
        logrus.Errorf("load config failed, err:%v", err)
        return
    }

    kafkaAddr := cfg.Section("kafka").Key("address").String()
    fmt.Println(kafkaAddr)

    // 1.初始化

    // 2.根据配置文件中的日志路径，使用 tailf 去收集日志

    // 3.把日志通过 sarama 发送到 kafka

}
