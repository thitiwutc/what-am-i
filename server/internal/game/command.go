package game

type CommandType int8

const (
	CmdTypeJoinRoom CommandType = iota + 1
)

// Command is used for deserializing message via websocket protocol.
type Command struct {
	Type CommandType `json:"type"`
}

type JoinRoomCommand struct {
	Command
	RoomID     string `json:"room_id" validate:"required,len=4"`
	PlayerName string `json:"player_name" validate:"max=20"`
}
