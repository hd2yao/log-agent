package main

import (
    "context"
    "fmt"
    "time"

    "go.etcd.io/etcd/client/v3"
)

// 代码连接 etcd

func main() {
    cli, err := clientv3.New(clientv3.Config{
        Endpoints:   []string{"127.0.0.1:2379"},
        DialTimeout: time.Second * 5,
    })
    if err != nil {
        fmt.Printf("connect to etcd failed, err:%v", err)
        return
    }
    defer cli.Close()

    // put
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    valueStr := `[{"path":"D:/Program/project/go/log-agent/logs/sy.log","topic":"sy"},{"path":"D:/Program/project/go/log-agent/logs/web_log.log","topic":"web_log"}]`
    _, err = cli.Put(ctx, "collect_log_conf", valueStr)
    if err != nil {
        fmt.Printf("put to etcd failed, err:%v", err)
        return
    }
    cancel()

    // get
    ctx, cancel = context.WithTimeout(context.Background(), time.Second)
    gr, err := cli.Get(ctx, "collect_log_conf")
    if err != nil {
        fmt.Printf("get from etcd failed, err:%v", err)
        return
    }
    for _, ev := range gr.Kvs {
        fmt.Printf("key:%s value:%s\n", ev.Key, ev.Value)
    }
    cancel()

}
