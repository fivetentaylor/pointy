package testutils

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/models"
	"github.com/fivetentaylor/pointy/pkg/utils"
)

func CreateUser(t *testing.T, ctx context.Context, options ...func(*models.User)) *models.User {
	user := &models.User{
		ID: uuid.NewString(),
		Email: fmt.Sprintf(
			"%s@example.com",
			utils.RandomSafeString(10),
		),
		Provider: "google",
	}

	for _, opt := range options {
		opt(user)
	}

	q := env.Query(ctx)
	userTbl := q.User
	err := userTbl.Create(user)
	if err != nil {
		t.Fatalf("CreateUser() failed to create user: error = %v", err)
	}

	return user
}

func CreateAdmin(t *testing.T, ctx context.Context, options ...func(*models.User)) *models.User {
	user := &models.User{
		ID: uuid.NewString(),
		Email: fmt.Sprintf(
			"%s@example.com",
			utils.RandomSafeString(10),
		),
		Provider: "google",
		Admin:    true,
	}

	for _, opt := range options {
		opt(user)
	}

	q := env.Query(ctx)
	userTbl := q.User
	err := userTbl.Create(user)
	if err != nil {
		t.Fatalf("CreateUser() failed to create user: error = %v", err)
	}

	return user
}
