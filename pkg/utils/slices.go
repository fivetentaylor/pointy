package utils

func SafeSlice[T any](data []T, start, end int) []T {
	length := len(data)

	// Adjust start and end indices to be within slice bounds
	if start < 0 {
		start = length + start + 1
	}
	if start < 0 {
		start = 0
	}
	if start > length {
		start = length
	}

	if end < 0 {
		end = length + end + 1
	}
	if end < 0 {
		end = 0
	}
	if end > length {
		end = length
	}
	if start > end {
		start = end
	}

	return data[start:end]
}

func Reverse[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func MakeCopy[T any](s []T) []T {
	copied := make([]T, len(s))
	copy(copied, s)
	return copied
}

func Difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}
