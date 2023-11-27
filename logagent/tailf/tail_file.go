package tailf

import (
    "github.com/hpcloud/tail"
    "github.com/sirupsen/logrus"
)

// tail 相关

var (
    TailTask *tail.Tail
)

func Init(filename string) (err error) {
    cfg := tail.Config{
        ReOpen:    true,
        Follow:    true,
        Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
        MustExist: false,
        Poll:      true,
    }

    // 打开文件开始读取数据
    TailTask, err = tail.TailFile(filename, cfg)
    if err != nil {
        logrus.Errorf("tailf:create tailTask for path:%s failed, err: %v", filename, err)
    }
    return
}
