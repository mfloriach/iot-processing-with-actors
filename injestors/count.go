package injestors

import (
	"datacollector/device/messages"
	"datacollector/mesures"
	"time"
)

type InjestorCount struct {
	start, end int
	submit     func(messages.Message)
}

func NewNewInjestor(start, end int) InjestorCount {
	return InjestorCount{start: start, end: end}
}

func (ij *InjestorCount) Run(submit func(messages.Message)) {
	ij.submit = submit

	go ij.injestorTelemetry(ij.start, ij.end)
	go ij.injestorAlarm(12)
	go ij.injestorCommand(12)
}

func (ij *InjestorCount) injestorTelemetry(start, end int) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	count := 0

	for range ticker.C {
		for i := start; i < end; i++ {
			ij.submit(messages.Message{
				Kind:     messages.MessageTelemetry,
				DeviceID: i,
				TTL:      time.Now(),
				Priority: messages.Priority(messages.PriorityLow),
				Telemetry: messages.Telemetry{
					Temperature: float64(count),
					Humidity:    float64(count),
					Battery:     float64(count),
					Noise:       float64(count),
				}})

			mesures.Stats.AddGenerate()
		}
		count++
	}
}

func (ij *InjestorCount) injestorCommand(deviceID int) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		ij.submit(messages.Message{
			Kind:     messages.MessageCommand,
			Priority: messages.Priority(messages.PriorityCritical),
			DeviceID: deviceID,
			TTL:      time.Now(),
			Command: messages.Command{
				Type:  messages.CommandEnable,
				Value: 3.12,
			}})

		mesures.Stats.AddGenerate()
	}
}

func (ij *InjestorCount) injestorAlarm(deviceID int) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		ij.submit(messages.Message{
			Kind:     messages.MessageAlarm,
			DeviceID: deviceID,
			TTL:      time.Now(),
			Priority: messages.Priority(messages.PriorityCritical),
			Alarm: messages.Alarm{
				Type:      messages.AlarmBatteryLow,
				Severity:  messages.SeverityWarning,
				Message:   "sdfdsfdfdsfs",
				Timestamp: time.Now(),
			}})

		mesures.Stats.AddGenerate()
	}
}
