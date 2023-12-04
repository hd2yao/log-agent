package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hd2yao/log-agent/logagent/tailf"
	"time"

	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/client/v3"

	"github.com/hd2yao/log-agent/logagent/common"
)

// etcd

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
func GetConf(key string) (collectEntryList []common.CollectEntry, err error) {
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

// WatchConf 监控 etcd 中日志收集项配置变化的函数
func WatchConf(key string) {
	for {
		watchChan := client.Watch(context.Background(), key)
		var newConf []common.CollectEntry
		for wresp := range watchChan {
			logrus.Info("get new conf from etcd!")
			for _, evt := range wresp.Events {
				fmt.Printf("type:%s key:%s value:%s\n", evt.Type, evt.Kv.Key, evt.Kv.Value)
				err := json.Unmarshal(evt.Kv.Value, &newConf)
				if err != nil {
					logrus.Errorf("json unmarshal new conf failed, err:%v", err)
					continue
				}
				// 告诉tailfile这个模块应该启用新的配置了!
				tailf.SendNewConf(newConf) // 没有任何接收就是阻塞的
			}
		}
	}
}
