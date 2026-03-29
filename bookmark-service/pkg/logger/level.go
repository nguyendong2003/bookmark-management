package logger

import (
	"os"

	"github.com/rs/zerolog"
)

const EnvLogLevel = "LOG_LEVEL"

// SetLogLevel sets the global log level based on the LOG_LEVEL environment variable
func SetLogLevel() {
	// Get log level from env variable, default to info level if not set or invalid
	level, err := zerolog.ParseLevel(os.Getenv(EnvLogLevel))
	if err != nil || level == zerolog.NoLevel {
		level = zerolog.InfoLevel
	}

	// Set global log level
	zerolog.SetGlobalLevel(level)
}
