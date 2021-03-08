package api

import (
	"github.com/Jannchie/pyobe-carrier/db"
	"github.com/Jannchie/pyobe-carrier/model"
	"github.com/gin-gonic/gin"
)

// PostStat is the callback function that posting a stat of probe.
func PostStat(c *gin.Context) {
	var stat = model.Stat{}
	if err := c.ShouldBind(&stat); err != nil {
		c.JSON(400, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}

	u, _ := c.Get("user")
	user, ok := u.(model.User)
	if !ok {
		c.JSON(400, gin.H{
			"code": -1,
			"msg":  "Get user info failed.",
		})
		return
	}
	stat.UID = user.ID
	if res := db.DB.Create(&stat); res.Error != nil {
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
