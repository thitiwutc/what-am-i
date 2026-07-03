package game

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/thitiwutc/what-am-i/server/internal/util"
)

type Room struct {
	ID      string
	State   State
	Players map[uuid.UUID]Player
}

func (r *Room) addPlayer(name string, notify chan<- *RoomState) *Player {
	p := Player{
		ID:     uuid.New(),
		Name:   name,
		Notify: notify,
		Role:   RoleParticipant,
	}
	if p.Name == "" {
		p.Name = util.GeneratePlayerName()
	}
	if len(r.Players) == 0 {
		p.Role = RoleHost
	}

	r.Players[p.ID] = p

	return &p
}

type RoomState struct {
	ID      string        `json:"id"`
	State   State         `json:"state"`
	Players []PlayerState `json:"players"`
}

type CreateRoomResponse struct {
	RoomID string `json:"room_id"`
}

func CreateRoomHandler(lgr *zerolog.Logger, rr *RoomManager) fiber.Handler {
	return func(c fiber.Ctx) error {
		room, err := rr.Create()
		if err != nil {
			return fmt.Errorf("create room: %w", err)
		}

		return c.JSON(BaseResponse[CreateRoomResponse]{
			Data: CreateRoomResponse{
				RoomID: room.ID,
			},
		})
	}
}

var ErrRoomNotFound = errors.New("room not found")

type RoomManager struct {
	rooms map[string]*Room
	lgr   *zerolog.Logger
}

func NewRoomRepository(lgr *zerolog.Logger) *RoomManager {
	return &RoomManager{
		rooms: make(map[string]*Room),
		lgr:   lgr,
	}
}

func (r *RoomManager) Create() (*Room, error) {
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
			State:   StateLobby,
			Players: make(map[uuid.UUID]Player),
		}
		r.rooms[rid] = &room

		return &room, nil
	}

	return nil, fmt.Errorf("unable to create room with unique id")
}

func (r *RoomManager) FindByID(id string) (*Room, error) {
	room, exists := r.rooms[id]
	if !exists {
		return nil, ErrRoomNotFound
	}

	return room, nil
}

func (r *RoomManager) Update(room *Room) error {
	_, exists := r.rooms[room.ID]
	if !exists {
		return ErrRoomNotFound
	}

	r.rooms[room.ID] = room

	ps := make([]PlayerState, 0, len(room.Players))
	for _, p := range room.Players {
		ps = append(ps, PlayerState{
			ID:                    p.ID,
			Name:                  p.Name,
			Role:                  p.Role,
			IsSecretWordSubmitted: p.SecretWord != "",
		})
	}

	// Change state when all player submit their secret words.
	allSubmitted := true
	for _, p := range room.Players {
		if p.SecretWord == "" {
			allSubmitted = false
			break
		}
	}

	r.lgr.Printf("All submitted: %t", allSubmitted)
	if allSubmitted {
		room.State = StateGuessing
	}

	rs := RoomState{
		ID:      room.ID,
		State:   room.State,
		Players: ps,
	}

	// Notify each player in the room.
	for _, p := range r.rooms[room.ID].Players {
		p.Notify <- &rs
	}

	return nil
}
