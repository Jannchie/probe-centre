package model

import "time"

type IPRecord struct {
	ID   uint64    `gorm:"primarykey"`
	IP   string    `gorm:"index"`
	Path string    ``
	Time time.Time `gorm:"index"`
}
