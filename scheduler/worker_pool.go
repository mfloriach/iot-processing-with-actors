package scheduler

import (
	"datacollector/device"
	"sync"
)

func worker(
	id int,
	manager device.DeviceManager,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for {
		manager.Store(id)
	}
}
