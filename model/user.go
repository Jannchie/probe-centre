package model

// User is the user of this system.
type User struct {
	BaseModel
	Name   string `gorm:"default:0"`
	Mail   string `gorm:"uniqueIndex"`
	Token  string `gorm:"uniqueIndex"`
	Key    []byte `json:"-"`
	Salt   []byte `json:"-"`
	Credit uint64 ``
}
