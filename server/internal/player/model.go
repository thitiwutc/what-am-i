package player

import "github.com/google/uuid"

type Player struct {
	ID   uuid.UUID
	Name string
}

type PlayerResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
