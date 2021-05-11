package util

import (
	"testing"
	"time"
)

func TestSpeedometer(t *testing.T) {
	s := NewSpeedometer()
	for range time.Tick(time.Second) {
		s.AddCount(1)
		if s.count == 3 {
			break
		}
	}
	speedStat := s.GetStat()
	t.Logf("%+v", speedStat)
	s.Stop()
}
