package v3

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
	"golang.org/x/net/html"
)

func (r *Rogue) DisplayHtml(startID, endID ID, includeIDs, smartQuote bool) string {
	html, err := r.GetHtml(startID, endID, includeIDs, smartQuote)
	if err != nil {
		return fmt.Sprintf("failed to get html: %v", err)
	}

	return html
}

func (r *Rogue) DisplayAllHtml(includeIDs, smartQuote bool) string {
	firstID, err := r.GetFirstID()
	if err != nil {
		return fmt.Sprintf("failed to get first id: %v", err)
	}

	lastID, err := r.GetLastID()
	if err != nil {
		return fmt.Sprintf("failed to get last id: %v", err)
	}

	html, err := r.GetHtml(firstID, lastID, includeIDs, smartQuote)
	if err != nil {
		return fmt.Sprintf("failed to get html: %v", err)
	}

	return html
}

func (r *Rogue) GetHtml(startID, endID ID, includeIDs, smartQuote bool) (string, error) {
	vis, spanNOS, lineNOS, err := r.ToIndexNos(startID, endID, nil, smartQuote)
	if err != nil {
		return "", fmt.Errorf("r.ToIndexNos(%v, %v): %w", startID, endID, err)
	}

	if vis == nil {
		return "", nil
	}

	fVis := &FugueVis{
		Text: vis.Text,
		IDs:  vis.IDs,
	}

	return ToHtml(fVis, spanNOS, lineNOS, includeIDs)
}

func (r *Rogue) GetFullHtml(includeIDs, smartQuote bool) (string, error) {
	firstID, err := r.GetFirstID()
	if err != nil {
		return "", fmt.Errorf("r.GetFirstID(): %w", err)
	}

	lastID, err := r.GetLastID()
	if err != nil {
		return "", fmt.Errorf("r.GetLastID(): %w", err)
	}

	return r.GetHtml(firstID, lastID, includeIDs, smartQuote)
}

func (r *Rogue) AfterIDToEndID(afterID ID) (ID, error) {
	c, err := r.GetTotCharByID(afterID)
	if err != nil {
		return NoID, err
	}

	if c == '\n' {
		return afterID, nil
	}

	endID, err := r.TotLeftOf(afterID)
	if err != nil {
		return NoID, err
	}

	return endID, nil
}

func (r *Rogue) GetHtmlAt(startID, endID ID, address *ContentAddress, includeIDs, smartQuote bool) (string, error) {
	vis, spanNOS, lineNOS, err := r.ToIndexNos(startID, endID, address, smartQuote)
	if err != nil {
		return "", fmt.Errorf("r.ToIndexNos(%v, %v): %w", startID, endID, err)
	}

	if vis == nil {
		return "", nil
	}

	fVis := &FugueVis{
		Text: vis.Text,
		IDs:  vis.IDs,
	}

	return ToHtml(fVis, spanNOS, lineNOS, includeIDs)
}

func _spanTag(k string) string {
	if k == "bold" || k == "b" {
		return "strong"
	} else if k == "italic" || k == "i" {
		return "em"
	} else if k == "strike" || k == "s" {
		return "s"
	} else if k == "underline" || k == "u" {
		return "u"
	} else if k == "del" {
		return "del"
	} else if k == "ins" {
		return "ins"
	} else if k == "link" || k == "a" {
		return "a"
	} else if k == "c" {
		return "code"
	} else if k == "ql" || k == "qr" || k == "author" {
		return "span"
	}

	return ""
}

