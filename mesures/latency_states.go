package mesures

import (
	"fmt"
	"log/slog"
	"runtime"
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
	if d > time.Second {
		fmt.Println("out of time")
	}
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

func (s *LatencyStats) PrintResults(m runtime.MemStats) {
	p50, p90, p99 := s.Percentiles()
	runtime.ReadMemStats(&m)
	g := Generated.Swap(0)
	p := Processed.Swap(0)

	backlog := int64(g) - int64(p)
	if backlog < 0 {
		backlog = 0
	}

	slog.Info(
		"telemetry processed",
		"p50", p50,
		"p90", p90,
		"p99", p99,
		"generated", g/30,
		"processed", p/30,
		"backlog", backlog,
		"gc", m.NumGC,
		"gc_cpu", m.GCCPUFraction,
		"heap_mb", m.HeapAlloc/1024/1024,
		"heap_objects", m.HeapObjects,
		"total_alloc_mb", m.TotalAlloc/1024/1024,
		"mallocs", m.Mallocs,
		"frees", m.Frees,
	)
}
