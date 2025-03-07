package jobs

import (
	"context"

	"github.com/fivetentaylor/pointy/pkg/ai"
	"github.com/fivetentaylor/pointy/pkg/background/wire"
	"github.com/fivetentaylor/pointy/pkg/env"
)

func RespondToThreadJob(ctx context.Context, args *wire.RespondToThread) error {
	log := env.Log(ctx)

	log.Info("responding to ai thread", "args", args)

	state := map[string]any{
		"docId":    args.DocId,
		"threadId": args.ThreadId,
		"authorId": args.AuthorId,
		"userId":   args.UserId,

		"inputMessageId":  args.InputMessageId,
		"outputMessageId": args.OutputMessageId,
	}

	dag := ai.ThreadDagV2()
	dag.ParentId = args.DocId
	err := dag.Run(ctx, state)
	if err != nil {
		log.Error("Failed to respond to thread", "error", err)
		return err
	}

	return nil
}
