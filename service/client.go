package service

import (
	"log"
	"sync"
	"time"

	"go.uber.org/atomic"

	cmd "github.com/Jannchie/probe-centre/constant"

	"github.com/Jannchie/probe-centre/model"
	"github.com/gorilla/websocket"
)

func WriteJsonWithLock(ws *websocket.Conn, m *sync.Mutex, v interface{}) error {
	m.Lock()
	defer m.Unlock()
	err := ws.WriteJSON(v)
	return err
}

func StartWebSocket(ws *websocket.Conn, user model.User) {
	grand := make(chan struct{})
	defer func() {
		<-grand
		ws.Close()
	}()
	m := sync.Mutex{}
	started := false
	pause := atomic.NewBool(false)

	go func() {
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
				if !started {
					_ = WriteJsonWithLock(ws, &m, model.StartMsg)
					started = true
					go func() {
						const duration = time.Second * 10
						ticker := time.NewTicker(duration)
						defer ticker.Stop()
						for {
							<-ticker.C
							if !pause.Load() {
								task := model.Task{}
								err = GetOneTask(&task)
								if err != nil {
									continue
								}
								_ = UpdatePend(&task)
								_ = WriteJsonWithLock(ws, &m, task)
							}
						}
					}()
				}
			case cmd.SWITCH:
				pause.Store(!pause.Load())
				var err error
				if pause.Load() {
					err = WriteJsonWithLock(ws, &m, model.PausedMsg)
				} else {
					err = WriteJsonWithLock(ws, &m, model.ResumedMsg)
				}
				if err != nil {
					log.Println(err)
				}
			case cmd.FINISH:
				err = WriteJsonWithLock(ws, &m, model.FinishedMsg)
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
