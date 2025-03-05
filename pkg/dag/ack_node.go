package dag

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/teamreviso/code/pkg/constants"
	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/models"
	"github.com/teamreviso/code/pkg/service/messaging"
	"github.com/teamreviso/code/pkg/stackerr"
	"github.com/teamreviso/code/pkg/storage/dynamo"
	"github.com/teamreviso/code/pkg/utils"
	"github.com/tmc/langchaingo/llms"
)

const ackPrompt = `<user_messages>
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
</message>
{{- end}}	
</user_messages>
`

var ackTemplate = template.Must(template.New("ack").Parse(ackPrompt))

type AckNode struct {
	Next Node
	BaseLLMNode
}

type ackNodePromptMessage struct {
	From            string
	SelectedContent *string
	Content         string
}

type AckNodeInputs struct {
	ThreadId        string `key:"threadId"`
	InputMessageId  string `key:"inputMessageId"`
	OutputMessageId string `key:"outputMessageId"`
}

func (n *AckNode) Run(ctx context.Context) (Node, error) {
	log := env.Log(ctx)

	input := &AckNodeInputs{}
	err := n.Hydrate(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error hydrating: %s", err)
	}

	log.Info("👾 running ack node", "input", input)

	outMsg, err := env.Dynamo(ctx).GetAiThreadMessage(input.ThreadId, input.OutputMessageId)
	if err != nil {
		return nil, stackerr.Errorf("error getting output message: %s", err)
	}

	messages, err := env.Dynamo(ctx).GetMessagesForThread(input.ThreadId)
	if err != nil {
		return nil, fmt.Errorf("error loading messages: %s", err)
	}

	data := make([]ackNodePromptMessage, 0, len(messages))

	for _, message := range messages {
		if message.LifecycleStage != dynamo.MessageLifecycleStageCompleted {
			continue
		}

		from := "assistant"
		if message.UserID != constants.RevisoUserID {
			from = "user"
		}

		tmplMsg := ackNodePromptMessage{
			From:    from,
			Content: message.FullContent(),
		}

		for _, attachment := range message.Attachments.Attachments {
			switch v := attachment.Value.(type) {
			case *models.Attachment_Document:
				tmplMsg.SelectedContent = &v.Document.Content
			}
		}

		data = append(data, tmplMsg)
	}

	var buf bytes.Buffer
	err = ackTemplate.Execute(&buf, map[string]any{
		"Messages": data,
	})
	if err != nil {
		return nil, fmt.Errorf("error executing template: %s", err)
	}

	content := &models.Content{}
	outMsg.LifecycleReason = "Thinking"
	outMsg.Attachments.Attachments = append(outMsg.Attachments.Attachments,
		&models.Attachment{
			Value: &models.Attachment_Content{
				Content: content,
			},
		},
	)
	err = messaging.UpdateMessage(ctx, outMsg)
	if err != nil {
		return nil, fmt.Errorf("error updating message: %s", err)
	}

	resp, err := n.GenerateStoredPrompt(
		ctx,
		input.ThreadId, "threadV2-acknowledge",
		map[string]string{
			"UserRequest": buf.String(),
		},
		llms.WithStreamingFunc(n.Stream(outMsg, content)))
	if err != nil {
		return nil, fmt.Errorf("error generating response: %s", err)
	}

	outMsg.LifecycleReason = "Thinking"
	err = messaging.UpdateMessage(ctx, outMsg)
	if err != nil {
		return nil, fmt.Errorf("error updating message: %s", err)
	}

	log.Info("[ack] response", "response", resp)
	log.Info("👾 completed ack node", "input", input)

	return nil, nil
}

func (n *AckNode) Stream(outMsg *dynamo.Message, content *models.Content) func(ctx context.Context, chunk []byte) error {
	var buf bytes.Buffer

	return func(ctx context.Context, chunk []byte) error {
		buf.Write(chunk)
		doc, err := utils.ParseIncompleteXML(buf.String())
		if err != nil {
			return fmt.Errorf("error parsing xml: %s", err)
		}

		outMsg.LifecycleReason = ""
		msg := doc.Find("message")
		if msg == nil {
			return nil
		}
		content.Text = msg.Value

		err = messaging.UpdateMessage(ctx, outMsg)
		if err != nil {
			return fmt.Errorf("error updating message: %s", err)
		}
		return nil
	}
}

type AckData struct {
	Message string
}
