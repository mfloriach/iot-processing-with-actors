package injestor

import (
	"datacollector/device"
	"iter"
	"math/rand"
	"time"
)

type Injestor interface {
	Run(string, time.Duration) iter.Seq[device.Message]
}

type InjestorRandom struct{}

func NewInjestorRandom() Injestor {
	return InjestorRandom{}
}

func (i InjestorRandom) Run(id string, duration time.Duration) iter.Seq[device.Message] {
	return func(yield func(device.Message) bool) {
		ticker := time.NewTicker(duration)
		defer ticker.Stop()

		for range ticker.C {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			ok := yield(device.Telemetry{
				Sample: device.Sample{
					DeviceID: id,
					TTL:      time.Now(),
				},

				Temperature: float64(r.Intn(101)),
				Humidity:    float64(r.Intn(71)),
				Battery:     float64(r.Intn(101)),
				Noise:       float64(r.Intn(21)),
			})

			if !ok {
				return
			}
		}
	}
}
