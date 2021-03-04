package model

// User is the user of this system.
type User struct {
	BaseModel
	Name  string `form:"Name"`
	Mail  string `form:"Mail"`
	Token string `gorm:"index"`
}
