package api

import (
	"net/http"

	"github.com/Jannchie/pyobe-carrier/db"
	"github.com/Jannchie/pyobe-carrier/model"
	"github.com/Jannchie/pyobe-carrier/util"
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

// GetMyProbeList is a function to get the probe that maintained by logged in user.
// It will return the latest stat of each probe.
func GetMyProbeList(c *gin.Context) {
	var err error
	var user model.User
	if user, err = util.GetUserFromCtx(c); err == nil {
		stats := []model.Stat{}
		res := db.DB.Raw(`
		WITH SUMMARY AS
			(SELECT *,
							ROW_NUMBER() OVER(PARTITION BY p.uuid
																ORDER BY p.id DESC) AS rk
			FROM stats p 
			WHERE uid = ?)
		SELECT s.*
		FROM SUMMARY s
		WHERE s.rk = 1
		`, user.ID).Find(&stats)
		if res.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": -1,
				"msg":  res.Error.Error(),
			})
			return
		}
		c.JSON(200, stats)
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"code": -1,
		"msg":  err.Error(),
	})
}
