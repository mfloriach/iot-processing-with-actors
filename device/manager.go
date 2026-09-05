package device

import (
	"datacollector/device/messages"
	"datacollector/libs"
)

type Device = libs.Actor[DeviceState, messages.Message]

type DeviceManager struct {
	devices map[int]*Device
}

func NewDeviceManager(max_sensor_num int) *DeviceManager {
	return &DeviceManager{
		devices: make(map[int]*Device, max_sensor_num),
	}
}

func (m *DeviceManager) AddDevice(device *Device) {
	m.devices[device.GetID()] = device
}

func (m *DeviceManager) GetDevice(id int) *Device {
	return m.devices[id]
}

func (m *DeviceManager) DeleteDevice(id int) {
	delete(m.devices, id)
}
