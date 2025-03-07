package admin

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/fivetentaylor/pointy/pkg/admin/templates"
	"github.com/fivetentaylor/pointy/pkg/constants"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/rogue"
	"github.com/fivetentaylor/pointy/pkg/service/document"
	v3 "github.com/fivetentaylor/pointy/rogue/v3"
)

func NewDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(ctx)
	log.Info("[admin] document new")

	templates.NewDocument().Render(r.Context(), w)
}

func CreateDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(ctx)
	currentUser, err := env.UserClaim(ctx)
	if err != nil {
		log.Error("[admin] error getting user claim", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method != http.MethodPost {
		log.Error("[admin] invalid request method", "method", r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form data
	err = r.ParseForm()
	if err != nil {
		log.Error("[admin] error parsing form", "error", err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Retrieve the document content from the form
	docContent := r.FormValue("doc")
	if docContent == "" {
		log.Error("[admin] document content cannot be empty")
		http.Error(w, "Document content cannot be empty", http.StatusBadRequest)
		return
	}

	sr := v3.SerializedRogue{}
	err = json.Unmarshal([]byte(docContent), &sr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a rogue and apply the snapshot in case it's in the older format
	rg := v3.NewRogue(constants.RevisoAuthorID)
	rg.MergeOp(v3.SnapshotOp{
		Snapshot: &sr,
	})

	ops, err := rg.ToOps()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userTbl := env.Query(ctx).User

	owner, err := userTbl.Where(userTbl.ID.Eq(currentUser.Id)).First()
	if err != nil {
		log.Error("[admin] error getting user", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a new Document
	doc, err := document.Create(ctx, owner.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	newTitle := "Untitled (from snapshot)"
	doc, err = document.UpdateDocument(ctx, doc, owner.ID, &newTitle, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s3 := env.S3(ctx)
	redis := env.Redis(ctx)
	query := env.Query(ctx)
	store := rogue.NewDocStore(s3, query, redis)
	for _, op := range ops {
		_, err = store.AddUpdate(ctx, doc.ID, op)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = store.SnapshotDoc(ctx, doc.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/drafts/%s", doc.ID), http.StatusSeeOther)
}

func GetDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(ctx)
	docID := chi.URLParam(r, "id")

	log.Info("[admin] document loading", "id", docID)

	s3 := env.S3(ctx)
	redis := env.Redis(ctx)
	query := env.Query(ctx)
	dynamo := env.Dynamo(ctx)
	store := rogue.NewDocStore(s3, query, redis)

	_, doc, err := store.GetCurrentDoc(ctx, docID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	snapshot, err := doc.NewSnapshotOp()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = doc.MergeOp(snapshot)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	startID := v3.RootID
	endID := v3.LastID

	params := r.URL.Query()
	if params.Get("startID") != "" {
		startID, err = v3.ParseID(params.Get("startID"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if params.Get("endID") != "" {
		endID, err = v3.ParseID(params.Get("endID"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	formatOps, err := doc.Formats.SearchOverlapping(startID, endID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Header.Get("Hx-Request") == "true" {
		templates.DocumentDebugger(docID, doc, formatOps, startID, endID).Render(r.Context(), w)
		return
	}

	// addressIDs, nextPage, err := dynamo.GetContentAddressIDs(docID, nil)
	addressIDs, _, err := dynamo.GetContentAddressIDs(docID, nil)
	if err != nil {
		log.Errorf("failed to get content address ids: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templates.Document(docID, doc, formatOps, startID, endID, addressIDs).Render(r.Context(), w)
}

func GetDocumentTree(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(ctx)
	docID := chi.URLParam(r, "id")

	log.Info("[admin] document tree loading", "id", docID)

	s3 := env.S3(ctx)
	redis := env.Redis(ctx)
	query := env.Query(ctx)
	store := rogue.NewDocStore(s3, query, redis)

	_, doc, err := store.GetCurrentDoc(ctx, docID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("[admin] document tree loaded", "id", docID, "doc", doc)

	// TODO: Fix this to represent multiple roots!
	templates.Tree(docID, doc.Roots[0]).Render(r.Context(), w)
}

func GetDocumentAI(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(ctx)
	docID := chi.URLParam(r, "id")

	s3 := env.S3(r.Context())
	keys, err := s3.List(s3.Bucket, fmt.Sprintf(constants.ConvosPrefix, docID), -1, -1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("[admin] document ai loaded", "id", docID)

	logFiles := make([]templates.LogFile, len(keys))
	for i, key := range keys {
		logFiles[i], err = templates.LogFileFromKey(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	templates.AI(docID, logFiles).Render(r.Context(), w)
}

func ShowLog(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	key64 := chi.URLParam(r, "key")

	key, err := base64.StdEncoding.DecodeString(key64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s3 := env.S3(ctx)
	object, err := s3.GetObject(s3.Bucket, string(key))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(object)
}

func GetDocumentEditor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(ctx)
	docID := chi.URLParam(r, "id")

	log.Info("[admin] document loading", "id", docID)

	s3 := env.S3(ctx)
	redis := env.Redis(ctx)
	query := env.Query(ctx)
	store := rogue.NewDocStore(s3, query, redis)

	_, doc, err := store.GetCurrentDoc(ctx, docID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	snapshot, err := doc.NewSnapshotOp()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = doc.MergeOp(snapshot)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("[admin] document loaded", "id", docID, "doc", doc)

	templates.Editor(docID, doc).Render(r.Context(), w)
}
