package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"datacollector/device"
	"datacollector/injestor"
	"datacollector/mesures"
	"datacollector/scheduler"

	_ "net/http/pprof"

	spretty "github.com/mickamy/slog-pretty"
)

func main() {
	runtime.GOMAXPROCS(1)
	var m runtime.MemStats

	go func() {
		slog.Error(
			"pprof",
			"error",
			http.ListenAndServe("localhost:6060", nil),
		)
	}()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Timeout after 5 minutes
	timeout := time.After(5 * time.Minute)

	logger := slog.New(spretty.NewHandler(os.Stdout, &spretty.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))
	slog.SetDefault(logger)

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
			mesures.Stats.PrintResults(m)
		case <-timeout:
			mesures.Stats.PrintResults(m)

			fmt.Println("Time is up! Stopping execution.")
			return
		default:
			time.Sleep(5 * time.Second) // Simulating work
		}
	}
}
