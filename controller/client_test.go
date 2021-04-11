package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/Jannchie/probe-centre/model"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestClientWebSocket(t *testing.T) {
	var resp *http.Response
	var s *httptest.Server
	var conn *websocket.Conn
	token := "00000000-0000-0000-0000-000000000000"
	s, conn, resp, _ = testWSHandleWithToken(ClientWebSocketHandle, token)
	assert.Equal(t, 101, resp.StatusCode)
	_ = conn.WriteJSON(model.ProbeCmd{CMD: "start"})
	conn.Close()
	s.Close()
	s, _, resp, _ = testWSHandleWithToken(ClientWebSocketHandle, "")
	assert.Equal(t, 400, resp.StatusCode)
	s.Close()
}

func TestClientsStatHandle(t *testing.T) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"base", args{ctx},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ClientsStatHandle(ctx)
		})
	}
}
