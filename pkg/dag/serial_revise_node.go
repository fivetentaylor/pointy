package dag

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/fivetentaylor/pointy/pkg/constants"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/models"
	"github.com/fivetentaylor/pointy/pkg/service/messaging"
	"github.com/fivetentaylor/pointy/pkg/storage/dynamo"
	"github.com/fivetentaylor/pointy/pkg/utils"
	v3 "github.com/fivetentaylor/pointy/rogue/v3"

	"github.com/tmc/langchaingo/llms"
)

type SerialReviseNode struct {
	Next Node
	BaseLLMNode
}

type serialReviseNodeInput struct {
	DocId           string `key:"docId"`
	InputMessageId  string `key:"inputMessageId"`
	OutputMessageId string `key:"outputMessageId"`
	ThreadId        string `key:"threadId"`
	AuthorId        string `key:"authorId"`
}

func (n *SerialReviseNode) Run(ctx context.Context) (Node, error) {
	log := env.Log(ctx)

	input := &serialReviseNodeInput{}
	err := n.Hydrate(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error hydrating: %s", err)
	}

	outMsg, err := env.Dynamo(ctx).GetAiThreadMessage(input.ThreadId, input.OutputMessageId)
	if err != nil {
		log.Error("error getting message", "error", err)
		return nil, fmt.Errorf("error getting message: %s", err)
	}

	document, err := GetDocumentAtMessageID(ctx, input.DocId, input.AuthorId, input.ThreadId, input.InputMessageId)
	if err != nil {
		return nil, fmt.Errorf("error getting document: %s", err)
	}

	chunks, err := ChunkDocument(document, constants.MaxTokenLength)
	if err != nil {
		return nil, fmt.Errorf("error chunking document: %s", err)
	}

	saveLogFile(ctx, "chunks.json", chunks)

	outMsg.LifecycleStage = dynamo.MessageLifecycleStageRevising

	for i, chunk := range chunks {
		outMsg.LifecycleReason = fmt.Sprintf("Revising %d of %d", i+1, len(chunks))
		err = messaging.UpdateMessage(ctx, outMsg)
		if err != nil {
			return nil, fmt.Errorf("error updating message before streaming: %s", err)
		}

		data, err := n.serialReviseData(ctx, i, len(chunks), input, document, chunk)
		if err != nil {
			log.Error("error generating prompt", "error", err)
			return nil, fmt.Errorf("error generating prompt: %s", err)
		}

		_, err = n.GenerateStoredPrompt(
			ctx,
			input.ThreadId, "threadv2-revise",
			data,
			llms.WithStreamingFunc(
				n.receiveStreamFunc(ctx, input, document, outMsg, chunk),
			),
		)
	}

	go saveLogFile(ctx, "reply.json", outMsg)
	outMsg.LifecycleStage = dynamo.MessageLifecycleStageRevised

	err = n.storeContentAddress(ctx, input, document, outMsg)
	if err != nil {
		return nil, fmt.Errorf("error storing content address: %s", err)
	}

	return n.Next, nil
}

func (n *SerialReviseNode) serialReviseData(
	ctx context.Context, i int, iTotal int,
	input *serialReviseNodeInput, doc *v3.Rogue, chunk EditTarget,
) (map[string]string, error) {
	mkdown, err := doc.GetFullMarkdown()
	if err != nil {
		return nil, fmt.Errorf("error getting markdown: %s", err)
	}

	data := ReviseNodePromptData{
		FullDocument:       mkdown,
		EditTargets:        []EditTarget{chunk},
		Messages:           make([]reviseNodePromptMessage, 0),
		RequestExplanation: false,
		RequestConclusion:  i == iTotal-1,
	}

	messages, err := env.Dynamo(ctx).GetMessagesForThread(input.ThreadId)
	if err != nil {
		return nil, fmt.Errorf("error loading messages: %s", err)
	}

	for _, message := range messages {
		if message.LifecycleStage != dynamo.MessageLifecycleStageCompleted {
			continue
		}

		tmplMsg := reviseNodePromptMessage{
			Content: message.FullContent(),
		}

		for _, attachment := range message.Attachments.Attachments {
			switch v := attachment.Value.(type) {
			case *models.Attachment_Document:
				tmplMsg.SelectedContent = &v.Document.Content
			}
		}

		data.Messages = append(data.Messages, tmplMsg)
	}

	docBuffer := &bytes.Buffer{}
	err = reviseDocumentTemplate.Execute(docBuffer, data)
	if err != nil {
		return nil, fmt.Errorf("error executing document template: %s", err)
	}
	msgBuffer := &bytes.Buffer{}
	err = reviseMessagesTemplate.Execute(msgBuffer, data)
	if err != nil {
		return nil, fmt.Errorf("error executing messages template: %s", err)
	}
	targetBuffer := &bytes.Buffer{}
	err = reviseSubsectionsTemplate.Execute(targetBuffer, data)
	if err != nil {
		return nil, fmt.Errorf("error executing subsections template: %s", err)
	}
	outputBuffer := &bytes.Buffer{}
	err = reviseOutputFormatTemplate.Execute(outputBuffer, data)
	if err != nil {
		return nil, fmt.Errorf("error executing output template: %s", err)
	}

	out := map[string]string{
		"FullDocument": docBuffer.String(),
		"Messages":     msgBuffer.String(),
		"Targets":      targetBuffer.String(),
		"OutputFormat": outputBuffer.String(),
	}

	return out, nil
}

