package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Jannchie/probe-centre/model"
	"github.com/gin-gonic/gin"
)

func GetRes(req *http.Request, fun func(c *gin.Context)) *httptest.ResponseRecorder {
	res := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(res)
	c.Request = req
	fun(c)
	return res
}

func GetResWithUser(req *http.Request, fun func(c *gin.Context), user model.User) *httptest.ResponseRecorder {
	res := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(res)
	c.Set("user", user)
	c.Request = req
	fun(c)
	return res
}

func TestPostTask(t *testing.T) {
	req, _ := http.NewRequest("GET", "test", strings.NewReader(
		`
		"URL": "www.baidu.com",
		"Interval": 0
		`))
	res := GetRes(req, PostTask)
	if res.Code == 200 {
		t.FailNow()
	}
	res = GetResWithUser(req, PostTask, model.User{})
	if res.Code == 200 {
		t.FailNow()
	}
	// req, _ = http.NewRequest("GET", "test", strings.NewReader(
	// 	`
	// 	"URL": "http://api.bilibili.com/x/web-interface/card?mid=1850091",
	// 	"Duration": 86400000
	// 	`))
	// res = GetRes(req, PostTask)
	// if res.Code != 200 {
	// 	t.Fail()
	// }
}
