package service

import (
	"log"
	"time"

	cmd "github.com/Jannchie/probe-centre/constant"

	"github.com/Jannchie/probe-centre/model"
	"github.com/gorilla/websocket"
)

func StartWebSocket(ws *websocket.Conn, user model.User) {
	grand := make(chan struct{})
	defer func() {
		<-grand
		ws.Close()
	}()
	go func() {
		started := false
		pause := false
		for {
			var probeCmd model.ProbeCmd
			err := ws.ReadJSON(&probeCmd)
			if err != nil {
				log.Println(err)
				return
			}
			log.Println(probeCmd.CMD)
			switch probeCmd.CMD {
			case "":
				return
			case cmd.START:
				go func() {
					if !started {
						_ = ws.WriteJSON(model.StartMsg)
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
								_ = UpdatePend(&task)
								_ = ws.WriteJSON(task)
							}
						}
					}
				}()
			case cmd.SWITCH:
				pause = !pause
				var err error
				if pause {
					err = ws.WriteJSON(model.PausedMsg)
				} else {
					err = ws.WriteJSON(model.ResumedMsg)
				}
				if err != nil {
					log.Println(err)
				}
			case cmd.FINISH:
				err = ws.WriteJSON(model.FinishedMsg)
				if err != nil {
					log.Println(err)
				}
				grand <- struct{}{}
			case cmd.COMMIT:
				if err != nil {
					log.Println(err)
					continue
				}
				err = SaveRawData(probeCmd.Data, user)
				if err != nil {
					log.Println(err)
					continue
				}
			}
		}
	}()
}
