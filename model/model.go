package model

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel is the base settings of the model.
type BaseModel struct {
	ID        uint64         `gorm:"primarykey"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
