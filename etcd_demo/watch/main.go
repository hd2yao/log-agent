package main

import (
    "context"
    "fmt"
    "time"

    "go.etcd.io/etcd/client/v3"
)

// watch

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

    // watch
    watchChan := cli.Watch(context.Background(), "sy")

    for watchResp := range watchChan {
        for _, event := range watchResp.Events {
            fmt.Printf("type:%s key:%s value:%s\n", event.Type, event.Kv.Key, event.Kv.Value)
        }
    }
}
