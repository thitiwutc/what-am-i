package game

type State int8

const (
	StateLobby State = iota + 1
	StateThinking
	StateAskOrGuess
	StateAnswering
	StateGuessChecking
	StateFinished
)

func (s State) String() string {
	switch s {
	case StateLobby:
		return "Lobby"
	case StateThinking:
		return "Thinking"
	case StateAskOrGuess:
		return "Ask or Guess"
	case StateAnswering:
		return "Answering"
	case StateGuessChecking:
		return "Guess Checking"
	case StateFinished:
		return "Finished"
	default:
		return "unknown state"
	}
}
