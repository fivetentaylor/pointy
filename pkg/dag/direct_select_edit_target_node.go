package dag

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/models"
	"github.com/fivetentaylor/pointy/pkg/service/messaging"
	"github.com/fivetentaylor/pointy/pkg/storage/dynamo"
	"github.com/fivetentaylor/pointy/pkg/utils"
	v3 "github.com/fivetentaylor/pointy/rogue/v3"
)

type DirectSelectEditTargetNode struct {
	Next    Node
	NoEdits Node
	BaseLLMNode
}

type DirectSelectEditTargetNodeInput struct {
	DocId          string `key:"docId"`
	ThreadId       string `key:"threadId"`
	AuthorId       string `key:"authorId"`
	ContentAddress string `key:"contentAddress"`

	InputMessage  *dynamo.Message `key:"inputMessage"`
	OutputMessage *dynamo.Message `key:"outputMessage"`
}

func (n *DirectSelectEditTargetNode) Run(ctx context.Context) (Node, error) {
	log := env.Log(ctx)

	input := &DirectSelectEditTargetNodeInput{}
	err := n.Hydrate(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error hydrating: %s", err)
	}
	if input.InputMessage == nil {
		return nil, fmt.Errorf("[dag.DirectSelectEditTargetNode.Run] no input message")
	}

	log.Info("Directly selecting focus", "docId", input.DocId, "threadId", input.ThreadId)

	outMsg := input.OutputMessage
	outMsg.LifecycleReason = "Reading draft..."
	err = messaging.UpdateMessage(ctx, outMsg)
	if err != nil {
		return nil, fmt.Errorf("error updating message: %s", err)
	}

	doc, err := GetDocumentAtMessage(ctx, input.DocId, input.AuthorId, input.InputMessage)
	if err != nil {
		return nil, fmt.Errorf("error getting document: %s", err)
	}

	targets := []EditTarget{}
	for _, attachment := range input.InputMessage.Attachments.Attachments {
		switch v := attachment.Value.(type) {
		case *models.Attachment_Document:
			startID, err := v3.ParseID(v.Document.Start)
			if err != nil {
				return nil, fmt.Errorf("error parsing id %q: %s", v.Document.Start, err)
			}

			endID, err := v3.ParseID(v.Document.End)
			if err != nil {
				return nil, fmt.Errorf("error parsing id %q: %s", v.Document.End, err)
			}

			ids := []v3.ID{startID, endID}
			startID, endID, err = doc.IDsToEnclosingSpan(ids, nil)
			if err != nil {
				return nil, fmt.Errorf("error getting enclosing span: %s", err)
			}

			markdown, err := doc.GetMarkdown(startID, endID)
			if err != nil {
				return nil, fmt.Errorf("error getting markdown: %s", err)
			}

			beforeID, err := doc.TotLeftOf(startID)
			if err != nil {
				if !errors.As(err, &v3.ErrorNoLeftTotSibling{}) {
					return nil, fmt.Errorf("error getting left of: %s", err)
				}
				beforeID = startID
			}

			targets = append(targets, EditTarget{
				ID:       uuid.NewString(),
				BeforeID: beforeID,
				AfterID:  endID,
				Action:   EditTargetActionReplace,
				Markdown: markdown,
			})
		}
	}

	if len(targets) > 0 {
		SetStateKey(ctx, "editTargets", targets)

		userSelections, err := json.Marshal(targets)
		if err != nil {
			return nil, fmt.Errorf("error marshalling user selections: %s", err)
		}

		go saveLogFile(ctx, fmt.Sprintf("user-selections-%d.txt", time.Now().Unix()), string(userSelections))

		return n.Next, nil
	}

	data, err := n.Data(ctx, doc, input)
	if err != nil {
		log.Error("error generating", "error", err)
		return nil, fmt.Errorf("error generating: %s", err)
	}

	resp, err := n.GenerateStoredPrompt(ctx, input.ThreadId, "threadv2-selectEditTarget", data)
	if err != nil {
		log.Error("error generating", "error", err)
		return nil, fmt.Errorf("error generating: %s", err)
	}

	if resp == nil || len(resp.Choices) == 0 {
		log.Error("error generating", "error", "no choices")
		return nil, fmt.Errorf("error generating: no choices")
	}

	llmResponse := html.UnescapeString(resp.Choices[0].Content)
	targets, err = n.EditTargets(ctx, doc, input, llmResponse)
	if err != nil {
		if err == NoSelectionError {
			return n.NoEdits, nil
		}

		// TODO: if this fails, we attach the whole document to the error
		log.Error("error updating state", "error", err)
		return nil, fmt.Errorf("error updating state: %s", err)
	}

	if len(targets) == 0 {
		outMsg.LifecycleReason = "Updating the draft..."
	} else {
		outMsg.LifecycleReason = "Making selected changes..."
	}
	err = messaging.UpdateMessage(ctx, outMsg)
	if err != nil {
		return nil, fmt.Errorf("error updating message: %s", err)
	}

	SetStateKey(ctx, "editTargets", targets)

	return n.Next, nil
}

