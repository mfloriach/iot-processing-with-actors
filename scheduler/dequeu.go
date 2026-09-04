package scheduler

import (
	"datacollector/device"
	"sync"
)

type Deque struct {
	mu    sync.Mutex
	items []*device.Actor
}

func (d *Deque) Push(t *device.Actor) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.items = append(d.items, t)
}

// Pop from the owner's side.
func (d *Deque) Pop() (*device.Actor, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	n := len(d.items)
	if n == 0 {
		return nil, false
	}

	t := d.items[n-1]
	d.items = d.items[:n-1]

	return t, true
}

// Steal from the opposite side.
func (d *Deque) Steal() (*device.Actor, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if len(d.items) == 0 {
		return nil, false
	}

	t := d.items[0]
	d.items = d.items[1:]

	return t, true
}

func (d *Deque) Len() int {
	d.mu.Lock()
	defer d.mu.Unlock()

	return len(d.items)
}
