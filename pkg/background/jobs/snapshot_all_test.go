package jobs_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/teamreviso/code/pkg/background/jobs"
	"github.com/teamreviso/code/pkg/background/wire"
	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/models"
	"github.com/teamreviso/code/pkg/rogue"
	"github.com/teamreviso/code/pkg/testutils"
)

func TestSnapshotAllRogue(t *testing.T) {
	testutils.EnsureStorage()
	ctx := testutils.TestContext()

	user := testutils.CreateUser(t, ctx)

	docIDs := []string{}
	for i := 0; i < 5; i++ {
		docID := uuid.New().String()
		fmt.Printf("DOC ID: %s\n", docID)
		docIDs = append(docIDs, docID)
		testutils.CreateTestDocument(t, ctx, docID, "i am a document")
		testutils.AddOwnerToDocument(t, ctx, docID, user.ID)

		_, ws, msgs, errors, cleanup := testutils.CreateTestRogueSessionServer(t, ctx, user, docID)
		defer cleanup()

		authorID := fmt.Sprintf("%s_TS", docID[0:4])

		if err := ws.WriteJSON(&rogue.Subscribe{Type: "subscribe", DocID: docID}); err != nil {
			t.Fatalf("Failed to write message: %v", err)
		}

		purgeMessages(t, msgs, errors)

		msg := &rogue.Operation{
			Type: "op",
			Op:   fmt.Sprintf("[0,[%q,100],\"a\",[%q,0],1]", authorID, "root"),
		}
		if err := ws.WriteJSON(msg); err != nil {
			t.Fatalf("Failed to write message: %v", err)
		}

		purgeMessages(t, msgs, errors)
	}

	ds := rogue.NewDocStore(env.S3(ctx), env.Query(ctx), env.Redis(ctx))

	/*count, err := ds.DeltaLogSize(ctx, docID)
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)*/

	err := jobs.SnapshotAllRogueJob(ctx, &wire.SnapshotAll{Version: "test"})
	if err != nil {
		t.Fatalf("Failed to snapshot all: %v", err)
	}

	for _, docID := range docIDs {
		count, err := ds.DeltaLogSize(ctx, docID)
		require.NoError(t, err)
		require.Equal(t, int64(0), count)
	}

	db := env.RawDB(ctx)
	documents := []models.Document{}
	result := db.Table("documents").Select("*").Find(&documents)
	require.NoError(t, result.Error)

	require.Equal(t, len(docIDs), len(documents))
	for _, doc := range documents {
		require.Equal(t, "test", doc.RogueVersion)
	}
}
