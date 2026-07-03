package game

import "github.com/google/uuid"

type EventType int8

const (
	EvtError EventType = iota
	EvtRoomState
	EvtJoinRoomResult
)

// Event is pushed to client via websocket connection.
type Event struct {
	Type EventType `json:"type"`
}

type ErrorEvent struct {
	Event
	Error string `json:"error"`
}

type RoomStateEvent struct {
	Event
	ID             string        `json:"id"`
	State          State         `json:"state"`
	Players        []PlayerState `json:"players"`
	ActivePlayerID uuid.UUID     `json:"active_player_id"`
}

type JoinRoomResultEvent struct {
	Event
	Player PlayerState `json:"player"`
}
