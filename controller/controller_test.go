package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/websocket"

	"github.com/gin-gonic/gin"

	"github.com/Jannchie/probe-centre/db"

	"github.com/Jannchie/probe-centre/service"

	"github.com/Jannchie/probe-centre/model"

	"github.com/Jannchie/probe-centre/test"
)

func TestMain(m *testing.M) {
	test.Init()
	user := model.User{
		Name:  "test",
		Mail:  "test@test.com",
		Token: "00000000-0000-0000-0000-000000000000",
	}
	passwd := "123456"
	user.Key, user.Salt = service.GenerateKeyAndSalt(passwd)
	db.DB.Create(&user)
	os.Exit(m.Run())
}
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
