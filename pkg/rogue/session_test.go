package rogue_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/rogue"
	"github.com/fivetentaylor/pointy/pkg/testutils"
	v3 "github.com/fivetentaylor/pointy/rogue/v3"
)

func TestSession_Subscribe(t *testing.T) {
	testutils.EnsureStorage()
	ctx := testutils.TestContext()
	docID := uuid.NewString()
	user := testutils.CreateUser(t, ctx)

	testutils.CreateTestDocument(t, ctx, docID, "test")
	testutils.AddOwnerToDocument(t, ctx, docID, user.ID)

	_, ws, msgs, _, cleanup := testutils.CreateTestRogueSessionServer(t, ctx, user, docID)
	defer cleanup()

	sendMessage(t, ws, &rogue.Subscribe{
		Type:  "subscribe",
		DocID: docID,
	})

	receivedMessage, err := pop[string](msgs, 1000*time.Millisecond)
	if err != nil {
		t.Fatalf("Failed to receive message: %v", err)
	}

	var authMsg map[string]interface{}
	err = json.Unmarshal([]byte(receivedMessage), &authMsg)
	if err != nil {
		t.Fatalf("Failed to unmarshal message: %v", err)
	}

	tp, ok := authMsg["type"]
	if !ok {
		t.Fatalf("Failed to get type from message: %v", err)
	}
	assert.Equal(t, "auth", tp)

	receivedMessage, err = pop[string](msgs, 1000*time.Millisecond)
	if err != nil {
		t.Fatalf("Failed to receive message: %v", err)
	}

	var outMsg v3.Message
	err = json.Unmarshal([]byte(receivedMessage), &outMsg)
	if err != nil {
		t.Fatalf("Failed to unmarshal message: %v", err)
	}

	_, ok = outMsg.Op.(v3.SnapshotOp)
	if !ok {
		t.Fatalf("Expected v2.SnapshotOp, got %T", outMsg.Op)
	}
}

func TestSession_Subscribe_with_existing_author(t *testing.T) {
	testutils.EnsureStorage()
	ctx := testutils.TestContext()
	docID := uuid.NewString()
	user := testutils.CreateUser(t, ctx)

	testutils.CreateTestDocument(t, ctx, docID, "test")
	testutils.AddOwnerToDocument(t, ctx, docID, user.ID)

	_, ws, msgs, _, cleanup := testutils.CreateTestRogueSessionServer(t, ctx, user, docID)
	defer cleanup()

	sendMessage(t, ws, &rogue.Subscribe{
		Type:     "subscribe",
		DocID:    docID,
		AuthorID: "0",
	})

	receivedMessage, err := pop(msgs, 1000*time.Millisecond)
	if err != nil {
		t.Fatalf("Failed to receive message: %v", err)
	}

	var authMsg map[string]interface{}
	err = json.Unmarshal([]byte(receivedMessage), &authMsg)
	if err != nil {
		t.Fatalf("Failed to unmarshal message: %v", err)
	}

	tp, ok := authMsg["type"]
	if !ok {
		t.Fatalf("Failed to get type from message: %v", err)
	}
	require.Equal(t, "auth", tp)

	receivedMessage, err = pop(msgs, 1000*time.Millisecond)
	if err != nil {
		t.Fatalf("Failed to receive message: %v", err)
	}

	var outMsg v3.Message
	err = json.Unmarshal([]byte(receivedMessage), &outMsg)
	if err != nil {
		t.Fatalf("Failed to unmarshal message: %v", err)
	}

	_, ok = outMsg.Op.(v3.SnapshotOp)
	if !ok {
		t.Fatalf("Expected v2.SnapshotOp, got %T", outMsg.Op)
	}
}

func TestSession_Subscribe_to_doc_you_dont_own(t *testing.T) {
	testutils.EnsureStorage()
	ctx := testutils.TestContext()
	user := testutils.CreateUser(t, ctx)
	otherDocID := uuid.NewString()

	testutils.CreateTestDocument(t, ctx, otherDocID, "test")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
		if err != nil {
			t.Fatalf("Failed to upgrade to WebSocket: %v", err)
		}

		log := env.Log(ctx)
		conn.SetCloseHandler(func(code int, text string) error {
			log.Info("[testUtils] client disconnected", "code", code, "text", text)
			return nil
		})

		ds := rogue.NewDocStore(env.S3(ctx), env.Query(ctx), env.Redis(ctx))
		_, err = rogue.NewSession(ctx, user, conn, ds, otherDocID)
		require.Error(t, err)
	}))

	u := "ws" + strings.TrimPrefix(server.URL, "http")
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("Failed to dial server: %v", err)
	}

	sendMessage(t, ws, &rogue.Subscribe{
		Type:  "subscribe",
		DocID: otherDocID,
	})

	defer server.Close()
}

func sendMessage(t *testing.T, ws *websocket.Conn, msg interface{}) {
	if err := ws.WriteJSON(msg); err != nil {
		t.Fatalf("Failed to write message: %v", err)
	}
}

func pop[T any](messages <-chan T, timeout time.Duration) (T, error) {
	select {
	case msg := <-messages:
		return msg, nil
	case <-time.After(timeout):
		var zeroValue T
		return zeroValue, errors.New("timeout waiting for message")
	}
}
