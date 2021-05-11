package util

import (
	"fmt"
	"testing"
	"time"
)

func TestSpeedometer(t *testing.T) {
	s := NewSpeedometer()
	for range time.Tick(time.Millisecond * 100) {
		s.AddCount(1)
		if s.count == 20 {
			break
		}
	}
	speedStat := s.GetStat()
	fmt.Printf("%+v", speedStat)
	s.Stop()
}
