package testutils

import (
	"context"
	"testing"

	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/models"
	"github.com/teamreviso/code/pkg/rogue"
	v3 "github.com/teamreviso/code/rogue/v3"
)

func CreateTestDocument(t *testing.T, ctx context.Context, docID string, content string) (*v3.Rogue, *rogue.DocStore) {
	docModel := &models.Document{
		ID:           docID,
		Title:        "TestDocument",
		RootParentID: docID,
	}

	q := env.Query(ctx)
	documentTbl := q.Document

	err := documentTbl.Create(docModel)
	if err != nil {
		t.Fatalf("CreateTestDocument() failed to create document: error = %v", err)
	}

	rd := v3.NewRogueForQuill("0")

	if content != "" {
		_, err = rd.Insert(0, content)
		if err != nil {
			t.Fatalf("CreateTestDocument() failed rogue.Insert: error = %v", err)
		}
	}

	s3 := env.S3(ctx)
	rc := env.Redis(ctx)
	qry := env.Query(ctx)
	ds := rogue.NewDocStore(s3, qry, rc)

	err = ds.SaveDocToS3(ctx, docID, 0, rd)
	if err != nil {
		t.Fatalf("CreateTestDocument() failed ds.SaveDocToS3: error = %v", err)
	}

	return rd, ds
}

func AddOwnerToDocument(t *testing.T, ctx context.Context, docID string, userID string) {
	q := env.Query(ctx)
	documentTbl := q.Document

	document, err := documentTbl.Where(documentTbl.ID.Eq(docID)).First()
	if err != nil {
		t.Fatalf("AddOwnerToDocument() failed to get document: error = %v", err)
	}

	docAccessTbl := q.DocumentAccess

	accessModel := &models.DocumentAccess{
		UserID:      userID,
		DocumentID:  document.ID,
		AccessLevel: "owner",
	}

	err = docAccessTbl.Create(accessModel)
	if err != nil {
		t.Fatalf("AddOwnerToDocument() failed to create document access: error = %v", err)
	}
}

func GetDocumentHtml(t *testing.T, ctx context.Context, docID string) string {
	log := env.Log(ctx)
	s3 := env.S3(ctx)
	q := env.Query(ctx)
	rds := env.Redis(ctx)

	docStore := rogue.NewDocStore(s3, q, rds)
	_, rog, err := docStore.GetCurrentDoc(ctx, docID)
	if err != nil {
		log.Errorf("error getting document from doc store %s", err)
		t.Fatalf("GetDocumentHtml() failed to get document: = %v", err)
	}

	html, err := rog.GetFullHtml(
		true,
		false,
	)
	if err != nil {
		t.Fatalf("GetDocumentHtml() failed to get html: = %v", err)
	}

	return html
}
