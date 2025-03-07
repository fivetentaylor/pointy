package jobs

import (
	"context"
	"fmt"

	"github.com/fivetentaylor/pointy/pkg/background/wire"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/service/notifications"
)

func NotifyFirstOpenJob(ctx context.Context, args *wire.NotifyFirstOpen) error {
	log := env.Log(ctx)
	log.Info("notifying user of first open", "docID", args.DocId, "readerID", args.ReaderId)

	err := notifications.SendFirstOpen(ctx, args.DocId, args.ReaderId)
	if err != nil {
		log.Errorf("error sending first open email(s): %s", err)
		return err
	}

	return nil
}

func NotifyNewTimelineCommentJob(ctx context.Context, args *wire.NotifyNewTimelineComment) error {
	log := env.Log(ctx)
	log.Info("notifying user of new timeline comment", "docID", args.DocId, "eventID", args.EventId, "excludeUserIds", args.ExcludeUserIds)

	err := notifications.SendNewComment(ctx, args.DocId, args.EventId, args.ExcludeUserIds)
	if err != nil {
		log.Errorf("error sending new comment email(s): %s", err)
		return err
	}

	return nil
}

func NotifyNewMentionShareJob(ctx context.Context, args *wire.NotifyNewMentionShare) error {
	log := env.Log(ctx)
	log.Info("notifying user of new mention share", "docID", args.DocId, "eventID", args.EventId, "recipientID", args.RecipientId)

	if args.DocId == "" || args.EventId == "" || args.RecipientId == "" {
		return fmt.Errorf("invalid arguments: DocId or EventId or RecipientId is empty")
	}

	err := notifications.SendNewMentionShare(ctx, args.DocId, args.EventId, args.RecipientId)
	if err != nil {
		log.Error("error sending new mention share notification", "error", err)
		return err
	}

	return nil
}
