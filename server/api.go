package server

import (
	"maa-server/config"
	"maa-server/scheduler"
	"maa-server/utils"
	"net/http"
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
	} else {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":             1,
			"msg": msg,
			"err":              err.Error(),
			"taskCluster":      taskCluster,
		})
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
	} else {
		c.JSON(http.StatusBadRequest, map[string]any{
			"code":             1,
			"msg": msg,
			"err":              err.Error(),
			"content":          content,
		})
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
	} else {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":             1,
			"msg": msg,
			"err":              err.Error(),
			"content":          content,
		})
	}
}

func GetProfiles(c *gin.Context) {
	c.JSON(200, gin.H{"code":0,"data":config.Profiles})
}

func UpdateProfiles(c *gin.Context){
	var data config.ProfilesStruct
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code":1,"err": err.Error()})
		return
	}
	config.Profiles = &data
	config.UpdateProfiles()

	result := utils.CheckGame()
	c.JSON(200, gin.H{"code":0,"data":result})
}
