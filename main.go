package main

import (
	"log/slog"
	"os"

	"datacollector/config"
	"datacollector/device"
	"datacollector/libs"
	"datacollector/mesures"
	"datacollector/scheduler"

	_ "net/http/pprof"

	spretty "github.com/mickamy/slog-pretty"
)

func main() {
	slog.SetDefault(slog.New(spretty.NewHandler(os.Stdout, &spretty.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})))

	deviceStatePool := libs.NewPool(func() *device.DeviceState {
		return new(device.DeviceState)
	})

	hooks := libs.ActorHooks{
		OnBackpressure: mesures.Stats.AddMailboxBackpressure,
	}

	manager := device.NewDeviceManager(config.SENSOR_PER_WORKER * config.NUM_CPUS)
	for i := 0; i < config.NUM_CPUS*config.SENSOR_PER_WORKER; i++ {
		manager.AddDevice(libs.NewActor[device.DeviceState, libs.Message[device.DeviceState]](i, config.MAILBOX_SIZE, deviceStatePool, hooks))
	}

	mesures.NewLatencyStats()

	sched := scheduler.NewScheduler(manager)
	go sched.Run()

	go StartServer(manager)

	mesures.Stats.Run()
}
