package v3_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	rogue "github.com/teamreviso/code/rogue/v3"
)

func TestFormatEquals(t *testing.T) {
	type testCase struct {
		name     string
		format1  rogue.FormatV3
		format2  rogue.FormatV3
		expected bool
	}

	tests := []testCase{
		{
			name:     "Test equal span",
			format1:  rogue.FormatV3Span{"b": "true"},
			format2:  rogue.FormatV3Span{"b": "true"},
			expected: true,
		},
		{
			name:     "Test not equal span",
			format1:  rogue.FormatV3Span{"b": "true"},
			format2:  rogue.FormatV3Span{"b": "false"},
			expected: false,
		},
		{
			name:     "Test equal span multiple keys",
			format1:  rogue.FormatV3Span{"b": "true", "i": "true"},
			format2:  rogue.FormatV3Span{"b": "true", "i": "true"},
			expected: true,
		},
		{
			name:     "Test not equal span multiple keys",
			format1:  rogue.FormatV3Span{"b": "true", "i": "true"},
			format2:  rogue.FormatV3Span{"b": "true", "i": "", "u": "true"},
			expected: false,
		},
		{
			name:     "Test equal bullet list",
			format1:  rogue.FormatV3BulletList(0),
			format2:  rogue.FormatV3BulletList(0),
			expected: true,
		},
		{
			name:     "Test not equal bullet list",
			format1:  rogue.FormatV3BulletList(0),
			format2:  rogue.FormatV3BulletList(1),
			expected: false,
		},
		{
			name:     "Test equal line",
			format1:  rogue.FormatV3Line{},
			format2:  rogue.FormatV3Line{},
			expected: true,
		},
		{
			name:     "Test not equal line",
			format1:  rogue.FormatV3Line{},
			format2:  rogue.FormatV3Span{"b": "true"},
			expected: false,
		},
		{
			name:     "Test not equal line",
			format1:  rogue.FormatV3Line{},
			format2:  rogue.FormatV3Header(1),
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, tc.format1.Equals(tc.format2))
		})
	}
}

