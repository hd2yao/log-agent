package tailf

import (
    "context"
    "time"

    "github.com/IBM/sarama"
    "github.com/hpcloud/tail"
    "github.com/sirupsen/logrus"

    "github.com/hd2yao/log-agent/logagent/kafka"
)

// tail 相关

type tailTask struct {
    path    string
    topic   string
    tailObj *tail.Tail
    ctx     context.Context
    cancel  context.CancelFunc
}

func newTailTask(path, topic string) *tailTask {
    ctx, cancel := context.WithCancel(context.Background())
    return &tailTask{
        path:   path,
        topic:  topic,
        ctx:    ctx,
        cancel: cancel,
    }
}

func (t *tailTask) Init() (err error) {
    cfg := tail.Config{
        ReOpen:    true,
        Follow:    true,
        Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
        MustExist: false,
        Poll:      true,
    }
    t.tailObj, err = tail.TailFile(t.path, cfg)
    return err
}

func (t *tailTask) run() {
    // 读取日志，发送到 kafka
    logrus.Infof("collect for path:%s is running...", t.path)
    for {
        select {
        case <-t.ctx.Done(): // 当 cancel() 被调用时，此处就会收到值
            logrus.Infof("path:%s is stopping...", t.path)
            return
        case line, ok := <-t.tailObj.Lines:
            if !ok {
                logrus.Warnf("tail file close reopen, fileName:%s\n", t.path)
                time.Sleep(time.Second) // 读取出错等一秒
                continue
            }
            // 利用通道将同步的代码改为异步的
            // 把读出来的一行日志包装成 kafka 中的 msg 类型
            msg := &sarama.ProducerMessage{}
            msg.Topic = t.topic
            msg.Value = sarama.StringEncoder(line.Text)
            // 放入通道
            kafka.ToMsgChan(msg)
        }
    }
}
