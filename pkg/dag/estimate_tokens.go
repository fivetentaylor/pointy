package dag

import (
	"strings"
	"unicode"
)

func estimateTokens(s string) int {
	// Split the string into words
	words := strings.FieldsFunc(s, func(r rune) bool {
		return unicode.IsSpace(r) || unicode.IsPunct(r)
	})

	// Count the number of words
	wordCount := len(words)

	// Estimate the number of tokens
	// This is a rough estimate: assume 1.3 tokens per word on average
	estimatedTokens := int(float64(wordCount) * 1.3)

	return estimatedTokens
}
