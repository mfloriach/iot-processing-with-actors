package messages

import (
	"time"
)

type MessageKind uint8

const (
	MessageTelemetry MessageKind = iota
	MessageCommand
	MessageAlarm
)

type Message struct {
	DeviceID  int
	TTL       time.Time `json:"-"`
	Kind      MessageKind
	Telemetry Telemetry
	Command   Command
	Alarm     Alarm
}
