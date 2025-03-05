package document

import (
	"context"
	"fmt"
	"strconv"

	"github.com/teamreviso/code/pkg/env"
)

const newAuthorIDSQL = `
INSERT INTO author_ids (author_id, document_id, user_id)
VALUES (increment_author_id(?), ?, ?)
RETURNING author_id
`

func NewAuthorID(ctx context.Context, docID, userID string) (string, error) {
	db := env.RawDB(ctx)

	var newAuthorID int
	if result := db.Raw(newAuthorIDSQL, docID, docID, userID).Scan(&newAuthorID); result.Error != nil {
		return "", result.Error
	}

	return strconv.FormatInt(int64(newAuthorID), 16), nil
}

func ValidateAuthorID(ctx context.Context, authorID, docID, userID string) (bool, error) {
	authorIDInt, err := strconv.ParseInt(authorID, 16, 32)
	if err != nil {
		return false, fmt.Errorf("failed to parse authorID: %w", err)
	}

	// find the author id
	q := env.Query(ctx)
	authorIDTbl := q.AuthorID

	authorIDRow, err := authorIDTbl.Where(
		authorIDTbl.AuthorID.Eq(int32(authorIDInt)),
		authorIDTbl.DocumentID.Eq(docID),
		authorIDTbl.UserID.Eq(userID)).First()
	if err != nil {
		if err.Error() == "record not found" {
			return false, nil
		}
		return false, err
	}

	return authorIDRow != nil, nil
}
