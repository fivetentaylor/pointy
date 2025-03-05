package utils

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/teamreviso/code/pkg/prettylog"
)

func NewSlogFromEnv() *slog.Logger {
	var slevel slog.Level
	err := json.Unmarshal([]byte(os.Getenv("LOG_LEVEL")), slevel)
	if err == nil {
		slevel = slog.LevelInfo
	}
	opts := &slog.HandlerOptions{
		Level: slevel,
	}

	env := os.Getenv("ENV")
	if env == "development" {
		return slog.New(
			prettylog.NewHandler(opts),
		)
	}

	return slog.New(
		slog.NewJSONHandler(os.Stdout, opts),
	)
}

func SetDefaultLogger() {
	slog.SetDefault(NewSlogFromEnv())
}
