package game

import (
	"encoding/json"
	"math"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/contrib/v3/websocket"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func GameHandler(
	validator *validator.Validate,
	lgr *zerolog.Logger,
	roomManager *RoomManager,
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
					room, err := roomManager.FindByID(roomId)
					if err != nil {
						lgr.Error().Err(err).Msg("Find room failed")
					} else {
						newPlayers, err := room.Players.DeleteByID(pid)
						if err != nil {
							lgr.Error().Err(err).Msg("Delete player failed")
						} else {
							room.Players = newPlayers
						}

						if len(room.Players) == 0 {
							roomManager.DeleteByID(roomId)
						}
					}

					err = roomManager.Update(room)
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

				room, err := roomManager.FindByID(cmd.RoomID)
				if err != nil {
					lgr.Error().Err(err).Msg("Find room failed")
					c.WriteJSON(ErrorEvent{
						Event: Event{Type: EvtError},
						Error: err.Error(),
					})
					continue
				}

				notify = make(chan *RoomState, 1)
				p := room.addPlayer(cmd.PlayerName, notify)
				roomId = room.ID
				pid = p.ID
				lgr := lgr.With().Str("pid", p.ID.String()).Logger()

				go func() {
					// Waiting for notification
					for r := range notify {
						lgr.Print("Room state changed")

						players := make([]PlayerState, 0, len(r.Players))
						for _, p := range r.Players {
							// Hide the secret word belong to the player themselves.
							if p.ID == pid && !p.Guesses.IsCorrect() {
								p.SecretWord = ""
							}
							players = append(players, p)
						}
						ev := RoomStateEvent{
							Event: Event{
								Type: EvtRoomState,
							},
							ID:           r.ID,
							State:        r.State,
							Players:      players,
							ActivePlayer: r.ActivePlayer,
							LastQuestion: r.LastQuestion,
							LastGuess:    r.LastGuess,
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

				err = roomManager.Update(room)
				if err != nil {
					lgr.Error().Err(err).Msg("Update room failed")
					c.WriteJSON(ErrorEvent{
						Event: Event{Type: EvtError},
						Error: err.Error(),
					})
					continue
				}

				c.WriteJSON(JoinRoomResultEvent{
					Event: Event{Type: EvtJoinRoomResult},
					Player: PlayerState{
						ID:   p.ID,
						Name: p.Name,
						Role: p.Role,
					},
				})
				lgr.Info().Msg("Sent join room result event")
			case CmdTypeStartGame:
				cmd := new(StartGameCommand)
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
					const errMsg = "Invalid command data"
					lgr.Error().Err(err).Msg(errMsg)
					c.WriteJSON(ErrorEvent{
						Event: Event{Type: EvtError},
						Error: errMsg,
					})
					continue
				}

				room, err := roomManager.FindByID(cmd.RoomID)
				if err != nil {
					lgr.Error().Err(err).Msg("Find room failed")
					c.WriteJSON(ErrorEvent{
						Event: Event{Type: EvtError},
						Error: err.Error(),
					})
					continue
				}

				if len(room.Players) < 2 {
					lgr.Warn().Msg("Host tries to start game with the number of players fewer than 2")
					c.WriteJSON(ErrorEvent{
						Event: Event{Type: EvtError},
						Error: "The number of players is less than 2",
					})
					continue
				}

				room.State = StateThinking
				err = roomManager.Update(room)
				if err != nil {
					lgr.Error().Err(err).Msg("Update room failed")
					c.WriteJSON(ErrorEvent{
						Event: Event{Type: EvtError},
						Error: err.Error(),
					})
					continue
				}

				lgr.Info().Str("rid", room.ID).Msg("Started game")
			case CmdTypeSubmitSecretWord:
				cmd := new(SubmitSecretWordCommand)
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
					const errMsg = "Invalid command data"
					lgr.Error().Err(err).Msg(errMsg)
					c.WriteJSON(ErrorEvent{
						Event: Event{Type: EvtError},
						Error: errMsg,
					})
					continue
				}

				lgr := lgr.With().
					Str("word", cmd.Word).
					Str("rid", cmd.RoomID).
					Str("pid", cmd.PlayerID.String()).
					Logger()

				lgr.Info().
					Msg("Received submit secret word command")

				room, err := roomManager.FindByID(cmd.RoomID)
				if err != nil {
					lgr.Error().Err(err).Msg("Find room failed")
					c.WriteJSON(ErrorEvent{
						Event: Event{Type: EvtError},
						Error: err.Error(),
					})
					continue
				}

				p, err := room.Players.FindByID(cmd.PlayerID)
				if err != nil {
					lgr.Error().Err(err).Msg("Player not in the room")
					c.WriteJSON(ErrorEvent{
						Event: Event{Type: EvtError},
						Error: "Player not in the room!",
					})
					continue
				}

				p.SecretWord = cmd.Word

				err = roomManager.Update(room)
				if err != nil {
					lgr.Error().Err(err).Msg("Update room failed")
					c.WriteJSON(ErrorEvent{
						Event: Event{Type: EvtError},
						Error: err.Error(),
					})
					continue
				}

			case CmdTypeAsk:
				cmd := new(AskCommand)
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
					const errMsg = "Invalid command data"
					lgr.Error().Err(err).Msg(errMsg)
					c.WriteJSON(ErrorEvent{
						Event: Event{Type: EvtError},
						Error: errMsg,
					})
					continue
				}

				room, err := roomManager.FindByID(cmd.RoomID)
				if err != nil {
					lgr.Error().Err(err).Msg("Find room failed")
					c.WriteJSON(ErrorEvent{
						Event: Event{Type: EvtError},
						Error: err.Error(),
					})
					continue
				}

				p, err := room.Players.FindByID(cmd.PlayerID)
				if err != nil {
					lgr.Error().Err(err).Msg("Player not in the room")
					c.WriteJSON(ErrorEvent{
						Event: Event{Type: EvtError},
						Error: "Player not in the room!",
					})
					continue
				}

				p.Questions = append(p.Questions, &Question{
					ID:          uuid.New(),
					Value:       cmd.Question,
					YesCount:    0,
					NoCount:     0,
					UnsureCount: 0,
					CreatedAt:   time.Now(),
				})
				room.State = StateAnswering
				roomManager.Update(room)
			case CmdTypeAnswerQuestion:
				cmd := new(AnswerCommand)
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
					const errMsg = "Invalid command data"
					lgr.Error().Err(err).Msg(errMsg)
					c.WriteJSON(ErrorEvent{
						Event: Event{Type: EvtError},
						Error: errMsg,
					})
					continue
				}

				room, err := roomManager.FindByID(cmd.RoomID)
				if err != nil {
					lgr.Error().Err(err).Msg("Find room failed")
					c.WriteJSON(ErrorEvent{
						Event: Event{Type: EvtError},
						Error: err.Error(),
					})
					continue
				}

				p, err := room.Players.FindByID(cmd.PlayerID)
				if err != nil {
					lgr.Error().Err(err).Msg("Player not in the room")
					c.WriteJSON(ErrorEvent{
						Event: Event{Type: EvtError},
						Error: "Player not in the room!",
					})
					continue
				}

				a, err := p.Questions.FindByID(cmd.QuestionID)
				if err != nil {
					lgr.Error().Err(err).Msg("Find question failed")
					c.WriteJSON(ErrorEvent{
						Event: Event{Type: EvtError},
						Error: err.Error(),
					})
					continue
				}

				switch cmd.Answer {
				case AnswerYes:
					a.YesCount++
				case AnswerNo:
					a.NoCount++
				case AnswerUnsure:
					a.UnsureCount++
				default:
					errMsg := "Unknown answer: " + strconv.Itoa(int(cmd.Answer))
					lgr.Info().Msg(errMsg)
					c.WriteJSON(ErrorEvent{
						Event: Event{Type: EvtError},
						Error: errMsg,
					})
					continue
				}

				if a.YesCount+a.NoCount+a.UnsureCount == len(room.Players)-1 {
					room.State = StateAskOrGuess
					room.NextTurn()
				}
				roomManager.Update(room)
			case CmdTypeGuess:
				cmd := new(GuessCommand)
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
					const errMsg = "Invalid command data"
					lgr.Error().Err(err).Msg(errMsg)
					c.WriteJSON(ErrorEvent{
						Event: Event{Type: EvtError},
						Error: errMsg,
					})
					continue
				}

				room, err := roomManager.FindByID(cmd.RoomID)
				if err != nil {
					lgr.Error().Err(err).Msg("Find room failed")
					c.WriteJSON(ErrorEvent{
						Event: Event{Type: EvtError},
						Error: err.Error(),
					})
					continue
				}

				p, err := room.Players.FindByID(cmd.PlayerID)
				if err != nil {
					lgr.Error().Err(err).Msg("Player not in the room")
					c.WriteJSON(ErrorEvent{
						Event: Event{Type: EvtError},
						Error: "Player not in the room!",
					})
					continue
				}

				p.Guesses = append(p.Guesses, &Guess{
					ID:           uuid.New(),
					Word:         cmd.Guess,
					CorrectCount: 0,
					WrongCount:   0,
					CreatedAt:    time.Now(),
				})
				room.State = StateGuessChecking
				roomManager.Update(room)
			case CmdTypeJudgeGuess:
				cmd := new(JudgeGuessCommand)
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
					const errMsg = "Invalid command data"
					lgr.Error().Err(err).Msg(errMsg)
					c.WriteJSON(ErrorEvent{
						Event: Event{Type: EvtError},
						Error: errMsg,
					})
					continue
				}

				room, err := roomManager.FindByID(cmd.RoomID)
				if err != nil {
					lgr.Error().Err(err).Msg("Find room failed")
					c.WriteJSON(ErrorEvent{
						Event: Event{Type: EvtError},
						Error: err.Error(),
					})
					continue
				}

				p, err := room.Players.FindByID(cmd.PlayerID)
				if err != nil {
					lgr.Error().Err(err).Msg("Player not in the room")
					c.WriteJSON(ErrorEvent{
						Event: Event{Type: EvtError},
						Error: "Player not in the room!",
					})
					continue
				}

				guess, err := p.Guesses.FindByID(cmd.GuessID)
				if err != nil {
					lgr.Error().Err(err).Msg("Find guess failed")
					c.WriteJSON(ErrorEvent{
						Event: Event{Type: EvtError},
						Error: err.Error(),
					})
					continue
				}

				switch cmd.Decision {
				case DecisionCorrect:
					guess.CorrectCount++
				case DesicionWrong:
					guess.WrongCount++
				default:
					errMsg := "Unknown guess decision: " + strconv.Itoa(int(cmd.Decision))
					lgr.Info().Msg(errMsg)
					c.WriteJSON(ErrorEvent{
						Event: Event{Type: EvtError},
						Error: errMsg,
					})
					continue
				}

				if guess.CorrectCount+guess.WrongCount == len(room.Players)-1 {
					majorPlayerCount := int(math.Round(float64(guess.CorrectCount+guess.WrongCount) / 2.0))
					isMajorCorrect := guess.CorrectCount >= majorPlayerCount
					lgr.Debug().
						Int("correct_count", guess.CorrectCount).
						Int("wrong_count", guess.WrongCount).
						Int("majority_player_count", majorPlayerCount).
						Msgf("Is correct for majority: %t", isMajorCorrect)
					guess.Correct = isMajorCorrect
					// Check if all players guessed correctly.
					allWordsGuessed := room.Players.IsAllWordsGuessed()
					if allWordsGuessed {
						room.State = StateFinished
					} else {
						room.State = StateAskOrGuess
						room.NextTurn()
					}
				}
				roomManager.Update(room)
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
