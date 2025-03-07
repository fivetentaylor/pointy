package dag

import (
	"context"
	"fmt"

	"github.com/fivetentaylor/pointy/pkg/models"
)

func ApplySuggestion(
	ctx context.Context,
	docId string,
	authorId string,
	attachment *models.Attachment_Suggestion,
) error {
	suggestion := attachment.Suggestion
	document, err := GetDocument(ctx, docId, authorId)
	if err != nil {
		return fmt.Errorf("error getting document: %s", err)
	}

	op, err := document.Insert(document.VisSize-1, "\n\n")
	if err != nil {
		return fmt.Errorf("error inserting new line: %s", err)
	}

	err = PublishOp(ctx, docId, op)
	if err != nil {
		return fmt.Errorf("error publishing op: %s", err)
	}

	startID := op.ID
	endID, err := document.VisRightOf(startID)
	if err != nil {
		return fmt.Errorf("error getting end id: %s", err)
	}

	mop, _, err := document.ApplyMarkdownDiff(authorId, suggestion.Content, startID, endID)
	if err != nil {
		return fmt.Errorf("error applying diff: %s", err)
	}

	err = PublishOp(ctx, docId, mop)
	if err != nil {
		return fmt.Errorf("error publishing op: %s", err)
	}

	return nil
}
