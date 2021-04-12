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
	j, err := json.Marshal(form.Data)
	var data []byte
	dataType := "STR"
	if err == nil {
		dataType = "JSON"
		data = j
	} else {
		data = []byte(form.Data)
	}
	now := time.Now().UTC()
	taskID := form.TaskID
	number := form.Number
	userID := user.ID
	db.DB.Model(&user).Update("credit", gorm.Expr("credit + 1"))

	var task = model.Task{}
	if res := db.DB.Where("id = ? AND pend > ? AND next < ?",
		taskID, now, now).
		Take(&task); res.Error != nil {
		return res.Error
	}

	if number != task.SeriesNumber {
		return errors.New("series number not match")
	}
	switch dataType {

	case "STR":
		if res := db.DB.Save(&model.RawData{
			Data:         data,
			TaskID:       taskID,
			UserID:       userID,
			SerialNumber: number,
			URL:          task.URL,
		}); res.Error != nil {
			return res.Error
		}
	case "JSON":
		if res := db.DB.Save(&model.RawJSONData{
			Data:         j,
			TaskID:       taskID,
			UserID:       userID,
			SerialNumber: number,
			URL:          task.URL,
		}); res.Error != nil {
			return res.Error
		}
	}
	db.DB.Where("id = ?", taskID).
		Updates(model.Task{
			SeriesNumber: number + 1,
			Next:         time.Now().UTC().Add(task.Interval),
		})
	return nil
}
