package admin

import (
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	awsS3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/teamreviso/code/pkg/admin/templates"
	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/rogue"
)

func GetStorage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	docID := chi.URLParam(r, "id")

	// get snapshots
	s3 := env.S3(ctx)
	redis := env.Redis(ctx)
	query := env.Query(ctx)
	store := rogue.NewDocStore(s3, query, redis)
	prefix := rogue.DocS3Path(docID)

	input := &awsS3.ListObjectsV2Input{
		Bucket:  aws.String(s3.Bucket),
		Prefix:  aws.String(prefix),
		MaxKeys: aws.Int64(40),
	}

	result, err := s3.Client.ListObjectsV2(input)
	if err != nil {
		log.Errorf("failed to list objects: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	keys := make([]templates.DocumentStorageSnapshot, 0, len(result.Contents))
	for _, obj := range result.Contents {
		keys = append(keys, templates.DocumentStorageSnapshot{
			Key:          *obj.Key,
			LastModified: *obj.LastModified,
		})
	}

	// get delta log
	ops, err := store.GetDeltaLog(ctx, docID, 0, -1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pendingOps := make([]templates.PendingOp, 0, len(ops))
	for _, op := range ops {
		pendingOps = append(pendingOps, templates.PendingOp{
			Op:    op.Member.(string),
			Score: op.Score,
		})
	}

	templates.DocumentStorage(docID, keys, pendingOps).Render(ctx, w)
}

func DeleteDeltaLogItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	docID := chi.URLParam(r, "id")
	opSeq, err := strconv.Atoi(chi.URLParam(r, "opSeq"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get snapshots
	s3 := env.S3(ctx)
	redis := env.Redis(ctx)
	query := env.Query(ctx)
	store := rogue.NewDocStore(s3, query, redis)

	err = store.DeleteDeltaLogItem(ctx, docID, int64(opSeq))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ops, err := store.GetDeltaLog(ctx, docID, 0, -1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pendingOps := make([]templates.PendingOp, 0, len(ops))
	for _, op := range ops {
		pendingOps = append(pendingOps, templates.PendingOp{
			Op:    op.Member.(string),
			Score: op.Score,
		})
	}

	// templates.DocumentStorage(docID, keys, pendingOps).Render(ctx, w)
	templates.BufferedOps(docID, pendingOps).Render(ctx, w)
}
