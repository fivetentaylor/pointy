package jobs

import (
	"context"
	"log/slog"

	"github.com/teamreviso/code/pkg/background/wire"
	"github.com/teamreviso/code/pkg/env"
)

// Ping simpliy pings the database
func PingJob(ctx context.Context, _ *wire.Ping) error {
	slog.Info("ping")
	err := env.RawDB(ctx).Select([]string{"1"}).Error
	if err != nil {
		return err
	}

	slog.Info("pong")

	return nil
}
