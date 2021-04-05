package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	dataStr := "{\"Mail\": \"test@test.com\",\"Password\":\"123456\"}"
	w := testHandle(CreateUser, dataStr)
	assert.Equal(t, 400, w.Code)
	dataStr = "{\"Mail\": \"test1@test.com\",\"Password\":\"123456\"}"
	w = testHandle(CreateUser, dataStr)
	assert.Equal(t, 200, w.Code)
	dataStr = "{\"Mail\": \"test1@test.com\",\"Password\":\"123456\"}"
	w = testHandle(CreateUser, dataStr)
	assert.Equal(t, 400, w.Code)
	dataStr = "1{\"Mail\": \"test1@test.com\".\"Password\":\"123456\"}"
	w = testHandle(CreateUser, dataStr)
	assert.Equal(t, 400, w.Code)
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
func testHandleWithToken(handle func(c *gin.Context), dataStr string, token string) *httptest.ResponseRecorder {
	r := gin.New()
	r.POST("/test", handle)
	w := httptest.NewRecorder()
	data := bytes.NewReader([]byte(dataStr))
	req, _ := http.NewRequest("POST", "/test", data)
	req.Header.Add("token", token)
	r.ServeHTTP(w, req)
	return w
}

func TestGetUser(t *testing.T) {
	w := testHandle(GetUser, `{"Token":"00000000-0000-0000-0000-000000000000"}`)
	assert.Equal(t, 200, w.Code)
	w = testHandle(GetUser, `{"Token":"00000000-0000-0000-0000-000000000001"}`)
	assert.Equal(t, 400, w.Code)
	w = testHandle(GetUser, `{"ID":1}`)
	assert.Equal(t, 200, w.Code)
	w = testHandle(GetUser, `{"ID":0}`)
	assert.Equal(t, 400, w.Code)
	w = testHandle(GetUser, `{"ID":100}`)
	assert.Equal(t, 400, w.Code)
	w = testHandle(GetUser, `{"other":123}`)
	assert.Equal(t, 400, w.Code)
	w = testHandle(GetUser, `123123`)
	assert.Equal(t, 400, w.Code)
}

func TestUpdateUser(t *testing.T) {
	token := "00000000-0000-0000-0000-000000000000"
	w := testHandleWithToken(UpdateUser, "", token)
	assert.Equal(t, 400, w.Code)
	w = testHandleWithToken(UpdateUser, `{"Name": "Temp"}`, token)
	assert.Equal(t, 200, w.Code)
}

func TestGetMe(t *testing.T) {
	token := "00000000-0000-0000-0000-000000000000"
	w := testHandleWithToken(GetMe, "", token)
	assert.Equal(t, 200, w.Code)
}

func TestRefreshToken(t *testing.T) {
	token := "00000000-0000-0000-0000-000000000000"
	w := testHandleWithToken(RefreshToken, "", token)
	assert.Equal(t, 200, w.Code)
	w = testHandleWithToken(RefreshToken, "", token)
	assert.Equal(t, 400, w.Code)
}

func TestLogin(t *testing.T) {
	w := testHandle(Login, `{"Mail":"test@test.com","Password":"123456"}`)
	assert.Equal(t, 200, w.Code)
	w = testHandle(Login, `{"Mail1":"test@test.com","Password":"123456"}`)
	assert.Equal(t, 400, w.Code)
	w = testHandle(Login, `{"Mail":"test@test.com","Password":"1234567"}`)
	assert.Equal(t, 400, w.Code)
	w = testHandle(Login, `{"Mail":"test@tes1t.com","Password":"123456"}`)
	assert.Equal(t, 400, w.Code)
}
