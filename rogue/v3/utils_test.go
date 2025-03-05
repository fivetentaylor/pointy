package v3

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/constraints"
)

func TestBisectLeft(t *testing.T) {
	type testCase[T any, K constraints.Ordered] struct {
		name         string
		array        []T
		key          K
		keyExtractor func(T) K
		expected     int
	}

	// Define your test cases
	testCases := []testCase[int, int]{
		{
			name:         "Empty array",
			array:        []int{},
			key:          5,
			keyExtractor: func(a int) int { return a },
			expected:     0,
		},
		{
			name:         "Single element less than key",
			array:        []int{1},
			key:          5,
			keyExtractor: func(a int) int { return a },
			expected:     1,
		},
		{
			name:         "Single element equal to key",
			array:        []int{5},
			key:          5,
			keyExtractor: func(a int) int { return a },
			expected:     0,
		},
		{
			name:         "Multiple elements",
			array:        []int{1, 3, 3, 5, 7},
			key:          4,
			keyExtractor: func(a int) int { return a },
			expected:     3,
		},
		{
			name:         "Multiple matching elements",
			array:        []int{1, 3, 4, 4, 4, 7},
			key:          4,
			keyExtractor: func(a int) int { return a },
			expected:     2,
		},
	}

	// Execute each test case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := bisectLeft(tc.array, tc.key, tc.keyExtractor)
			if result != tc.expected {
				t.Errorf("Test '%s' failed: expected %d, got %d", tc.name, tc.expected, result)
			}
		})
	}
}

func TestSortedInsert(t *testing.T) {
	compareInts := func(a, b int) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	}

	// Test case structure
	type testCase struct {
		name     string
		slice    []int
		element  int
		expected []int
		index    int
	}

	// Define test cases
	testCases := []testCase{
		{
			name:     "Insert into empty slice",
			slice:    []int{},
			element:  5,
			expected: []int{5},
			index:    0,
		},
		{
			name:     "Insert at the beginning",
			slice:    []int{2, 3, 4},
			element:  1,
			expected: []int{1, 2, 3, 4},
			index:    0,
		},
		{
			name:     "Insert in the middle",
			slice:    []int{1, 3, 4},
			element:  2,
			expected: []int{1, 2, 3, 4},
			index:    1,
		},
		{
			name:     "Insert at the end",
			slice:    []int{1, 2, 3},
			element:  4,
			expected: []int{1, 2, 3, 4},
			index:    3,
		},
		// Add more test cases as necessary
	}

	// Execute each test case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			originalSlice := make([]int, len(tc.slice))
			copy(originalSlice, tc.slice)

			index, slice := sortedInsert(tc.slice, tc.element, compareInts)

			if index != tc.index || !reflect.DeepEqual(slice, tc.expected) {
				t.Errorf("Test '%s' failed: expected index %d and slice %v, got index %d and slice %v",
					tc.name, tc.index, tc.expected, index, slice)
			}
		})
	}
}

func TestRandRunes(t *testing.T) {
	r := rand.New(rand.NewSource(1))

	s := RandomString(r, 10000)
	us := StrToUint16(s)

	for i := 0; i < len(us); i++ {
		if IsHighSurrogate(us[i]) {
			assert.True(t, IsLowSurrogate(us[i+1]))
		} else if IsLowSurrogate(us[i]) {
			assert.True(t, IsHighSurrogate(us[i-1]))
		}
	}
}

