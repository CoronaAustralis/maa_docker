package scheduler

import (
	"bufio"
	"context"
	"io"
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
	// 启动游戏
	if startGameFailed := MaaStartGame(); startGameFailed {
		handleTaskFailure(task)
		return
	}

	// 执行任务列表
	for _, v := range task.Tasks {
		log.Infoln("开始任务：", v)

		// 重试任务执行
		if !retryTaskExecution(v, task, 3) {
			log.Errorln("达到最大重试次数，任务失败")
			handleTaskFailure(task)
			return
		}
	}

	// 更新任务时间
	updateTaskTime(task)
}

// 重试任务执行逻辑
func retryTaskExecution(taskName string, task config.TaskCluster, maxRetries int) bool {
	for attempt := 1; attempt <= maxRetries; attempt++ {
		isCancel, err := MaaRun(taskName)
		if isCancel {
			handleTaskFailure(task)
			return false
		}
		if err != nil {
			log.Errorf("任务执行失败：%s，重试 %d/%d 次", taskName, attempt, maxRetries)
			log.Infoln("重启游戏...")
			if restartGameFailed := restartGame(); restartGameFailed {
				handleTaskFailure(task)
				return false
			}
		} else {
			// 任务执行成功
			return true
		}
	}
	return false
}

// 重启游戏
func restartGame() bool {
	if stopGameFailed := MaaStopGame(); stopGameFailed {
		return true
	}
	return MaaStartGame()
}

// 更新任务时间
func updateTaskTime(task config.TaskCluster) {
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

// 处理任务失败逻辑
func handleTaskFailure(task config.TaskCluster) {
	tmp, exists := config.Conf.TaskCluster[task.Hash]
	if !exists {
		return
	}
	tmp.IsEnable = false
	config.Conf.TaskCluster[task.Hash] = tmp
	config.UpdateConfig()
}

// 执行单个任务
func MaaRun(task string) (bool, error) {
	return ExecuteCommand("tmp/" + task)
}

// 启动游戏
func MaaStartGame() bool {
	log.Infoln("启动游戏")
	isCancel, err := ExecuteCommand("template/start")
	if err != nil {
		log.Errorln("启动游戏失败：", err)
		return true
	}
	return isCancel
}

// 停止游戏
func MaaStopGame() bool {
	log.Infoln("结束游戏")
	utils.StopGame()
	return false
}

// 执行命令
func ExecuteCommand(command string) (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ScheduleData.MaaCancelFunc = cancel
	cmd := exec.CommandContext(ctx, "maa", "run", command)

	// 设置标准输出和标准错误
	stdoutPipe, _ := cmd.StdoutPipe()
	stderrPipe, _ := cmd.StderrPipe()

	// 启动 goroutine 实时读取标准输出
	go logCommandOutput(stdoutPipe, "[MAA STDOUT]", "\x1b[36m")
	go logCommandOutput(stderrPipe, "[MAA STDERR]", "\x1b[33m")

	// 执行命令
	err := cmd.Run()

	// 判断是否被取消
	if ctx.Err() == context.Canceled {
		log.Infof("Command '%s' was canceled\n", command)
		return true, nil
	}

	// 返回错误
	if err != nil {
		log.Errorf("Command '%s' exited with error: %v\n", command, err)
		return false, err
	}

	// 正常完成
	return false, nil
}

// 实时日志输出
func logCommandOutput(pipe io.ReadCloser, prefix, color string) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		log.Infof("%s%s %s\x1b[0m", color, prefix, scanner.Text())
	}
}