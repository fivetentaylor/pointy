package jobs

import (
	"context"

	"github.com/teamreviso/code/pkg/background/wire"
	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/rogue"
)

func SnapshotRogueJob(ctx context.Context, arg *wire.SnapshotRogue) error {
	log := env.Log(ctx)

	log.Info("snapshot started", "arg", arg)

	docId := arg.DocId

	s3 := env.S3(ctx)
	redis := env.Redis(ctx)
	query := env.Query(ctx)

	ds := rogue.NewDocStore(s3, query, redis)
	size, err := ds.SizeDeltaLog(ctx, docId)
	if err != nil {
		log.Errorf("error getting delta log size: %s", err)
		return err
	}

	if size == 0 {
		log.Info("snapshot skipped (no changes)", "args", arg)
		return nil
	}

	err = ds.SnapshotDoc(ctx, docId)
	if err != nil {
		log.Errorf("error snapshotting doc: %s", err.Error())
		return err
	}

	log.Info("snapshot completed", "args", arg)

	return nil
}
