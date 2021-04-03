package controller

import (
	"os"
	"testing"

	"github.com/Jannchie/probe-centre/db"
	"github.com/Jannchie/probe-centre/model"

	"github.com/Jannchie/probe-centre/test"
)

func TestMain(m *testing.M) {
	test.Init()
	user := model.User{
		Name:  "test",
		Mail:  "test@test.com",
		Token: "00000000-0000-0000-0000-000000000000",
	}
	db.DB.Create(&user)
	os.Exit(m.Run())
}