package gmparse

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"

	xhtml "golang.org/x/net/html"
)

type Format map[string]interface{}
type FormatSpan struct {
	Start  int
	End    int
	Format Format
}

type parserState struct {
	plaintext strings.Builder
	spans     []FormatSpan
	source    []byte
}

func (f FormatSpan) String() string {
	return fmt.Sprintf("Start: %d, End: %d, Format: %v", f.Start, f.End, f.Format)
}

func SplitMarkdown(markdown string) (string, []FormatSpan, error) {
	if markdown == "" {
		return "", nil, nil
	}

	source := []byte(markdown)
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.Strikethrough,
			extension.TaskList,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithASTTransformers(
				util.Prioritized(&lineTransformer{}, 100),
			),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(), // necessary for parsing the HTML image blocks and <u> tags eventually
		),
	)

	reader := text.NewReader(source)
	parser := md.Parser()
	doc := parser.Parse(reader)

	state := &parserState{
		source: source,
		spans:  []FormatSpan{},
	}

	err := parseNode(doc, state)
	if err != nil {
		return "", nil, err
	}

	// If there's no line style on the last newline and there's no
	// newline in the original markdown, we should trim the last newline
	plainText := state.plaintext.String()
	endsInNewline := strings.HasSuffix(markdown, "\n")

	if !endsInNewline {
		lsp := len(state.spans)
		if lsp > 0 {
			lastSpan := state.spans[lsp-1]
			isLineFormat := lastSpan.Start == lastSpan.End
			if lastSpan.End < state.plaintext.Len()-1 || !isLineFormat {
				plainText = strings.TrimRight(plainText, "\n")
			}
		} else {
			plainText = strings.TrimRight(plainText, "\n")
		}
	}

	return plainText, state.spans, nil
}

func getHtml(node *ast.HTMLBlock, state *parserState) string {
	builder := strings.Builder{}
	lines := node.Lines()
	for i := 0; i < lines.Len(); i++ {
		line := lines.At(i)
		builder.Write(line.Value(state.source))
	}

	return builder.String()
}

func parseNode(n ast.Node, state *parserState) error {
	if n.Kind().String() == "Strikethrough" {
		start := state.plaintext.Len()
		parseChildren(n, state)
		end := state.plaintext.Len()
		state.spans = append(state.spans, FormatSpan{Start: start, End: end, Format: Format{"s": "true"}})

		return nil
	}

	if n.Kind().String() == "TaskList" {
		// TODO: add support for checklists

		return nil
	}

	switch node := n.(type) {
	case *ast.HTMLBlock:
		img, err := splitHtml(getHtml(node, state))
		if err != nil {
			return err
		}

		nestedPlaintext, spans, err := SplitMarkdown(img.Caption)
		if err != nil {
			return err
		}

		offset := state.plaintext.Len()
		for _, span := range spans {
			state.spans = append(state.spans, FormatSpan{
				Start:  span.Start + offset,
				End:    span.End + offset,
				Format: span.Format,
			})
		}

		unescapedPlaintext := xhtml.UnescapeString(nestedPlaintext)
		state.plaintext.WriteString(unescapedPlaintext)
		state.spans = append(state.spans, FormatSpan{
			Start: state.plaintext.Len(),
			End:   state.plaintext.Len(),
			Format: Format{
				"img":    img.Src,
				"alt":    img.Alt,
				"width":  img.Width,
				"height": img.Height,
			},
		})
		state.plaintext.WriteString("\n")

	case *ast.Text:
		content := string(node.Text(state.source))
		unescapedContent := xhtml.UnescapeString(content)
		state.plaintext.WriteString(unescapedContent)

	case *ast.Emphasis:
		start := state.plaintext.Len()
		parseChildren(node, state)
		end := state.plaintext.Len()

		format := Format{"i": "true"}
		if node.Level == 2 {
			format = Format{"b": "true"}
		}
		state.spans = append(state.spans, FormatSpan{Start: start, End: end, Format: format})

	case *ast.Link:
		start := state.plaintext.Len()
		parseChildren(node, state)
		end := state.plaintext.Len()
		state.spans = append(state.spans, FormatSpan{Start: start, End: end, Format: Format{"a": string(node.Destination)}})

	case *ast.Heading:
		parseChildren(node, state)
		state.plaintext.WriteString("\n")
		end := state.plaintext.Len() - 1
		state.spans = append(state.spans, FormatSpan{Start: end, End: end, Format: Format{"h": strconv.Itoa(node.Level)}})

	case *Line:
		parseChildren(node, state)
		state.plaintext.WriteString("\n")
		end := state.plaintext.Len() - 1

		if isBlockquote(node) {
			state.spans = append(state.spans, FormatSpan{Start: end, End: end, Format: Format{"bq": "true"}})
		}

		isBullet, isOrdered, isIndented := isList(node)
		if isBullet || isOrdered || isIndented {
			format := "ul"
			indent := calculateIndentDepth(node)
			if isOrdered {
				format = "ol"
			} else if isIndented {
				format = "il"
			}
			state.spans = append(state.spans, FormatSpan{Start: end, End: end, Format: Format{format: strconv.Itoa(indent)}})
		}

	case *ast.ThematicBreak:
		state.plaintext.WriteString("\n")
		end := state.plaintext.Len() - 1
		state.spans = append(state.spans, FormatSpan{Start: end, End: end, Format: Format{"r": "true"}})

	case *ast.FencedCodeBlock:
		language := string(node.Language(state.source))

		lines := node.Lines()
		for i := 0; i < lines.Len(); i++ {
			line := lines.At(i)
			state.plaintext.Write(line.Value(state.source))
			end := state.plaintext.Len() - 1
			state.spans = append(state.spans, FormatSpan{Start: end, End: end, Format: Format{"cb": language}})
		}

	default:
		parseChildren(node, state)
	}

	return nil
}

