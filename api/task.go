package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Jannchie/probe-centre/constant/code"
	"github.com/Jannchie/probe-centre/constant/msg"
	"github.com/Jannchie/probe-centre/constant/resp"
	"github.com/Jannchie/probe-centre/db"
	"github.com/Jannchie/probe-centre/model"
	"github.com/Jannchie/probe-centre/util"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func PostTask(c *gin.Context) {
	var form struct {
		URL      string        `form:"URL"`
		Interval time.Duration `form:"Interval"`
	}
	c.ShouldBind(&form)
	// Min intreval =  1 Hour
	user, err := util.GetUserFromCtx(c)
	if form.Interval < time.Hour/time.Second {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code.FAILED,
			"msg":  "Min intreval is 1 Hour!",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code.FAILED,
			"msg":  "Get user failed!",
		})
		return
	}
	if res := db.DB.Model(&model.Task{}).Create(&model.Task{
		URL:          form.URL,
		Interval:     form.Interval * time.Second,
		UserID:       user.ID,
		SeriesNumber: 0,
	}); res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code.FAILED,
			"msg":  res.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, resp.OK)
}

// GetTask get the task that is need to do and not pending.
func GetTask(c *gin.Context) {
	var task = model.Task{}
	if res := db.DB.Where("pend < ? AND next < ?", time.Now(), time.Now()).
		Take(&task); res.Error != nil {
		return
	}
	db.DB.Model(&task).Where("id = ?", task.ID).
		UpdateColumn("pend", time.Now().Add(time.Second*10))
	c.JSON(http.StatusOK, task)
}

// PostRaw update task and insert raw data
func PostRaw(c *gin.Context) {
	var form struct {
		Data   interface{} `form:"Data"`
		TaskID uint64      `form:"TaskID"`
		Number uint64      `form:"Number"`
	}

	user, err := util.GetUserFromCtx(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code.FAILED,
			"msg":  err.Error(),
		})
		return
	}
	c.ShouldBindJSON(&form)

	var j datatypes.JSON
	j, err = json.Marshal(form.Data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code.FAILED,
			"msg":  err.Error(),
		})
		return
	}
	taskID := form.TaskID
	number := form.Number
	userID := user.ID
	var task = model.Task{}
	if res := db.DB.Where("id = ? AND pend > ? AND next < ?", taskID, time.Now(), time.Now()).Take(&task); res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code.FAILED,
			"msg":  res.Error.Error(),
		})
		return
	}

	if number != task.SeriesNumber {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code.FAILED,
			"msg":  "Series number not match",
		})
	}
	fmt.Println(j)
	if res := db.DB.Save(&model.RawData{
		Data:         j,
		TaskID:       taskID,
		UserID:       userID,
		SerialNumber: number,
		URL:          task.URL,
	}); res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code.FAILED,
			"msg":  res.Error.Error(),
		})
		return
	}
	db.DB.Where("id = ?", taskID).
		Updates(model.Task{
			SeriesNumber: number + 1,
			Next:         task.Next.Add(task.Interval),
		})
	db.DB.Model(&user).Update("credit", gorm.Expr("credit + 1"))
	c.JSON(http.StatusOK, gin.H{
		"code": code.OK,
		"msg":  msg.OK,
	})
}
