package master

import (
    "encoding/json"
    "github.com/wulw1028/gcrontab/common"
    "net"
    "net/http"
    "time"
)

var(
    // 单例对象
    G_apiServer *ApiServer
)

// 任务的HTTP接口
type ApiServer struct {
    httpServer *http.Server
}

// 保存任务接口
// POST job={"name":"job1", "command": "echo hello", "cronExpr": "* * * * *"}
func handleJobSave(w http.ResponseWriter, r *http.Request){
    var(
        err error
        postJob string
        job common.Job
        oldJob *common.Job
        bytes []byte
    )
    // 解析POST表单
    if err = r.ParseForm();err != nil{
        goto ERR
    }

    postJob = r.PostForm.Get("job")
    if err = json.Unmarshal([]byte(postJob), &job);err != nil{
        goto ERR
    }

    // 保存到etcd
    if oldJob,err = G_jobMgr.SaveJob(&job);err != nil{
        goto ERR
    }

    if bytes,err = common.BuildResponse(0, "sucess", oldJob);err == nil{
        w.Write(bytes)
    }
    return

ERR:
    if bytes,err = common.BuildResponse(-1, err.Error(), nil);err == nil{
        w.Write(bytes)
    }
}

// 删除任务接口
// POST /job/delete name=job1
func handleJobDelete(w http.ResponseWriter, r *http.Request){
    var(
        err error
        name string
        oldJob *common.Job
        bytes []byte
    )

    if err = r.ParseForm();err != nil{
        goto ERR
    }

    name = r.PostForm.Get("name")

    if oldJob,err = G_jobMgr.DeleteJob(name);err != nil{
        goto ERR
    }

    if bytes,err = common.BuildResponse(0, "sucess", oldJob);err == nil{
        w.Write(bytes)
    }
    return

ERR:
    if bytes,err = common.BuildResponse(-1, err.Error(), nil);err == nil{
        w.Write(bytes)
    }
}

// 查询任务接口
// GET /job/list
func handleJobList(w http.ResponseWriter, r *http.Request){
    var(
        err error
        jobList []*common.Job
        bytes []byte
    )

    if jobList,err = G_jobMgr.ListJobs();err != nil{
        goto ERR
    }

    if bytes,err = common.BuildResponse(0, "sucess", jobList);err == nil{
        w.Write(bytes)
    }
    return

ERR:
    if bytes,err = common.BuildResponse(-1, err.Error(), nil);err == nil{
        w.Write(bytes)
    }
}

// 强制杀死某个任务
// POST /job/kill name=job1
func handleJobKill(w http.ResponseWriter, r *http.Request){
    var(
        name string
        err error
        bytes []byte
    )

    if err = r.ParseForm();err != nil{
        goto ERR
    }

    name = r.PostForm.Get("name")

    // 杀死任务
    if err = G_jobMgr.KillJob(name);err != nil{
        goto ERR
    }

    if bytes,err = common.BuildResponse(0, "sucess", nil);err == nil{
        w.Write(bytes)
    }
    return

ERR:
    if bytes,err = common.BuildResponse(-1, err.Error(), nil);err == nil{
        w.Write(bytes)
    }
}

// 初始化服务
func InitApiServer()(err error){
    var(
        mux *http.ServeMux
        listener net.Listener
        httpServer *http.Server
        staticDir http.Dir
        staticHander http.Handler
    )

    // 配置路由
    mux = http.NewServeMux()
    mux.HandleFunc("/job/save", handleJobSave)
    mux.HandleFunc("/job/delete", handleJobDelete)
    mux.HandleFunc("/job/list", handleJobList)
    mux.HandleFunc("/job/kill", handleJobKill)

    // 静态文件目录
    staticDir = http.Dir(G_config.WebRoot)
    staticHander = http.FileServer(staticDir)
    mux.Handle("/", http.StripPrefix("/", staticHander))


    // 监听服务
    if listener, err = net.Listen("tcp", G_config.ApiPort); err != nil{
        return
    }

    // 创建一个HTTP服务
    httpServer = &http.Server{
        ReadTimeout: time.Duration(G_config.ApiReadTimeout) * time.Millisecond,
        WriteTimeout: time.Duration(G_config.ApiWriteTimeout) * time.Millisecond,
        Handler: mux,
    }

    // 赋值单例
    G_apiServer = &ApiServer{
        httpServer: httpServer,
    }

    // 启动服务端
    go httpServer.Serve(listener)

    return
}
