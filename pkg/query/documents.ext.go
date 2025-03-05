package query

import (
	"log/slog"
	"math/rand"

	"github.com/charmbracelet/log"

	"github.com/teamreviso/code/pkg/constants"
	"github.com/teamreviso/code/pkg/models"
	"github.com/teamreviso/code/pkg/stackerr"
	"github.com/teamreviso/code/pkg/utils"
)

type AccessDeniedError struct {
	Message string
}

type DocumentExtended struct {
	models.Document
	OwnerID string `json:"owner_id" gorm:"-"`
}

func (e *AccessDeniedError) Error() string {
	return e.Message
}

func GetRandomDocument(db *Query) (*models.Document, error) {
	documentTbl := db.Document

	totalCount, err := documentTbl.Count()
	if err != nil {
		return nil, err
	}

	if totalCount == 0 {
		return nil, nil
	}

	doc, err := documentTbl.Offset(rand.Intn(int(totalCount))).First()
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func GetOwnerOfDocument(q *Query, docId string) (*models.User, error) {
	userTbl := q.User
	docAccessTbl := q.DocumentAccess

	user, err := userTbl.
		LeftJoin(docAccessTbl, userTbl.ID.EqCol(docAccessTbl.UserID)).
		Where(docAccessTbl.DocumentID.Eq(docId)).
		Where(docAccessTbl.AccessLevel.Eq("owner")).
		First()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetAccessibleDocumentsForUser(q *Query, documentIDs []string, userID string) ([]*models.Document, error) {
	slog.Info("getting documents for user", "ids", documentIDs, "userID", userID)
	documentTbl := q.Document
	docAccessTbl := q.DocumentAccess

	docs, err := documentTbl.
		LeftJoin(docAccessTbl, documentTbl.ID.EqCol(docAccessTbl.DocumentID)).
		Where(documentTbl.ID.In(documentIDs...)).
		Where(docAccessTbl.UserID.Eq(userID)).
		Find()
	if err != nil {
		log.Error("error fetching documents", "ids", documentIDs, "userID", userID, "err", err)
		return nil, stackerr.Wrap(err)
	}

	return docs, nil
}

func GetOwnedDocumentForUser(q *Query, documentID, userID string) (*models.Document, error) {
	slog.Info("getting owned document for user", "id", documentID, "userID", userID)
	documentTbl := q.Document
	docAccessTbl := q.DocumentAccess

	doc, err := documentTbl.
		LeftJoin(docAccessTbl, documentTbl.ID.EqCol(docAccessTbl.DocumentID)).
		Where(documentTbl.ID.Eq(documentID)).
		Where(docAccessTbl.AccessLevel.Eq("owner")).
		Where(docAccessTbl.UserID.Eq(userID)).
		First()

	if err != nil {
		return doc, err
	}

	return doc, err
}

func GetReadableDocumentForUser(q *Query, documentID, userID string) (*models.Document, error) {
	slog.Info("getting readable document for user", "id", documentID, "userID", userID)
	documentTbl := q.Document
	docAccessTbl := q.DocumentAccess

	if userID == constants.RevisoUserID {
		doc, err := documentTbl.
			LeftJoin(docAccessTbl, documentTbl.ID.EqCol(docAccessTbl.DocumentID)).
			Where(documentTbl.ID.Eq(documentID)).
			First()
		return doc, err
	}

	// Look for the document first
	doc, err := documentTbl.
		Where(documentTbl.ID.Eq(documentID)).
		First()
	// If there was an error fetching the document, return now
	if err != nil {
		log.Error("error fetching document", "id", documentID, "userID", userID, "err", err)
		return doc, stackerr.Wrap(err)
	}

	if doc.IsPublic {
		return doc, nil
	}

	// If no user ID is provided, that means the document needs to be public to be accessed
	if userID == "" {
		log.Error("unauthorized access to non-public document (anonymous access)", "id", documentID, "userID", userID)
		return nil, &AccessDeniedError{Message: "access denied: insufficient permissions"}
	}

	// If a user ID is provided, check their access level
	_, err = docAccessTbl.
		Where(docAccessTbl.UserID.Eq(userID)).
		Where(docAccessTbl.DocumentID.Eq(documentID)).
		First()
	if err != nil {
		log.Error("unauthorized access to non-public document (user does not have access)", "id", documentID, "userID", userID)
		return nil, &AccessDeniedError{Message: "access denied: insufficient permissions"}
	}

	return doc, nil
}

func GetEditableDocumentForUser(q *Query, id, userID string) (*models.Document, error) {
	slog.Info("getting editable document for user", "id", id, "userID", userID)
	documentTbl := q.Document
	docAccessTbl := q.DocumentAccess

	// Look for the document first
	doc, err := documentTbl.
		LeftJoin(docAccessTbl, documentTbl.ID.EqCol(docAccessTbl.DocumentID)).
		Where(documentTbl.ID.Eq(id)).
		First()
	if err != nil {
		return doc, err
	}

	access, err := docAccessTbl.
		Where(docAccessTbl.UserID.Eq(userID)).
		Where(docAccessTbl.DocumentID.Eq(id)).
		First()

	if err != nil || !utils.Contains(constants.AccessLevelsWithEdit, access.AccessLevel) {
		return nil, &AccessDeniedError{Message: "access denied: insufficient permissions"}
	}

	return doc, err
}
