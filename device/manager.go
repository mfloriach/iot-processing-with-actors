package device

import (
	"datacollector/config"
	"datacollector/mesures"
	"iter"
	"time"
)

type DeviceManager struct {
	devices map[string]*Actor
	ready   chan *Actor
}

func NewDeviceManager() DeviceManager {
	return DeviceManager{
		devices: make(map[string]*Actor, 30),
		ready:   make(chan *Actor, 30),
	}
}

func (m DeviceManager) AddDevice(device *Actor) {
	id := device.GetID()

	m.devices[id] = device
}

func (m DeviceManager) GetDevice(id string) *Actor {
	return m.devices[id]
}

func (m DeviceManager) Send(task Telemetry) {
	device := m.devices[task.DeviceID]
	if device.Send(task) {
		select {
		case m.ready <- device:
		default:
			// ready is full.
			start := time.Now()
			m.ready <- device
			mesures.Stats.AddReadyBackpressure(time.Since(start))
		}
	}
}

func (m DeviceManager) Next() iter.Seq[*Actor] {
	return func(yield func(*Actor) bool) {
		for t := range m.ready {
			if !yield(t) {
				break
			}
		}
	}
}

func (m DeviceManager) ProcessOne(t *Actor) (hasMore bool) {
	// Requeue the actor while it still has work so one busy mailbox does not
	// monopolize the shard and starve other devices.
	if t.Update(config.QUANTUM) {
		select {
		case m.ready <- t:
		default:
			// ready is full.
			start := time.Now()
			m.ready <- t
			mesures.Stats.AddReadyBackpressure(time.Since(start))
		}
		return true
	}

	return false
}
