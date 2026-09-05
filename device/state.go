package device

import (
	"datacollector/device/messages"
)

type DeviceState struct {
	Data   messages.Telemetry
	Online bool
}

func Apply(state *DeviceState, msg messages.Message) {
	switch msg.Kind {
	case messages.MessageTelemetry:
		state.Data = msg.Telemetry
		state.Online = true

		// case messages.MessageCommand:
		// 	fmt.Println("command")

		// case messages.MessageAlarm:
		// 	fmt.Println("alamr")
	}
}
