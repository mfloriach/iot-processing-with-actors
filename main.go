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
	mesures.NewLatencyStats()

	slog.SetDefault(slog.New(spretty.NewHandler(os.Stdout, &spretty.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})))

	deviceStatePool := libs.NewPool(func() *device.DeviceState {
		return new(device.DeviceState)
	})

	hooks := libs.ActorHooks[device.DeviceState, device.Telemetry]{
		UpdateState:    device.UpdateDeviceState,
		OnBackpressure: mesures.Stats.AddMailboxBackpressure,
	}

	manager := device.NewDeviceManager(config.SENSOR_PER_WORKER * config.NUM_CPUS)
	for i := 0; i < config.NUM_CPUS*config.SENSOR_PER_WORKER; i++ {
		manager.AddDevice(libs.NewActor(i, config.MAILBOX_SIZE, deviceStatePool, hooks))
	}

	sched := scheduler.NewScheduler(manager)
	go sched.Run()

	go StartServer(manager)

	mesures.Stats.Run()
}
