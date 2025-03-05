package heap

// MinHeap is a generic min-heap (priority queue) where T is the type of items in the heap, and K is the type of the key used for ordering.
type MinHeap[T any, K any] struct {
	Size     int
	items    []T
	keyFunc  func(T) K
	lessThan func(K, K) bool
}

// NewMinHeap creates a new MinHeap with the provided key extraction function.
func NewMinHeap[T any, K any](keyFunc func(T) K, lessThan func(K, K) bool) *MinHeap[T, K] {
	return &MinHeap[T, K]{
		items:    []T{},
		keyFunc:  keyFunc,
		lessThan: lessThan,
	}
}

// Push adds an element to the heap
func (h *MinHeap[T, K]) Push(item T) {
	h.items = append(h.items, item)
	h.up(len(h.items) - 1)
	h.Size++
}

// Pop removes and returns the smallest element from the heap
func (h *MinHeap[T, K]) Pop() (T, bool) {
	if len(h.items) == 0 {
		var zero T
		return zero, false
	}
	min := h.items[0]
	h.items[0] = h.items[len(h.items)-1]
	h.items = h.items[:len(h.items)-1]
	h.down(0)
	h.Size--
	return min, true
}

// up moves the element at index i up to its correct position in the heap
func (h *MinHeap[T, K]) up(i int) {
	for {
		parent := (i - 1) / 2
		k0 := h.keyFunc(h.items[parent])
		k1 := h.keyFunc(h.items[i])
		if i == 0 || h.lessThan(k0, k1) || !h.lessThan(k1, k0) {
			break
		}
		h.items[i], h.items[parent] = h.items[parent], h.items[i]
		i = parent
	}
}

// down moves the element at index i down to its correct position in the heap
func (h *MinHeap[T, K]) down(i int) {
	for {
		leftChild := 2*i + 1
		rightChild := 2*i + 2
		smallest := i

		if leftChild < len(h.items) {
			kSmallest := h.keyFunc(h.items[smallest])
			kLeft := h.keyFunc(h.items[leftChild])
			if h.lessThan(kLeft, kSmallest) {
				smallest = leftChild
			}
		}

		if rightChild < len(h.items) {
			kSmallest := h.keyFunc(h.items[smallest])
			kRight := h.keyFunc(h.items[rightChild])
			if h.lessThan(kRight, kSmallest) {
				smallest = rightChild
			}
		}

		if smallest == i {
			break
		}

		h.items[i], h.items[smallest] = h.items[smallest], h.items[i]
		i = smallest
	}
}

// MaxHeap is a generic max-heap (priority queue) where T is the type of items in the heap, and K is the type of the key used for ordering.
type MaxHeap[T any, K any] struct {
	Size     int
	items    []T
	keyFunc  func(T) K
	lessThan func(K, K) bool
}

// NewMaxHeap creates a new MaxHeap with the provided key extraction function.
func NewMaxHeap[T any, K any](keyFunc func(T) K, lessThan func(K, K) bool) *MaxHeap[T, K] {
	return &MaxHeap[T, K]{
		items:    []T{},
		keyFunc:  keyFunc,
		lessThan: lessThan,
	}
}

// Push adds an element to the heap
func (h *MaxHeap[T, K]) Push(item T) {
	h.items = append(h.items, item)
	h.up(len(h.items) - 1)
	h.Size++
}

// Pop removes and returns the largest element from the heap
func (h *MaxHeap[T, K]) Pop() (T, bool) {
	if len(h.items) == 0 {
		var zero T
		return zero, false
	}
	max := h.items[0]
	h.items[0] = h.items[len(h.items)-1]
	h.items = h.items[:len(h.items)-1]
	h.down(0)
	h.Size--
	return max, true
}

// up moves the element at index i up to its correct position in the heap
func (h *MaxHeap[T, K]) up(i int) {
	for {
		parent := (i - 1) / 2
		k0 := h.keyFunc(h.items[i])
		k1 := h.keyFunc(h.items[parent])

		if i == 0 || h.lessThan(k0, k1) || !h.lessThan(k1, k0) {
			break
		}
		h.items[i], h.items[parent] = h.items[parent], h.items[i]
		i = parent
	}
}

// down moves the element at index i down to its correct position in the heap
func (h *MaxHeap[T, K]) down(i int) {
	for {
		leftChild := 2*i + 1
		rightChild := 2*i + 2
		largest := i

		if leftChild < len(h.items) {
			kLargest := h.keyFunc(h.items[largest])
			kLeft := h.keyFunc(h.items[leftChild])
			if h.lessThan(kLargest, kLeft) {
				largest = leftChild
			}
		}
		if rightChild < len(h.items) {
			kLargest := h.keyFunc(h.items[largest])
			kRight := h.keyFunc(h.items[rightChild])
			if h.lessThan(kLargest, kRight) {
				largest = rightChild
			}
		}
		if largest == i {
			break
		}
		h.items[i], h.items[largest] = h.items[largest], h.items[i]
		i = largest
	}
}

func Identity[T any](x T) T {
	return x
}
