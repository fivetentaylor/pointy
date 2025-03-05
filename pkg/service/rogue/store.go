package rogue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/teamreviso/code/pkg/env"
	v3 "github.com/teamreviso/code/rogue/v3"
)

// SaveDocToS3 saves the doc to S3, it shouldn't normally be called directly. But is used directly in tests.
func SaveDocToS3(ctx context.Context, docID string, seq int64, doc *v3.Rogue) error {
	docBytes, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("error marshalling doc: %s", err)
	}

	s3path := SeqToS3Path(docID, seq)
	log.Info("Saving doc to S3", "s3path", s3path)
	s3 := env.S3(ctx)
	err = s3.PutObject(s3.Bucket, s3path, "text/plain", docBytes)
	if err != nil {
		return fmt.Errorf("error saving doc to S3: %s", err)
	}

	return nil
}
