package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Jannchie/probe-centre/service"

	"github.com/Jannchie/probe-centre/model"

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

// WsHandler is the handler of ws
func WsHandler(c *gin.Context) {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	user, _ := util.GetUserFromCtx(c)
	grand := make(chan struct{})
	if err != nil {
		util.ReturnError(c, err)
		return
	}
	defer ws.Close()
	go func() {
		started := false
		pause := false
		for {
			mt, p, err := ws.ReadMessage()
			if err != nil {
				log.Println(err)
			}
			msg := string(p)
			switch msg {
			case "start":
				go func() {
					if !started {
						err := ws.WriteMessage(mt, []byte("started"))
						if err != nil {
							log.Println(err)
						}
						started = true
						const duration = time.Second * 10
						ticker := time.NewTicker(duration)
						defer ticker.Stop()
						for {
							<-ticker.C
							if !pause {
								task := model.Task{}
								err = service.GetOneTask(&task)
								if err != nil {
									continue
								}
								log.Println(fmt.Sprintf("send task %s", task.URL))
								_ = service.UpdatePend(&task)
								_ = ws.WriteJSON(task)
							}
						}
					}
				}()
			case "switch":
				pause = !pause
				var err error
				if pause {
					err = ws.WriteMessage(mt, []byte("paused"))
				} else {
					err = ws.WriteMessage(mt, []byte("resume"))
				}
				if err != nil {
					log.Println(err)
				}
			case "fin":
				err := ws.WriteMessage(mt, []byte("finished"))
				if err != nil {
					log.Println(err)
				}
				grand <- struct{}{}
			default:
				data := RawDataForm{}
				_ = json.Unmarshal(p, &data)
				err = saveRawData(data, user)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}()
	<-grand
}
