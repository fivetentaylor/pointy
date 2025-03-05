package dag

import (
	"context"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
	"github.com/teamreviso/code/pkg/constants"
	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/service/messaging"
)

// Inputs:
// - docId
// - threadId
// - userId
type TitleThreadNode struct {
	Next Node
	Base
}

type TitleThreadInput struct {
	DocId    string `key:"docId"`
	ThreadId string `key:"threadId"`
	UserId   string `key:"userId"`
}

func (t *TitleThreadNode) Run(ctx context.Context) (Node, error) {
	log := env.Log(ctx)

	input := &TitleThreadInput{}
	err := t.Hydrate(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error hydrating: %s", err)
	}
	client := env.OpenAi(ctx)

	messages, err := env.Dynamo(ctx).GetMessagesForThread(input.ThreadId)
	if err != nil {
		return nil, fmt.Errorf("error getting messages: %s", err)
	}

	var chatThread strings.Builder
	for _, message := range messages {
		if message.AuthorID == constants.RevisoUserID {
			chatThread.WriteString("AI: ")
			chatThread.WriteString(message.Content)
			chatThread.WriteString(message.AIContent.ConcludingMessage)
			chatThread.WriteString("\n")
			continue
		}
		chatThread.WriteString(fmt.Sprintf("User: %s\n", message.Content))
	}

	thread, err := env.Dynamo(ctx).GetThreadForUser(input.DocId, input.ThreadId, input.UserId)
	if err != nil {
		log.Error("error getting thread", "error", err, "docId", input.DocId, "userId", input.UserId, "threadId", input.ThreadId)
		return nil, fmt.Errorf("error getting thread: %s", err)
	}

	//don't title thread if it already has a title
	if thread.Title != constants.DefaultThreadTitle {
		return t.Next, nil
	}

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a helpful assistant that summarizes chat threads.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: fmt.Sprintf("Summarize this chat conversation in a concise 3-10 word title, capturing the main topic discussed. Do not use quotations: %s", chatThread.String()),
				},
			},
		},
	)
	if err != nil {
		log.Error("error creating chat completion", "error", err)
		return nil, fmt.Errorf("error creating chat completion: %s", err)
	}

	thread.Title = strings.TrimSpace(
		strings.ReplaceAll(resp.Choices[0].Message.Content, "\n", ""))

	log.Info("saving thread =>", "thread", thread)

	err = messaging.UpdateThread(ctx, thread)
	if err != nil {
		log.Error("error updating thread", "error", err)
		return nil, fmt.Errorf("error updating thread: %s", err)
	}

	return t.Next, nil
}
