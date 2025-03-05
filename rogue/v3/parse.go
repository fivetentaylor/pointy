package v3

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/andybalholm/cascadia"
	"github.com/charmbracelet/log"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var newlinePattern = regexp.MustCompile(`\n`)
var spacePattern = regexp.MustCompile(`\s+`)

type TextSpan struct {
	StartIndex int
	EndIndex   int
	Format     FormatV3
}

func (s *TextSpan) MergeFormat(newFormat FormatV3) {
	if s.Format == nil || s.Format.Empty() {
		s.Format = newFormat
	}

	// current format is line and new format is span, just exit
	if !s.Format.IsSpan() && newFormat.IsSpan() {
		return
	}

	// current format is span and new format is line, overwrite
	if s.Format.IsSpan() && !newFormat.IsSpan() {
		s.Format = newFormat
	}

	// both formats are line formats
	if !s.Format.IsSpan() && !newFormat.IsSpan() {
		// only overwrite if the current format is a plain line
		if _, ok := s.Format.(FormatV3Line); ok {
			s.Format = newFormat
		}
	}

	// finally merge span formats
	if spanFormat, ok := s.Format.(FormatV3Span); ok {
		if newSpanFormat, ok := newFormat.(FormatV3Span); ok {
			for k, v := range newSpanFormat {
				spanFormat[k] = v
			}
		}

		s.Format = spanFormat
	}

	s.Format = s.Format.DropNull()
}

type PlainText []uint16

func (p *PlainText) String() string {
	return Uint16ToStr(*p)
}

func (p *PlainText) Len() int {
	return len(*p)
}

func (p *PlainText) WriteString(s string) {
	*p = append(*p, StrToUint16(s)...)
}

func ParseHtml(htmlContent string) (string, []TextSpan, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return "", nil, fmt.Errorf("error parsing HTML: %w", err)
	}

	stylesheets := extractStylesheets(doc)
	applyInlineStyles(doc, stylesheets)

	spans := []TextSpan{}
	plainText := PlainText{}
	traverse(doc, &plainText, &spans, "", -1)

	// Filter out the line formats
	for i := len(spans) - 1; i >= 0; i-- {
		span := spans[i]
		if _, ok := span.Format.(FormatV3Line); ok {
			spans = append(spans[:i], spans[i+1:]...)
		}
	}

	plaintext := plainText.String()
	if plaintext == "" {
		return "", nil, nil
	}

	if len(spans) == 0 && plaintext[len(plaintext)-1] == '\n' {
		plaintext = plaintext[:len(plaintext)-1]
	}

	return plaintext, spans, nil
}

func extractStylesheets(n *html.Node) []string {
	var stylesheets []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.DataAtom == atom.Style {
			var buf bytes.Buffer
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				html.Render(&buf, c)
			}
			stylesheets = append(stylesheets, buf.String())
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
	return stylesheets
}

