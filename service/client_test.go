package service

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
	test.Init()
	test.CreateTestUser()
	s := httptest.NewServer(http.HandlerFunc(handle))
	defer s.Close()
	u := "ws" + strings.TrimPrefix(s.URL, "http")
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
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
