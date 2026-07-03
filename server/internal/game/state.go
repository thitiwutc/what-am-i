package game

type State int8

const (
	StateLobby State = iota + 1
	StateThinking
	StateAsking
	StateAnswering
	StateFinished
)

func (s State) String() string {
	switch s {
	case StateLobby:
		return "Lobby"
	case StateThinking:
		return "Thinking"
	case StateAsking:
		return "Asking"
	case StateAnswering:
		return "Answering"
	case StateFinished:
		return "Finished"
	default:
		return "unknown state"
	}
}
