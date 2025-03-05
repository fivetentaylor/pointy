package dag

import (
	"context"
	"fmt"
	"time"

	"github.com/teamreviso/code/pkg/constants"
	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/graph/loaders"
	"github.com/teamreviso/code/pkg/service/timeline"
	"github.com/teamreviso/code/pkg/storage/dynamo"
	"github.com/tmc/langchaingo/llms"
)

type SummarizeCommentThreadNode struct {
	Next Node
	BaseLLMNode
}

type SummarizeCommentThreadNodeInputs struct {
	DocId         string `key:"docId"`
	EventId       string `key:"eventId"`
	ThreadEventID string `key:"threadEventId"`
}

func (n *SummarizeCommentThreadNode) Run(ctx context.Context) (Node, error) {
	log := env.Log(ctx)
	dydb := env.Dynamo(ctx)

	input := &SummarizeCommentThreadNodeInputs{}
	if err := n.Hydrate(ctx, input); err != nil {
		return nil, fmt.Errorf("error hydrating input: %w", err)
	}

	parentEvent, err := dydb.GetTimelineEvent(input.DocId, input.ThreadEventID)
	if err != nil {
		return nil, fmt.Errorf("error getting parent event: %w", err)
	}

	allReplies, err := dydb.GetDocumentTimelineReplies(input.DocId, input.ThreadEventID)
	if err != nil {
		return nil, fmt.Errorf("error getting document timeline replies: %s", err)
	}

	// create a new slice of replies that are messages
	replies := make([]*dynamo.TimelineEvent, 0, len(allReplies))
	for _, reply := range allReplies {
		if reply.Event.GetMessage() == nil {
			continue
		}
		replies = append(replies, reply)
	}

	events := []*dynamo.TimelineEvent{parentEvent}
	events = append(events, replies...)

	if len(events) == 0 {
		return nil, fmt.Errorf("no events found")
	}

	mentionedUserIds := parentEvent.Event.GetMessage().GetMentionedUserIds()
	log.Info("[dag] Parent event mentioned user ids", "mentionedUserIds", mentionedUserIds)
	for _, reply := range replies {
		mentionedUserIds = append(mentionedUserIds, reply.Event.GetMessage().GetMentionedUserIds()...)
	}
	log.Info("[dag] All mentioned user ids", "mentionedUserIds", mentionedUserIds)

	mentionedUsers, err := loaders.GetUsers(ctx, mentionedUserIds)
	if err != nil {
		return nil, fmt.Errorf("error loading mentioned users: %w", err)
	}
	log.Info("[dag] Mentioned users", "mentionedUsers", mentionedUsers)

	uniqueAuthors := make(map[string]string)
	for _, event := range events {
		uniqueAuthors[event.UserID] = ""
	}

	authorIds := make([]string, 0, len(uniqueAuthors))
	for authorId := range uniqueAuthors {
		authorIds = append(authorIds, authorId)
	}

	authors, err := loaders.GetUsers(ctx, authorIds)
	if err != nil {
		return nil, fmt.Errorf("error loading authors: %w", err)
	}

	for _, author := range authors {
		uniqueAuthors[author.ID] = author.Name
	}

	for _, event := range events {
		message := event.Event.GetMessage()
		log.Info("[dag] Message", "message", message, "event", event)
		message.Content = timeline.UnfurlBase64Mentions(message.Content, mentionedUsers)
		log.Info("[dag] Message after unfurling", "message", message)
	}

	messages := ""
	for _, event := range events {
		createdAt := time.Unix(event.CreatedAt, 0).Format("3:04 PM on Mon Jan 2")
		messages += fmt.Sprintf("%s at %s said %s", uniqueAuthors[event.UserID], createdAt, event.Event.GetMessage().GetContent())
		if event.Event.GetMessage().SelectionMarkdown != "" {
			messages += fmt.Sprintf(" with highlighted text %s", event.Event.GetMessage().SelectionMarkdown)
		}
		messages += "\n"
	}

	if len(messages) == 0 {
		return nil, fmt.Errorf("no messages to summarize")
	}

	aiMsg := "Comment messages:\n" + messages

	completion, err := n.GenerateContentForStoredPrompt(ctx, parentEvent.EventID, constants.PromptSummarizeCommentThread, []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, aiMsg),
	})
	if err != nil {
		return nil, fmt.Errorf("error generating content for freeplay prompt: %w", err)
	}

	if completion == nil || len(completion.Choices) == 0 || completion.Choices[0].Content == "" {
		return nil, fmt.Errorf("unexpected empty completion")
	}

	summary := completion.Choices[0].Content
	log.Info("Summarized comment thread", "summary", summary)

	// update the event with the summary
	resolutionEvent, err := dydb.GetTimelineEvent(input.DocId, input.EventId)
	if err != nil {
		return nil, fmt.Errorf("error getting timeline event: %w", err)
	}

	resolutionEvent.Event.GetResolution().ResolutionSummary = summary

	err = timeline.UpdateTimelineEvent(ctx, resolutionEvent)
	if err != nil {
		return nil, err
	}

	return n.Next, nil
}
