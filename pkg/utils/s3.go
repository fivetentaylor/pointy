package utils

import (
	"fmt"

	"github.com/fivetentaylor/pointy/pkg/constants"
)

func InvertSeq(seq int64) string {
	invertedSeq := constants.MaxSeqValue - seq
	return fmt.Sprintf("%016d", invertedSeq)
}
