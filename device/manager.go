package device

import "datacollector/mesures"

type DeviceActor interface {
	GetID() string
	Send(Telemetry) bool
	State() DeviceState
	Update(int) bool
}

type DeviceManager struct {
	devices map[string]DeviceActor
	ready   chan DeviceActor
	quantum int
}

func NewDeviceManager(quantum int) DeviceManager {
	return DeviceManager{
		devices: make(map[string]DeviceActor, 50),
		ready:   make(chan DeviceActor, 10),
		quantum: quantum,
	}
}

func (m DeviceManager) Add(device DeviceActor) {
	id := device.GetID()

	m.devices[id] = device
}

func (m DeviceManager) GetDevice(id string) DeviceActor {
	return m.devices[id]
}

func (m DeviceManager) Send(task Telemetry) {
	mesures.Stats.AddGenerate()
	id := task.GetDeviceID()

	device := m.devices[id]
	if device.Send(task) {
		m.ready <- device
	}
}

func (m DeviceManager) Process() {
	for t := range m.ready {
		// Requeue the actor while it still has work so one busy mailbox does not
		// monopolize the shard and starve other devices.
		if t.Update(m.quantum) {
			m.ready <- t
		}
	}
}
