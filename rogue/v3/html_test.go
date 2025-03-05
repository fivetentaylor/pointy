package v3_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	v3 "github.com/teamreviso/code/rogue/v3"
	"github.com/teamreviso/code/rogue/v3/testcases"
)

func TestGetHtml(t *testing.T) {
	t.Parallel()
	type testCase struct {
		name         string
		initialDoc   func(r *v3.Rogue) error
		start        v3.ID
		end          v3.ID
		expected     string
		expectedNoID string
	}

	testCases := []testCase{
		{
			name: "Only newline",
			initialDoc: func(r *v3.Rogue) error {
				return nil
			},
			start:        v3.ID{Author: "root", Seq: 0},
			end:          v3.ID{Author: "q", Seq: 1},
			expected:     "<p data-rid=\"q_1\"></p>",
			expectedNoID: "<p></p>",
		},
		{
			name: "Hello world!",
			initialDoc: func(r *v3.Rogue) error {
				_, err := r.Insert(0, "Hello World!")
				return err
			},
			start:        v3.ID{Author: "root", Seq: 0},
			end:          v3.ID{Author: "q", Seq: 1},
			expected:     "<p data-rid=\"q_1\"><span data-rid=\"auth0_3\">Hello World!</span></p>",
			expectedNoID: "<p><span>Hello World!</span></p>",
		},
		{
			name: "single line with span formatting",
			initialDoc: func(r *v3.Rogue) (err error) {
				_, err = r.Insert(0, "Hello World! awesome")
				require.NoError(t, err)
				_, err = r.Format(6, 5, v3.FormatV3Span{"b": "true"})
				require.NoError(t, err)
				_, err = r.Format(13, 7, v3.FormatV3Span{"i": "true"})
				require.NoError(t, err)
				_, err = r.Format(13, 7, v3.FormatV3Span{"u": "true"})
				require.NoError(t, err)
				return
			},
			start:        v3.ID{Author: "root", Seq: 0},
			end:          v3.ID{Author: "q", Seq: 1},
			expected:     "<p data-rid=\"q_1\"><span data-rid=\"auth0_3\">Hello </span><strong data-rid=\"auth0_9\">World</strong><span data-rid=\"auth0_14\">! </span><em><u data-rid=\"auth0_16\">awesome</u></em></p>",
			expectedNoID: "<p><span>Hello </span><strong>World</strong><span>! </span><em><u>awesome</u></em></p>",
		},
		{
			name: "subdoc with span formatting",
			initialDoc: func(r *v3.Rogue) (err error) {
				_, err = r.Insert(0, "Hello World! awesome")
				require.NoError(t, err)
				_, err = r.Format(6, 5, v3.FormatV3Span{"b": "true"})
				require.NoError(t, err)
				_, err = r.Format(13, 7, v3.FormatV3Span{"i": "true"})
				require.NoError(t, err)
				_, err = r.Format(13, 7, v3.FormatV3Span{"u": "true"})
				require.NoError(t, err)
				return
			},
			start:        v3.ID{Author: "auth0", Seq: 10},
			end:          v3.ID{Author: "q", Seq: 1},
			expected:     "<p data-rid=\"q_1\"><strong data-rid=\"auth0_10\">orld</strong><span data-rid=\"auth0_14\">! </span><em><u data-rid=\"auth0_16\">awesome</u></em></p>",
			expectedNoID: "<p><strong>orld</strong><span>! </span><em><u>awesome</u></em></p>",
		},
		{
			name: "subdoc with chopped line formatting",
			initialDoc: func(r *v3.Rogue) (err error) {
				_, err = r.Insert(0, "Hello World!\n")
				require.NoError(t, err)
				_, err = r.Format(12, 1, v3.FormatV3Header(4))
				require.NoError(t, err)
				return
			},
			start:        v3.ID{Author: "auth0", Seq: 8},
			end:          v3.ID{Author: "auth0", Seq: 11},
			expected:     "<p><span data-rid=\"auth0_8\"> Wor</span></p>",
			expectedNoID: "<p><span> Wor</span></p>",
		},
		{
			name: "single line with overlapping span formatting",
			initialDoc: func(r *v3.Rogue) (err error) {
				_, err = r.Insert(0, "Hello World! awesome")
				require.NoError(t, err)
				_, err = r.Format(6, 5, v3.FormatV3Span{"b": "true"})
				require.NoError(t, err)
				_, err = r.Format(13, 7, v3.FormatV3Span{"i": "true"})
				require.NoError(t, err)
				_, err = r.Format(8, 8, v3.FormatV3Span{"u": "true"})
				require.NoError(t, err)
				return
			},
			start:        v3.ID{Author: "root", Seq: 0},
			end:          v3.ID{Author: "q", Seq: 1},
			expected:     "<p data-rid=\"q_1\"><span data-rid=\"auth0_3\">Hello </span><strong data-rid=\"auth0_9\">Wo</strong><strong><u data-rid=\"auth0_11\">rld</u></strong><u data-rid=\"auth0_14\">! </u><em><u data-rid=\"auth0_16\">awe</u></em><em data-rid=\"auth0_19\">some</em></p>",
			expectedNoID: "<p><span>Hello </span><strong>Wo</strong><strong><u>rld</u></strong><u>! </u><em><u>awe</u></em><em>some</em></p>",
		},
		{
			name: "newline with deletes",
			initialDoc: func(r *v3.Rogue) (err error) {
				_, err = r.Insert(0, "Hello World!\nawesome")
				require.NoError(t, err)

				_, err = r.Format(6, 6, v3.FormatV3Span{"b": "true"})
				require.NoError(t, err)

				_, err = r.Format(12, 1, v3.FormatV3Header(2))
				require.NoError(t, err)
				return
			},
			start:        v3.ID{Author: "root", Seq: 0},
			end:          v3.ID{Author: "q", Seq: 1},
			expected:     "<h2 data-rid=\"auth0_15\"><span data-rid=\"auth0_3\">Hello </span><strong data-rid=\"auth0_9\">World!</strong></h2><p data-rid=\"q_1\"><span data-rid=\"auth0_16\">awesome</span></p>",
			expectedNoID: "<h2><span>Hello </span><strong>World!</strong></h2><p><span>awesome</span></p>",
		},
		{
			name: "single line with deleted span formatting",
			initialDoc: func(r *v3.Rogue) (err error) {
				_, err = r.Insert(0, "Hello World! awesome")
				require.NoError(t, err)
				_, err = r.Format(0, 11, v3.FormatV3Span{"b": "true"})
				require.NoError(t, err)
				_, err = r.Delete(6, 7)
				require.NoError(t, err)
				return
			},
			start:        v3.ID{Author: "root", Seq: 0},
			end:          v3.ID{Author: "q", Seq: 1},
			expected:     "<p data-rid=\"q_1\"><strong data-rid=\"auth0_3\">Hello </strong><span data-rid=\"auth0_16\">awesome</span></p>",
			expectedNoID: "<p><strong>Hello </strong><span>awesome</span></p>",
		},
		{
			name: "multiple lines with a list and indent",
			initialDoc: func(r *v3.Rogue) (err error) {
				_, err = r.Insert(0, "awesome\nstuff\nindented stuff")
				require.NoError(t, err)

				_, err = r.Format(7, 1, v3.FormatV3BulletList(0))
				require.NoError(t, err)

				_, err = r.Format(13, 1, v3.FormatV3BulletList(1))
				require.NoError(t, err)

				_, err = r.Format(28, 1, v3.FormatV3BulletList(2))
				require.NoError(t, err)
				return
			},
			start:        v3.ID{Author: "root", Seq: 0},
			end:          v3.ID{Author: "q", Seq: 1},
			expected:     "<ul><li data-rid=\"auth0_10\"><span data-rid=\"auth0_3\">awesome</span></li><ul><li data-rid=\"auth0_16\"><span data-rid=\"auth0_11\">stuff</span></li><ul><li data-rid=\"q_1\"><span data-rid=\"auth0_17\">indented stuff</span></li></ul></ul></ul>",
			expectedNoID: "<ul><li><span>awesome</span></li><ul><li><span>stuff</span></li><ul><li><span>indented stuff</span></li></ul></ul></ul>",
		},
		{
			name: "code block",
			initialDoc: func(r *v3.Rogue) (err error) {
				_, err = r.Insert(0, "print('Hello World!')\nx = 5\nprint(x)")
				require.NoError(t, err)

				_, err = r.Format(21, 1, v3.FormatV3CodeBlock("python"))
				require.NoError(t, err)

				_, err = r.Format(27, 1, v3.FormatV3CodeBlock("python"))
				require.NoError(t, err)

				_, err = r.Format(36, 1, v3.FormatV3CodeBlock("python"))
				require.NoError(t, err)

				return nil
			},
			start:        v3.ID{Author: "root", Seq: 0},
			end:          v3.ID{Author: "q", Seq: 1},
			expected:     "<pre><code data-rid=\"auth0_3\" class=\"language-python\" data-language=\"python\">print(&#39;Hello World!&#39;)\nx = 5\nprint(x)\n</code></pre>",
			expectedNoID: "<pre><code class=\"language-python\" data-language=\"python\">print(&#39;Hello World!&#39;)\nx = 5\nprint(x)\n</code></pre>",
		},
		{
			name: "blockquote",
			initialDoc: func(r *v3.Rogue) (err error) {
				_, err = r.Insert(0, "to be or not to be\nthat is the question")
				require.NoError(t, err)

				_, err = r.Format(18, 1, v3.FormatV3BlockQuote{})
				require.NoError(t, err)

				_, err = r.Format(39, 1, v3.FormatV3BlockQuote{})
				require.NoError(t, err)

				return nil
			},
			start:        v3.ID{Author: "root", Seq: 0},
			end:          v3.ID{Author: "q", Seq: 1},
			expected:     "<blockquote><p data-rid=\"auth0_21\"><span data-rid=\"auth0_3\">to be or not to be</span></p><p data-rid=\"q_1\"><span data-rid=\"auth0_22\">that is the question</span></p></blockquote>",
			expectedNoID: "<blockquote><p><span>to be or not to be</span></p><p><span>that is the question</span></p></blockquote>",
		},
		{
			name: "multiformat",
			initialDoc: func(r *v3.Rogue) (err error) {
				_, err = r.Insert(0, "Hello\nWorld!")
				require.NoError(t, err)

				_, err = r.Format(0, 12, v3.FormatV3Span{"b": "true", "i": "true", "u": "true"})
				require.NoError(t, err)

				return nil
			},
			start:        v3.ID{Author: "root", Seq: 0},
			end:          v3.ID{Author: "q", Seq: 1},
			expected:     "<p data-rid=\"auth0_8\"><strong><em><u data-rid=\"auth0_3\">Hello</u></em></strong></p><p data-rid=\"q_1\"><strong><em><u data-rid=\"auth0_9\">World!</u></em></strong></p>",
			expectedNoID: "<p><strong><em><u>Hello</u></em></strong></p><p><strong><em><u>World!</u></em></strong></p>",
		},
		{
			name: "mixed nested lists",
			initialDoc: func(r *v3.Rogue) (err error) {
				_, err = r.Insert(0, "awesome\nstuff\nindented stuff")
				require.NoError(t, err)

				_, err = r.Format(7, 1, v3.FormatV3OrderedList(0))
				require.NoError(t, err)

				_, err = r.Format(13, 1, v3.FormatV3BulletList(1))
				require.NoError(t, err)

				_, err = r.Format(28, 1, v3.FormatV3OrderedList(0))
				require.NoError(t, err)
				return
			},
			start:        v3.ID{Author: "root", Seq: 0},
			end:          v3.ID{Author: "q", Seq: 1},
			expected:     "<ol><li data-rid=\"auth0_10\"><span data-rid=\"auth0_3\">awesome</span></li><ul><li data-rid=\"auth0_16\"><span data-rid=\"auth0_11\">stuff</span></li></ul><li data-rid=\"q_1\"><span data-rid=\"auth0_17\">indented stuff</span></li></ol>",
			expectedNoID: "<ol><li><span>awesome</span></li><ul><li><span>stuff</span></li></ul><li><span>indented stuff</span></li></ol>",
		},
		{
			name: "crazy nested lists",
			initialDoc: func(r *v3.Rogue) (err error) {
				_, err = r.Insert(0, "a\nb\nc\nd\ne\nf\ng\nh")
				require.NoError(t, err)

				_, err = r.Format(0, 1, v3.FormatV3OrderedList(0))
				require.NoError(t, err)

				_, err = r.Format(2, 1, v3.FormatV3BulletList(1))
				require.NoError(t, err)

				_, err = r.Format(4, 1, v3.FormatV3OrderedList(0))
				require.NoError(t, err)

				_, err = r.Format(6, 1, v3.FormatV3OrderedList(1))
				require.NoError(t, err)

				_, err = r.Format(8, 1, v3.FormatV3OrderedList(0))
				require.NoError(t, err)

				_, err = r.Format(10, 1, v3.FormatV3OrderedList(1))
				require.NoError(t, err)

				_, err = r.Format(12, 1, v3.FormatV3OrderedList(2))
				require.NoError(t, err)

				_, err = r.Format(14, 1, v3.FormatV3OrderedList(2))
				require.NoError(t, err)
				return
			},
			start:        v3.ID{Author: "root", Seq: 0},
			end:          v3.ID{Author: "q", Seq: 1},
			expected:     "<ol><li data-rid=\"auth0_4\"><span data-rid=\"auth0_3\">a</span></li><ul><li data-rid=\"auth0_6\"><span data-rid=\"auth0_5\">b</span></li></ul><li data-rid=\"auth0_8\"><span data-rid=\"auth0_7\">c</span></li><ol><li data-rid=\"auth0_10\"><span data-rid=\"auth0_9\">d</span></li></ol><li data-rid=\"auth0_12\"><span data-rid=\"auth0_11\">e</span></li><ol><li data-rid=\"auth0_14\"><span data-rid=\"auth0_13\">f</span></li><ol><li data-rid=\"auth0_16\"><span data-rid=\"auth0_15\">g</span></li><li data-rid=\"q_1\"><span data-rid=\"auth0_17\">h</span></li></ol></ol></ol>",
			expectedNoID: "<ol><li><span>a</span></li><ul><li><span>b</span></li></ul><li><span>c</span></li><ol><li><span>d</span></li></ol><li><span>e</span></li><ol><li><span>f</span></li><ol><li><span>g</span></li><li><span>h</span></li></ol></ol></ol>",
		},
		{
			name: "indented text",
			initialDoc: func(r *v3.Rogue) (err error) {
				_, err = r.Insert(0, "a\nb\nc\nd")
				require.NoError(t, err)

				_, err = r.Format(1, 1, v3.FormatV3OrderedList(0))
				require.NoError(t, err)

				_, err = r.Format(3, 1, v3.FormatV3IndentedLine(0))
				require.NoError(t, err)

				_, err = r.Format(5, 1, v3.FormatV3IndentedLine(0))
				require.NoError(t, err)

				_, err = r.Format(7, 1, v3.FormatV3OrderedList(0))
				require.NoError(t, err)

				return nil
			},
			start:        v3.ID{Author: "root", Seq: 0},
			end:          v3.ID{Author: "q", Seq: 1},
			expected:     "<ol><li data-rid=\"auth0_4\"><span data-rid=\"auth0_3\">a</span></li><p data-rid=\"auth0_6\"><span data-rid=\"auth0_5\">b</span></p><p data-rid=\"auth0_8\"><span data-rid=\"auth0_7\">c</span></p><li data-rid=\"q_1\"><span data-rid=\"auth0_9\">d</span></li></ol>",
			expectedNoID: "<ol><li><span>a</span></li><p><span>b</span></p><p><span>c</span></p><li><span>d</span></li></ol>",
		},
		{
			name: "span null / empty",
			initialDoc: func(r *v3.Rogue) (err error) {
				_, err = r.Insert(0, "Hello World!")
				require.NoError(t, err)

				_, err = r.Format(0, 12, v3.FormatV3Span{"b": "true"})
				require.NoError(t, err)

				_, err = r.Format(0, 6, v3.FormatV3Span{"b": "null"})
				require.NoError(t, err)

				_, err = r.Format(6, 6, v3.FormatV3Span{"b": ""})
				require.NoError(t, err)

				return nil
			},
			start:        v3.ID{Author: "root", Seq: 0},
			end:          v3.ID{Author: "q", Seq: 1},
			expected:     "<p data-rid=\"q_1\"><span data-rid=\"auth0_3\">Hello World!</span></p>",
			expectedNoID: "<p><span>Hello World!</span></p>",
		},
		{
			name: "split list with non zero start indent",
			initialDoc: func(r *v3.Rogue) (err error) {
				_, err = r.Insert(0, "a\nb\nc\nd")
				require.NoError(t, err)

				_, err = r.Format(1, 1, v3.FormatV3OrderedList(0))
				require.NoError(t, err)

				_, err = r.Format(3, 1, v3.FormatV3OrderedList(1))
				require.NoError(t, err)

				_, err = r.Format(5, 1, v3.FormatV3OrderedList(1))
				require.NoError(t, err)

				_, err = r.Format(7, 1, v3.FormatV3OrderedList(0))
				require.NoError(t, err)

				_, err = r.Format(3, 1, v3.FormatV3Line{})
				require.NoError(t, err)

				return nil
			},
			start:        v3.ID{Author: "root", Seq: 0},
			end:          v3.ID{Author: "q", Seq: 1},
			expected:     "<ol><li data-rid=\"auth0_4\"><span data-rid=\"auth0_3\">a</span></li></ol><p data-rid=\"auth0_6\"><span data-rid=\"auth0_5\">b</span></p><ol><li data-rid=\"auth0_8\"><span data-rid=\"auth0_7\">c</span></li><li data-rid=\"q_1\"><span data-rid=\"auth0_9\">d</span></li></ol>",
			expectedNoID: "<ol><li><span>a</span></li></ol><p><span>b</span></p><ol><li><span>c</span></li><li><span>d</span></li></ol>",
		},
		{
			name: "get part of header",
			initialDoc: func(r *v3.Rogue) (err error) {
				_, err = r.Insert(0, "Hello World!")
				require.NoError(t, err)

				_, err = r.Format(1, 1, v3.FormatV3Header(1))
				require.NoError(t, err)

				return nil
			},
			start:        v3.ID{Author: "auth0", Seq: 3},
			end:          v3.ID{Author: "auth0", Seq: 7},
			expected:     "<p><span data-rid=\"auth0_3\">Hello</span></p>",
			expectedNoID: "<p><span>Hello</span></p>",
		},
		{
			name: "get part of header 2",
			initialDoc: func(r *v3.Rogue) (err error) {
				_, err = r.Insert(0, "Hello World!")
				require.NoError(t, err)

				_, err = r.Format(1, 1, v3.FormatV3Header(1))
				require.NoError(t, err)

				return nil
			},
			start:        v3.ID{Author: "auth0", Seq: 9},
			end:          v3.ID{Author: "auth0", Seq: 13},
			expected:     "<p><span data-rid=\"auth0_9\">World</span></p>",
			expectedNoID: "<p><span>World</span></p>",
		},
		{
			name: "image0",
			initialDoc: func(r *v3.Rogue) (err error) {
				_, err = r.Insert(0, "Hello World!")
				require.NoError(t, err)

				_, err = r.Format(0, 12, v3.FormatV3Image{Src: "https://example.com/image.jpg", Alt: "Hello World!", Width: "100px", Height: "100px"})
				require.NoError(t, err)

				return nil
			},
			start:        v3.ID{Author: "root", Seq: 0},
			end:          v3.ID{Author: "q", Seq: 1},
			expected:     `<figure><img src="https://example.com/image.jpg" style="width: 100px; height: 100px;" alt="Hello World!" /><figcaption data-rid="q_1"><span data-rid="auth0_3">Hello World!</span></figcaption></figure>`,
			expectedNoID: `<figure><img src="https://example.com/image.jpg" style="width: 100px; height: 100px;" alt="Hello World!" /><figcaption><span>Hello World!</span></figcaption></figure>`,
		},
		{
			name: "image1",
			initialDoc: func(r *v3.Rogue) (err error) {
				_, err = r.Insert(0, "Hello World!")
				require.NoError(t, err)

				_, err = r.Format(0, 12, v3.FormatV3Image{Src: "https://example.com/image.jpg"})
				require.NoError(t, err)

				_, err = r.Format(0, 5, v3.FormatV3Span{"b": "true"})
				require.NoError(t, err)

				_, err = r.Format(6, 5, v3.FormatV3Span{"u": "true"})
				require.NoError(t, err)

				return nil
			},
			start:        v3.ID{Author: "root", Seq: 0},
			end:          v3.ID{Author: "q", Seq: 1},
			expected:     `<figure><img src="https://example.com/image.jpg" /><figcaption data-rid="q_1"><strong data-rid="auth0_3">Hello</strong><span data-rid="auth0_8"> </span><u data-rid="auth0_9">World</u><span data-rid="auth0_14">!</span></figcaption></figure>`,
			expectedNoID: `<figure><img src="https://example.com/image.jpg" /><figcaption><strong>Hello</strong><span> </span><u>World</u><span>!</span></figcaption></figure>`,
		},
	}

	// Execute each test case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := v3.NewRogueForQuill("auth0")
			err := tc.initialDoc(r)
			if err != nil {
				t.Fatalf("Failed to initialize document: %v", err)
			}

			html, err := r.GetHtml(tc.start, tc.end, true, false)
			require.NoError(t, err)

			require.Equal(t, tc.expected, html)

			htmlNoID, err := r.GetHtml(tc.start, tc.end, false, false)
			require.NoError(t, err)

			require.Equal(t, tc.expectedNoID, htmlNoID)
		})
	}
}

