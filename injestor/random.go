package injestor

import (
	"datacollector/device"
	"iter"
	"math/rand"
	"time"
)

type Injestor interface {
	Run(string, *rand.Rand) iter.Seq[device.Telemetry]
}

type InjestorRandom struct{}

func NewInjestorRandom() Injestor {
	return InjestorRandom{}
}

func (i InjestorRandom) Run(
	deviceID string,
	r *rand.Rand,
) iter.Seq[device.Telemetry] {
	return func(yield func(device.Telemetry) bool) {
		for {
			ok := yield(device.Telemetry{
				Sample: device.Sample{
					DeviceID: deviceID,
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
