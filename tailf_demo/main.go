package main

import (
    "fmt"
    "time"

    "github.com/hpcloud/tail"
)

func main() {
    fileName := "D:\\Program\\project\\go\\log-agent\\tailf_demo\\my.log"
    config := tail.Config{
        ReOpen:    true,
        Follow:    true,
        Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
        MustExist: false,
        Poll:      true,
    }

    // 打开文件开始读取数据
    tails, err := tail.TailFile(fileName, config)
    if err != nil {
        fmt.Printf("tail %s failed, err:%v\n", fileName, err)
        return
    }

    // 开始读取数据
    for {
        msg, ok := <-tails.Lines // chan tail.Line
        if !ok {
            fmt.Printf("tail file close reopen, fileName:%s\n", tails.Filename)
            time.Sleep(time.Second) // 读取出错等一秒
            continue
        }
        fmt.Println("msg:", msg.Text)
    }
}
