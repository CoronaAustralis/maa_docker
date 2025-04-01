package scheduler

import (
	"bufio"
	"context"
	"maa-server/config"
	"maa-server/utils"

	// "os"
	"os/exec"
	"path/filepath"
	"sort"
	"time"

	log "github.com/sirupsen/logrus"
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
	if MaaStartGame() {
		tmp, exists := config.Conf.TaskCluster[task.Hash]
		if !exists {
			return
		}
		tmp.IsEnable = false
		config.Conf.TaskCluster[task.Hash] = tmp
		config.UpdateConfig()
		return
	}

	for _, v := range task.Tasks {
		log.Infoln("开始任务：", v)
		flag := 0
		for flag < 3 {
			flag++
			isCancel, err := MaaRun(v)
			if isCancel {
				return
			}
			if err != nil {
				log.Infoln("任务执行失败：", v)
				log.Infoln("重启游戏。。。")
				if MaaStopGame() {
					return
				}
				if MaaStartGame() {
					return
				}
			} else {
				break
			}
		}
		if flag == 3 {
			log.Errorln("达到最大重试次数")
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

	// 更新任务时间
	tmp, exists := config.Conf.TaskCluster[task.Hash]
	if !exists {
		return
	}
	if task.Time.Equal(tmp.Time) {
		switch task.Type {
		case "day":
			tmp.Time = tmp.Time.Add(24 * time.Hour)
		case "week":
			tmp.Time = tmp.Time.Add(24 * time.Hour * 7)
		case "month":
			tmp.Time = utils.AddOneMonth(tmp.Time)
		case "custom":
			tmp.Time = tmp.Time.Add(24 * time.Hour * 365)
		}
		config.Conf.TaskCluster[task.Hash] = tmp
		config.UpdateConfig()
	}
}

func MaaRun(task string) (bool, error) {
	return ExecuteCommand("tmp/" + task)
}

func MaaStartGame() bool {
	log.Infoln("启动游戏")
	isCancel, _ := ExecuteCommand("template/start")
	return isCancel
}

func MaaStopGame() bool {
	log.Infoln("结束游戏")
	// isCancel, _ := ExecuteCommand("template/end")
	// return isCancel
	utils.StopGame()
	return false
}

func ExecuteCommand(command string) (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // 确保上下文被正确释放

	ScheduleData.MaaCancelFunc = cancel
	cmd := exec.CommandContext(ctx, "maa", "run", command)

	stdoutPipe, _ := cmd.StdoutPipe()
	stderrPipe, _ := cmd.StderrPipe()

	// 启动 goroutine 实时读取标准输出
	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			log.Infof("\x1b[36m[MAA STDOUT]\x1b[0m %s", scanner.Text())
		}
	}()

	// 启动 goroutine 实时读取标准错误
	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			log.Warnf("\x1b[33m[MAA STDERR]\x1b[0m %s", scanner.Text())
		}
	}()
	// 执行命令
	err := cmd.Run()

	// 判断是否被取消
	if ctx.Err() == context.Canceled {
		log.Infof("Command '%s' was canceled\n", command)
		return true, nil
	}

	// 返回错误（如果有）
	if err != nil {
		log.Errorf("Command '%s' exited with error: %v\n", command, err)
		return false, err
	}

	// 正常完成
	return false, nil
}
