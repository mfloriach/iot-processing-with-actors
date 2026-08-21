package main

import (
	"datacollector/device"
	"encoding/json"
	"fmt"
	"log"
	"time"

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

	injestor := scheduler.NewInjestorRandom(50)
	scheduler := scheduler.NewScheduler(injestor.Run)
	go scheduler.Run(manager.Handler)

	// Reader 1
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

	// Reader 2
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

	time.Sleep(time.Second)
	// manager.Get("sensor-001").Send(device.Shutdown{})

	fmt.Println("device stopped")
}
