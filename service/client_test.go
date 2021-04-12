package service

import (
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Jannchie/probe-centre/db"

	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/assert"

	cmd "github.com/Jannchie/probe-centre/constant"
	"github.com/Jannchie/probe-centre/model"
	"github.com/Jannchie/probe-centre/test"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func handle(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()
	user := test.GetTestUser()
	StartWebSocket(c, user)
}
func TestStartWebSocket(t *testing.T) {
	s, ws, err := createWebSocketConnect()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer s.Close()
	defer ws.Close()
	var res model.CentreMsg

	_ = ws.WriteJSON(model.ProbeCmd{
		CMD: cmd.START,
	})
	_ = ws.ReadJSON(&res)
	assert.Equal(t, res.Code, model.StartMsg.Code)

	_ = ws.WriteJSON(model.ProbeCmd{
		CMD: cmd.SWITCH,
	})
	_ = ws.ReadJSON(&res)
	assert.Equal(t, res.Code, model.PausedMsg.Code)

	_ = ws.WriteJSON(model.ProbeCmd{
		CMD: cmd.SWITCH,
	})
	_ = ws.ReadJSON(&res)
	assert.Equal(t, res.Code, model.ResumedMsg.Code)

	exceptTask := test.CreateTestTask()
	task := model.Task{}
	_ = ws.ReadJSON(&task)
	assert.Equal(t, exceptTask.URL, task.URL)

	_ = ws.WriteJSON(model.ProbeCmd{
		CMD: cmd.FINISH,
	})
	_ = ws.ReadJSON(&res)
	assert.Equal(t, res.Code, model.FinishedMsg.Code)

}

func createWebSocketConnect() (*httptest.Server, *websocket.Conn, error) {
	s := httptest.NewServer(http.HandlerFunc(handle))
	u := "ws" + strings.TrimPrefix(s.URL, "http")
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	return s, ws, err
}

func TestGetClientConnectCount(t *testing.T) {

	s, ws, err := createWebSocketConnect()

	// need to wait record created
	time.Sleep(time.Second * 2)

	if err != nil {
		t.Fatalf("%v", err)
	}

	defer s.Close()
	defer ws.Close()

	var res model.CentreMsg
	tests := []struct {
		name      string
		wantCount int64
	}{
		{"With 1 client", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCount := GetClientConnectCount(); gotCount != tt.wantCount {
				t.Errorf("GetClientConnectCount() = %v, want %v", gotCount, tt.wantCount)
			}
		})
	}

	_ = ws.WriteJSON(model.ProbeCmd{
		CMD: cmd.START,
	})

	_ = ws.ReadJSON(&res)
}

func TestGetClientsStat(t *testing.T) {
	test.InitDB()
	var count int64
	db.DB.Model(&model.ClientRecord{}).Count(&count)
	tests := []struct {
		name string
		want gin.H
	}{
		{"0", gin.H{"Count": count}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetClientsStat(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetClientsStat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMain(m *testing.M) {
	test.InitDB()
	test.CreateTestUser()
	code := m.Run()
	os.Exit(code)

}
