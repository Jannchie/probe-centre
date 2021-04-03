package util

import (
	"net/http"

	"github.com/Jannchie/probe-centre/constant/code"
	"github.com/gin-gonic/gin"
)

// ReturnError add error json data
func ReturnError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"code": code.FAILED,
		"msg":  err.Error(),
	})
}
