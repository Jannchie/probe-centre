package repository

import "gorm.io/gorm"

var (
	User UserReposity
)

func Init(db *gorm.DB) {
	User = UserReposity{DB: db}
}
