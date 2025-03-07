package dag

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html"
	"text/template"

	"github.com/fivetentaylor/pointy/pkg/constants"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/models"
	"github.com/fivetentaylor/pointy/pkg/service/messaging"
	"github.com/fivetentaylor/pointy/pkg/stackerr"
	"github.com/fivetentaylor/pointy/pkg/storage/dynamo"
	"github.com/fivetentaylor/pointy/pkg/utils"
	v3 "github.com/fivetentaylor/pointy/rogue/v3"
	"github.com/tmc/langchaingo/llms"
)

type ReviseNode struct {
	Next Node
	BaseLLMNode
}

const reviseDocument = `<document>
{{.FullDocument}}
</document>
`

const reviseMessages = `<user_messages>
{{- range $val := .Messages}}
<message from="{{$val.From}}">
{{- if $val.SelectedContent }}
<selected_content>
{{ $val.SelectedContent }}
</selected_content>
{{- end }}
<content>
{{$val.Content}}
</content>
{{- range $doc := .AttachedDocs}}
<attached_document identifier="{{ $doc.Identifier }}">
{{ $doc.Content }}
</attached_document>
{{- end}}	
</message>
{{- end}}	
</user_messages>
`

const reviseSubsections = `{{- range $target := .EditTargets}}
{{-  if eq $target.Action "prepend" }}
To prepend content to the document use:
<updated_subsection id="{{$target.ID}}" action="{{$target.Action}}"></updated_subsection>
{{- else if eq $target.Action "append" }}
To append content to the document use:
<updated_subsection id="{{$target.ID}}" action="{{$target.Action}}"></updated_subsection>
{{- else }}
<updated_subsection id="{{$target.ID}}" action="{{$target.Action}}">
{{$target.Markdown}}
</updated_subsection>
{{- end }}

{{- end }}
`

const reviseOutputFormat = `{{ if .RequestExplanation }}
<explanation>
[Optional message to explain to the user that you lack access to an external resource required to meet their request. (For example, access to the internet.) Suggest that they add the content as an attachment to the message.]
</explanation>
{{ end }}

<updated_subsection id="[Id of the subsection to updated]" action="[replace, prepend, or append]">>
[Insert the revised subsection here in Markdown, incorporating user feedback and improvements]
</updated_subsection>

{{ if .RequestConclusion }}
<conclusion>
[Insert a conclusion message to the user, summarizing your actions. Be conversational, succinct, and helpful.]
</conclusion>
{{ end }}
`

var reviseDocumentTemplate = template.Must(template.New("revise").Parse(reviseDocument))

var reviseMessagesTemplate = template.Must(template.New("revise").Parse(reviseMessages))

var reviseSubsectionsTemplate = template.Must(template.New("revise").Parse(reviseSubsections))

var reviseOutputFormatTemplate = template.Must(template.New("revise").Parse(reviseOutputFormat))

type ReviseNodePromptData struct {
	FullDocument       string
	EditTargets        []EditTarget
	Messages           []reviseNodePromptMessage
	RequestExplanation bool
	RequestConclusion  bool
}

type reviseNodePromptMessage struct {
	From            string
	SelectedContent *string
	Content         string
	AttachedDocs    []*AttachedDocument
}

type reviseNodeAttachedDocument struct {
	Title   string
	Content string
}

type ReviseNodeInput struct {
	DocId           string       `key:"docId"`
	InputMessageId  string       `key:"inputMessageId"`
	OutputMessageId string       `key:"outputMessageId"`
	ThreadId        string       `key:"threadId"`
	AuthorId        string       `key:"authorId"`
	EditTargets     []EditTarget `key:"editTargets"`
}

func (n *ReviseNode) Run(ctx context.Context) (Node, error) {
	log := env.Log(ctx)

	n.AllowLLMOverride = true

	input := &ReviseNodeInput{}
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
		log.Error("error getting document at message", "error", err)
		return nil, fmt.Errorf("error getting document at message: %s", err)
	}

	data, err := n.ReviseData(ctx, input, document)
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

	appliedUpdateIds := make(map[string]bool)

	_, err = n.GenerateStoredPrompt(
		ctx,
		input.ThreadId, "threadv2-revise",
		data,
		llms.WithStreamingFunc(
			n.receiveStreamFunc(ctx, input, document, outMsg, appliedUpdateIds),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("error generating stored prompt: %s", err)
	}

	log.Info("👾 running revise node", "input", input, "output", outMsg)

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

func (n *ReviseNode) ReviseData(ctx context.Context, input *ReviseNodeInput, doc *v3.Rogue) (map[string]string, error) {
	log := env.SLog(ctx)
	mkdown, err := doc.GetFullMarkdown()
	if err != nil {
		return nil, fmt.Errorf("error getting markdown: %s", err)
	}

	escapedMkdown := html.EscapeString(mkdown)

	data := ReviseNodePromptData{
		FullDocument:       escapedMkdown,
		EditTargets:        input.EditTargets,
		Messages:           make([]reviseNodePromptMessage, 0),
		RequestExplanation: false,
		RequestConclusion:  true,
	}

	attachedDocuments, err := GetAttachedDocuments(ctx, input.ThreadId, input.InputMessageId)
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
		if message.MessageID == input.InputMessageId {
			tmplMsg.AttachedDocs = attachedDocuments
		}

		data.Messages = append(data.Messages, tmplMsg)
	}

	saveLogFile(ctx, "revise-data.json", data)

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

func (n *ReviseNode) receiveStreamFunc(
	ctx context.Context,
	input *ReviseNodeInput,
	document *v3.Rogue,
	msg *dynamo.Message,
	appliedUpdateIds map[string]bool,
) func(context.Context, []byte) error {
	log := env.SLog(ctx)

	buffer := bytes.NewBuffer([]byte{})

	var explanation *models.Content
	var conclusion *models.Content

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

		// log.Info("[dag] received stream", "buffer", buffer.String())

		xmlDoc, err := utils.ParseIncompleteXML(buffer.String())
		if err != nil {
			log.Error("[dag] error parsing xml", "error", err)
		}

		e := xmlDoc.Find("explanation")
		if e != nil {
			if explanation == nil {
				msg.LifecycleReason = ""
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
		}

		c := xmlDoc.Find("conclusion")
		if c != nil {
			if conclusion == nil {
				msg.LifecycleStage = dynamo.MessageLifecycleStageRevised
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

			log.Info("[dag] Updating subsection", "id", id)

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

func (n *ReviseNode) storeContentAddress(ctx context.Context, document *v3.Rogue, msg *dynamo.Message) error {
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
