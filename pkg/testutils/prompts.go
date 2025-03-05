package testutils

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"testing"

	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/service/prompts"
)

var promptsRunner = &OnceRunner{
	once: sync.Once{},
	wg:   sync.WaitGroup{},
}

func RefreshPrompts(t *testing.T, ctx context.Context) {
	err := prompts.Refresh(ctx)
	if err != nil {
		t.Fatal(fmt.Errorf("[RefreshPrompts] error refreshing prompts: %s", err))
	}

	return
}

func EnsurePrompts(t *testing.T, ctx context.Context) {
	promptsRunner.mu.Lock()
	promptsRunner.wg.Add(1) // Increment the WaitGroup counter

	insertPrompts := func() {
		rawDb := env.RawDB(ctx)

		// read prompts.sql
		promptSql, err := os.ReadFile("pkg/testutils/prompts.sql")
		if err != nil {
			t.Fatal(fmt.Errorf("[EnsurePrompts] error reading prompts.sql: %s", err))
		}

		err = rawDb.Exec(string(promptSql)).Error
		if err != nil {
			slog.Error("[EnsurePrompts] error inserting prompts", "error", err)
		}

		return
	}

	go func() {
		promptsRunner.once.Do(insertPrompts)
		promptsRunner.wg.Done()
	}()

	promptsRunner.wg.Wait()
	promptsRunner.mu.Unlock()
}
