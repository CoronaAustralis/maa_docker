package scheduler

import (
	"context"
	"log"
	"maa-server/config"
	"maa-server/utils"
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

func StartTask() {

}

func RunTask(task config.TaskCluster) {
	MaaStartGame()
	for _, v := range task.Tasks {
		log.Println("开始任务：", v)
		flag := 0
		for flag < 3 {
			flag++
			isCancel, err := MaaRun(v)
			if isCancel {
				return
			}
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
	// ScheduleData.FinishCallChan <- 1
}

func MaaRun(task string) (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	ScheduleData.MaaCancelFunc = cancel
	cmd := exec.CommandContext(ctx, "maa", "run", "tmp/"+task)
	err := cmd.Run()
	if ctx.Err() == context.Canceled {
		log.Println("Task was canceled")
		return true, nil
	}

	if err != nil {
		log.Printf("Task exited with error: %v\n", err)
		return false, err
	}
	return false, err
}

func MaaStartGame() (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	ScheduleData.MaaCancelFunc = cancel
	cmd := exec.CommandContext(ctx, "maa", "run", "template/start")
	err := cmd.Run()
	if ctx.Err() == context.Canceled {
		log.Println("start game was canceled")
		return true, nil
	}

	if err != nil {
		log.Printf("start game exited with error: %v\n", err)
		return false, err
	}
	return false, err
}

func MaaStopGame() (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	ScheduleData.MaaCancelFunc = cancel
	cmd := exec.CommandContext(ctx, "maa", "run", "template/end")
	err := cmd.Run()
	if ctx.Err() == context.Canceled {
		log.Println("stop game was canceled")
		return true, nil
	}

	if err != nil {
		log.Printf("stop game exited with error: %v\n", err)
		return false, err
	}
	return false, err
}
