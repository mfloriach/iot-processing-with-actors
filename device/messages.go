package device

import (
	"datacollector/mesures"
	"time"
)

func UpdateDeviceState(state *DeviceState, msg Telemetry) {
	state.Data = msg
	state.Online = true

	elapsed := time.Since(msg.TTL)
	mesures.Stats.Add(elapsed)
}

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
