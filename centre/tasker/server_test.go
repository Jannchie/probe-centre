package tasker

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/jannchie/probe/common"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

var router *gin.Engine

func TestMain(m *testing.M) {
	router = gin.Default()
	common.InitDB()
	Init(router)
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestGetOneTaskHandle(t *testing.T) {
	var w *httptest.ResponseRecorder
	var req *http.Request

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/task", strings.NewReader(`{"url":"www.test.com","interval":1}`))
	router.ServeHTTP(w, req)
	assert.Equal(t, "{\"code\":1,\"msg\":\"ok\"}", w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/task", strings.NewReader(`{"url":"www.test.com","interval":1}`))
	router.ServeHTTP(w, req)
	assert.Equal(t, `{"code":-1,"msg":"ERROR: duplicate key value violates unique constraint \"udx_url\" (SQLSTATE 23505)"}`, w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/task/stats", nil)
	router.ServeHTTP(w, req)
	res := gin.H{}
	_ = json.Unmarshal(w.Body.Bytes(), &res)
	assert.Equal(t, float64(2), res["sum"])

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/task?path=x", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
	assert.Equal(t, "{\"code\":-1,\"msg\":\"record not found\"}", w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/task?path=w", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	//assert.Equal(t, `{"id":1,"url":"www.test.com","interval":1,"next":"0001-01-01T00:00:00Z","wait":"2021-05-08T16:53:09.231354+09:00","number":1,"raw_data":null}`, w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/task", strings.NewReader(`{"url":"www.test.com"}`))
	router.ServeHTTP(w, req)
	assert.Equal(t, `{"code":1,"msg":"ok"}`, w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/task", strings.NewReader(`{"url":"www.test.com"}`))
	router.ServeHTTP(w, req)
	assert.Equal(t, `{"code":1,"msg":"ok"}`, w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/task/stats", nil)
	router.ServeHTTP(w, req)
	res = gin.H{}
	_ = json.Unmarshal(w.Body.Bytes(), &res)
	assert.Equal(t, float64(1), res["sum"])

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/task", strings.NewReader(`{"url":"www.test.com","interval":1}`))
	router.ServeHTTP(w, req)
	assert.Equal(t, `{"code":1,"msg":"ok"}`, w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/task/stats", nil)
	router.ServeHTTP(w, req)
	res = gin.H{}
	_ = json.Unmarshal(w.Body.Bytes(), &res)
	assert.Equal(t, float64(2), res["sum"])
}

func TestListTaskHandle(t *testing.T) {
	var w *httptest.ResponseRecorder
	var req *http.Request
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/task?count=2", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}
