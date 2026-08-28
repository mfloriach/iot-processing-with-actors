package mesures

import (
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

var (
	Stats     LatencyStats
	Generated atomic.Uint64
	Processed atomic.Uint64
)

type LatencyStats struct {
	mu      sync.Mutex
	samples []time.Duration
}

func (s *LatencyStats) Add(d time.Duration) {
	s.mu.Lock()
	s.samples = append(s.samples, d)
	s.mu.Unlock()
}

func (s *LatencyStats) Percentiles() (p50, p90, p99 time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.samples) == 0 {
		return 0, 0, 0
	}

	values := append([]time.Duration(nil), s.samples...)
	sort.Slice(values, func(i, j int) bool {
		return values[i] < values[j]
	})

	percentile := func(p float64) time.Duration {
		index := int(float64(len(values)-1) * p)
		return values[index]
	}

	return percentile(0.50),
		percentile(0.90),
		percentile(0.99)
}
