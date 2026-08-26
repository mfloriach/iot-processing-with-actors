package main

import (
	"log/slog"
	"os"
	"time"

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
	manager.Add(NewActor("sensor-000", device.DeviceState{}, device.Dispatch))
	manager.Add(NewActor("sensor-001", device.DeviceState{}, device.Dispatch))
	manager.Add(NewActor("sensor-002", device.DeviceState{}, device.Dispatch))
	manager.Add(NewActor("sensor-003", device.DeviceState{}, device.Dispatch))
	manager.Add(NewActor("sensor-004", device.DeviceState{}, device.Dispatch))
	manager.Add(NewActor("sensor-005", device.DeviceState{}, device.Dispatch))

	sched := scheduler.NewScheduler(4)
	go sched.Run(manager)

	go func() {
		injestor := injestor.NewInjestorRandom()
		for t := range injestor.Run("sensor-000", time.Second) {
			manager.Send(t)
		}
	}()

	go func() {
		injestor := injestor.NewInjestorRandom()
		for t := range injestor.Run("sensor-001", time.Millisecond) {
			manager.Send(t)
		}
	}()

	go func() {
		injestor := injestor.NewInjestorRandom()
		for t := range injestor.Run("sensor-002", time.Millisecond) {
			manager.Send(t)
		}
	}()

	go func() {
		injestor := injestor.NewInjestorRandom()
		for t := range injestor.Run("sensor-003", time.Millisecond) {
			manager.Send(t)
		}
	}()

	StartServer(manager)
}
