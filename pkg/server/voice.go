package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/query"
	"github.com/teamreviso/code/pkg/service/voice"
)

const (
	audioFileName = "streamed_audio.raw" // Raw audio file
)

func (s *Server) StreamingVoice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.SLog(ctx)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to upgrade WebSocket: %v\n", err)
		return
	}
	defer conn.Close()

	currentUser, err := env.UserClaim(ctx)
	if err != nil {
		log.Error("error getting current user", "error", err)
		http.Error(w, "error getting current user", http.StatusUnauthorized)
		return
	}

	documentId := chi.URLParam(r, "docID")
	threadId := chi.URLParam(r, "threadID")
	authorId := chi.URLParam(r, "authorID")

	err = ValidateDocumentAndThreadAccess(ctx, documentId, threadId, currentUser.Id)
	if err != nil {
		log.Error("error validating document and thread access", "error", err)
		http.Error(w, "error validating document and thread access", http.StatusUnauthorized)
		return
	}

	log.Info("🔈 Starting Streaming audio...", "documentId", documentId, "threadId", threadId, "authorId", authorId)

	session, err := voice.New(r.Context(), conn, documentId, threadId, currentUser.Id, authorId)
	if err != nil {
		fmt.Printf("Failed to create realtime session: %v\n", err)
		return
	}

	err = session.AddTool(ctx, voice.UpdateDocumentTool)
	if err != nil {
		fmt.Printf("Failed to create realtime session, tool error: %v\n", err)
		return
	}

	go session.Listen(ctx)
	defer session.Close(ctx)

	log.Info("🔊 Waiting for Streaming audio...", "documentId", documentId, "threadId", threadId, "authorId", authorId)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				fmt.Println("WebSocket connection closed.")
				break
			}
			fmt.Printf("Error reading WebSocket message: %v\n", err)
			break
		}

		session.AppendInputAudio(ctx, message)
	}

	log.Info("🔇 Streaming audio complete")
}

func ValidateDocumentAndThreadAccess(ctx context.Context, documentId, threadId, userId string) error {
	_, err := query.GetEditableDocumentForUser(env.Query(ctx), documentId, userId)
	if err != nil {
		return err
	}

	_, err = env.Dynamo(ctx).GetThreadForUser(documentId, threadId, userId)
	if err != nil {
		return err
	}

	return nil
}
