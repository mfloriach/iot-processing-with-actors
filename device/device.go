package device

type DeviceState struct {
	Data   Telemetry
	Online bool
}

type Message interface {
	IsMessage()
}
