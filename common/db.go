package common

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := os.Getenv("PROBE_DSN")
	DB, _ = gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
