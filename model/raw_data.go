package model

import "gorm.io/datatypes"

type RawData struct {
	BaseModel
	URL          string `gorm:"index"`
	UserID       uint64
	TaskID       uint64 `gorm:"index:idx_task_id_number,unique,priority:1"`
	Data         datatypes.JSON
	SerialNumber uint64 `gorm:"index:idx_task_id_number,unique,priority:2"`
}
