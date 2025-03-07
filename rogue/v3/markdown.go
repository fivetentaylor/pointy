package v3

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/sergi/go-diff/diffmatchpatch"
	mdParse "github.com/fivetentaylor/pointy/rogue/v3/gmparse"
	"github.com/fivetentaylor/pointy/rogue/v3/set"
	"golang.org/x/exp/constraints"
	"golang.org/x/net/html"
)

type MarkdownStyleDiff struct {
	Old *mdParse.FormatSpan
	New *mdParse.FormatSpan
}

func (rogueStyleDiff MarkdownStyleDiff) String() string {
	return fmt.Sprintf("%s -> %s", rogueStyleDiff.Old, rogueStyleDiff.New)
}

func spansToUTF16(plaintext string, spans []mdParse.FormatSpan) []mdParse.FormatSpan {
	out := make([]mdParse.FormatSpan, len(spans))
	for i, span := range spans {
		// fix for when markdown doesn't end in a newline
		// but has a line format on the last line
		startIx := min(span.Start, len(plaintext))
		endIx := min(span.End, len(plaintext))
		out[i] = mdParse.FormatSpan{
			Start:  UTF16Length(plaintext[:startIx]),
			End:    UTF16Length(plaintext[:endIx]),
			Format: span.Format,
		}
	}
	return out
}

type MarkdownDiff struct {
	TextDiffs    []diffmatchpatch.Diff
	OldFormats   []mdParse.FormatSpan
	NewFormats   []mdParse.FormatSpan
	OldPlaintext string
	NewPlaintext string
}

func DiffMarkdown(oldPlaintext, newMd string) (MarkdownDiff, error) {
	newPlaintext, newFormats, err := mdParse.SplitMarkdown(newMd)
	if err != nil {
		return MarkdownDiff{}, fmt.Errorf("SplitMarkdown(%q): %w", newMd, err)
	}

	newPlaintext, newFormats = mdParse.AlignWhitespace(oldPlaintext, newPlaintext, newFormats)

	diffs := DiffWords(oldPlaintext, newPlaintext)

	// oldFormats = spansToUTF16(oldPlaintext, oldFormats)
	newFormats = spansToUTF16(newPlaintext, newFormats)

	markdownDiff := MarkdownDiff{
		TextDiffs: diffs,
		// OldFormats:   oldFormats,
		NewFormats:   newFormats,
		OldPlaintext: oldPlaintext,
		NewPlaintext: newPlaintext,
	}

	return markdownDiff, nil
}

type Number interface {
	constraints.Integer | constraints.Float
}

var spanSymbols = map[string]string{
	"b": "**", // bold
	"i": "*",  // italic
	// "u": "__", // underline
	"s": "~~", // strike
	"a": "[",  // link
	"c": "`",  // code
}

func (r *Rogue) lineFormatOpToNOS(startIx, endIx int, op FormatOp) ([]NOSNode, error) {
	if op.StartID != op.EndID {
		log.Errorf("line format with non-zero length: %v", op)
		return nil, ErrorEmptySpan{}
	}

	formatStartIx, _, err := r.Rope.GetIndex(op.StartID)
	if err != nil {
		return nil, fmt.Errorf("GetIndex(%v): %w", op.StartID, err)
	}

	if formatStartIx < 0 {
		return nil, ErrorEmptySpan{}
	}

	return []NOSNode{{
		StartIx: formatStartIx - startIx,
		EndIx:   formatStartIx - startIx,
		Format:  op.Format,
	}}, nil
}

