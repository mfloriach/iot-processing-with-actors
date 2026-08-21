package device

type DeviceState struct {
	Data   Telemetry
	Online bool
}

func Dispatch(state *DeviceState, msg Message) {
	msg.Apply(state)
}
