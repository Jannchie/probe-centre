package api

import (
	"github.com/Jannchie/pyobe-carrier/db"
	"github.com/Jannchie/pyobe-carrier/model"
	"github.com/gin-gonic/gin"
)

// PostStat is the callback function that posting a stat of probe.
func PostStat(c *gin.Context) {
	var form = model.Stat{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(400, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	if res := db.DB.Create(&form); res.Error != nil {
		c.JSON(400, gin.H{
			"code": -1,
			"msg":  res.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 1,
		"msg":  "success",
	})
}
