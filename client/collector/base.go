package collector

import (
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/jannchie/probe/centre/collector"
	"github.com/jannchie/probe/centre/tasker"
	"github.com/jannchie/probe/common/util"
)

func RunBaseCollector() {
	s := util.NewSpeedometer()
	for {
		form := tasker.TaskForm{Count: 16}
		tasks, _ := tasker.ListTasks(form)
		if tasks != nil && len(tasks) == 0 {
			time.Sleep(time.Second)
			continue
		}
		for _, task := range tasks {
			taskURL := task.URL
			taskID := task.ID
			go SaveData(taskURL, taskID)
			s.AddCount(1)
		}
	}
}

func SaveData(targetUrl string, taskID uint64) {
	resp, _ := http.Get(targetUrl)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	if resp.StatusCode >= 300 {
		log.Println(errors.New("status code error"))
		return
	}
	data, _ := io.ReadAll(resp.Body)
	err := collector.SaveRawData(data, targetUrl, taskID)
	if err != nil {
		log.Println(err)
	}
}
