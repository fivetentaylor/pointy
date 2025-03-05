package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/query"
	"github.com/teamreviso/code/pkg/rogue"
	"github.com/teamreviso/code/pkg/views/editor"
)

func (s *Server) HtmlDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(ctx)
	q := env.Query(ctx)

	docID := chi.URLParam(r, "docID")

	var userID string
	currentUser, err := env.UserClaim(ctx)
	if err == nil {
		userID = currentUser.Id
	}

	doc, err := query.GetReadableDocumentForUser(q, docID, userID)
	if doc == nil || err != nil {
		if err != nil {
			switch e := err.(type) {
			case *query.AccessDeniedError:
				log.Errorf("error getting document (access denied): %s", e.Error())
				http.Error(w, err.Error(), http.StatusForbidden)
			default:
				log.Errorf("error getting document: %s", e.Error())
			}
		}
		http.NotFound(w, r)
		return
	}

	docStore := rogue.NewDocStore(s.S3, s.Query, s.Redis)
	_, rog, err := docStore.GetCurrentDoc(ctx, docID)
	if err != nil {
		log.Errorf("error getting document from doc store: %s", err)
		http.NotFound(w, r)
		return
	}

	html, err := rog.GetFullHtml(
		true,
		false,
	)
	if err != nil {
		log.Errorf("error getting html: %s", err)
		http.NotFound(w, r)
		return
	}

	w.Header().Set("X-Document-Title", doc.Title)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func (s *Server) DocumentEditor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(ctx)
	q := env.Query(ctx)

	docID := chi.URLParam(r, "docID")

	var userID string
	currentUser, err := env.UserClaim(ctx)
	if err == nil {
		userID = currentUser.Id
	}

	doc, err := query.GetEditableDocumentForUser(q, docID, userID)
	if doc == nil || err != nil {
		if err != nil {
			switch e := err.(type) {
			case *query.AccessDeniedError:
				log.Errorf("error getting document: %s", e.Error())
				http.Error(w, err.Error(), http.StatusForbidden)
			default:
				log.Errorf("error getting document: %s", e.Error())
			}
		}
		http.NotFound(w, r)
		return
	}

	docStore := rogue.NewDocStore(s.S3, s.Query, s.Redis)
	_, rog, err := docStore.GetCurrentDoc(ctx, docID)
	if err != nil {
		log.Errorf("error getting document from doc store: %s", err)
		http.NotFound(w, r)
		return
	}

	editor.Editor(docID, rog).Render(ctx, w)
}
