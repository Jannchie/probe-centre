package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Jannchie/probe-centre/model"
	"github.com/Jannchie/probe-centre/util"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Message is return msg
type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// TaskChan Save the tasks should be done.
var TaskChan = make(chan model.Task, 100)

// WsHandler is the handler of ws
func WsHandler(c *gin.Context) {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		util.ReturnError(c, err)
		return
	}
	defer ws.Close()
	for {
		//读取ws中的数据
		_, message, err := ws.ReadMessage()
		fmt.Println(string(message))
		if err != nil {
			break
		}
		if string(message) == "need" {
			task := <-TaskChan
			ws.WriteJSON(task)
		}
	}
}

func sendMsg(ws *websocket.Conn, mt int) {
	ticker := time.NewTicker(500 * time.Millisecond)
	done := make(chan bool)
	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			_ = ws.WriteMessage(mt, []byte(fmt.Sprint("Tick at %r", t)))
		}
	}
}
