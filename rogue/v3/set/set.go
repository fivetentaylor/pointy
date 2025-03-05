package set

import "encoding/json"

// Set type using a map with a generic type T
type Set[T comparable] map[T]struct{}

func NewSet[T comparable](elements ...T) Set[T] {
	s := Set[T]{}
	for _, e := range elements {
		s[e] = struct{}{}
	}
	return s
}

// Add an element to the set
func (s Set[T]) Add(item T) {
	s[item] = struct{}{}
}

// Remove an element from the set
func (s Set[T]) Pop(item T) bool {
	has := s.Has(item)
	delete(s, item)
	return has
}

// Check if an element is in the set
func (s Set[T]) Has(item T) bool {
	_, ok := s[item]
	return ok
}

func (s Set[T]) Clear() {
	for k := range s {
		delete(s, k)
	}
}

// Size of the set
func (s Set[T]) Size() int {
	return len(s)
}

func (s Set[T]) And(other Set[T]) Set[T] {
	out := Set[T]{}
	for k := range s {
		if other.Has(k) {
			out.Add(k)
		}
	}
	return out
}

func (s Set[T]) Or(other Set[T]) Set[T] {
	out := Set[T]{}
	for k := range s {
		out.Add(k)
	}
	for k := range other {
		out.Add(k)
	}
	return out
}

func (s Set[T]) Xor(other Set[T]) Set[T] {
	out := Set[T]{}
	for k := range s {
		if !other.Has(k) {
			out.Add(k)
		}
	}
	for k := range other {
		if !s.Has(k) {
			out.Add(k)
		}
	}
	return out
}

func (s Set[T]) Covers(other Set[T]) bool {
	for k := range other {
		if !s.Has(k) {
			return false
		}
	}
	return true
}

func (s Set[T]) Keys() []T {
	out := make([]T, 0, len(s))
	for k := range s {
		out = append(out, k)
	}
	return out
}

func (s Set[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Keys())
}

func (s *Set[T]) UnmarshalJSON(data []byte) error {
	var keys []T
	if err := json.Unmarshal(data, &keys); err != nil {
		return err
	}
	for _, k := range keys {
		s.Add(k)
	}
	return nil
}
