package service

import (
	"github.com/Jannchie/probe-centre/db"
	"github.com/Jannchie/probe-centre/model"
	"github.com/gin-gonic/gin"
)

func GetTaskStats() (gin.H, error) {
	var waitingCount int64
	var finishedCount int64
	var pendingCount int64
	// pending
	res := db.DB.Model(&model.Task{}).
		Where("pend > NOW() AND next < NOW()").
		Count(&pendingCount)
	if res.Error != nil {
		return nil, res.Error
	}
	// waiting
	res = db.DB.Model(&model.Task{}).
		Where("pend < NOW() AND next < NOW()").
		Count(&waitingCount)
	if res.Error != nil {
		return nil, res.Error
	}
	// finished
	res = db.DB.Model(&model.Task{}).Where("next > NOW()").
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
