package scheduler

import (
	"datacollector/device"
	"datacollector/device/messages"
	"datacollector/libs"
	"datacollector/mesures"
	"time"
)

type Worker struct {
	ID      int
	deque   libs.Deque[*libs.Actor[device.DeviceState, messages.Message]]
	manager *device.DeviceManager
}

func NewWorker(id int, manager *device.DeviceManager) *Worker {
	return &Worker{
		ID: id,
		deque: *libs.NewDeque[*libs.Actor[device.DeviceState, messages.Message]](
			100,
			libs.DequeHooks{},
		),
		manager: manager,
	}
}

func (w *Worker) Run(start, end int) {
	go w.injestorTelemetry(start, end)
	go w.injestorAlarm(12)
	go w.injestorCommand(12)
	go w.update()
}

func (w *Worker) submit(t messages.Message) {
	d := w.manager.GetDevice(t.DeviceID)

	if hasToEnque := d.Send(t); hasToEnque {
		w.deque.Push(d)
	}
}

func (w *Worker) injestorTelemetry(start, end int) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	count := 0

	for range ticker.C {
		for i := start; i < end; i++ {
			w.submit(messages.Message{
				Kind:     messages.MessageTelemetry,
				DeviceID: i,
				TTL:      time.Now(),
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

func (w *Worker) injestorCommand(deviceID int) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		w.submit(messages.Message{
			Kind:     messages.MessageCommand,
			DeviceID: deviceID,
			TTL:      time.Now(),
			Command: messages.Command{
				Type:  messages.CommandEnable,
				Value: 3.12,
			}})

		mesures.Stats.AddGenerate()
	}
}

func (w *Worker) injestorAlarm(deviceID int) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		w.submit(messages.Message{
			Kind:     messages.MessageAlarm,
			DeviceID: deviceID,
			TTL:      time.Now(),
			Alarm: messages.Alarm{
				Type:      messages.AlarmBatteryLow,
				Severity:  messages.SeverityWarning,
				Message:   "sdfdsfdfdsfs",
				Timestamp: time.Now(),
			}})

		mesures.Stats.AddGenerate()
	}
}

func (w *Worker) update() {
	for {
		a, ok := w.deque.Pop()
		if !ok {
			continue
		}

		if hasToEnque := a.Update(1); hasToEnque {
			w.deque.Push(a)
		}
	}
}
