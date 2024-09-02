package logger

import (
	"io"
	"time"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

func Init(level zerolog.Level, w io.Writer) {
	zerolog.SetGlobalLevel(level)
	zerolog.TimeFieldFormat = time.RFC3339
	Log = zerolog.New(w).With().Timestamp().Logger()
	Log.Debug().Msgf("Logger initialized with level %d", level)
}
