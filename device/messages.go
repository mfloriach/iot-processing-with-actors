package device

import (
	"datacollector/mesures"
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

func (t Telemetry) Apply(state *DeviceState) {
	state.Data = t
	state.Online = true

	elapsed := time.Since(t.TTL)
	mesures.Stats.Add(elapsed)
}

func (t Telemetry) GetDeviceID() int {
	return t.Sample.DeviceID
}
