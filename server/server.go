package server

import (
	"github.com/gin-gonic/gin"
)

func Init(e *gin.Engine) {
	api := e.Group("/api")
	
	api.GET("/GetClusters", GetTaskCluster)
	api.GET("/GetTemplateCluster", GetTemplateCluster)

	api.POST("/ChangeCluster",ChangeCluster)
	api.POST("/ChangeTask",ChangeTask)

	api.POST("/GetTaskFile",GetTaskFile)
	api.POST("/ChangeTaskFile",ChangeTaskFile)

	api.GET("/GetProfiles",GetProfiles)
	api.POST("/UpdateProfiles",UpdateProfile)
	
	api.POST("/UploadInfrastFile",UploadInfrastFile)

	api.GET("/CheckGame",CheckGame)

	api.GET("/GetRunningTask",GetRunningTask)
	api.GET("/ForceStopRunningTask",ForceStopRunningTask)

	e.GET("/ws", func(c *gin.Context) {
		WsHandler(c, hook)
	})

	e.StaticFile("/", "./dist/index.html")

	e.Static("/assets", "./dist/assets")
	e.StaticFile("/favicon.ico", "./dist/favicon.ico")
}

