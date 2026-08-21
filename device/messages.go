package device

import "fmt"

type Message interface {
	Apply(*DeviceState)
}

type Telemetry struct {
	DeviceID    string
	Temperature float64
	Humidity    float64
	Battery     float64
	Noise       float64
}

func (m Telemetry) Apply(state *DeviceState) {
	state.Data = m
	state.Online = true
}

type shutdown struct{}

func (m shutdown) Apply(state *DeviceState) {
	fmt.Println("shouting down the device ...")
}
