package device

import (
	"datacollector/libs"
)

type DeviceManager struct {
	devices map[int]*libs.Actor[DeviceState, Telemetry]
}

func NewDeviceManager(max_sensor_num int) *DeviceManager {
	return &DeviceManager{
		devices: make(map[int]*libs.Actor[DeviceState, Telemetry], max_sensor_num),
	}
}

func (m *DeviceManager) AddDevice(device *libs.Actor[DeviceState, Telemetry]) {
	m.devices[device.GetID()] = device
}

func (m *DeviceManager) GetDevice(id int) *libs.Actor[DeviceState, Telemetry] {
	return m.devices[id]
}
