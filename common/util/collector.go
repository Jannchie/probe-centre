package util

import (
	"fmt"
	"sync"
	"time"
)

type speedometer struct {
	count    uint64
	guard    chan struct{}
	duration time.Duration
	history  []uint64
	mutex    sync.RWMutex
}

type SpeedStat struct {
	Count uint64
	Speed uint64
}

func (s *speedometer) GetStat() SpeedStat {
	ss := SpeedStat{}
	ss.Count = s.count
	var delta uint64
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	count := len(s.history)
	if count <= 1 {
		ss.Speed = 0
		return ss
	} else {
		deltaTime := uint64(count-1) * uint64(s.duration)
		delta = s.history[count-1] - s.history[0]
		ss.Speed = delta * uint64(time.Minute) / deltaTime
		return ss
	}
}

func (s *speedometer) startTicker() {
	s.mutex.Lock()
	ticker := time.NewTicker(s.duration)
	s.guard = make(chan struct{})
	s.mutex.Unlock()

	l := 0
	for {
		select {
		case _, ok := <-s.guard:
			if !ok {
				ticker.Stop()
				return
			}
		case <-ticker.C:
			s.mutex.Lock()
			s.history = append(s.history, s.count)
			s.mutex.Unlock()
			if l < 60 {
				l += 1
			} else {
				s.history = s.history[1:]
			}
		}
	}
}

func (s *speedometer) AddCount(n uint64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.count += n
}

func (s *speedometer) AutoPrint() {
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ticker.C:
			stat := s.GetStat()
			fmt.Printf("Speed: %d/min Total: %d\n", stat.Speed, stat.Count)
		case _, ok := <-s.guard:
			if !ok {
				ticker.Stop()
				return
			}
		}
	}
}

func (s *speedometer) Stop() {
	s.mutex.Lock()
	s.guard <- struct{}{}
	s.mutex.Unlock()
}

func NewSpeedometer() *speedometer {
	s := &speedometer{}
	if s.duration == 0 {
		s.duration = time.Second * 1
	}
	go s.startTicker()
	go s.AutoPrint()
	return s
}