func _wrapTags(startID, endID *ID, text string, format FormatV3) string {
	if len(text) == 0 {
		return ""
	}

	attrs := ""
	if startID != nil {
		attrs = fmt.Sprintf(" data-rid=\"%s\"", startID)
	}

	if format == nil || !format.IsSpan() {
		return fmt.Sprintf("<span%s>%s</span>", attrs, html.EscapeString(text))
	}

	f := format.DropNull().(FormatV3Span)

	delete(f, "e")
	delete(f, "en")

	if _, ok := f["noid"]; ok {
		attrs = ""
		delete(f, "noid")
	}

	if len(f) == 0 {
		return fmt.Sprintf("<span%s>%s</span>", attrs, html.EscapeString(text))
	}

	lastIx := len(f) - 1
	openTags, closeTags := []string{}, []string{}
	for i, k := range MapSortedKeys(f) {
		tagAttrs := ""

		if k == "link" || k == "a" {
			tagAttrs = fmt.Sprintf("%s href=%q", tagAttrs, f[k])
		}

		if (k == "del" || k == "ins") && startID != nil && endID != nil {
			tagAttrs = fmt.Sprintf("%s data-delta-start=%q data-delta-end=%q", tagAttrs, startID, endID)
		}

		if (k == "author") && startID != nil && endID != nil {
			tagAttrs = fmt.Sprintf("%s data-author-prefix=%q", tagAttrs, f[k])
		}

		if k == "ql" {
			if text == "'" {
				text = "‘"
			} else if text == "\"" {
				text = "“"
			}
		}

		if k == "qr" {
			if text == "'" {
				text = "’"
			} else if text == "\"" {
				text = "”"
			}
		}

		tag := _spanTag(k)
		if tag == "" {
			// log.Errorf("unknown format: %s", k)
			continue
		}

		if i == lastIx {
			openTags = append(openTags, fmt.Sprintf("<%s%s%s>", tag, attrs, tagAttrs))
		} else {
			openTags = append(openTags, fmt.Sprintf("<%s%s>", tag, tagAttrs))

		}
		closeTags = append(closeTags, fmt.Sprintf("</%s>", tag))
	}

	Reverse(closeTags)
	text = strings.ReplaceAll(html.EscapeString(text), "\n", "<br>")
	return fmt.Sprintf("%s%s%s", strings.Join(openTags, ""), text, strings.Join(closeTags, ""))
}

func _printTags(tags []string, attributes []string) (open, close string) {
	if len(tags) == 0 {
		return "", ""
	}

	openTags, closeTags := []string{}, []string{}
	for i, t := range tags {
		if len(attributes) > 0 && i == len(tags)-1 {
			openTags = append(openTags, fmt.Sprintf("<%s %s>", t, strings.Join(attributes, " ")))
		} else {
			openTags = append(openTags, fmt.Sprintf("<%s>", t))
		}
		closeTags = append(closeTags, fmt.Sprintf("</%s>", t))
	}

	Reverse(closeTags)
	return strings.Join(openTags, ""), strings.Join(closeTags, "")
}

func isRawBlock(format FormatV3) bool {
	_, ok := format.(FormatV3CodeBlock)
	return ok
}

func _getBlockTags(startID *ID, format FormatV3) (open, close string) {
	var tags []string
	var attributes []string

	if format == nil {
		return "", ""
	}

	if format.IsSpan() {
		log.Errorf("this should be a line format %v", format)
		return "", ""
	}

	switch format := format.(type) {
	case FormatV3BulletList:
		tags = []string{"ul"}
	case FormatV3OrderedList:
		tags = []string{"ol"}
	case FormatV3BlockQuote:
		tags = []string{"blockquote"}
	case FormatV3CodeBlock:
		tags = []string{"pre", "code"}
		// only need to add id for a raw code block
		if startID != nil {
			attributes = append(attributes, fmt.Sprintf("data-rid=\"%s\"", startID))
		}

		language := string(format)
		if language != "" {
			attributes = append(attributes, fmt.Sprintf("class=\"language-%s\" data-language=\"%s\"", language, language))
		}
	case FormatV3Header, FormatV3Line, FormatV3Rule, FormatV3IndentedLine, FormatV3Image:
		return "", ""
	default:
		log.Errorf("unknown block format: %v %T", format, format)
		return "", ""
	}

	return _printTags(tags, attributes)
}

