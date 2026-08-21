package injestor

import (
	"context"
	"datacollector/device"
	"iter"
	"math/rand"
	"time"
)

type InjestorRandom struct {
	count int
}

func NewInjestorRandom(count int) InjestorRandom {
	return InjestorRandom{count: count}
}

func (i InjestorRandom) Run() iter.Seq[device.Telemetry] {
	return func(yield func(device.Telemetry) bool) {
		for {
			ctx, cancel := context.WithTimeout(
				context.TODO(),
				5*time.Second,
			)
			defer cancel()

			if !yield(device.Telemetry{
				DeviceID:    "sensor-001",
				Temperature: randomNumber(50),
				Humidity:    randomNumber(70),
				Battery:     randomNumber(100),
				Noise:       randomNumber(20),

				Context: ctx,
			}) {
				return
			}
		}
	}
}

func randomNumber(max int) float64 {
	return float64(rand.Intn(max + 1))
}
