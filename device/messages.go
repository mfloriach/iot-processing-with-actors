package device

import (
	"context"
	"log/slog"
	"time"
)

type Message interface {
	Apply(*DeviceState)
}

type Telemetry struct {
	DeviceID    string
	Temperature float64
	Humidity    float64
	Battery     float64
	Noise       float64

	Context context.Context `json:"-"`
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
	Battery float64

	Context context.Context `json:"-"`
}

func (m SetBattery) Apply(state *DeviceState) {
	state.Data.Battery = m.Battery
}
