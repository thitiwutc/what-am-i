package game

import (
	"time"

	"github.com/google/uuid"
)

type Guess struct {
	ID           uuid.UUID `json:"id"`
	Value        string    `json:"value"`
	CorrectCount int       `json:"correct_count"`
	WrongCount   int       `json:"wrong_count"`
	CreatedAt    time.Time `json:"created_at"`
}
