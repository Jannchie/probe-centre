package service

import (
	"log"
	"sync"
	"time"

	"github.com/Jannchie/probe-centre/db"

	"go.uber.org/atomic"

	cmd "github.com/Jannchie/probe-centre/constant"

	"github.com/Jannchie/probe-centre/model"
	"github.com/gorilla/websocket"
)

func GetClientConnectCount() (count int64) {
	db.DB.Model(&model.ClientRecord{}).Count(&count)
	return
}

func WriteJsonWithLock(ws *websocket.Conn, m *sync.Mutex, v interface{}) error {
	m.Lock()
	defer m.Unlock()
	err := ws.WriteJSON(v)
	return err
}

func GetClientRecord(ws *websocket.Conn, user model.User) model.ClientRecord {
	remoteIP := ws.RemoteAddr().String()
	userID := user.ID
	return model.ClientRecord{
		IP:     remoteIP,
		UserID: userID,
	}
}

func StartWebSocket(ws *websocket.Conn, user model.User) {
	clientRecord := GetClientRecord(ws, user)
	_ = db.DB.Create(&clientRecord)
	grand := make(chan struct{})
	defer func() {
		<-grand
		_ = db.DB.Delete(&clientRecord)
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
									_ = WriteJsonWithLock(ws, &m, model.EmptyMsg)
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
