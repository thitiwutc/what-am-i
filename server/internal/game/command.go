package game

import "github.com/google/uuid"

type CommandType int8

const (
	CmdTypeJoinRoom CommandType = iota + 1
	CmdTypeStartGame
	CmdTypeSubmitSecretWord
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

type StartGameCommand struct {
	Command
	RoomID string `json:"rid" validate:"required,len=4"`
}

type SubmitSecretWordCommand struct {
	Command
	RoomID   string    `json:"rid" validate:"required,len=4"`
	PlayerID uuid.UUID `json:"pid"`
	Word     string    `json:"word" validte:"required"`
}