func TestUTF16Length(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"ASCII", "Hello", 5},
		{"Chinese characters", "你好", 2},
		{"Emoji single code point", "😀", 2},
		{"Emoji ZWJ sequence", "👩‍❤️‍💋‍👩", 8},
		{"Mixed", "Hello, 世界!", 9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, r := range tt.input {
				fmt.Printf("%q\n", r)
			}
			got := UTF16Length(tt.input)
			// if got != tt.want {
			if got != len(StrToUint16(tt.input)) {
				t.Errorf("utf16CodePointsLength(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestAbs(t *testing.T) {
	tests := []struct {
		name     string
		input    any // Use 'any' since we're dealing with multiple types
		expected any
	}{
		{"int negative", -5, 5},
		{"int positive", 5, 5},
		{"int8 negative", int8(-5), int8(5)},
		{"int16 negative", int16(-20), int16(20)},
		{"int32 negative", int32(-100), int32(100)},
		{"int64 negative", int64(-200), int64(200)},
		{"float32 negative", float32(-1.23), float32(1.23)},
		{"float64 negative", float64(-4.56), float64(4.56)},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Since we're using 'any', we need to switch on type for each case to call 'abs' correctly
			switch v := tt.input.(type) {
			case int:
				if got := abs(v); got != tt.expected {
					t.Errorf("abs() = %v, want %v", got, tt.expected)
				}
			case int8:
				if got := abs(v); got != tt.expected {
					t.Errorf("abs() = %v, want %v", got, tt.expected)
				}
			case int16:
				if got := abs(v); got != tt.expected {
					t.Errorf("abs() = %v, want %v", got, tt.expected)
				}
			case int32:
				if got := abs(v); got != tt.expected {
					t.Errorf("abs() = %v, want %v", got, tt.expected)
				}
			case int64:
				if got := abs(v); got != tt.expected {
					t.Errorf("abs() = %v, want %v", got, tt.expected)
				}
			case float32:
				if got := abs(v); got != tt.expected {
					t.Errorf("abs() = %v, want %v", got, tt.expected)
				}
			case float64:
				if got := abs(v); got != tt.expected {
					t.Errorf("abs() = %v, want %v", got, tt.expected)
				}
			default:
				t.Errorf("Unsupported type %T", tt.input)
			}
		})
	}
}

// TestIsValidUTF16 runs table-driven tests on the IsValidUTF16 function
func TestIsValidUTF16(t *testing.T) {
	tests := []struct {
		name     string
		sequence []uint16
		isValid  bool
	}{
		{
			name:     "Valid sequence with no surrogates",
			sequence: []uint16{0x0061, 0x0062, 0x0063}, // "abc"
			isValid:  true,
		},
		{
			name:     "Valid sequence with a correctly formed surrogate pair",
			sequence: []uint16{0xD800, 0xDC00}, // A valid surrogate pair
			isValid:  true,
		},
		{
			name:     "Invalid sequence with high surrogate not followed by low",
			sequence: []uint16{0xD800, 0x0061}, // High surrogate followed by 'a'
			isValid:  false,
		},
		{
			name:     "Invalid sequence with low surrogate before high",
			sequence: []uint16{0xDC00, 0xD800}, // Low surrogate before high surrogate
			isValid:  false,
		},
		{
			name:     "Valid sequence with multiple surrogate pairs",
			sequence: []uint16{0xD800, 0xDC00, 0xD801, 0xDC01}, // Two valid surrogate pairs
			isValid:  true,
		},
		{
			name:     "Invalid sequence with high surrogate at end",
			sequence: []uint16{0x0061, 0xD800}, // 'a' followed by a lone high surrogate
			isValid:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.isValid, IsValidUTF16(tc.sequence))
		})
	}
}

func TestInterleave(t *testing.T) {
	rng := rand.New(rand.NewSource(42))

	tests := []struct {
		name     string
		slices   [][]int
		seed     int64
		expected []int
	}{
		{
			name:   "single slice",
			slices: [][]int{{1, 2, 3}},
		},
		{
			name:   "two slices",
			slices: [][]int{{1, 2, 3}, {10, 20, 30}},
		},
		{
			name:   "empty and non-empty slices",
			slices: [][]int{{}, {1, 2, 3}},
		},
		{
			name:   "all empty slices",
			slices: [][]int{{}, {}, {}},
		},
		{
			name:   "multiple types of slices",
			slices: [][]int{{1, 2, 3, 4, 5}, {}, {10, 20, 30, 40}, {400, 500, 600, 700}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Interleave(rng, tt.slices...)
			totalLength := 0
			for _, s := range tt.slices {
				totalLength += len(s)
			}

			require.Equal(t, totalLength, len(result))

			// check that all sequences are still sorted
			seqs := map[int][]int{}
			for _, x := range result {
				ix := int(math.Log10(float64(x)))
				seqs[ix] = append(seqs[ix], x)
			}

			for _, seq := range seqs {
				assert.True(t, sort.IntsAreSorted(seq))
			}
		})
	}
}

