package v3_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	v3 "github.com/teamreviso/code/rogue/v3"
)

func TestRichInsert(t *testing.T) {
	testCases := []struct {
		name             string
		fn               func(*v3.Rogue)
		visIx            int
		selLen           int
		text             string
		orgSpanFormat    v3.FormatV3Span
		expectedErr      error
		expectedMkd      string
		expectedHtml     *string
		expectedCursorID *v3.ID
		includeIDs       bool
	}{
		{
			name: "insert with no formatting",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "abc")
				require.NoError(t, err)
			},
			visIx:         3,
			text:          "def",
			orgSpanFormat: v3.FormatV3Span{},
			expectedErr:   nil,
			expectedMkd:   "abcdef\n\n",
		},
		{
			name: "insert \n in list",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "a\nb\nc\nd")
				require.NoError(t, err)
				_, err = r.Format(0, 8, v3.FormatV3BulletList(0))
				require.NoError(t, err)
			},
			visIx:         3,
			text:          "\n123",
			orgSpanFormat: v3.FormatV3Span{},
			expectedErr:   nil,
			expectedMkd:   "- a\n- b\n- 123\n- c\n- d\n",
		},
		{
			name: "insert two \n in list should break the list",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "a\nb\nc\n")
				require.NoError(t, err)
				_, err = r.Format(0, 5, v3.FormatV3OrderedList(0))
				require.NoError(t, err)
				_, _, err = r.RichInsert(3, 0, v3.FormatV3Span{}, "\n")
				require.NoError(t, err)
			},
			visIx:         4,
			text:          "\n",
			orgSpanFormat: v3.FormatV3Span{},
			expectedErr:   nil,
			expectedMkd:   "1. a\n1. b\n\n\n1. c\n\n\n",
		},
		{
			name: "insert \n after bold",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "abc")
				require.NoError(t, err)
				_, err = r.Format(0, 3, v3.FormatV3Span{"b": "true"})
				require.NoError(t, err)
			},
			visIx:         3,
			text:          "\n123",
			orgSpanFormat: v3.FormatV3Span{"b": "true"},
			expectedErr:   nil,
			expectedMkd:   "**abc**\n\n**123**\n\n",
		},
		{
			name: "insert newline after header, should drop header format",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "header")
				require.NoError(t, err)
				_, err = r.Format(0, 5, v3.FormatV3Header(2))
				require.NoError(t, err)
			},
			visIx:         6,
			text:          "\n123",
			orgSpanFormat: v3.FormatV3Span{},
			expectedErr:   nil,
			expectedMkd:   "## header\n\n123\n\n",
		},
		{
			name: "italic span markdown shortcut",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "*hello")
				require.NoError(t, err)
			},
			visIx:         6,
			text:          "*",
			orgSpanFormat: v3.FormatV3Span{},
			expectedErr:   nil,
			expectedMkd:   "*hello*\n\n",
		},
		{
			name: "bold span markdown shortcut",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "**hello*")
				require.NoError(t, err)
			},
			visIx:         8,
			text:          "*",
			orgSpanFormat: v3.FormatV3Span{},
			expectedErr:   nil,
			expectedMkd:   "**hello**\n\n",
		},
		{
			name: "strike span markdown shortcut",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "~hello")
				require.NoError(t, err)
			},
			visIx:         6,
			text:          "~",
			orgSpanFormat: v3.FormatV3Span{},
			expectedErr:   nil,
			expectedMkd:   "~~hello~~\n\n",
		},
		{
			name: "space after bold",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "hello")
				require.NoError(t, err)

				_, err = r.Format(0, 5, v3.FormatV3Span{"b": "true"})
				require.NoError(t, err)
			},
			visIx:         5,
			text:          " ",
			orgSpanFormat: v3.FormatV3Span{"b": "", "e": "true"},
			expectedErr:   nil,
			expectedMkd:   "**hello** \n\n",
		},
		{
			name: "code block shortcut",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "```python")
				require.NoError(t, err)
			},
			visIx:         9,
			text:          " ",
			orgSpanFormat: v3.FormatV3Span{},
			expectedErr:   nil,
			expectedMkd:   "```python\n\n```\n",
		},
		{
			name: "link shortcut",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "[hello](google.com")
				require.NoError(t, err)
			},
			visIx:         18,
			text:          ")",
			orgSpanFormat: v3.FormatV3Span{},
			expectedErr:   nil,
			expectedMkd:   "[hello](google.com)\n\n",
		},
		{
			name: "code snippet",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "`some code")
				require.NoError(t, err)
			},
			visIx:         10,
			text:          "`",
			orgSpanFormat: v3.FormatV3Span{},
			expectedErr:   nil,
			expectedMkd:   "`some code`\n\n",
		},
		{
			name: "insert tab at beginning of line",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "Hello World!")
				require.NoError(t, err)
			},
			visIx:         0,
			text:          "\t",
			orgSpanFormat: v3.FormatV3Span{},
			expectedErr:   nil,
			expectedMkd:   "\tHello World!\n\n",
		},
		{
			name: "indent multiple list items",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "a\nb\nc\nd")
				require.NoError(t, err)

				_, err = r.Format(0, 8, v3.FormatV3BulletList(0))
				require.NoError(t, err)
			},
			visIx:         2,
			selLen:        3,
			text:          "\t",
			orgSpanFormat: v3.FormatV3Span{},
			expectedErr:   nil,
			expectedMkd:   "- a\n  - b\n  - c\n- d\n",
		},
		{
			name: "newline from empty image caption",
			fn: func(r *v3.Rogue) {
				_, err := r.Format(0, 1, v3.FormatV3Image{
					Src: "https://www.fake.com/image.jpg",
				})
				require.NoError(t, err)
			},
			visIx:         0,
			selLen:        0,
			text:          "\n",
			orgSpanFormat: v3.FormatV3Span{},
			expectedErr:   nil,
			expectedMkd:   "<figure><img src=\"https://www.fake.com/image.jpg\" /><figcaption></figcaption></figure>\n\n\n\n",
			expectedHtml:  v3.PtrTo("<figure><img src=\"https://www.fake.com/image.jpg\" /><figcaption></figcaption></figure><p></p>"),
		},
		{
			name: "insert smart quote",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "\"")
				require.NoError(t, err)
			},
			visIx:            1,
			text:             "a",
			orgSpanFormat:    v3.FormatV3Span{},
			expectedMkd:      "&#34;a\n\n",
			expectedHtml:     v3.PtrTo(`<p data-rid="q_1"><span data-rid="0_3">“</span><span data-rid="0_4">a</span></p>`),
			expectedCursorID: &v3.ID{"q", 1},
			includeIDs:       true,
		},
		{
			name: "insert smart quote after newline",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "\n\"")
				require.NoError(t, err)
			},
			visIx:            2,
			text:             "a",
			orgSpanFormat:    v3.FormatV3Span{},
			expectedMkd:      "\n\n&#34;a\n\n",
			expectedHtml:     v3.PtrTo(`<p data-rid="0_3"></p><p data-rid="q_1"><span data-rid="0_4">“</span><span data-rid="0_5">a</span></p>`),
			expectedCursorID: &v3.ID{"q", 1},
			includeIDs:       true,
		},
		{
			name:             "check cursor",
			visIx:            0,
			text:             "abc",
			orgSpanFormat:    v3.FormatV3Span{},
			expectedMkd:      "abc\n\n",
			expectedHtml:     v3.PtrTo("<p><span>abc</span></p>"),
			expectedCursorID: &v3.ID{"q", 1},
		},
		{
			name: "bullet list shortcut",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "-a")
				require.NoError(t, err)
			},
			visIx:         1,
			text:          " ",
			orgSpanFormat: v3.FormatV3Span{},
			expectedErr:   nil,
			expectedMkd:   "- a\n",
			expectedHtml:  v3.PtrTo("<ul><li><span>a</span></li></ul>"),
		},
		{
			name: "bullet list shortcut only works at beginning of line",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "- a")
				require.NoError(t, err)
			},
			visIx:         3,
			text:          " ",
			orgSpanFormat: v3.FormatV3Span{},
			expectedErr:   nil,
			expectedMkd:   "\\- a \n\n",
			expectedHtml:  v3.PtrTo("<p><span>- a </span></p>"),
		},
		{
			name: "enter at start of nested list item dedents",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "a\nb")
				require.NoError(t, err)

				_, err = r.Format(0, 1, v3.FormatV3BulletList(0))
				require.NoError(t, err)

				_, err = r.Format(2, 1, v3.FormatV3BulletList(1))
				require.NoError(t, err)
			},
			visIx:         2,
			text:          "\n",
			orgSpanFormat: v3.FormatV3Span{},
			expectedErr:   nil,
			expectedMkd:   "- a\n- b\n",
			expectedHtml:  v3.PtrTo("<ul><li><span>a</span></li><li><span>b</span></li></ul>"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := v3.NewRogueForQuill("0")
			if tc.fn != nil {
				tc.fn(r)
			}

			_, cursorID, err := r.RichInsert(tc.visIx, tc.selLen, tc.orgSpanFormat, tc.text)
			if tc.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}

			if tc.expectedCursorID != nil {
				require.Equal(t, *tc.expectedCursorID, cursorID)
			}

			firstId, err := r.Rope.GetTotID(0)
			require.NoError(t, err)
			lastId, err := r.Rope.GetTotID(r.TotSize - 1)
			require.NoError(t, err)

			mkd, err := r.GetMarkdownBeforeAfter(firstId, lastId)
			require.NoError(t, err)

			if tc.expectedHtml != nil {
				html, err := r.GetHtml(firstId, lastId, tc.includeIDs, true)
				require.NoError(t, err)
				require.Equal(t, *tc.expectedHtml, html)
			}

			require.Equal(t, tc.expectedMkd, mkd)
		})
	}
}

