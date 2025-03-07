package v3_test

import (
	"fmt"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/stretchr/testify/require"
	v3 "github.com/fivetentaylor/pointy/rogue/v3"
	"github.com/fivetentaylor/pointy/rogue/v3/testcases"
)

func TestApplyMarkdownDiff(t *testing.T) {
	type testCase struct {
		name       string
		initialDoc func(r *v3.Rogue) error
		start      v3.ID
		end        v3.ID
		update     string

		expectedString  string
		expectedActions v3.Actions
		expectedOp      v3.MultiOp
	}

	testCases := []testCase{
		{
			name: "Hello world!",
			initialDoc: func(r *v3.Rogue) error {
				_, err := r.Insert(0, "Hello World!")
				return err
			},
			start:          v3.ID{Author: "root", Seq: 0},
			end:            v3.ID{Author: "q", Seq: 1},
			update:         "Hello Universe!\n\n",
			expectedString: "Hello Universe!\n\n",
			expectedActions: v3.Actions{
				v3.DeleteAction{Index: 6, Count: 6},
				v3.InsertAction{Index: 6, Text: "Universe!"},
			},
			expectedOp: v3.MultiOp{
				Mops: []v3.Op{
					v3.DeleteOp{ID: v3.ID{Author: "auth0", Seq: 15}, TargetID: v3.ID{Author: "auth0", Seq: 9}, SpanLength: 5},
					v3.InsertOp{ID: v3.ID{Author: "auth0", Seq: 16}, Text: "Universe", ParentID: v3.ID{Author: "auth0", Seq: 9}, Side: -1},
				},
			},
		},
		{
			name: "empty insert",
			initialDoc: func(r *v3.Rogue) error {
				_, err := r.Insert(0, "Hello  World!")
				return err
			},
			start:          v3.ID{Author: "auth0", Seq: 8},
			end:            v3.ID{Author: "auth0", Seq: 8},
			update:         "Cruel",
			expectedString: "Hello Cruel World!\n\n",
		},
		{
			name: "whitespace insert",
			initialDoc: func(r *v3.Rogue) error {
				_, err := r.Insert(0, "Hello  World!")
				return err
			},
			start:          v3.ID{Author: "auth0", Seq: 8},
			end:            v3.ID{Author: "auth0", Seq: 9},
			update:         "Cruel",
			expectedString: "Hello Cruel World!\n\n",
		},
		{
			name: "Hello world partial update",
			initialDoc: func(r *v3.Rogue) error {
				_, err := r.Insert(0, "Hello World!")
				return err
			},
			start:          v3.ID{Author: "auth0", Seq: 8},
			end:            v3.ID{Author: "q", Seq: 1},
			update:         "amazing Universe!\n\n",
			expectedString: "Hello amazing Universe!\n\n",
		},
		{
			name: "Hello world with formatting",
			initialDoc: func(r *v3.Rogue) error {
				_, err := r.Insert(0, "Hello World!")
				if err != nil {
					return err
				}
				_, err = r.Format(0, 5, v3.FormatV3Span{"b": "true"})
				return err
			},
			start:          v3.ID{Author: "root", Seq: 0},
			end:            v3.ID{Author: "auth0", Seq: 14},
			update:         "Hello Universe",
			expectedString: "Hello Universe!\n\n",
			expectedActions: v3.Actions{
				v3.FormatAction{Index: 0, Length: 5, Format: v3.FormatV3Span{"b": ""}},
				v3.DeleteAction{Index: 6, Count: 5},
				v3.InsertAction{Index: 6, Text: "Universe"},
			},
			expectedOp: v3.MultiOp{
				[]v3.Op{
					v3.FormatOp{
						ID:      v3.ID{Author: "auth0", Seq: 16},
						StartID: v3.ID{Author: "auth0", Seq: 3},
						EndID:   v3.ID{Author: "auth0", Seq: 8},
						Format:  v3.FormatV3Span{"e": "true"},
					},
					v3.DeleteOp{ID: v3.ID{Author: "auth0", Seq: 17}, TargetID: v3.ID{Author: "auth0", Seq: 9}, SpanLength: 5},
					v3.InsertOp{ID: v3.ID{Author: "auth0", Seq: 18}, Text: "Universe", ParentID: v3.ID{Author: "auth0", Seq: 9}, Side: -1},
				},
			},
		},
		{
			name: "line formatting",
			initialDoc: func(r *v3.Rogue) error {
				_, err := r.Insert(0, "Hello World!")
				require.NoError(t, err)

				_, err = r.Format(12, 1, v3.FormatV3Header(2))
				require.NoError(t, err)

				return nil
			},
			start:          v3.ID{Author: "root", Seq: 0},
			end:            v3.ID{Author: "q", Seq: 1},
			update:         "# Hello **Univ**erse!\n\n",
			expectedString: "# Hello **Univ**erse!\n\n",
		},
		{
			name: "real markdown",
			initialDoc: func(r *v3.Rogue) error {
				_, err := r.Insert(0, "This is a header\nThis is a list\nsubitem1")
				require.NoError(t, err)

				_, err = r.Format(16, 1, v3.FormatV3Header(1))
				require.NoError(t, err)

				_, err = r.Format(31, 1, v3.FormatV3BulletList(0))
				require.NoError(t, err)

				_, err = r.Format(40, 1, v3.FormatV3BulletList(1))
				require.NoError(t, err)

				return nil
			},
			start:          v3.ID{Author: "root", Seq: 0},
			end:            v3.ID{Author: "q", Seq: 1},
			update:         "# Hello Universe!\n\n- this is a list\n  - subitem2\n",
			expectedString: "# Hello Universe!\n\n- this is a list\n  - subitem2\n",
		},
		{
			name: "ordered list",
			initialDoc: func(r *v3.Rogue) error {
				_, err := r.Insert(0, "This is a header\nThis is a list\nsubitem1")
				require.NoError(t, err)

				_, err = r.Format(16, 1, v3.FormatV3Header(1))
				require.NoError(t, err)

				_, err = r.Format(31, 1, v3.FormatV3BulletList(0))
				require.NoError(t, err)

				_, err = r.Format(40, 1, v3.FormatV3BulletList(1))
				require.NoError(t, err)

				return nil
			},
			start:          v3.ID{Author: "root", Seq: 0},
			end:            v3.ID{Author: "q", Seq: 1},
			update:         "# Hello Universe!\n\n1. this is a list\n   1. subitem2\n",
			expectedString: "# Hello Universe!\n\n1. this is a list\n   1. subitem2\n",
		},
		{
			name: "deep nest",
			initialDoc: func(r *v3.Rogue) error {
				r.Insert(0, "one\ntwo\nthree\nfour")
				r.Format(3, 1, v3.FormatV3BulletList(0))
				r.Format(7, 1, v3.FormatV3BulletList(1))
				r.Format(13, 1, v3.FormatV3BulletList(2))
				r.Format(18, 1, v3.FormatV3BulletList(3))
				return nil
			},
			start:          v3.ID{Author: "root", Seq: 0},
			end:            v3.ID{Author: "q", Seq: 1},
			update:         "- one\n  - two\n    - three\n      - four\n- five\n",
			expectedString: "- one\n  - two\n    - three\n      - four\n- five\n",
		},
		{
			name: "span diffs",
			initialDoc: func(r *v3.Rogue) error {
				_, err := r.Insert(0, "Hello World!")
				require.NoError(t, err)

				_, err = r.Format(12, 1, v3.FormatV3BulletList(0))
				require.NoError(t, err)

				return nil
			},
			start:          v3.ID{Author: "root", Seq: 0},
			end:            v3.ID{Author: "q", Seq: 1},
			update:         "- **Hello** World!\n",
			expectedString: "- **Hello** World!\n",
		},
		{
			name: "middle format",
			initialDoc: func(r *v3.Rogue) error {
				_, err := r.Insert(0, "Hello World!\nitem1\nitem2")
				require.NoError(t, err)

				_, err = r.Format(12, 1, v3.FormatV3Header(1))
				require.NoError(t, err)

				_, err = r.Format(18, 1, v3.FormatV3BulletList(0))
				require.NoError(t, err)

				_, err = r.Format(24, 1, v3.FormatV3BulletList(0))
				require.NoError(t, err)

				return nil
			},
			start:          v3.ID{Author: "auth0", Seq: 15},
			end:            v3.ID{Author: "auth0", Seq: 21},
			update:         "- **item1**\n",
			expectedString: "# Hello World!\n\n- **item1**\n- item2\n",
		},
		{
			name: "middle content",
			initialDoc: func(r *v3.Rogue) error {
				_, err := r.Insert(0, "Hello World!\nitem1\nitem2")
				require.NoError(t, err)

				_, err = r.Format(12, 1, v3.FormatV3Header(1))
				require.NoError(t, err)

				_, err = r.Format(18, 1, v3.FormatV3BulletList(0))
				require.NoError(t, err)

				_, err = r.Format(24, 1, v3.FormatV3BulletList(0))
				require.NoError(t, err)

				return nil
			},
			start:          v3.ID{Author: "auth0", Seq: 15},
			end:            v3.ID{Author: "auth0", Seq: 21},
			update:         "- item7\n",
			expectedString: "# Hello World!\n\n- item7\n- item2\n",
		},
		{
			name: "middle complex content",
			initialDoc: func(r *v3.Rogue) error {
				_, err := r.Insert(0, "Hello World!\nthat is a great list\nyes it's the best")
				require.NoError(t, err)

				_, err = r.Format(12, 1, v3.FormatV3Header(1))
				require.NoError(t, err)

				_, err = r.Format(33, 1, v3.FormatV3BulletList(0))
				require.NoError(t, err)

				_, err = r.Format(51, 1, v3.FormatV3BulletList(1))
				require.NoError(t, err)

				return nil
			},
			start:          v3.ID{Author: "auth0", Seq: 15},
			end:            v3.ID{Author: "auth0", Seq: 36},
			update:         "- that's a **great** list\n",
			expectedString: "# Hello World!\n\n- that&#39;s a **great** list\n  - yes it&#39;s the best\n",
		},
		{
			name: "tail insert",
			initialDoc: func(r *v3.Rogue) error {
				_, err := r.Insert(0, "Hello World!\none\ntwo\nthree")
				require.NoError(t, err)

				_, err = r.Format(12, 1, v3.FormatV3Header(1))
				require.NoError(t, err)

				_, err = r.Format(16, 1, v3.FormatV3BulletList(0))
				require.NoError(t, err)

				_, err = r.Format(20, 1, v3.FormatV3BulletList(1))
				require.NoError(t, err)

				_, err = r.Format(26, 1, v3.FormatV3BulletList(2))
				require.NoError(t, err)

				return nil
			},
			start:          v3.ID{Author: "root", Seq: 0},
			end:            v3.ID{Author: "q", Seq: 1},
			update:         "# Hello World!\n\n- one\n  - two\n    - three\n- **four**\n",
			expectedString: "# Hello World!\n\n- one\n  - two\n    - three\n- **four**\n",
		},
		{
			name: "emojis",
			initialDoc: func(r *v3.Rogue) error {
				_, err := r.Insert(0, "Hello Wo💀rld!\nthat is a great list\nyes it's the best")
				require.NoError(t, err)

				_, err = r.Format(14, 1, v3.FormatV3Header(1))
				require.NoError(t, err)

				_, err = r.Format(35, 1, v3.FormatV3BulletList(0))
				require.NoError(t, err)

				_, err = r.Format(53, 1, v3.FormatV3BulletList(1))
				require.NoError(t, err)

				return nil
			},
			start:          v3.ID{Author: "auth0", Seq: 17},
			end:            v3.ID{Author: "auth0", Seq: 38},
			update:         "- that's a **gr👍eat** list\n",
			expectedString: "# Hello Wo💀rld!\n\n- that&#39;s a **gr👍eat** list\n  - yes it&#39;s the best\n",
		},
		{
			name: "emoji list",
			initialDoc: func(r *v3.Rogue) error {
				_, err := r.Insert(0, "Hello World!\none\ntwo\nthree")
				require.NoError(t, err)

				_, err = r.Format(12, 1, v3.FormatV3Header(1))
				require.NoError(t, err)

				_, err = r.Format(16, 1, v3.FormatV3BulletList(0))
				require.NoError(t, err)

				_, err = r.Format(20, 1, v3.FormatV3BulletList(1))
				require.NoError(t, err)

				_, err = r.Format(26, 1, v3.FormatV3BulletList(2))
				require.NoError(t, err)

				return nil
			},
			start:          v3.ID{Author: "root", Seq: 0},
			end:            v3.ID{Author: "q", Seq: 1},
			update:         "# Hello World!\n\n- one\n- two\n- three\n## List of Emojis\n\n- 😀 Grinning Face\n- 😂 Face with Tears of Joy\n- ❤️ Red Heart\n- 🚀 Rocket\n- 🌟 Glowing Star",
			expectedString: "# Hello World!\n\n- one\n- two\n- three\n## List of Emojis\n\n- 😀 Grinning Face\n- 😂 Face with Tears of Joy\n- ❤️ Red Heart\n- 🚀 Rocket\n- 🌟 Glowing Star\n",
		},
		{
			name: "replace list test",
			initialDoc: func(r *v3.Rogue) error {
				_, err := r.Insert(0, "Hello World!\none\ntwo\nthree")
				require.NoError(t, err)

				_, err = r.Format(12, 1, v3.FormatV3Header(1))
				require.NoError(t, err)

				_, err = r.Format(16, 1, v3.FormatV3BulletList(0))
				require.NoError(t, err)

				_, err = r.Format(20, 1, v3.FormatV3BulletList(1))
				require.NoError(t, err)

				_, err = r.Format(26, 1, v3.FormatV3BulletList(2))
				require.NoError(t, err)

				return nil
			},
			start:          v3.ID{Author: "root", Seq: 0},
			end:            v3.ID{Author: "q", Seq: 1},
			update:         "# Hello World!\n\nThis is **great!**\n\n",
			expectedString: "# Hello World!\n\nThis is **great!**\n\n",
		},
		{
			name: "Bad delete range real doc",
			initialDoc: func(r *v3.Rogue) error {
				*r = *testcases.Load(t, "bad_apply.json")
				return nil
			},
			start:          v3.ID{Author: "1", Seq: 20},
			end:            v3.ID{Author: "0000", Seq: 392},
			update:         "This is my ugly and whimsical document that stinks...\n\nAnd chickens are brown",
			expectedString: "# Hello World!\n\nThis is my ugly and whimsical document that stinks...\n\nAnd chickens are brown\n\n",
		},
		{
			name: "html escaping",
			initialDoc: func(r *v3.Rogue) error {
				r.Insert(0, "Hello <'html>'&\"World\"!")
				return nil
			},
			start:          v3.ID{Author: "root", Seq: 0},
			end:            v3.ID{Author: "q", Seq: 1},
			update:         "# Hello World!\n\n",
			expectedString: "# Hello World!\n\n",
		},
		{
			name: "real doc 2",
			initialDoc: func(r *v3.Rogue) error {
				*r = *testcases.Load(t, "invalid_format_ix.json")
				return nil
			},
			start:          v3.ID{Author: "root", Seq: 0},
			end:            v3.ID{Author: "1", Seq: 113},
			update:         "# Welcome to Our Project! 🌎\n\nHello World! This is the beginning of something great. Below, you'll find a quick guide on what to expect and how to get started.\n\n## Quick Start Guide\n\n1. **Introduction**\n   - Welcome to our project! We're glad you're here.\n1. **Setup**\n   - Ensure you have the necessary tools installed.\n   - For example, to check your Python version, use the following command in your terminal:\n```bash\n     python --version\n```\n3. **Contribution Guidelines**\n   - Read through our contribution guidelines to understand how you can contribute effectively.\n4. **Community**\n   - Join our community on [GitHub](https://github.com) or [Discord](https://discord.com) to stay updated and collaborate.\n\nWe're excited to have you on board and can't wait to see what we'll achieve together!",
			expectedString: "# Welcome to Our Project! 🌎\n\nHello World! This is the beginning of something great. Below, you&#39;ll find a quick guide on what to expect and how to get started.\n\n## Quick Start Guide\n\n1. **Introduction**\n  - Welcome to our project! We&#39;re glad you&#39;re here.\n1. **Setup**\n  - Ensure you have the necessary tools installed.\n  - For example, to check your Python version, use the following command in your terminal:\n```bash\n     python --version\n```\n1. **Contribution Guidelines**\n  - Read through our contribution guidelines to understand how you can contribute effectively.\n1. **Community**\n  - Join our community on [GitHub](https://github.com) or [Discord](https://discord.com) to stay updated and collaborate.\nWe&#39;re excited to have you on board and can&#39;t wait to see what we&#39;ll achieve together!\n\n",
		},
		{
			name: "real doc 3",
			initialDoc: func(r *v3.Rogue) error {
				*r = *testcases.Load(t, "bad_apply_2.json")
				return nil
			},
			start:          v3.ID{Author: "root", Seq: 0},
			end:            v3.ID{Author: "1", Seq: 113},
			update:         "Hello World!🌎",
			expectedString: "Hello [World](https://google.com)!🌎\n\n",
		},
		{
			name: "real doc 4",
			initialDoc: func(r *v3.Rogue) error {
				*r = *testcases.Load(t, "bad_span_format.json")
				return nil
			},
			start:          v3.ID{Author: "root", Seq: 0},
			end:            v3.ID{Author: "1", Seq: 113},
			update:         "Hello World!🌎",
			expectedString: "Hello [World](https://google.com)!🌎\n\n",
		},
		/*{
		      // The markdown parser checks indent based up nesting in other lists. This update being a single list item
		      // fails to register as a nested list item, so it get's inserted with indent == 0 and then causess the
		      // surrounding list to be converted from ordered to bullet
					name: "real doc 5",
					initialDoc: func(r *v3.Rogue) error {
						*r = *testcases.Load(t, "bad_apply_3.json")
						return nil
					},
					start:          v3.ID{Author: "0000", Seq: 5124},
					end:            v3.ID{Author: "0000", Seq: 5125},
					update:         "   - Embrace the philosophy of kung fu in your contributions: approach each task with dedication, patience, and a mindset geared towards continuous learning and mastery.",
					expectedString: "# Welcome to Our Project! 🌎\n\nHello World! This is the beginning of something great. Below, you'll find a quick guide on what to expect and how to get started.\n\n## Quick Start Guide\n\n1. **Introduction**\n  - Welcome to our project! We're glad you're here.\n1. **Setup**\n  - Ensure you have the necessary tools installed.\n  - For example, to check your Python version, use the following command in your terminal:\n```bash\n     python --version\n```\n1. **Contribution Guidelines**\n  - Embrace the philosophy of kung fu in your contributions: approach each task with dedication, patience, and a mindset geared towards continuous learning and mastery.\n  - Read through our contribution guidelines to understand how you can contribute effectively.\n1. **Community**\n  - Join our community on [GitHub](https://github.com) or [Discord](https://discord.com) to stay updated and collaborate.\nWe're excited to have you on board and can't wait to see what we'll achieve together!\n\n",
				},*/
		{
			name: "preexisting markdown",
			initialDoc: func(r *v3.Rogue) error {
				r.Insert(0, "# Hello <'html>'&\"World\"!\n\n- some list\n- markdown")
				return nil
			},
			start:          v3.ID{Author: "root", Seq: 0},
			end:            v3.ID{Author: "q", Seq: 1},
			update:         "# Hello World!\n\n- some list\n- markdown\n",
			expectedString: "# Hello World!\n\n- some list\n- markdown\n",
		},
		{
			name: "code",
			initialDoc: func(r *v3.Rogue) error {
				code := "def hello(name):\n\tprint(f\"Hello {name}\")\n\ndef goodbye(name):\n\tprint(f\"Goodbye {name}\")\n\nif __name__ == \"__main__\":\n\thello(\"Taylor\")\n\tgoodbye(\"Taylor\")"
				_, err := r.Insert(0, code)
				require.NoError(t, err)

				_, err = r.Format(16, 1, v3.FormatV3CodeBlock("python"))
				require.NoError(t, err)

				_, err = r.Format(40, 1, v3.FormatV3CodeBlock("python"))
				require.NoError(t, err)

				_, err = r.Format(41, 1, v3.FormatV3CodeBlock("python"))
				require.NoError(t, err)

				_, err = r.Format(60, 1, v3.FormatV3CodeBlock("python"))
				require.NoError(t, err)

				_, err = r.Format(86, 1, v3.FormatV3CodeBlock("python"))
				require.NoError(t, err)

				_, err = r.Format(87, 1, v3.FormatV3CodeBlock("python"))
				require.NoError(t, err)

				_, err = r.Format(114, 1, v3.FormatV3CodeBlock("python"))
				require.NoError(t, err)

				_, err = r.Format(131, 1, v3.FormatV3CodeBlock("python"))
				require.NoError(t, err)

				_, err = r.Format(150, 1, v3.FormatV3CodeBlock("python"))
				require.NoError(t, err)

				return nil
			},
			start:          v3.ID{Author: "root", Seq: 0},
			end:            v3.ID{Author: "q", Seq: 1},
			update:         "```python\n# Define a function named `hello` that greets the given name\ndef hello(name):\n\t# Print a greeting message including the name\n\tprint(f\"Hello {name}\")\n\n# Define a function named `goodbye` that says farewell to the given name\ndef goodbye(name):\n\t# Print a farewell message including the name\n\tprint(f\"Goodbye {name}\")\n\n# Check if this script is the main program and not being imported by another module\nif __name__ == \"__main__\":\n\t# Call the `hello` function with \"Taylor\" as the argument\n\thello(\"Taylor\")\n\t# Also call the `goodbye` function with \"Taylor\" as the argument\n\tgoodbye(\"Taylor\")\n```",
			expectedString: "```python\n# Define a function named `hello` that greets the given name\ndef hello(name):\n\t# Print a greeting message including the name\n\tprint(f&#34;Hello {name}&#34;)\n\n# Define a function named `goodbye` that says farewell to the given name\ndef goodbye(name):\n\t# Print a farewell message including the name\n\tprint(f&#34;Goodbye {name}&#34;)\n\n# Check if this script is the main program and not being imported by another module\nif __name__ == &#34;__main__&#34;:\n\t# Call the `hello` function with &#34;Taylor&#34; as the argument\n\thello(&#34;Taylor&#34;)\n\t# Also call the `goodbye` function with &#34;Taylor&#34; as the argument\n\tgoodbye(&#34;Taylor&#34;)\n```\n",
		},
		{
			name: "code reverse",
			initialDoc: func(r *v3.Rogue) error {
				*r = *testcases.Load(t, "code_block.json")

				// update := "```\ndef hello_iter(name, count):\n  for _ in range(count):\n```\\n    yield f\\\"Hello {name}\\\"\\n\\n\\n\\n```\\nfor name in hello_iter(\\\"Taylor\\\", 5):\\n  print(name)\\n```\\n"
				return nil
			},
			start:          v3.ID{Author: "root", Seq: 0},
			end:            v3.ID{Author: "q", Seq: 1},
			update:         "```python\n# Define a function named `hello` that greets the given name\ndef hello(name):\n\t# Print a greeting message including the name\n\tprint(f\"Hello {name}\")\n\n# Define a function named `goodbye` that says farewell to the given name\ndef goodbye(name):\n\t# Print a farewell message including the name\n\tprint(f\"Goodbye {name}\")\n\n# Check if this script is the main program and not being imported by another module\nif __name__ == \"__main__\":\n\t# Call the `hello` function with \"Taylor\" as the argument\n\thello(\"Taylor\")\n\t# Also call the `goodbye` function with \"Taylor\" as the argument\n\tgoodbye(\"Taylor\")\n```",
			expectedString: "```python\n# Define a function named `hello` that greets the given name\ndef hello(name):\n\t# Print a greeting message including the name\n\tprint(f&#34;Hello {name}&#34;)\n\n# Define a function named `goodbye` that says farewell to the given name\ndef goodbye(name):\n\t# Print a farewell message including the name\n\tprint(f&#34;Goodbye {name}&#34;)\n\n# Check if this script is the main program and not being imported by another module\nif __name__ == &#34;__main__&#34;:\n\t# Call the `hello` function with &#34;Taylor&#34; as the argument\n\thello(&#34;Taylor&#34;)\n\t# Also call the `goodbye` function with &#34;Taylor&#34; as the argument\n\tgoodbye(&#34;Taylor&#34;)\n```\n",
		},
		{
			name: "deleted span",
			initialDoc: func(r *v3.Rogue) error {
				_, err := r.Insert(0, "Hello Cruel World!")
				require.NoError(t, err)

				_, err = r.Format(0, 18, v3.FormatV3Span{"b": "true"})
				require.NoError(t, err)

				// Delete the span that the update is being applied to
				_, err = r.Delete(6, 5)
				require.NoError(t, err)

				return nil
			},
			start:          v3.ID{Author: "auth0", Seq: 9},
			end:            v3.ID{Author: "auth0", Seq: 13},
			update:         "~~Cruel~~*Beautiful*",
			expectedString: "**Hello ~~Cruel~~*Beautiful* World!**\n\n",
		},
		{
			name: "with images0",
			initialDoc: func(r *v3.Rogue) error {
				_, err := r.Insert(0, "Hello World!\nGoodbye World!\nNeato!")
				require.NoError(t, err)

				_, err = r.Format(34, 1, v3.FormatV3Image{
					Src: "https://example.com/image.png",
				})

				return nil
			},
			start:          v3.ID{Author: "root", Seq: 0},
			end:            v3.ID{Author: "q", Seq: 1},
			update:         "Hello World!\n\nGoodbye World!\n\n<figure><img src=\"https://example.com/image.png\" /><figcaption>NEATO!</figcaption></figure>\n\n",
			expectedString: "Hello World!\n\nGoodbye World!\n\n<figure><img src=\"https://example.com/image.png\" /><figcaption>NEATO!</figcaption></figure>\n\n",
		},
		{
			name: "with images1",
			initialDoc: func(r *v3.Rogue) error {
				_, err := r.Insert(0, "Hello World!")
				require.NoError(t, err)

				_, err = r.Format(12, 1, v3.FormatV3Image{
					Src: "https://example.com/image.png",
				})

				return nil
			},
			start:          v3.ID{Author: "auth0", Seq: 14},
			end:            v3.ID{Author: "q", Seq: 1},
			update:         "<figure><img src=\"https://example.com/image.png\" /><figcaption></figcaption></figure>\n\n# Cats\n\n- are cuddly\n- are cute\n",
			expectedString: "<figure><img src=\"https://example.com/image.png\" /><figcaption>Hello World!</figcaption></figure>\n\n# Cats\n\n- are cuddly\n- are cute\n",
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

			firstID, err := r.GetFirstTotID()
			require.NoError(t, err)

			lastID, err := r.GetLastTotID()
			require.NoError(t, err)

			mdBefore, err := r.GetMarkdownBeforeAfter(firstID, lastID)
			require.NoError(t, err)

			fmt.Printf("mdBefore: %q\n", mdBefore)

			md0, err := r.GetMarkdownBeforeAfter(tc.start, tc.end)
			require.NoError(t, err)

			fmt.Println("BEFORE")
			fmt.Printf("md: %q\n", md0)
			// fmt.Printf("r.GetText(): %q\n", r.GetText())

			r.Author = "000"
			mop, actions, err := r.ApplyMarkdownDiff("auth0", tc.update, tc.start, tc.end)
			if err != nil {
				t.Fatalf("ApplyMarkdownDiff(auth0, %q, %v, %v): %v", tc.update, tc.start, tc.end, err)
			}

			// v3.PrintJson("actions", actions)
			if false && tc.expectedActions != nil {
				require.Equal(t, tc.expectedActions, actions)
			}

			// v3.PrintJson("ops", ops)
			if len(tc.expectedOp.Mops) > 0 {
				require.Equal(t, tc.expectedOp, mop)
			}

			firstID, err = r.GetFirstTotID()
			require.NoError(t, err)

			lastID, err = r.GetLastTotID()
			require.NoError(t, err)

			md, err := r.GetMarkdownBeforeAfter(firstID, lastID)
			require.NoError(t, err)

			/*
				fmt.Println("AFTER")
				fmt.Printf("r.GetText(): %q\n", r.GetText())
				fmt.Printf("md: %q\n", md)
			*/

			require.Equal(t, tc.expectedString, md)
		})
	}
}

