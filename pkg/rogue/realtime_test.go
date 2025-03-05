package rogue_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/rogue"
	"github.com/teamreviso/code/pkg/testutils"
)

func TestRealtime_Subscribe(t *testing.T) {
	t.Parallel()
	testutils.EnsureStorage()
	cancelWorker := testutils.RunWorker(t)
	defer cancelWorker()
	ctx := testutils.TestContext()
	userID := uuid.NewString()
	authorID := userID[0:4] + "_TS"
	docID := uuid.NewString()
	testutils.CreateTestDocument(t, ctx, docID, "realtime subscribe test")

	rds := env.Redis(ctx)
	qry := env.Query(ctx)

	rt := rogue.NewRealtime(rds, qry, docID, userID, "Bob", "pink")

	cancel := rt.Subscribe(ctx, authorID, func(msg []byte) {
		t.Log(string(msg))
	})

	err := testutils.RetryOnError(5, 100*time.Millisecond, func() error {
		_, err := rogue.GetAuthorLastCursor(ctx, rds, docID, userID, authorID)
		return err
	})
	assert.NoError(t, err)

	connections, err := rogue.CurrentActiveConnections(ctx, rds, docID)
	assert.NoError(t, err)

	found := false
	for _, connection := range connections {
		if connection == fmt.Sprintf("%s:%s", userID, authorID) {
			found = true
			break
		}
	}
	assert.True(t, found)

	cancel()

	// User should be disconnected now
	r := testutils.RetryUntilTrue(5, 100*time.Millisecond, func() bool {
		_, err := rogue.GetAuthorLastCursor(ctx, rds, docID, userID, authorID)
		return err != nil
	})
	assert.True(t, r)

	connections, err = rogue.CurrentActiveConnections(ctx, rds, docID)
	assert.NoError(t, err)

	found = false
	for _, connection := range connections {
		if connection == fmt.Sprintf("%s:%s", userID, authorID) {
			found = true
			break
		}
	}
	assert.False(t, found)

}

func TestRealtime_Subscribe_wait_for_presence(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping long-running wait for presence in short mode.")
	}

	t.Parallel()
	testutils.EnsureStorage()
	cancelWorker := testutils.RunWorker(t)
	defer cancelWorker()
	ctx := testutils.TestContext()
	userID := uuid.NewString()
	authorID := userID[0:4] + "_TS"
	docID := uuid.NewString()
	testutils.CreateTestDocument(t, ctx, docID, "realtime subscribe test")

	rds := env.Redis(ctx)
	qry := env.Query(ctx)

	rt := rogue.NewRealtime(rds, qry, docID, userID, "Bob", "pink")

	cancel := rt.Subscribe(ctx, authorID, func(msg []byte) {
		t.Log(string(msg))
	})

	// Wait for presence check
	time.Sleep(rogue.PresenceCheckInterval + 1*time.Second)

	err := testutils.RetryOnError(5, 100*time.Millisecond, func() error {
		_, err := rogue.GetAuthorLastCursor(ctx, rds, docID, userID, authorID)
		return err
	})
	assert.NoError(t, err)

	cancel()
}
