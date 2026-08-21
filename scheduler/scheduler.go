package scheduler

import (
	"datacollector/device"
	"iter"
	"sync"
)

type Scheduler struct {
	injestor func() iter.Seq[device.Telemetry]
}

func NewScheduler(injestor func() iter.Seq[device.Telemetry]) Scheduler {
	return Scheduler{injestor: injestor}
}

func (s Scheduler) Run(handler func(id string, task device.Message)) {
	jobs := make(chan *device.Telemetry, 100)

	var wg sync.WaitGroup

	for i := range 2 {
		wg.Add(1)

		go worker(
			i,
			handler,
			jobs,
			&wg,
		)
	}

	for t := range s.injestor() {
		jobs <- &t
	}

	close(jobs)
}
