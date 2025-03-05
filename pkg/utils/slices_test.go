package utils

import (
	"reflect"
	"testing"
)

func TestSafeSlice(t *testing.T) {
	tests := []struct {
		name   string
		input  []int
		start  int
		end    int
		output []int
	}{
		{"Normal slice", []int{1, 2, 3, 4, 5}, 1, 3, []int{2, 3}},
		{"End before start", []int{1, 2, 3, 4, 5}, 3, 1, []int{}},
		{"Negative start", []int{1, 2, 3, 4, 5}, -5, 3, []int{2, 3}},
		{"Negative end", []int{1, 2, 3, 4, 5}, 1, -2, []int{2, 3, 4}},
		{"Both negative", []int{1, 2, 3, 4, 5}, -2, -3, []int{}},
		{"Both out of bounds", []int{1, 2, 3, 4, 5}, 7, 10, []int{}},
		{"Empty slice", []int{}, 0, 2, []int{}},
		{"Start equals end", []int{1, 2, 3, 4, 5}, 2, 2, []int{}},
		{"Start and end out of bounds", []int{1, 2, 3, 4, 5}, -10, 10, []int{1, 2, 3, 4, 5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SafeSlice(tt.input, tt.start, tt.end)
			if !reflect.DeepEqual(got, tt.output) {
				t.Errorf("SafeSlice(%v, %d, %d) = %v; want %v", tt.input, tt.start, tt.end, got, tt.output)
			}
		})
	}
}

func TestReverse(t *testing.T) {
	intTests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{"IntsEvenLength", []int{1, 2, 3, 4}, []int{4, 3, 2, 1}},
		{"IntsOddLength", []int{1, 2, 3, 4, 5}, []int{5, 4, 3, 2, 1}},
		{"EmptyInts", []int{}, []int{}},
		{"SingleInt", []int{1}, []int{1}},
	}

	for _, tt := range intTests {
		t.Run(tt.name, func(t *testing.T) {
			// We need reflection to call the generic function with various types dynamically
			Reverse(tt.input)

			if !reflect.DeepEqual(tt.input, tt.expected) {
				t.Errorf("got %v, want %v", tt.input, tt.expected)
			}
		})
	}

	stringTests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{"Strings", []string{"a", "b", "c"}, []string{"c", "b", "a"}},
		{"EmptyStrings", []string{}, []string{}},
		{"SingleString", []string{"a"}, []string{"a"}},
	}

	for _, tt := range stringTests {
		t.Run(tt.name, func(t *testing.T) {
			// We need reflection to call the generic function with various types dynamically
			Reverse(tt.input)

			if !reflect.DeepEqual(tt.input, tt.expected) {
				t.Errorf("got %v, want %v", tt.input, tt.expected)
			}
		})
	}
}
