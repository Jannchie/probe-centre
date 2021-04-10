package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Jannchie/probe-centre/db"
	"github.com/Jannchie/probe-centre/model"
)

func TestGetTask(t *testing.T) {
	var user model.User
	db.DB.Take(&user)
	var tempTask = model.Task{}
	w := testHandleWithToken(GetTask, "", user.Token)
	assert.Equal(t, 400, w.Code)
	db.DB.Create(&tempTask)
	w = testHandleWithToken(GetTask, "", user.Token)
	assert.Equal(t, 200, w.Code)
	db.DB.Delete(&tempTask)
}

func TestListTaskStats(t *testing.T) {
	w := testHandle(GetTaskStats, "")
	assert.Equal(t, 200, w.Code)
}

func TestPostTask(t *testing.T) {
	w := testHandle(PostTask, `{"URL":"www.baidu.com","Interval":"86400"}`)
	assert.Equal(t, 400, w.Code)
	var user model.User
	db.DB.Take(&user)
	w = testHandleWithToken(PostTask, `{"URL":"www.baidu.com","Interval":3599}`, user.Token)
	assert.Equal(t, 400, w.Code)
	w = testHandleWithToken(PostTask, `{"URL":"www.baidu.com","Interval":86400}`, user.Token)
	assert.Equal(t, 200, w.Code)
	w = testHandleWithToken(PostTask, `{"URL":"www.baidu.com","Interval":86400}`, user.Token)
	assert.Equal(t, 400, w.Code)
	w = testHandle(PostTask, `{"URL":"www.baidu.com","Interval":86400}`)
	assert.Equal(t, 400, w.Code)
}

func TestPostRaw(t *testing.T) {
	var user model.User
	db.DB.Take(&user)
	w := testHandle(PostRaw, "")
	assert.Equal(t, 400, w.Code)
	w = testHandleWithToken(PostRaw, "", user.Token)
	assert.Equal(t, 400, w.Code)
}
