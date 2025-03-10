package admin

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/fivetentaylor/pointy/pkg/admin/templates"
	"github.com/fivetentaylor/pointy/pkg/background/wire"
	"github.com/fivetentaylor/pointy/pkg/env"
)

func GetJobs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	templates.Jobs().Render(ctx, w)
}

func StartJobs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	jobKey := chi.URLParam(r, "key")

	switch jobKey {
	case "screenshot-all":
		_, err := env.Background(ctx).Enqueue(ctx, &wire.ScreenshotAll{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	case "snapshot":
		// get form data
		docID := r.FormValue("document_id")
		_, err := env.Background(ctx).Enqueue(ctx, &wire.SnapshotRogue{DocId: docID})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// write response string
		fmt.Fprint(w, fmt.Sprintf("snapshot job started for document: %s", docID))
		w.WriteHeader(http.StatusOK)
		return
	case "snapshot-all":
		// get form data
		version := r.FormValue("version")
		_, err := env.Background(ctx).Enqueue(ctx, &wire.SnapshotAll{Version: version})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// write response string
		fmt.Fprint(w, fmt.Sprintf("snapshot all job started with version: %s", version))
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Error(w, "no job found", http.StatusNotFound)
}
