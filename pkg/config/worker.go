package config

import (
	"log/slog"

	"github.com/teamreviso/code/pkg/constants"
)

type Worker struct {
	Addr            string
	Logger          *slog.Logger
	Concurrency     int
	OpenAIKey       string
	WorkerRedisAddr string
}

// implement github.com/jpoz/conveyor/config.WorkerConfig
func (w Worker) GetLogger() *slog.Logger {
	return w.Logger
}

func (w *Worker) SetLogger(logger *slog.Logger) {
	w.Logger = logger
}

func (w Worker) GetConcurrency() int {
	return w.Concurrency
}

func (w Worker) GetRedisURL() string {
	return w.WorkerRedisAddr
}

func (w Worker) GetNamespace() string {
	return constants.DefaultJobNamespace
}
