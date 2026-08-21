package device

type setOnline struct {
	Online bool
}

func (setOnline) IsMessage() {}

type shutdown struct{}

func (shutdown) IsMessage() {}
