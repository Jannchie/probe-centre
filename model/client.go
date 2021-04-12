package model

type ClientRecord struct {
	BaseModel
	IP      string
	UserID  uint64 `gorm:"index"`
	Count   uint64
	Success uint64
}
