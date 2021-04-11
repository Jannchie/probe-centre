package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Jannchie/probe-centre/test"

	"github.com/gorilla/websocket"

	"github.com/gin-gonic/gin"
)

func testHandle(handle func(c *gin.Context), dataStr string) *httptest.ResponseRecorder {
	r := gin.New()
	r.POST("/test", handle)
	w := httptest.NewRecorder()
	data := bytes.NewReader([]byte(dataStr))
	req, _ := http.NewRequest("POST", "/test", data)
	r.ServeHTTP(w, req)
	return w
}
func testHandleWithToken(handle func(c *gin.Context), dataStr string, token string) (w *httptest.ResponseRecorder) {
	r := gin.New()
	r.POST("/test", handle)
	w = httptest.NewRecorder()
	data := bytes.NewReader([]byte(dataStr))
	req, _ := http.NewRequest("POST", "/test", data)
	req.Header.Add("token", token)
	r.ServeHTTP(w, req)
	return
}

func testWSHandleWithToken(handle func(c *gin.Context), token string) (s *httptest.Server, conn *websocket.Conn, resp *http.Response, err error) {
	h := func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("token", token)
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = r
		handle(ctx)
	}
	s = httptest.NewServer(http.HandlerFunc(h))
	d := websocket.Dialer{}
	conn, resp, err = d.Dial("ws://"+s.Listener.Addr().String(), nil)
	return
}

func TestMain(m *testing.M) {
	test.InitDB()
	test.CreateTestUser()
	os.Exit(m.Run())
}
