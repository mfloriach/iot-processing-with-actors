package main

import (
	"log/slog"
	"os"
	"strconv"

	"datacollector/device"
	"datacollector/injestor"
	"datacollector/mesures"
	"datacollector/scheduler"

	_ "net/http/pprof"

	spretty "github.com/mickamy/slog-pretty"
)

func main() {
	mesures.NewLatencyStats()

	slog.SetDefault(slog.New(spretty.NewHandler(os.Stdout, &spretty.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})))

	manager := device.NewDeviceManager()
	for i := range 10 {
		manager.AddDevice(NewActor("sensor-"+strconv.Itoa(i), device.Dispatch))
	}

	injestor := injestor.NewInjestorRandom()
	sched := scheduler.NewScheduler(manager, injestor)
	go sched.Run()

	go StartServer(manager)

	mesures.Stats.Run()
}
