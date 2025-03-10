package jobs

import (
	"context"

	"github.com/fivetentaylor/pointy/pkg/background/wire"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/rogue"
)

func SnapshotAllRogueJob(ctx context.Context, arg *wire.SnapshotAll) error {
	log := env.Log(ctx)

	log.Info("snapshot all started")

	db := env.RawDB(ctx)
	s3 := env.S3(ctx)
	redis := env.Redis(ctx)
	query := env.Query(ctx)
	ds := rogue.NewDocStore(s3, query, redis)

	page, pageSize := 0, 10
	excludeValues := []interface{}{arg.Version, "error"}
	for {
		var docIds []string
		// TODO: Mark failed jobs documents maybe?
		result := db.Table("documents").
			Select("id").
			Where("rogue_version NOT IN ? OR rogue_version IS NULL", excludeValues).
			Where("deleted_at IS NULL").
			Limit(pageSize).
			Find(&docIds)

		if result.Error != nil {
			return result.Error
		}

		if len(docIds) == 0 {
			break // No more records
		}

		for _, docId := range docIds {
			err := ds.SnapshotDoc(ctx, docId)
			if err != nil {
				log.Errorf("error snapshotting doc: %s", err.Error())

				result := db.Table("documents").Where("id = ?", docId).Update("rogue_version", "error")
				if result.Error != nil {
					log.Errorf("failed to update document: %v", result.Error)
				}

				continue
			}

			result := db.Table("documents").Where("id = ?", docId).Update("rogue_version", arg.Version)
			if result.Error != nil {
				log.Errorf("failed to update document: %v", result.Error)
			}
		}

		page++
	}

	log.Info("snapshot all completed")

	return nil
}