func TestBisect(t *testing.T) {
	tests := []struct {
		name          string
		array         []int
		key           int
		keyExtractor  func(int) (int, error)
		expectedIndex int
		expectedErr   error
	}{
		{
			name:          "empty slice",
			array:         []int{},
			key:           5,
			keyExtractor:  func(x int) (int, error) { return x, nil },
			expectedIndex: 0,
			expectedErr:   nil,
		},
		{
			name:          "single element less than key",
			array:         []int{3},
			key:           5,
			keyExtractor:  func(x int) (int, error) { return x, nil },
			expectedIndex: 1,
			expectedErr:   nil,
		},
		{
			name:          "single element equal to key",
			array:         []int{5},
			key:           5,
			keyExtractor:  func(x int) (int, error) { return x, nil },
			expectedIndex: 1,
			expectedErr:   nil,
		},
		{
			name:          "single element more than key",
			array:         []int{6},
			key:           5,
			keyExtractor:  func(x int) (int, error) { return x, nil },
			expectedIndex: 0,
			expectedErr:   nil,
		},
		{
			name:          "multiple elements without key present",
			array:         []int{1, 3, 4, 6, 7},
			key:           5,
			keyExtractor:  func(x int) (int, error) { return x, nil },
			expectedIndex: 3,
			expectedErr:   nil,
		},
		{
			name:          "multiple elements with key present",
			array:         []int{1, 3, 5, 5, 5, 7, 9},
			key:           5,
			keyExtractor:  func(x int) (int, error) { return x, nil },
			expectedIndex: 5,
			expectedErr:   nil,
		},
		{
			name:  "key extractor returns error",
			array: []int{1, 2, 3},
			key:   2,
			keyExtractor: func(x int) (int, error) {
				if x == 2 {
					return 0, fmt.Errorf("error extracting key")
				}
				return x, nil
			},
			expectedIndex: 0,
			expectedErr:   fmt.Errorf("error extracting key"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index, err := bisect(tt.array, tt.key, tt.keyExtractor)
			if tt.expectedErr != nil {
				require.Error(t, err)
				require.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.expectedIndex, index)
		})
	}
}

