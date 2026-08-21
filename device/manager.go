package device

type DeviceManager struct {
	devices map[string]*actor
}

func NewDeviceManager() DeviceManager {
	return DeviceManager{
		devices: make(map[string]*actor),
	}
}

func (m DeviceManager) Add(id string) {
	m.devices[id] = newActor(id)
}

func (m DeviceManager) GetState(id string) DeviceState {
	return m.devices[id].State()
}

func (m DeviceManager) Send(id string, task Message) {
	m.devices[id].Send(task)
}

func (m DeviceManager) Store() {
	for _, d := range m.devices {
		d.Store()
	}
}

func (m DeviceManager) ShutDown(id string) {
	m.devices[id].Send(shutdown{})
}
