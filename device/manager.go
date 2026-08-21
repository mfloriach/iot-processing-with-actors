package device

type DeviceManager struct {
	devices map[string]*deviceActor
}

func NewDeviceManager() DeviceManager {
	return DeviceManager{
		devices: make(map[string]*deviceActor),
	}
}

func (m DeviceManager) Add(id string) {
	m.devices[id] = newDeviceActor(id)
}

func (m DeviceManager) GetState(id string) DeviceState {
	return m.devices[id].State()
}

func (m DeviceManager) Handler(id string, task Message) {
	m.devices[id].Send(task)
}