/*
// had to drop the gold.json snapshot from tests
// because it appears to have had a overlapping lamport
// timestamps between deletes and inserts for the same
// author
func TestGetHTML_gold_snapshot(t *testing.T) {
	r := testcases.Load(t, "gold.json")

	html, err := r.GetHtml(v3.RootID, v3.LastID, true)
	if err != nil {
		t.Fatalf("Failed to get HTML: %s", err)
	}
	require.NotEmpty(t, html)
}*/

func TestGetHTML_draft_snapshot(t *testing.T) {
	r := testcases.Load(t, "reviso_draft.json")

	html, err := r.GetFullHtml(true, false)
	if err != nil {
		t.Fatalf("Failed to get HTML: %s", err)
	}
	require.NotEmpty(t, html)

	startID := v3.ID{"891d-tl", 10061}
	endID := v3.ID{"891d-tl", 10237}

	startID, err = r.Rope.NearestVisRightOf(startID)
	require.NoError(t, err)

	endID, err = r.Rope.NearestVisLeftOf(endID)
	require.NoError(t, err)

	// this had overlapping formatting
	html, err = r.GetHtml(
		startID,
		endID,
		true,
		false,
	)
	if err != nil {
		t.Fatalf("Failed to get HTML: %s", err)
	}
	require.NotEmpty(t, html)
}

