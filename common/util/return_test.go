package util

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestReturnOK(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	type args struct {
		ctx *gin.Context
	}
	tests := []struct {
		name string
		args args
	}{
		{"1", args{c}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ReturnOK(tt.args.ctx)
		})
	}
}

func TestShouldReturn(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	type args struct {
		ctx *gin.Context
		err error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"1", args{c, nil}, false},
		{"2", args{c, errors.New("test")}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ShouldReturn(tt.args.ctx, tt.args.err); got != tt.want {
				t.Errorf("ShouldReturn() = %v, want %v", got, tt.want)
			}
		})
	}
}
