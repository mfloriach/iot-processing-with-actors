package device

import "datacollector/config"

type DeviceManager struct {
	devices map[int]*Actor
}

func NewDeviceManager() *DeviceManager {
	return &DeviceManager{
		devices: make(map[int]*Actor, config.SENSOR_PER_WORKER*config.NUM_CPUS),
	}
}

func (m *DeviceManager) AddDevice(device *Actor) {
	m.devices[device.id] = device
}

func (m *DeviceManager) GetDevice(id int) *Actor {
	return m.devices[id]
}