func parseChildren(n ast.Node, state *parserState) {
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		parseNode(c, state)
	}
}

func calculateIndentDepth(node ast.Node) int {
	depth := 0
	for p := node.Parent(); p != nil; p = p.Parent() {
		if _, ok := p.(*ast.List); ok {
			depth++
		}
	}
	return depth - 1 // Subtract 1 because the outermost list doesn't count as an indent
}

func isBlockquote(node ast.Node) bool {
	for node != nil {
		if _, ok := node.(*ast.Blockquote); ok {
			return true
		}
		node = node.Parent()
	}
	return false
}

func _leftSibIsLine(node ast.Node) bool {
	if node.PreviousSibling() != nil {
		_, ok := node.PreviousSibling().(*Line)
		return ok
	}

	return false
}

func isList(line *Line) (bullet, ordered, indented bool) {
	var node ast.Node = line
	isSoftReturn := _leftSibIsLine(node)
	for node != nil {
		if list, ok := node.(*ast.List); ok {
			// TODO: add support for checklists
			if isSoftReturn {
				return false, false, true
			} else if list.IsOrdered() {
				return false, true, false
			} else {
				return true, false, false
			}
		}
		node = node.Parent()

	}

	return false, false, false
}

func StripWhitespace(s string) (prefix, stripped, suffix string) {
	if len(s) == 0 {
		return "", "", ""
	}

	runes := []rune(s)

	prefixEnd := len(runes)
	for i, r := range runes {
		if !unicode.IsSpace(r) {
			prefixEnd = i
			break
		}
	}

	if prefixEnd == len(runes) {
		return s, "", ""
	}

	suffixStart := len(runes)
	for i := len(runes) - 1; i >= prefixEnd; i-- {
		if !unicode.IsSpace(runes[i]) {
			suffixStart = i + 1
			break
		}
	}

	return string(runes[:prefixEnd]), string(runes[prefixEnd:suffixStart]), string(runes[suffixStart:])
}

func AlignWhitespace(oldPlaintext, newPlaintext string, formats []FormatSpan) (updatedPlaintext string, updatedFormats []FormatSpan) {
	oldPrefix, _, oldSuffix := StripWhitespace(oldPlaintext)
	newPrefix, newStripped, newSuffix := StripWhitespace(newPlaintext)

	diffLen := len(oldPrefix) - len(newPrefix)

	suffix := oldSuffix
	if len(newSuffix) > len(oldSuffix) {
		suffix = newSuffix
	}

	out := fmt.Sprintf("%s%s%s", oldPrefix, newStripped, suffix)
	for _, span := range formats {
		updatedFormats = append(updatedFormats, FormatSpan{
			Start:  span.Start + diffLen,
			End:    span.End + diffLen,
			Format: span.Format,
		})
	}

	return out, updatedFormats
}
