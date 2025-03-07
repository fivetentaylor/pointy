package document

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/models"
	"github.com/fivetentaylor/pointy/pkg/query"
	"github.com/fivetentaylor/pointy/pkg/stackerr"
)

func GetReadableForUserByIds(ctx context.Context, userID string, ids []string) ([]*models.Document, error) {
	docs, err := query.GetAccessibleDocumentsForUser(env.Query(ctx), ids, userID)
	if err != nil {
		log.Errorf("error getting document: %s", stackerr.Wrap(err))
		return nil, fmt.Errorf("sorry, we could not load your document")
	}

	return docs, nil
}

func GetBranchCopies(ctx context.Context, userID, documentID string) ([]*models.Document, error) {
	var branchCopies []*models.Document
	err := env.Query(ctx).Transaction(func(tx *query.Query) error {
		documentTbl := tx.Document

		count, err := tx.DocumentAccess.
			Where(tx.DocumentAccess.DocumentID.Eq(documentID)).
			Where(tx.DocumentAccess.UserID.Eq(userID)).
			Count()
		if err != nil {
			return fmt.Errorf("error counting document access: %w", err)
		}

		if count == 0 {
			return stackerr.New(fmt.Errorf("user: %v has no access to document: %v", userID, documentID))
		}

		doc, err := documentTbl.
			Where(documentTbl.ID.Eq(documentID)).
			Where(documentTbl.DeletedAt.IsNull()).
			First()

		if err != nil {
			return fmt.Errorf("error querying branch copies: %w", err)
		}

		err = documentTbl.
			Join(tx.DocumentAccess, tx.DocumentAccess.DocumentID.EqCol(documentTbl.ID)).
			Where(documentTbl.RootParentID.Eq(doc.RootParentID)).
			Where(documentTbl.DeletedAt.IsNull()).
			Where(tx.DocumentAccess.UserID.Eq(userID)).
			Order(documentTbl.CreatedAt).
			Scan(&branchCopies)

		if err != nil {
			return fmt.Errorf("error scanning branch copies: %w", err)
		}

		return nil
	})

	if err != nil {
		log.Errorf("error retrieving branch copies for document %s: %v", documentID, err)
		return nil, fmt.Errorf("sorry, we could not retrieve the branch copies")
	}

	out := []*models.Document{}
	childTree := map[string][]*models.Document{}
	for _, doc := range branchCopies {
		if doc.ParentID == nil {
			continue
		}

		if _, ok := childTree[*doc.ParentID]; ok {
			childTree[*doc.ParentID] = append(childTree[*doc.ParentID], doc)
		} else {
			childTree[*doc.ParentID] = []*models.Document{doc}
		}
	}

	toVisit := []string{documentID}
	for len(toVisit) > 0 {
		parentID := toVisit[0]
		toVisit = toVisit[1:]

		if children, ok := childTree[parentID]; ok {
			for _, child := range children {
				out = append(out, child)
				toVisit = append(toVisit, child.ID)
			}
		}
	}

	return out, nil
}

func GetRootDocuments(ctx context.Context, folderId *string, userID string, limit, offset int) ([]*models.Document, error) {
	baseQuery := `
    WITH user_docs AS (
        SELECT d.*
        FROM documents d
        JOIN document_access da ON d.id = da.document_id
        WHERE da.user_id = ?
		AND da.access_level = 'owner'
		%s
    ), root_docs AS (
        SELECT ud0.*
        FROM user_docs ud0
        LEFT JOIN user_docs ud1 ON ud0.parent_id = ud1.id
        WHERE ud1.id IS NULL
    )
    SELECT r.*
    FROM root_docs r
    ORDER BY r.updated_at DESC
    LIMIT ? OFFSET ?
  `

	var (
		documents   []*models.Document
		queryParams []interface{}
	)

	folderFilter := " AND d.folder_id IS NULL"

	queryParams = append(queryParams, userID)
	if folderId != nil {
		folderFilter = " AND d.folder_id = ?"
		queryParams = append(queryParams, *folderId)
	}

	finalQuery := fmt.Sprintf(baseQuery, folderFilter)
	queryParams = append(queryParams, limit, offset)

	db := env.RawDB(ctx)
	err := db.Raw(finalQuery, queryParams...).Scan(&documents).Error
	if err != nil {
		return nil, err
	}

	return documents, nil
}

func GetAllDocuments(ctx context.Context, userID string, limit, offset int) ([]*models.Document, error) {
	baseQuery := `
    WITH user_docs AS (
        SELECT d.*
        FROM documents d
        JOIN document_access da ON d.id = da.document_id
        WHERE da.user_id = ?
		AND d.is_folder = false
    ), root_docs AS (
        SELECT ud0.*
        FROM user_docs ud0
        LEFT JOIN user_docs ud1 ON ud0.parent_id = ud1.id
        WHERE ud1.id IS NULL
    )
    SELECT r.*
    FROM root_docs r
    ORDER BY r.updated_at DESC
    LIMIT ? OFFSET ?
  `

	var (
		documents   []*models.Document
		queryParams []interface{}
	)

	queryParams = append(queryParams, userID)
	queryParams = append(queryParams, limit, offset)

	db := env.RawDB(ctx)
	err := db.Raw(baseQuery, queryParams...).Scan(&documents).Error
	if err != nil {
		return nil, err
	}

	return documents, nil
}

func GetSharedDocuments(ctx context.Context, userID string, limit, offset int) ([]*models.Document, error) {
	baseQuery := `
    WITH user_docs AS (
        SELECT d.*
        FROM documents d
        JOIN document_access da ON d.id = da.document_id
        WHERE da.user_id = ?
		AND d.is_folder = false
		AND da.access_level != 'owner'
    ), root_docs AS (
        SELECT ud0.*
        FROM user_docs ud0
        LEFT JOIN user_docs ud1 ON ud0.parent_id = ud1.id
        WHERE ud1.id IS NULL
    )
    SELECT r.*
    FROM root_docs r
    ORDER BY r.updated_at DESC
    LIMIT ? OFFSET ?
  `

	var (
		documents   []*models.Document
		queryParams []interface{}
	)

	queryParams = append(queryParams, userID)
	queryParams = append(queryParams, limit, offset)

	db := env.RawDB(ctx)
	err := db.Raw(baseQuery, queryParams...).Scan(&documents).Error
	if err != nil {
		return nil, err
	}

	return documents, nil
}

func SearchRootDocuments(ctx context.Context, query string, userID string, limit, offset int) ([]*models.Document, error) {
	baseQuery := `
    WITH user_docs AS (
        SELECT d.*
        FROM documents d
        JOIN document_access da ON d.id = da.document_id
        WHERE da.user_id = ?
		AND d.title ILIKE ?
		AND d.is_folder = false
    ), root_docs AS (
        SELECT ud0.*
        FROM user_docs ud0
        LEFT JOIN user_docs ud1 ON ud0.parent_id = ud1.id
        WHERE ud1.id IS NULL
    )
    SELECT r.*
    FROM root_docs r
    ORDER BY r.updated_at DESC
    LIMIT ? OFFSET ?
  `

	var (
		documents   []*models.Document
		queryParams []interface{}
	)

	queryParams = append(queryParams, userID)
	queryParams = append(queryParams, "%"+query+"%")
	queryParams = append(queryParams, limit, offset)

	db := env.RawDB(ctx)
	err := db.Raw(baseQuery, queryParams...).Scan(&documents).Error
	if err != nil {
		return nil, err
	}

	return documents, nil
}
