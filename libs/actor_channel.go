package libs

import (
	"sync/atomic"
)

type ActorWithChannel[S any, M any] struct {
	ID       int
	mailbox  chan M
	pool     *Pool[S]
	snapshot atomic.Pointer[S]
	queued   bool
	hooks    ActorHooks[M]
	update   func(*S, M)
}

func NewActorWithChannel[S any, M any](id int,
	mailboxSize uint,
	pool *Pool[S],
	update func(*S, M),
	options ...ActorOption[M]) *ActorWithChannel[S, M] {

	var hooks ActorHooks[M]

	for _, option := range options {
		option(&hooks)
	}

	actor := &ActorWithChannel[S, M]{
		ID:      id,
		mailbox: make(chan M, mailboxSize),
		pool:    pool,
		queued:  false,
		hooks:   hooks,
		update:  update,
	}
	actor.snapshot.Store(pool.Get())

	return actor
}

func (a *ActorWithChannel[S, M]) GetID() int {
	return a.ID
}

func (a *ActorWithChannel[S, M]) Update(quantum int) (hasNext bool) {
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
			return false
		}
	}

	return true
}

func (a *ActorWithChannel[S, M]) Send(msg M, priority int) (hasToEnqueu bool) {
	select {
	case a.mailbox <- msg:
	default:
		// Queue is full.
	}

	if a.queued {
		return false
	}
	a.queued = true
	return true
}

func (a *ActorWithChannel[S, M]) State() *S {
	snapshot := a.snapshot.Load()
	if snapshot == nil {
		return new(S)
	}

	return snapshot
}
