package jobs

import (
	"context"
	"fmt"

	"github.com/fivetentaylor/pointy/pkg/ai"
	"github.com/fivetentaylor/pointy/pkg/background/wire"
	"github.com/fivetentaylor/pointy/pkg/constants"
	"github.com/fivetentaylor/pointy/pkg/env"
)

func SummarizeSessionJob(ctx context.Context, args *wire.SummarizeSession) error {
	log := env.SLog(ctx)

	// Check if there is an open session for the user
	lastMsg, err := env.Redis(ctx).Get(ctx, fmt.Sprintf(constants.DocUserLastMessageKey, args.DocId, args.UserId)).Int64()
	if err != nil {
		log.Error("error getting last message time", "error", err)
	}

	// If there is a LastMessageTime set, and has been a messsage sent after the job's last message time, don't summarize
	if args.LastMessageTime > 0 && lastMsg > args.LastMessageTime {
		log.Info("Skipping summarize session", "lastMessageTime", lastMsg, "jobLastMessageTime", args.LastMessageTime)
		return nil
	}

	dag := ai.SummarizeDag()
	err = dag.Run(ctx, map[string]any{
		"sessionId": args.SessionId,
		"docId":     args.DocId,
		"userId":    args.UserId,
	})
	if err != nil {
		log.Error("Failed to summarize session", "error", err)
	}

	return err
}

func SummarizeCommentThreadJob(ctx context.Context, args *wire.SummarizeCommentThread) error {
	log := env.Log(ctx)

	dag := ai.SummarizeThreadDag()
	err := dag.Run(ctx, map[string]any{
		"docId":         args.DocId,
		"eventId":       args.EventId,
		"threadEventId": args.ThreadEventId,
	})

	if err != nil {
		log.Error("Failed to summarize comment thread", "error", err)
	}

	return err

}
