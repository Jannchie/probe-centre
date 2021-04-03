package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Jannchie/probe-centre/constant/resp"

	"github.com/Jannchie/probe-centre/service"

	"github.com/Jannchie/probe-centre/constant/code"
	"github.com/Jannchie/probe-centre/db"
	"github.com/Jannchie/probe-centre/model"
	"github.com/Jannchie/probe-centre/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetTask get the task that is need to do and not pending.
func GetTask(c *gin.Context) {
	var task = model.Task{}
	err := GetOneTask(&task)
	if err != nil {
		util.ReturnError(c, err)
		return
	}
	err = updatePend(&task)
	if util.ShouldReturn(c, err) {
		return
	}
	c.JSON(http.StatusOK, task)
}

func updatePend(task *model.Task) error {
	res := db.DB.Model(task).Where("id = ?", task.ID).
		UpdateColumn("pend", time.Now().UTC().Add(time.Second*10))
	err := res.Error
	return err
}

// GetOneTask is the way to get a task that should be done.
func GetOneTask(task *model.Task) error {
	var err error
	if res := db.DB.Where("pend < NOW() AND next < NOW()").
		Limit(1).Find(task); res.Error != nil {
		err = res.Error
	}
	return err
}

// PostRaw update task and insert data data
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
	err = saveRawData(c, form, user)
	if util.ShouldReturn(c, err) {
		return
	}
}

type RawDataForm struct {
	Data   interface{} `form:"Data"`
	TaskID uint64      `form:"TaskID"`
	Number uint64      `form:"Number"`
}

func saveRawData(c *gin.Context, form RawDataForm, user model.User) error {
	_ = c.ShouldBindJSON(&form)
	j, err := json.Marshal(form.Data)
	if err != nil {
		return err
	}
	taskID := form.TaskID
	number := form.Number
	userID := user.ID
	var task = model.Task{}
	if res := db.DB.Where("id = ? AND pend > NOW() AND next < NOW()",
		taskID).
		Take(&task); res.Error != nil {
		return res.Error
	}

	if number != task.SeriesNumber {
		return errors.New("series number not match")
	}
	if res := db.DB.Save(&model.RawData{
		Data:         j,
		TaskID:       taskID,
		UserID:       userID,
		SerialNumber: number,
		URL:          task.URL,
	}); res.Error != nil {
		return res.Error
	}
	db.DB.Where("id = ?", taskID).
		Updates(model.Task{
			SeriesNumber: number + 1,
			Next:         time.Now().UTC().Add(task.Interval * time.Second),
		})
	db.DB.Model(&user).Update("credit", gorm.Expr("credit + 1"))
	return nil
}

func ListTaskStats(c *gin.Context) {
	data, err := service.GetTaskStats()
	if util.ShouldReturn(c, err) {
		return
	}

	c.JSON(code.OK, data)
}

func PostTask(c *gin.Context) {
	var form struct {
		URL      string        `form:"URL"`
		Interval time.Duration `form:"Interval"`
	}
	_ = c.ShouldBind(&form)
	// Min interval =  1 Hour
	user, err := util.GetUserFromCtx(c)
	if form.Interval < time.Hour/time.Second {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code.FAILED,
			"msg":  "Min interval is 1 Hour!",
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
