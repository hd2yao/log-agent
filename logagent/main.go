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

    "github.com/hd2yao/log-agent/logagent/common"
    "github.com/hd2yao/log-agent/logagent/etcd"
    "github.com/hd2yao/log-agent/logagent/kafka"
    "github.com/hd2yao/log-agent/logagent/tailf"
)

type Config struct {
    KafkaConfig   `ini:"kafka"`
    EtcdConfig    `ini:"etcd"`
    CollectConfig `ini:"collect"`
}

type KafkaConfig struct {
    Address  string `ini:"address"`
    Topic    string `ini:"topic"`
    ChanSize int64  `ini:"chan_size"`
}

type EtcdConfig struct {
    Address    string `ini:"address"`
    CollectKey string `ini:"collect_key"`
}

type CollectConfig struct {
    LogFilePath string `ini:"logfile_path"`
}

func run() {
    select {}
}

func main() {
    // -1.获取本机 IP，为后续从 etcd 中获取配置文件做准备
    ip, err := common.GetOutboundIP()
    if err != nil {
        logrus.Errorf("get ip failed, err:%v", err)
        return
    }

    var config Config
    // 0.读配置文件:初始化配置文件，加载 kafka 和 collect 的配置项
    if err := ini.MapTo(&config, "./conf/conf.ini"); err != nil {
        logrus.Errorf("load config failed, err: %v", err)
        return
    }
    fmt.Printf("%#v\n", config)

    // 1.连接kafka，初始化 msgChan，起后台 goroutine 去往 kafka 中发送 msg
    if err := kafka.Init([]string{config.KafkaConfig.Address}, config.KafkaConfig.ChanSize); err != nil {
        logrus.Errorf("connect kafka failed, err: %v", err)
        return
    }
    logrus.Info("init kafka success!")

    // 初始化 etcd 连接
    if err := etcd.Init([]string{config.EtcdConfig.Address}); err != nil {
        logrus.Errorf("init etcd failed, err: %v", err)
        return
    }
    // 从 etcd 中拉取要收集日志的配置项
    collectKey := fmt.Sprintf(config.EtcdConfig.CollectKey, ip)
    allConf, err := etcd.GetConf(collectKey)
    if err != nil {
        logrus.Errorf("get conf from etcd failed, err:%v", err)
        return
    }
    fmt.Println(allConf)
    // 派一个小弟去监控 etcd 中 config.EtcdConfig.CollectKey 对应值的变化
    go etcd.WatchConf(collectKey)

    // 2.根据配置文件中的日志路径，使用 tailf 去收集日志
    if err := tailf.Init(allConf); err != nil { // 把从 etcd 中获取的配置项传到 Init
        logrus.Errorf("init tailf config failed, err: %v", err)
        return
    }
    logrus.Info("init tailf success!")

    //// 3.把日志通过 sarama 包装成 kafka.msg 发送到 msgChan 中
    //err = run()
    //if err != nil {
    //	logrus.Errorf("run failed, err: %v", err)
    //	return
    //}
    run()
}
