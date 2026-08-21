package main

import (
	"encoding/json"
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
	manager.Add("sensor-000", NewActor("sensor-000", device.DeviceState{}, device.Dispatch))
	manager.Add("sensor-001", NewActor("sensor-001", device.DeviceState{}, device.Dispatch))
	manager.Add("sensor-002", NewActor("sensor-002", device.DeviceState{}, device.Dispatch))
	manager.Add("sensor-003", NewActor("sensor-003", device.DeviceState{}, device.Dispatch))
	manager.Add("sensor-004", NewActor("sensor-004", device.DeviceState{}, device.Dispatch))
	manager.Add("sensor-005", NewActor("sensor-005", device.DeviceState{}, device.Dispatch))

	sched := scheduler.NewScheduler(4)
	go sched.Run(manager)

	go func() {
		ticker := time.NewTicker(time.Second)

		for t := range ticker.C {
			state := manager.State("sensor-001")
			jsonData, err := json.Marshal(state)
			if err != nil {
				slog.Error("Error marshaling to JSON", slog.Any("error", err))
			}

			slog.Info("lister 1",
				slog.String("state", string(jsonData)),
				slog.String("time", t.String()),
			)
		}
	}()

	go func() {
		ticker := time.NewTicker(time.Second)

		for t := range ticker.C {
			state := manager.State("sensor-001")
			jsonData, err := json.Marshal(state)
			if err != nil {
				slog.Error("Error marshaling to JSON", slog.Any("error", err))
			}

			slog.Info("lister 2",
				slog.String("state", string(jsonData)),
				slog.String("time", t.String()),
			)
		}
	}()

	go func() {
		injestor := injestor.NewInjestorRandom(10)
		for t := range injestor.Run() {
			manager.Send(t.DeviceID, t)
		}
	}()

	StartServer(manager)
}
