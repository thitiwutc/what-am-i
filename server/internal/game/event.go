package game

type EventType int8

const (
	EvtError EventType = iota
	EvtRoomState
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
