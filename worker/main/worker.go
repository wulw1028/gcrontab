package main

import(
    "github.com/wulw1028/gcrontab/worker"
    "runtime"
    "flag"
    "fmt"
)

var(
    confFile string // 配置文件路径
)

func initEnv(){
    runtime.GOMAXPROCS(runtime.NumCPU()-1)
}

// 解析命令行参数
func initArgs(){
    // worker -config ./worker.json
    flag.StringVar(&confFile, "config", "./worker.json", "指定worker配置文件")
    flag.Parse()
}

func main(){
    var(
        err error
    )
    // 初始化命令行参数
    initArgs()

    // 初始化线程
    initEnv()

    // 加载配置
    if err = worker.InitConfig(confFile);err != nil {
        goto ERR
    }

    // 启动调度器
    if err = worker.InitScheduler();err != nil{
        goto ERR
    }

    // 任务管理器
    if err = worker.InitJobMgr();err != nil{
        goto ERR
    }

    for{}

    return

ERR:
    fmt.Println(err)
}
