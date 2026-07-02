package game

import "github.com/google/uuid"

type Player struct {
	ID     uuid.UUID
	Name   string
	Notify chan<- *RoomState
}

type PlayerState struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
