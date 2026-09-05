package scheduler

import (
	"datacollector/config"
	"datacollector/device"
	"datacollector/device/messages"
	"datacollector/injestors"
	"datacollector/libs"
)

type Scheduler struct {
	manager *device.DeviceManager
	workers []*Worker
}

func NewScheduler(manager *device.DeviceManager) Scheduler {
	return Scheduler{
		manager: manager,
		workers: make([]*Worker, 0, config.NUM_CPUS),
	}
}

func (s *Scheduler) Run() {
	for i := 0; i < config.NUM_CPUS; i++ {
		deque := *libs.NewDeque[*libs.Actor[device.DeviceState, messages.Message]](
			i,
			100,
		)

		start := i * config.SENSOR_PER_WORKER
		end := ((i + 1) * config.SENSOR_PER_WORKER) - 1
		in := injestors.NewNewInjestor(start, end)

		w := NewWorker(i, s.manager, &deque, in, s.onStealActor)
		s.workers = append(s.workers, w)
		go w.Run()
	}
}

func (s *Scheduler) onStealActor(workerID int) (*libs.Actor[device.DeviceState, messages.Message], bool) {
	n := len(s.workers)

	for i := 0; i < n; i++ {
		victim := s.workers[(workerID+i+1)%n]

		actor, ok := victim.deque.Steal()
		if ok {
			return actor, true
		}
	}

	return nil, false
}
