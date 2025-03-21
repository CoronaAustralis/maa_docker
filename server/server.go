package server

import "github.com/gin-gonic/gin"

func Init(e *gin.Engine) {
	e.GET("/",Test)
	
	api := e.Group("/api")
	
	api.GET("/GetClusters", GetTaskCluster)
	api.GET("/GetTemplateCluster", GetTemplateCluster)

	api.POST("/ChangeCluster",ChangeCluster)
	api.POST("/ChangeTask",ChangeTask)

	api.POST("/GetTaskFile",GetTaskFile)
	api.POST("/ChangeTaskFile",ChangeTaskFile)

	api.POST("/GetProfiles",GetProfiles)
	api.POST("/UpdateProfiles",UpdateProfile)
	
	api.GET("/CheckGame",CheckGame)

	api.POST("UpdateClientType",UpdateClientType)
}

