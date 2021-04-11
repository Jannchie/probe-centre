package model

type ClientRecord struct {
	BaseModel
	IP     string
	UserID uint64 `gorm:"index"`
}
