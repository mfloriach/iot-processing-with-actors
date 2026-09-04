package device

type DeviceManager struct {
	devices map[string]*Actor
}

func NewDeviceManager() DeviceManager {
	return DeviceManager{
		devices: make(map[string]*Actor, 30),
	}
}

func (m DeviceManager) AddDevice(device *Actor) {
	id := device.GetID()

	m.devices[id] = device
}

func (m DeviceManager) GetDevice(id string) *Actor {
	device, ok := m.devices[id]
	if !ok {
		device = NewActor(id)
		m.AddDevice(device)
	}
	return device
}
