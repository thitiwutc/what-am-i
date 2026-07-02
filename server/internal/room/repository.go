package room

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type RoomRepository struct {
	rooms map[string]*Room
	lgr   *zerolog.Logger
}

func NewRoomRepository(lgr *zerolog.Logger) *RoomRepository {
	return &RoomRepository{
		rooms: make(map[string]*Room),
		lgr:   lgr,
	}
}

func (r *RoomRepository) Create() (*Room, error) {
	// Try to create room with unique ID 3 times.
	for range 3 {
		buf := make([]byte, 2)
		n, err := rand.Read(buf)
		if err != nil {
			return nil, fmt.Errorf("read secure random file: %w", err)
		}
		r.lgr.Printf("Read %d byte(s) from secure random file", n)

		rid := hex.EncodeToString(buf)
		if _, exists := r.rooms[rid]; exists {
			continue
		}
		room := Room{
			ID:      rid,
			Players: map[uuid.UUID]struct{}{},
		}
		r.rooms[rid] = &room

		return &room, nil
	}

	return nil, fmt.Errorf("unable to create room with unique id")
}
