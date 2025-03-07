package dag

import (
	"context"
	"fmt"

	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/models"
	v3 "github.com/fivetentaylor/pointy/rogue/v3"
)

type ApplyInput struct {
	DocId    string
	AuthorId string
	Update   string
	Document *v3.Rogue

	EditTarget *EditTarget
}

func ApplyUpdate(
	ctx context.Context,
	input ApplyInput,
) (*models.DocumentRevision, error) {
	log := env.Log(ctx)
	var (
		beforeId, afterId v3.ID
		err               error
	)

	if input.EditTarget == nil {
		beforeId, err = input.Document.GetFirstTotID()
		if err != nil {
			return nil, fmt.Errorf("error getting first id: %s", err)
		}

		afterId, err = input.Document.GetLastTotID()
		if err != nil {
			return nil, fmt.Errorf("error getting last id: %s", err)
		}
	} else {
		beforeId = input.EditTarget.BeforeID
		afterId = input.EditTarget.AfterID
	}

	updatedText := input.Update
	mop := v3.MultiOp{}
	if input.EditTarget != nil {
		log.Info("[dag] applying edit update", "beforeId", beforeId, "afterId", afterId, "action", input.EditTarget.Action)
		if input.EditTarget.Action == EditTargetActionAppend {
			visSize := input.Document.VisSize
			mop, _, err = input.Document.InsertMarkdown(visSize, input.Update)
			if err != nil {
				return nil, fmt.Errorf("error inserting markdown: %s", err)
			}
		} else if input.EditTarget.Action == EditTargetActionPrepend {
			mop, _, err = input.Document.InsertMarkdown(0, input.Update)
			if err != nil {
				return nil, fmt.Errorf("error inserting markdown: %s", err)
			}
		} else {
			mop, _, err = input.Document.ApplyMarkdownDiff(input.AuthorId, updatedText, beforeId, afterId)
			if err != nil {
				return nil, fmt.Errorf("error applying diff: %s", err)
			}
		}
	} else {
		log.Info("[dag] applying full update", "beforeId", beforeId, "afterId", afterId)
		mop, _, err = input.Document.ApplyMarkdownDiff(input.AuthorId, updatedText, beforeId, afterId)
		if err != nil {
			return nil, fmt.Errorf("error applying diff: %s", err)
		}
	}

	err = PublishOp(ctx, input.DocId, mop)
	if err != nil {
		return nil, fmt.Errorf("error publishing op: %s", err)
	}

	revision := &models.DocumentRevision{
		Start:   beforeId.String(),
		End:     afterId.String(),
		Updated: input.Update,
		Id:      "full_document",
	}

	if input.EditTarget != nil {
		revision.Id = input.EditTarget.ID
	}

	return revision, nil
}

func _appendUpdate(input ApplyInput) (beforeID, afterID v3.ID, md string, err error) {
	beforeID, afterID, md, err = input.Document.GetLineMarkdown(input.EditTarget.AfterID)
	if err != nil {
		return beforeID, afterID, "", fmt.Errorf("error getting line markdown: %s", err)
	}

	md = fmt.Sprintf("%s\n%s", md, input.Update)

	return beforeID, afterID, md, nil
}

func _prependUpdate(input ApplyInput) (beforeID, afterID v3.ID, md string, err error) {
	beforeID, afterID, md, err = input.Document.GetLineMarkdown(input.EditTarget.BeforeID)
	if err != nil {
		return beforeID, afterID, "", fmt.Errorf("error getting line markdown: %s", err)
	}

	md = fmt.Sprintf("%s\n%s", input.Update, md)

	return beforeID, afterID, md, nil
}
