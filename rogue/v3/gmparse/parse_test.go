package gmparse_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/fivetentaylor/pointy/rogue/v3/gmparse"
)

func TestSplitMarkdown(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedText    string
		expectedActions []gmparse.FormatSpan
	}{
		{
			name:         "escape",
			input:        "# ***He~~ll~~o***, *world*! How are you? ~~I'm fine~~. \\*\\*I'm not bold\n  - A list\n",
			expectedText: "Hello, world! How are you? I'm fine. \\*\\*I'm not bold\nA list\n",
		},
		{
			name:            "Basic",
			input:           "Hello World!",
			expectedText:    "Hello World!",
			expectedActions: []gmparse.FormatSpan{},
		},
		{
			name:         "Bold",
			input:        "**Text**",
			expectedText: "Text",
			expectedActions: []gmparse.FormatSpan{
				{Start: 0, End: 4, Format: gmparse.Format{"b": "true"}},
			},
		},
		{
			name:         "Italic",
			input:        "*Text*",
			expectedText: "Text",
			expectedActions: []gmparse.FormatSpan{
				{Start: 0, End: 4, Format: gmparse.Format{"i": "true"}},
			},
		},
		{
			name:         "Strikethrough",
			input:        "~~Text~~",
			expectedText: "Text",
			expectedActions: []gmparse.FormatSpan{
				{Start: 0, End: 4, Format: gmparse.Format{"s": "true"}},
			},
		},
		{
			name:         "Headers",
			input:        "# Title\n## Subtitle\n### Subsubtitle",
			expectedText: "Title\nSubtitle\nSubsubtitle\n",
			expectedActions: []gmparse.FormatSpan{
				{Start: 5, End: 5, Format: gmparse.Format{"h": "1"}},
				{Start: 14, End: 14, Format: gmparse.Format{"h": "2"}},
				{Start: 26, End: 26, Format: gmparse.Format{"h": "3"}},
			},
		},
		{
			name:         "Ordered List",
			input:        "1. Item 1\n2. Item 2\n3. Item 3",
			expectedText: "Item 1\nItem 2\nItem 3\n",
			expectedActions: []gmparse.FormatSpan{
				{Start: 6, End: 6, Format: gmparse.Format{"ol": "0"}},
				{Start: 13, End: 13, Format: gmparse.Format{"ol": "0"}},
				{Start: 20, End: 20, Format: gmparse.Format{"ol": "0"}},
			},
		},
		{
			name:         "Unordered List",
			input:        "- Item 1\n- Item 2\n- Item 3",
			expectedText: "Item 1\nItem 2\nItem 3\n",
			expectedActions: []gmparse.FormatSpan{
				{Start: 6, End: 6, Format: gmparse.Format{"ul": "0"}},
				{Start: 13, End: 13, Format: gmparse.Format{"ul": "0"}},
				{Start: 20, End: 20, Format: gmparse.Format{"ul": "0"}},
			},
		},
		{
			name:         "Code Block 2",
			input:        "```python\nprint('Hello, World!')\n```",
			expectedText: "print('Hello, World!')\n",
			expectedActions: []gmparse.FormatSpan{
				{Start: 22, End: 22, Format: gmparse.Format{"cb": "python"}},
			},
		},
		{
			name:         "Blockquote",
			input:        "> To be or not to be\n> That is the question",
			expectedText: "To be or not to be\nThat is the question\n",
			expectedActions: []gmparse.FormatSpan{
				{Start: 18, End: 18, Format: gmparse.Format{"bq": "true"}},
				{Start: 39, End: 39, Format: gmparse.Format{"bq": "true"}},
			},
		},
		{
			name:         "Links",
			input:        "This is [a link](https://example.com)",
			expectedText: "This is a link",
			expectedActions: []gmparse.FormatSpan{
				{Start: 8, End: 14, Format: gmparse.Format{"a": "https://example.com"}},
			},
		},
		{
			name:         "Complex",
			input:        "# Welcome to Our Project! 🌎\n\nHello [World](https://google.com)! This is the beginning of something great.\n\n## Quick Start Guide\n\n1. **Introduction**\n   - Welcome to our project!\n2. **Setup**\n   ```bash\n   python --version\n   ```\n3. **Contribution Guidelines**\n",
			expectedText: "Welcome to Our Project! 🌎\nHello World! This is the beginning of something great.\nQuick Start Guide\nIntroduction\nWelcome to our project!\nSetup\npython --version\nContribution Guidelines\n",
		},
		{
			name:            "Empty String",
			input:           "",
			expectedText:    "",
			expectedActions: nil,
		},
		{
			name: "complex",
			input: `# Title
## Subtitle

This is a **sentence** about something

1. ~Item 1~
1. Item 2
1. Item 3
1. Item 4
`,
			expectedText: "Title\nSubtitle\nThis is a sentence about something\nItem 1\nItem 2\nItem 3\nItem 4\n",
		},
		{
			name: "emojis",
			input: `# Title🙌
## Subtitle

This is a **sente😅nce** about something

1. ~~Item 1~~
1. Item 2💀
1. 🥴Item 3
1. Item 4
`,
			expectedText: "Title🙌\nSubtitle\nThis is a sente😅nce about something\nItem 1\nItem 2💀\n🥴Item 3\nItem 4\n",
		},
		{
			name:         "punctuation",
			input:        `This, is. Some! Text?`,
			expectedText: "This, is. Some! Text?",
		},
		{
			name:         "newlines",
			input:        "This\nis\nSome!\nText?",
			expectedText: "This\nis\nSome!\nText?",
		},
		{
			name:         "whitespace",
			input:        "This   is \n\n Some   Text   \n",
			expectedText: "This   is\nSome   Text\n",
		},
		{
			name:            "Empty with Spaces",
			input:           "   ",
			expectedText:    "",
			expectedActions: []gmparse.FormatSpan{},
		},
		{
			name:         "Empty Strikethrough",
			input:        " ~~ ~~ ",
			expectedText: "~~ ~~",
		},
		{
			name:            "Bold with Whitespace",
			input:           " ** abc ** ",
			expectedText:    "** abc **",
			expectedActions: []gmparse.FormatSpan{},
		},
		{
			name:         "Empty List Item",
			input:        "- Item 1\n- ",
			expectedText: "Item 1\n",
			expectedActions: []gmparse.FormatSpan{
				{Start: 6, End: 6, Format: gmparse.Format{"ul": "0"}},
			},
		},
		{
			name: "Complex Document Structure",
			input: `# Title
## Subtitle
### Subsubtitle
Hello just some normal text
1. Item 1
1. Item 2
1. Item 3
`,
			expectedText: "Title\nSubtitle\nSubsubtitle\nHello just some normal text\nItem 1\nItem 2\nItem 3\n",
			expectedActions: []gmparse.FormatSpan{
				{Start: 5, End: 5, Format: gmparse.Format{"h": "1"}},
				{Start: 14, End: 14, Format: gmparse.Format{"h": "2"}},
				{Start: 26, End: 26, Format: gmparse.Format{"h": "3"}},
				{Start: 61, End: 61, Format: gmparse.Format{"ol": "0"}},
				{Start: 68, End: 68, Format: gmparse.Format{"ol": "0"}},
				{Start: 75, End: 75, Format: gmparse.Format{"ol": "0"}}},
		},
		{
			name: "Nested Ordered List",
			input: `1. Item 1
1. Item 2
   1. Item 3
`,
			expectedText: "Item 1\nItem 2\nItem 3\n",
			expectedActions: []gmparse.FormatSpan{
				{Start: 6, End: 6, Format: gmparse.Format{"ol": "0"}},
				{Start: 13, End: 13, Format: gmparse.Format{"ol": "0"}},
				{Start: 20, End: 20, Format: gmparse.Format{"ol": "1"}},
			},
		},
		{
			name: "Nested Unordered List",
			input: `- Item 1   
- Item 2
  - SubItem 1
- Item 3
`,
			expectedText: "Item 1\nItem 2\nSubItem 1\nItem 3\n",
			expectedActions: []gmparse.FormatSpan{
				{Start: 6, End: 6, Format: gmparse.Format{"ul": "0"}},
				{Start: 13, End: 13, Format: gmparse.Format{"ul": "0"}},
				{Start: 23, End: 23, Format: gmparse.Format{"ul": "1"}},
				{Start: 30, End: 30, Format: gmparse.Format{"ul": "0"}},
			},
		},
		{
			// TODO: actually support checklists
			name: "Checklist",
			input: `- [ ] Item 1
- [x] Item 2
- [ ] Item 3
`,
			expectedText: "Item 1\nItem 2\nItem 3\n",
			expectedActions: []gmparse.FormatSpan{
				{Start: 6, End: 6, Format: gmparse.Format{"ul": "0"}},
				{Start: 13, End: 13, Format: gmparse.Format{"ul": "0"}},
				{Start: 20, End: 20, Format: gmparse.Format{"ul": "0"}},
			},
		},
		{
			name: "RealWorld Example",
			input: `# Daily Planner Template
## Date: \[Insert Date Here\]
## Main Goals
- Goal 1
- Goal 2
- Goal 3
## Schedule
Time | Activity
---|---
7:00 AM | Example: Morning Exercise
8:00 AM | Example: Breakfast
... | ...
## To-Do List
- \[ \] Task 1
- \[ \] Task 2
- \[ \] Task 3
## Notes
- Note 1
- Note 2
## Reflection
- Today's achievements:
- Areas for improvement:
`,
			expectedText: `Daily Planner Template
Date: \[Insert Date Here\]
Main Goals
Goal 1
Goal 2
Goal 3
Schedule
Time | Activity
---|---
7:00 AM | Example: Morning Exercise
8:00 AM | Example: Breakfast
... | ...
To-Do List
\[ \] Task 1
\[ \] Task 2
\[ \] Task 3
Notes
Note 1
Note 2
Reflection
Today's achievements:
Areas for improvement:
`,
		},
		{
			name:         "Code Block",
			input:        "```python\nfor i in range(10):\nprint('Hello, World!' + str(i))\n```\nHello World!\n# Hello Header\n",
			expectedText: "for i in range(10):\nprint('Hello, World!' + str(i))\nHello World!\nHello Header\n",
			expectedActions: []gmparse.FormatSpan{
				{Start: 19, End: 19, Format: gmparse.Format{"cb": "python"}},
				{Start: 51, End: 51, Format: gmparse.Format{"cb": "python"}},
				{Start: 77, End: 77, Format: gmparse.Format{"h": "1"}},
			},
		},
		{
			name:         "Blockquote with bold",
			input:        "> To **be**\n> or not to be\n",
			expectedText: "To be\nor not to be\n",
			expectedActions: []gmparse.FormatSpan{
				{Start: 3, End: 5, Format: gmparse.Format{"b": "true"}},
				{Start: 5, End: 5, Format: gmparse.Format{"bq": "true"}},
				{Start: 18, End: 18, Format: gmparse.Format{"bq": "true"}},
			},
		},
		{
			name:         "Links with Bold",
			input:        "This is [my **link**](https://google.com) and [another](https://bing.com)",
			expectedText: "This is my link and another",
			expectedActions: []gmparse.FormatSpan{
				{Start: 11, End: 15, Format: gmparse.Format{"b": "true"}},
				{Start: 8, End: 15, Format: gmparse.Format{"a": "https://google.com"}},
				{Start: 20, End: 27, Format: gmparse.Format{"a": "https://bing.com"}},
			},
		},
		{
			name:         "Complex RealWorld Example",
			input:        "# Welcome to Our Project! 🌎\n\nHello [World](https://google.com)! This is the beginning of something great. Below, you'll find a quick guide on what to expect and how to get started.\n\n## Quick Start Guide\n\n1. **Introduction**\n  - Welcome to our project! We're glad you're here.\n1. **Setup**\n  - Ensure you have the necessary tools installed.\n  - For example, to check your Python version, use the following command in your terminal:\n```bash\n    python --version\n```\n3. **Contribution Guidelines**\n\n  - Read through our contribution guidelines to understand how you can contribute effectively.\n4. **Community**\n\n  - Join our community on [GitHub](https://github.com) or [Discord](https://scord.com) to stay updated and collaborate.\nWe're excited to have you on board and can't wait to see what we'll achieve together!",
			expectedText: "Welcome to Our Project! 🌎\nHello World! This is the beginning of something great. Below, you'll find a quick guide on what to expect and how to get started.\nQuick Start Guide\nIntroduction\nWelcome to our project! We're glad you're here.\nSetup\nEnsure you have the necessary tools installed.\nFor example, to check your Python version, use the following command in your terminal:\n    python --version\nContribution Guidelines\nRead through our contribution guidelines to understand how you can contribute effectively.\nCommunity\nJoin our community on GitHub or Discord to stay updated and collaborate.\nWe're excited to have you on board and can't wait to see what we'll achieve together!\n",
		},
		{
			name:         "Collapsed List Item with Bold",
			input:        "  - **Cookie Crisp** - Let the simple pleasure of enjoying a bowl of\n",
			expectedText: "Cookie Crisp - Let the simple pleasure of enjoying a bowl of\n",
		},
		{
			name:         "All Formats",
			input:        "# Header\n\n- [ ] today's fast-paced, [[link ][]](http://google.com) world, [[ ][people.\n\n## Key ication\n\n- **Clarity**: understandable.\n- *Emphasis*: key points.\n- `Listening`: communication.\n\n> \"it has taken place.\" \n> — **George** Bernard Shaw\n\n### Challenges\n\n1. Information\n2. rapid\n3. barriers\n\n```For more insights into the art of communication, consider exploring various resources and engaging in conversations that challenge and expand your perspectives.```",
			expectedText: "Header\ntoday's fast-paced, [link ][] world, [[ ][people.\nKey ication\nClarity: understandable.\nEmphasis: key points.\nListening: communication.\n\"it has taken place.\"\n— George Bernard Shaw\nChallenges\nInformation\nrapid\nbarriers\nFor more insights into the art of communication, consider exploring various resources and engaging in conversations that challenge and expand your perspectives.",
			expectedActions: []gmparse.FormatSpan{
				{Start: 6, End: 6, Format: gmparse.Format{"h": "1"}},
				{Start: 27, End: 36, Format: gmparse.Format{"a": "http://google.com"}},
				{Start: 56, End: 56, Format: gmparse.Format{"ul": "0"}},
				{Start: 68, End: 68, Format: gmparse.Format{"h": "2"}},
				{Start: 69, End: 76, Format: gmparse.Format{"b": "true"}},
				{Start: 93, End: 93, Format: gmparse.Format{"ul": "0"}},
				{Start: 94, End: 102, Format: gmparse.Format{"i": "true"}},
				{Start: 115, End: 115, Format: gmparse.Format{"ul": "0"}},
				{Start: 141, End: 141, Format: gmparse.Format{"ul": "0"}},
				{Start: 163, End: 163, Format: gmparse.Format{"bq": "true"}},
				{Start: 168, End: 174, Format: gmparse.Format{"b": "true"}},
				{Start: 187, End: 187, Format: gmparse.Format{"bq": "true"}},
				{Start: 198, End: 198, Format: gmparse.Format{"h": "3"}},
				{Start: 210, End: 210, Format: gmparse.Format{"ol": "0"}},
				{Start: 216, End: 216, Format: gmparse.Format{"ol": "0"}},
				{Start: 225, End: 225, Format: gmparse.Format{"ol": "0"}},
			},
		},
		{
			name:            "Empty String",
			input:           "",
			expectedText:    "",
			expectedActions: nil,
		},
		{
			name:         "Horizontal Rule",
			input:        "# Header\n\n---\n\n## Subheader\n\n",
			expectedText: "Header\n\nSubheader\n",
			expectedActions: []gmparse.FormatSpan{
				{Start: 6, End: 6, Format: gmparse.Format{"h": "1"}},
				{Start: 7, End: 7, Format: gmparse.Format{"r": "true"}},
				{Start: 17, End: 17, Format: gmparse.Format{"h": "2"}},
			},
		},
		{
			name:         "Complex Inline Formatting",
			input:        "~*__one__*~[**link**](http://google.com)\n",
			expectedText: "onelink\n",
			expectedActions: []gmparse.FormatSpan{
				{Start: 0, End: 3, Format: gmparse.Format{"b": "true"}},
				{Start: 0, End: 3, Format: gmparse.Format{"i": "true"}},
				{Start: 0, End: 3, Format: gmparse.Format{"s": "true"}},
				{Start: 3, End: 7, Format: gmparse.Format{"b": "true"}},
				{Start: 3, End: 7, Format: gmparse.Format{"a": "http://google.com"}},
			},
		},
		{
			name:         "Multiple Blockquotes with Formatting",
			input:        ">one **bold**\n>two ~~strike~~\n>three *italic*\n\n",
			expectedText: "one bold\ntwo strike\nthree italic\n",
			expectedActions: []gmparse.FormatSpan{
				{Start: 4, End: 8, Format: gmparse.Format{"b": "true"}},
				{Start: 8, End: 8, Format: gmparse.Format{"bq": "true"}},
				{Start: 13, End: 19, Format: gmparse.Format{"s": "true"}},
				{Start: 19, End: 19, Format: gmparse.Format{"bq": "true"}},
				{Start: 26, End: 32, Format: gmparse.Format{"i": "true"}},
				{Start: 32, End: 32, Format: gmparse.Format{"bq": "true"}},
			},
		},
		{
			name:         "List with Formatted Items",
			input:        "- one **bold**\n- two ~~strike~~\n- three *italic*\n",
			expectedText: "one bold\ntwo strike\nthree italic\n",
			expectedActions: []gmparse.FormatSpan{
				{Start: 4, End: 8, Format: gmparse.Format{"b": "true"}},
				{Start: 8, End: 8, Format: gmparse.Format{"ul": "0"}},
				{Start: 13, End: 19, Format: gmparse.Format{"s": "true"}},
				{Start: 19, End: 19, Format: gmparse.Format{"ul": "0"}},
				{Start: 26, End: 32, Format: gmparse.Format{"i": "true"}},
				{Start: 32, End: 32, Format: gmparse.Format{"ul": "0"}},
			},
		},
		{
			name:         "Single Blockquote with Bold",
			input:        "> one **bold**\n",
			expectedText: "one bold\n",
			expectedActions: []gmparse.FormatSpan{
				{Start: 4, End: 8, Format: gmparse.Format{"b": "true"}},
				{Start: 8, End: 8, Format: gmparse.Format{"bq": "true"}},
			},
		},
		{
			name:         "Nested Ordered List 2",
			input:        "1. one\n   1. two\n      1. three\n",
			expectedText: "one\ntwo\nthree\n",
			expectedActions: []gmparse.FormatSpan{
				{Start: 3, End: 3, Format: gmparse.Format{"ol": "0"}},
				{Start: 7, End: 7, Format: gmparse.Format{"ol": "1"}},
				{Start: 13, End: 13, Format: gmparse.Format{"ol": "2"}},
			},
		},
		{
			name:            "Simple Multiline",
			input:           "one\ntwo\nthree\n",
			expectedText:    "one\ntwo\nthree\n",
			expectedActions: []gmparse.FormatSpan{},
		},
		{
			name:            "Double Newlines",
			input:           "one\n\ntwo\n\n\n\nthree\n",
			expectedText:    "one\ntwo\nthree\n",
			expectedActions: []gmparse.FormatSpan{},
		},
		{
			name:         "Indented list text",
			input:        "- one\n  two\n  three\n- four\n",
			expectedText: "one\ntwo\nthree\nfour\n",
			expectedActions: []gmparse.FormatSpan{
				{Start: 3, End: 3, Format: gmparse.Format{"ul": "0"}},
				{Start: 7, End: 7, Format: gmparse.Format{"il": "0"}},
				{Start: 13, End: 13, Format: gmparse.Format{"il": "0"}},
				{Start: 18, End: 18, Format: gmparse.Format{"ul": "0"}},
			},
		},
		{
			name:         "Horizontal Rule V2",
			input:        "## Goodbye\n\n---\n\nLoyal, wagging friend\nHearts sing with their love\n\n",
			expectedText: "Goodbye\n\nLoyal, wagging friend\nHearts sing with their love\n",
			expectedActions: []gmparse.FormatSpan{
				{Start: 7, End: 7, Format: gmparse.Format{"h": "2"}},
				{Start: 8, End: 8, Format: gmparse.Format{"r": "true"}},
			},
		},
		{
			name:         "image0 no caption",
			input:        "<figure><img src=\"https://www.fake.com/image.jpg\"><figcaption></figcaption></figure>\n\n",
			expectedText: "\n",
			expectedActions: []gmparse.FormatSpan{
				{
					Start: 0, End: 0, Format: gmparse.Format{
						"img":    "https://www.fake.com/image.jpg",
						"alt":    "",
						"height": "",
						"width":  "",
					},
				},
			},
		},
		{
			name:         "image1 with caption",
			input:        "<figure><img src=\"https://www.fake.com/image.jpg\"><figcaption>Hello World!</figcaption></figure>\n\n",
			expectedText: "Hello World!\n",
			expectedActions: []gmparse.FormatSpan{
				{
					Start: 12, End: 12, Format: gmparse.Format{
						"img":    "https://www.fake.com/image.jpg",
						"alt":    "",
						"height": "",
						"width":  "",
					},
				},
			},
		},
		{
			name:         "image2 with size, alt and caption",
			input:        "<figure><img src=\"https://www.fake.com/image.jpg\" style=\"width: 100px; height: 100px;\" alt=\"test\"><figcaption>Hello World!</figcaption></figure>\n\n",
			expectedText: "Hello World!\n",
			expectedActions: []gmparse.FormatSpan{
				{
					Start: 12, End: 12, Format: gmparse.Format{
						"img":    "https://www.fake.com/image.jpg",
						"alt":    "test",
						"height": "100px",
						"width":  "100px",
					},
				},
			},
		},
		{
			name:         "image3 with size, alt, caption and span formats",
			input:        "<figure><img src=\"https://www.fake.com/image.jpg\" style=\"width: 100px; height: 100px;\" alt=\"test\"><figcaption>**Hello** ~World!~</figcaption></figure>\n\n",
			expectedText: "Hello World!\n",
			expectedActions: []gmparse.FormatSpan{
				{
					Start: 0, End: 5, Format: gmparse.Format{"b": "true"},
				},
				{
					Start: 6, End: 12, Format: gmparse.Format{"s": "true"},
				},
				{
					Start: 12, End: 12, Format: gmparse.Format{
						"img":    "https://www.fake.com/image.jpg",
						"alt":    "test",
						"height": "100px",
						"width":  "100px",
					},
				},
			},
		},
		{
			name:         "image4 with other markdown",
			input:        "Hello World!\n\nGoodbye World!\n\n<figure><img src=\"https://example.com/image.png\" ><figcaption>NEATO!</figcaption></figure>\n\n",
			expectedText: "Hello World!\nGoodbye World!\nNEATO!\n",
			expectedActions: []gmparse.FormatSpan{
				{
					Start: 34, End: 34, Format: gmparse.Format{
						"img":    "https://example.com/image.png",
						"alt":    "",
						"height": "",
						"width":  "",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Printf("Test: %s\n", tt.name)
			fmt.Printf("Input: %q\n", tt.input)

			plaintext, actions, err := gmparse.SplitMarkdown(tt.input)
			fmt.Printf("Plaintext: %q\n", plaintext)

			require.NoError(t, err)
			require.Equal(t, tt.expectedText, plaintext)
			if tt.expectedActions != nil {
				require.Equal(t, tt.expectedActions, actions)
			}

			for _, a := range actions {
				if a.Start == a.End {
					fmt.Printf("%v: %q\n", a.Format, plaintext[a.Start:a.Start+1])
				} else {
					fmt.Printf("%v: %q\n", a.Format, plaintext[a.Start:a.End])
				}
			}
			fmt.Println()
		})
	}
}
