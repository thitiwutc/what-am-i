package room

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/thitiwutc/what-am-i/server/internal/api"
	"github.com/thitiwutc/what-am-i/server/internal/player"
	"github.com/thitiwutc/what-am-i/server/internal/util"
)

func CreateRoomHandler(lgr *zerolog.Logger, rr *RoomRepository) fiber.Handler {
	return func(c fiber.Ctx) error {
		room, err := rr.Create()
		if err != nil {
			return fmt.Errorf("create room: %w", err)
		}

		return c.JSON(api.BaseResponse[CreateRoomResponse]{
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

		p := player.Player{
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

		players := make([]player.PlayerResponse, 0, len(room.Players))
		for _, p := range room.Players {
			players = append(players, player.PlayerResponse{
				ID:   p.ID,
				Name: p.Name,
			})
		}

		return c.JSON(api.BaseResponse[JoinRoomResponse]{
			Data: JoinRoomResponse{
				Players: players,
			},
		})
	}
}
