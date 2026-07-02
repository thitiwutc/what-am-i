package game

type CommandType int8

const (
	Subscribe CommandType = iota + 1
)

// Command is used for deserializing message via websocket protocol.
type Command struct {
	Type CommandType `json:"type"`
}

type SubscribeCommand struct {
	Command
	Payload string `json:"payload"`
}
