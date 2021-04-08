package service

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/Jannchie/probe-centre/db"
	"github.com/Jannchie/probe-centre/model"
	"gorm.io/gorm"
)

func SaveRawData(form model.RawDataForm, user model.User) error {
	now := time.Now().UTC()
	j, err := json.Marshal(form.Data)
	if err != nil {
		return err
	}
	taskID := form.TaskID
	number := form.Number
	userID := user.ID
	var task = model.Task{}
	if res := db.DB.Where("id = ? AND pend > ? AND next < ?",
		taskID, now, now).
		Take(&task); res.Error != nil {
		return res.Error
	}

	if number != task.SeriesNumber {
		return errors.New("series number not match")
	}
	if res := db.DB.Save(&model.RawJSONData{
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
