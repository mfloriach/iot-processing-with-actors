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

type Shutdown struct{}

func (m Shutdown) Apply(state *DeviceState) {
	fmt.Println("shouting down the device ...")
}

type SetBattery struct {
	Battery float64
}

func (m SetBattery) Apply(state *DeviceState) {
	state.Data.Battery = m.Battery
}
