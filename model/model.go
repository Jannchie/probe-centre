package model

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel is the base settings of the model.
type BaseModel struct {
	ID        uint           `gorm:"primarykey"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
