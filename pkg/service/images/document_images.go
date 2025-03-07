package images

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log/slog"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/graph/model"
	"github.com/fivetentaylor/pointy/pkg/stackerr"
)

const MaxFileSize = 10 * 1024 * 1024

func UploadImage(ctx context.Context, file graphql.Upload, docID string) (img *model.Image, err error) {
	if file.Size > MaxFileSize {
		return img, stackerr.New(fmt.Errorf("file size exceeds the maximum limit of %d bytes", MaxFileSize))
	}

	imageID := uuid.NewString()
	s3 := env.S3(ctx)

	image := &model.Image{
		ID:        imageID,
		DocID:     docID,
		MimeType:  file.ContentType,
		CreatedAt: time.Now(),
		// URL:       fmt.Sprintf("%s/drafts/%s/images/%s", s3.AppHost, docID, imageID),
		URL:    fmt.Sprintf("%s/drafts/%s/images/loading.gif", s3.AppHost, docID),
		Status: model.StatusLoading,
	}

	go func() {
		s3 := env.S3(ctx)
		log := env.SLog(ctx)

		processedImage, mimeType, err := processImage(file.File, file.Filename)
		if err != nil {
			log.Error("error processing image", slog.Any("error", err), slog.String("stack", string(debug.Stack())))
		}

		key := fmt.Sprintf("%s/%s", docID, imageID)
		err = s3.PutObject(s3.ImagesBucket, key, mimeType, processedImage)
		if err != nil {
			log.Error("error uploading image", slog.Any("error", err), slog.String("stack", string(debug.Stack())))
		}
	}()

	return image, nil
}

func GetImage(ctx context.Context, docID string, imageID string) (*model.Image, error) {
	s3 := env.S3(ctx)

	imageKey := fmt.Sprintf("%s/%s", docID, imageID)
	exists, err := s3.Exists(s3.ImagesBucket, imageKey)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/drafts/%s/images/%s", s3.AppHost, docID, imageID)

	if exists {
		return &model.Image{
			ID:        imageID,
			DocID:     docID,
			URL:       url,
			MimeType:  "",
			CreatedAt: time.Now(),
			Status:    model.StatusSuccess,
		}, nil
	} else {
		return &model.Image{
			ID:        imageID,
			DocID:     docID,
			URL:       url,
			MimeType:  "",
			CreatedAt: time.Now(),
			Status:    model.StatusLoading,
		}, nil
	}
}

func ListDocumentImages(ctx context.Context, docID string) ([]*model.Image, error) {
	s3 := env.S3(ctx)
	log := env.SLog(ctx)

	keys, err := s3.List(s3.ImagesBucket, docID, -1, -1)
	if err != nil {
		return nil, err
	}

	images := make([]*model.Image, len(keys))
	for i, key := range keys {
		keyParts := strings.Split(key, "/")
		if len(keyParts) != 2 {
			log.Error("invalid key", slog.String("key", key), slog.String("stack", string(debug.Stack())))
			return nil, fmt.Errorf("internal error")
		}

		images[i] = &model.Image{
			ID:        keyParts[1],
			DocID:     docID,
			URL:       fmt.Sprintf("%s/drafts/%s/images/%s", s3.AppHost, docID, keyParts[1]),
			CreatedAt: time.Now(),
			MimeType:  "",
			Status:    model.StatusSuccess,
		}
	}

	return images, nil
}

func GetImageSignedURL(ctx context.Context, docID string, imageID string) (*model.SignedImageURL, error) {
	s3 := env.S3(ctx)

	url, err := s3.GetPresignedUrl(s3.ImagesBucket, fmt.Sprintf("%s/%s", docID, imageID), 10*time.Minute)
	if err != nil {
		return nil, err
	}

	return &model.SignedImageURL{
		URL:       url,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}, nil
}

const (
	MaxWidth  = 1920
	MaxHeight = 1080
)

func processImage(file io.Reader, filename string) ([]byte, string, error) {
	// Read the entire file into memory
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read file: %v", err)
	}

	// Detect if it's a GIF
	if isGIF(fileBytes) {
		return processGIF(bytes.NewReader(fileBytes))
	}

	// For non-GIF images
	img, imgFormat, err := image.Decode(bytes.NewReader(fileBytes))
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode image: %v", err)
	}

	// Resize the image while maintaining aspect ratio
	resized := imaging.Fit(img, MaxWidth, MaxHeight, imaging.Lanczos)

	var buf bytes.Buffer
	var mimeType string

	// Default to PNG for most cases
	err = png.Encode(&buf, resized)
	mimeType = "image/png"

	// Special handling for JPEG to avoid unnecessary quality loss
	if strings.ToLower(imgFormat) == "jpeg" || strings.ToLower(filepath.Ext(filename)) == ".jpg" || strings.ToLower(filepath.Ext(filename)) == ".jpeg" {
		buf.Reset()
		err = jpeg.Encode(&buf, resized, &jpeg.Options{Quality: 90})
		mimeType = "image/jpeg"
	}

	if err != nil {
		return nil, "", fmt.Errorf("failed to encode image: %v", err)
	}

	return buf.Bytes(), mimeType, nil
}

func isGIF(data []byte) bool {
	return len(data) > 3 && string(data[:3]) == "GIF"
}

func processGIF(file io.Reader) ([]byte, string, error) {
	gifData, err := gif.DecodeAll(file)
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode GIF: %v", err)
	}

	// Check if resizing is needed
	if gifData.Config.Width <= MaxWidth && gifData.Config.Height <= MaxHeight {
		// If the GIF is already small enough, return it as is
		var buf bytes.Buffer
		if err := gif.EncodeAll(&buf, gifData); err != nil {
			return nil, "", fmt.Errorf("failed to re-encode original GIF: %v", err)
		}
		return buf.Bytes(), "image/gif", nil
	}

	// Calculate new dimensions
	newWidth, newHeight := calculateNewDimensions(gifData.Config.Width, gifData.Config.Height, MaxWidth, MaxHeight)

	// Create a new GIF with resized frames
	newGif := &gif.GIF{
		Image:     make([]*image.Paletted, len(gifData.Image)),
		Delay:     make([]int, len(gifData.Image)),
		LoopCount: gifData.LoopCount,
	}

	for i, frame := range gifData.Image {
		// Resize the frame
		resized := imaging.Resize(frame, newWidth, newHeight, imaging.NearestNeighbor)

		// Create a new paletted image with the original palette
		newPaletted := image.NewPaletted(resized.Bounds(), frame.Palette)

		// Draw the resized image onto the new paletted image
		draw.Draw(newPaletted, newPaletted.Bounds(), resized, resized.Bounds().Min, draw.Src)

		newGif.Image[i] = newPaletted
		newGif.Delay[i] = gifData.Delay[i]
	}

	var buf bytes.Buffer
	err = gif.EncodeAll(&buf, newGif)
	if err != nil {
		return nil, "", fmt.Errorf("failed to encode resized GIF: %v", err)
	}

	return buf.Bytes(), "image/gif", nil
}

func calculateNewDimensions(width, height, maxWidth, maxHeight int) (newWidth, newHeight int) {
	aspectRatio := float64(width) / float64(height)

	if width > maxWidth {
		newWidth = maxWidth
		newHeight = int(float64(newWidth) / aspectRatio)
	} else if height > maxHeight {
		newHeight = maxHeight
		newWidth = int(float64(newHeight) * aspectRatio)
	} else {
		return width, height
	}

	if newHeight > maxHeight {
		newHeight = maxHeight
		newWidth = int(float64(newHeight) * aspectRatio)
	}

	return newWidth, newHeight
}
