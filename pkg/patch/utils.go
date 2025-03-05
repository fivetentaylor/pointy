package patch

import (
	"fmt"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

var sentenceEnders = map[rune]struct{}{
	'.':  {},
	'!':  {},
	'?':  {},
	'\n': {},
}

func SentenceBounds(s string, index int) (int, int) {
	if index < 0 {
		index = 0
	} else if index > len(s)-1 {
		index = max(0, len(s)-1)
	}

	// Convert byte index to rune index
	runeIndex := len([]rune(s[:index]))

	startIndex, endIndex := runeIndex, runeIndex

	// Convert string to slice of runes for easier indexing
	runes := []rune(s)

	// Scan backward
	for startIndex > 0 {
		if _, exists := sentenceEnders[runes[startIndex]]; exists {
			startIndex++ // move past the delimiter to exclude it in the result
			break
		}
		startIndex--
	}

	// Scan forward
	for endIndex < len(runes)-1 {
		if _, exists := sentenceEnders[runes[endIndex]]; exists {
			break
		}
		endIndex++
	}

	// Convert rune indices back to byte indices
	byteStart := len(string(runes[:startIndex]))
	byteEnd := len(string(runes[:endIndex])) // Not including endIndex

	return byteStart, byteEnd
}

func all(values []bool) bool {
	for _, v := range values {
		if !v {
			return false
		}
	}
	return true
}

func Apply(target, patch string) (string, error) {
	dmp := diffmatchpatch.New()
	patches, err := dmp.PatchFromText(patch)
	if err != nil {
		return "", fmt.Errorf("failed to build patch: %+v", err)
	}

	newContent, patchResults := dmp.PatchApply(patches, target)
	if !all(patchResults) {
		return "", fmt.Errorf("failed to apply patch: %+v", target)
	}

	return newContent, nil
}

func StartIndices(patch string) (int, int) {
	dmp := diffmatchpatch.New()

	patches, _ := dmp.PatchFromText(patch)
	if len(patches) > 0 {
		return patches[0].Start1, patches[0].Start2
	}

	return -1, -1
}
