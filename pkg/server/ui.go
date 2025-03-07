package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/fivetentaylor/pointy/pkg/background/wire"
	"github.com/fivetentaylor/pointy/pkg/constants"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/models"
	"github.com/fivetentaylor/pointy/pkg/query"
	"github.com/fivetentaylor/pointy/pkg/rogue"
	"github.com/fivetentaylor/pointy/pkg/service/document"
	"github.com/fivetentaylor/pointy/pkg/views"
	"github.com/fivetentaylor/pointy/pkg/views/drafts"
	"github.com/fivetentaylor/pointy/pkg/views/read"
	"github.com/fivetentaylor/pointy/pkg/views/utils"
)

func (s *Server) UI(r chi.Router) {
	r.Get("/not-found", s.NotFound)
	r.Get("/error", s.Error)
	r.Get("/logout", s.Logout)
	r.Get("/drafts", s.Drafts)
	r.Get("/drafts/", s.Drafts)
	r.Get("/drafts/{docID}", s.Draft)
	r.Get("/drafts/{docID}/{threadID}", s.Draft)
	r.Get("/drafts/{docID}/images/{imageID}", s.GetDocumentImage)
	r.Get("/read/{docID}", s.Read)
}

func (s *Server) _createDefaultIntroDoc(ctx context.Context, userID string) (*models.Document, error) {
	db := env.RawDB(ctx)

	var introDocID string
	result := db.Raw("SELECT doc_id FROM default_documents WHERE name = 'intro_doc'").Scan(&introDocID)

	if result.Error != nil {
		return nil, result.Error
	}

	if introDocID != "" {
		s3 := env.S3(ctx)
		redis := env.Redis(ctx)
		query := env.Query(ctx)
		ds := rogue.NewDocStore(s3, query, redis)

		copyDoc := &models.Document{
			IsPublic: true,
			Title:    "Welcome to Reviso",
		}

		dupDoc, err := document.CreateCustom(ctx, userID, copyDoc)
		if err != nil {
			return nil, fmt.Errorf("sorry, we could not duplicate your document")
		}

		err = ds.DuplicateDoc(ctx, introDocID, dupDoc.ID, nil)
		if err != nil {
			return nil, err
		}
	}

	doc, err := document.Create(ctx, userID)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (s *Server) Drafts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(ctx)

	currentUser, err := env.UserClaim(ctx)
	if err != nil {
		log.Error(fmt.Sprintf("Error getting user claim: %s", err.Error()))
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)

		return
	}

	// Get the last document that the user updated
	documentTbl := env.Query(ctx).Document
	docAccessTbl := env.Query(ctx).DocumentAccess

	usersDocuments, err := documentTbl.
		LeftJoin(docAccessTbl, documentTbl.ID.EqCol(docAccessTbl.DocumentID)).
		Where(docAccessTbl.UserID.Eq(currentUser.Id)).
		Order(documentTbl.UpdatedAt.Desc()).Limit(1).Find()
	if err != nil {
		log.Errorf("error getting documents: %s", err)
		http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
		return
	}

	if len(usersDocuments) == 0 {
		/*userTbl := env.Query(ctx).User
		user, err := userTbl.Where(userTbl.ID.Eq(currentUser.Id)).First()
		if err != nil {
			log.Errorf("error getting user: %s", err)
			http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
			return
		}

		doc, err := document.Create(ctx, user)
		if err != nil {
			log.Errorf("error creating document: %s", err)
			http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
			return
		}*/
		doc, err := s._createDefaultIntroDoc(ctx, currentUser.Id)
		if err != nil {
			log.Errorf("error creating default intro doc: %s", err)
			http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/drafts/%s", doc.ID), http.StatusTemporaryRedirect)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/drafts/%s", usersDocuments[0].ID), http.StatusTemporaryRedirect)
}

func (s *Server) Draft(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(ctx)
	q := env.Query(ctx)

	docID := chi.URLParam(r, "docID")

	currentUser, err := env.UserClaim(ctx)
	if err != nil {
		log.Error(fmt.Sprintf("Error getting user claim: %s", err.Error()))
		http.Redirect(w, r, fmt.Sprintf("/read/%s", docID), http.StatusTemporaryRedirect)

		return
	}

	log.Infof("Loading docID: %q for user: %s", docID, currentUser.Email)

	doc, err := query.GetEditableDocumentForUser(q, docID, currentUser.Id)
	if doc == nil || err != nil {
		if err != nil {
			log.Errorf("error getting document: %s", err)
		}
		http.Redirect(w, r, fmt.Sprintf("/read/%s", docID), http.StatusTemporaryRedirect)
		return
	}

	_, err = env.Background(ctx).Enqueue(ctx, &wire.AccessDoc{
		UserId:       currentUser.Id,
		DocId:        doc.ID,
		TimestampStr: time.Now().Format(time.RFC3339),
	})
	if err != nil {
		log.Errorf("error enqueuing job: %s", err)
		http.Redirect(w, r, fmt.Sprintf("/read/%s", docID), http.StatusTemporaryRedirect)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	drafts.Draft(doc.Title).Render(context.Background(), w)
}

func (s *Server) NotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.NotFound().Render(context.Background(), w)
}

func (s *Server) Error(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	utils.Error().Render(context.Background(), w)
}

func (s *Server) Logout(w http.ResponseWriter, r *http.Request) {
	expiredCookie := &http.Cookie{
		Name:     constants.CookieName,
		Domain:   os.Getenv("COOKIE_DOMAIN"),
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}

	// Set the new cookie, which will effectively remove the old one
	http.SetCookie(w, expiredCookie)

	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}

func (s *Server) Read(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(ctx)
	q := env.Query(ctx)

	docID := chi.URLParam(r, "docID")

	currentUser, _ := env.UserClaim(ctx)

	var userID string
	if currentUser != nil {
		userID = currentUser.Id
	}

	doc, err := query.GetReadableDocumentForUser(q, docID, userID)
	if doc == nil || err != nil {
		if err != nil {
			log.Errorf("error getting document: %s", err)
		}
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)

		return
	}

	docStore := rogue.NewDocStore(s.S3, s.Query, s.Redis)
	_, rog, err := docStore.GetCurrentDoc(ctx, docID)
	if err != nil {
		log.Errorf("error getting document from doc store: %s", err)
		http.Redirect(w, r, "/not-found", http.StatusTemporaryRedirect)
		return
	}

	html, err := rog.GetFullHtml(
		true,
		false,
	)
	if err != nil {
		log.Errorf("error getting html: %s", err)
		http.Redirect(w, r, "/not-found", http.StatusTemporaryRedirect)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	read.Read(doc.ID, doc.Title, userID, html).Render(context.Background(), w)
}
