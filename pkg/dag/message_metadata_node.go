package dag

import (
	"context"
	"fmt"

	"github.com/fivetentaylor/pointy/pkg/env"
)

type MessageMetadataNode struct {
	BaseLLMNode

	AllowEditsNode Node
	DefaultNode    Node
}

type AllowEditsNodeInputs struct {
	ThreadId       string `key:"threadId"`
	InputMessageId string `key:"inputMessageId"`
}

func (n *MessageMetadataNode) Run(ctx context.Context) (Node, error) {
	log := env.Log(ctx)
	log.Info("👾 running allow edits node")

	input := &AllowEditsNodeInputs{}
	err := n.Hydrate(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error hydrating: %s", err)
	}

	inMsg, err := env.Dynamo(ctx).GetAiThreadMessage(input.ThreadId, input.InputMessageId)
	if err != nil {
		return nil, fmt.Errorf("error getting input message: %s", err)
	}

	SetStateKey(ctx, "llmChoice", inMsg.MessageMetadata.Llm)

	if inMsg.MessageMetadata.AllowDraftEdits {
		return n.AllowEditsNode, nil
	}

	return n.DefaultNode, nil
}
