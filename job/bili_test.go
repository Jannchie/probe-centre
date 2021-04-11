package job

import (
	"testing"

	"github.com/Jannchie/probe-centre/test"
)

func TestCreateBiliTasks(t *testing.T) {
	test.InitDB()
	CreateBiliTasks()
}
