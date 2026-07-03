package game

type State int8

const (
	StateLobby State = iota + 1
	StateThinking
	StateGuessing
	StateFinished
)
