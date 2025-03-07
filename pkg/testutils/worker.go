package testutils

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/fivetentaylor/pointy/pkg/background/worker"
	"github.com/fivetentaylor/pointy/pkg/config"
	"github.com/fivetentaylor/pointy/pkg/env"
)

func RunWorker(t *testing.T) context.CancelFunc {
	EnsureStorage()
	ctx := TestContext()

	workerCfg := config.Worker{
		Addr:            "127.0.0.1:0",
		Concurrency:     1,
		OpenAIKey:       os.Getenv("OPENAI_API_KEY"),
		WorkerRedisAddr: os.Getenv("REDIS_URL"),
	}

	if workerCfg.OpenAIKey == "" {
		t.Fatal("OPENAI_API_KEY must be set")
	}
	if workerCfg.WorkerRedisAddr == "" {
		t.Fatal("REDIS_URL must be set")
	}

	gormdb := env.RawDB(ctx)
	rc := env.Redis(ctx)

	w, cancelWorker, err := worker.New(workerCfg, gormdb, rc)
	if err != nil {
		t.Fatal(fmt.Errorf("error creating worker: %w", err))
	}

	go func() {
		w.Run()
	}()

	return cancelWorker
}
