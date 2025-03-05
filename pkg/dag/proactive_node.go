package dag

import (
	"bytes"
	"context"
	"fmt"

	"github.com/teamreviso/code/pkg/constants"
	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/service/messaging"
	"github.com/teamreviso/code/pkg/storage/dynamo"
	"github.com/teamreviso/code/pkg/utils"
	"github.com/tmc/langchaingo/llms"
)

type ProactiveNode struct {
	Next Node

	BaseLLMNode
}

type proactiveMessage struct {
	Role    llms.ChatMessageType `json:"-"`
	Message string               `json:"message,omitempty"`
}

func NewProactiveNode(next Node) *ProactiveNode {
	return &ProactiveNode{
		Next: next,
	}
}

type ProactiveNodeInput struct {
	DocId       string `key:"docId"`
	ContainerId string `key:"containerId"`
	MessageId   string `key:"messageId"`
	ThreadId    string `key:"threadId"`
}

func (n *ProactiveNode) Run(ctx context.Context) (Node, error) {
	log := env.Log(ctx)

	input := &ProactiveNodeInput{}
	err := n.Hydrate(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error hydrating: %s", err)
	}

	msg, err := env.Dynamo(ctx).GetMessage(input.ContainerId, input.MessageId)
	if err != nil {
		log.Error("error getting message", "error", err)
		return nil, fmt.Errorf("error getting message: %s", err)
	}

	err = n.generate(ctx, input, msg)
	if err != nil {
		log.Error("error generating", "error", err)
		return nil, fmt.Errorf("error generating: %s", err)
	}

	return n.Next, nil
}

func (n *ProactiveNode) generate(
	ctx context.Context,
	input *ProactiveNodeInput,
	msg *dynamo.Message,
) error {
	log := env.Log(ctx)
	log.Info("[dag] generating proactive message", "threadId", input.ThreadId, "docId", input.DocId, "messageId", msg.MessageID)

	_, err := n.GenerateStoredPrompt(ctx, input.ThreadId, constants.PromptNewDoc, map[string]string{},
		llms.WithStreamingFunc(
			n.receiveStreamFunc(ctx, msg),
		),
	)
	if err != nil {
		log.Error("error generating content for freeplay prompt", "error", err)
		return fmt.Errorf("error generating content for freeplay prompt: %s", err)
	}

	msg.LifecycleStage = dynamo.MessageLifecycleStageCompleted
	err = messaging.UpdateMessage(ctx, msg)
	if err != nil {
		return fmt.Errorf("error updating message: %s", err)
	}

	return nil
}

func (n *ProactiveNode) receiveStreamFunc(
	ctx context.Context,
	msg *dynamo.Message,
) func(context.Context, []byte) error {
	log := env.Log(ctx)
	log.Info("[dag] [proactive] receive stream func")

	buffer := &bytes.Buffer{}

	return func(ctx context.Context, chunk []byte) error {
		_, err := buffer.Write(chunk)
		if err != nil {
			log.Error("error writing to buffer", "error", err)
			return err
		}

		_, content, err := utils.ParseIncompleteJSON(buffer.String())
		if err != nil {
			log.Error("error parsing content", "error", err)
			return err
		}

		contentSize := msg.GetContentSize()

		msg.Content = content["message"]

		if contentSize != msg.GetContentSize() {
			err = messaging.UpdateMessage(ctx, msg)
			if err != nil {
				return fmt.Errorf("error updating message: %s", err)
			}
		}

		return nil
	}
}
