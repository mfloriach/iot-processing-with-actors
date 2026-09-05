package main

import (
	"log/slog"
	"os"

	"datacollector/config"
	"datacollector/device"
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
	for i := 0; i < config.NUM_CPUS*config.SENSOR_PER_WORKER; i++ {
		manager.AddDevice(device.NewActor(i))
	}

	sched := scheduler.NewScheduler(manager)
	go sched.Run()

	go StartServer(manager)

	mesures.Stats.Run()
}
