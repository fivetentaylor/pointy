package v3

import (
	"fmt"
	"testing"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Assuming your FormatV3 and mergeFormats definitions are here

// TestMergeFormats tests the mergeFormats function
func TestNosMergeFormats(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		formatA  FormatV3
		formatB  FormatV3
		expected FormatV3
	}{
		{
			name:     "Merging distinct sets",
			formatA:  FormatV3Span{"color": "red", "size": "10"},
			formatB:  FormatV3Span{"font": "Arial", "strike": "true"},
			expected: FormatV3Span{"color": "red", "size": "10", "font": "Arial", "strike": "true"},
		},
		{
			name:     "Merging with overlapping keys",
			formatA:  FormatV3Span{"color": "red", "size": "10"},
			formatB:  FormatV3Span{"color": "blue", "strike": "true"},
			expected: FormatV3Span{"color": "blue", "size": "10", "strike": "true"},
		},
		{
			name:     "exclusive merge",
			formatA:  FormatV3Span{"color": "red", "size": "10"},
			formatB:  FormatV3Span{"e": "true"},
			expected: FormatV3Span{"e": "true"},
		},
		{
			name:     "exclusive merge with other stuff",
			formatA:  FormatV3Span{"color": "red", "size": "10"},
			formatB:  FormatV3Span{"e": "true", "b": "true"},
			expected: FormatV3Span{"e": "true", "b": "true"},
		},
		{
			name:     "sticky and no sticky merge",
			formatA:  FormatV3Span{"b": "true", "u": "true"},
			formatB:  FormatV3Span{"a": "google.com"},
			expected: FormatV3Span{"b": "true", "u": "true", "a": "google.com"},
		},
		{
			name:     "mixed with sticky exclusive",
			formatA:  FormatV3Span{"e": "true"},
			formatB:  FormatV3Span{"a": "google.com"},
			expected: FormatV3Span{"e": "true", "a": "google.com"},
		},
		{
			name:     "mixed with no sticky exclusive",
			formatA:  FormatV3Span{"a": "google.com", "b": "true"},
			formatB:  FormatV3Span{"en": "true"},
			expected: FormatV3Span{"b": "true", "en": "true"},
		},
		{
			name:     "mixed",
			formatA:  FormatV3Span{"a": "google.com", "b": "true"},
			formatB:  FormatV3Span{"a": "revi.so", "u": "true"},
			expected: FormatV3Span{"a": "revi.so", "b": "true", "u": "true"},
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := mergeFormats(tc.formatA, tc.formatB)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestNosDropNull(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		input    FormatV3
		expected FormatV3
	}{
		{
			name:     "Remove null values",
			input:    FormatV3Span{"color": "red", "size": "null", "width": "100"},
			expected: FormatV3Span{"color": "red", "width": "100"},
		},
		{
			name:     "No null values",
			input:    FormatV3Span{"font": "Arial", "height": "200"},
			expected: FormatV3Span{"font": "Arial", "height": "200"},
		},
		{
			name:     "All null values",
			input:    FormatV3Span{"color": "null", "size": "null"},
			expected: FormatV3Span{},
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.input.DropNull()
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestNosMergeNodes(t *testing.T) {
	// Define test cases
	tests := []struct {
		name      string
		cur       NOSNode
		new       NOSNode
		toInsert  []NOSNode
		toDelete  []NOSNode
		remaining *NOSNode
	}{
		{
			name:      "No overlap new after cur",
			cur:       NOSNode{StartIx: 1, EndIx: 3, Format: FormatV3Span{"color": "red"}},
			new:       NOSNode{StartIx: 5, EndIx: 7, Format: FormatV3Span{"color": "blue"}},
			toInsert:  nil,
			toDelete:  nil,
			remaining: &NOSNode{StartIx: 5, EndIx: 7, Format: FormatV3Span{"color": "blue"}},
		},
		{
			name:      "No overlap cur after new",
			cur:       NOSNode{StartIx: 5, EndIx: 7, Format: FormatV3Span{"color": "blue"}},
			new:       NOSNode{StartIx: 1, EndIx: 3, Format: FormatV3Span{"color": "red"}},
			toInsert:  []NOSNode{{StartIx: 1, EndIx: 3, Format: FormatV3Span{"color": "red"}}},
			toDelete:  nil,
			remaining: nil,
		},
		{
			name: "Partial overlap",
			cur:  NOSNode{StartIx: 1, EndIx: 5, Format: FormatV3Span{"color": "red"}},
			new:  NOSNode{StartIx: 4, EndIx: 6, Format: FormatV3Span{"size": "large"}},
			toInsert: []NOSNode{
				{StartIx: 1, EndIx: 3, Format: FormatV3Span{"color": "red"}},
				{StartIx: 4, EndIx: 5, Format: FormatV3Span{"color": "red", "size": "large"}},
			},
			toDelete:  []NOSNode{{StartIx: 1, EndIx: 5, Format: FormatV3Span{"color": "red"}}},
			remaining: &NOSNode{StartIx: 6, EndIx: 6, Format: FormatV3Span{"size": "large"}},
		},
		{
			name: "Complete overlap by new node",
			cur:  NOSNode{StartIx: 1, EndIx: 5, Format: FormatV3Span{"color": "red"}},
			new:  NOSNode{StartIx: 0, EndIx: 7, Format: FormatV3Span{"size": "large"}},
			toInsert: []NOSNode{
				{StartIx: 0, EndIx: 0, Format: FormatV3Span{"size": "large"}},
				{StartIx: 1, EndIx: 5, Format: FormatV3Span{"color": "red", "size": "large"}},
			},
			toDelete:  []NOSNode{{StartIx: 1, EndIx: 5, Format: FormatV3Span{"color": "red"}}},
			remaining: &NOSNode{StartIx: 6, EndIx: 7, Format: FormatV3Span{"size": "large"}},
		},
		{
			name: "Complete overlap by cur node",
			cur:  NOSNode{StartIx: 1, EndIx: 5, Format: FormatV3Span{"bold": "true"}},
			new:  NOSNode{StartIx: 2, EndIx: 4, Format: FormatV3Span{"strike": "true"}},
			toInsert: []NOSNode{
				{StartIx: 1, EndIx: 1, Format: FormatV3Span{"bold": "true"}},
				{StartIx: 2, EndIx: 4, Format: FormatV3Span{"bold": "true", "strike": "true"}},
				{StartIx: 5, EndIx: 5, Format: FormatV3Span{"bold": "true"}},
			},
			toDelete:  []NOSNode{{StartIx: 1, EndIx: 5, Format: FormatV3Span{"bold": "true"}}},
			remaining: nil,
		},
		{
			name: "Exact overlap",
			cur:  NOSNode{StartIx: 1, EndIx: 5, Format: FormatV3Span{"color": "red"}},
			new:  NOSNode{StartIx: 1, EndIx: 5, Format: FormatV3Span{"size": "large"}},
			toInsert: []NOSNode{
				{StartIx: 1, EndIx: 5, Format: FormatV3Span{"color": "red", "size": "large"}},
			},
			toDelete:  []NOSNode{{StartIx: 1, EndIx: 5, Format: FormatV3Span{"color": "red"}}},
			remaining: nil,
		},
		{
			name: "Left Exact overlap",
			cur:  NOSNode{StartIx: 1, EndIx: 5, Format: FormatV3Span{"color": "red"}},
			new:  NOSNode{StartIx: 1, EndIx: 8, Format: FormatV3Span{"size": "large"}},
			toInsert: []NOSNode{
				{StartIx: 1, EndIx: 5, Format: FormatV3Span{"color": "red", "size": "large"}},
			},
			toDelete:  []NOSNode{{StartIx: 1, EndIx: 5, Format: FormatV3Span{"color": "red"}}},
			remaining: &NOSNode{StartIx: 6, EndIx: 8, Format: FormatV3Span{"size": "large"}},
		},
		{
			name: "Right Exact overlap",
			cur:  NOSNode{StartIx: 3, EndIx: 8, Format: FormatV3Span{"color": "red"}},
			new:  NOSNode{StartIx: 1, EndIx: 8, Format: FormatV3Span{"size": "large"}},
			toInsert: []NOSNode{
				{StartIx: 1, EndIx: 2, Format: FormatV3Span{"size": "large"}},
				{StartIx: 3, EndIx: 8, Format: FormatV3Span{"size": "large", "color": "red"}},
			},
			toDelete:  []NOSNode{{StartIx: 3, EndIx: 8, Format: FormatV3Span{"color": "red"}}},
			remaining: nil,
		},
		/*{
			name: "Matching formats merge",
			cur:  NOSNode{StartIx: 3, EndIx: 8, Format: FormatV3Span{"color": "red"}},
			new:  NOSNode{StartIx: 7, EndIx: 12, Format: FormatV3Span{"color": "red"}},
			toInsert: []NOSNode{
				{StartIx: 3, EndIx: 12, Format: FormatV3Span{"color": "red"}},
			},
			toDelete:  []NOSNode{{StartIx: 3, EndIx: 8, Format: FormatV3Span{"color": "red"}}},
			remaining: nil,
		},*/
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			toInsert, toDelete, remaining, err := mergeNodes(tc.cur, tc.new)
			require.NoError(t, err)

			require.Equal(t, tc.toInsert, toInsert)
			require.Equal(t, tc.toDelete, toDelete)
			require.Equal(t, tc.remaining, remaining)
		})
	}
}

func TestNosDiffFormats(t *testing.T) {
	// Define test cases
	tests := []struct {
		name            string
		oldFormat       FormatV3
		newFormat       FormatV3
		expectedAdded   FormatV3
		expectedRemoved FormatV3
	}{
		{
			name:            "Elements added and removed",
			oldFormat:       FormatV3Span{"color": "red", "size": "medium"},
			newFormat:       FormatV3Span{"color": "red", "size": "large", "width": "100"},
			expectedAdded:   FormatV3Span{"size": "large", "width": "100"},
			expectedRemoved: FormatV3Span{"size": "medium"},
		},
		{
			name:            "No changes",
			oldFormat:       FormatV3Span{"color": "red"},
			newFormat:       FormatV3Span{"color": "red"},
			expectedAdded:   FormatV3Span{},
			expectedRemoved: FormatV3Span{},
		},
		{
			name:            "All elements removed",
			oldFormat:       FormatV3Span{"color": "red", "size": "medium"},
			newFormat:       FormatV3Span{},
			expectedAdded:   FormatV3Span{},
			expectedRemoved: FormatV3Span{"color": "red", "size": "medium"},
		},
		{
			name:            "All elements added",
			oldFormat:       FormatV3Span{},
			newFormat:       FormatV3Span{"color": "blue", "width": "200"},
			expectedAdded:   FormatV3Span{"color": "blue", "width": "200"},
			expectedRemoved: FormatV3Span{},
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			added, removed := diffFormats(tc.oldFormat, tc.newFormat)
			require.Equal(t, tc.expectedAdded, added)
			require.Equal(t, tc.expectedRemoved, removed)
		})
	}
}

func TestNosWalkNOS(t *testing.T) {
	nos := NewNOS()

	nos.tree.Put(&NOSNode{StartIx: 1})
	nos.tree.Put(&NOSNode{StartIx: 4})
	nos.tree.Put(&NOSNode{StartIx: 5})
	nos.tree.Put(&NOSNode{StartIx: 7})
	nos.tree.Put(&NOSNode{StartIx: 10})

	left, err := nos.tree.FindLeftSib(6)
	assert.NoError(t, err)
	assert.Equal(t, 5, left.StartIx)

	right, err := nos.tree.FindRightSib(6)
	assert.NoError(t, err)
	assert.Equal(t, 7, right.StartIx)

	start, err := nos.tree.FindLeftSibNode(6)
	require.NoError(t, err)

	traversal := []int{}
	start.WalkRight(func(node *NOSNode) error {
		traversal = append(traversal, node.StartIx)
		return nil
	})

	assert.Equal(t, []int{5, 7, 10}, traversal)
}

func TestNosTreeMap(t *testing.T) {
	tree := treemap.NewWithIntComparator()

	it := tree.Iterator()
	fmt.Printf("it.Next(): %v\n", it.Next())

	tree.Put(1, []string{"one"})
	res, _ := tree.Get(1)
	PrintJson("GET 1", res)
	res, ok := tree.Get(2)
	fmt.Printf("ok: %v\n", ok)
	PrintJson("GET 2", res)

}

func TestNosStringToLineNOS(t *testing.T) {
	nos := StringToLineNOS(StrToUint16("Hello\nnew\nline\nWorld💀😂!\nokay"))
	s := nos.tree.AsSlice()
	PrintJson("NOS", s)
	require.Equal(t, 5, len(s))

	nos = StringToLineNOS(StrToUint16("Hello World!"))
	s = nos.tree.AsSlice()
	PrintJson("NOS", s)
	require.Equal(t, 1, len(s))
}

func TestNosIterPairs(t *testing.T) {
	nos := NewNOS()

	nos.tree.Put(&NOSNode{StartIx: 1, EndIx: 3, Format: FormatV3Span{"color": "red"}})
	nos.tree.Put(&NOSNode{StartIx: 4, EndIx: 6, Format: FormatV3Span{"color": "blue"}})
	nos.tree.Put(&NOSNode{StartIx: 7, EndIx: 9, Format: FormatV3Span{"color": "green"}})
	nos.tree.Put(&NOSNode{StartIx: 10, EndIx: 12, Format: FormatV3Span{"color": "yellow"}})

	var out []interface{}
	err := nos.IterPairs(func(prev, next *NOSNode) error {
		out = append(out, prev)
		fmt.Printf("prev: %v, next: %v\n", prev, next)
		return nil
	})

	require.Equal(t, len(out), 5)
	require.NoError(t, err)
}
