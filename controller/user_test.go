package controller

import (
	"testing"

	"github.com/Jannchie/probe-centre/db"
	"github.com/Jannchie/probe-centre/model"

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
	var user model.User
	db.DB.Take(&user)
	w := testHandleWithToken(UpdateUser, "", user.Token)
	assert.Equal(t, 400, w.Code)
	w = testHandleWithToken(UpdateUser, `{"Name": "Temp"}`, user.Token)
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
