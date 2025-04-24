package server

import (
	"maa-server/config"
	"maa-server/scheduler"
	"maa-server/utils"
	"net/http"
	"path/filepath"
	"sort"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func Test(c *gin.Context) {
	log.Infoln(c.Request.Host)
}

func GetTemplateCluster(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "data": config.Conf.TemplateCluster})
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
	c.JSON(200, gin.H{"code": 0, "data": queue})
}

func ChangeCluster(c *gin.Context) {
	var data scheduler.ApiStruct

	if err := c.BindJSON(&data); err != nil {
		log.Errorln(err)
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg":"参数错误"})
		return
	}
	hasError, str := scheduler.ApiToUpdateCluster(data)
	if hasError {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": str})
		return
	} else {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code": 0,
			"msg":  str,
		})
		return
	}
}

func ChangeTask(c *gin.Context) {
	var data scheduler.ApiStruct

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "参数错误"})
		return
	}

	var msg string

	if scheduler.CheckIsCurrentTask(data) {
		msg = "当前任务正在运行, 该更改将在下次生效"
	} else {
		msg = ""
	}

	taskCluster, err := scheduler.ApiToUpdateTask(data)

	if err == nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code":        0,
			"msg":         msg,
			"taskCluster": taskCluster,
		})
		return
	} else {
		log.Errorln(err)
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":        1,
			"msg":         "参数错误",
			"taskCluster": taskCluster,
		})
		return
	}
}

func GetTaskFile(c *gin.Context) {
	var data scheduler.ApiStruct
	if err := c.BindJSON(&data); err != nil {
		log.Errorln(err)
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "参数错误"})
		return
	}
	var msg string

	if scheduler.CheckIsCurrentTask(data) {
		msg = "当前任务正在运行, 该更改将在下次生效"
	} else {
		msg = ""
	}
	content, err := scheduler.ReadTaskFile(data)
	if err == nil {
		c.JSON(http.StatusOK, map[string]any{
			"code":    0,
			"msg":     msg,
			"content": content,
		})
		return
	} else {
		log.Errorln(err)
		c.JSON(http.StatusBadRequest, map[string]any{
			"code":    1,
			"msg":     "参数错误",
			"content": content,
		})
		return
	}
}

func ChangeTaskFile(c *gin.Context) {
	var data scheduler.ApiStruct
	if err := c.BindJSON(&data); err != nil {
		log.Errorln(err)
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "参数错误"})
		return
	}
	var msg string

	if scheduler.CheckIsCurrentTask(data) {
		msg = "当前任务正在运行, 该更改将在下次生效"
	} else {
		msg = ""
	}
	content, err := scheduler.ModifyTaskFile(data)
	if err == nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code":    0,
			"msg":     msg,
			"content": content,
		})
		return
	} else {
		log.Errorln(err)
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    1,
			"msg":     "参数错误",
			"content": content,
		})
		return
	}
}

func GetProfiles(c *gin.Context) {
	content, err := scheduler.ReadProfile()
	if err == nil {
		c.JSON(http.StatusOK, map[string]any{
			"code":    0,
			"msg":     err,
			"content": content,
		})
		return
	} else {
		log.Errorln(err)
		c.JSON(http.StatusBadRequest, map[string]any{
			"code":    1,
			"msg":     "参数错误",
			"content": content,
		})
		return
	}
}

func UpdateProfile(c *gin.Context) {
	var data map[string]string
	if err := c.BindJSON(&data); err != nil {
		log.Errorln(err)
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg":"参数错误"})
		return
	}
	content, ok := data["content"]
	if !ok {
		c.JSON(400, gin.H{"code": 1, "msg": "content field is missing"})
		return
	}
	var tmpProfile config.ProfilesStruct
	if err := toml.Unmarshal([]byte(content), &tmpProfile); err != nil {
		log.Errorln(err)
		c.JSON(400, gin.H{"code": 1, "msg": "Error decoding TOML content"})
		return
	}

	var msg string

	if scheduler.ScheduleData.CurrentTaskCluster != nil {
		msg = "当前任务正在运行"
	} else {
		msg = ""
	}

	config.Profiles = &tmpProfile
	config.UpdateProfile()

	c.JSON(http.StatusOK, map[string]any{
		"code": 0,
		"msg":  msg,
	})
}

func CheckGame(c *gin.Context) {
	result := utils.IsGameReady()
	if result == ""{
		c.JSON(http.StatusBadRequest, map[string]any{
			"code": 1,
			"msg":  "游戏未准备就绪",
		})
		return
	}else{
		c.JSON(http.StatusOK, map[string]any{
			"code": 0,
			"msg":  result,
		})
		return
	}
}

func GetRunningTask(c *gin.Context) {
	if scheduler.ScheduleData.CurrentTaskCluster == nil {
		c.JSON(http.StatusOK, map[string]any{
			"code":        0,
			"taskCluster": "",
		})
	} else {
		c.JSON(http.StatusOK, map[string]any{
			"code":        0,
			"taskCluster": scheduler.ScheduleData.CurrentTaskCluster.Hash,
		})
	}
}

func ForceStopRunningTask(c *gin.Context) {
	if scheduler.ScheduleData.CurrentTaskCluster != nil {
		scheduler.ScheduleData.MaaCancelFunc()
	}
	c.JSON(http.StatusOK, map[string]any{
		"code": 0,
		"msg":  "任务已停止",
	})
}

func UploadInfrastFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		log.Errorln(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无法获取上传的文件",
			"code":  1,
		})
		return
	}
	uploadPath := filepath.Join(config.D.InfrastDir, "infrast.json") // 保存路径
	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		log.Errorln(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无法保存文件",
			"code":  1,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
	})
}

func UpdateEmailPush(c *gin.Context) {
	var data config.EmailPushStruct
	if err := c.BindJSON(&data); err != nil {
		log.Errorln(err)
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg":"参数错误"})
		return
	}
	if data.Token == "" && data.EmailAddress == "" {
		c.JSON(http.StatusOK, gin.H{"code": 0,"data":config.Conf.EmailPush})
		return
	}
	config.Conf.EmailPush = data
	config.UpdateConfig()
	c.JSON(http.StatusOK, gin.H{"code": 0,"data":config.Conf.EmailPush})
}
