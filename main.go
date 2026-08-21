package main

import (
	"log/slog"
	"os"

	"datacollector/device"
	"datacollector/injestor"
	"datacollector/scheduler"

	spretty "github.com/mickamy/slog-pretty"
)

func main() {
	logger := slog.New(spretty.NewHandler(os.Stdout, &spretty.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))
	slog.SetDefault(logger)

	manager := device.NewDeviceManager(4)
	manager.Add("sensor-000", NewActor("sensor-000", device.DeviceState{}, device.Dispatch))
	manager.Add("sensor-001", NewActor("sensor-001", device.DeviceState{}, device.Dispatch))
	manager.Add("sensor-002", NewActor("sensor-002", device.DeviceState{}, device.Dispatch))
	manager.Add("sensor-003", NewActor("sensor-003", device.DeviceState{}, device.Dispatch))
	manager.Add("sensor-004", NewActor("sensor-004", device.DeviceState{}, device.Dispatch))
	manager.Add("sensor-005", NewActor("sensor-005", device.DeviceState{}, device.Dispatch))

	sched := scheduler.NewScheduler(4)
	go sched.Run(manager)

	go func() {
		injestor := injestor.NewInjestorRandom()
		for t := range injestor.Run() {
			manager.Send(t.GetDeviceID(), t)
		}
	}()

	StartServer(manager)
}
