package game

import "github.com/google/uuid"

type Word struct {
	PlayerID uuid.UUID
	Value    string
}