func BenchmarkGetHTML(b *testing.B) {
	r := testcases.Load(b, "reviso_draft.json")

	for i := 0; i < b.N; i++ {
		_, err := r.GetHtml(v3.RootID, v3.LastID, true, false)
		if err != nil {
			b.Fatalf("Failed to get HTML: %s", err)
		}
	}
}

func TestGetHTMLV2(t *testing.T) {
	r := v3.NewRogueForQuill("auth0")

	r.Insert(0, "Hello World!\nGoodbye World!")
	r.Format(0, 5, v3.FormatV3Span{"b": "true"})
	r.Format(2, 5, v3.FormatV3Span{"i": "true"})
	r.Format(4, 5, v3.FormatV3Span{"u": "true"})

	html, err := r.GetHtml(v3.RootID, v3.LastID, true, false)
	require.NoError(t, err)

	fmt.Println(html)
}

func TestGetHTMLV2_newlines(t *testing.T) {
	r := v3.NewRogueForQuill("auth0")

	r.Insert(0, "Hello World!\nGoodbye World!")
	r.Format(0, 20, v3.FormatV3Span{"b": "true"})

	html, err := r.GetHtml(v3.RootID, v3.LastID, true, false)
	require.NoError(t, err)

	fmt.Println(html)
}

func TestGetHTMLV2_lineformat(t *testing.T) {
	r := v3.NewRogueForQuill("auth0")

	r.Insert(0, "Hello World!\nGoodbye World!")

	r.Format(0, 20, v3.FormatV3Span{"b": "true"})
	r.Format(12, 1, v3.FormatV3Header(1))
	r.Format(27, 1, v3.FormatV3BulletList(0))

	html, err := r.GetHtml(v3.RootID, v3.LastID, true, false)
	require.NoError(t, err)

	fmt.Println(html)
}

