package libs

import (
	"sync/atomic"
	"time"
)

type ActorHooks struct {
	OnBackpressure func(time.Duration)
}

type Actor[S any, T any] struct {
	ID          int
	mailbox     chan T
	pool        *Pool[S]
	snapshot    atomic.Pointer[S]
	updateState func(*S, T)
	queued      bool
	hooks       ActorHooks
}

func NewActor[S any, T any](id int, mailbox_size int, pool *Pool[S], hooks ActorHooks, updateState func(*S, T)) *Actor[S, T] {
	actor := &Actor[S, T]{
		ID:          id,
		mailbox:     make(chan T, mailbox_size),
		pool:        pool,
		queued:      false,
		hooks:       hooks,
		updateState: updateState,
	}
	actor.snapshot.Store(new(S))

	return actor
}

func (a *Actor[S, T]) GetID() int {
	return a.ID
}

func (a *Actor[S, T]) Update(quantum int) (hasNext bool) {
	for i := 0; i < quantum; i++ {
		select {
		case msg := <-a.mailbox:
			state := a.pool.Get()

			a.updateState(state, msg)

			// Publish a consistent snapshot for lock-free readers.
			oldState := a.snapshot.Swap(state)
			a.pool.Put(oldState)

		default:
			a.queued = false
			return false
		}
	}

	return true
}

func (a *Actor[S, T]) Send(msg T) bool {
	select {
	case a.mailbox <- msg:
		// Sent immediately.
	default:
		// Mailbox is full.
		start := time.Now()
		a.mailbox <- msg
		if a.hooks.OnBackpressure != nil {
			a.hooks.OnBackpressure(time.Since(start))
		}
	}

	if a.queued {
		return false
	}

	a.queued = true
	return true
}

func (a *Actor[S, T]) State() *S {
	snapshot := a.snapshot.Load()
	if snapshot == nil {
		return new(S)
	}

	return snapshot
}
