package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fivetentaylor/pointy/pkg/env"
)

func saveLogFile(ctx context.Context, id, filename string, content any) error {
	s3 := env.S3(ctx)

	key := fmt.Sprintf("convos/%s/%s", id, filename)

	bts, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling content: %s", err)
	}

	return s3.PutObject(s3.Bucket, key, "application/json", bts)
}
