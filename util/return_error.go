package util

import (
	"net/http"

	"github.com/Jannchie/probe-centre/constant/code"
	"github.com/gin-gonic/gin"
)

// ReturnError add error json data
func ReturnError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"Code": code.FAILED,
		"Msg":  err.Error(),
	})
}
