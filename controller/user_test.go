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
	r := gin.New()
	r.POST("/test", CreateUser)
	w := httptest.NewRecorder()
	data := bytes.NewReader([]byte("{\"Mail\": \"test@test.com\",\"Password\":\"123456\"}"))
	req, _ := http.NewRequest("POST", "/test", data)
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)

	w = httptest.NewRecorder()
	data = bytes.NewReader([]byte("{\"Mail\": \"test1@test.com\",\"Password\":\"123456\"}"))
	req, _ = http.NewRequest("POST", "/test", data)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}
