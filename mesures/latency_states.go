package mesures

import (
	"datacollector/config"
	"net/http"
	// "datacollector/mesures"
	"fmt"
	"log/slog"
	"math"
	"runtime"
	"sync/atomic"
	"time"
)

var (
	Stats LatencyStats
)

type LatencyStats struct {
	buckets     [64]atomic.Uint64
	PrevMallocs uint64
	generated   atomic.Uint64
	processed   atomic.Uint64
}

func NewLatencyStats() {
	runtime.GOMAXPROCS(config.NUM_CPUS)

	Stats = LatencyStats{}

	go func() {
		slog.Error(
			"pprof",
			"error",
			http.ListenAndServe("localhost:6060", nil),
		)
	}()
}

func (s *LatencyStats) Run() {
	var m runtime.MemStats

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	timeout := time.After(5 * time.Minute)

	for {
		select {
		case <-ticker.C:
			Stats.PrintResults(m)
		case <-timeout:
			Stats.PrintResults(m)

			fmt.Println("Time is up! Stopping execution.")
			return
		default:
			time.Sleep(5 * time.Second) // Simulating work
		}
	}
}

func (s *LatencyStats) AddGenerate() {
	s.generated.Add(1)
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
	s.processed.Add(1)
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

	g := s.generated.Swap(0)
	p := s.processed.Swap(0)

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

		"generated_sec", g/30,
		"processed_sec", p/30,
		"backlog", backlog,

		"gc", m.NumGC,
		"gc_cpu", m.GCCPUFraction,

		"heap_mb", m.HeapAlloc/1024/1024,
		"heap_objects", m.HeapObjects,

		"num_cpus", config.NUM_CPUS,
		"num_of_workers", config.NUM_OF_WORKERS_GENERETIC_NOISE,
		"num_of_workers_updating", config.NUM_OF_WORKERS_UPDATING,

		"total_alloc_mb", m.TotalAlloc/1024/1024,
		"mallocs", m.Mallocs,
		"frees", m.Frees,
		"allocs_per_message", math.Round(float64(mallocsDelta)/float64(g)),
	)

	s.Reset()
}
