package service

import (
	"testing"

	"github.com/Jannchie/probe-centre/test"

	"github.com/Jannchie/probe-centre/model"
)

func TestLoginByForm(t *testing.T) {
	test.InitDB()
	test.CreateTestUser()
	type args struct {
		form model.LoginForm
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"ok", args{model.LoginForm{Mail: "test@test.com", Password: "123456"}}, false},
		{"ok", args{model.LoginForm{Mail: "test@gmail.com", Password: "1234561"}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := LoginByForm(tt.args.form)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoginByForm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