func imageTags(startID *ID, format FormatV3Image) (openTags, closeTags string) {
	attributes := []string{}
	attributes = append(attributes, fmt.Sprintf("src=\"%s\"", format.Src))

	style := ""
	if format.Width != "" {
		style = fmt.Sprintf("width: %s;", format.Width)
	}

	if format.Height != "" {
		style = fmt.Sprintf("%s height: %s;", style, format.Height)
	}

	if style != "" {
		attributes = append(attributes, fmt.Sprintf("style=\"%s\"", style))
	}

	if format.Alt != "" {
		attributes = append(attributes, fmt.Sprintf("alt=\"%s\"", format.Alt))
	}

	if startID != nil {
		openTags = fmt.Sprintf("<figure><img %s /><figcaption data-rid=\"%s\">", strings.Join(attributes, " "), startID)
	} else {
		openTags = fmt.Sprintf("<figure><img %s /><figcaption>", strings.Join(attributes, " "))
	}

	closeTags = "</figcaption></figure>"
	return openTags, closeTags
}

func _getLineTag(startID *ID, format FormatV3) (openTags, closeTags string) {
	tag := "p"
	attributes := []string{}

	if f, ok := format.(FormatV3Image); ok && f.Src != "" {
		return imageTags(startID, f)
	}

	if startID != nil {
		attributes = append(attributes, fmt.Sprintf("data-rid=\"%s\"", startID))
	}

	switch f := format.(type) {
	case FormatV3Header:
		tag = fmt.Sprintf("h%d", f)
	case FormatV3CodeBlock:
		return "", ""
	case FormatV3BlockQuote:
		tag = "p"
	case FormatV3OrderedList:
		tag = "li"
	case FormatV3BulletList:
		tag = "li"
	case FormatV3IndentedLine:
		tag = "p"
	case FormatV3Line:
		tag = "p"
	case FormatV3Rule:
		tag = "hr"
	default:
		tag = "p" // default to giving spans a paragraph tag for now
	}

	return _printTags([]string{tag}, attributes)
}

type listMeta struct {
	isList, isBullet, isOrdered, isIndent bool
	indent                                int
}

func getListMeta(format FormatV3) listMeta {
	switch f := format.(type) {
	case FormatV3OrderedList:
		return listMeta{true, false, true, false, int(f)}
	case FormatV3BulletList:
		return listMeta{true, true, false, false, int(f)}
	case FormatV3IndentedLine:
		return listMeta{true, false, false, true, int(f)}
	}

	return listMeta{}
}

