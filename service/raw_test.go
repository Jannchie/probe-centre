package service

import (
	"testing"

	"github.com/Jannchie/probe-centre/db"

	"github.com/Jannchie/probe-centre/test"

	"github.com/Jannchie/probe-centre/model"
)

func TestSaveRawData(t *testing.T) {
	test.Init()
	test.CreateTestUser()
	var u model.User
	db.DB.Take(&u)
	db.DB.Create(&model.Task{
		URL:      "http://www.baidu.com",
		Interval: 3600,
	})
	type args struct {
		form model.RawDataForm
		user model.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"no user", args{model.RawDataForm{
			Data:   "",
			TaskID: 0,
			Number: 0,
		}, model.User{}}, true},
		{"user", args{model.RawDataForm{
			Data:   "",
			TaskID: 0,
			Number: 0,
		}, u}, true},
		{"ok", args{model.RawDataForm{
			Data:   "",
			TaskID: 1,
			Number: 0,
		}, u}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &model.Task{}
			if tt.name == "ok" {
				_ = GetOneTask(task)
				_ = UpdatePend(task)
			}
			if err := SaveRawData(tt.args.form, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("SaveRawData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