func (r *Rogue) formatOpToNOS(startIx, endIx int, op FormatOp) ([]NOSNode, error) {
	span, ok := op.Format.(FormatV3Span)
	if !ok {
		return r.lineFormatOpToNOS(startIx, endIx, op)
	}

	nos := []NOSNode{}
	sticky, noSticky := span.SplitSticky()

	formatStartIx, _, err := r.Rope.GetIndex(op.StartID)
	if err != nil {
		return nil, fmt.Errorf("GetIndex(%v): %w", op.StartID, err)
	}

	// if the start index is deleted
	if formatStartIx < 0 {
		startID, err := r.VisRightOf(op.StartID)
		if err != nil {
			// Span is empty,
			if errors.As(err, &ErrorNoRightVisSibling{}) {
				return nos, nil
			}
			return nil, fmt.Errorf("VisRightOf(%v): %w", op.StartID, err)
		}

		formatStartIx, _, err = r.Rope.GetIndex(startID)
		if err != nil {
			return nil, fmt.Errorf("GetIndex(%v): %w", startID, err)
		}
	}

	if len(sticky) > 0 {
		endID, err := r.VisLeftOf(op.EndID)
		if err != nil {
			if !errors.As(err, &ErrorNoLeftVisSibling{}) {
				return nil, fmt.Errorf("VisLeftOf(%v): %w", op.EndID, err)
			}
		} else {
			formatEndIx, _, err := r.Rope.GetIndex(endID)
			if err != nil {
				return nil, fmt.Errorf("GetIndex(%v): %w", endID, err)
			}

			formatStartIx = max(formatStartIx, startIx)
			formatEndIx = min(formatEndIx, endIx)

			if formatStartIx > formatEndIx {
				return nil, ErrorEmptySpan{}
			}

			nos = append(nos, NOSNode{
				StartIx: formatStartIx - startIx,
				EndIx:   formatEndIx - startIx,
				Format:  sticky,
			})
		}
	}

	if len(noSticky) > 0 {
		formatEndIx, _, err := r.Rope.GetIndex(op.EndID)
		if err != nil {
			return nil, fmt.Errorf("GetIndex(%v): %w", op.EndID, err)
		}

		if formatEndIx < 0 {
			endID, err := r.VisLeftOf(op.EndID)
			if err != nil {
				if errors.As(err, &ErrorNoLeftVisSibling{}) {
					return nos, nil
				}
				return nil, fmt.Errorf("VisLeftOf(%v): %w", op.EndID, err)
			}

			formatEndIx, _, err = r.Rope.GetIndex(endID)
			if err != nil {
				return nil, fmt.Errorf("GetIndex(%v): %w", endID, err)
			}
		}

		formatStartIx = max(formatStartIx, startIx)
		formatEndIx = min(formatEndIx, endIx)

		if formatStartIx > formatEndIx {
			return nil, ErrorEmptySpan{}
		}

		nos = append(nos, NOSNode{
			StartIx: formatStartIx - startIx,
			EndIx:   formatEndIx - startIx,
			Format:  noSticky,
		})
	}

	return nos, nil
}

