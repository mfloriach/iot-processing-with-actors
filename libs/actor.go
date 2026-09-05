package libs

import (
	"sync/atomic"
	"time"
)

type ActorHooks[M any] struct {
	OnApply        func(M)
	OnBackpressure func(time.Duration)
}

type Actor[S any, M any] struct {
	ID       int
	mailbox  chan M
	pool     *Pool[S]
	snapshot atomic.Pointer[S]
	queued   bool
	hooks    ActorHooks[M]
	update   func(*S, M)
}

func NewActor[S any, M any](id int, mailbox_size int, pool *Pool[S], hooks ActorHooks[M], update func(*S, M)) *Actor[S, M] {
	actor := &Actor[S, M]{
		ID:      id,
		mailbox: make(chan M, mailbox_size),
		pool:    pool,
		queued:  false,
		hooks:   hooks,
		update:  update,
	}
	actor.snapshot.Store(pool.Get())

	return actor
}

func (a *Actor[S, M]) GetID() int {
	return a.ID
}

func (a *Actor[S, M]) Update(quantum int) (hasNext bool) {
	for i := 0; i < quantum; i++ {
		select {
		case msg := <-a.mailbox:
			state := a.pool.Get()
			a.update(state, msg)

			// Publish a consistent snapshot for lock-free readers.
			oldState := a.snapshot.Swap(state)
			a.pool.Put(oldState)

			if a.hooks.OnApply != nil {
				a.hooks.OnApply(msg)
			}

		default:
			a.queued = false
			return false
		}
	}

	return true
}

func (a *Actor[S, M]) Send(msg M) (hasToEnqueu bool) {
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

func (a *Actor[S, M]) State() *S {
	snapshot := a.snapshot.Load()
	if snapshot == nil {
		return new(S)
	}

	return snapshot
}
