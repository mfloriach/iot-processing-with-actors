package main

import (
	"datacollector/device"
	"datacollector/mesures"
)

type actor struct {
	id       string
	mailbox  chan device.Telemetry
	state    device.DeviceState
	dispatch func(*device.DeviceState, device.Telemetry)
	queued   bool
}

func NewActor(id string, initial device.DeviceState, dispatch func(*device.DeviceState, device.Telemetry)) *actor {
	actor := &actor{
		id:       id,
		mailbox:  make(chan device.Telemetry, 5),
		dispatch: dispatch,
		queued:   false,
		state:    device.DeviceState{},
	}

	return actor
}

func (a *actor) Update(quantum int) bool {
	for range quantum {
		select {
		case msg := <-a.mailbox:
			mesures.Processed.Add(1)

			a.dispatch(&a.state, msg)

		default:
			a.queued = false
			return false
		}
	}

	return true
}

func (a *actor) Send(msg device.Telemetry) bool {
	a.mailbox <- msg

	if a.queued {
		return false
	}

	a.queued = true
	return true
}

func (a *actor) State() device.DeviceState {
	return a.state
}

func (a *actor) GetID() string {
	return a.id
}
