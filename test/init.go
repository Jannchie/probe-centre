package test

import (
	"github.com/Jannchie/probe-centre/db"
	"github.com/Jannchie/probe-centre/model"
	"github.com/Jannchie/probe-centre/service"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Init() {
	db.DB, _ = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	db.InitDB()
}

func CreateTestUser() {
	user := model.User{
		Name:  "test",
		Mail:  "test@test.com",
		Token: "00000000-0000-0000-0000-000000000000",
	}
	passwd := "123456"
	user.Key, user.Salt = service.GenerateKeyAndSalt(passwd)
	db.DB.Create(&user)
}
