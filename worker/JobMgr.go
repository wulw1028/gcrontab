package worker

import (
    "context"
    "github.com/coreos/etcd/clientv3"
    "github.com/coreos/etcd/mvcc/mvccpb"
    "github.com/wulw1028/gcrontab/common"
    "time"
)

var(
    G_jobMgr *JobMgr
)

type JobMgr struct{
    client *clientv3.Client
    kv clientv3.KV
    lease clientv3.Lease
    watcher clientv3.Watcher
}

func InitJobMgr()(err error){
    var(
        kv clientv3.KV
        lease clientv3.Lease
        config clientv3.Config
        client *clientv3.Client
        watcher clientv3.Watcher
    )

    // 初始化配置
    config = clientv3.Config{
        Endpoints: G_config.EtcdEndpoints,
        DialTimeout: time.Duration(G_config.EtcdDialTimeout) * time.Millisecond,
    }

    // 建立连接
    if client, err = clientv3.New(config);err != nil{
        return
    }

    // 得到KV和Lease的API子集
    kv = clientv3.NewKV(client)
    lease = clientv3.NewLease(client)
    watcher = clientv3.NewWatcher(client)

    // 赋值单例
    G_jobMgr = &JobMgr{
        client: client,
        kv: kv,
        lease: lease,
        watcher: watcher,
    }

    // 启动监听
    G_jobMgr.watchJobs()

    return
}

// 监听任务变化
func (jobMgr *JobMgr)watchJobs()(err error)  {
    var(
        getResp *clientv3.GetResponse
        kvPair *mvccpb.KeyValue
        job *common.Job
        watchStartRevision int64
        watchChan clientv3.WatchChan
        watchResp clientv3.WatchResponse
        watchEvent *clientv3.Event
        jobName string
        jobEvent *common.JobEvent
    )

    if getResp,err = jobMgr.kv.Get(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithPrefix());err != nil{
        return
    }

    // 获取当前的任务
    for _, kvPair = range getResp.Kvs{
        if job, err = common.UnpackJob(kvPair.Value);err == nil{
            // TODO：把job同步给scheduler调度协程
            jobEvent = common.BuildJobEvent(common.JOB_ENENT_SAVE, job)
            G_scheduler.PushJobEvent(jobEvent)
        }
    }

    // 从该revision向后监听变化事件
    go func() {
        // 从GET时刻的后续版本开始监听变化
        watchStartRevision = getResp.Header.Revision + 1

        watchChan = jobMgr.watcher.Watch(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithRev(watchStartRevision), clientv3.WithPrefix())
        // 处理监听事件
        for watchResp = range watchChan{
            for _, watchEvent = range watchResp.Events{
                switch watchEvent.Type {
                case mvccpb.PUT:    // 任务保存事件
                    // TODO：反序列化Job，推送一个更新事件给scheduler
                    if job, err = common.UnpackJob(watchEvent.Kv.Value);err != nil{
                        continue
                    }

                    jobEvent = common.BuildJobEvent(common.JOB_ENENT_SAVE, job)

                case mvccpb.DELETE: // 任务被删除
                    // TODO：推送一个删除事件给scheduler
                    jobName= common.ExtractJobName(string(watchEvent.Kv.Key))

                    job = &common.Job{Name: jobName}

                    jobEvent = common.BuildJobEvent(common.JOB_ENENT_DELETE, job)
                }
                // TODO：更新和删除都需要推送给scheduler，提取到后面统一处理
                // G_scheduler PushJobEvent(jobEvent)
                G_scheduler.PushJobEvent(jobEvent)
            }
        }
    }()

    return
}