func TestDiffWords(t *testing.T) {
	tests := []struct {
		name     string
		text1    string
		text2    string
		expected []diffmatchpatch.Diff
	}{
		{
			name:  "same words",
			text1: "Hello World",
			text2: "Hello World",
			expected: []diffmatchpatch.Diff{
				{Type: diffmatchpatch.DiffEqual, Text: "Hello World"},
			},
		},
		{
			name:  "word change",
			text1: "Hello World",
			text2: "Goodbye World",
			expected: []diffmatchpatch.Diff{
				{Type: diffmatchpatch.DiffDelete, Text: "Hello"},
				{Type: diffmatchpatch.DiffInsert, Text: "Goodbye"},
				{Type: diffmatchpatch.DiffEqual, Text: " World"},
			},
		},
		{
			name:  "punctuation",
			text1: "Hello, World!",
			text2: "Hello. World?",
			expected: []diffmatchpatch.Diff{
				{Type: 0, Text: "Hello"},
				{Type: -1, Text: ","},
				{Type: 1, Text: "."},
				{Type: 0, Text: " World"},
				{Type: -1, Text: "!"},
				{Type: 1, Text: "?"},
			},
		},
		{
			name:  "multiline0",
			text1: "Hello\nWorld\nHow are you?",
			text2: "Hello\nWorld friends\nHow are you!",
			expected: []diffmatchpatch.Diff{
				{Type: 0, Text: "Hello"},
				{Type: 0, Text: "\n"},
				{Type: 0, Text: "World"},
				{Type: 1, Text: " friends"},
				{Type: 0, Text: "\n"},
				{Type: 0, Text: "How are you"},
				{Type: -1, Text: "?"},
				{Type: 1, Text: "!"},
			},
		},
		{
			name:  "multiline1",
			text1: "Hello\nWorld\nHow are you?",
			text2: "Hello\nEarth\nHow are you!",
			expected: []diffmatchpatch.Diff{
				{Type: 0, Text: "Hello"},
				{Type: 0, Text: "\n"},
				{Type: -1, Text: "World"},
				{Type: 1, Text: "Earth"},
				{Type: 0, Text: "\n"},
				{Type: 0, Text: "How are you"},
				{Type: -1, Text: "?"},
				{Type: 1, Text: "!"},
			},
		},
		{
			name:  "multiple spaces",
			text1: "Hello   World",
			text2: "Hello World",
			expected: []diffmatchpatch.Diff{
				{Type: 0, Text: "Hello "},
				{Type: -1, Text: "  "},
				{Type: 0, Text: "World"},
			},
		},
		{
			name:  "apostrophes and contractions",
			text1: "It's a beautiful day, isn't it?",
			text2: "It is a beautiful day, is it not?",
			expected: []diffmatchpatch.Diff{
				{Type: 0, Text: "It"},
				{Type: -1, Text: "'s"},
				{Type: 1, Text: " is"},
				{Type: 0, Text: " a beautiful day, "},
				{Type: 1, Text: "is it not?"},
				{Type: -1, Text: "isn't it?"},
			},
		},
		{
			name:  "mixed case and numbers",
			text1: "There are 10 APPLES and 5 oranges.",
			text2: "There are ten apples and five ORANGES.",
			expected: []diffmatchpatch.Diff{
				{Type: 0, Text: "There are "},
				{Type: -1, Text: "10 APPLES"},
				{Type: 1, Text: "ten apples"},
				{Type: 0, Text: " and "},
				{Type: -1, Text: "5 oranges"},
				{Type: 1, Text: "five ORANGES"},
				{Type: 0, Text: "."},
			},
		},
		{
			name:  "special characters",
			text1: "Email: user@example.com, Phone: 123-456-7890",
			text2: "Email: admin@example.com, Tel: (123) 456-7890",
			expected: []diffmatchpatch.Diff{
				{Type: 0, Text: "Email: "},
				{Type: -1, Text: "user@example.com, Phone"},
				{Type: 1, Text: "admin@example.com, Tel"},
				{Type: 0, Text: ": "},
				{Type: 1, Text: "("},
				{Type: 0, Text: "123"},
				{Type: -1, Text: "-"},
				{Type: 1, Text: ") "},
				{Type: 0, Text: "456-7890"},
			},
		},
		{
			name:  "tabs and newlines",
			text1: "Line 1\n\tIndented\nLine 3",
			text2: "Line 1\n    Indented\nLine 3",
			expected: []diffmatchpatch.Diff{
				{Type: 0, Text: "Line 1"},
				{Type: 0, Text: "\n"},
				{Type: -1, Text: "\t"},
				{Type: 1, Text: "    "},
				{Type: 0, Text: "Indented"},
				{Type: 0, Text: "\n"},
				{Type: 0, Text: "Line 3"},
			},
		},
		{
			name:  "adds a newline in middle of text",
			text1: "Hello WorldFriends",
			text2: "Hello World\nFriends",
			expected: []diffmatchpatch.Diff{
				{Type: diffmatchpatch.DiffEqual, Text: "Hello World"},
				{Type: diffmatchpatch.DiffInsert, Text: "\n"},
				{Type: diffmatchpatch.DiffEqual, Text: "Friends"},
			},
		},
		{
			name:  "equals with newline",
			text1: "\n",
			text2: "\n# Hello World\nFriends\n",
			expected: []diffmatchpatch.Diff{
				{Type: diffmatchpatch.DiffEqual, Text: "\n"},
				{Type: diffmatchpatch.DiffInsert, Text: "# Hello World\nFriends\n"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diffs := v3.DiffWords(tt.text1, tt.text2)
			require.Equal(t, tt.expected, diffs)
		})
	}
}

func TestDiffLines(t *testing.T) {
	tests := []struct {
		name     string
		text1    string
		text2    string
		expected []diffmatchpatch.Diff
	}{
		{
			name:  "same text",
			text1: "Hello\nWorld",
			text2: "Hello\nWorld",
			expected: []diffmatchpatch.Diff{
				{Type: diffmatchpatch.DiffEqual, Text: "Hello\nWorld"},
			},
		},
		{
			name:  "added line",
			text1: "Hello\n",
			text2: "Hello\nWorld",
			expected: []diffmatchpatch.Diff{
				{Type: diffmatchpatch.DiffEqual, Text: "Hello\n"},
				{Type: diffmatchpatch.DiffInsert, Text: "World"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diffs := v3.DiffLines(tt.text1, tt.text2)
			require.Equal(t, tt.expected, diffs)
		})
	}
}

func TestSplitIncludingWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected []string
	}{
		{
			name:     "empty string",
			text:     "",
			expected: nil,
		},
		{
			name:     "single line",
			text:     "Hello World",
			expected: []string{"Hello ", "World"},
		},
		{
			name:     "multiple lines",
			text:     "Hello\nWorld",
			expected: []string{"Hello", "\n", "World"},
		},
		{
			name:     "multiple lines with trailing newline",
			text:     "Hello\nWorld\n",
			expected: []string{"Hello", "\n", "World", "\n"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v3.SplitIncludingWhitespace(tt.text)
			require.Equal(t, tt.expected, result)
		})
	}
}

var mdWithImages = `# Feline Whispers

## Whiskers in Moonlight

Soft paws on moonbeams,

Green eyes gleam in darkness,

Purrs echo softly.

<figure><img src="https://app.reviso.dev:9090/drafts/76399538-0aa4-4a30-9d96-449681595b9d/images/5c762946-5401-4ebb-bd49-5091c0ff7c2c"><figcaption>My pussy</figcaption></figure>

## Cat Virtues

1. Independence - Cats are self-reliant and confident
2. Grace - Felines move with elegance and poise
3. Cleanliness - They are meticulous in their grooming habits
4. Curiosity - Always exploring and learning about their environment
5. Affection - Capable of deep bonds and showing love on their own terms
6. Patience - Masters of the art of waiting and observing
7. Agility - Quick reflexes and impressive acrobatic abilities
8. Intuition - Keen sense of their surroundings and people's emotions
9. Playfulness - Retaining a kitten-like joy throughout their lives
10. Adaptability - Able to thrive in various environments and situations`
