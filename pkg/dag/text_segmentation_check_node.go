package dag

import (
	"context"
	"fmt"
	"strings"

	"github.com/teamreviso/code/pkg/constants"
)

type TextSegmentationCheckNode struct {
	LargeText Node
	Default   Node

	Base
}

type TextSegmentationCheckNodeInput struct {
	DocId          string `key:"docId"`
	AuthorId       string `key:"authorId"`
	ThreadId       string `key:"threadId"`
	InputMessageId string `key:"inputMessageId"`

	EditTargets []EditTarget `key:"editTargets"`
}

func (n *TextSegmentationCheckNode) Run(ctx context.Context) (Node, error) {
	input := &TextSegmentationCheckNodeInput{}
	err := n.Hydrate(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error hydrating: %s", err)
	}

	var fullText strings.Builder

	if len(input.EditTargets) == 0 {
		// We're editing the entire document
		document, err := GetDocumentAtMessageID(ctx, input.DocId, input.AuthorId, input.ThreadId, input.InputMessageId)
		if err != nil {
			return nil, fmt.Errorf("error getting document: %s", err)
		}

		mkdown, err := document.GetFullMarkdown()
		if err != nil {
			return nil, fmt.Errorf("error getting markdown: %s", err)
		}

		fullText.WriteString(mkdown)
	}

	for _, editTarget := range input.EditTargets {
		fullText.WriteString(editTarget.Markdown)
	}

	tokens := estimateTokens(fullText.String())

	SetStateKey(ctx, "estimatedTokens", tokens)

	if tokens > constants.MaxTokenLength {
		SetStateKey(ctx, "segmentationChoice", "largeText")
		return n.LargeText, nil
	}

	SetStateKey(ctx, "segmentationChoice", "default")
	return n.Default, nil
}
