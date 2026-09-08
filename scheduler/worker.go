package scheduler

import (
	"datacollector/device"
	"datacollector/device/messages"
	"datacollector/injestors"
	"datacollector/libs"
	"runtime"
)

type Worker struct {
	ID           int
	deque        *libs.Deque[*libs.Actor[device.DeviceState, messages.Message]]
	manager      *device.DeviceManager
	in           injestors.InjestorCount
	onStealActor func(int) (*libs.Actor[device.DeviceState, messages.Message], bool)
}

func NewWorker(
	id int,
	manager *device.DeviceManager,
	deque *libs.Deque[*libs.Actor[device.DeviceState, messages.Message]],
	injestor injestors.InjestorCount,
	onStealActor func(int) (*libs.Actor[device.DeviceState, messages.Message], bool),
) *Worker {
	return &Worker{
		ID:           id,
		deque:        deque,
		manager:      manager,
		in:           injestor,
		onStealActor: onStealActor,
	}
}

func (w *Worker) Run() {
	w.in.Run(w.submit)

	go w.update()
}

func (w *Worker) submit(t messages.Message) {
	d := w.manager.GetDevice(t.DeviceID)

	if hasToEnque := d.Send(t, int(t.Priority)); hasToEnque {
		w.deque.Push(d)
	}
}

func (w *Worker) update() {
	for {
		a, ok := w.deque.Pop()
		if !ok {
			a, ok = w.onStealActor(w.ID)
			if !ok {
				runtime.Gosched()
				continue
			}
		}

		if hasToEnque := a.Update(1); hasToEnque {
			w.deque.Push(a)
		}
	}
}
