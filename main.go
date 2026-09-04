package main

import (
	"log/slog"
	"os"

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

	injestor := injestor.NewInjestorRandom()
	sched := scheduler.NewScheduler(&manager, injestor)
	go sched.Run()

	go StartServer(manager)

	mesures.Stats.Run()
}
