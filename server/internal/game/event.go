package game

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
	RoomState RoomState `json:"room_state"`
}

type JoinRoomResultEvent struct {
	Event
	Player PlayerState `json:"player"`
}
