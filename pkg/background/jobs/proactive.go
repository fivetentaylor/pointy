package jobs

import (
	"context"

	"github.com/fivetentaylor/pointy/pkg/ai"
	"github.com/fivetentaylor/pointy/pkg/background/wire"
	"github.com/fivetentaylor/pointy/pkg/env"
)

func ProactiveAiMessageJob(ctx context.Context, args *wire.ProactiveAiMessage) error {
	log := env.Log(ctx)

	log.Info("[dag] responding to ai thread", "args", args)

	dag := ai.ProactiveDag()

	err := dag.Run(ctx, map[string]interface{}{
		"type":        args.Type,
		"docId":       args.DocId,
		"containerId": args.ContainerId,
		"messageId":   args.MessageId,
		"threadId":    args.ThreadId,
	})
	if err != nil {
		log.Error("Failed to respond to thread", "error", err)
	}

	return err
}
