package scheduler

import (
	"log"
	"maa-server/config"
	"maa-server/utils"
	"path/filepath"
	"sort"
	"time"
)

type ByTime []config.TaskCluster

func (a ByTime) Len() int           { return len(a) }
func (a ByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTime) Less(i, j int) bool { return a[i].Time.Before(a[j].Time) }

func GetTask() config.TaskCluster {
	now := time.Now()
	queue := map[string][]config.TaskCluster{"day": {}, "week": {}, "month": {}, "custom": {}}
	for _, v := range config.Conf.TaskCluster {
		if now.After(v.Time) && v.IsEnable {
			queue[v.Type] = append(queue[v.Type], v)
		} else {
		}
	}
	typePriority := []string{"month", "week", "day", "custom"}
	for _, i := range typePriority {
		if len(queue[i]) > 0 {
			sort.Sort(ByTime(queue[i]))
			tmp := queue[i][0]
			tmp.Tasks = make([]string, len(queue[i][0].Tasks))
			copy(tmp.Tasks, queue[i][0].Tasks)
			if err := utils.DeleteDirSub(filepath.Join(filepath.Join(config.D.FightDir, "tmp"))); err != nil {
				return config.TaskCluster{}
			}
			if err := utils.CopyDir(filepath.Join(filepath.Join(config.D.FightDir, tmp.Hash)), filepath.Join(filepath.Join(config.D.FightDir, "tmp"))); err != nil {
				return config.TaskCluster{}
			}
			return tmp
		}
	}
	return config.TaskCluster{}
}

func RunTask(task config.TaskCluster) {
	log.Println("开始任务：", task.Alias)
	time.Sleep(10 * time.Second)
	log.Println(task.Tasks)
	log.Println("结束任务：", task.Alias)
	tmp := config.Conf.TaskCluster[task.Hash]
	if task.Time.Equal(tmp.Time) {
		if task.Type == "day" {
			tmp.Time = tmp.Time.Add(24 * time.Hour)
			config.Conf.TaskCluster[task.Hash] = tmp
		} else if task.Type == "week" {
			tmp.Time = tmp.Time.Add(24 * time.Hour * 7)
			config.Conf.TaskCluster[task.Hash] = tmp
		} else if task.Type == "month" {
			tmp.Time = utils.AddOneMonth(tmp.Time)
			config.Conf.TaskCluster[task.Hash] = tmp
		} else if task.Type == "custom" {
			tmp.Time = tmp.Time.Add(24 * time.Hour * 365)
			config.Conf.TaskCluster[task.Hash] = tmp
		}
		config.UpdateConfig()
	}
	ScheduleData.CurrentTaskCluster = nil
	ScheduleData.FinishCallChan <- 1
}
