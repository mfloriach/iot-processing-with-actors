package scheduler

import (
	"datacollector/device"
	"fmt"
	"sync"
)

func worker(
	id int,
	manager device.DeviceManager,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for {
		if err := manager.Store(id); err != nil {
			fmt.Println(err)
		}
	}
}