func TestRichDelete(t *testing.T) {
	testCases := []struct {
		name         string
		fn           func(*v3.Rogue)
		visIx        int
		length       int
		expectedErr  error
		expectedHtml string
	}{
		{
			name: "delete with no formatting",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "abc")
				require.NoError(t, err)
			},
			visIx:        2,
			length:       1,
			expectedErr:  nil,
			expectedHtml: "<p><span>ab</span></p>",
		},
		{
			name: "delete \n in list",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "a\nb\nc\nd")
				require.NoError(t, err)
				_, err = r.Format(0, 8, v3.FormatV3BulletList(0))
				require.NoError(t, err)
			},
			visIx:        3,
			length:       1,
			expectedErr:  nil,
			expectedHtml: "<ul><li><span>a</span></li><li><span>b</span></li></ul><p><span>c</span></p><ul><li><span>d</span></li></ul>",
		},
		{
			name: "delete at start of doc clears formatting",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "a")
				require.NoError(t, err)
				_, err = r.Format(0, 1, v3.FormatV3Header(1))
				require.NoError(t, err)
			},
			visIx:        -1,
			length:       1,
			expectedErr:  nil,
			expectedHtml: "<p><span>a</span></p>",
		},
		{
			name: "delete at start of indented list item dedents",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "a\nb")
				require.NoError(t, err)

				_, err = r.Format(0, 1, v3.FormatV3BulletList(0))
				require.NoError(t, err)

				_, err = r.Format(2, 1, v3.FormatV3BulletList(1))
				require.NoError(t, err)
			},
			visIx:        1,
			length:       1,
			expectedErr:  nil,
			expectedHtml: "<ul><li><span>a</span></li><li><span>b</span></li></ul>",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := v3.NewRogueForQuill("0")
			if tc.fn != nil {
				tc.fn(r)
			}

			_, _, err := r.RichDelete(tc.visIx, tc.length)
			if tc.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.Equal(t, tc.expectedErr, err)
			}

			firstId, err := r.Rope.GetTotID(0)
			require.NoError(t, err)
			lastId, err := r.Rope.GetTotID(r.TotSize - 1)
			require.NoError(t, err)

			html, err := r.GetHtml(firstId, lastId, false, false)
			require.NoError(t, err)

			require.Equal(t, tc.expectedHtml, html)
		})
	}
}

