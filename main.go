package main

import (
	"log/slog"
	"os"
	"strconv"
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

	manager := device.NewDeviceManager(4, 30)
	for i := range 1000 {
		id := "sensor-" + strconv.Itoa(i)
		manager.Add(NewActor(id, device.DeviceState{}, device.Dispatch))
	}

	sched := scheduler.NewScheduler(4)
	go sched.Run(manager)

	injestor := injestor.NewInjestorRandom()

	go func() {
		for t := range injestor.Run("sensor-0", time.Second) {
			manager.Send(t)
		}
	}()

	go func() {
		for t := range injestor.Run("sensor-1", time.Millisecond) {
			manager.Send(t)
		}
	}()

	go func() {
		for t := range injestor.Run("sensor-2", time.Millisecond) {
			manager.Send(t)
		}
	}()

	go func() {
		for t := range injestor.Run("sensor-3", time.Millisecond) {
			manager.Send(t)
		}
	}()

	StartServer(manager)
}
