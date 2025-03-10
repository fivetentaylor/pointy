package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/fivetentaylor/pointy/pkg/constants"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/query"
	"github.com/fivetentaylor/pointy/pkg/service/images"
)

const maxImageSize = 5 * 1024 * 1024
const docWidth = 585

// CreateDocumentImage
// note: this hasn't been implemented with the frontend yet
func (s *Server) CreateDocumentImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(ctx)
	q := env.Query(ctx)

	docID := chi.URLParam(r, "docID")

	var userID string
	currentUser, err := env.UserClaim(ctx)
	if err == nil {
		userID = currentUser.Id
	}

	doc, err := query.GetReadableDocumentForUser(q, docID, userID)
	if doc == nil || err != nil {
		if err != nil {
			log.Errorf("error getting document: %s", err)
		}
		http.NotFound(w, r)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error uploading file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileType := handler.Header.Get("Content-Type")
	if fileType != "image/jpeg" && fileType != "image/png" {
		log.Error("File must be an image")
		http.Error(w, "File must be an image", http.StatusBadRequest)
		return
	}

	if handler.Size > maxImageSize {
		log.Error(fmt.Sprintf("File must be smaller than %d bytes", maxImageSize))
		http.Error(w, "File must be smaller than 5MB", http.StatusBadRequest)
		return
	}

	img, _, err := image.Decode(file)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to decode image: %s", err.Error()))
		http.Error(w, "Failed to decode image", http.StatusInternalServerError)
		return
	}

	scaledImg := images.FitImageToWidth(img, docWidth)

	var outputBuffer bytes.Buffer
	err = png.Encode(&outputBuffer, scaledImg)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to save resized image: %s", err.Error()))
		http.Error(w, "Failed to save resized image", http.StatusInternalServerError)
		return
	}

	imageId := uuid.NewString()
	imageKey := fmt.Sprintf(constants.DocumentImageKeyFormat, docID, imageId)
	err = s.S3.PutObject(s.S3.Bucket, imageKey, "image/png", outputBuffer.Bytes())
	if err != nil {
		log.Error(fmt.Sprintf("Failed to save scalled image to S3: %s", err.Error()))
		http.Error(w, "Failed to save image", http.StatusInternalServerError)
		return
	}

	originalKey := fmt.Sprintf(constants.OriginalDocumentImageKeyFormat, docID, imageId)
	err = s.S3.PutObjectWithSeeker(s.S3.Bucket, originalKey, fileType, file)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to save scalled image to S3: %s", err.Error()))
		http.Error(w, "Failed to save image", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(map[string]string{
		"imageId": imageId,
	})
	if err != nil {
		log.Error(fmt.Sprintf("Failed to marshal response: %s", err.Error()))
		http.Error(w, "Failed to build response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (s *Server) GetUserAvatar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(ctx)
	q := env.Query(ctx)

	userID := chi.URLParam(r, "userID")
	user, err := q.User.Where(q.User.ID.Eq(userID)).First()
	if err != nil {
		log.Error(fmt.Sprintf("Error getting user: %s", err.Error()))
		http.Error(w, "Error getting user", http.StatusUnauthorized)
		return
	}

	avatarUrl, err := images.AvatarUrlForUser(ctx, user)
	if err != nil {
		log.Error(fmt.Sprintf("Error getting avatar url: %s", err.Error()))
		http.Error(w, "Error getting avatar url", http.StatusUnauthorized)
		return
	}
	if avatarUrl == nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, *avatarUrl, http.StatusTemporaryRedirect)
}

func (s *Server) UpdateUserAvatar(w http.ResponseWriter, r *http.Request) {
	currentUser, err := env.UserClaim(r.Context())
	if err != nil {
		log.Error(fmt.Sprintf("Error getting user claim: %s", err.Error()))
		http.Error(w, "Error getting user claim", http.StatusUnauthorized)
		return
	}

	userTbl := env.Query(r.Context()).User

	user, err := userTbl.Where(userTbl.ID.Eq(currentUser.Id)).First()
	if err != nil {
		log.Error(fmt.Sprintf("Error getting user: %s", err.Error()))
		http.Error(w, "Error getting user", http.StatusUnauthorized)
		return
	}

	// Get file from form
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Error(fmt.Sprintf("Error getting file: %s", err.Error()))
		http.Error(w, "Error uploading file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileType := handler.Header.Get("Content-Type")
	if fileType != "image/jpeg" && fileType != "image/png" {
		log.Error("File must be an image")
		http.Error(w, "File must be an image", http.StatusBadRequest)
		return
	}

	if handler.Size > maxImageSize {
		log.Error(fmt.Sprintf("File must be smaller than %d bytes", maxImageSize))
		http.Error(w, "File must be smaller than 5MB", http.StatusBadRequest)
		return
	}

	img, _, err := image.Decode(file)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to decode image: %s", err.Error()))
		http.Error(w, "Failed to decode image", http.StatusInternalServerError)
		return
	}

	err = images.UpdateUserAvatar(r.Context(), user, img)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to update user avatar: %s", err.Error()))
		http.Error(w, "Failed to update user avatar", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) GetDocumentImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := env.Log(ctx)
	q := env.Query(ctx)

	docID := chi.URLParam(r, "docID")
	imageID := chi.URLParam(r, "imageID")

	var userID string
	currentUser, err := env.UserClaim(ctx)
	if err == nil {
		userID = currentUser.Id
	}

	doc, err := query.GetReadableDocumentForUser(q, docID, userID)
	if doc == nil || err != nil {
		if err != nil {
			log.Errorf("error getting document: %s", err)
		}
		http.NotFound(w, r)
		return
	}

	if imageID == "loading.gif" {
		signedURL, err := images.GetImageSignedURL(ctx, "default", "loading.gif")
		if err != nil {
			log.Errorf("error getting signed url: %s", err)
			http.Error(w, "Error getting signed url", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, signedURL.URL, http.StatusTemporaryRedirect)
	}

	// get signed link and redirect
	signedURL, err := images.GetImageSignedURL(ctx, docID, imageID)
	if err != nil {
		log.Errorf("error getting signed url: %s", err)
		http.Error(w, "Error getting signed url", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, signedURL.URL, http.StatusTemporaryRedirect)
}
