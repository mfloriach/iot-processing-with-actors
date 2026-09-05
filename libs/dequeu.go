package libs

import (
	"fmt"
	"sync"
)

type DequeHooks struct {
	OnBackpressure func(int)
}

type DequeOption[T any] func(*Deque[T])

func WithOnBackpressure[T any](fn func(int)) DequeOption[T] {
	return func(d *Deque[T]) {
		d.hooks.OnBackpressure = fn
	}
}

type Deque[T any] struct {
	ID    int
	mu    sync.Mutex
	items []T
	hooks DequeHooks
}

func NewDeque[T any](id, size int, options ...DequeOption[T]) *Deque[T] {
	d := &Deque[T]{
		ID:    id,
		items: make([]T, 0, size),
	}

	for _, option := range options {
		option(d)
	}

	return d
}

func (d *Deque[T]) Push(t T) {
	// start := time.Now()

	d.mu.Lock()
	defer d.mu.Unlock()

	// wait := time.Since(start)
	// if wait > 0 && d.hooks.OnBackpressure != nil {
	// 	d.hooks.OnBackpressure(d.ID)
	// }

	d.items = append(d.items, t)
}

// Pop from the owner's side.
func (d *Deque[T]) Pop() (T, bool) {
	// start := time.Now()

	d.mu.Lock()
	defer d.mu.Unlock()

	// wait := time.Since(start)
	// if wait > 0 && d.hooks.OnBackpressure != nil {
	// 	d.hooks.OnBackpressure(d.ID)
	// }

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
	// start := time.Now()

	d.mu.Lock()
	defer d.mu.Unlock()

	// wait := time.Since(start)
	// if wait > 0 && d.hooks.OnBackpressure != nil {
	// 	d.hooks.OnBackpressure(d.ID)
	// }

	if len(d.items) == 0 {
		fmt.Println(d.ID)
		var zero T
		return zero, false
	}

	fmt.Println("____________________")

	fmt.Println("I can not steal")

	t := d.items[0]
	d.items = d.items[1:]

	return t, true
}

func (d *Deque[T]) Len() int {
	d.mu.Lock()
	defer d.mu.Unlock()

	return len(d.items)
}
