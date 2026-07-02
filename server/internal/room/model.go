package room

import (
	"github.com/google/uuid"
	"github.com/thitiwutc/what-am-i/server/internal/player"
)

type Room struct {
	ID      string
	Players map[uuid.UUID]player.Player
}

type CreateRoomResponse struct {
	RoomID string `json:"room_id"`
}

type JoinRoomRequest struct {
	PlayerName string `json:"player_name" validate:"max=20"`
}

type JoinRoomResponse struct {
	PlayerID   string `json:"player_id"`
	PlayerName string `json:"player_name"`
}
