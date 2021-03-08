package model

// User is the user of this system.
type User struct {
	BaseModel
	Name  string
	Mail  string `gorm:"index"`
	Token string `gorm:"index"`
	Key   []byte `json:"-"`
	Salt  []byte `json:"-"`
}
