package jobs

import (
	"context"
	"time"

	"github.com/teamreviso/code/pkg/background/wire"
	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/service/notifications"
	"github.com/teamreviso/code/pkg/service/sharing"
)

func AccessDocJob(ctx context.Context, args *wire.AccessDoc) error {
	log := env.Log(ctx)

	userID := args.UserId
	docID := args.DocId
	timestampStr := args.TimestampStr

	timestamp, err := time.Parse(time.RFC3339, timestampStr)
	if err != nil {
		log.Errorf("error parsing timestamp %s: %s", timestampStr, err)
		return err
	}

	q := env.Query(ctx)

	docAccessTbl := q.DocumentAccess

	access, err := docAccessTbl.Where(docAccessTbl.UserID.Eq(userID)).Where(docAccessTbl.DocumentID.Eq(docID)).First()
	if err != nil {
		log.Errorf("error getting access: %s", err)
		return err
	}
	if access == nil {
		log.Errorf("no access found")
		return nil
	}

	if access.LastAccessedAt.IsZero() {
		err = notifications.EnqueueFirstOpen(ctx, docID, userID)
		if err != nil {
			log.Errorf("error notifying first read: %s", err)
			return err
		}

		doc, err := q.Document.Where(q.Document.ID.Eq(docID)).First()
		if err != nil {
			log.Errorf("error getting document: %s", err)
			return err
		}

		err = sharing.AddTimelineJoin(ctx, doc, userID)
		if err != nil {
			log.Errorf("error adding to timeline: %s", err)
		}
	}

	access.LastAccessedAt = timestamp
	err = docAccessTbl.Save(access)
	if err != nil {
		log.Errorf("error saving access: %s", err)
		return err
	}

	return nil
}
