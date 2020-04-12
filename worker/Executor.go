package worker

import "github.com/wulw1028/gcrontab/common"

// 任务执行器
type Executor struct {
	
}

var(
	G_executor *Executor
)

// 执行一个任务
func (executor *Executor)ExecutorJob(job *common.Job)  {
	
}

// 初始化执行器
func InitExecutor()(err error)  {

}