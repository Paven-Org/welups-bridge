package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

var logger *zerolog.Logger

func Get() *zerolog.Logger {
	if logger == nil {
		Init(false)
	}
	return logger
}

func Init(structured bool) {

	if logger != nil {
		logger.Info().Str("in", "InitLog").Msg("Logger already initialized")
		return
	}

	var output io.Writer
	if structured {
		output = os.Stdout
	} else {
		output = zerolog.NewConsoleWriter(
			func(w *zerolog.ConsoleWriter) {
				w.Out = os.Stdout
				w.TimeFormat = time.RFC3339
			})
	}

	logger = new(zerolog.Logger)
	*logger = zerolog.
		New(output).
		With().
		Timestamp().
		Caller().
		Logger()

	logger.Info().Msg("logger initialized")
}
