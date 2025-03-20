package scheduler

import (
	"maa-server/config"
	"maa-server/utils"
	"path/filepath"
)

type ApiStruct struct {
	ApiType        string
	NewTaskCluster config.TaskCluster
	Content        string
}

var ClusterMap = map[string]func(as ApiStruct)error{
	"AddCluster":AddCluster,
	"ModifyCluster":ModifyCluster,
	"DeleteCluster":DeleteCluster,
}

func ApiToUpdateCluster(as ApiStruct)(bool ,string){
	err := ClusterMap[as.ApiType](as)
	if(err != nil){
		return true,err.Error()
	}
	if(ScheduleData.CurrentTaskCluster != nil && as.NewTaskCluster.Hash == ScheduleData.CurrentTaskCluster.Hash ){
			return false, "当前任务正在运行, 该更改将在下次生效"
	}
	return false,""
}

func AddCluster(as ApiStruct) error {
	err:=utils.CreateNestedDirectory(filepath.Join(config.D.FightDir ,as.NewTaskCluster.Hash))
	if(err != nil){
		return err
	}
	config.Conf.TaskCluster[as.NewTaskCluster.Hash] = as.NewTaskCluster
	config.UpdateConfig()
	return nil
}

func ModifyCluster(as ApiStruct) error {
	config.Conf.TaskCluster[as.NewTaskCluster.Hash] = as.NewTaskCluster
	config.UpdateConfig()
	return nil
}

func DeleteCluster(as ApiStruct) error {
	for k := range config.Conf.TaskCluster {
		if k == as.NewTaskCluster.Hash {
			delete(config.Conf.TaskCluster, k)
			err := utils.DeleteFileOrDir(filepath.Join(filepath.Join(config.D.FightDir, as.NewTaskCluster.Hash)))
			if(err != nil){
				return err
			}
			config.UpdateConfig()
			break
		}
	}
	return nil
}