func TestGetHTMLV2_emojilists(t *testing.T) {
	r := testcases.Load(t, "emoji_lists.json")

	html, err := r.GetHtml(v3.ID{"0000", 83}, v3.ID{"0000", 109}, true, false)
	// html, err := r.GetHtml(v3.RootID, v3.LastID)
	require.NoError(t, err)

	fmt.Println(html)
}

func TestGetHtmlAtAddress(t *testing.T) {
	t.Parallel()
	type testCase struct {
		name       string
		initialDoc func(r *v3.Rogue) error
		start      v3.ID
		end        v3.ID
		address    v3.ContentAddress

		expected string
	}

	testCases := []testCase{
		{
			name: "Hello World",
			initialDoc: func(r *v3.Rogue) error {
				_, err := r.Insert(0, "Hello World ")
				require.NoError(t, err)

				address, err := r.GetAddress(v3.RootID, v3.LastID)
				fmt.Printf("Address: %v\n", address)
				require.NoError(t, err)

				_, err = r.Insert(12, "Goodbye World ")
				require.NoError(t, err)

				return nil
			},
			start: v3.ID{Author: "root", Seq: 0},
			end:   v3.ID{Author: "q", Seq: 1},
			address: v3.ContentAddress{
				StartID: v3.ID{Author: "root", Seq: 0},
				EndID:   v3.ID{Author: "q", Seq: 1},
				MaxIDs:  map[string]int{"0": 14, "q": 2, "root": 0},
			},
			expected: `<p data-rid="q_1"><span data-rid="0_3">Hello World </span></p>`,
		},
		{
			name: "Hello World mid insert",
			initialDoc: func(r *v3.Rogue) error {
				_, err := r.Insert(0, "Hello World ")
				require.NoError(t, err)

				address, err := r.GetAddress(v3.RootID, v3.LastID)
				fmt.Printf("Address: %v\n", address)
				require.NoError(t, err)

				_, err = r.Insert(12, "Goodbye World ")
				require.NoError(t, err)

				return nil
			},
			start: v3.ID{Author: "root", Seq: 0},
			end:   v3.ID{Author: "q", Seq: 1},
			address: v3.ContentAddress{
				StartID: v3.ID{Author: "root", Seq: 0},
				EndID:   v3.ID{Author: "q", Seq: 1},
				MaxIDs:  map[string]int{"0": 4, "q": 2, "root": 0},
			},
			expected: `<p data-rid="q_1"><span data-rid="0_3">He</span></p>`,
		},
		{
			name: "Multiple authors",
			initialDoc: func(r *v3.Rogue) error {
				_, err := r.Insert(0, "Hello ")
				require.NoError(t, err)

				r.Author = "1"
				_, err = r.Insert(6, "World!")
				require.NoError(t, err)

				r.Author = "0"

				address, err := r.GetAddress(v3.RootID, v3.LastID)
				fmt.Printf("Address: %v\n", address)
				require.NoError(t, err)

				_, err = r.Insert(12, "Goodbye World ")
				require.NoError(t, err)

				return nil
			},
			start: v3.ID{Author: "root", Seq: 0},
			end:   v3.ID{Author: "q", Seq: 1},
			address: v3.ContentAddress{
				StartID: v3.ID{Author: "root", Seq: 0},
				EndID:   v3.ID{Author: "q", Seq: 1},
				MaxIDs:  map[string]int{"0": 8, "1": 14, "q": 2, "root": 0},
			},
			expected: `<p data-rid="q_1"><span data-rid="0_3">Hello World!</span></p>`,
		},
		{
			name: "Formats",
			initialDoc: func(r *v3.Rogue) error {
				_, err := r.Insert(0, "Hello World!")
				require.NoError(t, err)

				_, err = r.Format(0, 12, v3.FormatV3Span{"b": "true"})
				require.NoError(t, err)

				address, err := r.GetAddress(v3.RootID, v3.LastID)
				fmt.Printf("Address: %v\n", address)
				require.NoError(t, err)

				_, err = r.Format(0, 12, v3.FormatV3Span{"b": ""})
				require.NoError(t, err)

				return nil
			},
			start: v3.ID{Author: "root", Seq: 0},
			end:   v3.ID{Author: "q", Seq: 1},
			address: v3.ContentAddress{
				StartID: v3.ID{Author: "root", Seq: 0},
				EndID:   v3.ID{Author: "q", Seq: 1},
				MaxIDs:  map[string]int{"0": 15, "q": 2, "root": 0},
			},
			expected: `<p data-rid="q_1"><strong data-rid="0_3">Hello World!</strong></p>`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := v3.NewRogueForQuill("0")
			err := tc.initialDoc(r)
			if err != nil {
				t.Fatalf("Failed to initialize document: %v", err)
			}

			html, err := r.GetHtmlAt(tc.start, tc.end, &tc.address, true, false)
			require.NoError(t, err)

			require.Equal(t, tc.expected, html)

		})
	}
}

