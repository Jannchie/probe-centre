package service

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Jannchie/probe-centre/model"
	"github.com/gorilla/websocket"
)

type probeCmd struct {
	CMD  string `json:"CMD"`
	Data string `json:"Data,omitempty"`
}

type CentreMsg struct {
	Msg string `json:"Msg"`
	URL string `json:"URL,omitempty"`
}

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
			var cmd probeCmd
			err := ws.ReadJSON(&cmd)
			if err != nil {
				log.Println(err)
				return
			}
			log.Println(cmd)
			switch cmd.CMD {
			case "":
				return
			case "start":
				go func() {
					if !started {
						err := ws.WriteJSON(CentreMsg{"start", ""})
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
					err = ws.WriteJSON(CentreMsg{"paused", ""})
				} else {
					err = ws.WriteJSON(CentreMsg{"resume", ""})
				}
				if err != nil {
					log.Println(err)
				}
			case "fin":
				err = ws.WriteJSON(CentreMsg{"finished", ""})
				if err != nil {
					log.Println(err)
				}
				grand <- struct{}{}
			case "commit":
				data := model.RawDataForm{}
				err = json.Unmarshal([]byte(cmd.Data), &data)
				if err != nil {
					log.Println(err)
					continue
				}
				err = SaveRawData(data, user)
				if err != nil {
					log.Println(err)
					continue
				}
			}
		}
	}()
}
