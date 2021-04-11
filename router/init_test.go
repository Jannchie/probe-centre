package router

import (
	"net/http/httptest"
	"testing"

	"github.com/Jannchie/probe-centre/db"
	"github.com/Jannchie/probe-centre/model"

	"github.com/Jannchie/probe-centre/test"

	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
)

func TestInit(t *testing.T) {
	Init()
}

func TestAuthRequired(t *testing.T) {
	test.InitDB()
	test.CreateTestUser()

	type args struct {
		c *gin.Context
	}
	ctx1, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx1.Request = httptest.NewRequest("get", "/test", nil)
	ctx2, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx2.Request = httptest.NewRequest("get", "/test", nil)
	ctx2.Request.Header.Set("token", "00000000-0000-0000-0000-000000000000")
	ctx3, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx3.Request = httptest.NewRequest("get", "/test", nil)
	ctx3.Request.Header.Set("token", "00000000-0000-0000-0000-000000000001")
	tests := []struct {
		name string
		args args
	}{
		{"without token", args{ctx1}},
		{"with token", args{ctx2}},
		{"with err token", args{ctx3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AuthRequired(tt.args.c)
			_, exists := tt.args.c.Get("user")
			switch tt.name {
			case "without token":
				assert.Equal(t, false, exists)
			case "with token":
				assert.Equal(t, true, exists)
			case "with err token":
				assert.Equal(t, false, exists)
			}
		})
	}
}

func TestRecordIP(t *testing.T) {

	test.InitDB()
	test.CreateTestUser()

	ctx1, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx1.Request = httptest.NewRequest("GET", "/test", nil)
	ctx2, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx2.Request = httptest.NewRequest("GET", "/test", nil)
	ctx2.Request.Header.Set("token", "00000000-0000-0000-0000-000000000000")

	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		args args
	}{
		{"without token", args{ctx1}},
		{"with token", args{ctx2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RecordIP(tt.args.c)
			switch tt.name {
			case "with token":
				ip := tt.args.c.ClientIP()
				var user model.User
				db.DB.Take(&user)
				assert.Equal(t, ip, user.IP)
			}
		})
	}
}

func TestCors(t *testing.T) {
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	ctx.Request = httptest.NewRequest("GET", "/test", nil)
	Cors(ctx)
	assert.Equal(t, "Origin, X-Requested-With, Content-Type, Accept, Authorization, token", r.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(t, "*", r.Header().Get("Access-Control-Allow-Origin"))
	ctx2, _ := gin.CreateTestContext(r)
	ctx2.Request = httptest.NewRequest("OPTIONS", "/test", nil)
	Cors(ctx2)
}
