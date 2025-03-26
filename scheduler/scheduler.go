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
	FinishCallChan     chan int
}

var ScheduleData = ScheduleStruct{
	FinishCallChan: make(chan int),
}

func Schedule() {
	for {
		_, ok := utils.IsGameReady()
		if !ok {
			time.Sleep(time.Second * 5)
			continue
		}
		task := GetTask()
		if task.Hash == "" {
			MaaStopGame()
			time.Sleep(time.Second * 3)
		} else {
			ScheduleData.CurrentTaskCluster = &task
			RunTask(task)
			ScheduleData.CurrentTaskCluster = nil
			// <-ScheduleData.FinishCallChan
		}
	}
}
