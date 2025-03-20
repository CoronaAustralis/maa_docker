package scheduler

import (
	"maa-server/config"
	"time"
)

type ScheduleStruct struct {
	GameReadyFlag	chan bool
	CurrentTaskCluster *config.TaskCluster
	FinishCallChan     chan int
}

var ScheduleData = ScheduleStruct{
	FinishCallChan: make(chan int),
	GameReadyFlag: make(chan bool),
}

func Schedule() {
	<- ScheduleData.GameReadyFlag
	for {
		task := GetTask()
		if task.Hash == "" {
			time.Sleep(time.Second * 3)
		} else {
			ScheduleData.CurrentTaskCluster = &task
			go RunTask(task)
			<-ScheduleData.FinishCallChan
		}
	}
}
