package scheduler

import (
	"datacollector/config"
	"datacollector/device"
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
		w := NewWorker(i, s.manager)
		s.workers = append(s.workers, w)
		go w.Run(i*config.SENSOR_PER_WORKER, ((i+1)*config.SENSOR_PER_WORKER)-1)
	}
}
