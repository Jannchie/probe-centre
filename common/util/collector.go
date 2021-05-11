package util

import (
	"time"
)

type speedometer struct {
	count    uint64
	ticker   *time.Ticker
	guard    chan struct{}
	duration time.Duration
	history  []uint64
}

type SpeedStat struct {
	Count uint64
	Speed uint64
}

func (s *speedometer) GetStat() SpeedStat {
	ss := SpeedStat{}
	ss.Count = s.count
	var delta uint64
	if len(s.history) <= 1 {
		ss.Speed = 0
		return ss
	} else {
		count := len(s.history) - 1
		deltaTime := uint64(count) * uint64(s.duration)
		delta = s.history[count] - s.history[0]
		ss.Speed = delta * uint64(time.Minute) / deltaTime
		return ss
	}
}

func (s *speedometer) startTicker() {

	s.ticker = time.NewTicker(s.duration)
	s.guard = make(chan struct{})
	l := 0
	for {
		select {
		case <-s.guard:
			break
		case <-s.ticker.C:
			s.history = append(s.history, s.count)
			if l < 60 {
				l += 1
			} else {
				s.history = s.history[1:]
			}
		}
	}
}

func (s *speedometer) AddCount(n uint64) {
	s.count += n
}

func (s *speedometer) Stop() {
	s.ticker.Stop()
	s.guard <- struct{}{}
}

func NewSpeedometer() *speedometer {
	s := &speedometer{}
	if s.duration == 0 {
		s.duration = time.Second * 1
	}
	go s.startTicker()
	return s
}
