package heap_test

import (
	"testing"

	heap "github.com/teamreviso/code/rogue/v3/heap"
)

func TestMinHeap(t *testing.T) {
	tests := []struct {
		name     string
		inputs   []int
		expected []int
		expectOk bool
	}{
		{
			name:     "Ascending order test",
			inputs:   []int{3, 2, 1, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
			expectOk: true,
		},
		{
			name:     "Descending order test",
			inputs:   []int{5, 4, 3, 2, 1},
			expected: []int{1, 2, 3, 4, 5},
			expectOk: true,
		},
		{
			name:     "Single element test",
			inputs:   []int{1},
			expected: []int{1},
			expectOk: true,
		},
		{
			name:     "Empty heap test",
			inputs:   []int{},
			expected: nil,
			expectOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := heap.NewMinHeap(heap.Identity, func(a, b int) bool { return a < b })
			for _, num := range tt.inputs {
				h.Push(num)
			}

			for _, exp := range tt.expected {
				min, ok := h.Pop()
				if !ok || min != exp {
					t.Errorf("expected %d, got %d", exp, min)
				}
			}

			// Check if pop operation on empty heap behaves correctly
			_, ok := h.Pop()
			if ok != false {
				t.Errorf("expected empty heap error state to be %v, got %v", false, ok)
			}
		})
	}
}
