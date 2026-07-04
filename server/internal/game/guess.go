package game

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrGuessNotFound = errors.New("guess not found")

type Guess struct {
	ID           uuid.UUID `json:"id"`
	Word         string    `json:"word"`
	CorrectCount int       `json:"correct_count"`
	WrongCount   int       `json:"wrong_count"`
	CreatedAt    time.Time `json:"created_at"`
}

type Guesses []*Guess

func (gs Guesses) FindByID(id uuid.UUID) (*Guess, error) {
	for _, g := range gs {
		if g.ID == id {
			return g, nil
		}
	}

	return nil, ErrGuessNotFound
}
