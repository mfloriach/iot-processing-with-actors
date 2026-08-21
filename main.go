package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"datacollector/device"
	"datacollector/injestor"
	"datacollector/scheduler"
)

func main() {
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
		for i := 0; i < 10; i++ {
			state := manager.State("sensor-001")
			jsonData, err := json.Marshal(state)
			if err != nil {
				log.Fatalf("Error marshaling to JSON: %s", err)
			}

			fmt.Println(string(jsonData))

			time.Sleep(50 * time.Millisecond)
		}
	}()

	go func() {
		for i := 0; i < 10; i++ {
			state := manager.State("sensor-001")
			jsonData, err := json.Marshal(state)
			if err != nil {
				log.Fatalf("Error marshaling to JSON: %s", err)
			}

			fmt.Println(string(jsonData))

			time.Sleep(80 * time.Millisecond)
		}
	}()

	injestor := injestor.NewInjestorRandom(10)
	for t := range injestor.Run() {
		manager.Send(t.DeviceID, t)
	}

	time.Sleep(2 * time.Second)
}
