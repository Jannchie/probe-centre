package model

import "time"

// User is the user of this system.
type Task struct {
	BaseModel
	URL          string        `gorm:"uniqueIndex"`
	Last         time.Time     `gorm:"index"`
	Next         time.Time     `gorm:"index"`
	Pend         time.Time     `gorm:"index"`
	Interval     time.Duration ``
	UserID       uint64        `gorm:"index"`
	SeriesNumber uint64
}
