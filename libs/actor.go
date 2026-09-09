package libs

import (
	"sync/atomic"
	"time"
)

type ActorHooks[M any] struct {
	OnApply        func(M)
	OnBackpressure func(time.Duration)
}

type ActorOption[M any] func(*ActorHooks[M])

func ActorWithOnApply[M any](fn func(M)) ActorOption[M] {
	return func(h *ActorHooks[M]) {
		h.OnApply = fn
	}
}

func ActorWithOnBackpressure[M any](fn func(time.Duration)) ActorOption[M] {
	return func(h *ActorHooks[M]) {
		h.OnBackpressure = fn
	}
}

type Actor[S any, M any] struct {
	ID       int
	mailbox  Mailbox[M]
	pool     *Pool[S]
	snapshot atomic.Pointer[S]
	queued   bool
	hooks    ActorHooks[M]
	update   func(*S, M)
}

func NewActorWithPriorityQueues[S any, M any](id int,
	mailboxSize uint,
	pool *Pool[S],
	update func(*S, M),
	options ...ActorOption[M]) *Actor[S, M] {

	var hooks ActorHooks[M]

	for _, option := range options {
		option(&hooks)
	}

	actor := &Actor[S, M]{
		ID:      id,
		mailbox: NewMailbox[M](mailboxSize),
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
		msg, _ := a.mailbox.Pop()

		state := a.pool.Get()
		a.update(state, msg)

		// Publish a consistent snapshot for lock-free readers.
		oldState := a.snapshot.Swap(state)
		a.pool.Put(oldState)

		if a.hooks.OnApply != nil {
			a.hooks.OnApply(msg)
		}

		if a.mailbox.Len() == 0 {
			a.queued = false
			return false
		}
	}

	return true
}

func (a *Actor[S, M]) Send(msg M, priority int) (hasToEnqueu bool) {
	a.mailbox.Push(msg, priority)

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
