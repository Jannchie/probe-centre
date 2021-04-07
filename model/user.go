package model

// User is the user of this system.
type User struct {
	BaseModel
	Name   string ``
	Mail   string `gorm:"uniqueIndex"`
	Token  string `gorm:"uniqueIndex"`
	Key    []byte `json:"-"`
	Salt   []byte `json:"-"`
	IP     string `json:"-"`
	Credit uint64 ``
}
