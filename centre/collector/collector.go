package collector

import (
	"io"
	"net/http"
	"time"

	"github.com/jannchie/probe/centre/tasker"
)

func RunBaseCollector() {
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
		}
	}
}

func SaveData(targetUrl string, taskID uint64) {
	resp, _ := http.Get(targetUrl)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	data, _ := io.ReadAll(resp.Body)
	_ = SaveRawData(data, targetUrl, taskID)
}
