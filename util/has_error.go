package util

import (
	"github.com/gin-gonic/gin"
)

func ShouldReturn(c *gin.Context, err error) bool {
	if err != nil {
		ReturnError(c, err)
		return true
	}
	return false
}
