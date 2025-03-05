package admin

import (
	"encoding/base64"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/teamreviso/code/pkg/admin/templates"
	"github.com/teamreviso/code/pkg/env"
	s3client "github.com/teamreviso/code/pkg/storage/s3"
	"github.com/teamreviso/code/pkg/utils"
)

var (
	s3keys   = utils.NewThreadSafeSlice[string]() // Thread-safe slice to hold S3 objects.
	loadOnce sync.Once                            // Ensure S3 loading is initiated only once.
)

func loadS3DataInBackground(s *s3client.S3, bucket string) {
	loadOnce.Do(func() {
		go func() {
			input := &s3.ListObjectsV2Input{
				Bucket:  aws.String(bucket),
				MaxKeys: aws.Int64(1000),
			}

			err := s.Client.ListObjectsV2Pages(input,
				func(page *s3.ListObjectsV2Output, lastPage bool) bool {
					for _, obj := range page.Contents {
						s3keys.Add(*obj.Key)
					}

					return !lastPage
				})

			if err != nil {
				log.Errorf("failed to list objects: %s", err)
			}
		}()
	})
}

func GetS3(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(ctx)

	log.Info("[admin] document storage loading")

	s3 := env.S3(ctx)
	// redis := env.Redis(ctx)

	loadS3DataInBackground(s3, s3.Bucket)
	for s3keys.Length() < 2000 {
		log.Info("[admin] waiting for s3 data to load")
		time.Sleep(100 * time.Millisecond)
	}

	log.Info("[admin] document storage loaded")

	objects := s3keys.Slice(0, 1000)
	pageCount := s3keys.Length() / 1000
	templates.S3(s3.Bucket, pageCount, objects).Render(r.Context(), w)
}

func filterS3keys(w http.ResponseWriter, pattern string) *utils.ThreadSafeSlice[string] {
	if pattern == "" {
		return s3keys
	}

	r, err := regexp.Compile(pattern)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	objects, err := s3keys.Select(func(key string) (bool, error) {
		return r.MatchString(key), nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	out := utils.NewThreadSafeSlice[string]()
	out.Objects = objects

	return out
}

func GetS3ObjectList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(ctx)

	// Retrieve the value of the form field "prefix"
	pattern := r.FormValue("pattern")
	action := r.FormValue("action")
	pageIx, err := strconv.Atoi(r.FormValue("pageIx"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s3 := env.S3(ctx)
	loadS3DataInBackground(s3, s3.Bucket)
	for s3keys.Length() < 2000 {
		log.Info("[admin] waiting for s3 data to load")
		time.Sleep(100 * time.Millisecond)
	}

	log.Info("[admin] s3 objects listed")

	var objects []string
	filteredKeys := filterS3keys(w, pattern)
	pageCount := filteredKeys.Length() / 1000

	log.Infof("pattern: %s, action: %s, pageIx: %d, pageCount: %d", pattern, action, pageIx, pageCount)

	if action == "prev" {
		pageIx--
	} else if action == "next" {
		pageIx++
	} else {
		pageIx = 1
	}

	ix := pageIx - 1
	objects = filteredKeys.Slice(ix*1000, pageIx*1000)
	templates.S3ObjectList(pattern, objects, pageIx, pageCount+1).Render(r.Context(), w)
}

func GetS3Object(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(ctx)

	key, err := base64.StdEncoding.DecodeString(chi.URLParam(r, "key"))
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

	log.Info("[admin] s3 object loaded", "key", key)

	w.Header().Set("Content-Type", "text/plain")
	w.Write(object)
}
