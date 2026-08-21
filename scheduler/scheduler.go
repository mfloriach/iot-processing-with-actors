package scheduler

import (
	"datacollector/device"
	"sync"
)

type Scheduler struct {
}

func NewScheduler() Scheduler {
	return Scheduler{}
}

func (s Scheduler) Run(manager device.DeviceManager) {

	var wg sync.WaitGroup

	for i := range 1 {
		wg.Add(1)

		go worker(
			i,
			manager,
			&wg,
		)
	}
}
