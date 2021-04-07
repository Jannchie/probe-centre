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
			Data:   nil,
			TaskID: 0,
			Number: 0,
		}, model.User{}}, true},
		{"user", args{model.RawDataForm{
			Data:   nil,
			TaskID: 0,
			Number: 0,
		}, u}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SaveRawData(tt.args.form, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("SaveRawData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
