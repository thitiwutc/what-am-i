package main

import (
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/thitiwutc/what-am-i/server/internal/room"
)

func main() {
	lgr := zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.TimeOnly,
	}).
		With().
		Timestamp().
		Logger()
	app := fiber.New(fiber.Config{
		StructValidator: &structValidator{validate: validator.New()},
	})

	// Middlewares
	app.Use(logger.New())
	app.Use(requestid.New())
	app.Get(healthcheck.LivenessEndpoint, healthcheck.New())
	app.Get(healthcheck.ReadinessEndpoint, healthcheck.New())
	app.Get(healthcheck.StartupEndpoint, healthcheck.New())

	// TODO: Use for non-prod only
	app.Use(cors.New())

	api := app.Group("/api")

	roomRepo := room.NewRoomRepository(&lgr)

	roomGroup := api.Group("/rooms")
	roomGroup.Post("/", room.CreateRoomHandler(&lgr, roomRepo))

	log.Fatal().Msg(app.Listen(":3000").Error())
}

type structValidator struct {
	validate *validator.Validate
}

// Validator needs to implement the Validate method
func (v *structValidator) Validate(out any) error {
	return v.validate.Struct(out)
}
