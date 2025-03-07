package v3_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	v3 "github.com/fivetentaylor/pointy/rogue/v3"
)

func TestGetCurSpanFormat(t *testing.T) {
	type cond struct {
		startIx int
		endIx   int
		format  v3.FormatV3Span
	}

	type testCase struct {
		name      string
		setupFunc func() *v3.Rogue
		conds     []cond
	}

	tests := []testCase{
		{
			name: "basic format test",
			setupFunc: func() *v3.Rogue {
				r := v3.NewRogueForQuill("1")
				_, err := r.Insert(0, "Hello World!")
				require.NoError(t, err)
				_, err = r.Format(0, 5, v3.FormatV3Span{"b": "true"})
				require.NoError(t, err)
				return r
			},
			conds: []cond{
				{
					startIx: 0,
					endIx:   4,
					format:  v3.FormatV3Span{"b": "true", "e": "true", "en": "true"},
				},
				{
					startIx: 0,
					endIx:   3,
					format:  v3.FormatV3Span{"b": "true", "e": "true", "en": "true"},
				},
				{
					startIx: 2,
					endIx:   4,
					format:  v3.FormatV3Span{"b": "true", "e": "true", "en": "true"},
				},
				{
					startIx: 4,
					endIx:   4,
					format:  v3.FormatV3Span{"b": "true", "e": "true", "en": "true"},
				},
				{
					startIx: 4,
					endIx:   6,
					format:  v3.FormatV3Span{"e": "true", "en": "true"},
				},
				{
					startIx: 6,
					endIx:   9,
					format:  v3.FormatV3Span{"e": "true", "en": "true"},
				},
			},
		},
		{
			name: "multi format test",
			setupFunc: func() *v3.Rogue {
				r := v3.NewRogueForQuill("1")

				_, err := r.Insert(0, "Hello World!")
				require.NoError(t, err)

				_, err = r.Format(0, 5, v3.FormatV3Span{"b": "true"})
				require.NoError(t, err)

				_, err = r.Format(3, 8, v3.FormatV3Span{"i": "true"})
				require.NoError(t, err)

				return r
			},
			conds: []cond{
				{
					startIx: 3,
					endIx:   4,
					format:  v3.FormatV3Span{"b": "true", "i": "true", "e": "true", "en": "true"},
				},
				{
					startIx: 3,
					endIx:   6,
					format:  v3.FormatV3Span{"i": "true", "e": "true", "en": "true"},
				},
				{
					startIx: 0,
					endIx:   4,
					format:  v3.FormatV3Span{"b": "true", "e": "true", "en": "true"},
				},
				{
					startIx: 0,
					endIx:   3,
					format:  v3.FormatV3Span{"b": "true", "e": "true", "en": "true"},
				},
				{
					startIx: 2,
					endIx:   4,
					format:  v3.FormatV3Span{"b": "true", "e": "true", "en": "true"},
				},
				{
					startIx: 5,
					endIx:   5,
					format:  v3.FormatV3Span{"i": "true", "b": "true", "e": "true", "en": "true"},
				},
				{
					startIx: 5,
					endIx:   6,
					format:  v3.FormatV3Span{"i": "true", "e": "true", "en": "true"},
				},
				{
					startIx: 6,
					endIx:   9,
					format:  v3.FormatV3Span{"i": "true", "e": "true", "en": "true"},
				},
			},
		},
		{
			name: "sticky and not sticky",
			setupFunc: func() *v3.Rogue {
				r := v3.NewRogueForQuill("1")

				_, err := r.Insert(0, "Hello World!")
				require.NoError(t, err)

				_, err = r.Format(0, 5, v3.FormatV3Span{"b": "true"})
				require.NoError(t, err)

				_, err = r.Format(3, 8, v3.FormatV3Span{"a": "http://revi.so"})
				require.NoError(t, err)

				return r
			},
			conds: []cond{
				{
					startIx: 3,
					endIx:   4,
					format:  v3.FormatV3Span{"b": "true", "a": "http://revi.so", "e": "true", "en": "true"},
				},
				{
					startIx: 3,
					endIx:   6,
					format:  v3.FormatV3Span{"a": "http://revi.so", "e": "true", "en": "true"},
				},
				{
					startIx: 0,
					endIx:   4,
					format:  v3.FormatV3Span{"b": "true", "e": "true", "en": "true"},
				},
				{
					startIx: 0,
					endIx:   3,
					format:  v3.FormatV3Span{"b": "true", "e": "true", "en": "true"},
				},
				{
					startIx: 2,
					endIx:   4,
					format:  v3.FormatV3Span{"b": "true", "e": "true", "en": "true"},
				},
				{
					startIx: 5,
					endIx:   5,
					format:  v3.FormatV3Span{"a": "http://revi.so", "b": "true", "e": "true", "en": "true"},
				},
				{
					startIx: 5,
					endIx:   6,
					format:  v3.FormatV3Span{"a": "http://revi.so", "e": "true", "en": "true"},
				},
				{
					startIx: 6,
					endIx:   9,
					format:  v3.FormatV3Span{"a": "http://revi.so", "e": "true", "en": "true"},
				},
			},
		},
	}

	for _, tc := range tests {
		r := tc.setupFunc()
		for i, c := range tc.conds {
			name := fmt.Sprintf("%s %d", tc.name, i)

			t.Run(name, func(t *testing.T) {
				// Convert positions to IDs
				startID, err := r.Rope.GetVisID(c.startIx)
				require.NoError(t, err)

				endID, err := r.Rope.GetVisID(c.endIx)
				require.NoError(t, err)

				// Get current span format
				format, err := r.GetCurSpanFormat(startID, endID)
				require.NoError(t, err)

				firstID, err := r.Rope.GetTotID(0)
				require.NoError(t, err)

				lastID, err := r.Rope.GetTotID(r.TotSize - 1)
				require.NoError(t, err)

				html, err := r.GetHtml(firstID, lastID, true, false)
				require.NoError(t, err)

				fmt.Printf("%q\n", html)

				// Validate the expected format
				require.Equalf(t, c.format, format, "Test '%s' failed for cond %d with startIx: %d, endIx: %d", tc.name, i, c.startIx, c.endIx)
			})
		}
	}
}