func (r *Rogue) FormatSelectionToNOS(fs *FormatSelection) (span, line *NOS, err error) {
	span, err = r.FormatOpsToNOS(fs)
	if err != nil {
		return nil, nil, fmt.Errorf("FormatOpsToNOS(fs): %w", err)
	}

	line = StringToLineNOS(fs.Vis.Text)

	// inesert empty format for unformatted newlines or fix newlines
	// with an invalid line format, this helps later when we render
	// the document as html or markdown
	err = line.tree.Dft(func(n *NOSNode) error {
		sn, err := span.tree.Get(n.EndIx)
		if err != nil {
			return fmt.Errorf("Get(%v): %w", n.EndIx, err)
		}

		// if there's no line format, or the line format is not a valid line format
		if sn == nil || sn.Format.IsSpan() {
			err := span.Insert(NOSNode{n.EndIx, n.EndIx, FormatV3Line{}})
			if err != nil {
				return fmt.Errorf("Insert(%v): %w", n.EndIx, err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, nil, fmt.Errorf("line.tree.Dft(...): %w", err)
	}

	return span, line, nil
}

func (r *Rogue) FormatOpsToNOS(fs *FormatSelection) (*NOS, error) {
	startIx, endIx := fs.StartIx, fs.EndIx
	formatOps := fs.FormatOps
	text := fs.Vis.Text

	nos := NewNOS()

	for _, op := range formatOps {
		newNodes, err := r.formatOpToNOS(startIx, endIx, op)
		if err != nil {
			if errors.As(err, &ErrorEmptySpan{}) {
				continue
			}
			return nil, fmt.Errorf("formatOpToNOS(%v, %v, %v): %w", startIx, endIx, op, err)
		}

		for _, newNode := range newNodes {
			if newNode.IsLineFormat() && text[newNode.EndIx] != '\n' {
				continue
			}

			nos.Insert(newNode)
		}
	}

	return nos, nil
}

func (r *Rogue) DisplayMarkdown(startID, endID ID) string {
	md, err := r.GetMarkdownBeforeAfter(startID, endID)
	if err != nil {
		return fmt.Sprintf("failed to get markdown: %v", err)
	}
	return md
}

type FormatSelection struct {
	BeforeID  ID
	AfterID   ID
	StartID   ID
	EndID     ID
	StartIx   int
	EndIx     int
	FormatOps []FormatOp
	Vis       *FugueVis
}

func (r *Rogue) _getOriginalSelection(beforeID, afterID ID) (startID, endID ID, err error) {
	startID, err = r.TotRightOf(beforeID)
	if err != nil {
		if errors.As(err, &ErrorNoRightTotSibling{}) {
			return beforeID, afterID, nil
		}
		return NoID, NoID, fmt.Errorf("VisRightOf(%v): %w", beforeID, err)
	}

	// If the cursor is at the end of the line, we need to include the newline
	endID = afterID
	c, err := r.GetCharByID(afterID)
	if err != nil || c != '\n' {
		endID, err = r.TotLeftOf(afterID)
		if err != nil {
			if errors.As(err, &ErrorNoLeftTotSibling{}) {
				return beforeID, afterID, nil
			}
			return NoID, NoID, fmt.Errorf("VisLeftOf(%v): %w", afterID, err)
		}
	}

	return startID, endID, nil
}

func diffFormats(old, new FormatV3) (added, removed FormatV3Span) {
	added = FormatV3Span{}
	removed = FormatV3Span{}

	oldSpan, ok := old.(FormatV3Span)
	if !ok {
		return added, removed
	}

	newSpan, ok := new.(FormatV3Span)
	if !ok {
		return added, removed
	}

	for k, v := range oldSpan {
		if newV, ok := newSpan[k]; !ok || newV != v {
			removed[k] = v
		}
	}

	for k, v := range newSpan {
		if oldV, ok := oldSpan[k]; !ok || oldV != v {
			added[k] = v
		}
	}

	return added, removed
}

func formatToMd(f FormatV3) (out set.Set[string]) {
	out = set.Set[string]{}

	span, ok := f.(FormatV3Span)
	if !ok {
		return out
	}

	for k := range span {
		if s, ok := spanSymbols[k]; ok {
			out.Add(s)
		} else {
			// log.Errorf("unknown format: %s", k)
		}
	}

	return out
}

type MarkdownBuilder struct {
	actions *NOSActions
	pending *PendingActions
}

func NewMarkdownBuilder() *MarkdownBuilder {
	return &MarkdownBuilder{
		actions: NewNOSActions(),
		pending: NewPendingActions(),
	}
}

func (mb *MarkdownBuilder) toAddRemove(prev, next *NOSNode) (set.Set[string], set.Set[string]) {
	var addedSymbols set.Set[string]
	var removedSymbols set.Set[string]

	if next == nil {
		// last node
		removedSymbols = mb.pending.Symbols()
	} else if prev == nil {
		// first node
		addedSymbols = formatToMd(next.Format)
	} else if next.IsLineFormat() {
		// close all pending spans if the next node is a line format
		removedSymbols = mb.pending.Symbols()
	} else if prev.IsLineFormat() {
		// open all new spans if the previous node is a line format
		addedSymbols = formatToMd(next.Format)
	} else if prev.EndIx == next.StartIx-1 {
		// format spans touch
		added, removed := diffFormats(prev.Format, next.Format)
		addedSymbols = formatToMd(added)
		removedSymbols = formatToMd(removed)
	} else {
		// format spans are disjoint
		addedSymbols = formatToMd(next.Format)
		removedSymbols = formatToMd(prev.Format)
	}

	return addedSymbols, removedSymbols
}

func (mb *MarkdownBuilder) render(startIx, endIx int, text []uint16, builder *strings.Builder) error {
	prevIx := startIx

	it := mb.actions.Tree.Iterator()
	for it.Next() {
		ix, syms := it.Key().(int), it.Value().([]string)
		content := Uint16ToStr(text[prevIx:ix])
		content = EscapeMarkdownSyms(content)
		content = html.EscapeString(content)
		_, err := builder.WriteString(content)
		if err != nil {
			return fmt.Errorf("builder.WriteString(%v): %w", text[prevIx:ix], err)
		}

		sym := strings.Join(syms, "")
		_, err = builder.WriteString(sym)
		if err != nil {
			return fmt.Errorf("builder.WriteString(%v): %w", sym, err)
		}
		prevIx = ix
	}

	content := Uint16ToStr(text[prevIx:endIx])
	content = EscapeMarkdownSyms(content)
	content = html.EscapeString(content)
	_, err := builder.WriteString(content)
	if err != nil {
		return fmt.Errorf("builder.WriteString(%v): %w", text[prevIx:endIx], err)
	}

	return nil
}

func findFirstLowercaseMatch(input string) string {
	pattern := `[a-z]*` // Matches one or more lowercase letters.
	re := regexp.MustCompile(pattern)
	match := re.FindString(input)
	return match
}

func getMDBlockTag(format FormatV3) (open, close string, err error) {
	switch f := format.(type) {
	case FormatV3CodeBlock:
		ls := findFirstLowercaseMatch(string(f))
		return fmt.Sprintf("```%s\n", ls), "```\n", nil
	}

	return "", "", nil
}

func getMDLineTag(format FormatV3) (open, close string, err error) {
	if format == nil {
		return "", "\n\n", nil
	}

	switch f := format.(type) {
	case FormatV3Line:
		return "", "\n\n", nil
	case FormatV3CodeBlock:
		return "", "\n", nil
	case FormatV3BlockQuote:
		return "> ", "\n", nil
	case FormatV3Header:
		return fmt.Sprintf("%s ", strings.Repeat("#", int(f))), "\n\n", nil
	case FormatV3OrderedList:
		open := strings.Repeat("   ", int(f)) + "1. "
		return open, "\n", nil
	case FormatV3BulletList:
		open := strings.Repeat("  ", int(f)) + "- "
		return open, "\n", nil
	case FormatV3IndentedLine:
		open := strings.Repeat("  ", int(f)) + "   " // 3 spaces which will work for either "-" or "1."
		return open, "\n", nil
	case FormatV3Rule:
		return "---", "\n\n", nil
	case FormatV3Image:
		openTags, closeTags := imageTags(nil, f)
		return openTags, fmt.Sprintf("%s\n\n", closeTags), nil
	default:
		return "", "", fmt.Errorf("unknown line format: %v", format)
	}
}

func (mb *MarkdownBuilder) addActionsV2(prev, next *NOSNode) error {
	addedSymbols, removedSymbols := mb.toAddRemove(prev, next)
	toAdd := []pendAction{}
	toRemove := []pendAction{}

	// add span format symbols
	it := mb.pending.Tree.Iterator()
	for it.End(); it.Prev(); {
		if removedSymbols.Size() == 0 {
			break
		}

		eix := prev.EndIx + 1

		// syms represents the symbols added at a previous index
		// removedSymbols represents spans that are getting closed
		// at the current index. If syms covers removedSymbols,
		// then we only wan't to close the spans represented by
		// removedSymbols
		six, syms := it.Key().(int), it.Value().(set.Set[string])
		if syms.Covers(removedSymbols) {
			syms = removedSymbols
		}

		for s := range syms {
			// add the new spans to the actions
			mb.actions.Insert(six, s)
			if s == "[" {
				// handle links
				if f, ok := prev.Format.(FormatV3Span); ok {
					if href, ok := f["a"]; ok {
						ls := strings.Trim(href, "\"")
						linkClose := fmt.Sprintf("](%s)", ls)
						mb.actions.Insert(eix, linkClose)
					}
				}
			} else {
				mb.actions.Insert(eix, s)
			}

			// if two spans have disjoint overlap then we need
			// to close the latter one and reopen after the former
			if !removedSymbols.Pop(s) && next != nil {
				toAdd = append(toAdd, pendAction{next.StartIx, s})
			}

			// cleanup any symbols that have been written
			toRemove = append(toRemove, pendAction{six, s})
		}
	}

	// apply all the changes to pending
	for _, a := range toAdd {
		mb.pending.Add(a.ix, a.s)
	}

	for _, a := range toRemove {
		mb.pending.Remove(a.ix, a.s)
	}

	for s := range addedSymbols {
		mb.pending.Add(next.StartIx, s)
	}

	return nil
}

func (r *Rogue) GetFullMarkdown() (string, error) {
	startID, err := r.GetFirstID()
	if err != nil {
		return "", fmt.Errorf("r.GetFirstID(): %w", err)
	}

	endID, err := r.GetLastID()
	if err != nil {
		return "", fmt.Errorf("r.GetLastID(): %w", err)
	}

	return r.GetMarkdown(startID, endID)
}

func (r *Rogue) GetMarkdown(startID, endID ID) (string, error) {
	vis, spanNOS, lineNOS, err := r.ToIndexNos(startID, endID, nil, false)
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

	return ToMarkdown(fVis, spanNOS, lineNOS)
}

func (r *Rogue) GetMarkdownVis(startIdx, endIdx int) (string, error) {
	startId, err := r.Rope.GetVisID(startIdx)
	if err != nil {
		return "", fmt.Errorf("r.Rope.GetVisID(%v): %w", startIdx, err)
	}

	endID, err := r.Rope.GetVisID(endIdx)
	if err != nil {
		return "", fmt.Errorf("r.Rope.GetVisID(%v): %w", endIdx, err)
	}

	return r.GetMarkdown(startId, endID)
}

func (r *Rogue) GetMarkdownBeforeAfter(beforeID, afterID ID) (string, error) {
	startID, endID, err := r._getOriginalSelection(beforeID, afterID)
	if err != nil {
		return "", fmt.Errorf("GetOriginalSelection(%v, %v): %v", beforeID, afterID, err)
	}

	return r.GetMarkdown(startID, endID)
}

func (r *Rogue) GetMarkdownAt(startID, endID ID, address ContentAddress) (string, error) {
	vis, spanNOS, lineNOS, err := r.ToIndexNos(startID, endID, &address, false)
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

	return ToMarkdown(fVis, spanNOS, lineNOS)
}

func ToMarkdown(vis *FugueVis, spanNOS, lineNOS *NOS) (string, error) {
	var err error
	prevIx := 0
	builder := strings.Builder{}
	curBlockClose := ""

	err = lineNOS.tree.Dft(func(line *NOSNode) error {
		// write the opening tag of the block if there is one
		newBlockOpen, newBlockClose, err := getMDBlockTag(line.Format)
		if err != nil {
			return fmt.Errorf("getBlockTag(%v): %w", line.Format, err)
		}

		if curBlockClose != newBlockClose {
			builder.WriteString(curBlockClose)
			builder.WriteString(newBlockOpen)
			curBlockClose = newBlockClose
		}

		// just print the raw line contents if we're in a raw block
		// which is just a code-block for now
		if isRawBlock(line.Format) {
			content := Uint16ToStr(vis.Text[line.StartIx : line.EndIx+1])
			// content = EscapeMarkdownSyms(content) // don't escape code blocks
			content = html.EscapeString(content)
			builder.WriteString(content)
			prevIx = line.EndIx + 1
			return nil
		}

		// write the opening tag of the line if there is one
		openTag, closeTag, err := getMDLineTag(line.Format)
		if err != nil {
			return fmt.Errorf("getMDLineTag(%v): %w", line.Format, err)
		}

		_, err = builder.WriteString(openTag)
		if err != nil {
			return fmt.Errorf("builder.WriteString(%v): %w", openTag, err)
		}

		// iterate over each span within the line
		mb := NewMarkdownBuilder()
		err = spanNOS.betweenPairs(line.StartIx, line.EndIx-1, func(prev, next *NOSNode) error {
			return mb.addActionsV2(prev, next)
		})
		if err != nil {
			if !errors.As(err, &ErrorStopIteration{}) {
				return fmt.Errorf("betweenPairs(%v, %v): %w", line.StartIx, line.EndIx, err)
			}
		}

		// render the line
		err = mb.render(line.StartIx, line.EndIx, vis.Text, &builder)
		if err != nil {
			return fmt.Errorf("render(%v, %v, %v): %w", line.StartIx, line.EndIx, vis.Text, err)
		}

		// write the closing tag
		_, err = builder.WriteString(closeTag)
		if err != nil {
			return fmt.Errorf("builder.WriteString(%v): %w", closeTag, err)
		}

		prevIx = line.EndIx + 1

		return nil
	})

	if err != nil {
		if !errors.As(err, &ErrorStopIteration{}) {
			return "", fmt.Errorf("lineNOS.tree.Dft(...): %w", err)
		}
	}

	// close the last block
	_, err = builder.WriteString(curBlockClose)
	if err != nil {
		return "", fmt.Errorf("builder.WriteString(%v): %w", curBlockClose, err)
	}

	// write any remaining content without a trailing newline
	if prevIx < len(vis.Text) {
		endIx := len(vis.Text)

		mb := NewMarkdownBuilder()
		err = spanNOS.betweenPairs(prevIx, endIx-1, func(prev, next *NOSNode) error {
			mb.addActionsV2(prev, next)
			return nil
		})

		if err != nil {
			return "", fmt.Errorf("between(%v, %v): %w", prevIx, endIx, err)
		}

		err = mb.render(prevIx, endIx, vis.Text, &builder)
		if err != nil {
			return "", fmt.Errorf("render(%v, %v, %v): %w", prevIx, endIx, vis.Text, err)
		}

	}

	return builder.String(), nil
}

func EscapeMarkdownSyms(md string) string {
	var escaped strings.Builder
	escaped.Grow(len(md) * 2)

	// escape span special chars
	for _, r := range md {
		switch r {
		case '\\', '*', '_', '~':
			escaped.WriteRune('\\')
		}
		escaped.WriteRune(r)
	}

	// escape line formats
	pattern := "(?m)^( *- | *\\d+\\. |> |#{1,6} |[ \\t]*```[a-zA-Z]*)"
	r := regexp.MustCompile(pattern)

	return r.ReplaceAllStringFunc(escaped.String(), func(match string) string {
		return "\\" + match
	})
}

func UnescapeMarkdownSyms(escapedMd string) string {
	var unescaped strings.Builder
	unescaped.Grow(len(escapedMd))

	for i := 0; i < len(escapedMd); i++ {
		if escapedMd[i] == '\\' && i+1 < len(escapedMd) {
			i++
		}
		unescaped.WriteByte(escapedMd[i])
	}

	return unescaped.String()
}

func (r *Rogue) GetLineMarkdown(id ID) (beforeID, afterID ID, md string, err error) {
	startID, endID, _, err := r.GetLineAt(id, nil)
	if err != nil {
		return NoID, NoID, "", err
	}

	md, err = r.GetMarkdown(startID, endID)
	if err != nil {
		return NoID, NoID, "", err
	}

	beforeID, err = r.TotLeftOf(startID)
	if err != nil {
		return NoID, NoID, "", err
	}

	afterID, err = r.TotRightOf(endID)
	if err != nil {
		if !errors.As(err, &ErrorNoRightTotSibling{}) {
			return NoID, NoID, "", err
		}

		afterID = endID
	}

	return beforeID, afterID, md, nil
}

func (r *Rogue) InsertMarkdown(ix int, md string) (MultiOp, Actions, error) {
	mop := MultiOp{}
	offset := ix

	plaintext, formats, err := mdParse.SplitMarkdown(md)
	if err != nil {
		return mop, nil, err
	}

	newFormats := spansToUTF16(plaintext, formats)

	actions := make(Actions, 0)

	op, err := r.Insert(ix, plaintext)
	if err != nil {
		return mop, nil, err
	}
	mop = mop.Append(op)
	actions = append(actions, InsertAction{
		Index: ix,
		Text:  plaintext,
	})

	lineOps, lineActions, err := r._applyFormatDiffs(offset, newFormats)
	if err != nil {
		return mop, nil, err
	}

	mop = mop.Append(lineOps)
	actions = append(actions, lineActions...)
	r.OpIndex.Put(mop)

	return mop, actions, nil
}
