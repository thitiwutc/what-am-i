package game

import "github.com/google/uuid"

type PlayerRole int8

const (
	RoleHost PlayerRole = iota + 1
	RoleParticipant
)

type Player struct {
	ID         uuid.UUID
	Name       string
	Notify     chan<- *RoomState
	Role       PlayerRole
	SecretWord string
}

type PlayerState struct {
	ID                    uuid.UUID  `json:"id"`
	Name                  string     `json:"name"`
	Role                  PlayerRole `json:"role"`
	IsSecretWordSubmitted bool       `json:"isws"`
}
