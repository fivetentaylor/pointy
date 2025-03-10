package query

import (
	"fmt"

	"github.com/fivetentaylor/pointy/pkg/constants"
)

func AccessLevelForDocument(q *Query, docId, userID string) (string, error) {
	docAccessTbl := q.DocumentAccess

	access, err := docAccessTbl.
		Where(docAccessTbl.UserID.Eq(userID)).
		Where(docAccessTbl.DocumentID.Eq(docId)).
		First()

	if err != nil {
		return constants.AccessLevelNone, fmt.Errorf("error getting document access: %s", err)
	}

	return access.AccessLevel, err
}
