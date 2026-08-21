package device

import (
	"errors"
	"hash/fnv"
)

type DeviceActor interface {
	Send(Message)
	State() DeviceState
	Store()
}

type DeviceManager struct {
	shards map[int]map[string]DeviceActor
}

func NewDeviceManager(numShards int) DeviceManager {
	shards := make(map[int]map[string]DeviceActor)
	for i := range numShards {
		shards[i] = make(map[string]DeviceActor)
	}

	return DeviceManager{
		shards: shards,
	}
}

func (m DeviceManager) Add(id string, device DeviceActor) {
	s := m.shardIndex(id)
	m.shards[s][id] = device
}

func (m DeviceManager) State(id string) DeviceState {
	s := m.shardIndex(id)
	return m.shards[s][id].State()
}

func (m DeviceManager) Send(id string, task Message) {
	s := m.shardIndex(id)
	m.shards[s][id].Send(task)
}

func (m DeviceManager) Store(shardID int) error {
	if shardID < 0 || shardID > len(m.shards) {
		return errors.New("shard numer does not exist")
	}

	for _, d := range m.shards[shardID] {
		d.Store()
	}

	return nil
}

func (m DeviceManager) shardIndex(id string) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(id))

	return int(h.Sum32() % uint32(len(m.shards)))
}
