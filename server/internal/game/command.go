package game

type CommandType int8

const (
	CmdTypeJoinRoom CommandType = iota + 1
	CmdTypeStartGame
)

// Command is used for deserializing message via websocket protocol.
type Command struct {
	Type CommandType `json:"type"`
}

type JoinRoomCommand struct {
	Command
	RoomID     string `json:"rid" validate:"required,len=4"`
	PlayerName string `json:"pname" validate:"max=20"`
}
