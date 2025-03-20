package scheduler

import (
	"io"
	"log"
	"maa-server/config"
	"maa-server/utils"
	"os"
	"path/filepath"
)

var TaskMap = map[string]func(as ApiStruct) (config.TaskCluster, error){
	"AddTask":    AddTask,
	"RenameTask": RenameTask,
	"DeleteTask": DeleteTask,
}

func CheckIsCurrentTask(as ApiStruct) bool {
	if ScheduleData.CurrentTaskCluster != nil && as.NewTaskCluster.Hash == ScheduleData.CurrentTaskCluster.Hash {
		return true
	} else {
		return false
	}
}

func ApiToUpdateTask(as ApiStruct) (config.TaskCluster, error) {
	return TaskMap[as.ApiType](as)
}

func AddTask(as ApiStruct) (config.TaskCluster, error) {
	oldCluster := config.Conf.TaskCluster[as.NewTaskCluster.Hash]
	taskMap := make(map[string]int)
	for _, i := range oldCluster.Tasks {
		taskMap[i] = 0
	}
	for _, i := range as.NewTaskCluster.Tasks {
		if _, exists := taskMap[i]; !exists {
			oldCluster.Tasks = append(oldCluster.Tasks, i)
			taskMap[i] = 0
			err := utils.CopyFile(filepath.Join(config.D.FightDir, "template", i+".toml"), filepath.Join(config.D.FightDir, as.NewTaskCluster.Hash, i+".toml"))
			if err != nil {
				return config.TaskCluster{}, err
			}
		} else {
			tmp := i
			for {
				tmp += "-a"
				if _, exists := taskMap[tmp]; !exists {
					oldCluster.Tasks = append(oldCluster.Tasks, tmp)
					taskMap[tmp] = 0
					err := utils.CopyFile(filepath.Join(config.D.FightDir, "template", i+".toml"), filepath.Join(config.D.FightDir, as.NewTaskCluster.Hash, tmp+".toml"))
					if err != nil {
						return config.TaskCluster{}, err
					}
					break
				}
			}
		}
	}
	config.Conf.TaskCluster[as.NewTaskCluster.Hash] = oldCluster
	config.UpdateConfig()
	return oldCluster, nil
}

func DeleteTask(as ApiStruct) (config.TaskCluster, error) {
	name := as.Content
	for index, ele := range config.Conf.TaskCluster[as.NewTaskCluster.Hash].Tasks {
		if name == ele {
			tmp := config.Conf.TaskCluster[as.NewTaskCluster.Hash]
			utils.PopSlice(&tmp.Tasks, index)
			err := utils.DeleteFileOrDir(filepath.Join(config.D.FightDir, as.NewTaskCluster.Hash, ele+".toml"))
			if err != nil {
				return config.TaskCluster{}, err
			}
			config.Conf.TaskCluster[as.NewTaskCluster.Hash] = tmp
			config.UpdateConfig()
			break
		}
	}
	return config.Conf.TaskCluster[as.NewTaskCluster.Hash], nil
}

func RenameTask(as ApiStruct) (config.TaskCluster, error) {
	newName := as.Content
	oldName := as.NewTaskCluster.Tasks[0]
	tmpTaskCluster := config.Conf.TaskCluster[as.NewTaskCluster.Hash]
	err := utils.RenameFile(filepath.Join(config.D.FightDir, as.NewTaskCluster.Hash, oldName+".toml"), filepath.Join(config.D.FightDir, as.NewTaskCluster.Hash, newName+".toml"))
	if err != nil {
		return config.TaskCluster{}, err
	}
	for i,ele := range tmpTaskCluster.Tasks{
		if(ele == oldName){
			tmpTaskCluster.Tasks[i] = newName
			break
		}
	}
	config.Conf.TaskCluster[as.NewTaskCluster.Hash] = tmpTaskCluster
	config.UpdateConfig()
	return as.NewTaskCluster, nil
}

func ReadTaskFile(as ApiStruct) (string,error) {
	fs, err := os.Open(filepath.Join(config.D.FightDir, as.NewTaskCluster.Hash, as.Content+".toml"))
	if(err != nil){
		return "",err
	}
	defer fs.Close()
	bytes, err := io.ReadAll(fs)
	if(err != nil){
		return "",err
	}
	return string(bytes),nil
}

func ModifyTaskFile(as ApiStruct) (string,error) {
	fs, err := os.OpenFile(filepath.Join(config.D.FightDir, as.NewTaskCluster.Hash, as.NewTaskCluster.Tasks[0]+".toml"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if(err != nil){
		return "",err
	}
	defer fs.Close()
	log.Println()
	if _,err := fs.WriteString(as.Content);err !=nil{
		return "",err
	}
	return "",nil
}

func ReadProfile() (string,error) {
	fs, err := os.Open(filepath.Join(config.D.ProfilesDir, "default.toml"))
	if(err != nil){
		return "",err
	}
	defer fs.Close()
	bytes, err := io.ReadAll(fs)
	if(err != nil){
		return "",err
	}
	return string(bytes),nil
}