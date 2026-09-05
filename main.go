package main

import (
	"log/slog"
	"os"
	"time"

	"datacollector/config"
	"datacollector/device"
	"datacollector/device/messages"
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

	hooks := libs.ActorHooks[messages.Message]{
		OnApply: func(msg messages.Message) {
			elapsed := time.Since(msg.TTL)
			mesures.Stats.Add(elapsed)
		},
		OnBackpressure: mesures.Stats.AddMailboxBackpressure,
	}

	manager := device.NewDeviceManager(config.SENSOR_PER_WORKER * config.NUM_CPUS)
	for i := 0; i < config.NUM_CPUS*config.SENSOR_PER_WORKER; i++ {
		manager.AddDevice(libs.NewActor(i, config.MAILBOX_SIZE, deviceStatePool, hooks, device.Apply))
	}

	mesures.NewLatencyStats()

	sched := scheduler.NewScheduler(manager)
	go sched.Run()

	go StartServer(manager)

	mesures.Stats.Run()
}
