package v3

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/charmbracelet/log"
	"github.com/sergi/go-diff/diffmatchpatch"
	mdParse "github.com/fivetentaylor/pointy/rogue/v3/gmparse"
)

func _mdStyleToNull(ms *mdParse.FormatSpan) (FormatV3, error) {
	format, err := MapToFormatV3(ms.Format)
	if err != nil {
		return nil, fmt.Errorf("MapToFormatV3(%v): %w", ms.Format, err)
	}

	switch f := format.(type) {
	case FormatV3Span:
		for k := range f {
			f[k] = ""
		}
		return f, nil
	default:
		return FormatV3Line{}, nil
	}
}

func (r *Rogue) _applyFormatDiffs(offset int, formats []mdParse.FormatSpan) (Op, Actions, error) {
	actions := Actions{}

	mop := MultiOp{}
	for _, d := range formats {
		fLen := max(d.End-d.Start, 1) // handle line formats with min(1)
		rStyle, err := MapToFormatV3(d.Format)
		if err != nil {
			return nil, nil, fmt.Errorf("mapToFormatV3(%v): %w", d.Format, err)
		}

		fop, err := r.Format(offset+d.Start, fLen, rStyle)
		if err != nil {
			return nil, nil, fmt.Errorf("Format(%v, %v, %v): %w", d.Start, fLen, rStyle, err)
		}
		mop = mop.Append(fop)
		actions = append(actions, FormatAction{Index: offset + d.Start, Length: fLen, Format: rStyle})
	}

	return FlattenMop(mop), actions, nil
}

// hack to maintain [q, 1] as always the end of the document
// I might make this standard behavior for the core rogue.Insert function
func (r *Rogue) _handleTailInsert(ix int, text string) (InsertOp, error) {
	if ix == r.VisSize {
		if len(text) > 0 {
			if text[len(text)-1] == '\n' {
				return r.Insert(ix-1, fmt.Sprintf("\n%s", text[0:len(text)-1]))
			} else {
				return r.Insert(ix-1, fmt.Sprintf("\n%s", text))
			}
		}
	}

	return r.Insert(ix, text)
}

func (r *Rogue) _trimForClearFormats(startIx, endIx int, toTrimStart, toTrimEnd ID) (tStartID, tEndID ID, isVis bool, err error) {
	_, tStartIx, err := r.Rope.GetIndex(toTrimStart)
	if err != nil {
		return NoID, NoID, false, fmt.Errorf("GetIndex(%v): %w", toTrimStart, err)
	}

	_, tEndIx, err := r.Rope.GetIndex(toTrimEnd)
	if err != nil {
		return NoID, NoID, false, fmt.Errorf("GetIndex(%v): %w", toTrimEnd, err)
	}

	startIx = max(startIx, tStartIx)
	endIx = min(endIx, tEndIx)

	if endIx < startIx {
		return NoID, NoID, false, fmt.Errorf("totEndIx < totStartIx: %d < %d", endIx, startIx)
	}

	tStartID, err = r.Rope.GetTotID(startIx)
	if err != nil {
		return NoID, NoID, false, fmt.Errorf("GetTotID(%d): %w", startIx, err)
	}

	tEndID, err = r.Rope.GetTotID(endIx)
	if err != nil {
		return NoID, NoID, false, fmt.Errorf("GetTotID(%d): %w", endIx, err)
	}

	isVis = true
	visStartID, err := r.VisRightOf(tStartID)
	if err != nil {
		if !errors.As(err, &ErrorNoRightVisSibling{}) {
			return NoID, NoID, false, fmt.Errorf("VisRightOf(%v): %w", tStartID, err)
		}
		isVis = false
	}

	visEndID, err := r.VisLeftOf(tEndID)
	if err != nil {
		if !errors.As(err, &ErrorNoLeftVisSibling{}) {
			return NoID, NoID, false, fmt.Errorf("VisLeftOf(%v): %w", tEndID, err)
		}
		isVis = false
	}

	if isVis {
		visStartIx, _, err := r.Rope.GetIndex(visStartID)
		if err != nil {
			return NoID, NoID, false, fmt.Errorf("GetIndex(%v): %w", visStartID, err)
		}

		visEndIx, _, err := r.Rope.GetIndex(visEndID)
		if err != nil {
			return NoID, NoID, false, fmt.Errorf("GetIndex(%v): %w", visEndID, err)
		}

		isVis = visStartIx <= visEndIx
	}

	return tStartID, tEndID, isVis, nil
}