func TestGetCurLineFormat(t *testing.T) {
	type cond struct {
		startIx int
		endIx   int
		format  v3.FormatV3
	}

	type testCase struct {
		name      string
		setupFunc func() *v3.Rogue
		conds     []cond
	}

	tests := []testCase{
		{
			name: "basic format test",
			setupFunc: func() *v3.Rogue {
				r := v3.NewRogue("1")

				_, err := r.Insert(0, "Hello\nWorld!\n")
				require.NoError(t, err)

				_, err = r.Format(0, 1, v3.FormatV3Header(1))
				require.NoError(t, err)

				return r
			},
			conds: []cond{
				{
					startIx: 0,
					endIx:   0,
					format:  v3.FormatV3Header(1),
				},
				{
					startIx: 5,
					endIx:   5,
					format:  v3.FormatV3Header(1),
				},
				{
					startIx: 3,
					endIx:   5,
					format:  v3.FormatV3Header(1),
				},
				{
					startIx: 3,
					endIx:   6,
					format:  v3.FormatV3Line{},
				},
			},
		},
		{
			name: "multi line format test",
			setupFunc: func() *v3.Rogue {
				r := v3.NewRogue("1")

				_, err := r.Insert(0, "Hello\nWorld!\n")
				require.NoError(t, err)

				_, err = r.Format(0, 6, v3.FormatV3Header(1))
				require.NoError(t, err)

				return r
			},
			conds: []cond{
				{
					startIx: 0,
					endIx:   0,
					format:  v3.FormatV3Header(1),
				},
				{
					startIx: 3,
					endIx:   5,
					format:  v3.FormatV3Header(1),
				},
				{
					startIx: 3,
					endIx:   6,
					format:  v3.FormatV3Line{},
				},
			},
		},
	}

	for _, tc := range tests {
		r := tc.setupFunc()
		for i, c := range tc.conds {
			name := fmt.Sprintf("%s %d", tc.name, i)

			t.Run(name, func(t *testing.T) {
				// Convert positions to IDs
				startID, err := r.Rope.GetVisID(c.startIx)
				require.NoError(t, err)

				endID, err := r.Rope.GetVisID(c.endIx)
				require.NoError(t, err)

				// Get current span format
				format, err := r.GetCurLineFormat(startID, endID)
				require.NoError(t, err)

				// Validate the expected format
				require.Equalf(t, c.format, format, "Test '%s' failed for cond %d with startIx: %d, endIx: %d", tc.name, i, c.startIx, c.endIx)
			})
		}
	}
}

