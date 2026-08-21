package device

import (
	"hash/fnv"
)

type DeviceManager struct {
	shards map[int]map[string]*actor[DeviceState, Message]
}

func NewDeviceManager(numShards int) DeviceManager {
	shards := make(map[int]map[string]*actor[DeviceState, Message])
	for i := range numShards {
		shards[i] = make(map[string]*actor[DeviceState, Message])
	}

	return DeviceManager{
		shards: shards,
	}
}

func (m DeviceManager) Add(id string) {
	s := m.shardIndex(id)
	m.shards[s][id] = newActor(id, DeviceState{}, Dispatch)
}

func (m DeviceManager) State(id string) DeviceState {
	s := m.shardIndex(id)
	return m.shards[s][id].State()
}

func (m DeviceManager) Send(id string, task Message) {
	s := m.shardIndex(id)
	m.shards[s][id].Send(task)
}

func (m DeviceManager) Store(shardID int) {
	// TODO shardID no in range
	for _, d := range m.shards[shardID] {
		d.Store()
	}
}

func (m DeviceManager) ShutDown(id string) {
	s := m.shardIndex(id)
	m.shards[s][id].Send(shutdown{})
}

func (m DeviceManager) shardIndex(id string) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(id))

	return int(h.Sum32() % uint32(len(m.shards)))
}