func (n *DirectSelectEditTargetNode) Data(ctx context.Context, doc *v3.Rogue, input *DirectSelectEditTargetNodeInput) (map[string]string, error) {
	msg := input.InputMessage
	content := msg.FullContent()

	html, err := doc.GetFullHtml(true, false)
	if err != nil {
		return nil, fmt.Errorf("error getting markdown: %s", err)
	}

	data := SelectEditTargetPromptTemplateData{
		Document: html,
		Message: SelectEditTargetPromptTemplateUserMessage{
			Content:    content,
			Selections: []*models.DocumentSelection{},
		},
	}

	for _, attachment := range msg.Attachments.Attachments {
		switch v := attachment.Value.(type) {
		case *models.Attachment_Document:
			data.Message.Selections = append(data.Message.Selections, v.Document)
		}
	}

	var documentString strings.Builder
	err = SelectEditTargetDocumentTemplate.Execute(&documentString, data)
	if err != nil {
		return nil, fmt.Errorf("error executing document template: %s", err)
	}

	var messageString strings.Builder
	err = SelectEditTargetMessagesTemplate.Execute(&messageString, data)
	if err != nil {
		return nil, fmt.Errorf("error executing document template: %s", err)
	}

	out := map[string]string{
		"Document": documentString.String(),
		"Message":  messageString.String(),
	}

	return out, nil
}

func (n *DirectSelectEditTargetNode) EditTargets(ctx context.Context, doc *v3.Rogue, input *DirectSelectEditTargetNodeInput, llmResponse string) ([]EditTarget, error) {
	xml, err := utils.ParseIncompleteXML(llmResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing xml: %s", err)
	}

	response := xml.FindDeep("response")
	if response == nil {
		return nil, fmt.Errorf("no valid response from LLM: %+v", xml)
	}

	if response.FindDeep("full_document") != nil {
		return nil, nil
	}

	if response.FindDeep("no_selection") != nil {
		return nil, NoSelectionError
	}

	if response.FindDeep("append") != nil {
		lastId, err := doc.GetLastID()
		if err != nil {
			return nil, fmt.Errorf("error getting last id: %s", err)
		}

		return []EditTarget{
			{
				ID:       uuid.NewString(),
				BeforeID: lastId,
				AfterID:  lastId,
				Action:   EditTargetActionAppend,
			},
		}, nil
	}

	if response.FindDeep("prepend") != nil {
		firstId, err := doc.GetFirstID()
		if err != nil {
			return nil, fmt.Errorf("error getting last id: %s", err)
		}

		return []EditTarget{
			{
				ID:       uuid.NewString(),
				BeforeID: firstId,
				AfterID:  firstId,
				Action:   EditTargetActionPrepend,
			},
		}, nil
	}

	relevantSections := response.FindAllDeep("relevant_section")
	targets := []EditTarget{}
	for _, section := range relevantSections {
		ts, err := n.findRelevantSection(doc, input.InputMessage, section)
		if err != nil {
			return nil, fmt.Errorf("error finding relevant section: %s", err)
		}
		targets = append(targets, ts...)
	}

	return targets, nil
}

func (n *DirectSelectEditTargetNode) findRelevantSection(document *v3.Rogue, msg *dynamo.Message, section *utils.Tag) ([]EditTarget, error) {
	var contentAddress v3.ContentAddress

	err := json.Unmarshal(
		[]byte(msg.MessageMetadata.ContentAddress),
		&contentAddress,
	)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling content address: %s", err)
	}

	ids := []v3.ID{}
	stack := []*utils.Tag{section}
	for len(stack) > 0 {
		tag := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		rid := tag.Attr("data-rid")
		if rid != "" {
			id, err := v3.ParseID(rid)
			if err != nil {
				return nil, fmt.Errorf("error parsing id %q: %s", rid, err)
			}

			ids = append(ids, id)
		}
		stack = append(stack, tag.Children...)
	}

	if len(ids) == 0 {
		return []EditTarget{}, nil
	}

	startID, endID, err := document.IDsToEnclosingSpan(ids, &contentAddress)
	if err != nil {
		return nil, err
	}

	markdown, err := document.GetMarkdown(startID, endID)
	if err != nil {
		return nil, fmt.Errorf("error getting markdown: %s", err)
	}

	beforeID, err := document.TotLeftOf(startID)
	if err != nil {
		if !errors.As(err, &v3.ErrorNoLeftTotSibling{}) {
			return nil, fmt.Errorf("error getting left of: %s", err)
		}
		beforeID = startID
	}

	return []EditTarget{
		{
			ID:       uuid.NewString(),
			BeforeID: beforeID,
			AfterID:  endID,
			Action:   EditTargetActionReplace,
			Markdown: markdown,
		},
	}, nil
}
