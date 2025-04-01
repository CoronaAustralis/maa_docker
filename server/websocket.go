package server

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"maa-server/utils"
	"net/http"
	"time"
)

var wsUpgrader = websocket.Upgrader{
	HandshakeTimeout: time.Second * 5,
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var hook *utils.WebSocketHook

func init() {
	hook = utils.NewWebSocketHook()
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.AddHook(hook)

	// go func(){
	// 	for{
	// 		log.Println("for test")
	// 		time.Sleep(time.Second*3)
	// 	}
	// }()
}

func WsHandler(c *gin.Context, hook *utils.WebSocketHook) {

	ws, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Errorln(err)
		return
	}
	defer ws.Close()

	hook.AddClient(ws)
	defer hook.RemoveClient(ws)

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
}
