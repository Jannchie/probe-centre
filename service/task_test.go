package service

import (
	"reflect"
	"testing"

	"github.com/Jannchie/probe-centre/db"

	"github.com/Jannchie/probe-centre/model"

	"github.com/Jannchie/probe-centre/test"

	"github.com/gin-gonic/gin"
)

func TestGetTaskStats(t *testing.T) {
	test.InitDB()
	db.DB.Exec("DELETE FROM tasks")
	tests := []struct {
		name    string
		want    gin.H
		wantErr bool
	}{
		{"ok", gin.H{"Pending": int64(0), "Finished": int64(0), "Waiting": int64(0)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTaskStats()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTaskStats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTaskStats() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetOneTask(t *testing.T) {
	test.InitDB()
	type args struct {
		task *model.Task
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"err", args{&model.Task{}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetOneTask(tt.args.task); (err != nil) != tt.wantErr {
				t.Errorf("GetOneTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
