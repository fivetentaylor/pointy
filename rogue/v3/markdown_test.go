package v3_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	rogue "github.com/teamreviso/code/rogue/v3"
	"github.com/teamreviso/code/rogue/v3/testcases"
)

func TestGetMarkdown(t *testing.T) {
	type testCase struct {
		name           string
		setupDoc       func(doc *rogue.Rogue)
		startID        rogue.ID
		endID          rogue.ID
		expectedOutput string
	}

	tests := []testCase{
		{
			name: "nested list and header",
			setupDoc: func(doc *rogue.Rogue) {
				_, err := doc.Insert(0, "Hello World!\nIsn't this great?\nAnd a header")
				assert.NoError(t, err)

				_, err = doc.Format(12, 1, rogue.FormatV3BulletList(0))
				assert.NoError(t, err)

				_, err = doc.Format(30, 1, rogue.FormatV3BulletList(1))
				assert.NoError(t, err)

				_, err = doc.Format(43, 1, rogue.FormatV3Header(2))
				assert.NoError(t, err)
			},
			startID:        rogue.ID{"root", 0},
			endID:          rogue.ID{"q", 1},
			expectedOutput: "- Hello World!\n  - Isn&#39;t this great?\n## And a header\n\n",
		},
		{
			name: "spans and header",
			setupDoc: func(doc *rogue.Rogue) {
				_, err := doc.Insert(0, "Hello World!\nIsn't this great?")
				assert.NoError(t, err)

				_, err = doc.Format(0, 5, rogue.FormatV3Span{"b": "true"})
				assert.NoError(t, err)

				_, err = doc.Format(3, 7, rogue.FormatV3Span{"s": "true"})
				assert.NoError(t, err)

				_, err = doc.Format(30, 1, rogue.FormatV3Header(1))
				assert.NoError(t, err)
			},
			startID:        rogue.ID{"root", 0},
			endID:          rogue.ID{"q", 1},
			expectedOutput: "**Hel~~lo~~**~~ Worl~~d!\n\n# Isn&#39;t this great?\n\n",
		},
		{
			name: "partial span",
			setupDoc: func(doc *rogue.Rogue) {
				_, err := doc.Insert(0, "Hello World!")
				assert.NoError(t, err)

				_, err = doc.Format(0, 5, rogue.FormatV3Span{"b": "true"})
				assert.NoError(t, err)
			},
			startID:        rogue.ID{"auth0", 4},
			endID:          rogue.ID{"auth0", 13},
			expectedOutput: "**llo** Worl\n\n",
		},
		{
			name: "half bullet",
			setupDoc: func(doc *rogue.Rogue) {
				_, err := doc.Insert(0, "Hello World!\ncool beans")
				assert.NoError(t, err)

				_, err = doc.Format(12, 1, rogue.FormatV3BulletList(0))
				assert.NoError(t, err)

				_, err = doc.Format(23, 1, rogue.FormatV3BulletList(0))
				assert.NoError(t, err)
			},
			startID:        rogue.ID{"auth0", 8},
			endID:          rogue.ID{"auth0", 22},
			expectedOutput: "- World!\n- cool b\n",
		},
		{
			name: "emojis",
			setupDoc: func(doc *rogue.Rogue) {
				_, err := doc.Insert(0, "💀Hello World!\nIsn't this😅 great?\nAnd a header🤗")
				assert.NoError(t, err)

				_, err = doc.Format(14, 1, rogue.FormatV3BulletList(0))
				assert.NoError(t, err)

				_, err = doc.Format(21, 13, rogue.FormatV3Span{"b": "true"})
				assert.NoError(t, err)

				_, err = doc.Format(34, 1, rogue.FormatV3BulletList(1))
				assert.NoError(t, err)

				_, err = doc.Format(49, 1, rogue.FormatV3Header(2))
				assert.NoError(t, err)
			},
			startID:        rogue.ID{"root", 0},
			endID:          rogue.ID{"q", 1},
			expectedOutput: "- 💀Hello World!\n  - Isn&#39;t **this😅 great?**\n## And a header🤗\n\n",
		},
		{
			name: "partial deleted span styles",
			setupDoc: func(doc *rogue.Rogue) {
				_, err := doc.Insert(0, "Hello World!\nIsn't this great!")
				require.NoError(t, err)

				_, err = doc.Format(19, 4, rogue.FormatV3Span{"b": "true"})
				require.NoError(t, err)

				_, err = doc.Delete(14, 8)
				require.NoError(t, err)
			},
			startID:        rogue.ID{"root", 0},
			endID:          rogue.ID{"q", 1},
			expectedOutput: "Hello World!\n\nI**s** great!\n\n",
		},
		{
			name: "fully deleted span styles",
			setupDoc: func(doc *rogue.Rogue) {
				_, err := doc.Insert(0, "Hello World!\nIsn't this great!")
				require.NoError(t, err)

				_, err = doc.Format(19, 4, rogue.FormatV3Span{"b": "true"})
				require.NoError(t, err)

				_, err = doc.Delete(13, 11)
				require.NoError(t, err)
			},
			startID:        rogue.ID{"root", 0},
			endID:          rogue.ID{"q", 1},
			expectedOutput: "Hello World!\n\ngreat!\n\n",
		},
		{
			name: "spans over newlines",
			setupDoc: func(doc *rogue.Rogue) {
				_, err := doc.Insert(0, "Hello\nWorld!")
				require.NoError(t, err)

				_, err = doc.Format(0, 12, rogue.FormatV3Span{"b": "true"})
				require.NoError(t, err)

			},
			startID:        rogue.ID{"root", 0},
			endID:          rogue.ID{"q", 1},
			expectedOutput: "**Hello**\n\n**World!**\n\n",
		},
		{
			name: "improperly formatted newline",
			setupDoc: func(doc *rogue.Rogue) {
				_, err := doc.Insert(0, "Hello\nWorld!")
				require.NoError(t, err)

				_, err = doc.Format(5, 1, rogue.FormatV3Span{"b": "true"})
				require.NoError(t, err)

			},
			startID:        rogue.ID{"root", 0},
			endID:          rogue.ID{"q", 1},
			expectedOutput: "Hello\n\nWorld!\n\n",
		},
		{
			name: "bold test",
			setupDoc: func(doc *rogue.Rogue) {
				_, err := doc.Insert(0, "Hello World!")
				require.NoError(t, err)

				_, err = doc.Format(0, 12, rogue.FormatV3Span{"b": "true"})
				require.NoError(t, err)

				_, err = doc.Format(0, 6, rogue.FormatV3Span{"b": "null"})
				require.NoError(t, err)
			},
			startID:        rogue.ID{"root", 0},
			endID:          rogue.ID{"q", 1},
			expectedOutput: "Hello **World!**\n\n",
		},
		{
			name: "format newline bug",
			setupDoc: func(doc *rogue.Rogue) {
				*doc = *testcases.Load(t, "format_newline_bug.json")
			},
			startID:        rogue.ID{"root", 0},
			endID:          rogue.ID{"q", 1},
			expectedOutput: "Call me Ishmael.\n\n\n\n- ❤️\n- 🤗\n- 💀\n- 😅\n\n\n",
		},
		{
			name: "span format list",
			setupDoc: func(doc *rogue.Rogue) {
				_, err := doc.Insert(0, "one\ntwo\nthree")
				require.NoError(t, err)

				for _, i := range []int{3, 7, 13} {
					_, err = doc.Format(i, 1, rogue.FormatV3OrderedList(0))
					require.NoError(t, err)
				}

				_, err = doc.Format(0, 13, rogue.FormatV3Span{"b": "true"})
				require.NoError(t, err)
			},
			startID:        rogue.ID{"root", 0},
			endID:          rogue.ID{"q", 1},
			expectedOutput: "1. **one**\n1. **two**\n1. **three**\n",
		},
		{
			name: "code block",
			setupDoc: func(doc *rogue.Rogue) {
				_, err := doc.Insert(0, "for i in range(10):\n    print(i)")
				require.NoError(t, err)

				_, err = doc.Format(19, 1, rogue.FormatV3CodeBlock("python"))
				require.NoError(t, err)

				_, err = doc.Format(32, 1, rogue.FormatV3CodeBlock("python"))
				require.NoError(t, err)
			},
			startID:        rogue.ID{"root", 0},
			endID:          rogue.ID{"q", 1},
			expectedOutput: "```python\nfor i in range(10):\n    print(i)\n```\n",
		},
		{
			name: "blockquote",
			setupDoc: func(doc *rogue.Rogue) {
				_, err := doc.Insert(0, "to be\nor not to be\nthat is the question")
				require.NoError(t, err)

				_, err = doc.Format(5, 1, rogue.FormatV3BlockQuote{})
				require.NoError(t, err)

				_, err = doc.Format(18, 1, rogue.FormatV3BlockQuote{})
				require.NoError(t, err)

				_, err = doc.Format(39, 1, rogue.FormatV3BlockQuote{})
				require.NoError(t, err)
			},
			startID:        rogue.ID{"root", 0},
			endID:          rogue.ID{"q", 1},
			expectedOutput: "> to be\n> or not to be\n> that is the question\n",
		},
		{
			name: "neighboring spans",
			setupDoc: func(doc *rogue.Rogue) {
				_, err := doc.Insert(0, "HelloWorld!")
				require.NoError(t, err)

				_, err = doc.Format(0, 5, rogue.FormatV3Span{"s": "true"})
				require.NoError(t, err)

				_, err = doc.Format(5, 6, rogue.FormatV3Span{"i": "true"})
				require.NoError(t, err)
			},
			startID:        rogue.ID{"root", 0},
			endID:          rogue.ID{"q", 1},
			expectedOutput: "~~Hello~~*World!*\n\n",
		},
		{
			name: "XML content",
			setupDoc: func(doc *rogue.Rogue) {
				_, err := doc.Insert(0, "Hello <b>World!</b> World!")
				require.NoError(t, err)
			},
			startID:        rogue.ID{"root", 0},
			endID:          rogue.ID{"q", 1},
			expectedOutput: "Hello &lt;b&gt;World!&lt;/b&gt; World!\n\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			doc := rogue.NewRogueForQuill("auth0")

			// Setup the document
			tc.setupDoc(doc)

			// Get the Markdown
			md, err := doc.GetMarkdownBeforeAfter(tc.startID, tc.endID)
			assert.NoError(t, err)

			// Assert the Markdown output
			assert.Equal(t, tc.expectedOutput, md)

			// Optional: Print the Markdown for manual verification
			/*fmt.Println("MARKDOWN")
			fmt.Println(md)*/
		})
	}
}

func TestEscapeUnescapeMarkdownSyms(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		escaped string
	}{
		{
			name:    "Simple",
			input:   "*markdown* _text_",
			escaped: "\\*markdown\\* \\_text\\_",
		},
		{
			name:    "Complex",
			input:   "Example *markdown* text with _various_ ~characters~ [to] (escape), including #, +, -, ., and !",
			escaped: "Example \\*markdown\\* text with \\_various\\_ \\~characters\\~ [to] (escape), including #, +, -, ., and !",
		},
		{
			name:    "Complex Lines",
			input:   "# Hello World!\n  - Item 1\n  - Item 2\n  - Item 3\n> a quote > with a quote\n\n```python\nfor i in range(10):\n    print(i)\n```\n",
			escaped: "\\# Hello World!\n\\  - Item 1\n\\  - Item 2\n\\  - Item 3\n\\> a quote > with a quote\n\n\\```python\nfor i in range(10):\n    print(i)\n\\```\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			escaped := rogue.EscapeMarkdownSyms(tt.input)
			require.Equal(t, tt.escaped, escaped)

			unescaped := rogue.UnescapeMarkdownSyms(escaped)
			require.Equal(t, tt.input, unescaped)
		})
	}
}
