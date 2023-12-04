package tailf

import (
	"time"

	"github.com/IBM/sarama"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"

	"github.com/hd2yao/log-agent/logagent/common"
	"github.com/hd2yao/log-agent/logagent/kafka"
)

// tail 相关

var confChan chan []common.CollectEntry

type tailTask struct {
	path    string
	topic   string
	tailObj *tail.Tail
}

func newTailTask(path, topic string) *tailTask {
	return &tailTask{
		path:  path,
		topic: topic,
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
		line, ok := <-t.tailObj.Lines
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

func Init(allConf []common.CollectEntry) (err error) {
	// allConf 中存了若干个日志的收集项
	// 针对每一个日志收集项创建一个对应的 tailObj
	for _, conf := range allConf {
		tt := newTailTask(conf.Path, conf.Topic) // 创建一个日志收集任务
		if err := tt.Init(); err != nil {
			logrus.Errorf("tailf:create tailTask for path:%s failed, err: %v", conf.Path, err)
			continue
		}
		logrus.Infof("create a tail task for path:%s success\n", conf.Path)
		// 去收集日志
		go tt.run()
	}

	// 初始化新配置的 channel
	confChan = make(chan []common.CollectEntry) // 做一个阻塞的 channel
	// 派一个小弟等着新配置来
	newConf := <-confChan // 取到值说明有新的配置
	// 新配置来了之后应该管理一下之前启动的那些 tailTask
	logrus.Infof("get new conf from etcd, conf:%v\n", newConf)
	return
}

func SendNewConf(newConf []common.CollectEntry) {
	confChan <- newConf
}
