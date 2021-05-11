package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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
