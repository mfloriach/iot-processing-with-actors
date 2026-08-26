package main

import (
	"sync"
	"sync/atomic"
)

type actor[S any, M any] struct {
	id       string
	mailbox  chan M
	snapshot atomic.Value
	dispatch func(*S, M)
	mu       sync.Mutex
	queued   bool
}

func NewActor[S any, M any](id string, initial S, dispatch func(*S, M)) *actor[S, M] {
	actor := &actor[S, M]{
		id:       id,
		mailbox:  make(chan M, 100),
		dispatch: dispatch,
	}

	actor.snapshot.Store(initial)

	return actor
}

func (a *actor[S, M]) Update() bool {
	msg := <-a.mailbox

	state := a.snapshot.Load().(S)
	a.dispatch(&state, msg)
	a.snapshot.Store(state)

	a.mu.Lock()
	defer a.mu.Unlock()

	if len(a.mailbox) == 0 {
		a.queued = false
		return false
	}

	return true
}

func (a *actor[S, M]) Send(msg M) bool {
	a.mailbox <- msg

	a.mu.Lock()
	defer a.mu.Unlock()

	if a.queued {
		return false
	}

	a.queued = true
	return true
}

func (a *actor[S, M]) State() S {
	return a.snapshot.Load().(S)
}

func (a *actor[S, M]) GetID() string {
	return a.id
}
