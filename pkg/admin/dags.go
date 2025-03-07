package admin

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/fivetentaylor/pointy/pkg/admin/templates"
	"github.com/fivetentaylor/pointy/pkg/ai"
	"github.com/fivetentaylor/pointy/pkg/constants"
	"github.com/fivetentaylor/pointy/pkg/dag"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/service/rogue"
	"github.com/fivetentaylor/pointy/pkg/storage/dynamo"
	v3 "github.com/fivetentaylor/pointy/rogue/v3"
)

var DagMap = map[string]func() *dag.Dag{
	ai.ThreadDagV2().Name: ai.ThreadDagV2,
}

func GetDags(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(ctx)
	docID := chi.URLParam(r, "id")

	s3 := env.S3(r.Context())
	keys, err := s3.List(s3.Bucket, fmt.Sprintf(constants.DagsDirByParentId, docID), -1, -1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("[admin] document dags loaded", "id", docID)

	dagFiles := make([]*templates.DagRunFile, len(keys))
	for i, key := range keys {
		dagFiles[i], err = templates.DagRunFileFromKey(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	templates.Dag(docID, dagFiles).Render(r.Context(), w)
}

func GetDagsForName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(ctx)
	docID := chi.URLParam(r, "id")
	name := chi.URLParam(r, "dagName")
	dagId := chi.URLParam(r, "dagId")

	s3 := env.S3(r.Context())
	keys, err := s3.List(s3.Bucket, fmt.Sprintf(constants.DagsDirByParentId, docID), -1, -1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filteredKeys := make([]string, 0)
	for _, key := range keys {
		if strings.Contains(key, name) && strings.Contains(key, dagId) {
			filteredKeys = append(filteredKeys, key)
		}
	}

	log.Info("[admin] filtered document dags loaded", "id", docID, "keys", filteredKeys, "name", name, "dagId", dagId)

	dagFiles := make([]*templates.DagRunFile, len(filteredKeys))
	for i, key := range filteredKeys {
		dagFiles[i], err = templates.DagRunFileFromKey(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	templates.Dag(docID, dagFiles).Render(r.Context(), w)
}

func GetDagLogFile(w http.ResponseWriter, r *http.Request) {
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

	escapedString := html.EscapeString(string(object))

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(escapedString))
}

func GetDagChecks(w http.ResponseWriter, r *http.Request) {
	templates.DagChecks(DagMap).Render(r.Context(), w)
}

func GetChecksForDag(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(r.Context())
	key64 := chi.URLParam(r, "key")
	key, err := base64.StdEncoding.DecodeString(key64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("[admin] loading dag", "key", string(key))

	dagFn, ok := DagMap[string(key)]
	if !ok {
		http.Error(w, fmt.Sprintf("dag %q not found", string(key)), http.StatusNotFound)
		return
	}

	files, err := dag.ListFuncationalChecks(ctx, string(key))
	if err != nil {
		log.Error("error listing files", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("[admin] loading dag checks", "key", string(key), "files", len(files))

	templates.ChecksForDag(dagFn(), files).Render(r.Context(), w)
}

func GetResultsForDagCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(r.Context())
	key64 := chi.URLParam(r, "key")
	id := chi.URLParam(r, "id")
	key, err := base64.StdEncoding.DecodeString(key64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("[admin] loading dag", "key", string(key))

	dagFn, ok := DagMap[string(key)]
	if !ok {
		http.Error(w, fmt.Sprintf("dag %q not found", string(key)), http.StatusNotFound)
		return
	}

	file, err := dag.GetFuncationalCheck(ctx, string(key), id)
	if err != nil {
		log.Error("error listing files", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("[admin] loading dag check", "key", string(key), "id", id)

	results, err := dag.ListFuncationalCheckResults(ctx, id)
	if err != nil {
		log.Error("error listing files", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	examples, err := dag.ListFuncationalCheckExamples(ctx, file.DagName, id)
	if err != nil {
		log.Error("error listing files", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templates.DagCheck(dagFn(), file, results, examples).Render(r.Context(), w)
}

func ViewExampleForDagCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(r.Context())
	key64 := chi.URLParam(r, "key")
	id := chi.URLParam(r, "id")
	exampleId := chi.URLParam(r, "exampleId")
	key, err := base64.StdEncoding.DecodeString(key64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("[admin] viewing example for dag", "key", string(key), "id", id, "exampleId", exampleId)

	dagFn, ok := DagMap[string(key)]
	if !ok {
		http.Error(w, fmt.Sprintf("dag %q not found", string(key)), http.StatusNotFound)
		return
	}

	file, err := dag.GetFuncationalCheck(ctx, string(key), id)
	if err != nil {
		log.Error("error getting file", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	example, err := dag.GetFuncationalCheckExample(ctx, file.DagName, id, exampleId)
	if err != nil {
		log.Error("error getting result", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var checkDoc *v3.Rogue
	err = json.Unmarshal([]byte(file.SerializedRogue), &checkDoc)
	if err != nil {
		log.Error("error unmarshalling check doc", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var beforeDoc *v3.Rogue
	err = json.Unmarshal([]byte(example.SerializedRogueBefore), &beforeDoc)
	if err != nil {
		log.Error("error unmarshalling before result doc", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var resultDoc *v3.Rogue
	err = json.Unmarshal([]byte(example.SerializedRogueResult), &resultDoc)
	if err != nil {
		log.Error("error unmarshalling result doc", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var messages []*dynamo.Message
	err = json.Unmarshal([]byte(example.SerializedThreadResult), &messages)
	if err != nil {
		log.Error("error unmarshalling result thread", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := templates.ViewDagCheckExampleData{
		Dag:       dagFn(),
		Check:     file,
		CheckDoc:  checkDoc,
		Messages:  messages,
		Example:   example,
		ResultDoc: resultDoc,
		BeforeDoc: beforeDoc,
	}

	templates.ViewDagCheckExample(data).Render(r.Context(), w)
}

func ViewResultForDagCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(r.Context())
	key64 := chi.URLParam(r, "key")
	id := chi.URLParam(r, "id")
	resultId := chi.URLParam(r, "resultId")
	key, err := base64.StdEncoding.DecodeString(key64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("[admin] viewing results for dag", "key", string(key), "id", id, "resultId", resultId)

	dagFn, ok := DagMap[string(key)]
	if !ok {
		http.Error(w, fmt.Sprintf("dag %q not found", string(key)), http.StatusNotFound)
		return
	}

	file, err := dag.GetFuncationalCheck(ctx, string(key), id)
	if err != nil {
		log.Error("error getting file", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := dag.GetFuncationalCheckResult(ctx, id, resultId)
	if err != nil {
		log.Error("error getting result", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var checkDoc *v3.Rogue
	err = json.Unmarshal([]byte(file.SerializedRogue), &checkDoc)
	if err != nil {
		log.Error("error unmarshalling check doc", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var beforeDoc *v3.Rogue
	err = json.Unmarshal([]byte(result.SerializedRogueBefore), &beforeDoc)
	if err != nil {
		log.Error("error unmarshalling before result doc", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var resultDoc *v3.Rogue
	err = json.Unmarshal([]byte(result.SerializedRogueResult), &resultDoc)
	if err != nil {
		log.Error("error unmarshalling result doc", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var messages []*dynamo.Message
	err = json.Unmarshal([]byte(result.SerializedThreadResult), &messages)
	if err != nil {
		log.Error("error unmarshalling result thread", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	example, err := dag.GetFuncationalCheckExample(ctx, result.DagName, id, resultId)
	if err != nil {
		log.Error("error getting example", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := templates.ViewDagCheckResultData{
		Dag:             dagFn(),
		Check:           file,
		CheckDoc:        checkDoc,
		Messages:        messages,
		Result:          result,
		ResultDoc:       resultDoc,
		BeforeDoc:       beforeDoc,
		ExistingExample: example,
	}

	templates.ViewDagCheckResult(data).Render(r.Context(), w)
}

func NewDagCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	docID := chi.URLParam(r, "docId")
	threadID := chi.URLParam(r, "threadId")

	doc, err := rogue.CurrentDocument(ctx, docID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templates.NewDagCheck(doc, docID, threadID, DagMap).Render(ctx, w)
}

func CreateDagCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form data.", http.StatusBadRequest)
		return
	}
	documentID := r.FormValue("document_id")
	threadID := r.FormValue("thread_id")
	dagName := r.FormValue("dag")
	checkName := r.FormValue("checkName")

	if documentID == "" || threadID == "" || dagName == "" || checkName == "" {
		fmt.Fprintf(w, "All fields are required.")
		return
	}

	fc, err := dag.CreateFunctionalCheckFile(ctx, documentID, threadID, dagName, checkName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = fc.Save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	key64 := base64.StdEncoding.EncodeToString([]byte(fc.DagName))
	http.Redirect(w, r, fmt.Sprintf("/admin/checks/dags/%s", key64), http.StatusFound)
}

func CreateDagCheckExample(w http.ResponseWriter, r *http.Request) {
	log := env.Log(r.Context())
	log.Info("[admin] creating dag check example")

	ctx := r.Context()
	if err := r.ParseForm(); err != nil {
		log.Error("error parsing form data", "error", err)
		http.Error(w, "Error parsing form data.", http.StatusBadRequest)
		return
	}
	checkId := chi.URLParam(r, "id")
	resultId := r.FormValue("result_id")
	approved := r.FormValue("approved")
	if resultId == "" {
		fmt.Fprintf(w, "All fields are required.")
		return
	}

	result, err := dag.GetFuncationalCheckResult(ctx, checkId, resultId)
	if err != nil {
		log.Error("error getting result", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	example, err := dag.CreateFuncationalCheckExampleFile(ctx, result)
	if err != nil {
		log.Error("error creating example", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if approved == "true" {
		example.Approved = true
	}

	err = example.Save()
	if err != nil {
		log.Error("error saving example", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	key64 := base64.StdEncoding.EncodeToString([]byte(example.DagName))
	http.Redirect(w, r, fmt.Sprintf("/admin/checks/dags/%s/%s/results/%s", key64, checkId, resultId), http.StatusFound)
}

func RunDagCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(ctx)
	key64 := chi.URLParam(r, "key")
	key, err := base64.StdEncoding.DecodeString(key64)
	if err != nil {
		log.Error("error decoding key", "key", key64, "error", err)
		templates.RunFailed(err).Render(ctx, w)
		return
	}

	id := chi.URLParam(r, "id")

	fnCheck, err := dag.GetFuncationalCheck(ctx, string(key), id)
	if err != nil {
		log.Error("error getting check", "key", string(key), "id", id, "error", err)
		templates.RunFailed(err).Render(ctx, w)
		return
	}

	currentUser, err := env.UserClaim(ctx)
	if err != nil {
		log.Error("error getting user", "key", string(key), "id", id, "error", err)
		templates.RunFailed(err).Render(ctx, w)
		return
	}

	dagFn, ok := DagMap[string(key)]
	if !ok {
		log.Error("dag not found", "key", string(key), "id", id)
		templates.RunFailed(fmt.Errorf("dag %s not found", string(key))).Render(ctx, w)
		return
	}

	result, err := fnCheck.Run(ctx, currentUser.Id, dagFn())
	if err != nil {
		log.Error("error running check", "key", string(key), "id", id, "error", err)
		templates.RunFailed(err).Render(ctx, w)
		return
	}

	templates.ViewButton(key64, result).Render(ctx, w)
}