func TestLinksNoSticky(t *testing.T) {
	r := v3.NewRogue("0")

	_, err := r.Insert(0, "google.com")
	require.NoError(t, err)

	_, err = r.Format(0, 10, v3.FormatV3Span{"a": "https://www.google.com"})
	require.NoError(t, err)

	_, _, err = r.RichDelete(6, 4)
	require.NoError(t, err)

	_, err = r.Insert(6, " search")
	require.NoError(t, err)

	firstID, err := r.Rope.GetTotID(0)
	require.NoError(t, err)

	lastID, err := r.Rope.GetTotID(r.TotSize - 1)
	require.NoError(t, err)

	html, err := r.GetHtml(firstID, lastID, false, false)
	require.NoError(t, err)

	require.Equal(t, "<p><a href=\"https://www.google.com\">google</a><span> search</span></p>", html)
}

func TestHorizontalRule(t *testing.T) {
	r := v3.NewRogueForQuill("0")

	_, err := r.Insert(0, "--")
	require.NoError(t, err)

	_, cursorID, err := r.RichInsert(2, 0, v3.FormatV3Span{}, "-")
	require.NoError(t, err)

	fmt.Printf("r.GetText(): %q\n", r.GetText())

	require.Equal(t, v3.ID{"0", 8}, cursorID)

	firstID, err := r.Rope.GetTotID(0)
	require.NoError(t, err)

	lastID, err := r.Rope.GetTotID(r.TotSize - 1)
	require.NoError(t, err)

	html, err := r.GetHtml(firstID, lastID, true, false)
	require.NoError(t, err)

	require.Equal(t, `<hr data-rid="q_1"/><p data-rid="0_8"></p>`, html)
}

