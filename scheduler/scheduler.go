package scheduler

import (
	"datacollector/device"
	"datacollector/injestor"
	"fmt"
	"math/rand"
	"time"
)

const (
	NUM_OF_WORKERS_UPDATING        = 5
	NUM_OF_WORKERS_GENERETIC_NOISE = 20

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
	s.receiveNoise()
	s.updateStatus()
}

func (s Scheduler) updateStatus() {
	for range NUM_OF_WORKERS_UPDATING {
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
	for i := 1; i <= NUM_OF_WORKERS_GENERETIC_NOISE; i++ {
		go func() {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			deviceID := fmt.Sprintf("sensor-%d", r.Intn(10))
			for d := range s.injestor.Run(deviceID, r) {
				s.manager.Send(d)
			}
		}()
	}
}

func (s Scheduler) receiveAnalysis(sensorID string) {
	go func() {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for range ticker.C {
			d := device.Telemetry{
				Sample: device.Sample{
					DeviceID: sensorID,
					TTL:      time.Now(),
				},

				Temperature: float64(r.Intn(101)),
				Humidity:    float64(r.Intn(71)),
				Battery:     float64(r.Intn(101)),
				Noise:       float64(r.Intn(21)),
			}

			s.manager.Send(d)
		}
	}()
}
