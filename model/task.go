package model

import "time"

// Task is the task for client to crawl data.
type Task struct {
	BaseModel
	URL          string        `gorm:"uniqueIndex"`
	Last         time.Time     `gorm:"index"`
	Next         time.Time     `gorm:"index"`
	Pend         time.Time     `gorm:"index"`
	Interval     time.Duration ``
	UserID       uint64        `gorm:"index"`
	SeriesNumber uint64        ``
}
