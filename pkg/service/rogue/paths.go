package rogue

import (
	"fmt"
	"path"

	"github.com/fivetentaylor/pointy/pkg/constants"
)

func SeqToS3Path(docID string, seq int64) string {
	return path.Join(DocS3Path(docID), invertSeq(seq))
}

func invertSeq(seq int64) string {
	invertedSeq := constants.MaxSeqValue - seq
	return fmt.Sprintf("%016d", invertedSeq)
}

func DocS3Path(docID string) string {
	return fmt.Sprintf(constants.DocumentSnapshotPrefix, constants.S3Prefix, docID)
}
