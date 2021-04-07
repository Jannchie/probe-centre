package service

import (
	"time"

	"github.com/Jannchie/probe-centre/db"
	"github.com/Jannchie/probe-centre/model"
	"github.com/gin-gonic/gin"
)

func GetTaskStats() (gin.H, error) {
	now := time.Now().UTC()
	var waitingCount int64
	var finishedCount int64
	var pendingCount int64
	// pending
	res := db.DB.Model(&model.Task{}).
		Where("pend > ? AND next < ?", now, now).
		Count(&pendingCount)
	if res.Error != nil {
		return nil, res.Error
	}
	// waiting
	res = db.DB.Model(&model.Task{}).
		Where("pend < ? AND next < ?", now, now).
		Count(&waitingCount)
	if res.Error != nil {
		return nil, res.Error
	}
	// finished
	res = db.DB.Model(&model.Task{}).Where("next > ?", now).
		Count(&finishedCount)
	if res.Error != nil {
		return nil, res.Error
	}
	data := gin.H{
		"waiting":  waitingCount,
		"Pending":  pendingCount,
		"finished": finishedCount,
	}
	return data, nil
}

// GetOneTask is the way to get a task that should be done.
func GetOneTask(task *model.Task) error {
	var err error
	now := time.Now().UTC()
	if res := db.DB.Where("pend < ? AND next < ?", now, now).Take(task); res.Error != nil {
		err = res.Error
		return err

	}
	return nil
}

func UpdatePend(task *model.Task) error {
	res := db.DB.Model(task).Where("id = ?", task.ID).
		UpdateColumn("pend", time.Now().UTC().Add(time.Second*10))
	err := res.Error
	return err
}
