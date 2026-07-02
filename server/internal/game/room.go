package game

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Room struct {
	ID      string
	State   State
	Players map[uuid.UUID]Player
}

type RoomState struct {
	ID      string        `json:"id"`
	State   State         `json:"state"`
	Players []PlayerState `json:"players"`
}

type CreateRoomResponse struct {
	RoomID string `json:"room_id"`
}

type State int8

const (
	StatePreGame State = iota + 1
)

func CreateRoomHandler(lgr *zerolog.Logger, rr *RoomRepository) fiber.Handler {
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
			State:   StatePreGame,
			Players: make(map[uuid.UUID]Player),
		}
		r.rooms[rid] = &room

		return &room, nil
	}

	return nil, fmt.Errorf("unable to create room with unique id")
}

func (r *RoomRepository) FindByID(id string) (*Room, error) {
	room, exists := r.rooms[id]
	if !exists {
		return nil, ErrRoomNotFound
	}

	return room, nil
}

func (r *RoomRepository) Update(room *Room) error {
	_, exists := r.rooms[room.ID]
	if !exists {
		return ErrRoomNotFound
	}

	r.rooms[room.ID] = room

	ps := make([]PlayerState, 0, len(room.Players))
	for _, p := range room.Players {
		ps = append(ps, PlayerState{
			ID:   p.ID,
			Name: p.Name,
		})
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
