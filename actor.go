package main

import (
	"datacollector/device"
	"datacollector/mesures"
	"sync/atomic"
	"time"
)

type actor struct {
	id       string
	mailbox  chan device.Telemetry
	state    device.DeviceState
	snapshot atomic.Value
	dispatch func(*device.DeviceState, device.Telemetry)
	queued   bool
}

func NewActor(id string, dispatch func(*device.DeviceState, device.Telemetry)) *actor {
	actor := &actor{
		id:       id,
		mailbox:  make(chan device.Telemetry, 5),
		dispatch: dispatch,
		queued:   false,
		state:    device.DeviceState{},
	}
	actor.snapshot.Store(actor.state)

	return actor
}

func (a *actor) Update() bool {
	// for range quantum {
	select {
	case msg := <-a.mailbox:
		a.dispatch(&a.state, msg)
		// Publish a consistent snapshot for lock-free readers.
		a.snapshot.Store(a.state)
		elapsed := time.Since(msg.TTL)
		mesures.Stats.Add(elapsed)
	default:
		a.queued = false
		return false
	}
	// }

	return true
}

func (a *actor) Send(msg device.Telemetry) bool {
	a.mailbox <- msg
	mesures.Stats.AddGenerate()

	if a.queued {
		return false
	}

	a.queued = true
	return true
}

func (a *actor) State() device.DeviceState {
	snapshot := a.snapshot.Load()
	if snapshot == nil {
		return device.DeviceState{}
	}

	return snapshot.(device.DeviceState)
}

func (a *actor) GetID() string {
	return a.id
}
