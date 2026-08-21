package device

type Telemetry struct {
	DeviceID    string
	Temperature float64
	Humidity    float64
	Battery     float64
	Noise       float64
}

func (Telemetry) IsMessage() {}

func (m Telemetry) Apply(state *DeviceState) {
	state.Data = m
	state.Online = true
}
