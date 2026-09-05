package libs

import (
	"sync"
	"time"
)

type DequeHooks struct {
	OnBackpressure func(time.Duration)
}

type Deque[T any] struct {
	mu    sync.Mutex
	items []T
	hooks DequeHooks
}

func NewDeque[T any](size uint, hooks DequeHooks) *Deque[T] {
	return &Deque[T]{
		hooks: hooks,
		items: make([]T, 0, size),
	}
}

func (d *Deque[T]) Push(t T) {
	start := time.Now()

	d.mu.Lock()
	defer d.mu.Unlock()

	wait := time.Since(start)
	if wait > 0 && d.hooks.OnBackpressure != nil {
		d.hooks.OnBackpressure(wait)
	}

	d.items = append(d.items, t)
}

// Pop from the owner's side.
func (d *Deque[T]) Pop() (T, bool) {
	start := time.Now()

	d.mu.Lock()
	defer d.mu.Unlock()

	wait := time.Since(start)
	if wait > 0 && d.hooks.OnBackpressure != nil {
		d.hooks.OnBackpressure(wait)
	}

	n := len(d.items)
	if n == 0 {
		var zero T
		return zero, false
	}

	t := d.items[n-1]
	d.items = d.items[:n-1]

	return t, true
}

// Steal from the opposite side.
func (d *Deque[T]) Steal() (T, bool) {
	start := time.Now()

	d.mu.Lock()
	defer d.mu.Unlock()

	wait := time.Since(start)
	if wait > 0 && d.hooks.OnBackpressure != nil {
		d.hooks.OnBackpressure(wait)
	}

	if len(d.items) == 0 {
		var zero T
		return zero, false
	}

	t := d.items[0]
	d.items = d.items[1:]

	return t, true
}

func (d *Deque[T]) Len() int {
	d.mu.Lock()
	defer d.mu.Unlock()

	return len(d.items)
}