func (r *Rogue) ClearFormats(startID, endID ID) (mop MultiOp, err error) {
	_, startIx, err := r.Rope.GetIndex(startID)
	if err != nil {
		return mop, fmt.Errorf("GetIndex(%v): %w", startID, err)
	}

	_, endIx, err := r.Rope.GetIndex(endID)
	if err != nil {
		return mop, fmt.Errorf("GetIndex(%v): %w", endID, err)
	}

	err = r.NOS.Sticky.Tree.Slice(startIx, endIx, func(node *NOSV2Node) error {
		tStartID, tEndID, isVis, err := r._trimForClearFormats(startIx, endIx, node.StartID, node.EndID)
		if err != nil {
			log.Errorf("trim failed for node: %v with err: %v\n", node, err)
			return nil
		}

		if !isVis {
			return nil
		}

		// these are sticky spans, so if it was trimmed make sure you maintain stickiness
		if tEndID != node.EndID {
			tEndID, err = r.TotRightOf(tEndID)
			if err != nil {
				return fmt.Errorf("TotRightOf(%v): %w", tEndID, err)
			}
		}

		mop = mop.Append(FormatOp{
			ID:      r.NextID(1),
			StartID: tStartID,
			EndID:   tEndID,
			Format:  FormatV3Span{"e": "true"},
		})

		return nil
	})

	if err != nil {
		return mop, fmt.Errorf("NOS.sticky.tree.Slice(%d, %d): %w", startIx, endIx, err)
	}

	err = r.NOS.NoSticky.Tree.Slice(startIx, endIx, func(node *NOSV2Node) error {
		tStartID, tEndID, isVis, err := r._trimForClearFormats(startIx, endIx, node.StartID, node.EndID)
		if err != nil {
			log.Errorf("trim failed for node: %v with err: %v\n", node, err)
			return nil
		}

		if !isVis {
			return nil
		}

		mop = mop.Append(FormatOp{
			ID:      r.NextID(1),
			StartID: tStartID,
			EndID:   tEndID,
			Format:  FormatV3Span{"e": "true"},
		})

		return nil
	})

	if err != nil {
		return mop, fmt.Errorf("NOS.noSticky.tree.Slice(%d, %d): %w", startIx, endIx, err)
	}

	err = r.NOS.Line.Tree.Slice(startIx, endIx, func(node *NOSV2Node) error {
		isDel, err := r.IsDeleted(node.StartID)
		if err != nil {
			return fmt.Errorf("IsDeleted(%v): %w", node.StartID, err)
		}

		if isDel {
			return nil
		}

		mop = mop.Append(FormatOp{
			ID:      r.NextID(1),
			StartID: node.StartID,
			EndID:   node.EndID,
			Format:  FormatV3Line{},
		})

		return nil
	})

	if err != nil {
		return mop, fmt.Errorf("NOS.line.tree.Slice(%d, %d): %w", startIx, endIx, err)
	}

	_, err = r.MergeOp(mop)
	if err != nil {
		return mop, fmt.Errorf("MergeOp(%v): %w", mop, err)
	}

	return mop, nil
}

func (r *Rogue) _startIx(beforeID ID) (int, error) {
	startID, err := r.VisRightOf(beforeID)
	if err != nil {
		if errors.As(err, &ErrorNoRightVisSibling{}) {
			return r.VisSize, nil
		}
		return -1, fmt.Errorf("VisRightOf(%v): %w", beforeID, err)
	}

	startIx, _, err := r.Rope.GetIndex(startID)
	if err != nil {
		return -1, fmt.Errorf("Rope.GetIndex(%v): %w", startID, err)
	}

	return startIx, nil
}

