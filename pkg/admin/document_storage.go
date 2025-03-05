package admin

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/teamreviso/code/pkg/admin/templates"
	"github.com/teamreviso/code/pkg/background/wire"
	"github.com/teamreviso/code/pkg/constants"
	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/rogue"
	"github.com/teamreviso/code/pkg/service/document"
	v3 "github.com/teamreviso/code/rogue/v3"
)

func GetDocumentSnapshot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(ctx)
	docID := chi.URLParam(r, "id")

	log.Info("[admin] document snapshot loading", "id", docID)

	s3 := env.S3(ctx)
	redis := env.Redis(ctx)
	query := env.Query(ctx)
	store := rogue.NewDocStore(s3, query, redis)

	// Maybe snapshot the doc to s3 when this is clicked to clear out redis?
	// store.SnapshotDoc(ctx, docID)

	_, doc, err := store.GetCurrentDoc(ctx, docID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("[admin] document snapshot loaded", "id", docID, "doc", doc)

	docBytes, err := json.Marshal(doc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	_, err = w.Write([]byte(docBytes))
	if err != nil {
		log.Error("Error writing response", "error", err)

	}
}

func GetDocumentSnapshotHTMLByKey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(ctx)
	docID := chi.URLParam(r, "id")
	key := chi.URLParam(r, "key")

	deckey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("[admin] document snapshot html", "id", docID, "key", key)

	s3 := env.S3(ctx)
	object, err := s3.GetObject(s3.Bucket, string(deckey))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rogue := v3.NewRogue(constants.RevisoAuthorID)

	sr := v3.SerializedRogue{}
	err = json.Unmarshal(object, &sr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rogue.MergeOp(v3.SnapshotOp{
		Snapshot: &sr,
	})

	html, err := rogue.GetHtml(v3.RootID, v3.LastID, true, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templates.DocumentStorageSnapshotView(docID, string(deckey), html).Render(r.Context(), w)
}

func SnapshotDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(ctx)
	docID := chi.URLParam(r, "id")

	store := rogue.NewDocStore(env.S3(ctx), env.Query(ctx), env.Redis(ctx))

	err := store.SnapshotDoc(ctx, docID)
	if err != nil {
		log.Error("Error snapshotting document", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = env.Background(ctx).Enqueue(ctx, &wire.Screenshot{DocId: docID})
	if err != nil {
		log.Error("Error snapshotting document", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Redirect", fmt.Sprintf("/admin/documents/%s/storage", docID))
	w.WriteHeader(http.StatusOK)
}

func DownloadSnapshotFromKey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(ctx)
	docID := chi.URLParam(r, "id")
	key := chi.URLParam(r, "key")

	deckey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("[admin] document snapshot html", "id", docID, "key", key)

	s3 := env.S3(ctx)
	object, err := s3.GetObject(s3.Bucket, string(deckey))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sr := v3.SerializedRogue{}
	err = json.Unmarshal(object, &sr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(sr)
	if err != nil {
		log.Error("Error writing response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func NewDocumentFromSnapshot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(ctx)
	docID := chi.URLParam(r, "id")
	key := chi.URLParam(r, "key")

	docTbl := env.Query(ctx).Document
	userTbl := env.Query(ctx).User
	docAccessTbl := env.Query(ctx).DocumentAccess

	parentDoc, err := docTbl.Where(docTbl.ID.Eq(docID)).First()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	owner, err := userTbl.
		LeftJoin(docAccessTbl, userTbl.ID.EqCol(docAccessTbl.UserID)).
		Where(docAccessTbl.DocumentID.Eq(parentDoc.ID)).
		Where(docAccessTbl.AccessLevel.Eq("owner")).
		First()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	deckey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("[admin] document snapshot html", "id", docID, "key", key)

	s3 := env.S3(ctx)
	object, err := s3.GetObject(s3.Bucket, string(deckey))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sr := v3.SerializedRogue{}
	err = json.Unmarshal(object, &sr)
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

	// Create a new Document
	doc, err := document.Create(ctx, owner.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	newTitle := fmt.Sprintf("Copy - %s", parentDoc.Title)
	doc, err = document.UpdateDocument(ctx, doc, owner.ID, &newTitle, &parentDoc.IsPublic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

	_, err = env.Background(ctx).Enqueue(ctx, &wire.Screenshot{DocId: docID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Redirect", fmt.Sprintf("/admin/documents/%s", doc.ID))
	w.WriteHeader(http.StatusOK)
}

func RevertDocumentToSnapshot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(ctx)
	docID := chi.URLParam(r, "id")
	key := chi.URLParam(r, "key")

	deckey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("[admin] document snapshot html", "id", docID, "key", key)

	s3 := env.S3(ctx)
	object, err := s3.GetObject(s3.Bucket, string(deckey))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sr := v3.SerializedRogue{}
	err = json.Unmarshal(object, &sr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rg := v3.NewRogue(constants.RevisoAuthorID)
	_, err = rg.MergeOp(v3.SnapshotOp{Snapshot: &sr})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ops, err := rg.ToOps()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	redis := env.Redis(ctx)
	query := env.Query(ctx)
	store := rogue.NewDocStore(s3, query, redis)

	// Get lastest snapshot
	seq, _, err := store.GetLastSnapshot(docID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ops = append(ops, v3.FormatOp{
		ID:      v3.ID{Author: constants.RevisoAuthorID, Seq: int(seq) + 1},
		StartID: v3.RootID,
		EndID:   v3.RootID,
		Format:  v3.FormatV3Span{},
	})

	err = store.SaveDocToS3(ctx, docID, seq+1, rg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = env.Background(ctx).Enqueue(ctx, &wire.Screenshot{DocId: docID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Redirect", fmt.Sprintf("/admin/documents/%s", docID))
	w.WriteHeader(http.StatusOK)
}
