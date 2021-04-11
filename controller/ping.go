package controller

import (
	"net/http"

	"github.com/Jannchie/probe-centre/constant/code"
	"github.com/Jannchie/probe-centre/constant/msg"
	"github.com/gin-gonic/gin"
)

// Ping is the callback function to check whether the server is online.
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": msg.OK,
		"Code":    code.OK,
	})
}
