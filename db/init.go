package db

import (
	"os"

	"github.com/Jannchie/pyobe-carrier/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the instance of the database.
var DB *gorm.DB

// InitDB initializes the mysql database.
func InitDB() {
	// Get Database DSN From System Environment Variable
	dsn := os.Getenv("PG_DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Stat{})
	if err != nil {
		panic(err)
	}
	DB = db
}
