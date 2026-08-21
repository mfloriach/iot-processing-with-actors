package device

import (
	"sync/atomic"
)

type actor struct {
	id       string
	mailbox  chan Message
	snapshot atomic.Value
}

func newActor(id string) *actor {
	actor := &actor{
		id:      id,
		mailbox: make(chan Message, 100),
	}

	actor.snapshot.Store(DeviceState{})

	return actor
}

func (a *actor) Store() {
	state := DeviceState{}

	select {
	case msg := <-a.mailbox:
		state.Dispatch(msg)
		a.snapshot.Store(state)
	default:
		return
	}
}

func (a *actor) Send(msg Message) {
	a.mailbox <- msg
}

func (a *actor) State() DeviceState {
	return a.snapshot.Load().(DeviceState)
}
