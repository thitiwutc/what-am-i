package game

import (
	"strconv"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/rs/zerolog"
)

func GameHandler(lgr *zerolog.Logger) func(*websocket.Conn) {
	return func(c *websocket.Conn) {
		lgr.Info().Msg("Upgraded protocol to websocket")

		for {
			cmd := new(Command)
			err := c.ReadJSON(&cmd)
			if err != nil {
				lgr.Error().Err(err).Msg("Parse websocket message as JSON failed")
				c.WriteMessage(websocket.TextMessage, []byte(err.Error()))
				continue
			}

			lgr.Printf("Received command: %+v", cmd)

			switch cmd.Type {
			case Subscribe:
				lgr.Printf("Subscribe command received: %+v", cmd)
				c.WriteMessage(websocket.TextMessage, []byte("OK"))
			default:
				errMsg := "unknown command type: " + strconv.Itoa(int(cmd.Type))
				lgr.Error().Msg(errMsg)
				c.WriteMessage(websocket.TextMessage, []byte(errMsg))
				continue
			}
		}
	}
}
