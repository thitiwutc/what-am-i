package room

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
	"github.com/thitiwutc/what-am-i/server/internal/api"
)

func CreateRoomHandler(lgr *zerolog.Logger, rr *RoomRepository) fiber.Handler {
	return func(c fiber.Ctx) error {
		room, err := rr.Create()
		if err != nil {
			return fmt.Errorf("create room: %w", err)
		}

		return c.JSON(api.BaseResponse[string]{
			Data: room.ID,
		})
	}
}