func ToHtml(vis *FugueVis, spanNOS, lineNOS *NOS, includeIDs bool) (string, error) {
	prevIx := 0
	builder := strings.Builder{}
	prevLineFmts := []FormatV3{}
	openLineTag, closeLineTag := "", ""

	lines := lineNOS.tree.AsSlice()

	listStartIndent := 0
	// iteratre over each line span
	for _, line := range lines {
		// write the opening tag of the block if there is one
		var curLineFmt FormatV3 = line.Format

		var startID *ID = nil
		var lineID *ID = nil
		if includeIDs {
			startID = &vis.IDs[line.StartIx]
			if line.EndIx < len(vis.IDs) {
				lineID = &vis.IDs[line.EndIx]
			}
		}

		// HANDLE BLOCKS
		if len(prevLineFmts) == 0 {
			blockOpen, _ := _getBlockTags(startID, curLineFmt)
			builder.WriteString(blockOpen)
			if blockOpen != "" {
				prevLineFmts = append(prevLineFmts, curLineFmt)
			}

			switch cf := curLineFmt.(type) {
			case FormatV3OrderedList:
				listStartIndent = int(cf)
			case FormatV3BulletList:
				listStartIndent = int(cf)
			}
		} else {
			prevLineFmt := prevLineFmts[len(prevLineFmts)-1]

			if !curLineFmt.Equals(prevLineFmt) {
				prev := getListMeta(prevLineFmt)
				cur := getListMeta(curLineFmt)

				if !cur.isList || !prev.isList {
					listStartIndent = 0

					// close any open blocks since these are different line format types
					for i := len(prevLineFmts) - 1; i >= 0; i-- {
						prevLineFmt := prevLineFmts[i]
						_, blockClose := _getBlockTags(nil, prevLineFmt)
						builder.WriteString(blockClose)
						prevLineFmts = prevLineFmts[:i]
					}

					blockOpen, _ := _getBlockTags(startID, curLineFmt)
					builder.WriteString(blockOpen)
					if blockOpen != "" {
						prevLineFmts = append(prevLineFmts, curLineFmt)
					}
				} else if prev.indent < cur.indent {
					// further indent so hold existing list block open
					blockOpen, _ := _getBlockTags(startID, curLineFmt)
					builder.WriteString(blockOpen)
					if blockOpen != "" {
						prevLineFmts = append(prevLineFmts, curLineFmt)
					}
				} else if max(cur.indent, listStartIndent) <= max(prev.indent, listStartIndent) {
					// indent less than existing so close blocks until we reach the same or less indent
					var prevLineFmt FormatV3
					for i := len(prevLineFmts) - 1; i >= 0; i-- {
						prevLineFmt = prevLineFmts[i]
						prev := getListMeta(prevLineFmt)

						if max(cur.indent, listStartIndent) >= max(prev.indent, listStartIndent) {
							break
						}

						_, blockClose := _getBlockTags(nil, prevLineFmt)
						builder.WriteString(blockClose)
						prevLineFmts = prevLineFmts[:i]
					}

					// if adjacent lists with different types, close prev and open new
					if len(prevLineFmts) > 0 {
						i := len(prevLineFmts) - 1
						prevLineFmt := prevLineFmts[i]
						prev := getListMeta(prevLineFmt)
						if !cur.isIndent && prev.isBullet != cur.isBullet {
							_, blockClose := _getBlockTags(nil, prevLineFmt)
							builder.WriteString(blockClose)
							prevLineFmts = prevLineFmts[:i]

							if len(prevLineFmts) == 0 {
								listStartIndent = cur.indent // reset if this is a new list
							}
							blockOpen, _ := _getBlockTags(startID, curLineFmt)
							builder.WriteString(blockOpen)
							if blockOpen != "" {
								prevLineFmts = append(prevLineFmts, curLineFmt)
							}
						}
					}
				}
			}
		}
		// END HANDLE BLOCKS

		// just print the raw line contents if we're in a raw block
		// which is just a code-block for now
		if isRawBlock(curLineFmt) {
			eix := min(line.EndIx+1, len(vis.Text))
			content := html.EscapeString(Uint16ToStr(vis.Text[line.StartIx:eix]))
			builder.WriteString(content)
			prevIx = line.EndIx + 1
			continue
		}

		if _, ok := curLineFmt.(FormatV3Rule); ok {
			if includeIDs {
				builder.WriteString(fmt.Sprintf("<hr data-rid=\"%s\"/>", lineID))
			} else {
				builder.WriteString("<hr/>")
			}
			prevIx = line.EndIx + 1
			continue
		}

		// write the opening tag of the line if there is one
		openLineTag, closeLineTag = _getLineTag(lineID, curLineFmt)
		builder.WriteString(openLineTag)

		err := _htmlWriteLine(line, vis, spanNOS, includeIDs, &builder)
		if err != nil {
			return "", fmt.Errorf("_htmlWriteLine(%v, %v, %v, %v, %v): %w", line, vis, spanNOS, includeIDs, &builder, err)
		}

		// write the closing tag
		builder.WriteString(closeLineTag)
		prevIx = line.EndIx + 1
	}

	// close the last blocks
	for i := len(prevLineFmts) - 1; i >= 0; i-- {
		_, blockClose := _getBlockTags(nil, prevLineFmts[i])
		builder.WriteString(blockClose)
	}

	// write any remaining content without a trailing newline
	if prevIx < len(vis.Text) {
		var startID *ID = nil
		var endID *ID = nil
		if includeIDs {
			startID = &vis.IDs[prevIx]
			endID = &vis.IDs[len(vis.Text)-1]
		}

		content := Uint16ToStr(vis.Text[prevIx:])
		t := _wrapTags(startID, endID, content, nil)
		builder.WriteString(t)
	}

	return builder.String(), nil
}

