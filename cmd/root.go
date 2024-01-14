package cmd

import (
	"context"
	"errors"
	"os"
	"syscall"

	"github.com/gofiber/fiber/v3"
	"github.com/oklog/run"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Execute() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if err := execute(); err != nil {
		log.Error().Err(err).Msg("service exitted abnormally")
		os.Exit(1)
	}
}

func execute() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var g run.Group
	{
		app := fiber.New()
		g.Add(func() error {
			return app.Listen(":3000")
		}, func(error) {
			app.Shutdown()
		})
	}

	g.Add(run.SignalHandler(ctx,
		os.Interrupt,
		os.Kill,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGHUP,
	))

	var se run.SignalError
	if err := g.Run(); err != nil && !errors.As(err, &se) {
		return err
	}

	return nil
}
