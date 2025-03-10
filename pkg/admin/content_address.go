package admin

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/fivetentaylor/pointy/pkg/admin/templates"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/rogue"

	v3 "github.com/fivetentaylor/pointy/rogue/v3"
)

func PostDocumentAddress(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(ctx)

	docID := chi.URLParam(r, "id")

	log.Info("[admin] document tree loading", "id", docID)

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

	startIDstr := r.FormValue("startID")
	endIDstr := r.FormValue("endID")

	startID, err := v3.ParseID(startIDstr)
	if err != nil {
		log.Errorf("failed to parse startID: %s", startIDstr)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	endID, err := v3.ParseID(endIDstr)
	if err != nil {
		log.Errorf("failed to parse endID: %s", endIDstr)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	address, err := doc.GetAddress(startID, endID)
	if err != nil {
		log.Errorf("failed to get address: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	addressBytes, err := json.Marshal(address)
	if err != nil {
		log.Errorf("failed to marshal address: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Infof("startID: %s, endID: %s address: %v", startID, endID, address)

	addressID, err := dynamo.CreateContentAddress(docID, addressBytes)
	if err != nil {
		log.Errorf("failed to create content address: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Infof("addressID: %s", addressID)
}

func GetContentAddressIDs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(ctx)

	docID := chi.URLParam(r, "id")

	log.Info("[admin] document tree loading", "id", docID)

	dynamo := env.Dynamo(ctx)

	addressIDs, _, err := dynamo.GetContentAddressIDs(docID, nil)
	if err != nil {
		log.Errorf("failed to get content address ids: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Infof("addressIDs: %v", addressIDs)
	templates.ContentAddresses(docID, addressIDs).Render(ctx, w)
}

func GetContentAddress(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(ctx)

	docID := chi.URLParam(r, "id")
	addressID := chi.URLParam(r, "addressID")

	dynamo := env.Dynamo(ctx)
	addressBytes, err := dynamo.GetContentAddress(docID, addressID)
	if err != nil {
		log.Errorf("failed to get content address: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(addressBytes)
}

const tpl = `
<!DOCTYPE html>
<html>
<head>
    <title>Document Diff</title>
    <style>
        /* Add your CSS styles here */
        body { font-family: Arial, sans-serif; }

		del {background-color: #FFBABA;}
		ins {background-color: #BAFFC9;}

		li[ql-indent="1"] {
			margin-left: 20px;
		}

		li[ql-indent="2"] {
			margin-left: 40px;
		}

		li[ql-indent="3"] {
			margin-left: 60px;
		}
    </style>
</head>
<body>
    {{.DiffHTML}}
</body>
</html>
`

func GetDocumentDiff(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(ctx)

	docID := chi.URLParam(r, "id")
	addressID := chi.URLParam(r, "addressID")

	s3 := env.S3(ctx)
	redis := env.Redis(ctx)
	query := env.Query(ctx)
	store := rogue.NewDocStore(s3, query, redis)

	_, doc, err := store.GetCurrentDoc(ctx, docID)
	if err != nil {
		log.Errorf("failed to get current doc: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dynamo := env.Dynamo(ctx)
	addressBytes, err := dynamo.GetContentAddress(docID, addressID)
	if err != nil {
		log.Errorf("failed to get content address: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	address := v3.ContentAddress{}
	err = json.Unmarshal(addressBytes, &address)
	if err != nil {
		log.Errorf("failed to unmarshal content address: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	firstID, err := doc.Rope.GetTotID(0)
	if err != nil {
		log.Errorf("failed to get first ID: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	lastID, err := doc.Rope.GetTotID(doc.TotSize - 1)
	if err != nil {
		log.Errorf("failed to get last ID: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	html, err := doc.GetHtmlDiff(firstID, lastID, &address, true, false)
	if err != nil {
		log.Errorf("failed to get diff html: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// w.Header().Set("Content-Type", "text/html")
	// w.Write([]byte(html))

	t, err := template.New("webpage").Parse(tpl)
	if err != nil {
		http.Error(w, "Error creating template", http.StatusInternalServerError)
		return
	}

	data := struct {
		DiffHTML template.HTML // Use template.HTML to avoid auto-escaping
	}{
		DiffHTML: template.HTML(html),
	}

	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}
