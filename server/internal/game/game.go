package game

import (
	"encoding/json"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/contrib/v3/websocket"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/thitiwutc/what-am-i/server/internal/util"
)

func GameHandler(
	validator *validator.Validate,
	lgr *zerolog.Logger,
	rr *RoomRepository,
) func(*websocket.Conn) {
	return func(c *websocket.Conn) {
		lgr.Info().Msg("Upgraded protocol to websocket")

		var notify chan *RoomState
		var roomId string
		var pid uuid.UUID

		for {
			mt, msg, err := c.ReadMessage()
			// Terminate read loop on closed connection.
			if websocket.IsCloseError(
				err,
				websocket.CloseNormalClosure,
				websocket.CloseGoingAway,
				websocket.CloseNoStatusReceived,
				websocket.CloseAbnormalClosure,
			) {
				lgr.Error().Err(err).Msg("Websocket connection closed")
				if roomId != "" {
					lgr := lgr.With().Str("room_id", roomId).Logger()
					room, err := rr.FindByID(roomId)
					if err != nil {
						lgr.Error().Err(err).Msg("Find room failed")
					} else {
						delete(room.Players, pid)
					}

					err = rr.Update(room)
					if err != nil {
						lgr.Error().Err(err).Msg("Update room failed")
					}
					lgr.Info().Str("pid", pid.String()).Msg("Removed player from the room")
				}
				close(notify)
				break
			} else if err != nil {
				lgr.Error().Err(err).Msg("Read message failed")
				c.WriteJSON(ErrorEvent{
					Event: Event{Type: EvtError},
					Error: err.Error(),
				})
				continue
			}

			lgr.Debug().Int("msg_type", mt).Msgf("Received message: %s", msg)

			cmd := new(Command)
			err = json.Unmarshal(msg, cmd)
			if err != nil {
				lgr.Error().Err(err).Msg("JSON Unmarshal websocket message failed")
				c.WriteJSON(ErrorEvent{
					Event: Event{Type: EvtError},
					Error: err.Error(),
				})
				continue
			}

			switch cmd.Type {
			case CmdTypeJoinRoom:
				cmd := new(JoinRoomCommand)
				err = json.Unmarshal(msg, cmd)
				if err != nil {
					lgr.Error().Err(err).Msg("JSON Unmarshal websocket message failed")
					c.WriteJSON(ErrorEvent{
						Event: Event{Type: EvtError},
						Error: err.Error(),
					})
					continue
				}

				err = validator.Struct(cmd)
				if err != nil {
					const errMsg = "Invalid JoinRoom command data"
					lgr.Error().Err(err).Msg(errMsg)
					c.WriteJSON(ErrorEvent{
						Event: Event{Type: EvtError},
						Error: errMsg,
					})
					continue
				}

				lgr.Printf("JoinRoom command received: %+v", cmd)

				room, err := rr.FindByID(cmd.RoomID)
				if err != nil {
					lgr.Error().Err(err).Msg("Find room failed")
					c.WriteJSON(ErrorEvent{
						Event: Event{Type: EvtError},
						Error: err.Error(),
					})
					continue
				}

				notify = make(chan *RoomState, 1)
				p := Player{
					ID:     uuid.New(),
					Name:   cmd.PlayerName,
					Notify: notify,
				}
				pid = p.ID
				roomId = room.ID
				if p.Name == "" {
					p.Name = util.GeneratePlayerName()
				}
				room.Players[p.ID] = p

				go func() {
					lgr := lgr.With().Str("pid", p.ID.String()).Logger()
					// Waiting for notification
					for r := range notify {
						lgr.Print("Room state changed")
						ev := RoomStateEvent{
							Event: Event{
								Type: EvtRoomState,
							},
							RoomState: *r,
						}
						b, err := json.Marshal(ev)
						if err != nil {
							lgr.Error().Err(err).Msg("JSON marshal room state event failed")
							c.WriteJSON(ErrorEvent{
								Event: Event{Type: EvtError},
								Error: err.Error(),
							})
							continue
						}

						c.WriteMessage(websocket.TextMessage, b)
					}
				}()

				err = rr.Update(room)
				if err != nil {
					lgr.Error().Err(err).Msg("Update room failed")
					c.WriteJSON(ErrorEvent{
						Event: Event{Type: EvtError},
						Error: err.Error(),
					})
					continue
				}
			default:
				errMsg := "unknown command type: " + strconv.Itoa(int(cmd.Type))
				lgr.Error().Msg(errMsg)
				c.WriteJSON(ErrorEvent{
					Event: Event{Type: EvtError},
					Error: errMsg,
				})
				continue
			}
		}
	}
}