func TestEnclosingSpanID(t *testing.T) {
	type cond struct {
		ix             int
		expectedOffset int
		expectedID     v3.ID
	}

	type testCase struct {
		name      string
		setupFunc func() *v3.Rogue
		conds     []cond
	}

	tests := []testCase{
		{
			name: "only newlines",
			setupFunc: func() *v3.Rogue {
				r := v3.NewRogueForQuill("1")

				_, err := r.Insert(0, "\n")
				require.NoError(t, err)

				return r
			},
			conds: []cond{
				{
					ix:             0,
					expectedOffset: 0,
					expectedID:     v3.ID{"1", 3},
				},
				{
					ix:             1,
					expectedOffset: 0,
					expectedID:     v3.ID{"q", 1},
				},
			},
		},
		{
			name: "consecutive newlines",
			setupFunc: func() *v3.Rogue {
				r := v3.NewRogueForQuill("1")

				_, err := r.Insert(0, "hello\n\nworld")
				require.NoError(t, err)

				return r
			},
			conds: []cond{
				{
					ix:             6,
					expectedOffset: 0,
					expectedID:     v3.ID{"1", 9},
				},
			},
		},
		{
			name: "more basic",
			setupFunc: func() *v3.Rogue {
				r := v3.NewRogueForQuill("1")

				_, err := r.Insert(0, "Hello World!\n")
				require.NoError(t, err)

				_, err = r.Format(2, 3, v3.FormatV3Span{"b": "true"})
				require.NoError(t, err)

				return r
			},
			conds: []cond{
				{
					ix:             0,
					expectedOffset: 0,
					expectedID:     v3.ID{"1", 3},
				},
				{
					ix:             2,
					expectedOffset: 0,
					expectedID:     v3.ID{"1", 5},
				},
				{
					ix:             4,
					expectedOffset: 2,
					expectedID:     v3.ID{"1", 5},
				},
				{
					ix:             5,
					expectedOffset: 3,
					expectedID:     v3.ID{"1", 5},
				},
				{
					ix:             6,
					expectedOffset: 1,
					expectedID:     v3.ID{"1", 8},
				},
			},
		},
		{
			name: "basic",
			setupFunc: func() *v3.Rogue {
				r := v3.NewRogueForQuill("1")

				_, err := r.Insert(0, "Hello World!\n")
				require.NoError(t, err)

				_, err = r.Format(0, 5, v3.FormatV3Span{"b": "true"})
				require.NoError(t, err)

				_, err = r.Format(3, 8, v3.FormatV3Span{"i": "true"})
				require.NoError(t, err)

				return r
			},
			conds: []cond{
				{
					ix:             0,
					expectedOffset: 0,
					expectedID:     v3.ID{"1", 3},
				},
				{
					ix:             3,
					expectedOffset: 0,
					expectedID:     v3.ID{"1", 6},
				},
				{
					ix:             4,
					expectedOffset: 1,
					expectedID:     v3.ID{"1", 6},
				},
				{
					ix:             5,
					expectedOffset: 0,
					expectedID:     v3.ID{"1", 8},
				},
				{
					ix:             9,
					expectedOffset: 4,
					expectedID:     v3.ID{"1", 8},
				},
				{
					ix:             11,
					expectedOffset: 6,
					expectedID:     v3.ID{"1", 8},
				},
				{
					ix:             12,
					expectedOffset: 1,
					expectedID:     v3.ID{"1", 14},
				},
			},
		},
		{
			name: "multiline",
			setupFunc: func() *v3.Rogue {
				r := v3.NewRogueForQuill("1")

				_, err := r.Insert(0, "Hello World!\nHow are you?\nI am fine!")
				require.NoError(t, err)

				_, err = r.Format(0, 5, v3.FormatV3Header(1))
				require.NoError(t, err)

				_, err = r.Format(0, 20, v3.FormatV3Span{"b": "true"})
				require.NoError(t, err)

				_, err = r.Format(3, 8, v3.FormatV3Span{"i": "true"})
				require.NoError(t, err)

				return r
			},
			conds: []cond{
				{
					ix:             0,
					expectedOffset: 0,
					expectedID:     v3.ID{"1", 3},
				},
				{
					ix:             3,
					expectedOffset: 0,
					expectedID:     v3.ID{"1", 6},
				},
				{
					ix:             10,
					expectedOffset: 7,
					expectedID:     v3.ID{"1", 6},
				},
				{
					ix:             11,
					expectedOffset: 0,
					expectedID:     v3.ID{"1", 14},
				},
				{
					ix:             13,
					expectedOffset: 0,
					expectedID:     v3.ID{"1", 16},
				},
				{
					ix:             19,
					expectedOffset: 6,
					expectedID:     v3.ID{"1", 16},
				},
				{
					ix:             20,
					expectedOffset: 7,
					expectedID:     v3.ID{"1", 16},
				},
				{
					ix:             24,
					expectedOffset: 4,
					expectedID:     v3.ID{"1", 23},
				},
				{
					ix:             26,
					expectedOffset: 0,
					expectedID:     v3.ID{"1", 29},
				},
				{
					ix:             35,
					expectedOffset: 9,
					expectedID:     v3.ID{"1", 29},
				},
			},
		},
		{
			name: "code block",
			setupFunc: func() *v3.Rogue {
				r := v3.NewRogueForQuill("1")

				_, err := r.Insert(0, "for i in range(10):\n    print(i)")
				require.NoError(t, err)

				_, err = r.Format(18, 5, v3.FormatV3CodeBlock("python"))
				require.NoError(t, err)

				return r
			},
			conds: []cond{
				{
					ix:             0,
					expectedOffset: 0,
					expectedID:     v3.ID{"1", 3},
				},
				{
					ix:             20,
					expectedOffset: 20,
					expectedID:     v3.ID{"1", 3},
				},
			},
		},
		{
			name: "mixed code block",
			setupFunc: func() *v3.Rogue {
				r := v3.NewRogueForQuill("1")

				_, err := r.Insert(0, "Hello World!\nfor i in range(10):\n    print(i)\ncool")
				require.NoError(t, err)

				_, err = r.Format(27, 8, v3.FormatV3CodeBlock("python"))
				require.NoError(t, err)

				return r
			},
			conds: []cond{
				{
					ix:             0,
					expectedOffset: 0,
					expectedID:     v3.ID{"1", 3},
				},
				{
					ix:             13,
					expectedOffset: 0,
					expectedID:     v3.ID{"1", 16},
				},
				{
					ix:             31,
					expectedOffset: 18,
					expectedID:     v3.ID{"1", 16},
				},
				{
					ix:             32,
					expectedOffset: 19,
					expectedID:     v3.ID{"1", 16},
				},
				{
					ix:             33,
					expectedOffset: 20,
					expectedID:     v3.ID{"1", 16},
				},
				{
					ix:             45,
					expectedOffset: 32,
					expectedID:     v3.ID{"1", 16},
				},
				{
					ix:             46,
					expectedOffset: 0,
					expectedID:     v3.ID{"1", 49},
				},
			},
		},
		{
			name: "with quotes",
			setupFunc: func() *v3.Rogue {
				r := v3.NewRogueForQuill("1")

				_, err := r.Insert(0, "\n\"hello\"")
				require.NoError(t, err)

				return r
			},
			conds: []cond{
				{
					ix:             0,
					expectedOffset: 0,
					expectedID:     v3.ID{"1", 3},
				},
				{
					ix:             1,
					expectedOffset: 0,
					expectedID:     v3.ID{"1", 4},
				},
				{
					ix:             2,
					expectedOffset: 1,
					expectedID:     v3.ID{"1", 4},
				},
				{
					ix:             3,
					expectedOffset: 1,
					expectedID:     v3.ID{"1", 5},
				},
			},
		},
	}

	for _, tc := range tests {
		r := tc.setupFunc()
		for i, c := range tc.conds {
			name := fmt.Sprintf("%s %d", tc.name, i)

			t.Run(name, func(t *testing.T) {
				firstID, err := r.Rope.GetTotID(0)
				require.NoError(t, err)

				lastID, err := r.Rope.GetTotID(r.TotSize - 1)
				require.NoError(t, err)

				// print html with IDs for debugging
				html, err := r.GetHtml(firstID, lastID, true, true)
				require.NoError(t, err)

				fmt.Println(html)

				targetID, err := r.Rope.GetVisID(c.ix)
				require.NoError(t, err)

				id, offset, err := r.EnclosingSpanID(targetID, nil, true)
				require.NoError(t, err)

				require.Equal(t, c.expectedID, id)
				require.Equal(t, c.expectedOffset, offset)
			})
		}
	}
}

