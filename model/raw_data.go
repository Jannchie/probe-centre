package model

import "gorm.io/datatypes"

// RawJSONData is the data Json data.
type RawJSONData struct {
	BaseModel
	URL          string `gorm:"index"`
	UserID       uint64
	TaskID       uint64 `gorm:"index:unique,priority:1"`
	Data         datatypes.JSON
	SerialNumber uint64 `gorm:"index:unique,priority:2"`
}

// RawData is the raw text data.
type RawData struct {
	BaseModel
	URL          string `gorm:"index"`
	UserID       uint64
	TaskID       uint64 `gorm:"index:idx_task_id_number,unique,priority:1"`
	SerialNumber uint64 `gorm:"index:idx_task_id_number,unique,priority:2"`
	Data         []byte
}
