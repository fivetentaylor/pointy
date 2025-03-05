package utils

import (
	"sync"
)

// ThreadSafeSlice is a generic type that holds a slice of any type and a mutex for thread-safe operations.
type ThreadSafeSlice[T any] struct {
	sync.Mutex
	Objects []T
}

// NewThreadSafeSlice creates a new instance of a thread-safe slice.
func NewThreadSafeSlice[T any]() *ThreadSafeSlice[T] {
	return &ThreadSafeSlice[T]{}
}

// Add appends an object to the slice in a thread-safe manner.
func (t *ThreadSafeSlice[T]) Add(object T) {
	t.Lock()
	defer t.Unlock()
	t.Objects = append(t.Objects, object)
}

// Length returns the length of the slice in a thread-safe manner.
func (t *ThreadSafeSlice[T]) Length() int {
	t.Lock()
	defer t.Unlock()
	return len(t.Objects)
}

// Slice returns a slice of the objects from start to end indices in a thread-safe manner.
func (t *ThreadSafeSlice[T]) Slice(start, end int) []T {
	t.Lock()
	defer t.Unlock()
	if start < 0 {
		start = 0
	}
	if end > len(t.Objects) {
		end = len(t.Objects)
	}
	if start > end {
		return []T{}
	}
	sliceCopy := make([]T, end-start)
	copy(sliceCopy, t.Objects[start:end])
	return sliceCopy
}

// Select returns a slice of objects that match the predicate in a thread-safe manner.
func (t *ThreadSafeSlice[T]) Select(predicate func(T) (bool, error)) ([]T, error) {
	t.Lock()
	defer t.Unlock()
	var selected []T
	for _, obj := range t.Objects {
		match, err := predicate(obj)
		if err != nil {
			return nil, err // Return immediately if an error is encountered
		}
		if match {
			selected = append(selected, obj)
		}
	}
	return selected, nil
}
