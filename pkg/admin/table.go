package admin

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/fivetentaylor/pointy/pkg/admin/templates"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/rogue"
	v3 "github.com/fivetentaylor/pointy/rogue/v3"
)

func GetTable(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(ctx)
	docID := chi.URLParam(r, "id")

	log.Info("[admin] document table loading", "id", docID)

	s3 := env.S3(ctx)
	redis := env.Redis(ctx)
	query := env.Query(ctx)
	store := rogue.NewDocStore(s3, query, redis)

	_, doc, err := store.GetCurrentDoc(ctx, docID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("[admin] document table loaded", "id", docID)

	templates.Table(docID, doc).Render(r.Context(), w)
}

func HistoryTableForm(w http.ResponseWriter, r *http.Request) {
	docID := chi.URLParam(r, "id")

	// check if there's a git param "address"
	addressStr := r.URL.Query().Get("address")
	if addressStr != "" {
		GetHistoryTable(w, r)
		return
	}

	templates.HistoryForm(docID).Render(r.Context(), w)
}

func GetHistoryTable(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(ctx)
	docID := chi.URLParam(r, "id")
	addressStr := r.URL.Query().Get("address")

	log.Info("[admin] document table loading", "id", docID)

	s3 := env.S3(ctx)
	redis := env.Redis(ctx)
	query := env.Query(ctx)
	store := rogue.NewDocStore(s3, query, redis)

	_, doc, err := store.GetCurrentDoc(ctx, docID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("[admin] document history table loaded", "id", docID, "address", addressStr)

	ca := v3.ContentAddress{}
	err = json.Unmarshal([]byte(addressStr), &ca)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	startID := ca.StartID
	endID := ca.EndID

	vis, span, line, err := doc.ToIndexNos(startID, endID, &ca, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fVis := &v3.FugueVis{
		Text: vis.Text,
		IDs:  vis.IDs,
	}

	html, err := v3.ToHtml(fVis, span, line, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templates.HistoryTable(docID, addressStr, vis, span, line, html).Render(r.Context(), w)
}
