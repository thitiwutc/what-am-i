package game

import "github.com/google/uuid"

type CommandType int8

const (
	CmdTypeJoinRoom CommandType = iota + 1
	CmdTypeStartGame
	CmdTypeSubmitSecretWord
	CmdTypeAsk
	CmdTypeAnswerQuestion
	CmdTypeGuess
	CmdTypeJudgeGuess
)

// Command is sent by client.
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
	Word     string    `json:"word" validate:"required"`
}

type AskCommand struct {
	Command
	RoomID   string    `json:"rid" validate:"required,len=4"`
	PlayerID uuid.UUID `json:"pid"`
	Question string    `json:"question" validate:"required"`
}

type AnswerCommand struct {
	Command
	QuestionID uuid.UUID `json:"qid"`
	RoomID     string    `json:"rid" validate:"required,len=4"`
	PlayerID   uuid.UUID `json:"pid"`
	Answer     Answer    `json:"answer"`
}

type GuessCommand struct {
	Command
	RoomID   string    `json:"rid" validate:"required,len=4"`
	PlayerID uuid.UUID `json:"pid"`
	Guess    string    `json:"guess"`
}

type JudgeGuessCommand struct {
	Command
	GuessID  uuid.UUID `json:"gid"`
	RoomID   string    `json:"rid" validate:"required,len=4"`
	PlayerID uuid.UUID `json:"pid"`
	Decision Decision  `json:"decision"`
}
