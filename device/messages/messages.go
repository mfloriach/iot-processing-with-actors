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

type Priority uint8

const (
	PriorityCritical MessageKind = iota
	PriorityNormal
	PriorityLow
)

type Message struct {
	DeviceID  int
	TTL       time.Time `json:"-"`
	Priority  Priority
	Kind      MessageKind
	Telemetry Telemetry
	Command   Command
	Alarm     Alarm
}
