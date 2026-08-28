package device

type DeviceState struct {
	Data   Telemetry
	Online bool
}

func Dispatch(state *DeviceState, msg Telemetry) {
	msg.Apply(state)
}
