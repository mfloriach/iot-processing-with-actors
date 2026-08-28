package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"runtime/trace"
	"strconv"
	"time"

	"datacollector/device"
	"datacollector/injestor"
	"datacollector/scheduler"

	spretty "github.com/mickamy/slog-pretty"
)

func main() {
	var m runtime.MemStats

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

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
		case <-ctx.Done():
			// Read the current memory statistics
			runtime.ReadMemStats(&m)

			// m.NumGC holds the total number of completed GC cycles
			fmt.Printf("The GC has triggered %d times.\n", m.NumGC)
			fmt.Println("Time is up! Stopping execution.")
			return
		default:
			// Place your working code here
			fmt.Println("Doing work...")
			time.Sleep(5 * time.Second) // Simulating work
		}
	}
}
