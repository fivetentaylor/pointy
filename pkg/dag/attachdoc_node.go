package dag

import (
	"context"
	"fmt"

	"github.com/fivetentaylor/pointy/pkg/models"
)

const fullDocId = "full_doc"
const attachedDocsKey = "attached_docs"

type AttachDocNode struct {
	Next Node
	Base
}

type AttachDocInput struct {
	DocId    string `key:"docId"`
	AuthorId string `key:"authorId"`
}

func (n *AttachDocNode) Run(ctx context.Context) (Node, error) {
	input := &AttachDocInput{}
	err := n.Hydrate(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error hydrating: %s", err)
	}

	document, err := GetDocument(ctx, input.DocId, input.AuthorId)
	if err != nil {
		return nil, fmt.Errorf("error getting document: %s", err)
	}

	attachedDocuments, err := GetStateKey[map[string]*models.Attachment_Document](ctx, attachedDocsKey)
	if err != nil {
		return nil, err
	}
	if attachedDocuments == nil {
		attachedDocuments = map[string]*models.Attachment_Document{}
	}

	startID, err := document.Rope.GetTotID(0)
	if err != nil {
		return nil, fmt.Errorf("error getting start id: %s", err)
	}

	endID, err := document.Rope.GetTotID(document.TotSize - 1)
	if err != nil {
		return nil, fmt.Errorf("error getting end id: %s", err)
	}

	markdown, err := document.GetMarkdown(startID, endID)
	if err != nil {
		return nil, fmt.Errorf("error getting markdown: %s", err)
	}

	attachedDocuments[fullDocId] = &models.Attachment_Document{
		Document: &models.DocumentSelection{
			Id:      fullDocId,
			Content: markdown,
			Start:   startID.String(),
			End:     endID.String(),
		},
	}

	SetStateKey(ctx, attachedDocsKey, attachedDocuments)

	return n.Next, nil
}
