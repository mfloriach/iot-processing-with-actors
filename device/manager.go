package device

import (
	"hash/fnv"
)

type DeviceActor interface {
	GetID() string
	Send(Message) bool
	State() DeviceState
	Update() bool
}

type shard struct {
	devices map[string]DeviceActor
	ready   chan DeviceActor
}

type DeviceManager struct {
	shards map[int]shard
}

func NewDeviceManager(numShards int) DeviceManager {
	shards := make(map[int]shard)
	for i := range numShards {
		shards[i] = shard{devices: make(map[string]DeviceActor), ready: make(chan DeviceActor, 1024)}
	}

	return DeviceManager{
		shards: shards,
	}
}

func (m DeviceManager) Add(device DeviceActor) {
	id := device.GetID()

	m.getShard(id).devices[id] = device
}

func (m DeviceManager) GetDevice(id string) DeviceActor {
	return m.getShard(id).devices[id]
}

func (m DeviceManager) Send(task Message) {
	id := task.GetDeviceID()

	device := m.getShard(id).devices[id]
	if device.Send(task) {
		m.getShard(id).ready <- device
	}
}

func (m DeviceManager) Process(shardID int) {
	ready := m.shards[shardID].ready

	for t := range ready {
		// Requeue the actor while it still has work so one busy mailbox does not
		// monopolize the shard and starve other devices.
		if t.Update() {
			ready <- t
		}
	}
}

func (m DeviceManager) getShard(id string) shard {
	h := fnv.New32a()
	_, _ = h.Write([]byte(id))

	return m.shards[int(h.Sum32()%uint32(len(m.shards)))]
}
