package model

// Stat is the user of this system.
type Stat struct {
	BaseModel
	UUID     string `form:"UUID" binding:"required" gorm:"index"`
	Name     string `form:"Name" binding:"required" gorm:"index"`
	URLCount uint   `form:"UrlCount"`
	ResCount uint   `form:"ResCount"`
	ErrCount uint   `form:"ErrCount"`
	Status   uint8  `form:"Status"`
	UID      uint64 `gorm:"index"`
}
