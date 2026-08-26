package injestor

import (
	"datacollector/device"
	"iter"
	"math/rand"
	"time"
)

type InjestorRandom struct{}

func NewInjestorRandom() InjestorRandom {
	return InjestorRandom{}
}

func (i InjestorRandom) Run(id string, duration time.Duration) iter.Seq[device.Message] {
	return func(yield func(device.Message) bool) {
		ticker := time.NewTicker(duration)
		defer ticker.Stop()

		for range ticker.C {
			ok := yield(device.Telemetry{
				Sample: device.Sample{
					DeviceID: id,
					TTL:      time.Now(),
				},

				Temperature: randomNumber(50),
				Humidity:    randomNumber(70),
				Battery:     randomNumber(100),
				Noise:       randomNumber(20),
			})

			if !ok {
				return
			}
		}
	}
}

func randomNumber(max int) float64 {
	return float64(rand.Intn(max + 1))
}
