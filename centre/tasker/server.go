package tasker

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	. "github.com/jannchie/probe/common"
	. "github.com/jannchie/probe/common/model"
	. "github.com/jannchie/probe/common/util"

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
	if ShouldReturn(c, err) {
		return
	}
	var res = TaskStats{}
	DB.Model(&Task{}).Where("url LIKE ?", fmt.Sprintf("%s%%", form.Path)).Count(&res.Sum)
	c.JSON(http.StatusOK, res)
}

var mutex sync.Mutex

type TaskForm struct {
	Path      string `form:"path"`
	Count     int    `form:"count"`
	Collector string `form:"collector"`
}

func GetTaskHandle(c *gin.Context) {
	var form TaskForm
	err := c.ShouldBindQuery(&form)
	if ShouldReturn(c, err) {
		return
	}
	if form.Count <= 1 {
		task, err := GetOneTask(form)
		if ShouldReturn(c, err) {
			return
		}
		c.JSON(http.StatusOK, task)
	} else {
		tasks, err := ListTasks(form)
		if ShouldReturn(c, err) {
			return
		}
		c.JSON(http.StatusOK, tasks)
	}
}

func GetOneTask(form TaskForm) (Task, error) {
	var task Task
	mutex.Lock()
	defer mutex.Unlock()
	now := time.Now().UTC()
	res := DB.Where("url LIKE ? and created_at < ?", fmt.Sprintf("%s%%", form.Path), now).Take(&task)
	if res.Error != nil {
		return task, res.Error
	}
	DB.Model(&task).UpdateColumn("created_at", now.Add(time.Second*20))
	return task, nil
}

func ListTasks(form TaskForm) ([]Task, error) {
	var tasks []Task
	mutex.Lock()
	defer mutex.Unlock()
	res := DB.Where("url LIKE ? AND collector = ?",
		fmt.Sprintf("%s%%", form.Path), form.Collector).Limit(form.Count).
		Find(&tasks)
	if res.Error != nil {
		return nil, res.Error
	}
	updateTasksCreatedAt(tasks)
	return tasks, nil
}

func updateTasksCreatedAt(tasks []Task) {
	now := time.Now().UTC()
	idList := make([]uint64, len(tasks))
	for index, item := range tasks {
		idList[index] = item.ID
	}
	DB.Model(&Task{}).Where("id IN ?", idList).UpdateColumn("created_at", now.Add(time.Second*20))
}

func PostOneTaskHandle(c *gin.Context) {
	var task Task
	err := c.ShouldBindJSON(&task)
	if ShouldReturn(c, err) {
		return
	}
	res := DB.Create(&task)
	if ShouldReturn(c, res.Error) {
		return
	}
	ReturnOK(c)
}

func DeleteOneTaskHandle(c *gin.Context) {
	var task Task
	err := c.ShouldBindJSON(&task)
	if ShouldReturn(c, err) {
		return
	}
	res := DB.Delete(&task, "url = ?", task.URL)
	if ShouldReturn(c, res.Error) {
		return
	}
	ReturnOK(c)
}

func Init(engine *gin.Engine) {
	_ = DB.AutoMigrate(&Task{})
	engine.Group("/task").
		GET("", GetTaskHandle).
		POST("", PostOneTaskHandle).
		DELETE("", DeleteOneTaskHandle).
		GET("/stats", GetTaskStatsHandle)
}
