package common

import (
    "encoding/json"
    "fmt"
    "strings"
)

// 定时任务
type Job struct{
    Name string `json:"name"`
    Command string `json:"command"`
    CronExpr string `json:"cronExpr"`
}

// 定时任务添加到cron中，添加Run方法
func (job Job)Run(){
    fmt.Println(job.Name, job.Command)
}

// HTTP接口应答
type Response struct{
    Errno int `json:"errno"`
    Msg string `json:"msg"`
    Data interface{} `json:"data"`
}

// 变化事件
type JobEvent struct {
    EventType int   // SAVE，DELETE
    Job *Job
}

// 应答方法
func BuildResponse(errno int, msg string, data interface{})(resp []byte, err error){
    var(
        respone Response
    )
    
    respone.Errno = errno
    respone.Msg = msg
    respone.Data = data

    resp, err = json.Marshal(respone)

    return
}

// 反序列化Job
func UnpackJob(value []byte)(ret *Job, err error)  {
    ret = &Job{}

    if err = json.Unmarshal(value, &ret);err != nil{
        return
    }

    return
}

// 从etcd的key中提取任务名
func ExtractJobName(jobKey string)(string){
    return strings.TrimPrefix(jobKey, JOB_SAVE_DIR)
}

// 任务变化事件有2种：1）更新任务 2）删除任务
func BuildJobEvent(eventType int, job *Job)(jobEvent *JobEvent)  {
    return &JobEvent{
        EventType: eventType,
        Job: job,
    }
}