// ApplyMarkdownDiff applies a diff to the current document
func (r *Rogue) ApplyMarkdownDiff(author, newMd string, beforeID, afterID ID) (MultiOp, Actions, error) {
	mop := MultiOp{}

	authorBefore := r.Author
	defer func() { r.Author = authorBefore }()
	r.Author = author

	cidx, err := r._startIx(beforeID)
	if err != nil {
		return mop, nil, fmt.Errorf("SelectionStartIx(%v): %w", beforeID, err)
	}
	offset := cidx

	startID, endID, err := r._getOriginalSelection(beforeID, afterID)
	if err != nil {
		return mop, nil, fmt.Errorf("GetOriginalSelection(%v, %v): %v", beforeID, afterID, err)
	}

	vis, err := r.Rope.GetBetween(startID, endID)
	if err != nil {
		return mop, nil, fmt.Errorf("GetBetween(%v, %v): %v", startID, endID, err)
	}

	plaintext := Uint16ToStr(vis.Text)
	if len(plaintext) > 0 && plaintext[len(plaintext)-1] == '\n' && len(newMd) > 0 && newMd[len(newMd)-1] != '\n' {
		newMd = newMd + "\n"
	}

	mdDiff, err := DiffMarkdown(plaintext, newMd)
	if err != nil {
		return mop, nil, fmt.Errorf("DiffMarkdown(%q, %q): %w", plaintext, newMd, err)
	}

	// First remove old formats
	fop, err := r.ClearFormats(startID, endID)
	if err != nil {
		return mop, nil, fmt.Errorf("ClearFormats(%v, %v): %w", startID, endID, err)
	}
	mop = mop.Append(fop)

	// Apply the text diffs
	actions := make(Actions, 0) // TODO: calculate and return actions consistently
	var op Op
	if plaintext == "" {
		op, err := r.InsertRightOf(beforeID, mdDiff.NewPlaintext)
		if err != nil {
			return mop, nil, fmt.Errorf("InsertRightOf(%v, %q): %w", beforeID, newMd, err)
		}
		mop = mop.Append(op)

		offset, _, err = r.Rope.GetIndex(op.ID)
		if err != nil {
			return mop, nil, fmt.Errorf("GetIndex(%v): %w", op.ID, err)
		}

		actions = append(actions, InsertAction{
			Index: offset,
			Text:  mdDiff.NewPlaintext,
		})
	} else {
		for _, d := range mdDiff.TextDiffs {
			switch d.Type {
			case diffmatchpatch.DiffDelete:
				delLen := UTF16Length(d.Text)
				dop, err := r.Delete(cidx, delLen)
				if err != nil {
					return mop, nil, fmt.Errorf("Delete(%v, %v): %w", cidx, delLen, err)
				}
				mop = mop.Append(dop)
				actions = append(actions, DeleteAction{
					Index: cidx,
					Count: delLen,
				})
			case diffmatchpatch.DiffInsert:
				// op, err = r._handleTailInsert(cidx, d.Text)
				op, err = r.Insert(cidx, d.Text)
				if err != nil {
					return mop, nil, fmt.Errorf("handleTailInsert(%v, %q): %w", cidx, d.Text, err)
				}
				mop = mop.Append(op)
				actions = append(actions, InsertAction{
					Index: cidx,
					Text:  d.Text,
				})
				cidx += UTF16Length(d.Text)
			case diffmatchpatch.DiffEqual:
				cidx += UTF16Length(d.Text)
			}
		}
	}

	lineOps, lineActions, err := r._applyFormatDiffs(offset, mdDiff.NewFormats)
	if err != nil {
		return mop, nil, fmt.Errorf("_applyFormatDiffs: %w", err)
	}

	mop = mop.Append(lineOps)
	actions = append(actions, lineActions...)
	r.OpIndex.Put(mop)

	return mop, actions, nil
}

func DiffLines(text1, text2 string) []diffmatchpatch.Diff {
	dmp := diffmatchpatch.New()
	lines1, lines2, lineArray := dmp.DiffLinesToChars(text1, text2)
	diffs := dmp.DiffMain(lines1, lines2, false)
	return dmp.DiffCharsToLines(diffs, lineArray)
}

func DiffWords(text1, text2 string) []diffmatchpatch.Diff {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(text1, text2, false)
	return MergeAdjacentInsertsDeletes(MergeSplitWords(SplitAtNewlines(diffs)))
}

var (
	runeSep       = ','
	runeSepOffset = int(runeSep) + 1
)