func TestSpanDifference(t *testing.T) {
	tests := []struct {
		name     string
		s        rogue.Span
		b        rogue.Span
		expected []rogue.Span
	}{
		{
			name:     "Case 1: b is completely before s",
			s:        rogue.Span{StartIx: 10, EndIx: 20},
			b:        rogue.Span{StartIx: 5, EndIx: 8},
			expected: []rogue.Span{{StartIx: 10, EndIx: 20}},
		},
		{
			name:     "Case 1: b is completely after s",
			s:        rogue.Span{StartIx: 10, EndIx: 20},
			b:        rogue.Span{StartIx: 25, EndIx: 30},
			expected: []rogue.Span{{StartIx: 10, EndIx: 20}},
		},
		{
			name:     "Case 2: b completely covers s",
			s:        rogue.Span{StartIx: 10, EndIx: 20},
			b:        rogue.Span{StartIx: 5, EndIx: 25},
			expected: nil,
		},
		{
			name:     "Case 3: b overlaps with the start of s",
			s:        rogue.Span{StartIx: 10, EndIx: 20},
			b:        rogue.Span{StartIx: 5, EndIx: 15},
			expected: []rogue.Span{{StartIx: 16, EndIx: 20}},
		},
		{
			name:     "Case 4: b is completely inside s",
			s:        rogue.Span{StartIx: 10, EndIx: 20},
			b:        rogue.Span{StartIx: 12, EndIx: 18},
			expected: []rogue.Span{{StartIx: 10, EndIx: 11}, {StartIx: 19, EndIx: 20}},
		},
		{
			name:     "Case 4: b is single value completely inside s",
			s:        rogue.Span{StartIx: 10, EndIx: 20},
			b:        rogue.Span{StartIx: 18, EndIx: 18},
			expected: []rogue.Span{{StartIx: 10, EndIx: 17}, {StartIx: 19, EndIx: 20}},
		},
		{
			name:     "Case 5: b overlaps with the end of s",
			s:        rogue.Span{StartIx: 10, EndIx: 20},
			b:        rogue.Span{StartIx: 15, EndIx: 25},
			expected: []rogue.Span{{StartIx: 10, EndIx: 14}},
		},
		{
			name:     "Edge case: b starts at s.StartIx",
			s:        rogue.Span{StartIx: 10, EndIx: 20},
			b:        rogue.Span{StartIx: 10, EndIx: 15},
			expected: []rogue.Span{{StartIx: 16, EndIx: 20}},
		},
		{
			name:     "Edge case: b ends at s.EndIx",
			s:        rogue.Span{StartIx: 10, EndIx: 20},
			b:        rogue.Span{StartIx: 15, EndIx: 20},
			expected: []rogue.Span{{StartIx: 10, EndIx: 14}},
		},
		{
			name:     "Edge case: b is single value at end of s",
			s:        rogue.Span{StartIx: 7, EndIx: 12},
			b:        rogue.Span{StartIx: 12, EndIx: 12},
			expected: []rogue.Span{{7, 11}},
		},
		{
			name:     "Edge case: b is single value at start of s",
			s:        rogue.Span{StartIx: 7, EndIx: 12},
			b:        rogue.Span{StartIx: 7, EndIx: 7},
			expected: []rogue.Span{{8, 12}},
		},
		{
			name:     "Edge case: s is single value inside of b",
			s:        rogue.Span{StartIx: 9, EndIx: 9},
			b:        rogue.Span{StartIx: 7, EndIx: 12},
			expected: nil,
		},
		{
			name:     "Edge case: s is single value at end of b",
			s:        rogue.Span{StartIx: 12, EndIx: 12},
			b:        rogue.Span{StartIx: 7, EndIx: 12},
			expected: nil,
		},
		{
			name:     "Edge case: s is single value at start of b",
			s:        rogue.Span{StartIx: 7, EndIx: 7},
			b:        rogue.Span{StartIx: 7, EndIx: 12},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.s.Difference(tt.b)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestSpanIntersection(t *testing.T) {
	tests := []struct {
		name     string
		s        rogue.Span
		b        rogue.Span
		expected *rogue.Span
	}{
		{
			name:     "Case 1: b is completely before s",
			s:        rogue.Span{StartIx: 10, EndIx: 20},
			b:        rogue.Span{StartIx: 5, EndIx: 8},
			expected: nil,
		},
		{
			name:     "Case 1: b is completely after s",
			s:        rogue.Span{StartIx: 10, EndIx: 20},
			b:        rogue.Span{StartIx: 25, EndIx: 30},
			expected: nil,
		},
		{
			name:     "Case 2: b completely covers s",
			s:        rogue.Span{StartIx: 10, EndIx: 20},
			b:        rogue.Span{StartIx: 5, EndIx: 25},
			expected: &rogue.Span{StartIx: 10, EndIx: 20},
		},
		{
			name:     "Case 3: b overlaps with the start of s",
			s:        rogue.Span{StartIx: 10, EndIx: 20},
			b:        rogue.Span{StartIx: 5, EndIx: 15},
			expected: &rogue.Span{StartIx: 10, EndIx: 15},
		},
		{
			name:     "Case 4: b is completely inside s",
			s:        rogue.Span{StartIx: 10, EndIx: 20},
			b:        rogue.Span{StartIx: 12, EndIx: 18},
			expected: &rogue.Span{StartIx: 12, EndIx: 18},
		},
		{
			name:     "Case 4: b is single value completely inside s",
			s:        rogue.Span{StartIx: 10, EndIx: 20},
			b:        rogue.Span{StartIx: 18, EndIx: 18},
			expected: &rogue.Span{StartIx: 18, EndIx: 18},
		},
		{
			name:     "Case 5: b overlaps with the end of s",
			s:        rogue.Span{StartIx: 10, EndIx: 20},
			b:        rogue.Span{StartIx: 15, EndIx: 25},
			expected: &rogue.Span{StartIx: 15, EndIx: 20},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.s.Intersection(tt.b)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestSearchOverlapping(t *testing.T) {
	type testCase struct {
		name           string
		setup          func(doc *rogue.Rogue)
		startID        rogue.ID
		endID          rogue.ID
		expectedLength int
	}

	tests := []testCase{
		{
			name: "Test with valid range",
			setup: func(doc *rogue.Rogue) {
				doc.Insert(0, "Hello World!")
				doc.Format(0, 5, rogue.FormatV3Span{"b": "true"})
			},
			startID:        rogue.ID{"auth0", 5},
			endID:          rogue.ID{"auth0", 6},
			expectedLength: 1,
		},
		{
			name: "Test with no overlap",
			setup: func(doc *rogue.Rogue) {
				doc.Insert(0, "Hello World!")
				doc.Format(0, 5, rogue.FormatV3Span{"b": "true"})
				doc.Format(12, 1, rogue.FormatV3BulletList(0))
			},
			startID:        rogue.ID{"auth0", 9},
			endID:          rogue.ID{"auth0", 11},
			expectedLength: 0,
		},
		{
			name: "Test with entire range",
			setup: func(doc *rogue.Rogue) {
				doc.Insert(0, "Hello World!")
				doc.Format(0, 5, rogue.FormatV3Span{"b": "true"})
				doc.Format(12, 1, rogue.FormatV3BulletList(0))
			},
			startID:        rogue.ID{"root", 0},
			endID:          rogue.ID{"q", 1},
			expectedLength: 2,
		},
		{
			name: "Test with single id",
			setup: func(doc *rogue.Rogue) {
				doc.Insert(0, "Hello World!")
				doc.Format(0, 5, rogue.FormatV3Span{"b": "true"})
				doc.Format(3, 8, rogue.FormatV3Span{"u": "true"})
				doc.Format(12, 1, rogue.FormatV3BulletList(0))
			},
			startID:        rogue.ID{"auth0", 6},
			endID:          rogue.ID{"auth0", 6},
			expectedLength: 2,
		},
		{
			name: "Test with 3 overlapping",
			setup: func(doc *rogue.Rogue) {
				doc.Insert(0, "Hello World!")
				doc.Format(0, 5, rogue.FormatV3Span{"b": "true"})
				doc.Format(3, 8, rogue.FormatV3Span{"u": "true"})
				doc.Format(12, 1, rogue.FormatV3BulletList(0))
			},
			startID:        rogue.ID{"auth0", 6},
			endID:          rogue.ID{"q", 1},
			expectedLength: 4,
		},
		{
			name: "Test with end of underline",
			setup: func(doc *rogue.Rogue) {
				doc.Insert(0, "Hello World!")
				doc.Format(0, 5, rogue.FormatV3Span{"b": "true"})
				doc.Format(3, 8, rogue.FormatV3Span{"u": "true"})
				doc.Format(12, 1, rogue.FormatV3BulletList(0))
			},
			startID:        rogue.ID{"auth0", 9},
			endID:          rogue.ID{"auth0", 10},
			expectedLength: 1,
		},
		{
			name: "Test lots of formats",
			setup: func(doc *rogue.Rogue) {
				doc.Insert(0, "Hello World!")
				doc.Format(0, 5, rogue.FormatV3Span{"b": "true"})
				doc.Format(3, 8, rogue.FormatV3Span{"u": "true"})
				doc.Format(12, 1, rogue.FormatV3BulletList(0))
				doc.Format(2, 8, rogue.FormatV3Span{"b": "true"})
				doc.Format(3, 3, rogue.FormatV3Span{"b": "null"})
				doc.Format(2, 8, rogue.FormatV3Span{"u": "null"})
				doc.Format(5, 5, rogue.FormatV3Span{"b": "true"})
				doc.Format(0, 10, rogue.FormatV3Span{"s": "true"})
				doc.Format(0, 12, rogue.FormatV3Span{"b": "null"})
			},
			startID:        rogue.ID{"root", 0},
			endID:          rogue.ID{"q", 1},
			expectedLength: 17,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			doc := rogue.NewRogueForQuill("auth0")
			tc.setup(doc)

			formatOps, err := doc.Formats.SearchOverlapping(tc.startID, tc.endID)
			require.NoError(t, err)
			require.Len(t, formatOps, tc.expectedLength)

			isBalanced, height := doc.Formats.Sticky.Root.IsBalanced()
			fmt.Printf("sticky isBalanced: %v, height: %d\n", isBalanced, height)
			if !isBalanced {
				doc.Formats.Sticky.Print()
			}
			require.True(t, isBalanced, "Sticky tree is not balanced")

			isBalanced, height = doc.Formats.NoSticky.Root.IsBalanced()
			fmt.Printf("no sticky isBalanced: %v, height: %d\n", isBalanced, height)
			if !isBalanced {
				doc.Formats.NoSticky.Print()
			}
			require.True(t, isBalanced, "No sticky tree is not balanced")

			fmt.Printf("formatOps: %+v\n", formatOps)
		})
	}
}