func TestGetContainingLine(t *testing.T) {
	type cond struct {
		ix   int
		text string
	}

	type testCase struct {
		name      string
		setupFunc func() *v3.Rogue
		conds     []cond
	}

	tests := []testCase{
		{
			name: "basic",
			setupFunc: func() *v3.Rogue {
				r := v3.NewRogue("1")

				_, err := r.Insert(0, "Hello World!\nHow are you?\nI am fine!\n")
				require.NoError(t, err)

				return r
			},
			conds: []cond{
				{
					ix:   0,
					text: "Hello World!\n",
				},
				{
					ix:   12,
					text: "Hello World!\n",
				},
				{
					ix:   13,
					text: "How are you?\n",
				},
				{
					ix:   25,
					text: "How are you?\n",
				},
				{
					ix:   26,
					text: "I am fine!\n",
				},
				{
					ix:   30,
					text: "I am fine!\n",
				},
				{
					ix:   36,
					text: "I am fine!\n",
				},
			},
		},
	}

	for _, tc := range tests {
		r := tc.setupFunc()
		for i, c := range tc.conds {
			name := fmt.Sprintf("%s %d", tc.name, i)

			t.Run(name, func(t *testing.T) {
				id, err := r.Rope.GetVisID(c.ix)
				require.NoError(t, err)

				startID, endID, _, err := r.GetLineAt(id, nil)
				require.NoError(t, err)

				vis, err := r.Rope.GetBetween(startID, endID)
				require.NoError(t, err)

				require.Equal(t, c.text, v3.Uint16ToStr(vis.Text))
			})
		}
	}
}
