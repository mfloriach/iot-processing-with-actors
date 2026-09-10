package libs

import (
	"sync/atomic"
)

type ActorWithPriorityQueue[S any, M any] struct {
	ID       int
	mailbox  *Deque[M]
	pool     *Pool[S]
	snapshot atomic.Pointer[S]
	queued   bool
	hooks    ActorHooks[M]
	update   func(*S, M)
}

func NewActorWithPriorityQueue[S any, M any](id int,
	mailboxSize uint,
	pool *Pool[S],
	update func(*S, M),
	options ...ActorOption[M]) *ActorWithPriorityQueue[S, M] {

	var hooks ActorHooks[M]

	for _, option := range options {
		option(&hooks)
	}

	actor := &ActorWithPriorityQueue[S, M]{
		ID:      id,
		mailbox: NewDeque[M](1, int(mailboxSize)),
		pool:    pool,
		queued:  false,
		hooks:   hooks,
		update:  update,
	}
	actor.snapshot.Store(pool.Get())

	return actor
}

func (a *ActorWithPriorityQueue[S, M]) GetID() int {
	return a.ID
}

func (a *ActorWithPriorityQueue[S, M]) Update(quantum int) (hasNext bool) {
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

func (a *ActorWithPriorityQueue[S, M]) Send(msg M, priority int) (hasToEnqueu bool) {
	a.mailbox.Push(msg)

	if a.queued {
		return false
	}

	a.queued = true
	return true
}

func (a *ActorWithPriorityQueue[S, M]) State() *S {
	snapshot := a.snapshot.Load()
	if snapshot == nil {
		return new(S)
	}

	return snapshot
}
