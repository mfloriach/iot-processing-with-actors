package device

import "iter"

type DeviceActor interface {
	GetID() string
	Send(Telemetry) bool
	State() DeviceState
	Update() bool
}

type DeviceManager struct {
	devices map[string]DeviceActor
	ready   chan DeviceActor
}

func NewDeviceManager() DeviceManager {
	return DeviceManager{
		devices: make(map[string]DeviceActor, 50),
		ready:   make(chan DeviceActor, 10),
	}
}

func (m DeviceManager) AddDevice(device DeviceActor) {
	id := device.GetID()

	m.devices[id] = device
}

func (m DeviceManager) GetDevice(id string) DeviceActor {
	return m.devices[id]
}

func (m DeviceManager) Send(task Telemetry) {
	id := task.GetDeviceID()

	device := m.devices[id]
	if device.Send(task) {
		m.ready <- device
	}
}

func (m DeviceManager) Next() iter.Seq[DeviceActor] {
	return func(yield func(DeviceActor) bool) {
		for t := range m.ready {
			if !yield(t) {
				break
			}
		}
	}
}

func (m DeviceManager) ProcessOne(t DeviceActor) (hasMore bool) {
	// Requeue the actor while it still has work so one busy mailbox does not
	// monopolize the shard and starve other devices.
	if t.Update() {
		m.ready <- t
		return true
	}

	return false
}
