package game

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrActionNotFound = errors.New("action not found")

type ActionType int8

const (
	ActionTypeQuestion ActionType = iota + 1
	ActionTypeGuessing
)

type Action struct {
	ID          uuid.UUID  `json:"id"`
	Type        ActionType `json:"type"`
	Value       string     `json:"value"`
	YesCount    int        `json:"yes_count"`
	NoCount     int        `json:"no_count"`
	UnsureCount int        `json:"unsure_count"`
	CreatedAt   time.Time  `json:"created_at"`
}

type Actions []*Action

func (as Actions) FindByID(id uuid.UUID) (*Action, error) {
	for _, a := range as {
		if a.ID == id {
			return a, nil
		}
	}

	return nil, ErrActionNotFound
}
