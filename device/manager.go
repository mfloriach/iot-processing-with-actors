package device

import "datacollector/config"

type DeviceManager struct {
	devices map[int]*Actor[DeviceState]
}

func NewDeviceManager() *DeviceManager {
	return &DeviceManager{
		devices: make(map[int]*Actor[DeviceState], config.SENSOR_PER_WORKER*config.NUM_CPUS),
	}
}

func (m *DeviceManager) AddDevice(device *Actor[DeviceState]) {
	m.devices[device.id] = device
}

func (m *DeviceManager) GetDevice(id int) *Actor[DeviceState] {
	return m.devices[id]
}
