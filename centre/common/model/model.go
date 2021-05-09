package model

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID        uint64    `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"-"`
	URL       string    `gorm:"uniqueIndex:udx_url" json:"url"`
	Collector string    `json:"collector"`
}
type RawData struct {
	ID        uint64         `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Data      string         `json:"data"`
	URL       string         `gorm:"index" json:"url"`
}

type IPRecord struct {
	ID   uint64    `json:"-" gorm:"primarykey"`
	IP   string    `json:"ip" gorm:"index"`
	Path string    `json:"path"`
	Time time.Time `json:"time" gorm:"index"`
}
