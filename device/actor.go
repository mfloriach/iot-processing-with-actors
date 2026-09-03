package device

import (
	"datacollector/mesures"
	"sync"
	"sync/atomic"
	"time"
)

var statePool = sync.Pool{
	New: func() any {
		return new(DeviceState)
	},
}

type DeviceState struct {
	Data   Telemetry
	Online bool
}

type Actor struct {
	id       string
	mailbox  chan Telemetry
	snapshot atomic.Pointer[DeviceState]
	queued   bool
}

func NewActor(id string) *Actor {
	actor := &Actor{
		id:      id,
		mailbox: make(chan Telemetry, 100000),
		queued:  false,
	}
	actor.snapshot.Store(&DeviceState{})

	return actor
}

func (a *Actor) Update() bool {
	select {
	case msg := <-a.mailbox:
		state := statePool.Get().(*DeviceState)

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

	return true
}

func (a *Actor) Send(msg Telemetry) bool {
	select {
	case a.mailbox <- msg:
		// Sent immediately.
	default:
		// // Mailbox is full.
		// slog.Debug("backpressure",
		// 	"queue", "mailbox",
		// 	"device", a.id,
		// 	"len", len(a.mailbox),
		// 	"cap", cap(a.mailbox),
		// )

		a.mailbox <- msg
	}

	mesures.Stats.AddGenerate()

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

func (a *Actor) GetID() string {
	return a.id
}
