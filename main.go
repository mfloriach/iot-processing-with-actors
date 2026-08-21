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
	manager := device.NewDeviceManager()
	manager.Add("sensor-000")
	manager.Add("sensor-001")
	manager.Add("sensor-002")
	manager.Add("sensor-003")
	manager.Add("sensor-004")
	manager.Add("sensor-005")

	sched := scheduler.NewScheduler()
	go sched.Run(manager)

	go func() {
		for i := 0; i < 10; i++ {
			state := manager.GetState("sensor-001")
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
			state := manager.GetState("sensor-001")
			jsonData, err := json.Marshal(state)
			if err != nil {
				log.Fatalf("Error marshaling to JSON: %s", err)
			}

			fmt.Println(string(jsonData))

			time.Sleep(80 * time.Millisecond)
		}
	}()

	injestor := injestor.NewInjestorRandom(50)
	for t := range injestor.Run() {
		manager.Send(t.DeviceID, t)
	}

	time.Sleep(time.Second)
	manager.ShutDown("sensor-001")

	fmt.Println("device stopped")
}
