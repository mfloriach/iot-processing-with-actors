package device

type Telemetry struct {
	DeviceID    string
	Temperature float64
	Humidity    float64
	Battery     float64
	Noise       float64
}

func (Telemetry) IsMessage() {}
