package device

type DeviceState struct {
	Data            Telemetry
	Online          bool
	FirmwareVersion string
}

type Message interface {
	Apply(*DeviceState)
}

func (s *DeviceState) Dispatch(msg Message) {
	msg.Apply(s)
}
