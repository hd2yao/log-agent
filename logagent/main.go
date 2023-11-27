package main

// 日志收集的客户端
// 类似的开源项目还有 filebeat
// 收集指定目录下的日志问价，发送到 kafka 中

// 当前已掌握：
// 向 kafka 中发送数据
// 使用 tailf 读日志文件

import (
    "fmt"
    "github.com/hd2yao/log-agent/logagent/tailf"
    "time"

    "github.com/sirupsen/logrus"
    "gopkg.in/ini.v1"

    "github.com/hd2yao/log-agent/logagent/kafka"
)

type Config struct {
    KafkaConfig   `ini:"kafka"`
    CollectConfig `ini:"collect"`
}

type KafkaConfig struct {
    Address string `ini:"address"`
    Topic   string `ini:"topic"`
}

type CollectConfig struct {
    LogFilePath string `ini:"logfile_path"`
}

// 真正的业务逻辑
func run() (err error) {
    // TailTask --> log --> Client --> kafka
    for {
        msg, ok := <-tailf.TailTask.Lines
        if !ok {
            logrus.Warnf("tail file close reopen, fileName:%s\n", tailf.TailTask.Filename)
            time.Sleep(time.Second) // 读取出错等一秒
            continue
        }
        fmt.Println("msg:", msg.Text)
    }
    return
}

func main() {
    // 0.读配置文件
    var config Config
    //cfg, err := ini.Load("./conf/conf.ini")
    //if err != nil {
    //    logrus.Errorf("load config failed, err:%v", err)
    //    return
    //}
    //
    //kafkaAddr := cfg.Section("kafka").Key("address").String()
    //fmt.Println(kafkaAddr)
    if err := ini.MapTo(&config, "./conf/conf.ini"); err != nil {
        logrus.Errorf("load config failed, err: %v", err)
        return
    }
    fmt.Printf("%#v\n", config)

    // 1.初始化（连接kafka）
    if err := kafka.Init([]string{config.KafkaConfig.Address}); err != nil {
        logrus.Errorf("connect kafka failed, err: %v", err)
        return
    }
    logrus.Info("init kafka success!")
    // 2.根据配置文件中的日志路径，使用 tailf 去收集日志
    if err := tailf.Init(config.CollectConfig.LogFilePath); err != nil {
        logrus.Errorf("init tailf config failed, err: %v", err)
        return
    }
    logrus.Info("init tailf success!")
    // 3.把日志通过 sarama 发送到 kafka

    err := run()
    if err != nil {
        logrus.Errorf("run failed, err: %v", err)
        return
    }

}
