package collector

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jannchie/probe/common"
	"github.com/jannchie/probe/common/model"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/jannchie/probe/centre/tasker"
)

var router *gin.Engine

func TestMain(m *testing.M) {
	router = gin.Default()
	common.InitDB()
	tasker.Init(router)
	Init(router)
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestGetRawData(t *testing.T) {
	task := model.Task{
		ID:  120,
		URL: "www.jannchie.com",
	}
	common.DB.Save(&task)
	rawData := model.RawData{
		ID:   801,
		Data: []byte("test"),
	}
	common.DB.Save(&rawData)
	var w *httptest.ResponseRecorder
	var req *http.Request

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/raw?path=www.jannchie.c1om&gt=800&count=2", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, "[]", w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/raw?path=www.jannchie.com&gt=801&count=2", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, "[]", w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/raw?path=www.jannchie.com&gt=800&count=2", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Result().StatusCode)

}

func TestPostRawDataHandle(t *testing.T) {
	var tasks []model.Task
	common.DB.Find(&tasks)
	common.DB.Create(&model.Task{
		URL: "www.abc.com",
		ID:  123,
	})
	rawData := gin.H{
		"data":    "xxx",
		"task_id": 123,
		"number":  1,
	}
	byteData, _ := json.Marshal(rawData)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/raw", bytes.NewBuffer(byteData))
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Result().StatusCode)
}
