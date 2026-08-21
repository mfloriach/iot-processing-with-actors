package device

import "fmt"

type setOnline struct {
	Online bool
}

func (m setOnline) Apply(state *DeviceState) {
	state.Online = m.Online
}

type shutdown struct{}

func (m shutdown) Apply(state *DeviceState) {
	fmt.Println("shouting down the device ...")
}
