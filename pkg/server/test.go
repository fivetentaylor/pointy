package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/teamreviso/code/pkg/constants"
	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/models"
	"github.com/teamreviso/code/pkg/query"
	"github.com/teamreviso/code/pkg/rogue"
	"github.com/teamreviso/code/pkg/server/auth"
	"github.com/teamreviso/code/pkg/service/document"
)

type CreateTestDocumentRequest struct {
	Email string `json:"email"`
}

// CreateTestDocument is only used for testing and is only available in test mode
// It will create a document and return a token to be used in testing
func (s *Server) CreateTestDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(ctx)

	request := &CreateTestDocumentRequest{}
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		log.Error(fmt.Sprintf("Error decoding request body: %s", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := &models.User{
		Email:       request.Email,
		Name:        "Test User",
		DisplayName: "Test User",
		Provider:    "google",
	}

	userTbl := env.Query(ctx).User
	err = userTbl.Create(user)
	if err != nil {
		log.Error(fmt.Sprintf("Error creating user: %s", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	doc, err := document.Create(ctx, user.ID)
	if err != nil {
		log.Error(fmt.Sprintf("Error creating document: %s", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := s.Auth.JWT().GenerateUserToken(user)
	if err != nil {
		log.Error(fmt.Sprintf("Error generating token: %s", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// set cookie
	cookie := &http.Cookie{
		Name:     constants.CookieName,
		Domain:   os.Getenv("COOKIE_DOMAIN"),
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(w, cookie)

	w.Write([]byte(fmt.Sprintf(`{"documentId": "%s"}`, doc.ID)))
}

const ViewTestDocumentTemplate = `
<!DOCTYPE html>
<html>
  <head>
	<title>Test Document</title>
	<link rel="stylesheet" href="/src/style/main.css"/>
  </head>
  <body>
	<h1>Test Document</h1>

	<div style="max-width: 65ch; margin: 0 auto; shadow: 0 2px 4px 0 rgba(0, 0, 0, 0.2), 0 25px 50px 0 rgba(0, 0, 0, 0.1);">
		<rogue-editor id=%q docid=%q apiHost=%q wsHost=%q>
			<div class="content p-4 ring-0 focus:outline-none">
			</div>
		</rogue-editor>
	</div>

	<script src="/src/rogue.js"></script>
  </body>
</html>
`

// ViewTestDocument is only used for testing and is only available in test mode
func (s *Server) ViewTestDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(ctx)
	q := env.Query(ctx)

	docID := chi.URLParam(r, "docID")

	user, err := query.GetOwnerOfDocument(q, docID)
	if err != nil {
		log.Error(fmt.Sprintf("Error getting user: %s", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jwt := auth.JWT(ctx)
	token, err := jwt.GenerateUserToken(user)
	if err != nil {
		log.Error(fmt.Sprintf("Error generating token: %s", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	atCookie := http.Cookie{
		Name:     constants.CookieName,
		Domain:   os.Getenv("COOKIE_DOMAIN"),
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(45 * time.Minute),
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}

	host := r.Host
	wsHost := fmt.Sprintf("ws://%s", host)
	apiHost := fmt.Sprintf("http://%s", host)

	if r.TLS != nil {
		wsHost = fmt.Sprintf("wss://%s", host)
		apiHost = fmt.Sprintf("https://%s", host)
	}

	outputHtml := fmt.Sprintf(ViewTestDocumentTemplate, docID, docID, apiHost, wsHost)

	http.SetCookie(w, &atCookie)
	w.Write([]byte(outputHtml))
}

func (s *Server) TestDocumentHTML(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(ctx)
	docID := chi.URLParam(r, "docID")
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

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func (s *Server) TestDocumentText(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(ctx)
	docID := chi.URLParam(r, "docID")
	docStore := rogue.NewDocStore(s.S3, s.Query, s.Redis)
	_, rog, err := docStore.GetCurrentDoc(ctx, docID)
	if err != nil {
		log.Errorf("error getting document from doc store: %s", err)
		http.NotFound(w, r)
		return
	}

	txt := rog.GetText()
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(txt))
}
