package injestor

import (
	"datacollector/device"
	"iter"
)

type InjestorRandom struct {
	count int
}

func NewInjestorRandom(count int) InjestorRandom {
	return InjestorRandom{count: count}
}

func (i InjestorRandom) Run() iter.Seq[device.Telemetry] {
	return func(yield func(device.Telemetry) bool) {
		for i := i.count; i >= 1; i-- {
			// id := strconv.Itoa(i % 5)
			if !yield(device.Telemetry{
				DeviceID:    "sensor-001",
				Temperature: 20,
				Humidity:    50,
				Battery:     10,
				Noise:       5,
			}) {
				return
			}
		}
	}
}
