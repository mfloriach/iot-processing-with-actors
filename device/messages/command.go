package messages

type CommandType uint8

const (
	CommandRestart CommandType = iota
	CommandSetInterval
	CommandSetThreshold
	CommandEnable
	CommandDisable
)

type Command struct {
	Type  CommandType
	Value float64
}
