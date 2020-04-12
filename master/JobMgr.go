package master

import(
    "github.com/coreos/etcd/clientv3"
    "github.com/coreos/etcd/mvcc/mvccpb"
    "github.com/wulw1028/gcrontab/common"
    "encoding/json"
    "time"
    "context"
)

var(
    G_jobMgr *JobMgr
)

type JobMgr struct{
    client *clientv3.Client
    kv clientv3.KV
    lease clientv3.Lease
}

func InitJobMgr()(err error){
    var(
        kv clientv3.KV
        lease clientv3.Lease
        config clientv3.Config
        client *clientv3.Client
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

    // 赋值单例
    G_jobMgr = &JobMgr{
        client: client,
        kv: kv,
        lease: lease,
    }

    return
}

// 保存任务，保存到etcd /cron/jobs/人物名 -> json
func (jobMgr *JobMgr)SaveJob(job *common.Job)(oldJob *common.Job, err error){
    var(
        jobKey string
        jobValue []byte
        putResp *clientv3.PutResponse
        oldJobObj common.Job
    )

    jobKey = common.JOB_SAVE_DIR + job.Name
    if jobValue, err = json.Marshal(job);err != nil{
        return
    }

    // 保存到etcd
    if putResp, err = jobMgr.kv.Put(context.TODO(), jobKey, string(jobValue), clientv3.WithPrevKV());err != nil{
        return 
    }

    // 如果是跟新，那么返回旧值
    if putResp.PrevKv != nil{
        if err = json.Unmarshal(putResp.PrevKv.Value, &oldJobObj);err != nil{
            err = nil
            return
        }
        oldJob = &oldJobObj
    }
    return
}

// 删除任务
func (jobMgr *JobMgr)DeleteJob(name string)(oldJob *common.Job, err error){
    var(
        jobKey string
        delResp *clientv3.DeleteResponse
        oldJobObj common.Job
    )

    jobKey = common.JOB_SAVE_DIR + name
    
    if delResp,err = jobMgr.kv.Delete(context.TODO(), jobKey, clientv3.WithPrevKV()); err != nil{
        return
    }

    if len(delResp.PrevKvs) != 0{
        if err = json.Unmarshal(delResp.PrevKvs[0].Value, &oldJobObj);err != nil{
            err = nil
            return
        }
        oldJob = &oldJobObj
    }

    return
}

func (jobMgr *JobMgr)ListJobs()(jobList []*common.Job, err error){
    var(
        jobKey string
        getResp *clientv3.GetResponse
        kvPair *mvccpb.KeyValue
        job *common.Job
    )

    jobKey = common.JOB_SAVE_DIR

    if getResp,err = jobMgr.kv.Get(context.TODO(), jobKey, clientv3.WithPrefix()); err != nil{
        return
    }

    jobList = make([]*common.Job, 0)
    for _, kvPair = range getResp.Kvs {
        job = &common.Job{}
        if err = json.Unmarshal(kvPair.Value, job);err != nil{
            err = nil
            continue
        }
        jobList = append(jobList, job)
    }

    return
}

func (jobMgr *JobMgr)KillJob(name string)(err error){
    var(
        killerKey string
        leaseGrantResp *clientv3.LeaseGrantResponse
        leaseId clientv3.LeaseID
    )

    killerKey = common.JOB_KILLER_DIR + name

    // 让worker监听到一次put操作，创建一个租约让其自动过期即可
    if leaseGrantResp, err = jobMgr.lease.Grant(context.TODO(), 1);err != nil{
        return
    }

    // 租约ID
    leaseId = leaseGrantResp.ID

    // 设置killer标记
    if _,err = jobMgr.kv.Put(context.TODO(), killerKey, "", clientv3.WithLease(leaseId));err != nil{
        return
    }

    return
}
