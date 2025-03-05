package ai

import (
	"context"

	"github.com/teamreviso/code/pkg/dag"
	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/models"
	"github.com/teamreviso/code/pkg/service/messaging"
	"github.com/teamreviso/code/pkg/stackerr"
	"github.com/teamreviso/code/pkg/storage/dynamo"
)

// BaseLLMNode: dag.BaseLLMNode{
// 	Model:    "llama-3.2-3b-preview",
// 	Provider: dag.Groq,
// },

// BaseLLMNode: dag.BaseLLMNode{
// 	Model:    "llama-3.1-8b-instant",
// 	Provider: dag.Groq,
// },

func SummarizeDag() *dag.Dag {
	return dag.New("summarize", &dag.SummarizeUpdateNode{})
}

func SummarizeThreadDag() *dag.Dag {
	return dag.New("summarize-thread", &dag.SummarizeCommentThreadNode{})
}

func ProactiveDag() *dag.Dag {
	d := dag.New("proactive", &dag.ProactiveNode{})
	d.OnError = failure

	return d
}

func ThreadDagV2() *dag.Dag {
	postNode := &dag.TitleThreadNode{}
	askBranch := &dag.AskNode{
		BaseLLMNode: dag.BaseLLMNode{
			AllowLLMOverride: true,
		},
		Next: postNode,
	}

	d := dag.New("threadv2", &dag.MessageMetadataNode{
		AllowEditsNode: &dag.SelectEditTargetNode{
			Next: &dag.TextSegmentationCheckNode{
				LargeText: &dag.SerialReviseNode{
					Next: postNode,
				},
				Default: &dag.ReviseNode{
					Next: postNode,
				},
			},
			NoEdits: askBranch,
		},

		DefaultNode: askBranch,
	})

	d.OnError = failure
	d.OnComplete = complete

	return d
}

func complete(ctx context.Context, d *dag.Dag) {
	log := env.Log(ctx)
	log.Info("dag complete", "name", d.Name)

	threadID, err := dag.GetStateKey[string](ctx, "threadId")
	if err != nil {
		log.Error("error getting container id", "error", stackerr.Wrap(err))
		return
	}

	msgId, err := dag.GetStateKey[string](ctx, "outputMessageId")
	if err != nil {
		log.Error("error getting message id", "error", stackerr.Wrap(err))
		return
	}

	msg, ierr := env.Dynamo(ctx).GetAiThreadMessage(threadID, msgId)
	if ierr != nil {
		log.Error("error getting message", "error", err)
		return
	}

	msg.LifecycleStage = dynamo.MessageLifecycleStageCompleted
	err = messaging.UpdateMessage(ctx, msg)
	if err != nil {
		log.Error("error updating message", "error", err)
		return
	}
}

func failure(ctx context.Context, node dag.Node, err error) {
	if err == nil {
		return
	}
	log := env.Log(ctx)
	log.Error("error running dag", "error", stackerr.Wrap(err))

	threadID, ierr := dag.GetStateKey[string](ctx, "threadId")
	if ierr != nil {
		log.Error("error getting container id", "error", err)
		return
	}

	msgId, ierr := dag.GetStateKey[string](ctx, "outputMessageId")
	if ierr != nil {
		log.Error("error getting message id", "error", err)
		return
	}

	msg, ierr := env.Dynamo(ctx).GetAiThreadMessage(threadID, msgId)
	if ierr != nil {
		log.Error("error getting message", "error", err)
		return
	}

	// If no updates were applied, add error to message and return
	errorAttachment := &models.Error{
		Title: "Error",
		Text:  "Sorry our system was unable to respond to your message. Please try again.",
		Error: err.Error(),
	}
	msg.Attachments.Attachments = append(msg.Attachments.Attachments, &models.Attachment{
		Value: &models.Attachment_Error{
			Error: errorAttachment,
		},
	})
	msg.LifecycleStage = dynamo.MessageLifecycleStageCompleted
	err = messaging.UpdateMessage(ctx, msg)
	if err != nil {
		log.Error("error updating message", "error", err)
		return
	}
}
