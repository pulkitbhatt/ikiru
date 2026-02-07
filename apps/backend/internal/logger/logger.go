package logger

import (
	"io"
	"os"
	"time"

	"github.com/pulkitbhatt/ikiru/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

func New(cfg *config.Config) zerolog.Logger {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = time.RFC3339Nano

	level, err := zerolog.ParseLevel(cfg.Primary.LogLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}

	var output io.Writer = os.Stdout

	if cfg.Primary.Env == config.EnvDev {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	}

	return zerolog.New(output).
		Level(level).
		With().
		Timestamp().
		Logger()
}