func SplitIncludingWhitespace(input string) []string {
	var result []string
	word := ""
	for _, char := range input {
		if char == ' ' {
			if word != "" {
				result = append(result, word)
				word = ""
			}
			if len(result) > 0 {
				result[len(result)-1] = result[len(result)-1] + " "
			} else {
				result = append(result, " ")
			}
		} else if char == '\t' {
			if word != "" {
				result = append(result, word)
				word = ""
			}
			result = append(result, "\t")
		} else if char == '\n' {
			if word != "" {
				result = append(result, word)
				word = ""
			}
			result = append(result, "\n")
		} else {
			word += string(char)
		}
	}
	if word != "" {
		result = append(result, word)
	}
	return result
}

func _splitKeepNewlines(s string) []string {
	var result []string
	var builder strings.Builder

	for _, r := range s {
		if r == '\n' {
			if builder.Len() > 0 {
				result = append(result, builder.String())
				builder.Reset()
			}
			result = append(result, string(r))
		} else {
			builder.WriteRune(r)
		}
	}

	if builder.Len() > 0 {
		result = append(result, builder.String())
	}

	return result
}

func SplitAtNewlines(diffs []diffmatchpatch.Diff) []diffmatchpatch.Diff {
	result := make([]diffmatchpatch.Diff, 0, len(diffs)*2)
	for _, diff := range diffs {
		for _, line := range _splitKeepNewlines(diff.Text) {
			result = append(result, diffmatchpatch.Diff{Type: diff.Type, Text: line})
		}
	}

	// fmt.Printf("SplitAtNewlines: %v\n", result)
	return result
}

func _splitFrontNonspace(s string) (begin, end string) {
	for i, r := range s {
		if _isSpace(r) {
			return s[:i], s[i:]
		}
	}
	return s, ""
}

func _splitBackNonspace(s string) (begin, end string) {
	for i := len(s) - 1; i >= 0; i-- {
		if _isSpace(rune(s[i])) {
			return s[:i+1], s[i+1:]
		}
	}

	return "", s
}

func _isSpace(r rune) bool {
	return unicode.IsSpace(r) && r != '\n'
}

func _isAllSpace(s string) bool {
	for _, r := range s {
		if !_isSpace(r) {
			return false
		}
	}

	return true
}

func _isAllPunctuation(s string) bool {
	for _, r := range s {
		if !unicode.IsPunct(r) {
			return false
		}
	}
	return true
}

