package room

import "github.com/google/uuid"

type Room struct {
	ID      string
	Players map[uuid.UUID]struct{}
}
