package attachments

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"code.sajari.com/docconv/v2"

	goose "github.com/advancedlogic/GoOse"
	"github.com/fivetentaylor/pointy/pkg/constants"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/models"
)

const ExtractedTextFilename = "extracted.txt"

var SupportedMimeTypes = []string{
	"text/markdown", "text/csv",
	"application/msword", "application/vnd.ms-word",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	"application/vnd.openxmlformats-officedocument.presentationml.presentation",
	"application/vnd.oasis.opendocument.text",
	// "application/vnd.apple.pages", "application/x-iwork-pages-sffpages",
	"application/pdf",
	"application/rtf", "application/x-rtf", "text/rtf", "text/richtext",
	"text/html",
	"text/url",
	"text/xml", "application/xml",
	// "image/jpeg", "image/png", "image/tif", "image/tiff",
	"text/plain",
}

func ExtractedText(ctx context.Context, attachment *models.DocumentAttachment) (string, error) {
	log := env.SLog(ctx)
	if existing := checkForExtractedText(ctx, attachment); existing != nil {
		log.Info("extracted text already exists for attachment", "attachment_id", attachment.ID)
		return *existing, nil
	}

	var text string
	var err error
	switch attachment.ContentType {
	case "text/markdown":
		text, err = extractedTextFromMarkdown(ctx, attachment)
	case "text/csv":
		text, err = extractedTextFromCSV(ctx, attachment)
	case "text/url":
		text, err = extractedTextFromURL(ctx, attachment)
	default:
		reader, err := originalReader(ctx, attachment)
		if err != nil {
			log.Error("error getting original reader", "error", err)
			return "", fmt.Errorf("error getting original reader: %w", err)
		}

		res, err := docconv.Convert(reader, attachment.ContentType, true)
		if err != nil {
			log.Error("error extracting text from attachment", "error", err)
			return "", fmt.Errorf("error extracting text with docconv: %w", err)
		}

		text = res.Body
	}
	if err != nil {
		log.Error("error extracting text from attachment", "error", err, "attachment_id", attachment.ID)
		return "", err
	}

	if strings.TrimSpace(text) == "" {
		log.Info("no text extracted from attachment", "attachment_id", attachment.ID)
		return "", fmt.Errorf("no text extracted from attachment")
	}

	go saveExtractedText(ctx, attachment, text)

	return text, nil
}

func saveExtractedText(ctx context.Context, attachment *models.DocumentAttachment, text string) {
	log := env.SLog(ctx)
	s3 := env.S3(ctx)

	outKey := fmt.Sprintf(constants.DocumentAttachmentFileKey, attachment.S3ID, ExtractedTextFilename)
	err := s3.PutObject(s3.Bucket, outKey, "text/plain", []byte(text))
	if err != nil {
		log.Error("error putting extracted text", "error", err, "key", outKey)
	}
}

func extractedTextFromURL(ctx context.Context, attachment *models.DocumentAttachment) (string, error) {
	input, err := originalString(ctx, attachment)
	if err != nil {
		return "", err
	}

	// Original URL
	g := goose.New()
	article, err := g.ExtractFromURL(input)
	if err != nil {
		return extractedTextFromScrapingbeeURL(ctx, attachment, err)
	}

	return article.CleanedText, nil
}

func extractedTextFromScrapingbeeURL(ctx context.Context, attachment *models.DocumentAttachment, originalError error) (string, error) {
	log := env.SLog(ctx)
	apiKey := os.Getenv("SCRAPINGBEE_API_KEY")
	if apiKey == "" {
		return "", originalError
	}

	input, err := originalString(ctx, attachment)
	if err != nil {
		return "", err
	}
	url := strings.TrimSpace(input)
	wrappedUrl := fmt.Sprintf("https://app.scrapingbee.com/api/v1/?api_key=%s&url=%s", apiKey, url)

	log.Info("extracting text from scrapingbee url", "url", url, "wrapped_url", wrappedUrl)
	// Original URL
	g := goose.New()
	article, err := g.ExtractFromURL(wrappedUrl)
	if err != nil {
		log.Error("error extracting text from scrapingbee url", "error", err, "wrapped_url", wrappedUrl)
		return "", err
	}

	log.Info("extracted text from scrapingbee url", "url", url, "wrapped_url", wrappedUrl, "text", article.CleanedText)
	return article.CleanedText, nil
}

func extractedTextFromMarkdown(ctx context.Context, attachment *models.DocumentAttachment) (string, error) {
	bts, err := originalBytes(ctx, attachment)
	return string(bts), err
}

func extractedTextFromCSV(ctx context.Context, attachment *models.DocumentAttachment) (string, error) {
	log := env.SLog(ctx)
	r, err := originalReader(ctx, attachment)
	if err != nil {
		log.Error("error getting original bytes", "error", err)
		return "", fmt.Errorf("error getting original bytes: %w", err)
	}

	// Create CSV reader
	reader := csv.NewReader(r)

	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		log.Error("failed to parse CSV", "error", err)
		return "", fmt.Errorf("failed to parse CSV: %w", err)
	}

	var builder strings.Builder
	for _, row := range records {
		builder.WriteString(strings.Join(row, " "))
		builder.WriteString("\n")
	}

	return builder.String(), nil
}

func checkForExtractedText(ctx context.Context, attachment *models.DocumentAttachment) *string {
	log := env.SLog(ctx)
	s3 := env.S3(ctx)
	key := fmt.Sprintf(constants.DocumentAttachmentFileKey, attachment.S3ID, ExtractedTextFilename)
	exists, err := s3.ObjectExists(s3.Bucket, key)
	if err != nil {
		log.Error("error checking if extracted text exists", "error", err, "key", key)
		return nil
	}

	if exists {
		extractedText, err := s3.GetObject(s3.Bucket, key)
		if err != nil {
			log.Error("error getting extracted text existing key", "error", err, "key", key)
			return nil
		}
		extractedTextString := string(extractedText)

		return &extractedTextString
	}

	return nil
}

func originalString(ctx context.Context, attachment *models.DocumentAttachment) (string, error) {
	bytes, err := originalBytes(ctx, attachment)
	if err != nil {
		return "", fmt.Errorf("error getting original string: %w", err)
	}
	return string(bytes), nil
}

func originalBytes(ctx context.Context, attachment *models.DocumentAttachment) ([]byte, error) {
	log := env.SLog(ctx)
	s3 := env.S3(ctx)
	originalKey := fmt.Sprintf(constants.DocumentAttachmentOriginalKey, attachment.S3ID)
	originalBytes, err := s3.GetObject(s3.Bucket, originalKey)
	if err != nil {
		log.Error("error getting original bytes", "error", err)
		return nil, fmt.Errorf("error getting original bytes: %w", err)
	}
	return originalBytes, nil
}

func originalReader(ctx context.Context, attachment *models.DocumentAttachment) (io.Reader, error) {
	log := env.SLog(ctx)
	originalBytes, err := originalBytes(ctx, attachment)
	if err != nil {
		log.Error("error getting original bytes", "error", err)
		return nil, fmt.Errorf("error getting original bytes: %w", err)
	}

	return bytes.NewReader(originalBytes), nil
}
