package db

import (
	"github.com/Jannchie/probe-centre/model"
	"github.com/Jannchie/probe-centre/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the instance of the database.
var (
	DB *gorm.DB
)

// Init initializes the mysql database.
func Init(dsn string) {
	// Get Database DSN From System Environment Variable
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	DB = db
	if err != nil {
		panic(err)
	}
	InitDB()
}

func InitDB() {
	_ = DB.AutoMigrate(
		&model.Task{},
		&model.User{},
		&model.RawJSONData{},
		&model.IPRecord{},
	)
	repository.Init(DB)
}
