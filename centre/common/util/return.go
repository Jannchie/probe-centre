package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ShouldReturnWithCode(ctx *gin.Context, err error, code int) bool {
	if err != nil {
		ctx.JSON(code, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return true
	}
	return false
}

func ShouldReturn(ctx *gin.Context, err error) bool {
	if err != nil {
		errText := err.Error()
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  errText,
		})
		return true
	}
	return false
}

func ReturnOK(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "ok",
	})
}
