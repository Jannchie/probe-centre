package controller

import (
	"net/http"
	"time"

	"github.com/Jannchie/probe-centre/constant/resp"

	"github.com/Jannchie/probe-centre/service"

	"github.com/Jannchie/probe-centre/constant/code"
	"github.com/Jannchie/probe-centre/db"
	"github.com/Jannchie/probe-centre/model"
	"github.com/Jannchie/probe-centre/util"
	"github.com/gin-gonic/gin"
)

// GetTask get the task that is need to do and not pending.
func GetTask(c *gin.Context) {
	var task = model.Task{}
	err := service.GetOneTask(&task)
	if err != nil {
		util.ReturnError(c, err)
		return
	}
	err = service.UpdatePend(&task)
	if util.ShouldReturn(c, err) {
		return
	}
	c.JSON(http.StatusOK, task)
}

// PostRaw update task and insert data data
func PostRaw(c *gin.Context) {
	var form model.RawDataForm
	user, err := util.GetUserFromCtx(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Code": code.FAILED,
			"Msg":  err.Error(),
		})
		return
	}
	err = c.ShouldBindJSON(&form)
	if util.ShouldReturn(c, err) {
		return
	}
	err = service.SaveRawData(form, user)
	if util.ShouldReturn(c, err) {
		return
	}
}

func GetTaskStats(c *gin.Context) {
	data, err := service.GetTaskStats()
	if util.ShouldReturn(c, err) {
		return
	}
	c.JSON(code.OK, data)
}

func PostTask(c *gin.Context) {
	var form struct {
		URL      string        `form:"URL"`
		Subject  string        `form:"Subject"`
		Interval time.Duration `form:"Interval"`
	}
	err := c.ShouldBindJSON(&form)
	if util.ShouldReturn(c, err) {
		return
	}
	// Min interval = 1 Hour
	user, err := util.GetUserFromCtx(c)
	if form.Interval < time.Hour {
		c.JSON(http.StatusBadRequest, gin.H{
			"Code": code.FAILED,
			"Msg":  "Min interval is 1 Hour!",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Code": code.FAILED,
			"Msg":  err.Error(),
		})
		return
	}
	if res := db.DB.Model(&model.Task{}).Create(&model.Task{
		URL:          form.URL,
		Interval:     form.Interval,
		UserID:       user.ID,
		SeriesNumber: 0,
	}); res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Code": code.FAILED,
			"Msg":  res.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, resp.OK)
}
