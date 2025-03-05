package utils

import (
	"math"
	"strings"
)

func IsBreak(u uint16) bool {
	return u == 0x000D || u == 0x000A || u == 0x0085 || u == 0x2028 || u == 0x2029
}

func AllIndicesOfString(s, substr string) []int {
	var indices []int
	index := 0

	for {
		// Find the next index of substr in s starting from index
		i := strings.Index(s[index:], substr)
		if i == -1 {
			break // No more occurrences found
		}
		// Append the index of the occurrence to the slice
		indices = append(indices, index+i)
		// Move index forward to continue searching
		index += i + 1
	}

	return indices
}

func AllIndicesOfUint16(text, pattern []uint16) []int {
	result := []int{}
	patternLength := len(pattern)
	textLength := len(text)

	for i := 0; i <= textLength-patternLength; i++ {
		substring := text[i : i+patternLength]
		for j := 0; j < patternLength; j++ {
			if substring[j] != pattern[j] {
				break
			}
			if j == patternLength-1 {
				result = append(result, i)
				i += patternLength - 1
			}
		}
	}
	return result
}

// FindSimilarSubstrings finds all indices in 'text' where a substring similar to 'pattern' starts.
// 'maxDistance' specifies the maximum Levenshtein distance allowed.
func FindSimilarSubstrings(text, pattern string, maxDistance int) []int {
	var result []int
	patternLength := len(pattern)
	textLength := len(text)

	for i := 0; i <= textLength-patternLength; i++ {
		substring := text[i : i+patternLength]
		distance := LevenshteinDistanceMax(pattern, substring, maxDistance)
		if distance <= maxDistance {
			result = append(result, i)
		}
	}
	return result
}

// LevenshteinDistanceMax computes the Levenshtein distance between two strings.
// It returns a value greater than 'maxDistance' if the distance exceeds 'maxDistance'.
func LevenshteinDistanceMax(s, t string, maxDistance int) int {
	m := len(s)
	n := len(t)

	// Early exit if length difference exceeds maxDistance
	if int(math.Abs(float64(m-n))) > maxDistance {
		return maxDistance + 1
	}

	// Ensure s is the shorter string
	if m > n {
		s, t = t, s
		m, n = n, m
	}

	previousRow := make([]int, n+1)
	currentRow := make([]int, n+1)

	for j := 0; j <= n; j++ {
		previousRow[j] = j
	}

	for i := 1; i <= m; i++ {
		currentRow[0] = i
		minInRow := currentRow[0]

		for j := 1; j <= n; j++ {
			cost := 0
			if s[i-1] != t[j-1] {
				cost = 1
			}

			currentRow[j] = Min(
				currentRow[j-1]+1,     // Insertion
				previousRow[j]+1,      // Deletion
				previousRow[j-1]+cost, // Substitution
			)

			if currentRow[j] < minInRow {
				minInRow = currentRow[j]
			}
		}

		// Early exit if the minimum distance in the current row exceeds maxDistance
		if minInRow > maxDistance {
			return maxDistance + 1
		}

		previousRow, currentRow = currentRow, previousRow
	}

	if previousRow[n] > maxDistance {
		return maxDistance + 1
	}
	return previousRow[n]
}

// FindSimilarSubstringsUint16 finds all indices in 'text' where a substring similar to 'pattern' starts.
// Both 'text' and 'pattern' are []uint16 slices.
// 'maxDistance' specifies the maximum Levenshtein distance allowed.
func FindSimilarSubstringsUint16(text, pattern []uint16, maxDistance int) []int {
	result := []int{}
	patternLength := len(pattern)
	textLength := len(text)

	if patternLength == 0 || textLength == 0 {
		return result
	}

	for i := 0; i <= textLength-patternLength; i++ {
		substring := text[i : i+patternLength]
		distance := LevenshteinDistanceMaxUint16(pattern, substring, maxDistance)
		if distance <= maxDistance {
			result = append(result, i)
			i += patternLength - 1
		}
	}
	return result
}

// LevenshteinDistanceMaxUint16 computes the Levenshtein distance between two []uint16 slices.
// It returns a value greater than 'maxDistance' if the distance exceeds 'maxDistance'.
func LevenshteinDistanceMaxUint16(s, t []uint16, maxDistance int) int {
	m := len(s)
	n := len(t)

	// Early exit if length difference exceeds maxDistance
	if int(math.Abs(float64(m-n))) > maxDistance {
		return maxDistance + 1
	}

	// Ensure s is the shorter slice
	if m > n {
		s, t = t, s
		m, n = n, m
	}

	previousRow := make([]int, n+1)
	currentRow := make([]int, n+1)

	for j := 0; j <= n; j++ {
		previousRow[j] = j
	}

	for i := 1; i <= m; i++ {
		currentRow[0] = i
		minInRow := currentRow[0]

		for j := 1; j <= n; j++ {
			cost := 0
			if s[i-1] != t[j-1] {
				cost = 1
			}

			currentRow[j] = min(
				currentRow[j-1]+1,     // Insertion
				previousRow[j]+1,      // Deletion
				previousRow[j-1]+cost, // Substitution
			)

			if currentRow[j] < minInRow {
				minInRow = currentRow[j]
			}
		}

		// Early exit if the minimum distance in the current row exceeds maxDistance
		if minInRow > maxDistance {
			return maxDistance + 1
		}

		previousRow, currentRow = currentRow, previousRow
	}

	if previousRow[n] > maxDistance {
		return maxDistance + 1
	}
	return previousRow[n]
}
