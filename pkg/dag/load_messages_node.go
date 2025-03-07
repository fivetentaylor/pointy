package dag

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/fivetentaylor/pointy/pkg/env"
)

type LoadMessagesNode struct {
	Next Node
	Base
}

type LoadMessagesNodeInput struct {
	ThreadId        string `key:"threadId"`
	InputMessageId  string `key:"inputMessageId"`
	OutputMessageId string `key:"outputMessageId"`
}

func (n *LoadMessagesNode) Run(ctx context.Context) (Node, error) {
	input := &LoadMessagesNodeInput{}
	err := n.Hydrate(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error hydrating: %s", err)
	}

	dydb := env.Dynamo(ctx)

	inputMessage, err := dydb.GetAiThreadMessage(input.ThreadId, input.InputMessageId)
	if err != nil {
		log.Error("[LoadMessagesNode] error getting input message", "error", err)
		return nil, fmt.Errorf("error getting input message: %s", err)
	}
	SetStateKey(ctx, "inputMessage", inputMessage)

	outputMessage, err := dydb.GetAiThreadMessage(input.ThreadId, input.OutputMessageId)
	if err != nil {
		log.Error("[LoadMessagesNode] error getting output message", "error", err)
		return nil, fmt.Errorf("error getting output message: %s", err)
	}
	SetStateKey(ctx, "outputMessage", outputMessage)

	return n.Next, nil
}
