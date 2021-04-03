package test

import (
	"github.com/Jannchie/probe-centre/db"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Init() {
	db.DB, _ = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	db.InitDB()
}
