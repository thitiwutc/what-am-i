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

type Player struct {
	ID   uuid.UUID
	Name string
}

type PlayerResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Room struct {
	ID      string
	State   State
	Players map[uuid.UUID]Player
}

type CreateRoomResponse struct {
	RoomID string `json:"room_id"`
}

type JoinRoomRequest struct {
	PlayerName string `json:"player_name" validate:"max=20"`
}

type JoinRoomResponse struct {
	Players []PlayerResponse `json:"players"`
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

func JoinRoomHandler(lgr *zerolog.Logger, rr *RoomRepository) fiber.Handler {
	return func(c fiber.Ctx) error {
		reqBody := new(JoinRoomRequest)
		err := c.Bind().Body(reqBody)
		if err != nil {
			return fmt.Errorf("bind json request body: %w", err)
		}

		rid := c.Params("room_id")
		lgr.Debug().Msg("Room ID: " + rid)
		if len(rid) != 4 {
			return errors.New("invalid room id")
		}

		room, err := rr.FindByID(rid)
		if err != nil {
			return fmt.Errorf("find room %s: %w", rid, err)
		}

		p := Player{
			ID:   uuid.New(),
			Name: reqBody.PlayerName,
		}
		if p.Name == "" {
			p.Name = util.GeneratePlayerName()
		}

		room.Players[p.ID] = p
		err = rr.Update(room)
		if err != nil {
			return fmt.Errorf("update room %s: %w", rid, err)
		}

		players := make([]PlayerResponse, 0, len(room.Players))
		for _, p := range room.Players {
			players = append(players, PlayerResponse{
				ID:   p.ID,
				Name: p.Name,
			})
		}

		return c.JSON(BaseResponse[JoinRoomResponse]{
			Data: JoinRoomResponse{
				Players: players,
			},
		})
	}
}

var ErrRoomNotFound = errors.New("room not found")

type RoomRepository struct {
	rooms map[string]*Room
	subs  map[string]map[string]<-chan *Room
	lgr   *zerolog.Logger
}

func NewRoomRepository(lgr *zerolog.Logger) *RoomRepository {
	return &RoomRepository{
		rooms: make(map[string]*Room),
		subs:  make(map[string]map[string]<-chan *Room),
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

	return nil
}

func (r *RoomRepository) Subscribe(rid string, pid string, c <-chan *Room) error {
	_, exists := r.subs[rid]
	if !exists {
		return ErrRoomNotFound
	}

	r.subs[rid][pid] = c

	return nil
}

func (r *RoomRepository) Unsubscribe(rid string, pid string, c <-chan *Room) error {
	_, exists := r.subs[rid]
	if !exists {
		return ErrRoomNotFound
	}

	delete(r.subs[rid], pid)

	return nil
}