func TestIsSpanAllLists(t *testing.T) {
	testCases := []struct {
		name     string
		fn       func(*v3.Rogue)
		visIx    int
		length   int
		expected bool
	}{
		{
			name: "plain bullet list",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "a\nb\nc\nd")
				require.NoError(t, err)
				_, err = r.Format(0, 8, v3.FormatV3BulletList(0))
				require.NoError(t, err)
			},
			visIx:    3,
			length:   1,
			expected: true,
		},
		{
			name: "full plain bullet list",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "a\nb\nc\nd")
				require.NoError(t, err)
				_, err = r.Format(0, 8, v3.FormatV3BulletList(0))
				require.NoError(t, err)
			},
			visIx:    0,
			length:   7,
			expected: true,
		},
		{
			name: "mixed bullet and ordered list",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "a\nb\nc\nd\ne\nf")
				require.NoError(t, err)
				_, err = r.Format(0, 4, v3.FormatV3BulletList(0))
				require.NoError(t, err)
				_, err = r.Format(4, 7, v3.FormatV3OrderedList(0))
				require.NoError(t, err)
			},
			visIx:    0,
			length:   11,
			expected: true,
		},
		{
			name: "list with non-list content at start",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "start\na\nb\nc")
				require.NoError(t, err)
				_, err = r.Format(6, 6, v3.FormatV3BulletList(0))
				require.NoError(t, err)
			},
			visIx:    0,
			length:   11,
			expected: false,
		},
		{
			name: "list with non-list content at end",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "a\nb\nc\nend")
				require.NoError(t, err)
				_, err = r.Format(0, 6, v3.FormatV3BulletList(0))
				require.NoError(t, err)
			},
			visIx:    0,
			length:   9,
			expected: false,
		},
		{
			name: "list with non-list content in middle",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "a\nb\nmiddle\nc\nd")
				require.NoError(t, err)
				_, err = r.Format(0, 4, v3.FormatV3BulletList(0))
				require.NoError(t, err)
				_, err = r.Format(11, 4, v3.FormatV3BulletList(0))
				require.NoError(t, err)
			},
			visIx:    0,
			length:   15,
			expected: false,
		},
		{
			name: "nested lists",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "a\n  b\n  c\nd")
				require.NoError(t, err)
				_, err = r.Format(0, 10, v3.FormatV3BulletList(0))
				require.NoError(t, err)
				_, err = r.Format(2, 6, v3.FormatV3BulletList(1))
				require.NoError(t, err)
			},
			visIx:    0,
			length:   10,
			expected: true,
		},
		{
			name: "single character non-list",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "a")
				require.NoError(t, err)
			},
			visIx:    0,
			length:   1,
			expected: false,
		},
		{
			name: "single item list",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "a")
				require.NoError(t, err)
				_, err = r.Format(0, 1, v3.FormatV3BulletList(0))
				require.NoError(t, err)
			},
			visIx:    0,
			length:   1,
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := v3.NewRogueForQuill("0")
			if tc.fn != nil {
				tc.fn(r)
			}
			startID, err := r.Rope.GetVisID(tc.visIx)
			require.NoError(t, err)
			endID, err := r.Rope.GetVisID(tc.visIx + tc.length - 1)
			require.NoError(t, err)

			actual, err := r.IsSpanAllLists(startID, endID)
			require.NoError(t, err)
			require.Equal(t, tc.expected, actual)
		})
	}
}
