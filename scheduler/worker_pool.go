package scheduler

import (
	"datacollector/device"
	"fmt"
	"sync"
)

func worker(
	id int,
	handler func(id string, task device.Message),
	jobs <-chan *device.Telemetry,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for j := range jobs {
		handler(j.DeviceID, *j)

		fmt.Printf(
			"worker=%d processing device=%s\n",
			id,
			j.DeviceID,
		)
	}
}
