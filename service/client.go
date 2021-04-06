package service

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Jannchie/probe-centre/model"
	"github.com/gorilla/websocket"
)

func StartWebSocket(ws *websocket.Conn, user model.User) {
	grand := make(chan struct{})
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
			case "":
				return
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
								err = GetOneTask(&task)
								if err != nil {
									continue
								}
								log.Println(fmt.Sprintf("send task %s", task.URL))
								_ = UpdatePend(&task)
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
				data := model.RawDataForm{}
				_ = json.Unmarshal(p, &data)
				err = SaveRawData(data, user)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}()
	<-grand
}
