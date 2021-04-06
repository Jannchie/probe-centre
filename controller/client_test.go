package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientWebSocket(t *testing.T) {
	var resp *http.Response
	var s *httptest.Server
	//var conn *websocket.Conn
	//token := "00000000-0000-0000-0000-000000000000"
	//s, conn, resp, _ = testWSHandleWithToken(ClientWebSocketHandle, token)
	//assert.Equal(t, 101, resp.StatusCode)
	//_ = conn.WriteMessage(1, []byte("start"))
	//conn.Close()
	//s.Close()
	s, _, resp, _ = testWSHandleWithToken(ClientWebSocketHandle, "")
	assert.Equal(t, 400, resp.StatusCode)
	s.Close()
}
