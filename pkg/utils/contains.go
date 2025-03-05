package utils

// Contains checks if the slice s contains all of the elements in es.
func Contains[T comparable](s []T, es ...T) bool {
	if len(es) == 0 {
		return false
	}
	for _, e := range es {
		found := false
		for _, v := range s {
			if v == e {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
