package scheduler

import (
	"datacollector/device"
	"datacollector/injestor"
	"datacollector/mesures"
	"strconv"
	"time"
)

type Worker struct {
	ID      int
	deque   Deque
	in      injestor.Injestor
	manager *device.DeviceManager
}

func NewWorker(id int, manager *device.DeviceManager, injestor injestor.Injestor) *Worker {
	return &Worker{
		ID:      id,
		deque:   Deque{},
		in:      injestor,
		manager: manager,
	}
}

func (w *Worker) Run(start, end int) {
	go w.injestor(start, end)
	go w.update()
}

func (w *Worker) submit(t device.Telemetry) {
	d := w.manager.GetDevice(t.DeviceID)

	d.Send(t)
	w.deque.Push(d)
}

func (w *Worker) injestor(start, end int) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	count := 0

	for range ticker.C {
		for i := start; i < end; i++ {
			d := device.Telemetry{
				Sample: device.Sample{
					DeviceID: strconv.Itoa(i),
					TTL:      time.Now(),
				},

				Temperature: float64(count),
				Humidity:    float64(count),
				Battery:     float64(count),
				Noise:       float64(count),
			}

			w.submit(d)
			mesures.Stats.AddGenerate()
			count++
		}
	}
}

func (w *Worker) update() {
	for {
		a, ok := w.deque.Pop()
		if ok {
			a.Update(1)
		}
	}
}
