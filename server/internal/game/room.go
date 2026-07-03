package game

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	mrand "math/rand/v2"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/thitiwutc/what-am-i/server/internal/util"
)

type Room struct {
	ID              string
	State           State
	Players         Players
	ActivePlayerIdx int
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

	r.Players = append(r.Players, &p)

	return &p
}

type RoomState struct {
	ID             string                    `json:"id"`
	State          State                     `json:"state"`
	Players        map[uuid.UUID]PlayerState `json:"players"`
	ActivePlayerID uuid.UUID                 `json:"active_player_id"`
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
			ID:              rid,
			State:           StateLobby,
			Players:         nil,
			ActivePlayerIdx: 0,
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

	allSubmitted := true
	for _, p := range room.Players {
		if p.SecretWord == "" {
			allSubmitted = false
			break
		}
	}

	r.lgr.Printf("All submitted: %t", allSubmitted)
	// Start guessing state
	if allSubmitted {
		room.State = StateGuessing

		shufWords := make([]Word, 0)
		for _, p := range room.Players {
			shufWords = append(shufWords, Word{
				PlayerID: p.ID,
				Value:    p.SecretWord,
			})
		}
		oriWords := make([]Word, len(shufWords))
		copy(oriWords, shufWords)

		shufOk := false
		// Shuffle so that no word belong to the player who submitted it.
	SHUFFLE:
		for range 100 {
			mrand.Shuffle(len(shufWords), func(i, j int) {
				shufWords[i], shufWords[j] = shufWords[j], shufWords[i]
			})

			// Check if the word owner get their own word.
			for i, w := range shufWords {
				if w.PlayerID == oriWords[i].PlayerID {
					continue SHUFFLE
				}
			}
			shufOk = true
			break
		}

		if !shufOk {
			return errors.New("shuffle words failed")
		}

		// Actual word shuffling
		for i, w := range shufWords {
			oriWords[i].Value = w.Value
		}

		log.Printf("Pre shuffle: %+v", room.Players)
		for i, w := range oriWords {
			room.Players[i].SecretWord = w.Value
		}
		log.Printf("Post shuffle: %+v", room.Players)

		// Shuffle players for random turn order.
		mrand.Shuffle(len(room.Players), func(i, j int) {
			room.Players[i], room.Players[j] = room.Players[j], room.Players[i]
		})
		r.lgr.Printf("Shuffled players: %+v", room.Players)
	}

	ps := make(map[uuid.UUID]PlayerState)
	for _, p := range room.Players {
		ps[p.ID] = PlayerState{
			ID:                    p.ID,
			Name:                  p.Name,
			Role:                  p.Role,
			IsSecretWordSubmitted: p.SecretWord != "",
			SecretWord:            p.SecretWord,
		}
	}

	rs := RoomState{
		ID:             room.ID,
		State:          room.State,
		Players:        ps,
		ActivePlayerID: room.Players[room.ActivePlayerIdx].ID,
	}

	// Notify each player in the room.
	for _, p := range r.rooms[room.ID].Players {
		p.Notify <- &rs
	}

	return nil
}
