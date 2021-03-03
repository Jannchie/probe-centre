package api

import "github.com/gin-gonic/gin"

// Ping is the callback function to check whether the server is online.
func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "success",
		"code":    1,
	})
}