func MergeSplitWords(diffs []diffmatchpatch.Diff) []diffmatchpatch.Diff {
	result := make([]diffmatchpatch.Diff, 0, len(diffs))

	i := 0
	for i < len(diffs) {
		diff := diffs[i]

		if diff.Type != diffmatchpatch.DiffEqual || diff.Text == "\n" || _isAllPunctuation(diff.Text) {
			i++
			result = append(result, diff)
			continue
		}

		runes := []rune(diff.Text)

		// Check start of equal diff
		if !_isSpace(runes[0]) && len(result) > 0 && !_isAllPunctuation(diff.Text) {
			prevDiff := result[len(result)-1]
			prevRunes := []rune(prevDiff.Text)

			isSplitWord := !_isSpace(prevRunes[len(prevRunes)-1]) && !_isAllPunctuation(prevDiff.Text) && prevDiff.Text != "\n"

			if isSplitWord && prevDiff.Type == diffmatchpatch.DiffInsert {
				beg, end := _splitFrontNonspace(diff.Text)

				if beg != "" && !_isAllPunctuation(beg) {
					result[len(result)-1].Text = prevDiff.Text + beg

					result = append(result, diffmatchpatch.Diff{
						Type: diffmatchpatch.DiffDelete,
						Text: beg,
					})

					diff.Text = end
				}
			} else if isSplitWord && prevDiff.Type == diffmatchpatch.DiffDelete {
				beg, end := _splitFrontNonspace(diff.Text)

				if beg != "" && !_isAllPunctuation(beg) {
					result[len(result)-1].Text = prevDiff.Text + beg

					result = append(result, diffmatchpatch.Diff{
						Type: diffmatchpatch.DiffInsert,
						Text: beg,
					})

					diff.Text = end
				}
			}
		}

		if diff.Text == "" {
			i++
			continue
		}

		// diff text may be different now
		runes = []rune(diff.Text)

		// Check end of equal diff
		if !_isSpace(runes[len(runes)-1]) && i < len(diffs)-1 && !_isAllPunctuation(diff.Text) {
			nextDiff := diffs[i+1]
			nextRunes := []rune(nextDiff.Text)
			isSplitWord := !_isSpace(nextRunes[0]) && !_isAllPunctuation(nextDiff.Text) && nextDiff.Text != "\n"

			if !isSplitWord || nextDiff.Type == diffmatchpatch.DiffEqual {
				i++
				result = append(result, diff)
				continue
			}

			if nextDiff.Type == diffmatchpatch.DiffInsert {
				beg, end := _splitBackNonspace(diff.Text)

				if end != "" && !_isAllPunctuation(end) {
					if beg != "" {
						result = append(result, diffmatchpatch.Diff{
							Type: diffmatchpatch.DiffEqual,
							Text: beg,
						})
					}

					result = append(result, diffmatchpatch.Diff{
						Type: diffmatchpatch.DiffDelete,
						Text: end,
					})

					result = append(result, diffmatchpatch.Diff{
						Type: diffmatchpatch.DiffInsert,
						Text: end + nextDiff.Text,
					})

					i++
				} else {
					result = append(result, diff)
				}
			} else if isSplitWord && nextDiff.Type == diffmatchpatch.DiffDelete {
				beg, end := _splitBackNonspace(diff.Text)

				if end != "" && !_isAllPunctuation(end) {
					if beg != "" {
						result = append(result, diffmatchpatch.Diff{
							Type: diffmatchpatch.DiffEqual,
							Text: beg,
						})
					}

					result = append(result, diffmatchpatch.Diff{
						Type: diffmatchpatch.DiffInsert,
						Text: end,
					})

					result = append(result, diffmatchpatch.Diff{
						Type: diffmatchpatch.DiffDelete,
						Text: end + nextDiff.Text,
					})

					i++
				} else {
					result = append(result, diff)
				}
			}
		} else if diff.Text != "" {
			result = append(result, diff)
		}

		i++
	}

	// fmt.Printf("MergeSplitWords: %v\n", result)
	return result
}

func MergeAdjacentInsertsDeletes(diffs []diffmatchpatch.Diff) []diffmatchpatch.Diff {
	result := make([]diffmatchpatch.Diff, 0, len(diffs))
	lastInsertIdx := -1
	lastDeleteIdx := -1

	for i, diff := range diffs {
		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			if lastInsertIdx == -1 {
				lastInsertIdx = len(result)
				result = append(result, diff)
			} else {
				result[lastInsertIdx].Text += diff.Text
			}
		case diffmatchpatch.DiffDelete:
			if lastDeleteIdx == -1 {
				lastDeleteIdx = len(result)
				result = append(result, diff)
			} else {
				result[lastDeleteIdx].Text += diff.Text
			}
		case diffmatchpatch.DiffEqual:
			var nextDiff *diffmatchpatch.Diff
			if i < len(diffs)-1 {
				nextDiff = &diffs[i+1]
			}
			nextNotEqual := nextDiff != nil && nextDiff.Type != diffmatchpatch.DiffEqual

			if _isAllSpace(diff.Text) && nextNotEqual && (lastInsertIdx != -1 || lastDeleteIdx != -1) {
				if lastInsertIdx == -1 {
					lastInsertIdx = len(result)
					result = append(result, diffmatchpatch.Diff{
						Type: diffmatchpatch.DiffInsert,
						Text: diff.Text,
					})
				} else {
					result[lastInsertIdx].Text += diff.Text
				}

				if lastDeleteIdx == -1 {
					lastDeleteIdx = len(result)
					result = append(result, diffmatchpatch.Diff{
						Type: diffmatchpatch.DiffDelete,
						Text: diff.Text,
					})
				} else {
					result[lastDeleteIdx].Text += diff.Text
				}
			} else {
				result = append(result, diff)
				lastInsertIdx = -1
				lastDeleteIdx = -1
			}
		default:
			log.Errorf("invalid diff type: %v", diff.Type)
		}
	}

	// fmt.Printf("MergeAdjacentInsertsDeletes: %v\n", result)
	return result
}
