package document

import (
	"context"
	"fmt"
	"time"

	"github.com/teamreviso/code/pkg/constants"
	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/models"
	"github.com/teamreviso/code/pkg/query"
	"github.com/teamreviso/code/pkg/service/pubsub"
)

func CreateFolder(ctx context.Context, user *models.User) (*models.Document, error) {
	log := env.SLog(ctx)
	doc := &models.Document{
		Title:    constants.DefaultFolderTitle,
		IsFolder: true,
	}

	err := env.Query(ctx).Transaction(func(tx *query.Query) error {
		documentTbl := tx.Document
		docAccessTbl := tx.DocumentAccess

		err := documentTbl.Omit(documentTbl.RootParentID).Create(doc)
		if err != nil {
			log.Error("error creating db folder", "error", err)
			return err
		}

		err = docAccessTbl.Create(&models.DocumentAccess{
			UserID:         user.ID,
			DocumentID:     doc.ID,
			AccessLevel:    constants.DefaultDocumentRole,
			LastAccessedAt: time.Now(),
		})
		if err != nil {
			log.Error("error creating folder access", "error", err)
			return err
		}

		return nil
	})
	if err != nil {
		log.Error("error creating db folder document", "error", err)
		return nil, fmt.Errorf("sorry, we could not create your folder")
	}

	err = pubsub.PublishNewDocument(ctx, user.ID, doc.ID)
	if err != nil {
		log.Error("error publishing folder document", "error", err)
	}

	log.Info("folder created", "doc_id", doc.ID, "event", "folder_created")
	return doc, nil
}