func TestGetHtmlDiff(t *testing.T) {
	tests := []struct {
		name         string
		setupFunc    func(r *v3.Rogue, ca *v3.ContentAddress)
		expectedHtml string
	}{
		{
			name: "Basic Insert and Delete",
			setupFunc: func(r *v3.Rogue, ca *v3.ContentAddress) {
				_, err := r.Insert(0, "Hello World!")
				require.NoError(t, err)

				firstID, err := r.Rope.GetTotID(0)
				require.NoError(t, err)

				lastID, err := r.Rope.GetTotID(r.TotSize - 1)
				require.NoError(t, err)

				addr, err := r.GetAddress(firstID, lastID)
				require.NoError(t, err)

				*ca = *addr

				_, err = r.Delete(6, 6)
				require.NoError(t, err)

				_, err = r.Insert(6, "Friends!")
				require.NoError(t, err)
			},
			expectedHtml: "<p data-rid=\"q_1\"><span data-rid=\"0_3\">Hello </span><del data-delta-start=\"0_9\" data-delta-end=\"0_13\">World</del><ins data-rid=\"0_16\" data-delta-start=\"0_16\" data-delta-end=\"0_22\">Friends</ins><span data-rid=\"0_23\">!</span></p>",
		},
		{
			name: "With some formats",
			setupFunc: func(r *v3.Rogue, ca *v3.ContentAddress) {
				_, err := r.Insert(0, "Hello World!")
				require.NoError(t, err)

				_, err = r.Format(0, 5, v3.FormatV3Span{"b": "true"})
				require.NoError(t, err)

				firstID, err := r.Rope.GetTotID(0)
				require.NoError(t, err)

				lastID, err := r.Rope.GetTotID(r.TotSize - 1)
				require.NoError(t, err)

				addr, err := r.GetAddress(firstID, lastID)
				require.NoError(t, err)

				*ca = *addr

				_, err = r.Delete(6, 6)
				require.NoError(t, err)

				_, err = r.Insert(6, "Friends!")
				require.NoError(t, err)

				_, err = r.Format(6, 8, v3.FormatV3Span{"u": "true"})
				require.NoError(t, err)

			},
			expectedHtml: "<p data-rid=\"q_1\"><strong data-rid=\"0_3\">Hello</strong><span data-rid=\"0_8\"> </span><del data-delta-start=\"0_9\" data-delta-end=\"0_13\">World</del><ins data-delta-start=\"0_17\" data-delta-end=\"0_23\"><u data-rid=\"0_17\">Friends</u></ins><u data-rid=\"0_24\">!</u></p>",
		},
		{
			name: "Near deletes and inserts",
			setupFunc: func(r *v3.Rogue, ca *v3.ContentAddress) {
				_, err := r.Insert(0, "The quick brown fox jumps over the lazy dog.")
				require.NoError(t, err)

				firstID, err := r.Rope.GetTotID(0)
				require.NoError(t, err)

				lastID, err := r.Rope.GetTotID(r.TotSize - 1)
				require.NoError(t, err)

				addr, err := r.GetAddress(firstID, lastID)
				require.NoError(t, err)

				*ca = *addr

				_, _, err = r.ApplyMarkdownDiff("!1", "The slow green fox jumps over the excited dog.\n", firstID, lastID)
				require.NoError(t, err)

			},
			expectedHtml: "<p data-rid=\"q_1\"><span data-rid=\"0_3\">The </span><del data-delta-start=\"0_7\" data-delta-end=\"0_17\">quick brown</del><ins data-rid=\"!1_48\" data-delta-start=\"!1_48\" data-delta-end=\"!1_57\">slow green</ins><span data-rid=\"0_18\"> fox jumps over the </span><del data-delta-start=\"0_38\" data-delta-end=\"0_41\">lazy</del><ins data-rid=\"!1_59\" data-delta-start=\"!1_59\" data-delta-end=\"!1_65\">excited</ins><span data-rid=\"0_42\"> dog.</span></p>",
		},
		{
			name: "newline and added text",
			setupFunc: func(r *v3.Rogue, ca *v3.ContentAddress) {
				_, err := r.Insert(0, "Hello World! pal")
				require.NoError(t, err)

				firstID, err := r.Rope.GetTotID(0)
				require.NoError(t, err)

				lastID, err := r.Rope.GetTotID(r.TotSize - 1)
				require.NoError(t, err)

				addr, err := r.GetAddress(firstID, lastID)
				require.NoError(t, err)

				*ca = *addr

				_, err = r.Delete(12, 4)
				require.NoError(t, err)

				_, err = r.Insert(5, "\ncool")
				require.NoError(t, err)
			},
			expectedHtml: `<p data-rid="0_20"><span data-rid="0_3">Hello</span></p><p data-rid="q_1"><ins data-rid="0_21" data-delta-start="0_21" data-delta-end="0_24">cool</ins><span data-rid="0_8"> World!</span><del data-delta-start="0_15" data-delta-end="0_18"> pal</del></p>`,
		},
		{
			name: "delete lines",
			setupFunc: func(r *v3.Rogue, ca *v3.ContentAddress) {
				_, err := r.Insert(0, "Hello\nWorld!\nGoodbye\nWorld!")
				require.NoError(t, err)

				firstID, err := r.Rope.GetTotID(0)
				require.NoError(t, err)

				lastID, err := r.Rope.GetTotID(r.TotSize - 1)
				require.NoError(t, err)

				addr, err := r.GetAddress(firstID, lastID)
				require.NoError(t, err)

				*ca = *addr

				_, err = r.Delete(20, 7)
				require.NoError(t, err)

				_, err = r.Delete(5, 7)
				require.NoError(t, err)

			},
			expectedHtml: `<p data-rid="0_15"><span data-rid="0_3">Hello</span></p><p data-rid="q_1"><del data-delta-start="0_9" data-delta-end="0_15">World!<br></del><span data-rid="0_16">Goodbye</span></p><p data-rid="q_1"><del data-delta-start="0_24" data-delta-end="q_1">World!<br></del></p>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := v3.NewRogueForQuill("0")
			ca := v3.ContentAddress{}

			tt.setupFunc(r, &ca)

			firstID, err := r.Rope.GetTotID(0)
			require.NoError(t, err)

			lastID, err := r.Rope.GetTotID(r.TotSize - 1)
			require.NoError(t, err)

			fmt.Println(r.OpIndex)

			output, err := r.GetHtmlDiff(firstID, lastID, &ca, true, false)
			require.NoError(t, err)

			require.Equal(t, tt.expectedHtml, output)

			fmt.Printf("%q\n", output)
		})
	}
}

func TestGetHtmlDiffBetween(t *testing.T) {
	tests := []struct {
		name         string
		fn           func(r *v3.Rogue)
		startAddress v3.ContentAddress
		endAddress   v3.ContentAddress
		expectedHTML string
	}{
		{
			name: "From empty doc",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "Hello World!")
				require.NoError(t, err)
			},
			startAddress: v3.ContentAddress{StartID: v3.RootID, EndID: v3.LastID, MaxIDs: map[string]int{"q": 2, "root": 0}},
			endAddress:   v3.ContentAddress{StartID: v3.RootID, EndID: v3.LastID, MaxIDs: map[string]int{"1": 14, "q": 2, "root": 0}},
			expectedHTML: "<p data-rid=\"q_1\"><ins data-rid=\"1_3\" data-delta-start=\"1_3\" data-delta-end=\"1_14\">Hello World!</ins></p>",
		},
		{
			name: "Basic Insert and Delete",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "Hello World!")
				require.NoError(t, err)

				// add, err := r.GetFullAddress()
				// require.NoError(t, err)
				// fmt.Printf("%#v\n", add)

				_, err = r.Delete(6, 6)
				require.NoError(t, err)

				_, err = r.Insert(6, "Friends!")
				require.NoError(t, err)
			},
			startAddress: v3.ContentAddress{StartID: v3.RootID, EndID: v3.LastID, MaxIDs: map[string]int{"1": 14, "q": 2, "root": 0}},
			endAddress:   v3.ContentAddress{StartID: v3.RootID, EndID: v3.LastID, MaxIDs: map[string]int{"1": 23, "q": 2, "root": 0}},
			expectedHTML: "<p data-rid=\"q_1\"><span data-rid=\"1_3\">Hello </span><del data-delta-start=\"1_9\" data-delta-end=\"1_13\">World</del><ins data-rid=\"1_16\" data-delta-start=\"1_16\" data-delta-end=\"1_22\">Friends</ins><span data-rid=\"1_23\">!</span></p>",
		},
		{
			name: "Basic Insert and Delete after more inserts",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "Hello World!")
				require.NoError(t, err)

				_, err = r.Delete(6, 6)
				require.NoError(t, err)

				_, err = r.Insert(6, "Friends!")
				require.NoError(t, err)

				_, err = r.Insert(6, "I shouldn't be shown!")
				require.NoError(t, err)
			},
			startAddress: v3.ContentAddress{StartID: v3.RootID, EndID: v3.LastID, MaxIDs: map[string]int{"1": 14, "q": 2, "root": 0}},
			endAddress:   v3.ContentAddress{StartID: v3.RootID, EndID: v3.LastID, MaxIDs: map[string]int{"1": 23, "q": 2, "root": 0}},
			expectedHTML: `<p data-rid="q_1"><span data-rid="1_3">Hello </span><del data-delta-start="1_9" data-delta-end="1_13">World</del><ins data-rid="1_16" data-delta-start="1_16" data-delta-end="1_22">Friends</ins><span data-rid="1_23">!</span></p>`,
		},
		{
			name: "Formats",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "Hello World!")
				require.NoError(t, err)

				_, err = r.Format(0, 11, v3.FormatV3Header(2))
				require.NoError(t, err)

				_, err = r.Delete(6, 6)
				require.NoError(t, err)

				_, err = r.Insert(6, "Friends!")
				require.NoError(t, err)
			},
			startAddress: v3.ContentAddress{StartID: v3.RootID, EndID: v3.LastID, MaxIDs: map[string]int{"1": 14, "q": 2, "root": 0}},
			endAddress:   v3.ContentAddress{StartID: v3.RootID, EndID: v3.LastID, MaxIDs: map[string]int{"1": 24, "q": 2, "root": 0}},
			expectedHTML: `<h2 data-rid="q_1"><span data-rid="1_3">Hello </span><del data-delta-start="1_9" data-delta-end="1_13">World</del><ins data-rid="1_17" data-delta-start="1_17" data-delta-end="1_23">Friends</ins><span data-rid="1_24">!</span></h2>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := v3.NewRogueForQuill("1")
			tt.fn(r)

			firstID, err := r.GetFirstID()
			require.NoError(t, err)

			lastID, err := r.GetLastID()
			require.NoError(t, err)

			htmlBefore, err := r.GetHtmlAt(firstID, lastID, &tt.startAddress, false, false)
			require.NoError(t, err)
			fmt.Printf("htmlBefore: %q\n", htmlBefore)

			htmlAfter, err := r.GetHtmlAt(firstID, lastID, &tt.endAddress, false, false)
			require.NoError(t, err)
			fmt.Printf("htmlAfter: %q\n", htmlAfter)

			html, err := r.GetHtmlDiffBetween(v3.RootID, v3.LastID, &tt.startAddress, &tt.endAddress, true, false)
			require.NoError(t, err)
			require.Equal(t, tt.expectedHTML, html)
		})
	}

	t.Run("RW example", func(t *testing.T) {
		t.Skip()
		snapshot := `{"ops":[[0,["root",0],"x",["",0],0],[0,["q",1],"\n",["root",0],1],[1,["q",2],["root",0],1],[6,["1",3],[[0,["1",3],"Untitled\n",["q",1],-1],[6,["1",12],[[2,["1",12],["1",11],["1",11],{"h":"1"}]]]]],[6,["1",13],[[1,["1",13],["1",3],8],[0,["1",14],"H",["1",3],-1],[6,["1",15],[[2,["1",15],["1",14],["1",3],{"e":"true"}]]]]],[6,["1",16],[[0,["1",16],"e",["1",14],1],[6,["1",17],[[2,["1",17],["1",16],["1",3],{"e":"true"}]]]]],[6,["1",18],[[0,["1",18],"l",["1",16],1],[6,["1",19],[[2,["1",19],["1",18],["1",3],{"e":"true"}]]]]],[6,["1",20],[[0,["1",20],"l",["1",18],1],[6,["1",21],[[2,["1",21],["1",20],["1",3],{"e":"true"}]]]]],[6,["1",22],[[0,["1",22],"o",["1",20],1],[6,["1",23],[[2,["1",23],["1",22],["1",3],{"e":"true"}]]]]]]}`

		var r v3.Rogue
		err := json.Unmarshal([]byte(snapshot), &r)
		require.NoError(t, err)

		fmt.Printf("%q\n", r.GetText())

		firstID, err := r.GetFirstID()
		require.NoError(t, err)
		lastID, err := r.GetLastID()
		require.NoError(t, err)

		start := v3.ContentAddress{StartID: firstID, EndID: lastID, MaxIDs: map[string]int{"q": 2, "root": 0}}
		end := v3.ContentAddress{StartID: firstID, EndID: lastID, MaxIDs: map[string]int{"1": 23, "q": 2, "root": 0}}

		html, err := r.GetHtmlDiffBetween(firstID, lastID, &start, &end, true, false)
		require.NoError(t, err)
		require.NotEmpty(t, html)

		require.Equal(t, "<h1 data-rid=\"1_11\"><ins data-rid=\"1_14\" data-delta-start=\"1_14\" data-delta-end=\"1_22\">Hello</ins></h1><p data-rid=\"q_1\"></p>", html)
	})
}
