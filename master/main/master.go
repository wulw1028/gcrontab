package main

import(
    "runtime"
    "github.com/wulw1028/go-crontab/master"
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
    // master -config ./master.json
    flag.StringVar(&confFile, "config", "./master.json", "指定master配置文件")
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
    if err = master.InitConfig(confFile);err != nil {
        goto ERR
    }

    // 任务管理器
    if err = master.InitJobMgr();err != nil{
        goto ERR
    }

    // 启动HttpApi服务
    if err = master.InitApiServer(); err != nil {
        goto ERR
    }

    for{}

    return

ERR:
    fmt.Println(err)
}
