package controller

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientWebSocket(t *testing.T) {
	var resp *http.Response
	token := "00000000-0000-0000-0000-000000000000"
	resp = testWSHandleWithToken(ClientWebSocketHandle, token)
	assert.Equal(t, 101, resp.StatusCode)
	resp = testWSHandleWithToken(ClientWebSocketHandle, "")
	assert.Equal(t, 400, resp.StatusCode)
}
