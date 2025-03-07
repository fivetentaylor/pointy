package dag

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/fivetentaylor/pointy/pkg/constants"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/models"
	"github.com/fivetentaylor/pointy/pkg/service/messaging"
	"github.com/fivetentaylor/pointy/pkg/storage/dynamo"
	"github.com/fivetentaylor/pointy/pkg/utils"
	"github.com/tmc/langchaingo/llms"
)

const askDocument = `<document>
{{.FullDocument}}
</document>`

const askMessages = `<user_messages>
{{- range $val := .Messages}}
<message from="{{$val.From}}">
{{- if $val.SelectedContent }}
<selected_content>
{{ $val.SelectedContent }}
</selected_content>
{{- end }}
{{- if $val.UpdatedContent }}
<updated_content>
{{ $val.UpdatedContent }}
</updated_content>
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
</user_messages>`

type AskPromptTemplateData struct {
	FullDocument string
	Messages     []askNodePromptMessage
}

type askNodePromptMessage struct {
	From            string
	SelectedContent *string
	UpdatedContent  *string
	Content         string
	AttachedDocs    []*AttachedDocument
}

type askNodeAttachedDocument struct {
	Title   string
	Content string
}

var askDocumentTemplate = template.Must(template.New("ask").Parse(askDocument))
var askMessagesTemplate = template.Must(template.New("ask").Parse(askMessages))

type AskNode struct {
	Next Node
	BaseLLMNode
}

type AskNodeInput struct {
	DocId           string `key:"docId"`
	InputMessageId  string `key:"inputMessageId"`
	OutputMessageId string `key:"outputMessageId"`
	ThreadId        string `key:"threadId"`
	AuthorId        string `key:"authorId"`
}

func (n *AskNode) Run(ctx context.Context) (Node, error) {
	log := env.Log(ctx)

	input := &AskNodeInput{}
	err := n.Hydrate(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error hydrating: %s", err)
	}

	outMsg, err := env.Dynamo(ctx).GetAiThreadMessage(input.ThreadId, input.OutputMessageId)
	if err != nil {
		log.Error("error getting output message", "error", err)
		return nil, fmt.Errorf("error getting output message: %s", err)
	}

	err = messaging.UpdateMessage(ctx, outMsg)
	if err != nil {
		return nil, fmt.Errorf("error updating message: %s", err)
	}

	log.Info("answering", "docId", input.DocId, "authorID", input.AuthorId, "threadId", input.ThreadId)

	data, err := n.buildPromptData(ctx, input)
	if err != nil {
		log.Error("error building prompt data", "error", err)
		return nil, fmt.Errorf("error building prompt data: %s", err)
	}

	buffer := bytes.NewBuffer([]byte{})
	_, err = n.GenerateStoredPrompt(
		ctx,
		input.ThreadId,
		"threadv2-ask",
		data,
		llms.WithStreamingFunc(
			n.receiveStreamFunc(ctx, buffer, outMsg),
		),
	)
	if err != nil {
		log.Error("[dag] error generating", "error", err)
		return nil, fmt.Errorf("error generating: %s", err)
	}

	return n.Next, nil
}

func (n *AskNode) receiveStreamFunc(
	ctx context.Context,
	buffer *bytes.Buffer,
	msg *dynamo.Message,
) func(context.Context, []byte) error {
	log := env.Log(ctx)

	attachment := &models.Content{
		Role: "answer",
	}
	msg.Attachments.Attachments = append(msg.Attachments.Attachments, &models.Attachment{
		Value: &models.Attachment_Content{
			Content: attachment,
		},
	})

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

		m := xmlDoc.Find("message")
		if m == nil {
			return nil
		}

		if c := m.Find("content"); c != nil {
			attachment.Text = c.Value
		} else {
			attachment.Text = m.Value
		}

		msg.LifecycleReason = ""

		err = messaging.UpdateMessage(ctx, msg)
		if err != nil {
			log.Error("[dag] error publishing message", "error", err)
		}

		return nil
	}
}

func (n *AskNode) buildPromptData(ctx context.Context, input *AskNodeInput) (map[string]string, error) {
	log := env.Log(ctx)

	attachedDocuments, err := GetAttachedDocuments(ctx, input.ThreadId, input.InputMessageId)
	if err != nil {
		log.Error("error getting attached documents", "error", err)
		return nil, fmt.Errorf("error getting attached documents: %s", err)
	}

	messages, err := env.Dynamo(ctx).GetMessagesForThread(input.ThreadId)
	if err != nil {
		log.Error("error loading messages", "error", err)
		return nil, fmt.Errorf("error loading messages: %s", err)
	}

	data := &AskPromptTemplateData{
		Messages: make([]askNodePromptMessage, 0, len(messages)),
	}

	for _, message := range messages {
		if message.LifecycleStage != dynamo.MessageLifecycleStageCompleted {
			continue
		}

		from := "assistant"
		if message.UserID != constants.RevisoUserID {
			from = "user"
		}

		tmplMsg := askNodePromptMessage{
			From:    from,
			Content: message.FullContent(),
		}

		for _, attachment := range message.Attachments.Attachments {
			switch v := attachment.Value.(type) {
			case *models.Attachment_Document:
				tmplMsg.SelectedContent = &v.Document.Content
			case *models.Attachment_Revision:
				tmplMsg.UpdatedContent = &v.Revision.Updated
			}
		}

		// only show attached documents if it's the input message
		if message.MessageID == input.InputMessageId {
			tmplMsg.AttachedDocs = attachedDocuments
		}

		data.Messages = append(data.Messages, tmplMsg)
	}

	if len(data.Messages) == 0 {
		return nil, fmt.Errorf("no messages found")
	}

	doc, err := GetDocument(ctx, input.DocId, "reviso")
	if err != nil {
		log.Error("error loading document", "error", err)
		return nil, fmt.Errorf("error loading document: %s", err)
	}
	mkdown, err := doc.GetFullMarkdown()
	if err != nil {
		log.Error("error getting markdown", "error", err)
		return nil, fmt.Errorf("error getting markdown: %s", err)
	}

	data.FullDocument = mkdown

	docBuffer := &bytes.Buffer{}
	err = askDocumentTemplate.Execute(docBuffer, data)
	if err != nil {
		return nil, fmt.Errorf("error executing document template: %s", err)
	}
	msgBuffer := &bytes.Buffer{}
	err = askMessagesTemplate.Execute(msgBuffer, data)
	if err != nil {
		return nil, fmt.Errorf("error executing messages template: %s", err)
	}

	out := map[string]string{
		"FullDocument": docBuffer.String(),
		"Messages":     msgBuffer.String(),
	}

	return out, nil

}
