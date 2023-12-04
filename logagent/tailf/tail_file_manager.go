package tailf

import (
    "github.com/sirupsen/logrus"

    "github.com/hd2yao/log-agent/logagent/common"
)

// tailTask 的管理者

type tailTaskManager struct {
    tailTaskMap      map[string]*tailTask       // 所有的 tailTask 任务
    collectEntryList []common.CollectEntry      // 所有配置项
    confChan         chan []common.CollectEntry // 等待新配置的通道
}

var tailTaskMgr *tailTaskManager

func Init(allConf []common.CollectEntry) (err error) {
    // allConf 中存了若干个日志的收集项
    // 针对每一个日志收集项创建一个对应的 tailObj
    tailTaskMgr = &tailTaskManager{
        tailTaskMap:      make(map[string]*tailTask, 20),
        collectEntryList: allConf,
        confChan:         make(chan []common.CollectEntry), // 做一个阻塞的 channel
    }
    for _, conf := range allConf {
        tt := newTailTask(conf.Path, conf.Topic) // 创建一个日志收集任务
        if err := tt.Init(); err != nil {
            logrus.Errorf("tailf:create tailTask for path:%s failed, err: %v", conf.Path, err)
            continue
        }
        logrus.Infof("create a tail task for path:%s success\n", conf.Path)
        tailTaskMgr.tailTaskMap[tt.path] = tt // 把创建的 tailTask 任务登记在册，方便后续管理
        // 起一个后台的goroutine去收集日志
        go tt.run()
    }

    go tailTaskMgr.watch() // 在后台等新的配置来
    return
}

func (t *tailTaskManager) watch() {
    for {
        // 派一个小弟等着新配置来,
        newConf := <-t.confChan // 取到值说明新的配置来啦
        // 新配置来了之后应该管理一下我之前启动的那些tailTask
        logrus.Infof("get new conf from etcd, conf:%v, start manage tailTask...", newConf)
        for _, conf := range newConf {
            // 1.原来已经存在的任务就不用动
            if t.isExist(conf) {
                continue
            }

            // 2.原来没有的我要新创建一个 tailTask 任务
            tt := newTailTask(conf.Path, conf.Topic)
            if err := tt.Init(); err != nil {
                logrus.Errorf("tailf:create tailTask for path:%s failed, err: %v", conf.Path, err)
                continue
            }
            logrus.Infof("create a tail task for path:%s success\n", conf.Path)
            tailTaskMgr.tailTaskMap[tt.path] = tt // 把创建的 tailTask 任务登记在册，方便后续管理
            // 起一个后台的goroutine去收集日志
            go tt.run()

            // 3.原来有的现在没有的要停掉 tailTask
        }
    }
}

func (t *tailTaskManager) isExist(conf common.CollectEntry) bool {
    _, ok := t.tailTaskMap[conf.Path]
    return ok
}

func SendNewConf(newConf []common.CollectEntry) {
    tailTaskMgr.confChan <- newConf
}
