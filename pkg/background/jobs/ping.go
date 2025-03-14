package jobs

import (
	"context"
	"log/slog"

	"github.com/fivetentaylor/pointy/pkg/background/wire"
	"github.com/fivetentaylor/pointy/pkg/env"
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
