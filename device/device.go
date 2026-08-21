package device

type DeviceState struct {
	Data   Telemetry
	Online bool
}

func (s *DeviceState) Dispatch(msg Message) {
	msg.Apply(s)
}
