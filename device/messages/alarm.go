package messages

import "time"

type AlarmType uint8

const (
	AlarmTemperatureHigh AlarmType = iota
	AlarmTemperatureLow
	AlarmOffline
	AlarmBatteryLow
	AlarmSensorFailure
)

type Severity uint8

const (
	SeverityInfo Severity = iota
	SeverityWarning
	SeverityCritical
)

type Alarm struct {
	Type      AlarmType
	Severity  Severity
	Message   string
	Timestamp time.Time
}
