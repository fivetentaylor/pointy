package loaders

import (
	"context"

	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/models"
)

type ctxKey string

const (
	LoadersKey = ctxKey("dataloaders")
)

type userReader struct{}

var MissingUser = &models.User{
	ID:          "missing",
	Name:        "missing",
	DisplayName: "missing",
	Email:       "missing@revi.so",
	Admin:       false,
}

// getUsers implements a batch function that can retrieve many users by ID,
// for use in a dataloader
func (u *userReader) getUsers(ctx context.Context, userIDs []string) ([]*models.User, []error) {
	usrTbl := env.Query(ctx).User
	users, err := usrTbl.Where(usrTbl.ID.In(userIDs...)).Find()
	if err != nil {
		return nil, []error{err}
	}

	out := make([]*models.User, len(userIDs))
	for i, id := range userIDs {
		for _, user := range users {
			if user.ID == id {
				out[i] = user
				break
			}
		}
	}

	return out, nil
}

// GetUser returns single user by id efficiently
func GetUser(ctx context.Context, userID string) (*models.User, error) {
	//env.Log(ctx).Infof("GetUser: %s", userID)
	loaders := For(ctx)
	return loaders.UserLoader.Load(ctx, userID)
}

// GetUsers returns many users by ids efficiently
func GetUsers(ctx context.Context, userIDs []string) ([]*models.User, error) {
	//env.Log(ctx).Infof("GetUsers: %s", userIDs)
	loaders := For(ctx)
	return loaders.UserLoader.LoadAll(ctx, userIDs)
}
