package device

import (
	"hash/fnv"
	"iter"
)

type DeviceActor interface {
	Send(Message)
	State() DeviceState
	Update()
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
		shards[i] = shard{devices: make(map[string]DeviceActor), ready: make(chan DeviceActor)}
	}

	return DeviceManager{
		shards: shards,
	}
}

func (m DeviceManager) Add(id string, device DeviceActor) {
	m.getShard(id).devices[id] = device
}

func (m DeviceManager) State(id string) DeviceState {
	return m.getShard(id).devices[id].State()
}

func (m DeviceManager) Send(id string, task Message) {
	deviceID := m.getShard(id).devices[id]

	deviceID.Send(task)
	m.getShard(id).ready <- deviceID
}

func (m DeviceManager) Next(shardID int) iter.Seq[DeviceActor] {
	if shardID < 0 || shardID > len(m.shards) {
		// return errors.New("shard numer does not exist")
	}

	ready := m.shards[shardID].ready

	return func(yield func(DeviceActor) bool) {
		for {
			if !yield(<-ready) {
				return
			}
		}
	}
}

func (m DeviceManager) getShard(id string) shard {
	h := fnv.New32a()
	_, _ = h.Write([]byte(id))

	return m.shards[int(h.Sum32()%uint32(len(m.shards)))]
}
