package admin

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/fivetentaylor/pointy/pkg/admin/templates"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/rogue"
)

func GetDocumentLogs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(ctx)
	docID := chi.URLParam(r, "id")

	log.Info("[admin] loading document logs", "id", docID)

	s3 := env.S3(ctx)
	logger := rogue.NewLogger(log, s3, docID, "")

	entries, err := logger.Read()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templates.Logs(docID, entries).Render(ctx, w)
}
