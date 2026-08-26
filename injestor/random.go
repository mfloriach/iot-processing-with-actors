package injestor

import (
	"context"
	"datacollector/device"
	"iter"
	"math/rand"
	"time"
)

type InjestorRandom struct {
}

func NewInjestorRandom() InjestorRandom {
	return InjestorRandom{}
}

func (i InjestorRandom) Run() iter.Seq[device.Message] {
	return func(yield func(device.Message) bool) {
		for {
			ctx, cancel := context.WithTimeout(
				context.TODO(),
				5*time.Second,
			)
			defer cancel()

			if !yield(device.Telemetry{
				Sample: device.Sample{
					DeviceID: "sensor-001",
					Context:  ctx,
				},

				Temperature: randomNumber(50),
				Humidity:    randomNumber(70),
				Battery:     randomNumber(100),
				Noise:       randomNumber(20),
			}) {
				return
			}
		}
	}
}

func randomNumber(max int) float64 {
	return float64(rand.Intn(max + 1))
}
