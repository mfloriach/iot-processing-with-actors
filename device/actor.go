package device

import (
	"datacollector/mesures"
	"sync/atomic"
	"time"
)

type DeviceState struct {
	Data   Telemetry
	Online bool
}

type Actor struct {
	id       string
	mailbox  chan Telemetry
	state    DeviceState
	snapshot atomic.Value
	queued   bool
}

func NewActor(id string) *Actor {
	actor := &Actor{
		id:      id,
		mailbox: make(chan Telemetry, 5),
		queued:  false,
		state:   DeviceState{},
	}
	actor.snapshot.Store(actor.state)

	return actor
}

func (a *Actor) Update() bool {
	select {
	case msg := <-a.mailbox:
		// Publish a consistent snapshot for lock-free readers.
		a.snapshot.Store(DeviceState{
			Data:   msg,
			Online: true,
		})
		elapsed := time.Since(msg.TTL)
		mesures.Stats.Add(elapsed)
	default:
		a.queued = false
		return false
	}

	return true
}

func (a *Actor) Send(msg Telemetry) bool {
	a.mailbox <- msg
	mesures.Stats.AddGenerate()

	if a.queued {
		return false
	}

	a.queued = true
	return true
}

func (a *Actor) State() DeviceState {
	snapshot := a.snapshot.Load()
	if snapshot == nil {
		return DeviceState{}
	}

	return snapshot.(DeviceState)
}

func (a *Actor) GetID() string {
	return a.id
}
