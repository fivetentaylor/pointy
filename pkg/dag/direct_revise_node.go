package dag

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/teamreviso/code/pkg/constants"
	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/models"
	"github.com/teamreviso/code/pkg/service/messaging"
	"github.com/teamreviso/code/pkg/stackerr"
	"github.com/teamreviso/code/pkg/storage/dynamo"
	"github.com/teamreviso/code/pkg/utils"
	v3 "github.com/teamreviso/code/rogue/v3"
	"github.com/tmc/langchaingo/llms"
)

type DirectReviseNode struct {
	Next Node
	BaseLLMNode
}

type DirectReviseNodeInput struct {
	DocId       string       `key:"docId"`
	ThreadId    string       `key:"threadId"`
	AuthorId    string       `key:"authorId"`
	EditTargets []EditTarget `key:"editTargets"`

	InputMessage  *dynamo.Message `key:"inputMessage"`
	OutputMessage *dynamo.Message `key:"outputMessage"`
}

func (n *DirectReviseNode) Run(ctx context.Context) (Node, error) {
	log := env.Log(ctx)

	n.AllowLLMOverride = true

	input := &DirectReviseNodeInput{}
	err := n.Hydrate(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error hydrating: %s", err)
	}

	outMsg := input.OutputMessage

	document, err := GetDocumentAtMessage(ctx, input.DocId, input.AuthorId, input.InputMessage)
	if err != nil {
		log.Error("error getting document at message", "error", err)
		return nil, fmt.Errorf("error getting document at message: %s", err)
	}

	data, err := n.DirectReviseData(ctx, input, document)
	if err != nil {
		log.Error("error generating prompt", "error", err)
		return nil, fmt.Errorf("error generating prompt: %s", err)
	}

	if len(input.EditTargets) == 0 {
		outMsg.LifecycleReason = "Updating the draft..."
	} else {
		outMsg.LifecycleReason = "Making selected changes..."
	}
	outMsg.LifecycleStage = dynamo.MessageLifecycleStagePending
	err = messaging.UpdateMessage(ctx, outMsg)
	if err != nil {
		return nil, fmt.Errorf("error updating message before streaming: %s", err)
	}

	state := GetState(ctx)
	appliedUpdateIds := make(map[string]bool)

	log.Info("👾 running direct revise node", "input", input, "output", outMsg)

	_, err = n.GenerateStoredPrompt(
		ctx,
		input.ThreadId, "threadv2-revise",
		data,
		llms.WithStreamingFunc(
			n.receiveStreamFunc(ctx, input, document, outMsg, state, appliedUpdateIds),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("error generating stored prompt: %s", err)
	}

	go saveLogFile(ctx, "reply.json", outMsg)

	outMsg.LifecycleStage = dynamo.MessageLifecycleStageRevised

	if len(appliedUpdateIds) == 0 {
		// If no updates were applied, add error to message and return
		errorAttachment := &models.Error{
			Title: "Error",
			Text:  "Sorry our system was unable to respond to your message. Please try again.",
			Error: "appliedUpdateIds is empty",
		}
		outMsg.Attachments.Attachments = append(outMsg.Attachments.Attachments, &models.Attachment{
			Value: &models.Attachment_Error{
				Error: errorAttachment,
			},
		})
		err = messaging.UpdateMessage(ctx, outMsg)
		if err != nil {
			return nil, fmt.Errorf("error updating message on error: %s", err)
		}

		return nil, nil
	}

	err = n.storeContentAddress(ctx, document, outMsg)
	if err != nil {
		return nil, fmt.Errorf("error storing content address: %s", err)
	}

	err = messaging.UpdateMessage(ctx, outMsg)
	if err != nil {
		return nil, fmt.Errorf("error updating message: %s", err)
	}

	return n.Next, nil
}

func (n *DirectReviseNode) DirectReviseData(ctx context.Context, input *DirectReviseNodeInput, doc *v3.Rogue) (map[string]string, error) {
	log := env.SLog(ctx)
	mkdown, err := doc.GetFullMarkdown()
	if err != nil {
		log.Error("[dag.DirectReviseNode.DirectReviseData] error getting markdown", "error", err)
		return nil, fmt.Errorf("error getting markdown: %s", err)
	}

	data := ReviseNodePromptData{
		FullDocument:       mkdown,
		EditTargets:        input.EditTargets,
		Messages:           make([]reviseNodePromptMessage, 0),
		RequestExplanation: false,
		RequestConclusion:  true,
	}

	attachedDocuments, err := GetAttachedDocuments(ctx, input.ThreadId, input.InputMessage.MessageID)
	if err != nil {
		log.Error("error getting attached documents", "error", err)
		return nil, fmt.Errorf("error getting attached documents: %s", err)
	}

	messages, err := env.Dynamo(ctx).GetMessagesForThread(input.ThreadId)
	if err != nil {
		return nil, fmt.Errorf("error loading messages: %s", err)
	}

	for _, message := range messages {
		if message.LifecycleStage != dynamo.MessageLifecycleStageCompleted {
			continue
		}

		from := "assistant"
		if message.UserID != constants.RevisoUserID {
			from = "user"
		}

		tmplMsg := reviseNodePromptMessage{
			From:    from,
			Content: message.FullContent(),
		}

		for _, attachment := range message.Attachments.Attachments {
			switch v := attachment.Value.(type) {
			case *models.Attachment_Document:
				tmplMsg.SelectedContent = &v.Document.Content
			}
		}

		// only show attached documents if it's the input message
		if message.MessageID == input.InputMessage.MessageID {
			tmplMsg.AttachedDocs = attachedDocuments
		}

		data.Messages = append(data.Messages, tmplMsg)
	}

	saveLogFile(ctx, "directUpdate-data.json", data)

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

func (n *DirectReviseNode) receiveStreamFunc(
	ctx context.Context,
	input *DirectReviseNodeInput,
	document *v3.Rogue,
	msg *dynamo.Message,
	state *State,
	appliedUpdateIds map[string]bool,
) func(context.Context, []byte) error {
	log := env.SLog(ctx)

	buffer := bytes.NewBuffer([]byte{})

	editTargets := make(map[string]EditTarget)
	for _, et := range input.EditTargets {
		editTargets[et.ID] = et
	}

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
			state.Set("revision_explanation", e.Value)
		}

		c := xmlDoc.Find("conclusion")
		if c != nil {
			state.Set("revision_conclusion", c.Value)
		}

		updates := xmlDoc.FindAll("updated_subsection")
		for _, update := range updates {
			id := update.Attr("id")

			if id != "" && appliedUpdateIds[id] {
				continue
			}

			msg.LifecycleReason = ""
			msg.LifecycleStage = dynamo.MessageLifecycleStageRevising

			if !update.Complete {
				continue
			}

			var revision *models.DocumentRevision
			if id == "full_document" {
				beforeID, err := document.GetFirstTotID()
				if err != nil {
					log.Error("[dag] error getting first id", "error", err)
					return stackerr.Wrap(err)
				}

				afterID, err := document.GetLastTotID()
				if err != nil {
					log.Error("[dag] error getting last id", "error", err)
					return stackerr.Wrap(err)
				}

				editTarget := &EditTarget{
					ID:       id,
					BeforeID: beforeID,
					AfterID:  afterID,
				}

				action := update.Attr("action")
				if action == "prepend" || action == "append" {
					editTarget.Action = EditTargetAction(action)
				} else {
					editTarget.Action = EditTargetActionReplace
				}

				revision, err = ApplyUpdate(ctx, ApplyInput{
					DocId:      input.DocId,
					AuthorId:   input.AuthorId,
					Update:     update.RawValue,
					Document:   document,
					EditTarget: editTarget,
				})
				if err != nil {
					log.Error("[dag] error applying full update", "error", err)
					continue
				}
			} else {
				et, ok := editTargets[id]
				if !ok {
					log.Error("[dag] edit target not found!!", "id", id)
					continue
				}

				revision, err = ApplyUpdate(ctx, ApplyInput{
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
			}

			msg.Attachments.Attachments = append(
				msg.Attachments.Attachments,
				&models.Attachment{
					Value: &models.Attachment_Revision{
						Revision: revision,
					},
				})

			appliedUpdateIds[id] = true
		}

		err = messaging.UpdateMessage(ctx, msg)
		if err != nil {
			log.Error("[dag] error publishing message", "error", err)
		}

		return nil
	}
}

func (n *DirectReviseNode) storeContentAddress(ctx context.Context, document *v3.Rogue, msg *dynamo.Message) error {
	log := env.SLog(ctx)
	if document.VisSize <= 1 {
		return nil
	}

	firstId, err := document.GetFirstTotID()
	if err != nil {
		log.Error("[dag] error getting first id", "error", err)
		return fmt.Errorf("error getting first id: %s", err)
	}

	lastId, err := document.GetLastTotID()
	if err != nil {
		log.Error("[dag] error getting last id", "error", err)
		return fmt.Errorf("error getting last id: %s", err)
	}

	address, err := document.GetAddress(firstId, lastId)
	if err != nil {
		log.Error("[dag] error getting address", "error", err)
		return fmt.Errorf("error getting address: %s", err)
	}

	bts, err := json.Marshal(address)
	if err != nil {
		log.Error("[dag] error marshaling address", "error", err)
		return fmt.Errorf("error marshaling address: %s", err)
	}

	msg.MessageMetadata.ContentAddress = string(bts)

	return nil
}
