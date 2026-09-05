package device

import (
	"datacollector/config"
	"datacollector/libs"
	"datacollector/mesures"
	"sync/atomic"
	"time"
)

var statePool = libs.NewPool(func() *DeviceState {
	return new(DeviceState)
})

type DeviceState struct {
	Data   Telemetry
	Online bool
}

type Actor struct {
	id       int
	mailbox  chan Telemetry
	snapshot atomic.Pointer[DeviceState]
	queued   bool
}

func NewActor(id int) *Actor {
	actor := &Actor{
		id:      id,
		mailbox: make(chan Telemetry, config.MAILBOX_SIZE),
		queued:  false,
	}
	actor.snapshot.Store(&DeviceState{})

	return actor
}

func (a *Actor) Update(quantum int) bool {
	for i := 0; i < quantum; i++ {
		select {
		case msg := <-a.mailbox:
			state := statePool.Get()

			state.Data = msg
			state.Online = true

			// Publish a consistent snapshot for lock-free readers.
			oldState := a.snapshot.Swap(state)
			statePool.Put(oldState)

			elapsed := time.Since(msg.TTL)
			mesures.Stats.Add(elapsed)
		default:
			a.queued = false
			return false
		}
	}

	return true
}

func (a *Actor) Send(msg Telemetry) bool {
	select {
	case a.mailbox <- msg:
		// Sent immediately.
	default:
		// Mailbox is full.
		start := time.Now()
		a.mailbox <- msg
		mesures.Stats.AddMailboxBackpressure(time.Since(start))
	}

	if a.queued {
		return false
	}

	a.queued = true
	return true
}

func (a *Actor) State() *DeviceState {
	snapshot := a.snapshot.Load()
	if snapshot == nil {
		return &DeviceState{}
	}

	return snapshot
}
