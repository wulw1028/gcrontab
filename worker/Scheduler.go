package worker

import (
	"fmt"
	"github.com/robfig/cron"
	"github.com/wulw1028/gcrontab/common"
)

// 任务调度
type Scheduler struct {
	jobEventChan chan *common.JobEvent	// etcd任务事件队列
	jobPlanTable map[string]*common.Job	// 任务调度计划表
}

var(
	G_scheduler *Scheduler
	G_cron *cron.Cron
)

// 调度协程
func (scheduler *Scheduler)scheduleLoop(){
	var(
		jobEvent *common.JobEvent
	)
	for{
		select {
		case jobEvent = <- scheduler.jobEventChan:	// 监听任务变化事件
			// 对内存中维护的任务列表做增删改查
			scheduler.handleJobEvent(jobEvent)
			G_cron.Start()
		}
	}
}

// 推送任务变化事件
func (scheduler *Scheduler)PushJobEvent(jobEvent *common.JobEvent)  {
	scheduler.jobEventChan <- jobEvent
}

// 处理任务事件
func (scheduler *Scheduler)handleJobEvent(jobEvent *common.JobEvent){
	var(
		jobExisted bool
	)
	switch jobEvent.EventType {
	case common.JOB_ENENT_SAVE:	// 保存任务事件
		scheduler.jobPlanTable[jobEvent.Job.Name] = jobEvent.Job
		G_cron.AddJob(jobEvent.Job.CronExpr, jobEvent.Job)
		fmt.Println(jobEvent.Job)
	case common.JOB_ENENT_DELETE:	// 删除任务事件
		if _, jobExisted = scheduler.jobPlanTable[jobEvent.Job.Name];jobExisted{
			delete(scheduler.jobPlanTable, jobEvent.Job.Name)
		}
	}
	fmt.Println(scheduler.jobPlanTable)
}

// 尝试执行任务
func (scheduler *Scheduler)TtrStartJob(){

}

// 初始化调度器
func InitScheduler()(err error)  {
	G_scheduler = &Scheduler{
		jobEventChan: make(chan *common.JobEvent, 1000),
		jobPlanTable: make(map[string]*common.Job),
	}
	G_cron = cron.New()

	// 启动调度协程
	go G_scheduler.scheduleLoop()
	return
}