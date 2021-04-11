package job

import (
	"testing"

	"github.com/Jannchie/probe-centre/test"
)

func TestRemoveIpRecord(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"1"},
	}
	test.InitDB()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			removeIPTask()
		})
	}
}
