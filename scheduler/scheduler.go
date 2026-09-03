package scheduler

import (
	"datacollector/config"
	"datacollector/device"
	"datacollector/injestor"
	"fmt"
	"math/rand"
	"time"
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
	s.receiveAnalysis(config.SENSOR_OF_ANALYSIS_ID)
	s.receiveNoise()
	s.updateStatus()
}

func (s Scheduler) updateStatus() {
	for range config.NUM_OF_WORKERS_UPDATING {
		go func() {
			for {
				for d := range s.manager.Next() {
					if hasNext := s.manager.ProcessOne(d); !hasNext {
						break
					}
				}
				time.Sleep(time.Millisecond * 10)
			}
		}()
	}
}

func (s Scheduler) receiveNoise() {
	workers := config.NUM_OF_WORKERS_GENERETIC_NOISE
	sensors := config.NUM_OF_SENSORS

	for i := 0; i < workers; i++ {
		go func() {
			start := i * sensors / workers
			end := (i + 1) * sensors / workers

			r := rand.New(rand.NewSource(time.Now().UnixNano()))

			for deviceID := start; deviceID < end; deviceID++ {
				id := fmt.Sprintf("sensor-%d", deviceID+1)

				for d := range s.injestor.Run(id, r) {
					s.manager.Send(d)
				}
			}
		}()
	}
}

func (s Scheduler) receiveAnalysis(sensorID string) {
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		count := 0

		for range ticker.C {
			d := device.Telemetry{
				Sample: device.Sample{
					DeviceID: sensorID,
					TTL:      time.Now(),
				},

				Temperature: float64(count),
				Humidity:    float64(count),
				Battery:     float64(count),
				Noise:       float64(count),
			}

			s.manager.Send(d)
			count++
		}
	}()
}
