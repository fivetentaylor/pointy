package dag

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fivetentaylor/pointy/pkg/constants"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/models"
	"github.com/fivetentaylor/pointy/pkg/service/timeline"
	"github.com/fivetentaylor/pointy/pkg/storage/dynamo"
	"github.com/fivetentaylor/pointy/pkg/utils"
	v3 "github.com/fivetentaylor/pointy/rogue/v3"
	"github.com/tmc/langchaingo/llms"
)

var ErrorNoDifference = fmt.Errorf("no difference")

type SummarizeUpdateNode struct {
	Next Node
	BaseLLMNode
}

type SummarizeUpdateNodeInputs struct {
	SessionId string `key:"sessionId"`
	DocId     string `key:"docId"`
	UserId    string `key:"userId"`
}

func (n *SummarizeUpdateNode) Run(ctx context.Context) (Node, error) {
	log := env.Log(ctx)
	dydb := env.Dynamo(ctx)

	input := &SummarizeUpdateNodeInputs{}
	err := n.Hydrate(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error hydrating: %s", err)
	}

	log.Info(
		"[dag] summarizing update",
		"sessionId", input.SessionId,
		"docId", input.DocId,
	)

	document, err := GetDocument(ctx, input.DocId, "")
	if err != nil {
		log.Error("error getting document", "error", err)
		return nil, fmt.Errorf("error getting document: %s", err)
	}

	endingAddress, err := document.GetFullAddress()
	if err != nil {
		return nil, fmt.Errorf("error getting full address: %w", err)
	}

	startingAddress, err := document.GetEmptyAddress()
	if err != nil {
		return nil, fmt.Errorf("error getting empty address: %s", err)
	}

	lastCompletedUpdate, err := dydb.GetLastCompletedUserUpdate(input.DocId, input.UserId)
	if err != nil {
		return nil, fmt.Errorf("error getting last update: %s", err)
	}
	if lastCompletedUpdate != nil {
		// if there was a last update, get the content address from it it is where the current update will start
		lastAddress := lastCompletedUpdate.Event.GetUpdate().GetEndingContentAddress()
		err = json.Unmarshal([]byte(lastAddress), startingAddress)
		if err != nil {

			return nil, fmt.Errorf("error unmarshaling content address: %s", err)
		}
	}

	diff, err := SummarizeDocDifference(ctx, input.DocId, input.UserId, startingAddress, endingAddress)
	if err != nil {
		if err == ErrorNoDifference {
			return n.Next, nil
		}
		log.Error("error generating summary", "error", err)
		return nil, fmt.Errorf("error generating summary: %s", err)
	}

	event := &dynamo.TimelineEvent{
		DocID:  input.DocId,
		UserID: input.UserId,
		Event: &models.TimelineEventPayload{
			Payload: &models.TimelineEventPayload_Update{
				Update: &models.TimelineDocumentUpdateV1{
					State: models.UpdateState_SUMMARIZING_STATE,
				},
			},
		},
	}

	err = timeline.CreateTimelineEvent(ctx, event)
	if err != nil {
		return nil, fmt.Errorf("error saving timeline event: %s", err)
	}

	completion, err := n.GenerateContentForStoredPrompt(ctx, input.SessionId, constants.PromptSummarizeDocDiff, []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, diff),
	})
	if err != nil {
		log.Error("error generating content for freeplay prompt", "error", err)
		return nil, fmt.Errorf("error generating content for freeplay prompt: %s", err)
	}

	log.Info("[dag] got summary", "summary", completion.Choices[0].Content)

	startAddr, err := json.Marshal(startingAddress)
	if err != nil {
		log.Error("error parsing content address", "error", err)
		return nil, fmt.Errorf("error parsing content address: %s", err)
	}
	endAddr, err := json.Marshal(endingAddress)
	if err != nil {
		log.Error("error parsing content address", "error", err)
		return nil, fmt.Errorf("error parsing content address: %s", err)
	}

	event.Event.GetUpdate().StartingContentAddress = string(startAddr)
	event.Event.GetUpdate().EndingContentAddress = string(endAddr)
	event.Event.GetUpdate().Content = completion.Choices[0].Content
	event.Event.GetUpdate().State = models.UpdateState_COMPLETE_STATE

	err = timeline.UpdateTimelineEvent(ctx, event)
	if err != nil {
		log.Error("error updating timeline event", "error", err)
		return nil, fmt.Errorf("error updating timeline event: %s", err)
	}

	return n.Next, nil
}

func SummarizeDocDifference(ctx context.Context, docId, userId string, startAddress, endingAddress *v3.ContentAddress) (string, error) {
	log := env.Log(ctx)
	document, err := GetDocument(ctx, docId, "")
	if err != nil {
		log.Error("error getting document", "error", err)
		return "", fmt.Errorf("error getting document: %s", err)
	}

	aIdTbl := env.Query(ctx).AuthorID
	usersAuthorIds, err := aIdTbl.Where(
		aIdTbl.UserID.Eq(userId),
		aIdTbl.DocumentID.Eq(docId),
	).Find()
	if err != nil {
		log.Error("error getting author ids", "error", err)
		return "", fmt.Errorf("error getting author ids: %s", err)
	}

	authorIds := make([]string, 0, len(usersAuthorIds))
	for _, authorId := range usersAuthorIds {
		authorIds = append(authorIds, authorId.AuthorIDString())
	}

	for a, v := range endingAddress.MaxIDs {
		// do not change the author's max ID or it's AI author's max ID
		if utils.Contains(authorIds, a) || utils.Contains(authorIds, strings.TrimPrefix(a, "!")) {
			continue
		}
		startAddress.MaxIDs[a] = v
	}

	log.Info("comparing", "start", startAddress, "end", endingAddress)

	beforeMkd, err := document.GetMarkdownAt(startAddress.StartID, startAddress.EndID, *startAddress)
	if err != nil {
		return "", fmt.Errorf("error getting before mkd: %s", err)
	}

	// log.Info("before mkd", "mkd", beforeMkd)

	afterMkd, err := document.GetMarkdownAt(endingAddress.StartID, endingAddress.EndID, *endingAddress)
	if err != nil {
		return "", fmt.Errorf("error getting after mkd: %s", err)
	}

	// log.Info("after mkd", "mkd", afterMkd)

	if beforeMkd == afterMkd {
		log.Info("no difference")
		return "", ErrorNoDifference
	}

	return fmt.Sprintf("Before:\n```\n%s\n```\n\nAfter:\n```\n%s\n```", beforeMkd, afterMkd), nil
}
