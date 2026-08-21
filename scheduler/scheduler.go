package scheduler

import (
	"datacollector/device"
	"sync"
)

type Scheduler struct {
	numOfWorks int
}

func NewScheduler(numOfWorks int) Scheduler {
	return Scheduler{numOfWorks: numOfWorks}
}

func (s Scheduler) Run(manager device.DeviceManager) {
	var wg sync.WaitGroup

	for i := range s.numOfWorks {
		wg.Add(1)

		go worker(
			i,
			manager,
			&wg,
		)
	}
}
