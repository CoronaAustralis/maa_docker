package main

import (
	"maa-server/scheduler"
	"maa-server/server"

	"github.com/gin-gonic/gin"
)

func main() {
	
	r := gin.New()
	server.Init(r)
	go scheduler.Schedule()
	r.Run("0.0.0.0:8080")
}