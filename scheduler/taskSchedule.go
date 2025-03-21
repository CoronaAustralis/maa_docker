package scheduler

import (
	"log"
	"maa-server/config"
	"maa-server/utils"
	"os"
	"os/exec"
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
		if now.After(v.Time) && v.IsEnable && len(v.Tasks) > 0 {
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
	MaaStartGame()
	for _, v := range task.Tasks {
		log.Println("开始任务：", v)
		flag := 0
		for flag < 3 {
			flag++
			err := MaaRun(v)
			if err != nil {
				log.Println("任务执行失败：", v)
				log.Println("重试。。。")
				MaaStopGame()
				MaaStartGame()
			} else {
				break
			}
		}
		if flag == 3 {
			log.Println("达到最大重试次数")
			tmp, exists := config.Conf.TaskCluster[task.Hash]
			if !exists {
				return
			}
			tmp.IsEnable = false
			config.Conf.TaskCluster[task.Hash] = tmp
			config.UpdateConfig()
			return
		}
	}
	MaaStopGame()
	tmp, exists := config.Conf.TaskCluster[task.Hash]
	if !exists {
		return
	}
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

func MaaRun(task string) error {
	cmd := exec.Command("maa", "run", "tmp/"+task)

	// 将子进程的输出和错误重定向到当前进程的标准输出和错误
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 启动子进程
	return cmd.Run()
}

func MaaStartGame(){
	cmd := exec.Command("maa", "run", "template/start")

	// 将子进程的输出和错误重定向到当前进程的标准输出和错误
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 启动子进程
	cmd.Run()
}

func MaaStopGame(){
	cmd := exec.Command("maa", "run", "template/end")

	// 将子进程的输出和错误重定向到当前进程的标准输出和错误
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 启动子进程
	cmd.Run()
}
