package models

import (
	"encoding/json"
	"fmt"

	"github.com/hoisie/mustache"
	"github.com/teamreviso/freeplay"
)

func (p *Prompt) Content() ([]*freeplay.Message, error) {
	msgs := []*freeplay.Message{}
	err := json.Unmarshal([]byte(p.ContentJSON), &msgs)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal content: %v", err)
	}
	return msgs, nil
}

func (p *Prompt) RenderContent(data any) ([]*freeplay.Message, error) {
	msgs, err := p.Content()
	if err != nil {
		return nil, err
	}

	outMsgs := make([]*freeplay.Message, 0, len(msgs))
	for _, msg := range msgs {
		tmpl, err := mustache.ParseString(msg.Content)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template %q: %v", msg.Content, err)
		}
		formattedContent := tmpl.Render(msg.Content, data)

		outMsgs = append(outMsgs, &freeplay.Message{
			Content: formattedContent,
			Role:    msg.Role,
		})
	}

	return outMsgs, nil
}
