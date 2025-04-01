package scheduler

import (
	"context"
	"maa-server/config"
	"maa-server/utils"
	"time"
)

type ScheduleStruct struct {
	CurrentTaskCluster *config.TaskCluster
	MaaCancelFunc     context.CancelFunc
	IsStop			bool
	FinishCallChan     chan int
}

var ScheduleData = ScheduleStruct{
	FinishCallChan: make(chan int),
	IsStop: false,
}

func Schedule() {
	for {
		ok := utils.IsDeviceReady()
		if !ok {
			time.Sleep(time.Second * 5)
			continue
		}else{
			if utils.IsGameReady() == ""{
				time.Sleep(time.Second * 5)
				continue
			}
		}
		task := GetTask()
		if task.Hash == "" {
			if(!ScheduleData.IsStop){
				MaaStopGame()
				ScheduleData.IsStop = true
			}
			time.Sleep(time.Second * 3)
		} else {
			ScheduleData.CurrentTaskCluster = &task
			RunTask(task)
			ScheduleData.IsStop = false
			ScheduleData.CurrentTaskCluster = nil
		}
	}
}
