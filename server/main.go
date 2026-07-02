package main

import (
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/thitiwutc/what-am-i/server/internal/room"
)

func main() {
	lgr := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).
		With().
		Timestamp().
		Logger()
	app := fiber.New()

	// Middlewares
	app.Use(logger.New())
	app.Use(requestid.New())
	app.Get(healthcheck.LivenessEndpoint, healthcheck.New())
	app.Get(healthcheck.ReadinessEndpoint, healthcheck.New())
	app.Get(healthcheck.StartupEndpoint, healthcheck.New())

	roomRepo := room.NewRoomRepository(&lgr)

	roomGroup := app.Group("/rooms")
	roomGroup.Post("/", room.CreateRoomHandler(&lgr, roomRepo))

	log.Fatal().Msg(app.Listen(":3000").Error())
}
