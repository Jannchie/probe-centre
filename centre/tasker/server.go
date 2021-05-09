package tasker

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/jannchie/probe/centre/common"
	"github.com/jannchie/probe/centre/common/model"
	"github.com/jannchie/probe/centre/common/util"

	"github.com/gin-gonic/gin"
)

type TaskStats struct {
	Sum int64 `json:"sum"`
}

func GetTaskStatsHandle(c *gin.Context) {
	var form struct {
		Path string `form:"path"`
	}
	err := c.ShouldBindQuery(&form)
	if util.ShouldReturn(c, err) {
		return
	}
	var res = TaskStats{}
	common.DB.Model(&model.Task{}).Where("url LIKE ?", fmt.Sprintf("%s%%", form.Path)).Count(&res.Sum)
	c.JSON(http.StatusOK, res)
}

var mutex sync.Mutex

type TaskForm struct {
	Path  string `form:"path"`
	Count int    `form:"count"`
}

func GetTaskHandle(c *gin.Context) {
	var form TaskForm
	err := c.ShouldBindQuery(&form)
	if util.ShouldReturn(c, err) {
		return
	}
	if form.Count <= 1 {
		getOneTask(c, form)
	} else {
		listTasks(c, form)
	}
}

func getOneTask(c *gin.Context, form struct {
	Path  string `form:"path"`
	Count int    `form:"count"`
}) {
	err := c.ShouldBindQuery(&form)
	if util.ShouldReturn(c, err) {
		return
	}
	mutex.Lock()
	defer mutex.Unlock()
	var task model.Task
	now := time.Now().UTC()
	res := common.DB.Where("url LIKE ? and created_at < ?", fmt.Sprintf("%s%%", form.Path), now).Take(&task)
	if util.ShouldReturnWithCode(c, res.Error, http.StatusNotFound) {
		return
	}
	common.DB.Model(&task).UpdateColumn("created_at", now.Add(time.Second*20))
	c.JSON(http.StatusOK, task)
}

func listTasks(c *gin.Context, form TaskForm) {
	var tasks []model.Task
	res := common.DB.Where("url LIKE ?",
		fmt.Sprintf("%s%%", form.Path)).Offset(form.Count).
		Find(&tasks)
	if util.ShouldReturnWithCode(c, res.Error, http.StatusNotFound) {
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func PostOneTaskHandle(c *gin.Context) {
	var task model.Task
	err := c.ShouldBindJSON(&task)
	if util.ShouldReturn(c, err) {
		return
	}
	res := common.DB.Create(&task)
	if util.ShouldReturn(c, res.Error) {
		return
	}
	util.ReturnOK(c)
}

func DeleteOneTaskHandle(c *gin.Context) {
	var task model.Task
	err := c.ShouldBindJSON(&task)
	if util.ShouldReturn(c, err) {
		return
	}
	res := common.DB.Delete(&task, "url = ?", task.URL)
	if util.ShouldReturn(c, res.Error) {
		return
	}
	util.ReturnOK(c)
}

func Init(engine *gin.Engine) {
	_ = common.DB.AutoMigrate(&model.Task{})
	engine.Group("/task").
		GET("", GetTaskHandle).
		POST("", PostOneTaskHandle).
		DELETE("", DeleteOneTaskHandle).
		GET("/stats", GetTaskStatsHandle)
}
