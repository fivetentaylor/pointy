package dag

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/service/attachments"
	"github.com/fivetentaylor/pointy/pkg/service/document"
)

type AttachedDocument struct {
	Identifier string
	Type       string
	Content    string
}

func GetAttachedDocuments(ctx context.Context, threadID string, messageID string) ([]*AttachedDocument, error) {
	inputMessage, err := env.Dynamo(ctx).GetAiThreadMessage(threadID, messageID)
	if err != nil {
		log.Error("error getting input message", "error", err)
		return nil, fmt.Errorf("error getting input message: %s", err)
	}

	attachedRevisoDocumentIds := make([]string, 0)
	attachedFileIds := make([]string, 0)
	for _, msgAttachment := range inputMessage.Attachments.Attachments {
		revisoDoc := msgAttachment.GetRevisoDocument()
		if revisoDoc != nil {
			attachedRevisoDocumentIds = append(attachedRevisoDocumentIds, revisoDoc.Id)
			continue
		}

		file := msgAttachment.GetFile()
		if file != nil {
			attachedFileIds = append(attachedFileIds, file.Id)
			continue
		}
	}

	attachedRevisoDocuments, err := document.GetReadableForUserByIds(
		ctx,
		inputMessage.UserID,
		attachedRevisoDocumentIds,
	)
	if err != nil {
		log.Error("error getting attached reviso documents", "error", err)
		return nil, fmt.Errorf("error getting attached reviso documents: %s", err)
	}

	fileAttachments, err := attachments.GetForUser(ctx, inputMessage.UserID, attachedFileIds)
	if err != nil {
		log.Error("error getting attached files", "error", err)
		return nil, fmt.Errorf("error getting attached files: %s", err)
	}

	out := make([]*AttachedDocument, 0, len(attachedRevisoDocuments)+len(fileAttachments))

	for _, doc := range attachedRevisoDocuments {
		rogueDoc, err := GetDocument(ctx, doc.ID, "reviso")
		if err != nil {
			log.Error("error loading rogue document", "error", err, "doc_id", doc.ID)
			return nil, fmt.Errorf("error loading rogue document: %s", err)
		}

		markdown, err := rogueDoc.GetFullMarkdown()
		if err != nil {
			log.Error("error getting markdown", "error", err)
			return nil, fmt.Errorf("error getting markdown: %s", err)
		}

		out = append(out, &AttachedDocument{
			Identifier: doc.Title,
			Content:    markdown,
			Type:       "reviso",
		})
	}

	for _, file := range fileAttachments {
		extractedText, err := attachments.ExtractedText(ctx, file)
		if err != nil {
			log.Error("error getting extracted text", "error", err)
			return nil, fmt.Errorf("error getting extracted text: %s", err)
		}

		out = append(out, &AttachedDocument{
			Identifier: file.Filename,
			Content:    extractedText,
			Type:       "file",
		})
	}

	return out, nil
}
