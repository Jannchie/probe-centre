package controller

import (
	"net/http"

	"github.com/Jannchie/probe-centre/service"

	"github.com/Jannchie/probe-centre/util"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	EnableCompression: true,
}

// ClientWebSocketHandle is the handler of ws
func ClientWebSocketHandle(c *gin.Context) {
	user, err := util.GetUserFromCtx(c)
	if util.ShouldReturn(c, err) {
		return
	}
	ws, _ := upGrader.Upgrade(c.Writer, c.Request, nil)
	service.StartWebSocket(ws, user)
}

func ClientsStatHandle(c *gin.Context) {
	res := service.GetClientsStat()
	c.JSON(http.StatusOK, res)
}