func TestUtf8ToUtf16Index(t *testing.T) {
	tests := []struct {
		name      string
		s         string
		utf8Index int
		want      int
	}{
		{"One Unicode", "🌍a", 0, 0},
		{"One Unicode", "🌍a", 1, 0},
		{"One Unicode", "🌍a", 2, 0},
		{"One Unicode", "🌍a", 3, 0},
		{"One Unicode", "🌍a", 4, 2},
		{"Empty string", "", 0, -1},
		{"ASCII only", "Hello", 0, 0},
		{"ASCII only", "Hello", 3, 3},
		{"ASCII only, out of range", "Hello", 10, -1},
		{"With Unicode", "Hello, 世界", 7, 7},
		{"Unicode character", "Hello, 世界", 8, 7},
		{"After Unicode character", "Hello, 世界", 9, 7},
		{"In Emoji", "Hello! 🌍", 7, 7},
		{"In Emoji", "Hello! 🌍", 8, 7},
		{"In Emoji", "Hello! 🌍", 9, 7},
		{"In Emoji", "Hello! 🌍", 10, 7},
		{"Multiple Emojis", "🌍🌎🌏", 0, 0},
		{"Multiple Emojis", "🌍🌎🌏", 4, 2},
		{"Multiple Emojis", "🌍🌎🌏", 8, 4},
		{"Multiple Emojis", "🌍🌎🌏", 10, 4},
		{"Multiple Emojis", "🌍🌎🌏", 11, 4},
		{"Multiple Emojis", "🌍🌎🌏", 12, -1},
		{"Invalid UTF-8", string([]byte{0xff, 0xfe, 0xfd}), 1, -1},
		{"Negative index", "Hello", -1, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Utf8ToUtf16Ix(tt.s, tt.utf8Index)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestSliceMinIx(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected int
	}{
		{
			name:     "empty int slice",
			input:    []int{},
			expected: -1,
		},
		{
			name:     "single int element",
			input:    []int{5},
			expected: 0,
		},
		{
			name:     "multiple ints, min at start",
			input:    []int{1, 2, 3, 4, 5},
			expected: 0,
		},
		{
			name:     "multiple ints, min in middle",
			input:    []int{5, 4, 1, 3, 2},
			expected: 2,
		},
		{
			name:     "multiple ints, min at end",
			input:    []int{5, 4, 3, 2, 1},
			expected: 4,
		},
		{
			name:     "negative numbers",
			input:    []int{-1, -5, -3, -2, -4},
			expected: 1,
		},
		{
			name:     "duplicate min values",
			input:    []int{2, 1, 3, 1, 4},
			expected: 1, // should return first occurrence
		},
		{
			name:     "float64 values",
			input:    []float64{3.14, 2.71, 1.41, 1.73},
			expected: 2,
		},
		{
			name:     "string values",
			input:    []string{"banana", "apple", "cherry", "date"},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch input := tt.input.(type) {
			case []int:
				if got := SliceMinIx(input); got != tt.expected {
					t.Errorf("SliceMinIx() = %v, want %v", got, tt.expected)
				}
			case []float64:
				if got := SliceMinIx(input); got != tt.expected {
					t.Errorf("SliceMinIx() = %v, want %v", got, tt.expected)
				}
			case []string:
				if got := SliceMinIx(input); got != tt.expected {
					t.Errorf("SliceMinIx() = %v, want %v", got, tt.expected)
				}
			}
		})
	}
}

func TestSliceMaxIx(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected int
	}{
		{
			name:     "empty int slice",
			input:    []int{},
			expected: -1,
		},
		{
			name:     "single int element",
			input:    []int{5},
			expected: 0,
		},
		{
			name:     "multiple ints, max at start",
			input:    []int{5, 4, 3, 2, 1},
			expected: 0,
		},
		{
			name:     "multiple ints, max in middle",
			input:    []int{1, 2, 5, 3, 4},
			expected: 2,
		},
		{
			name:     "multiple ints, max at end",
			input:    []int{1, 2, 3, 4, 5},
			expected: 4,
		},
		{
			name:     "negative numbers",
			input:    []int{-1, -5, -3, -2, -4},
			expected: 0,
		},
		{
			name:     "duplicate max values",
			input:    []int{2, 4, 3, 4, 1},
			expected: 1, // should return first occurrence
		},
		{
			name:     "float64 values",
			input:    []float64{3.14, 2.71, 1.41, 1.73},
			expected: 0,
		},
		{
			name:     "string values",
			input:    []string{"banana", "apple", "cherry", "date"},
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch input := tt.input.(type) {
			case []int:
				if got := SliceMaxIx(input); got != tt.expected {
					t.Errorf("SliceMaxIx() = %v, want %v", got, tt.expected)
				}
			case []float64:
				if got := SliceMaxIx(input); got != tt.expected {
					t.Errorf("SliceMaxIx() = %v, want %v", got, tt.expected)
				}
			case []string:
				if got := SliceMaxIx(input); got != tt.expected {
					t.Errorf("SliceMaxIx() = %v, want %v", got, tt.expected)
				}
			}
		})
	}
}
