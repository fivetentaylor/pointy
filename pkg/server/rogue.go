package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/models"
	"github.com/fivetentaylor/pointy/pkg/rogue"
)

func (s *Server) RogueWebSocket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.SLog(ctx)
	query := env.Query(ctx)

	currentUser, err := env.UserClaim(ctx)
	if err != nil {
		log.Error("error getting current user", "error", err)
		http.Error(w, "error getting current user", http.StatusUnauthorized)
		return
	}

	docID := chi.URLParam(r, "docID")
	documentTbl := query.Document

	_, err = documentTbl.
		Where(documentTbl.ID.Eq(docID)).
		First()
	if err != nil {
		log.Error("error getting document", "error", err)
		http.NotFound(w, r)
		return
	}

	docAccessTbl := query.DocumentAccess
	var user *models.User

	if currentUser == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	count, err := docAccessTbl.
		Where(docAccessTbl.DocumentID.Eq(docID)).
		Where(docAccessTbl.UserID.Eq(currentUser.Id)).
		Count()
	if err != nil {
		log.Error("error checking document access", "error", err)
		// Decide on the appropriate response, e.g., internal server error
		http.Error(w, "error checking access", http.StatusInternalServerError)
		return
	}
	hasWriteAccess := count > 0
	if !hasWriteAccess {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userTbl := query.User
	user, err = userTbl.Where(userTbl.ID.Eq(currentUser.Id)).First()
	if err != nil {
		log.Error("error querying user", "error", err)
		// Decide on the appropriate response, e.g., internal server error
		http.Error(w, "error querying user", http.StatusInternalServerError)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading to websocket:", err)
		return
	}
	defer conn.Close()

	docStore := rogue.NewDocStore(s.S3, s.Query, s.Redis)
	session, err := rogue.NewSession(ctx, user, conn, docStore, docID)
	if err != nil {
		log.Error("error creating session", "error", err)
		http.Error(w, "error creating session", http.StatusInternalServerError)
		return
	}

	defer func() {
		log.Info("🛑 editor disconnected, enqueuing snapshot", "docID", docID, "event", "editor_disconnected")
		session.Close(ctx)
	}()

	log.Info("📡 editor connected", "docID", docID, "event", "editor_connected")
	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			return
		}

		if messageType == websocket.TextMessage {
			err = session.Message(ctx, msg)
			if err != nil {
				log.Error("error handling message", "err", err)
				return
			}

			continue
		}

		log.Warn("unexpected message type", "messageType", messageType, "message", msg)
	}
}
