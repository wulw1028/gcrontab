package main

import (
	"fmt"
	"github.com/robfig/cron"
	"log"
)

type Job1 struct {
	Name	string
}

func (this Job1)Run() {
	fmt.Println(this.Name)
}

type Job2 struct {
	Name 	string
}

func (this Job2)Run() {
	fmt.Println(this.Name)
}

func main()  {
	var(
		c *cron.Cron
		spec string
	)
	c = cron.New()
	spec = "*/5 * * * * ?"
	c.AddFunc(spec, func() {
		log.Println("cron runing")
	})

	c.AddJob(spec, Job1{"wulewei job1"})
	c.AddJob(spec, Job2{"wulewei job2"})

	c.Start()

	defer c.Stop()

	select {}
}