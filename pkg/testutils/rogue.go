package testutils

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/models"
	"github.com/fivetentaylor/pointy/pkg/rogue"
)

func CreateTestRogueSessionServer(t *testing.T, ctx context.Context, user *models.User, docID string) (*httptest.Server, *websocket.Conn, chan string, chan error, func()) {
	messages := make(chan string, 1)
	errors := make(chan error, 1)
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

		session, err := rogue.NewSession(ctx, user, conn, ds, docID)
		require.NoError(t, err)
		session.DeactivateDocLogger() // don't log to S3 in tests

		// Start a goroutine to read messages from the session's conn
		go func() {
			log.Info("[testUtils] test server connected", "docID", docID)
			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					log.Error(fmt.Errorf("[testUtils] error reading message: %w", err))
					break
				}
				log.Info("[testUtils] message received", "docID", docID, "message", string(message))
				err = session.Message(ctx, message)
				if err != nil {
					log.Error(fmt.Errorf("[testUtils] error processing message: %w", err))
					errors <- err
					break
				}
			}

			log.Info("test server disconnected", "docID", docID)
		}()
	}))

	u := "ws" + strings.TrimPrefix(server.URL, "http")
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("Failed to dial server: %v", err)
	}

	go func() {
		for {
			_, msg, err := ws.ReadMessage()
			if err != nil {
				log := env.Log(ctx)
				log.Error(fmt.Errorf("error reading message: %w", err))
				return
			}
			messages <- string(msg)
		}
	}()

	cleanup := func() {
		ws.Close()
		server.Close()
		close(messages)
	}

	return server, ws, messages, errors, cleanup
}
