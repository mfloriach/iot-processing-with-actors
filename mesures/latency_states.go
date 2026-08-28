package mesures

import (
	"log/slog"
	"math"
	"runtime"
	"sync/atomic"
	"time"
)

var (
	Stats     LatencyStats
	Generated atomic.Uint64
	Processed atomic.Uint64
)

type LatencyStats struct {
	buckets     [64]atomic.Uint64
	PrevMallocs uint64
}

// Add registra una latencia.
// Cada bucket representa un rango de latencia.
func (s *LatencyStats) Add(d time.Duration) {
	if d < 0 {
		return
	}

	// Convertimos a nanosegundos.
	ns := uint64(d)

	// Bucket logarítmico.
	var bucket int

	for ns > 0 {
		bucket++
		ns >>= 1
	}

	if bucket >= len(s.buckets) {
		bucket = len(s.buckets) - 1
	}

	s.buckets[bucket].Add(1)
}

func (s *LatencyStats) Percentiles() (
	p50, p90, p99 time.Duration,
) {
	total := uint64(0)

	for i := range s.buckets {
		total += s.buckets[i].Load()
	}

	if total == 0 {
		return 0, 0, 0
	}

	p50 = s.percentile(total, 0.50)
	p90 = s.percentile(total, 0.90)
	p99 = s.percentile(total, 0.99)

	return
}

func (s *LatencyStats) percentile(
	total uint64,
	percent float64,
) time.Duration {

	target := uint64(float64(total-1) * percent)

	var count uint64

	for i := range s.buckets {
		count += s.buckets[i].Load()

		if count > target {
			// El bucket representa aproximadamente:
			// [2^(i-1), 2^i)
			if i == 0 {
				return 0
			}

			return time.Duration(uint64(1) << (i - 1))
		}
	}

	return time.Duration(math.MaxInt64)
}

func (s *LatencyStats) Reset() {
	for i := range s.buckets {
		s.buckets[i].Store(0)
	}
}

func (s *LatencyStats) PrintResults(m runtime.MemStats) {
	p50, p90, p99 := s.Percentiles()

	runtime.ReadMemStats(&m)

	g := Generated.Swap(0)
	p := Processed.Swap(0)

	mallocsDelta := m.Mallocs - s.PrevMallocs
	s.PrevMallocs = m.Mallocs

	var backlog uint64

	if g > p {
		backlog = g - p
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
		"allocs_per_message", math.Round(float64(mallocsDelta)/float64(g)),
	)

	s.Reset()
}
