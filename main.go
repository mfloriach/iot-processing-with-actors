package main

import (
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"runtime/trace"
	"strconv"
	"time"

	"datacollector/device"
	"datacollector/injestor"
	"datacollector/mesures"
	"datacollector/scheduler"

	spretty "github.com/mickamy/slog-pretty"
)

func main() {
	var m runtime.MemStats

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Timeout after 5 minutes
	timeout := time.After(5 * time.Minute)

	logger := slog.New(spretty.NewHandler(os.Stdout, &spretty.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))
	slog.SetDefault(logger)

	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := trace.Start(f); err != nil {
		panic(err)
	}
	defer trace.Stop()

	manager := device.NewDeviceManager(30)
	for i := range 1000 {
		id := "sensor-" + strconv.Itoa(i)
		manager.Add(NewActor(id, device.DeviceState{}, device.Dispatch))
	}

	injestor := injestor.NewInjestorRandom()
	sched := scheduler.NewScheduler(manager, injestor)
	go sched.Run()

	go StartServer(manager)

	for {
		select {
		case <-ticker.C:
			p50, p90, p99 := mesures.Stats.Percentiles()
			runtime.ReadMemStats(&m)
			g := mesures.Generated.Swap(0)
			p := mesures.Processed.Swap(0)

			backlog := int64(g) - int64(p)
			if backlog < 0 {
				backlog = 0
			}

			slog.Info(
				"telemetry processed",
				"p50", p50,
				"p90", p90,
				"p99", p99,
				"generated", g,
				"processed", p,
				"backlog", backlog,
				"garbage", m.NumGC,
			)
		case <-timeout:
			p50, p90, p99 := mesures.Stats.Percentiles()
			runtime.ReadMemStats(&m)

			slog.Info(
				"telemetry processed",
				"p50", p50,
				"p90", p90,
				"p99", p99,
				"garbage", m.NumGC,
			)

			fmt.Println("Time is up! Stopping execution.")
			return
		default:
			time.Sleep(5 * time.Second) // Simulating work
		}
	}
}
