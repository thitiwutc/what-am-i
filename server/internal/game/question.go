package game

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrQuestionNotFound = errors.New("question not found")

type Question struct {
	ID          uuid.UUID `json:"id"`
	Value       string    `json:"value"`
	YesCount    int       `json:"yes_count"`
	NoCount     int       `json:"no_count"`
	UnsureCount int       `json:"unsure_count"`
	CreatedAt   time.Time `json:"created_at"`
}

type Questions []*Question

func (as Questions) FindByID(id uuid.UUID) (*Question, error) {
	for _, a := range as {
		if a.ID == id {
			return a, nil
		}
	}

	return nil, ErrQuestionNotFound
}
