package service

import (
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

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
	if err != nil {
		log.Println(err)
	}
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

type saveClientRecord struct {
	v model.ClientRecord
	m sync.Mutex
}

func (s *saveClientRecord) Get() model.ClientRecord {
	s.m.Lock()
	defer s.m.Unlock()
	return s.v
}
func (s *saveClientRecord) Set(newVal model.ClientRecord) {
	s.m.Lock()
	defer s.m.Unlock()
	s.v = newVal
}

func (s *saveClientRecord) Save() {
	s.m.Lock()
	defer s.m.Unlock()
	val := s.v
	db.DB.Save(&val)
}

func (s *saveClientRecord) Delete() {
	s.m.Lock()
	defer s.m.Unlock()
	_ = db.DB.Delete(&s.v)
}

func StartWebSocket(ws *websocket.Conn, user model.User) {
	clientRecord := GetClientRecord(ws, user)
	_ = db.DB.Create(&clientRecord)
	scr := saveClientRecord{}
	scr.Set(clientRecord)
	go func() {
		for range time.Tick(time.Second) {
			scr.Save()
		}
	}()
	grand := make(chan struct{})
	defer func() {
		<-grand
		scr.Delete()
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
						const duration = time.Second * 5
						ticker := time.NewTicker(duration)
						defer ticker.Stop()
						for {
							<-ticker.C
							if !pause.Load() {
								task := model.Task{}
								err = GetOneTask(&task)
								if err != nil {
									_ = WriteJsonWithLock(ws, &m, model.EmptyMsg)
									continue
								}
								_ = UpdatePend(&task)
								_ = WriteJsonWithLock(ws, &m, task)
								clientRecord.Count += 1
								scr.Set(clientRecord)
							}
						}
					}()
				}
			case cmd.SWITCH:
				pause.Store(!pause.Load())
				if pause.Load() {
					_ = WriteJsonWithLock(ws, &m, model.PausedMsg)
				} else {
					_ = WriteJsonWithLock(ws, &m, model.ResumedMsg)
				}
			case cmd.FINISH:
				_ = WriteJsonWithLock(ws, &m, model.FinishedMsg)
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
				clientRecord.Success += 1
				scr.Set(clientRecord)
			}
		}
	}()
}

func GetClientsStat() gin.H {
	return gin.H{
		"Count": GetClientConnectCount(),
	}
}
