package logging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// Logger is the global application logger.
var Logger zerolog.Logger

// Setup initializes the global logger with the given level string.
// Valid levels: "trace", "debug", "info", "warn", "error", "fatal", "panic".
// Defaults to "info" if the level string is invalid.
func Setup(level string) {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}

	Logger = zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339},
	).Level(lvl).With().Timestamp().Logger()
}

func init() {
	Setup("info")
}
