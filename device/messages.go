package device

import (
	"time"
)

type Sample struct {
	DeviceID string
	TTL      time.Time `json:"-"`
}

type Telemetry struct {
	Sample

	Temperature float64
	Humidity    float64
	Battery     float64
	Noise       float64
}
