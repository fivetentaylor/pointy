package admin

import (
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jpoz/conveyor/config"
	"github.com/jpoz/conveyor/hub"

	"github.com/fivetentaylor/pointy/pkg/admin/src"
	"github.com/fivetentaylor/pointy/pkg/admin/templates"
	"github.com/fivetentaylor/pointy/pkg/constants"
	"github.com/fivetentaylor/pointy/pkg/env"
)

//go:embed static/*
var staticContent embed.FS

const adminRowQuery = `SELECT d.id, d.title, u.email, d.created_at, d.updated_at
FROM documents d
JOIN document_access da ON d.id = da.document_id
JOIN users u ON da.user_id = u.id
WHERE da.access_level = 'owner'
ORDER BY d.updated_at DESC
;`

func Routes(r chi.Router) {
	slogger := slog.Default()

	cSrv, err := hub.NewServer(&config.Hub{
		RedisURL:  os.Getenv("REDIS_URL"),
		Namespace: constants.DefaultJobNamespace,
		Logger:    slogger,
	})
	if err != nil {
		slogger.Error("error initializing hub server", "error", err)
		os.Exit(1)
	}

	r.Get("/", GetAdmin)
	r.Get("/users", GetUsers)
	r.Delete("/users/{userID}", DeleteUser)
	r.Get("/documents/{id}", GetDocument)
	r.Get("/documents/{id}/ai", GetDocumentAI)
	r.Get("/documents/{id}/ai/{key}", ShowLog)
	r.Get("/documents/{id}/dags", GetDags)
	r.Get("/documents/{id}/dags/{key}", GetDagLogFile)
	r.Get("/documents/{id}/dag/{dagName}/{dagId}", GetDagsForName)
	r.Get("/documents/{id}/snapshot", GetDocumentSnapshot)
	r.Post("/documents/{id}/snapshot", SnapshotDocument)
	r.Get("/documents/{id}/snapshots/{key}", GetDocumentSnapshotHTMLByKey)
	r.Get("/documents/{id}/snapshots/{key}/download", DownloadSnapshotFromKey)
	r.Post("/documents/{id}/snapshots/{key}/revert", RevertDocumentToSnapshot)
	r.Post("/documents/{id}/snapshots/{key}/new_document", NewDocumentFromSnapshot)
	r.Get("/documents/{id}/tree", GetDocumentTree)
	r.Get("/documents/{id}/table", GetTable)
	r.Get("/documents/{id}/table/history", HistoryTableForm)
	r.Get("/documents/{id}/table/history/{address}", GetHistoryTable)
	r.Get("/documents/{id}/log", GetDocumentLogs)
	r.Get("/documents/{id}/storage", GetStorage)
	r.Get("/documents/{id}/storage/{key}", GetStorage)
	r.Delete("/documents/{id}/storage/{opSeq}", DeleteDeltaLogItem)
	r.Get("/documents/{id}/edit", GetDocumentEditor)
	r.Get("/documents/{id}/address/{addressID}/diff", GetDocumentDiff)
	r.Get("/documents/{id}/address/{addressID}", GetContentAddress)
	r.Post("/documents/{id}/address", PostDocumentAddress)
	r.Get("/documents/{id}/address", GetContentAddressIDs)
	r.Get("/documents/{id}/messaging/threads", GetThreads)
	r.Get("/documents/{id}/messaging/threads/{threadId}", GetThread)
	r.Get("/documents/new", NewDocument)
	r.Post("/documents", CreateDocument)
	r.Get("/s3", GetS3)
	r.Post("/s3/list", GetS3ObjectList)
	r.Get("/s3/object/{key}", GetS3Object)
	r.Get("/jobs", GetJobs)
	r.Post("/jobs/start/{key}", StartJobs)
	r.Get("/prompts", GetPrompts)
	r.Get("/waitlist", WaitlistHandler)
	r.Post("/waitlist/update", UpdateWaitlistAccessHandler)
	r.Post("/prompts/refresh", RefreshPrompts)
	r.Post("/payment/subscription/plans", SyncSubscriptionPlans)
	r.Get("/payment/subscription/plans", SubscriptionPlans)
	r.Get("/intro_doc", GetIntroDoc)
	r.Post("/intro_doc", PostIntroDoc)

	// check dags
	r.Get("/checks/dags", GetDagChecks)
	r.Get("/checks/dags/new/document/{docId}/thread/{threadId}", NewDagCheck)
	r.Get("/checks/dags/{key}", GetChecksForDag)
	r.Get("/checks/dags/{key}/{id}/results/{resultId}", ViewResultForDagCheck)
	r.Get("/checks/dags/{key}/{id}/results", GetResultsForDagCheck)
	r.Get("/checks/dags/{key}/{id}/examples/{exampleId}", ViewExampleForDagCheck)
	r.Post("/checks/dags", CreateDagCheck)
	r.Post("/checks/dags/check/{id}/examples", CreateDagCheckExample)
	r.Post("/checks/dags/{key}/run/{id}", RunDagCheck)

	r.Handle("/src/*", src.BuildHandler("/admin/src/"))
	r.Handle("/static/*", Static("/admin/static/"))

	r.Route("/conveyor", func(r chi.Router) {
		r.Handle("/*", cSrv.Handler(
			hub.HandlerOpts{
				Prefix:  "/admin/conveyor",
				HomeURL: "/admin",
			},
		))
	})
}

func GetAdmin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	docs := []templates.AdminDocument{}

	// load last 10 docs
	err := env.RawDB(ctx).Raw(adminRowQuery).Scan(&docs).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templates.Admin(docs).Render(ctx, w)
}

func Static(remove string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		trimmedPath := strings.TrimPrefix(r.URL.Path, remove)
		// Ensure the path is sanitized to avoid directory traversal issues.
		path := filepath.Join("static", filepath.Clean("/"+trimmedPath))
		fmt.Println(path)

		// Avoid serving directory paths.
		if strings.HasSuffix(r.URL.Path, "/") {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		// Read the file content from the embedded file system.
		data, err := fs.ReadFile(staticContent, path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		// Determine the content type of the file.
		contentType := mime.TypeByExtension(filepath.Ext(path))
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		// Set the Content-Type header.
		w.Header().Set("Content-Type", contentType)

		// Write the content to the response writer.
		// Note: This is not using http.ServeContent as it requires io.ReadSeeker.
		// If caching or range requests are not a concern, this method is simpler.
		w.Write(data)
	}
}
