package scheduler

import (
	"datacollector/config"
	"datacollector/device"
	"datacollector/injestor"
)

type Scheduler struct {
	injestor injestor.Injestor
	manager  *device.DeviceManager
}

func NewScheduler(manager *device.DeviceManager, injestor injestor.Injestor) Scheduler {
	return Scheduler{
		injestor: injestor,
		manager:  manager,
	}
}

func (s Scheduler) Run() {
	for i := 0; i < config.NUM_OF_WORKERS; i++ {
		w := NewWorker(i, s.manager, s.injestor)
		go w.Run(i*config.SENSOR_PER_WORKER, ((i+1)*config.SENSOR_PER_WORKER)-1)
	}
}
