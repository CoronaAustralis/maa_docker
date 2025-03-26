package server

import (
	"maa-server/config"
	"maa-server/scheduler"
	"maa-server/utils"
	"net/http"
	"path/filepath"
	"sort"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func Test(c *gin.Context) {
	log.Println(c.Request.Host)
}

func GetTemplateCluster(c *gin.Context) {
	c.JSON(200, gin.H{"code":0,"data":config.Conf.TemplateCluster})
}

func GetTaskCluster(c *gin.Context) {
	queue := map[string][]config.TaskCluster{"day": {}, "week": {}, "month": {}, "custom": {}}
	for _, v := range config.Conf.TaskCluster {
		queue[v.Type] = append(queue[v.Type], v)
	}
	typePriority := []string{"month", "week", "day", "custom"}
	for _, i := range typePriority {
		if len(queue[i]) > 0 {
			sort.Sort(scheduler.ByTime(queue[i]))
		}
	}
	c.JSON(200, gin.H{"code":0,"data":queue})
}

func ChangeCluster(c *gin.Context) {
	var data scheduler.ApiStruct

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code":1,"err": err.Error()})
		return
	}
	hasError, str := scheduler.ApiToUpdateCluster(data)
	if hasError {
		c.JSON(http.StatusBadRequest, gin.H{"code":1,"err": str})
		return
	} else {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code":   0,
			"msg": str,
		})
		return
	}
}

func ChangeTask(c *gin.Context) {
	var data scheduler.ApiStruct

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code":1,"err": err.Error()})
		return
	}

	var msg string;

	if scheduler.CheckIsCurrentTask(data){
		msg = "当前任务正在运行, 该更改将在下次生效"
	}else{
		msg = ""
	}

	taskCluster, err := scheduler.ApiToUpdateTask(data)

	if err == nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code":             0,
			"msg": msg,
			"err":              err,
			"taskCluster":      taskCluster,
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":             1,
			"msg": msg,
			"err":              err.Error(),
			"taskCluster":      taskCluster,
		})
		return
	}
}

func GetTaskFile(c *gin.Context) {
	var data scheduler.ApiStruct
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code":1,"err": err.Error()})
		return
	}
	var msg string;

	if scheduler.CheckIsCurrentTask(data){
		msg = "当前任务正在运行, 该更改将在下次生效"
	}else{
		msg = ""
	}
	content, err := scheduler.ReadTaskFile(data)
	if err == nil {
		c.JSON(http.StatusOK, map[string]any{
			"code":             0,
			"msg": msg,
			"err":              err,
			"content":          content,
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, map[string]any{
			"code":             1,
			"msg": msg,
			"err":              err.Error(),
			"content":          content,
		})
		return
	}
}

func  ChangeTaskFile(c *gin.Context) {
	var data scheduler.ApiStruct
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code":1,"err": err.Error()})
		return
	}
	var msg string;

	if scheduler.CheckIsCurrentTask(data){
		msg = "当前任务正在运行, 该更改将在下次生效"
	}else{
		msg = ""
	}
	content, err := scheduler.ModifyTaskFile(data)
	if err == nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code":             0,
			"msg": msg,
			"err":              err,
			"content":          content,
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":             1,
			"msg": msg,
			"err":              err.Error(),
			"content":          content,
		})
		return
	}
}

func GetProfiles(c *gin.Context) {
	content, err := scheduler.ReadProfile()
	if err == nil {
		c.JSON(http.StatusOK, map[string]any{
			"code":             0,
			"err":              err,
			"content":          content,
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, map[string]any{
			"code":             1,
			"err":              err.Error(),
			"content":          content,
		})
		return
	}
}

func UpdateProfile(c *gin.Context){
	var data config.ProfilesStruct
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code":1,"err": err.Error()})
		return
	}

	var msg string;

	if scheduler.ScheduleData.CurrentTaskCluster != nil {
		msg = "当前任务正在运行"
	}else{
		msg = ""
	}
	
	config.Profiles = &data
	config.UpdateProfile()

	c.JSON(http.StatusOK, map[string]any{
		"code":             0,
		"msg": msg,
	})
	return
}

func CheckGame(c *gin.Context){
	res,flag := utils.IsGameReady()
	if flag {
		c.JSON(http.StatusOK, map[string]any{
			"code":             0,
			"msg": res,
			"err":res,
		})
		return
	}else{
		c.JSON(http.StatusBadRequest, map[string]any{
			"code":             1,
			"msg": res,
			"err":res,
		})
		return
	}
}

func UploadInfrastFile(c *gin.Context){
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无法获取上传的文件",
			"code":1,
		})
		return
	}
	uploadPath := filepath.Join(config.D.InfrastDir, "infrast.json") // 保存路径
	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无法保存文件",
			"code":1,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":0,
	})
}