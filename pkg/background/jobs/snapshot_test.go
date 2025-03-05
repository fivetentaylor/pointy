package jobs_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/teamreviso/code/pkg/background/jobs"
	"github.com/teamreviso/code/pkg/background/wire"
	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/rogue"
	"github.com/teamreviso/code/pkg/testutils"
)

func TestSnapshotRogue(t *testing.T) {
	testutils.EnsureStorage()
	ctx := testutils.TestContext()

	docID := uuid.New().String()

	user := testutils.CreateUser(t, ctx)
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

	ds := rogue.NewDocStore(env.S3(ctx), env.Query(ctx), env.Redis(ctx))
	count, err := ds.DeltaLogSize(ctx, docID)
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `<html><body>
	<div id="screenshot-page">hello</div>
</body></html>`)
	}))

	// SnapshotRogue uses chrome to hit the app server and screenshot
	// the page.
	// Since we don't have a real app server we need to set the
	// WEB_HOST environment variable to our fake httptest.Server
	fmt.Printf("WEB_HOST: %s\n", svr.URL)
	os.Setenv("WEB_HOST", svr.URL)

	err = jobs.SnapshotRogueJob(ctx, &wire.SnapshotRogue{DocId: docID})
	if err != nil {
		t.Fatalf("Failed to snapshot: %v", err)
	}

	count, err = ds.DeltaLogSize(ctx, docID)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

func purgeMessages[T any](t *testing.T, msgs <-chan T, errs <-chan error) {
	for n, v, err := next[T](msgs, errs); n; n, v, err = next[T](msgs, errs) {
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		fmt.Println("websocket returned:", v)
	}
}

func next[T any](messages <-chan T, errs <-chan error) (bool, T, error) {
	var zeroValue T
	select {
	case msg := <-messages:
		return true, msg, nil
	case err := <-errs:
		return false, zeroValue, err
	case <-time.After(100 * time.Millisecond):
		return false, zeroValue, errors.New("timeout waiting for message")
	}
}
