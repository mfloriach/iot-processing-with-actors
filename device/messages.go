package device

import (
	"time"
)

type DeviceState struct {
	Data   Telemetry
	Online bool
}

type Sample struct {
	DeviceID int
	TTL      time.Time `json:"-"`
}

type Telemetry struct {
	Sample

	Temperature float64
	Humidity    float64
	Battery     float64
	Noise       float64
}
