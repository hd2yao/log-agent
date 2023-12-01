package etcd

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/sirupsen/logrus"
    "go.etcd.io/etcd/client/v3"
)

// etcd

type CollectEntry struct {
    Path  string `json:"path"`
    Topic string `json:"topic"`
}

var (
    client *clientv3.Client
)

func Init(address []string) (err error) {
    client, err = clientv3.New(clientv3.Config{
        Endpoints:   address,
        DialTimeout: time.Second * 5,
    })
    if err != nil {
        fmt.Printf("connect to etcd failed, err:%v\n", err)
        return
    }
    return
}

// GetConf 拉取日志收集配置项的函数
func GetConf(key string) (collectEntryList []CollectEntry, err error) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
    defer cancel()
    resp, err := client.Get(ctx, key)
    if err != nil {
        logrus.Errorf("get conf from etcd by key:%s failed, err:%v\n", key, err)
        return
    }
    if len(resp.Kvs) == 0 {
        logrus.Warningf("get len:0 conf from etcd by key:%s\n", key)
        return
    }
    ret := resp.Kvs[0]
    // ret.Value // json 格式字符串
    err = json.Unmarshal(ret.Value, &collectEntryList)
    if err != nil {
        logrus.Errorf("json unmarshal failed, err:%v", err)
        return
    }
    return
}