func applyInlineStyles(n *html.Node, stylesheets []string) {
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			inlineStyle := getFullStyle(n, stylesheets)
			if inlineStyle != "" {
				// Add the style as an inline style attribute
				existingStyle := getAttribute(n, "style")
				if existingStyle != "" {
					inlineStyle = existingStyle + "; " + inlineStyle
				}
				setAttribute(n, "style", inlineStyle)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
}

func getAttribute(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

func setAttribute(n *html.Node, key, val string) {
	for i, a := range n.Attr {
		if a.Key == key {
			n.Attr[i].Val = val
			return
		}
	}
	n.Attr = append(n.Attr, html.Attribute{Key: key, Val: val})
}

func getFullStyle(n *html.Node, stylesheets []string) string {
	var styles string
	for _, sheet := range stylesheets {
		rules := parseCSS(sheet)
		for selector, declarations := range rules {
			if matchesSelector(n, selector) {
				styles += declarations
			}
		}
	}
	return styles
}

func parseCSS(css string) map[string]string {
	rules := make(map[string]string)
	css = strings.TrimSpace(css)
	lines := strings.Split(css, "}")
	for _, line := range lines {
		parts := strings.SplitN(line, "{", 2)
		if len(parts) == 2 {
			selector := strings.TrimSpace(parts[0])
			declarations := strings.TrimSpace(parts[1])
			rules[selector] = declarations
		}
	}
	return rules
}

func matchesSelector(n *html.Node, selector string) bool {
	compiledSelector, err := cascadia.Compile(selector)
	if err != nil {
		log.Printf("Error compiling selector %s: %v", selector, err)
		return false
	}
	return compiledSelector.Match(n)
}

func _containsLineFormat(span TextSpan, spans []TextSpan) (int, bool) {
	ix := bisectLeft(spans, span.StartIndex, func(s TextSpan) int { return s.StartIndex })

	for _, s := range spans[ix:] {
		if s.StartIndex >= span.EndIndex {
			break
		}

		if !s.Format.IsSpan() {
			return ix, true
		}
	}

	return ix, false
}

func traverse(node *html.Node, plainText *PlainText, spans *[]TextSpan, curType string, listDepth int) {
	if node.Type == html.ElementNode {
		if node.Data == "head" || node.Data == "meta" || node.Data == "style" || node.Data == "script" || node.Data == "title" {
			return
		}
	}

	// set block type here
	if node.Type == html.ElementNode {
		if node.Data == "ul" {
			curType = "ul"
			listDepth++
		} else if node.Data == "ol" {
			curType = "ol"
			listDepth++
		} else if node.Data == "blockquote" {
			curType = "blockquote"
		} else if node.Data == "pre" {
			curType = "pre"
		}
	}

	if node.Type == html.TextNode {
		if curType == "pre" {
			plainText.WriteString(node.Data)
		} else {
			isWhitespaceOnly, text := normalizeWhitespace(node.Data)
			if !isWhitespaceOnly || plainText.Len() > 0 {
				plainText.WriteString(text)
			}
		}
	}

	var span *TextSpan
	if node.Type == html.ElementNode {
		span = &TextSpan{StartIndex: plainText.Len()}
		applyElementStyles(node, span, curType, listDepth)
		applyCssStyles(node, span)
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		traverse(c, plainText, spans, curType, listDepth)
	}

	if node.Type == html.ElementNode {
		if span != nil && span.Format != nil && !span.Format.Empty() {
			span.EndIndex = plainText.Len()

			ix, yes := _containsLineFormat(*span, *spans)
			if !yes {
				if span.StartIndex == span.EndIndex && node.Data == "hr" {
					span.EndIndex++ // Need to point to after the newline for the format call in paste.go
				}

				*spans = InsertAt(*spans, ix, *span)

				_, isCodeBlock := span.Format.(FormatV3CodeBlock)

				if !span.Format.IsSpan() && !isCodeBlock {
					plainText.WriteString("\n")
				}
			}
		}
	}
}

func applyElementStyles(node *html.Node, span *TextSpan, curType string, listDepth int) {
	if node.Data == "b" || node.Data == "strong" {
		span.Format = FormatV3Span{"b": "true"}
	} else if node.Data == "i" || node.Data == "em" {
		span.Format = FormatV3Span{"i": "true"}
	} else if node.Data == "s" || node.Data == "del" {
		span.Format = FormatV3Span{"s": "true"}
	} else if node.Data == "u" {
		span.Format = FormatV3Span{"u": "true"}
	} else if node.Data == "a" {
		span.Format = FormatV3Span{"a": getAttribute(node, "href")}
	} else if node.Data == "li" && curType == "ul" {
		span.Format = FormatV3BulletList(max(0, listDepth))
	} else if node.Data == "li" && curType == "ol" {
		span.Format = FormatV3OrderedList(max(0, listDepth))
	} else if node.Data == "h1" {
		span.Format = FormatV3Header(1)
	} else if node.Data == "h2" {
		span.Format = FormatV3Header(2)
	} else if node.Data == "h3" {
		span.Format = FormatV3Header(3)
	} else if node.Data == "h4" {
		span.Format = FormatV3Header(4)
	} else if node.Data == "h5" {
		span.Format = FormatV3Header(5)
	} else if node.Data == "h6" {
		span.Format = FormatV3Header(6)
	} else if node.Data == "blockquote" {
		span.Format = FormatV3BlockQuote{}
	} else if node.Data == "code" && curType == "pre" {
		language := getAttribute(node, "data-language")
		span.Format = FormatV3CodeBlock(language)
	} else if node.Data == "code" {
		span.Format = FormatV3Span{"c": "true"}
	} else if node.Data == "p" && (curType == "ol" || curType == "ul") {
		span.Format = FormatV3IndentedLine(listDepth)
	} else if node.Data == "p" && curType == "blockquote" {
		span.Format = FormatV3BlockQuote{}
	} else if node.Data == "br" && curType == "" {
		span.Format = FormatV3Line{}
	} else if node.Data == "p" && curType == "" {
		span.Format = FormatV3Line{}
	} else if node.Data == "hr" {
		span.Format = FormatV3Rule{}
	}
}

type fontCategory struct {
	min    float64
	max    float64
	format FormatV3
}

var fontCategoriesPX = []fontCategory{
	{min: 0, max: 15, format: FormatV3Span{}},
	{min: 15, max: 20, format: FormatV3Header(3)},
	{min: 20, max: 26, format: FormatV3Header(2)},
	{min: 26, max: 50, format: FormatV3Header(1)},
}

var fontCategoriesPT = []fontCategory{
	{min: 0, max: 15, format: FormatV3Span{}},
	{min: 15, max: 20, format: FormatV3Header(3)},
	{min: 20, max: 26, format: FormatV3Header(2)},
	{min: 26, max: 50, format: FormatV3Header(1)},
}

func applyCssStyles(node *html.Node, span *TextSpan) {
	styles := parseStyles(getAttribute(node, "style"))
	if len(styles) == 0 {
		return
	}

	if textDecoration := styles["text-decoration"]; textDecoration != "" {
		if textDecoration == "underline" {
			span.MergeFormat(FormatV3Span{"u": "true"})
		} else if textDecoration == "line-through" {
			span.MergeFormat(FormatV3Span{"s": "true"})
		} else if textDecoration == "none" {
			span.MergeFormat(FormatV3Span{"u": "null", "s": "null"})
		}
	}

	if fontWeight := styles["font-weight"]; fontWeight != "" {
		if fontWeight == "bold" {
			span.MergeFormat(FormatV3Span{"b": "true"})
		} else if fontWeight == "normal" {
			span.MergeFormat(FormatV3Span{"b": "null"})
		}
	}

	if fontStyle := styles["font-style"]; fontStyle != "" {
		if fontStyle == "italic" {
			span.MergeFormat(FormatV3Span{"i": "true"})
		} else if fontStyle == "normal" {
			span.MergeFormat(FormatV3Span{"i": "null"})
		}
	}

	if fontSize := styles["font-size"]; fontSize != "" {
		if strings.HasSuffix(fontSize, "px") {
			applyFontSize(span, strings.TrimSuffix(fontSize, "px"), fontCategoriesPX)
		}
		if strings.HasSuffix(fontSize, "pt") {
			applyFontSize(span, strings.TrimSuffix(fontSize, "pt"), fontCategoriesPT)
		}
	}

	if font := styles["font"]; font != "" {
		parts := strings.Split(font, " ")
		for _, part := range parts {
			if strings.HasSuffix(part, "px") {
				applyFontSize(span, strings.TrimSuffix(part, "px"), fontCategoriesPX)
			}
			if strings.HasSuffix(part, "pt") {
				applyFontSize(span, strings.TrimSuffix(part, "pt"), fontCategoriesPT)
			}
		}
	}
}

func normalizeWhitespace(input string) (bool, string) {
	input = newlinePattern.ReplaceAllString(input, "")
	isOnlyWhitespace := len(strings.TrimSpace(input)) == 0
	return isOnlyWhitespace, spacePattern.ReplaceAllString(input, " ")
}

func applyFontSize(span *TextSpan, sizeStr string, categories []fontCategory) {
	size, err := strconv.Atoi(sizeStr)
	floatSize := float64(size)
	if err != nil {
		floatSize, err = strconv.ParseFloat(sizeStr, 32)
		if err != nil {
			return
		}
	}

	for _, category := range categories {
		if floatSize >= category.min && floatSize < category.max {
			span.MergeFormat(category.format)
			break
		}
	}
}

func parseStyles(style string) map[string]string {
	styles := make(map[string]string)
	rules := strings.Split(style, ";")
	for _, rule := range rules {
		rule = strings.TrimSpace(rule)

		if rule == "" {
			continue
		}

		parts := strings.SplitN(rule, ":", 2)
		if len(parts) != 2 {
			log.Warnf("Invalid style rule: %q", rule)
			continue
		}
		property := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		styles[property] = value
	}
	return styles
}

func ParsePaste(data []PasteItem) (string, []TextSpan, error) {

	return "", nil, nil
}
