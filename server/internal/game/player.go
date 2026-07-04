package game

import (
	"errors"

	"github.com/google/uuid"
)

var ErrPlayerNotFound = errors.New("player not found")

type PlayerRole int8

const (
	RoleHost PlayerRole = iota + 1
	RoleParticipant
)

type Player struct {
	ID     uuid.UUID
	Name   string
	Notify chan<- *RoomState
	Role   PlayerRole

	// SecretWord is a player summitted word in thinking phase
	// and guessing word in guessing phase
	SecretWord string

	Questions Questions
	Guesses   Guesses
}

type PlayerState struct {
	ID                    uuid.UUID   `json:"id"`
	Name                  string      `json:"name"`
	Role                  PlayerRole  `json:"role"`
	IsSecretWordSubmitted bool        `json:"isws"`
	SecretWord            string      `json:"secret_word,omitempty"`
	Questions             []*Question `json:"questions"`
	Guesses               Guesses     `json:"guesses"`
	IsWordGuessed         bool        `json:"is_word_guessed"`
}

type Players []*Player

func (ps Players) FindByID(id uuid.UUID) (*Player, error) {
	for _, p := range ps {
		if p.ID == id {
			return p, nil
		}
	}

	return nil, ErrPlayerNotFound
}

func (ps Players) DeleteByID(id uuid.UUID) (Players, error) {
	foundIdx := -1
	for i, p := range ps {
		if p.ID == id {
			foundIdx = i
			break
		}
	}

	if foundIdx < 0 {
		return nil, ErrPlayerNotFound
	}

	ps = append(ps[:foundIdx], ps[foundIdx+1:]...)

	return ps, nil
}
func (ps Players) IsAllWordsGuessed() bool {
	// Check if all words are guessed correctly.
	for _, p := range ps {
		if !p.Guesses.IsCorrect() {
			return false
		}
	}

	return true
}
