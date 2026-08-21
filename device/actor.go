package device

import (
	"sync/atomic"
)

type deviceActor struct {
	id       string
	mailbox  chan Message
	snapshot atomic.Value
}

func newDeviceActor(id string) *deviceActor {
	actor := &deviceActor{
		id:      id,
		mailbox: make(chan Message, 100),
	}

	actor.snapshot.Store(DeviceState{})

	go actor.run()

	return actor
}

func (a *deviceActor) run() {
	state := DeviceState{}

	for msg := range a.mailbox {
		switch msg := msg.(type) {

		case Telemetry:
			state.Data = msg
			state.Online = true

			a.snapshot.Store(state)

		case setOnline:
			state.Online = msg.Online
			a.snapshot.Store(state)

		case shutdown:
			return
		}
	}
}

func (a *deviceActor) Send(msg Message) {
	a.mailbox <- msg
}

func (a *deviceActor) State() DeviceState {
	return a.snapshot.Load().(DeviceState)
}