func _htmlWriteLine(line *NOSNode, vis *FugueVis, spanNOS *NOS, includeIDs bool, builder *strings.Builder) error {
	// iterate over each span within the line
	prevIx := line.StartIx
	err := spanNOS.between(line.StartIx, line.EndIx-1, func(n *NOSNode) error {
		var startID *ID = nil
		var endID *ID = nil
		if includeIDs {
			startID = &vis.IDs[prevIx]
			endID = &vis.IDs[n.StartIx]
		}
		content := Uint16ToStr(vis.Text[prevIx:n.StartIx])
		t := _wrapTags(startID, endID, content, nil)
		builder.WriteString(t)

		startID = nil
		if includeIDs {
			startID = &vis.IDs[n.StartIx]
			endID = &vis.IDs[n.EndIx]
		}

		content = Uint16ToStr(vis.Text[n.StartIx : n.EndIx+1])
		t = _wrapTags(startID, endID, content, n.Format)
		builder.WriteString(t)

		prevIx = n.EndIx + 1

		return nil
	})

	if err != nil {
		return fmt.Errorf("between(%v, %v): %w", line.StartIx, line.EndIx, err)
	}

	// write any remaining unformated text before end of line
	if prevIx < line.EndIx {
		var startID *ID = nil
		var endID *ID = nil
		if includeIDs {
			startID = &vis.IDs[prevIx]
			if line.EndIx < len(vis.IDs) {
				endID = &vis.IDs[line.EndIx]
			}
		}

		content := Uint16ToStr(vis.Text[prevIx:line.EndIx])
		t := _wrapTags(startID, endID, content, nil)
		builder.WriteString(t)
	}

	return nil
}

// We made a change, show you what happend
func (r *Rogue) GetHtmlDiff(startID, endID ID, address *ContentAddress, includeIDs, smartQuote bool) (string, error) {
	filteredNOS, err := r.GetFilteredNOS(startID, endID, address, smartQuote)
	if err != nil {
		return "", err
	}

	diffText := FugueVis{
		IDs:  filteredNOS.IDs,
		Text: filteredNOS.Text,
	}

	return ToHtml(&diffText, filteredNOS.spanNos, filteredNOS.lineNos, includeIDs)
}

func (r *Rogue) GetHtmlXRay(startID, endID ID, includeIDs, smartQuote bool) (string, error) {
	vis, spanNOS, lineNOS, err := r.ToIndexNos(startID, endID, nil, smartQuote)
	if err != nil {
		return "", fmt.Errorf("r.ToIndexNos(%v, %v): %w", startID, endID, err)
	}

	if vis == nil {
		return "", nil
	}

	i := 1
	start := 0
	curr := vis.IDs[start].Author[0:1]
	for i < len(vis.IDs) {
		id := vis.IDs[i]
		if id.Author[0:1] == curr {
			i++
			continue
		}

		spanNOS.Insert(NOSNode{
			StartIx: start,
			EndIx:   i - 1,
			Format:  FormatV3Span{"author": curr},
		})
		start = i
		curr = id.Author[0:1]
	}

	if start < len(vis.IDs) {
		spanNOS.Insert(NOSNode{
			StartIx: start,
			EndIx:   len(vis.IDs) - 1,
			Format:  FormatV3Span{"author": curr},
		})
	}

	fVis := &FugueVis{
		Text: vis.Text,
		IDs:  vis.IDs,
	}

	return ToHtml(fVis, spanNOS, lineNOS, includeIDs)
}

func (r *Rogue) GetHtmlDiffBetween(startID, endID ID, fromAddress, toAddress *ContentAddress, includeIDs, smartQuote bool) (string, error) {
	filteredNOS, err := r.GetFilteredNOSBetween(startID, endID, fromAddress, toAddress, smartQuote)
	if err != nil {
		return "", fmt.Errorf("GetFilteredNOSBetween(%v, %v, %v, %v): %w", startID, endID, fromAddress, toAddress, err)
	}

	diffText := FugueVis{
		IDs:  filteredNOS.IDs,
		Text: filteredNOS.Text,
	}

	return ToHtml(&diffText, filteredNOS.spanNos, filteredNOS.lineNos, includeIDs)
}
