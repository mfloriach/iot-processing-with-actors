package device

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

type Message interface {
	GetDeviceID() string
	Apply(*DeviceState)
}

type Sample struct {
	DeviceID string
	Context  context.Context `json:"-"`
}

type Telemetry struct {
	Sample

	Temperature float64
	Humidity    float64
	Battery     float64
	Noise       float64
}

func (m Telemetry) GetDeviceID() string {
	return m.DeviceID
}

func (m Telemetry) Apply(state *DeviceState) {
	select {
	case <-time.After(time.Second):
		state.Data = m
		state.Online = true
	case <-m.Context.Done():
		slog.Error("error on update device", slog.Any("error", m.Context.Err()))
	}
}

type Shutdown struct{}

func (m Shutdown) Apply(state *DeviceState) {
	slog.Info("shouting down the device ...")
}

type SetBattery struct {
	Sample

	Battery float64
}

func (m SetBattery) Apply(state *DeviceState) {
	state.Data.Battery = m.Battery
}

type SetAlarm struct {
	Sample

	Value float64
}

func (m SetAlarm) GetDeviceID() string {
	return m.DeviceID
}

func (m SetAlarm) Apply(state *DeviceState) {
	fmt.Println("Alarm triggered.....")
}
