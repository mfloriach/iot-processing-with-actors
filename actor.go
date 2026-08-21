package main

import (
	"sync/atomic"
)

type actor[S any, M any] struct {
	id       string
	mailbox  chan M
	snapshot atomic.Value
	dispatch func(*S, M)
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

func (a *actor[S, M]) Update() {
	select {
	case msg := <-a.mailbox:
		state := a.snapshot.Load().(S)
		a.dispatch(&state, msg)
		a.snapshot.Store(state)
	default:
		return
	}
}

func (a *actor[S, M]) Send(msg M) {
	a.mailbox <- msg
}

func (a *actor[S, M]) State() S {
	return a.snapshot.Load().(S)
}
