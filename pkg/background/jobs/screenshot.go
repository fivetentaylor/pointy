package jobs

import (
	"context"

	"github.com/fivetentaylor/pointy/pkg/background/wire"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/rogue"
)

func ScreenshotAllJob(ctx context.Context, args *wire.ScreenshotAll) error {
	log := env.Log(ctx)
	docTlb := env.Query(ctx).Document

	log.Info("screenshot all started")

	allDocs, err := docTlb.Find()
	if err != nil {
		log.Error(err)
		return err
	}

	for _, doc := range allDocs {
		_, err = env.Background(ctx).Enqueue(ctx, &wire.Screenshot{DocId: doc.ID})
		if err != nil {
			log.Error(err)
			return err
		}
	}

	return nil
}

func ScreenshotJob(ctx context.Context, args *wire.Screenshot) error {
	log := env.Log(ctx)
	log.Info("screenshot started", "args", args)

	docId := args.DocId

	s3 := env.S3(ctx)
	redis := env.Redis(ctx)
	query := env.Query(ctx)

	ds := rogue.NewDocStore(s3, query, redis)
	err := ds.ScreenshotDoc(ctx, docId)
	if err != nil {
		log.Errorf("error screenshotting doc: %s", err)
		return err
	}

	log.Info("screenshot completed", "args", args)

	return nil
}
