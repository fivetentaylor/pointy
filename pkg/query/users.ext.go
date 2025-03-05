package query

import (
	"math/rand"

	"github.com/teamreviso/code/pkg/models"
)

func GetRandomUser(db *Query) (*models.User, error) {
	tbl := db.User

	totalCount, err := tbl.Count()
	if err != nil {
		return nil, err
	}

	if totalCount == 0 {
		return nil, nil
	}

	doc, err := tbl.Offset(rand.Intn(int(totalCount))).First()
	if err != nil {
		return nil, err
	}

	return doc, nil
}
