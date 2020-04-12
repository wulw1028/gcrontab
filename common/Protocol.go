package common

import "encoding/json"

// 定时任务
type Job struct{
    Name string `json:"name"`
    Command string `json:"command"`
    CronExpr string `json:"cronExpr"`
}

// HTTP接口应答
type Response struct{
    Errno int `json:"errno"`
    Msg string `json:"msg"`
    Data interface{} `json:"data"`
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
