package device

import (
	"iter"
)

type DeviceManager struct {
	devices map[string]*Actor
	ready   chan *Actor
}

func NewDeviceManager() DeviceManager {
	return DeviceManager{
		devices: make(map[string]*Actor, 50),
		ready:   make(chan *Actor, 10),
	}
}

func (m DeviceManager) AddDevice(device *Actor) {
	id := device.GetID()

	m.devices[id] = device
}

func (m DeviceManager) GetDevice(id string) *Actor {
	return m.devices[id]
}

func (m DeviceManager) Send(task Telemetry) {
	id := task.GetDeviceID()

	device := m.devices[id]
	if device.Send(task) {
		m.ready <- device
	}
}

func (m DeviceManager) Next() iter.Seq[*Actor] {
	return func(yield func(*Actor) bool) {
		for t := range m.ready {
			if !yield(t) {
				break
			}
		}
	}
}

func (m DeviceManager) ProcessOne(t *Actor) (hasMore bool) {
	// Requeue the actor while it still has work so one busy mailbox does not
	// monopolize the shard and starve other devices.
	if t.Update() {
		m.ready <- t
		return true
	}

	return false
}
