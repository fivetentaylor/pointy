package admin

import (
	"fmt"
	"net/http"

	"github.com/fivetentaylor/pointy/pkg/admin/templates"
	"github.com/fivetentaylor/pointy/pkg/env"
)

func GetIntroDoc(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := env.RawDB(ctx)

	var docID string
	result := db.Raw("SELECT doc_id FROM default_documents WHERE name = 'intro_doc'").Scan(&docID)

	if result.Error != nil {
		http.Error(w, "Failed to fetch intro doc", http.StatusInternalServerError)
		return
	}

	templates.IntroDoc(docID).Render(ctx, w)
}

func PostIntroDoc(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := env.RawDB(ctx)

	introDocID := r.FormValue("intro_doc_id")
	fmt.Printf("intro_doc_id: %s\n", introDocID)

	query := `INSERT INTO default_documents (name, doc_id)
            VALUES ('intro_doc', ?)
            ON CONFLICT (name)
            DO UPDATE SET doc_id = EXCLUDED.doc_id;`

	result := db.Exec(query, introDocID)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	templates.CurrentIntroDocID(introDocID).Render(ctx, w)
}