func (n *SerialReviseNode) receiveStreamFunc(
	ctx context.Context,
	input *serialReviseNodeInput,
	document *v3.Rogue,
	msg *dynamo.Message,
	et EditTarget,
) func(context.Context, []byte) error {
	log := env.Log(ctx)

	buffer := bytes.NewBuffer([]byte{})

	var explanation *models.Content
	var conclusion *models.Content

	return func(ctx context.Context, chunk []byte) error {
		_, err := buffer.Write(chunk)
		if err != nil {
			log.Error("[dag] error writing to buffer", "error", err)
			return err
		}

		xmlDoc, err := utils.ParseIncompleteXML(buffer.String())
		if err != nil {
			log.Error("[dag] error parsing xml", "error", err)
		}

		e := xmlDoc.Find("explanation")
		if e != nil {
			if explanation == nil {
				explanation = &models.Content{
					Role: "explanation",
				}
				msg.Attachments.Attachments = append(msg.Attachments.Attachments, &models.Attachment{
					Value: &models.Attachment_Content{
						Content: explanation,
					},
				})
			}

			explanation.Text = e.Value

			err = messaging.UpdateMessage(ctx, msg)
			if err != nil {
				log.Error("[dag] error publishing message", "error", err)
			}
		}

		c := xmlDoc.Find("conclusion")
		if c != nil {
			if conclusion == nil {
				conclusion = &models.Content{
					Role: "conclusion",
				}
				msg.Attachments.Attachments = append(msg.Attachments.Attachments, &models.Attachment{
					Value: &models.Attachment_Content{
						Content: conclusion,
					},
				})
			}

			conclusion.Text = c.Value

			err = messaging.UpdateMessage(ctx, msg)
			if err != nil {
				log.Error("[dag] error publishing message", "error", err)
			}
		}

		updates := xmlDoc.FindAll("updated_subsection")
		for _, update := range updates {
			if !update.Complete {
				continue
			}

			revision, err := ApplyUpdate(ctx, ApplyInput{
				DocId:      input.DocId,
				AuthorId:   input.AuthorId,
				Update:     update.RawValue,
				Document:   document,
				EditTarget: &et,
			})
			if err != nil {
				log.Error("[dag] error applying update", "error", err)
				continue
			}

			msg.Attachments.Attachments = append(
				msg.Attachments.Attachments,
				&models.Attachment{
					Value: &models.Attachment_Revision{
						Revision: revision,
					},
				})
		}

		return nil
	}
}

func (n *SerialReviseNode) applyUpdate(
	ctx context.Context,
	input *serialReviseNodeInput,
	document *v3.Rogue,
	editTarget EditTarget,
	update string,
) (*models.DocumentRevision, error) {
	// TODO: look into why we need to go vis left of the selection
	beforeID := editTarget.BeforeID
	vis, _, err := document.Rope.GetIndex(beforeID)
	if err != nil {
		return nil, fmt.Errorf("error getting index: %s", err)
	}
	if vis > 1 {
		beforeID, err = document.VisLeftOf(beforeID)
		if err != nil {
			return nil, fmt.Errorf("error getting vis left: %s", err)
		}
	}

	mop, _, err := document.ApplyMarkdownDiff(input.AuthorId, update, beforeID, editTarget.AfterID)
	if err != nil {
		return nil, fmt.Errorf("error applying diff: %s", err)
	}

	err = PublishOp(ctx, input.DocId, mop)
	if err != nil {
		return nil, fmt.Errorf("error publishing op: %s", err)
	}

	revision := &models.DocumentRevision{
		Start:   beforeID.String(),
		End:     editTarget.AfterID.String(),
		Updated: update,
		Id:      editTarget.ID,
	}

	return revision, nil
}

func (n *SerialReviseNode) storeContentAddress(ctx context.Context, input *serialReviseNodeInput, document *v3.Rogue, msg *dynamo.Message) error {
	if document.VisSize <= 1 {
		return nil
	}

	firstId, err := document.GetFirstTotID()
	if err != nil {
		return fmt.Errorf("error getting first id: %s", err)
	}

	lastId, err := document.GetLastTotID()
	if err != nil {
		return fmt.Errorf("error getting last id: %s", err)
	}

	address, err := document.GetAddress(firstId, lastId)
	if err != nil {
		return fmt.Errorf("error getting address: %s", err)
	}

	bts, err := json.Marshal(address)
	if err != nil {
		return fmt.Errorf("error marshaling address: %s", err)
	}

	msg.MessageMetadata.ContentAddress = string(bts)

	err = messaging.UpdateMessage(ctx, msg)
	if err != nil {
		return fmt.Errorf("error updating message: %s", err)
	}

	return nil
}
