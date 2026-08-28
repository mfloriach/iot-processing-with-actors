package scheduler

import (
	"datacollector/device"
	"datacollector/injestor"
	"datacollector/mesures"
	"fmt"
	"sync"
	"time"
)

const (
	REQUEST_PER_SECOND_NOISE = 120_000
	MAX_REQUEST_PER_WORKER   = 12_000

	NUM_OF_WORKERS_UPDATING        = 10
	NUM_OF_WORKERS_GENERETIC_NOISE = 10

	SENSOR_OF_ANALYSIS_ID = "sensor-0"
)

type Scheduler struct {
	manager  device.DeviceManager
	injestor injestor.Injestor
}

func NewScheduler(manager device.DeviceManager, injestor injestor.Injestor) Scheduler {
	return Scheduler{
		manager:  manager,
		injestor: injestor,
	}
}

func (s Scheduler) Run() {
	s.receiveAnalysis(SENSOR_OF_ANALYSIS_ID)
	s.receiveNoise(REQUEST_PER_SECOND_NOISE)
	s.updateStatus()
}

func (s Scheduler) updateStatus() {
	var wg sync.WaitGroup

	for range NUM_OF_WORKERS_UPDATING {
		wg.Add(1)

		go func() {
			defer wg.Done()

			s.manager.Process()
		}()
	}
}

func (s Scheduler) receiveNoise(packetsPerSecond int) {
	num_workers := packetsPerSecond / MAX_REQUEST_PER_WORKER
	if num_workers > NUM_OF_WORKERS_GENERETIC_NOISE {
		fmt.Println("too much goroutines")
	}

	for i := 1; i <= num_workers; i++ {
		deviceID := fmt.Sprintf("sensor-%d", i)

		go func(deviceID string) {
			for t := range s.injestor.Run(deviceID, time.Second/MAX_REQUEST_PER_WORKER) {
				mesures.Generated.Add(1)
				s.manager.Send(t)
			}
		}(deviceID)
	}
}

func (s Scheduler) receiveAnalysis(sensorID string) {
	go func() {
		for t := range s.injestor.Run(sensorID, time.Second) {
			mesures.Generated.Add(1)
			s.manager.Send(t)
		}
	}()
}
