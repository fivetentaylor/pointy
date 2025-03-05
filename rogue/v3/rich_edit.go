package v3

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/teamreviso/code/pkg/stackerr"
)

var (
	_bulletRe   *regexp.Regexp
	_orderedRe  *regexp.Regexp
	_headerRe   *regexp.Regexp
	_quoteRe    *regexp.Regexp
	_ruleRe     *regexp.Regexp
	_strikeRe   *regexp.Regexp
	_boldRe     *regexp.Regexp
	_italicRe   *regexp.Regexp
	_linkRe     *regexp.Regexp
	_codeSnipRe *regexp.Regexp
)

func init() {
	_bulletPattern := `^ *(-|\*) $`
	_bulletRe = regexp.MustCompile(_bulletPattern)

	_orderedPattern := `^ *(\d+)\. $`
	_orderedRe = regexp.MustCompile(_orderedPattern)

	_headerPattern := `^#{1,6} $`
	_headerRe = regexp.MustCompile(_headerPattern)

	_quotePattern := `^> $`
	_quoteRe = regexp.MustCompile(_quotePattern)

	_rulePattern := `^--$`
	_ruleRe = regexp.MustCompile(_rulePattern)

	_strikePattern := `~(\S(?:[^\s~]| [^\s~])+?\S)~`
	_strikeRe = regexp.MustCompile(_strikePattern)

	_boldPattern := `\*\*(\S(?:[^\s*]| [^\s*])+?\S)\*\*`
	_boldRe = regexp.MustCompile(_boldPattern)

	_italicPattern := `\*([^*\s](?:[^*]+[^*\s])?)\*`
	_italicRe = regexp.MustCompile(_italicPattern)

	_linkPattern := `\[(?:[^[\]\\]|\\.)+\]\((?:[^()\\]|\\.)+\)`
	_linkRe = regexp.MustCompile(_linkPattern)

	_codeSnipPattern := "`(\\S(?:[^\\s`]| [^\\s`])+?\\S)`"
	_codeSnipRe = regexp.MustCompile(_codeSnipPattern)
}

func (r *Rogue) _getVisLine(visIx int) (*FugueVis, error) {
	c, err := r.GetChar(visIx)
	if err != nil {
		return nil, fmt.Errorf("GetChar(%d): %w", visIx, err)
	}

	if visIx == 0 && c == '\n' {
		return &FugueVis{IDs: []ID{}, Text: []uint16{}}, nil
	}

	lix, rix := visIx, visIx
	if visIx > 0 && c == '\n' {
		lix = visIx - 1
	}

	lid, err := r.Rope.GetVisID(lix)
	if err != nil {
		return nil, fmt.Errorf("Rope.GetVisID(%d): %w", lix, err)
	}

	rid := lid
	if lix == rix {
		rid, err = r.Rope.GetVisID(visIx)
		if err != nil {
			return nil, fmt.Errorf("Rope.GetVisID(%d): %w", visIx, err)
		}
	}

	startID, err := r.VisScanLeftOf(lid, '\n')
	if err != nil {
		return nil, fmt.Errorf("VisScanLeftOf(%v, '\\n'): %w", lid, err)
	}

	if startID == nil {
		sid, err := r.Rope.GetVisID(0)
		if err != nil {
			return nil, fmt.Errorf("Rope.GetVisID(0): %w", err)
		}

		startID = &sid
	} else {
		sid, err := r.VisRightOf(*startID)
		if err != nil {
			return nil, fmt.Errorf("VisRightOf(%v): %w", *startID, err)
		}

		startID = &sid
	}

	endID, err := r.VisScanRightOf(rid, '\n')
	if err != nil {
		return nil, fmt.Errorf("VisScanRightOf(%v, '\\n'): %w", rid, err)
	}

	if endID == nil {
		eid, err := r.Rope.GetVisID(r.VisSize - 1)
		if err != nil {
			return nil, fmt.Errorf("Rope.GetVisID(%d): %w", r.VisSize-1, err)
		}

		endID = &eid
	} else {
		eid, err := r.VisLeftOf(*endID)
		if err != nil {
			return nil, fmt.Errorf("VisLeftOf(%v): %w", *endID, err)
		}

		endID = &eid
	}

	vis, err := r.Rope.GetBetween(*startID, *endID)
	if err != nil {
		return nil, fmt.Errorf("Rope.GetBetween(%v, %v): %w", *startID, *endID, err)
	}

	return vis, nil
}

func (r *Rogue) _getStartVisLine(visIx int) (*FugueVis, error) {
	c, err := r.GetChar(visIx)
	if err != nil {
		return nil, fmt.Errorf("GetChar(%d): %w", visIx, err)
	}

	if visIx == 0 && c == '\n' {
		return &FugueVis{IDs: []ID{}, Text: []uint16{}}, nil
	}

	endID, err := r.Rope.GetVisID(visIx)
	if err != nil {
		return nil, fmt.Errorf("Rope.GetVisID(%d): %w", visIx, err)
	}

	lix := visIx
	if visIx > 0 && c == '\n' {
		lix = visIx - 1
	}

	lid, err := r.Rope.GetVisID(lix)
	if err != nil {
		return nil, fmt.Errorf("Rope.GetVisID(%d): %w", lix, err)
	}

	startID, err := r.VisScanLeftOf(lid, '\n')
	if err != nil {
		return nil, fmt.Errorf("VisScanLeftOf(%v, '\\n'): %w", lid, err)
	}

	if startID == nil {
		lid, err = r.Rope.GetVisID(0)
		if err != nil {
			return nil, fmt.Errorf("Rope.GetVisID(0): %w", err)
		}
	} else {
		lid, err = r.VisRightOf(*startID)
		if err != nil {
			return nil, fmt.Errorf("VisRightOf(%v): %w", *startID, err)
		}
	}

	startID = &lid

	vis, err := r.Rope.GetBetween(*startID, endID)
	if err != nil {
		return nil, fmt.Errorf("Rope.GetBetween(%v, %v): %w", *startID, endID, err)
	}

	return vis, nil
}

func (r *Rogue) RichInsert(visIx, selLen int, orgSpanFormat FormatV3Span, sText string) (ops []Op, cursorID ID, err error) {
	mop := MultiOp{}

	orgSpanFormat["e"] = "true"
	orgSpanFormat["en"] = "true"

	// HANDLE INSERT LINK
	if selLen > 0 && IsValidURL(sText) {
		fop, err := r.Format(visIx, selLen, FormatV3Span{"a": sText})
		if err != nil {
			return nil, cursorID, err
		}

		cursorID, err = r.Rope.GetVisID(visIx + selLen)
		if err != nil {
			return nil, cursorID, err
		}

		ops = append(ops, fop)
		return ops, cursorID, nil
	}

	selStartID, err := r.Rope.GetVisID(visIx)
	if err != nil {
		return nil, NoID, err
	}

	selEndID, err := r.Rope.GetVisID(visIx + selLen)
	if err != nil {
		return nil, NoID, err
	}

	// HANDLE TABS OF LISTS WITH SELECTION
	if selLen > 0 && (sText == "\t" || sText == "(1+4cT5lP9") {
		isSpanAllLists, err := r.IsSpanAllLists(selStartID, selEndID)
		if err != nil {
			return nil, NoID, err
		}

		if isSpanAllLists {
			prevLineFormat, err := r.GetPrevLineFormat(selStartID)
			if err != nil {
				return nil, NoID, err
			}
			prev := getListMeta(prevLineFormat)

			curLineFormat, err := r.GetCurLineFormat(selStartID, selStartID)
			if err != nil {
				return nil, NoID, err
			}
			cur := getListMeta(curLineFormat)

			delta := -1
			if sText == "\t" {
				delta = 1
			}

			canIndent := delta == 1 && prev.isList && cur.indent <= prev.indent
			canDedent := delta == -1 && prev.isList && cur.indent >= prev.indent

			if canIndent || canDedent {
				mop, err := r.IndentSpan(selStartID, selEndID, delta)
				if err != nil {
					return nil, NoID, err
				}

				ops = append(ops, mop)
				return ops, selEndID, nil
			}

			return nil, selEndID, nil
		}
	}

	// DELETE SELECTION IF ANY
	if selLen > 0 {
		dop, visIx, err := r.RichDelete(visIx, selLen)
		if err != nil {
			return nil, NoID, fmt.Errorf("RichDelete(%d, %d): %w", visIx, selLen, err)
		}

		selStartID, err = r.Rope.GetVisID(visIx)
		if err != nil {
			return nil, NoID, fmt.Errorf("Rope.GetVisID(%d): %w", visIx, err)
		}

		mop = mop.Append(dop)
	}

	isStartOfLine := false
	if visIx == 0 {
		isStartOfLine = true
	} else {
		leftChar, err := r.GetChar(visIx - 1)
		if err != nil {
			return nil, NoID, fmt.Errorf("GetChar(%d): %w", visIx-1, err)
		}

		if leftChar == '\n' {
			isStartOfLine = true
		}
	}

	isEndOfLine := false
	if visIx == r.VisSize-1 {
		isEndOfLine = true
	} else {
		rightChar, err := r.GetChar(visIx)
		if err != nil {
			return nil, NoID, fmt.Errorf("GetChar(%d): %w", visIx, err)
		}

		if rightChar == '\n' {
			isEndOfLine = true
		}
	}

	// HANDLE SHIFT ENTER FOR SOFT RETURN
	if sText == "oLEcI0yPY9" {
		curFormat, err := r.GetCurLineFormat(selStartID, selStartID)
		if err != nil {
			return nil, NoID, fmt.Errorf("GetCurLineFormat(%v, %v): %w", selStartID, selStartID, err)
		}

		cur := getListMeta(curFormat)
		if cur.isList {
			op, err := r.Insert(visIx, "\n")
			if err != nil {
				return nil, NoID, fmt.Errorf("Insert(%d, %q): %w", visIx, "\n", err)
			}
			mop = mop.Append(op)

			fop, err := r.Format(visIx, 1, curFormat)
			if err != nil {
				return nil, NoID, fmt.Errorf("Format(%d, 1, %v): %w", visIx, curFormat, err)
			}
			mop = mop.Append(fop)

			fop, err = r.Format(visIx+1, 1, FormatV3IndentedLine(cur.indent))
			if err != nil {
				return nil, NoID, fmt.Errorf("Format(%d, 1, %v): %w", visIx+1, FormatV3IndentedLine(cur.indent), err)
			}
			mop = mop.Append(fop)

			r.OpIndex.Put(mop)
			ops = append(ops, mop)

			cursorID, err = r.Rope.GetVisID(visIx + 1)
			if err != nil {
				return nil, cursorID, fmt.Errorf("Rope.GetVisID(%d): %w", visIx+1, err)
			}

			return ops, cursorID, nil
		}

		sText = "\n" // fallthrough to regular newline
	}

	// HANDLE TAB AND SHIFT TAB FOR LIST INDENT
	if selLen == 0 && (sText == "\t" || sText == "(1+4cT5lP9" || sText == "\n") {
		// "(1+4cT5lP9" is a random string to represent a Shift + Tab keypress
		if isStartOfLine {
			curFormat, err := r.GetCurLineFormat(selStartID, selStartID)
			if err != nil {
				return nil, NoID, err
			}
			cur := getListMeta(curFormat)

			cond0 := sText == "\n" && cur.isList && cur.indent > 0
			cond1 := sText != "\n" && cur.isList

			if cond0 || cond1 {
				prevFormat, err := r.GetPrevLineFormat(selStartID)
				if err != nil {
					return nil, NoID, err
				}
				prev := getListMeta(prevFormat)

				indent := 0

				if prev.isList {
					if sText == "\t" {
						if cur.indent < prev.indent {
							indent = prev.indent
						} else {
							indent = prev.indent + 1
						}
					} else {
						if prev.indent < cur.indent {
							indent = prev.indent
						} else {
							indent = max(0, prev.indent-1)
						}
					}
				}

				var format FormatV3
				switch curFormat.(type) {
				case FormatV3BulletList:
					format = FormatV3BulletList(indent)
				case FormatV3OrderedList:
					format = FormatV3OrderedList(indent)
				case FormatV3IndentedLine:
					format = FormatV3IndentedLine(indent)
				}

				fops, err := r.Format(visIx, 1, format)
				if err != nil {
					return nil, NoID, err
				}

				mop = mop.Append(fops)
				op := FlattenMop(mop)
				r.OpIndex.Put(op)
				ops = append(ops, op)
				return ops, selStartID, nil
			}
		}
	}

	// this represents a Shift+Tab keypress, just ignore it if we got to this point
	if sText == "(1+4cT5lP9" {
		return nil, selStartID, nil
	}

	// HANDLE MARKDOWN SHORTCUTS
	prevChar := uint16(' ')
	if visIx > 0 {
		prevChar, err = r.GetChar(visIx - 1)
		if err != nil {
			return nil, NoID, fmt.Errorf("GetChar(%d): %w", visIx-1, err)
		}
	}

	if sText == " " && selLen == 0 {
		iop, err := r.Insert(visIx, sText)
		if err != nil {
			return nil, NoID, fmt.Errorf("Insert(%d, %q): %w", visIx, sText, err)
		}
		ops = append(ops, iop)

		vis, err := r._getStartVisLine(visIx)
		if err != nil {
			return nil, NoID, fmt.Errorf("_getStartVisLine(%d): %w", visIx, err)
		}

		lineStartIx, _, err := r.Rope.GetIndex(vis.IDs[0])
		if err != nil {
			return nil, NoID, fmt.Errorf("GetIndex(%v): %w", vis.IDs[0], err)
		}

		lineText := Uint16ToStr(vis.Text[0 : visIx-lineStartIx+1])

		// BULLET LIST SHORTCUT
		if _bulletRe.MatchString(lineText) {
			curFormat, err := r.GetCurLineFormat(selStartID, selStartID)
			if err != nil {
				return nil, NoID, err
			}

			cur := getListMeta(curFormat)

			if !cur.isBullet {
				six, _, err := r.Rope.GetIndex(vis.IDs[0])
				if err != nil {
					return nil, NoID, err
				}

				mop = MultiOp{}
				dop, err := r.Delete(six, visIx-six+1)
				if err != nil {
					return nil, NoID, err
				}
				mop = mop.Append(dop)

				fop, err := r.Format(six, 1, FormatV3BulletList(cur.indent))
				if err != nil {
					return nil, NoID, err
				}
				mop = mop.Append(fop)

				cursorID, err = r.Rope.GetVisID(six)
				if err != nil {
					return nil, cursorID, fmt.Errorf("Rope.GetVisID(%d): %w", six, err)
				}

				r.OpIndex.Put(mop)
				ops = append(ops, mop)

				return ops, cursorID, nil
			}
		}

		// ORDERED LIST SHORTCUT
		if _orderedRe.MatchString(lineText) {
			curFormat, err := r.GetCurLineFormat(selStartID, selStartID)
			if err != nil {
				return nil, NoID, fmt.Errorf("GetCurLineFormat(%v, %v): %w", selStartID, selStartID, err)
			}

			cur := getListMeta(curFormat)

			if !cur.isOrdered {
				six, _, err := r.Rope.GetIndex(vis.IDs[0])
				if err != nil {
					return nil, NoID, fmt.Errorf("GetIndex(%v): %w", vis.IDs[0], err)
				}

				// reset multiop to make undo history cleaner
				mop = MultiOp{}
				dop, err := r.Delete(six, visIx-six+1)
				if err != nil {
					return nil, NoID, fmt.Errorf("Delete(%d, %d): %w", six, visIx-six, err)
				}
				mop = mop.Append(dop)

				fop, err := r.Format(six, 1, FormatV3OrderedList(cur.indent))
				if err != nil {
					return nil, NoID, fmt.Errorf("Format(%d, 1, %v): %w", six, FormatV3OrderedList(0), err)
				}
				mop = mop.Append(fop)

				cursorID, err = r.Rope.GetVisID(six)
				if err != nil {
					return nil, cursorID, fmt.Errorf("Rope.GetVisID(%d): %w", six, err)
				}

				r.OpIndex.Put(mop)
				ops = append(ops, mop)

				return ops, cursorID, nil
			}
		}

		// HEADER SHORTCUT
		if _headerRe.MatchString(lineText) {
			six, _, err := r.Rope.GetIndex(vis.IDs[0])
			if err != nil {
				return nil, NoID, fmt.Errorf("GetIndex(%v): %w", vis.IDs[0], err)
			}

			// reset multiop to make undo history cleaner
			mop = MultiOp{}
			dop, err := r.Delete(six, visIx-six+1)
			if err != nil {
				return nil, NoID, fmt.Errorf("Delete(%d, %d): %w", six, visIx-six, err)
			}
			mop = mop.Append(dop)

			headerIndent := CountLeadingChar(lineText, '#')
			fop, err := r.Format(six, 1, FormatV3Header(headerIndent))
			if err != nil {
				return nil, NoID, fmt.Errorf("Format(%d, 1, %v): %w", six, FormatV3Header(0), err)
			}
			mop = mop.Append(fop)

			cursorID, err = r.Rope.GetVisID(six)
			if err != nil {
				return nil, cursorID, fmt.Errorf("Rope.GetVisID(%d): %w", six, err)
			}

			r.OpIndex.Put(mop)
			ops = append(ops, mop)
			return ops, cursorID, nil
		}

		// QUOTE SHORTCUT
		if _quoteRe.MatchString(lineText) {
			six, _, err := r.Rope.GetIndex(vis.IDs[0])
			if err != nil {
				return nil, NoID, fmt.Errorf("GetIndex(%v): %w", vis.IDs[0], err)
			}

			// reset multiop to make undo history cleaner
			mop = MultiOp{}
			dop, err := r.Delete(six, visIx-six+1)
			if err != nil {
				return nil, NoID, fmt.Errorf("Delete(%d, %d): %w", six, visIx-six, err)
			}
			mop = mop.Append(dop)

			fop, err := r.Format(six, 1, FormatV3BlockQuote{})
			if err != nil {
				return nil, NoID, fmt.Errorf("Format(%d, 1, %v): %w", six, FormatV3BlockQuote{}, err)
			}
			mop = mop.Append(fop)

			cursorID, err = r.Rope.GetVisID(six)
			if err != nil {
				return nil, cursorID, fmt.Errorf("Rope.GetVisID(%d): %w", six, err)
			}

			r.OpIndex.Put(mop)
			ops = append(ops, mop)
			return ops, cursorID, nil
		}

		// CODE BLOCK SHORTCUT
		language := _isCodeBlock(lineText)
		if language != nil {
			six, _, err := r.Rope.GetIndex(vis.IDs[0])
			if err != nil {
				return nil, NoID, fmt.Errorf("GetIndex(%v): %w", vis.IDs[0], err)
			}

			// reset multiop to make undo history cleaner
			mop = MultiOp{}
			dop, err := r.Delete(six, visIx-six+1)
			if err != nil {
				return nil, NoID, fmt.Errorf("Delete(%d, %d): %w", six, visIx-six, err)
			}
			mop = mop.Append(dop)

			fop, err := r.Format(six, 1, FormatV3CodeBlock(*language))
			if err != nil {
				return nil, NoID, fmt.Errorf("Format(%d, 1, %v): %w", six, FormatV3BlockQuote{}, err)
			}
			mop = mop.Append(fop)

			cursorID, err = r.Rope.GetVisID(six)
			if err != nil {
				return nil, cursorID, fmt.Errorf("Rope.GetVisID(%d): %w", six, err)
			}

			r.OpIndex.Put(mop)
			ops = append(ops, mop)
			return ops, cursorID, nil
		}

		curSpanFormat, err := r.GetCurSpanFormat(selStartID, selStartID)
		if err != nil {
			return nil, NoID, fmt.Errorf("GetCurSpanFormat(%v, %v): %w", selStartID, selStartID, err)
		}

		if !curSpanFormat.Equals(orgSpanFormat) {
			fop, err := r.Format(visIx, 1, orgSpanFormat)
			if err != nil {
				return nil, NoID, fmt.Errorf("Format(%d, 1, %v): %w", visIx, orgSpanFormat, err)
			}

			ops = append(ops, fop)
		}

		cursorID, err = r.Rope.GetVisID(visIx + 1)
		if err != nil {
			return nil, cursorID, fmt.Errorf("Rope.GetVisID(%d): %w", visIx, err)
		}

		return ops, cursorID, nil
	}

	// MARKDOWN SPAN FORMATS
	if selLen == 0 && (sText == "~" || sText == "*" || sText == "`") {
		// Insert first then user can undo back to original text
		op, err := r.Insert(visIx, sText)
		if err != nil {
			return nil, NoID, fmt.Errorf("Insert(%d, %q): %w", visIx, sText, err)
		}
		ops = append(ops, op)
		mop = MultiOp{}

		cursorID, err = r.Rope.GetVisID(visIx + 1)
		if err != nil {
			return nil, cursorID, fmt.Errorf("Rope.GetVisID(%d): %w", visIx+1, err)
		}

		vis, err := r._getVisLine(visIx)
		if err != nil {
			return nil, NoID, fmt.Errorf("_getVisLine(%d): %w", visIx, err)
		}

		lineText := Uint16ToStr(vis.Text)
		lineStartIx, _, err := r.Rope.GetIndex(vis.IDs[0])
		if err != nil {
			return nil, NoID, fmt.Errorf("GetIndex(%v): %w", vis.IDs[0], err)
		}

		if sText == "~" {
			// HANDLE STRIKETHROUGH SHORTCUT ~ ~
			startIx, endIx := _spanMatchIndices(_strikeRe, lineText, visIx, lineStartIx)

			if startIx >= 0 {
				dop, err := r.Delete(endIx, 1)
				if err != nil {
					return nil, NoID, fmt.Errorf("Delete(%d, 1): %w", endIx, err)
				}
				mop = mop.Append(dop)

				fop, err := r.Format(startIx+1, endIx-startIx-1, FormatV3Span{"s": "true"})
				if err != nil {
					return nil, NoID, fmt.Errorf("Format(%d, %d, %v): %w", startIx+1, endIx-startIx-1, FormatV3Span{"s": "true"}, err)
				}
				mop = mop.Append(fop)

				dop, err = r.Delete(startIx, 1)
				if err != nil {
					return nil, NoID, fmt.Errorf("Delete(%d, 1): %w", startIx, err)
				}
				mop = mop.Append(dop)

				r.OpIndex.Put(mop)
				ops = append(ops, mop)

				cursorID, err = r.Rope.GetVisID(endIx - 1)
				if err != nil {
					return nil, cursorID, fmt.Errorf("Rope.GetVisID(%d): %w", endIx, err)
				}
			}
		} else if sText == "*" {
			// HANDLE BOLD SHORTCUT ** **
			startIx, endIx := _spanMatchIndices(_boldRe, lineText, visIx, lineStartIx)

			if startIx >= 0 {
				dop, err := r.Delete(endIx-1, 2)
				if err != nil {
					return nil, NoID, fmt.Errorf("Delete(%d, 2): %w", endIx-1, err)
				}
				mop = mop.Append(dop)

				fop, err := r.Format(startIx+2, endIx-startIx-3, FormatV3Span{"b": "true"})
				if err != nil {
					return nil, NoID, fmt.Errorf("Format(%d, %d, %v): %w", startIx+2, endIx-startIx-2, FormatV3Span{"b": "true"}, err)
				}
				mop = mop.Append(fop)

				dop, err = r.Delete(startIx, 2)
				if err != nil {
					return nil, NoID, fmt.Errorf("Delete(%d, 2): %w", startIx, err)
				}
				mop = mop.Append(dop)

				r.OpIndex.Put(mop)
				ops = append(ops, mop)

				cursorID, err = r.Rope.GetVisID(endIx - 3)
				if err != nil {
					return nil, cursorID, fmt.Errorf("Rope.GetVisID(%d): %w", endIx-3, err)
				}

				return ops, cursorID, nil
			}

			// HANDLE ITALIC SHORTCUT * *
			allIndices := _italicRe.FindAllStringIndex(lineText, -1)

			startIx, endIx = -1, -1
			for _, indices := range allIndices {
				start, end := indices[0], indices[1]
				// make sure this isn't a bold shortcut
				if start > 0 && lineText[start-1] == '*' {
					break
				}

				start = Utf8ToUtf16Ix(lineText, start)
				end = Utf8ToUtf16Ix(lineText, end-1)
				if end+lineStartIx == visIx {
					startIx = start + lineStartIx
					endIx = end + lineStartIx
					break
				}
			}

			if startIx >= 0 {
				dop, err := r.Delete(endIx, 1)
				if err != nil {
					return nil, NoID, fmt.Errorf("Delete(%d, 1): %w", endIx, err)
				}
				mop = mop.Append(dop)

				fop, err := r.Format(startIx+1, endIx-startIx-1, FormatV3Span{"i": "true"})
				if err != nil {
					return nil, NoID, fmt.Errorf("Format(%d, %d, %v): %w", startIx+1, endIx-startIx-1, FormatV3Span{"i": "true"}, err)
				}
				mop = mop.Append(fop)

				dop, err = r.Delete(startIx, 1)
				if err != nil {
					return nil, NoID, fmt.Errorf("Delete(%d, 1): %w", startIx, err)
				}
				mop = mop.Append(dop)

				r.OpIndex.Put(mop)
				ops = append(ops, mop)

				cursorID, err = r.Rope.GetVisID(endIx - 1)
				if err != nil {
					return nil, cursorID, fmt.Errorf("Rope.GetVisID(%d): %w", endIx-1, err)
				}

				return ops, cursorID, nil
			}
		} else if sText == "`" {
			// HANDLE CODE SHORTCUT `` ``
			startIx, endIx := _spanMatchIndices(_codeSnipRe, lineText, visIx, lineStartIx)

			if startIx >= 0 {
				dop, err := r.Delete(endIx, 1)
				if err != nil {
					return nil, NoID, fmt.Errorf("Delete(%d, 1): %w", endIx, err)
				}
				mop = mop.Append(dop)

				fop, err := r.Format(startIx+1, endIx-startIx-1, FormatV3Span{"c": "true"})
				if err != nil {
					return nil, NoID, fmt.Errorf("Format(%d, %d, %v): %w", startIx+1, endIx-startIx-1, FormatV3Span{"c": "true"}, err)
				}
				mop = mop.Append(fop)

				dop, err = r.Delete(startIx, 1)
				if err != nil {
					return nil, NoID, fmt.Errorf("Delete(%d, 1): %w", startIx, err)
				}
				mop = mop.Append(dop)

				r.OpIndex.Put(mop)
				ops = append(ops, mop)

				cursorID, err = r.Rope.GetVisID(endIx - 1)
				if err != nil {
					return nil, cursorID, fmt.Errorf("Rope.GetVisID(%d): %w", endIx-1, err)
				}

				return ops, cursorID, nil
			}
		}

		return ops, cursorID, nil
	}

	// HANDLE HORIZONTAL RULE SHORTCUT
	if selLen == 0 && sText == "-" && prevChar == '-' {
		vis, err := r._getVisLine(visIx)
		if err != nil {
			return nil, NoID, fmt.Errorf("_getVisLine(%d): %w", visIx, err)
		}

		lineText := Uint16ToStr(vis.Text)

		if _ruleRe.MatchString(lineText) {
			six, _, err := r.Rope.GetIndex(vis.IDs[0])
			if err != nil {
				return nil, NoID, fmt.Errorf("GetIndex(%v): %w", vis.IDs[0], err)
			}

			op, err := r.Insert(visIx, sText)
			if err != nil {
				return nil, NoID, fmt.Errorf("Insert(%d, %q): %w", visIx, sText, err)
			}
			ops = append(ops, op)

			// reset multiop to make undo history cleaner
			mop = MultiOp{}
			dop, err := r.Delete(six, visIx-six+1)
			if err != nil {
				return nil, NoID, fmt.Errorf("Delete(%d, %d): %w", six, visIx-six, err)
			}
			mop = mop.Append(dop)

			fop, err := r.Format(six, 1, FormatV3Rule{})
			if err != nil {
				return nil, NoID, fmt.Errorf("Format(%d, 1, %v): %w", six, FormatV3Rule{}, err)
			}
			mop = mop.Append(fop)

			iop, err := r.Insert(six+1, "\n")
			if err != nil {
				return nil, NoID, fmt.Errorf("Insert(%d, %q): %w", six+1, "\n", err)
			}
			mop = mop.Append(iop)

			cursorID, err = r.Rope.GetVisID(six + 1)
			if err != nil {
				return nil, cursorID, fmt.Errorf("Rope.GetVisID(%d): %w", six+1, err)
			}

			r.OpIndex.Put(mop)
			ops = append(ops, mop)

			return ops, cursorID, nil
		}
	}

	// HANDLE LINK SHORTCUT
	if selLen == 0 && sText == ")" {
		op, err := r.Insert(visIx, sText)
		if err != nil {
			return nil, NoID, fmt.Errorf("Insert(%d, %q): %w", visIx, sText, err)
		}
		ops = append(ops, op)

		vis, err := r._getVisLine(visIx)
		if err != nil {
			return nil, NoID, fmt.Errorf("_getVisLine(%d): %w", visIx, err)
		}

		lineText := Uint16ToStr(vis.Text)

		if _linkRe.MatchString(lineText) {
			ix0, ix1, ix2, url, err := _parseMarkdownLink(lineText)
			if err != nil {
				return nil, NoID, fmt.Errorf("_parseMarkdownLink(%q): %w", lineText, err)
			}

			lineStartIx, _, err := r.Rope.GetIndex(vis.IDs[0])
			if err != nil {
				return nil, NoID, fmt.Errorf("GetIndex(%v): %w", vis.IDs[0], err)
			}

			mop := MultiOp{}

			fop, err := r.Format(lineStartIx+ix0+1, ix1-ix0-1, FormatV3Span{"a": url})
			if err != nil {
				return nil, NoID, fmt.Errorf("Format(%d, %d, %v): %w", lineStartIx+ix0+1, ix1-ix0-1, FormatV3Span{"a": url}, err)
			}
			mop = mop.Append(fop)

			dop, err := r.Delete(lineStartIx+ix1, ix2-ix1+1)
			if err != nil {
				return nil, NoID, fmt.Errorf("Delete(%d, %d): %w", lineStartIx+ix1, ix2-ix1+1, err)
			}
			mop = mop.Append(dop)

			dop, err = r.Delete(lineStartIx+ix0, 1)
			if err != nil {
				return nil, NoID, fmt.Errorf("Delete(%d, 1): %w", lineStartIx+ix0, err)
			}
			mop = mop.Append(dop)

			cursorID, err = r.Rope.GetVisID(lineStartIx + ix1 - 1)
			if err != nil {
				return nil, cursorID, fmt.Errorf("Rope.GetVisID(%d): %w", lineStartIx+ix1-1, err)
			}

			r.OpIndex.Put(mop)
			ops = append(ops, mop)
			return ops, cursorID, nil
		}

		cursorID, err = r.Rope.GetVisID(visIx + 1)
		if err != nil {
			return nil, cursorID, fmt.Errorf("Rope.GetVisID(%d): %w", visIx+1, err)
		}

		return ops, cursorID, nil
	}

	text := StrToUint16(sText)

	// CAN JUST INSERT TEXT IF IT DOESN'T CONTAIN A NEWLINE
	if !strings.Contains(sText, "\n") {
		iop, err := r.Insert(visIx, sText)
		if err != nil {
			return nil, NoID, fmt.Errorf("Insert(%d, %q): %w", visIx, sText, err)
		}
		mop = mop.Append(iop)

		cursorID, err = r.Rope.GetVisID(visIx + len(text))
		if err != nil {
			return nil, cursorID, fmt.Errorf("Rope.GetVisID(%d): %w", visIx, err)
		}

		visID, err := r.Rope.GetVisID(visIx)
		if err != nil {
			return nil, NoID, fmt.Errorf("Rope.GetVisID(%d): %w", visIx, err)
		}

		curSpanFormat, err := r.GetCurSpanFormat(visID, visID)
		if err != nil {
			return nil, NoID, fmt.Errorf("GetCurSpanFormat(%v, %v): %w", visID, visID, err)
		}

		if !curSpanFormat.Equals(orgSpanFormat) {
			fop, err := r.Format(visIx, len(text), orgSpanFormat)
			if err != nil {
				return nil, NoID, fmt.Errorf("Format(%d, 1, %v): %w", visIx, orgSpanFormat, err)
			}

			mop = mop.Append(fop)
		}

		op := FlattenMop(mop)
		r.OpIndex.Put(op)
		ops = append(ops, op)
		return ops, cursorID, nil
	}

	// HANDLE INSERTING TEXT THAT CONTAINS NEWLINES
	format, err := r.GetCurLineFormat(selStartID, selStartID)
	if err != nil {
		return nil, NoID, fmt.Errorf("GetCurLineFormat(%v, %v): %w", selStartID, selStartID, err)
	}

	isEmpty, err := r.isEmptyLine(visIx)
	if err != nil {
		return nil, NoID, fmt.Errorf("isEmptyLine(%d): %w", visIx, err)
	}

	_, isImage := format.(FormatV3Image)
	if isImage && isEmpty && sText == "\n" {
		iop, err := r.Insert(visIx+1, sText)
		if err != nil {
			return nil, NoID, err
		}
		mop = mop.Append(iop)

		cursorID, err = r.Rope.GetVisID(visIx + 1)
		if err != nil {
			return nil, cursorID, err
		}

		r.OpIndex.Put(mop)
		ops = append(ops, mop)
		return ops, cursorID, nil
	}

	_, isPlainLine := format.(FormatV3Line)
	if !isPlainLine && !isImage && isEmpty && sText == "\n" {
		// SINGLE NEWLINE IN EMPTY LINE WITH FORMAT
		fop, err := r.Format(visIx, 1, FormatV3Line{})
		if err != nil {
			return nil, NoID, fmt.Errorf("Format(%d, 1, %v): %w", visIx, FormatV3Line{}, err)
		}

		cursorID, err = r.Rope.GetVisID(visIx)
		if err != nil {
			return nil, cursorID, fmt.Errorf("Rope.GetVisID(%d): %w", visIx, err)
		}

		ops = append(ops, fop)
		return ops, cursorID, nil
	}

	_, isList := format.(FormatV3BulletList)
	if !isList {
		_, isList = format.(FormatV3OrderedList)
	}
	if isStartOfLine && !isList && sText == "\n" {
		// SINGLE NEWLINE AT START OF LINE IN NON LIST
		fop, err := r.Format(visIx, 1, FormatV3Line{})
		if err != nil {
			return nil, NoID, fmt.Errorf("Format(%d, 1, %v): %w", visIx, FormatV3Line{}, err)
		}
		mop = mop.Append(fop)

		iop, err := r.Insert(visIx, sText)
		if err != nil {
			return nil, NoID, fmt.Errorf("Insert(%d, %q): %w", visIx, sText, err)
		}
		mop = mop.Append(iop)

		fop2, err := r.Format(visIx+1, 1, format)
		if err != nil {
			return nil, NoID, fmt.Errorf("Format(%d, 1, %v): %w", visIx+1, format, err)
		}
		mop = mop.Append(fop2)

		cursorID, err = r.Rope.GetVisID(visIx + 1)
		if err != nil {
			return nil, cursorID, fmt.Errorf("Rope.GetVisID(%d): %w", visIx+1, err)
		}

		r.OpIndex.Put(mop)
		ops = append(ops, mop)
		return ops, cursorID, nil
	}

	_, isHeader := format.(FormatV3Header)
	if isHeader && isEndOfLine && sText == "\n" {
		iop, err := r.Insert(visIx, sText)
		if err != nil {
			return nil, NoID, fmt.Errorf("Insert(%d, %q): %w", visIx, sText, err)
		}
		mop = mop.Append(iop)

		fop, err := r.Format(visIx, 1, format)
		if err != nil {
			return nil, NoID, fmt.Errorf("Format(%d, 1, %v): %w", visIx, format, err)
		}
		mop = mop.Append(fop)

		fop, err = r.Format(visIx+1, 1, FormatV3Line{})
		if err != nil {
			return nil, NoID, fmt.Errorf("Format(%d, 1, %v): %w", visIx+1, FormatV3Line{}, err)
		}
		mop = mop.Append(fop)

		cursorID, err = r.Rope.GetVisID(visIx + 1)
		if err != nil {
			return nil, cursorID, fmt.Errorf("Rope.GetVisID(%d): %w", visIx+1, err)
		}

		r.OpIndex.Put(mop)
		ops = append(ops, mop)
		return ops, cursorID, nil
	}

	// HANDLE INSERTING LONGER TEXT THAT CONTAINS NEWLINES
	iop, err := r.Insert(visIx, sText)
	if err != nil {
		return nil, NoID, fmt.Errorf("Insert(%d, %q): %w", visIx, sText, err)
	}
	mop = mop.Append(iop)

	cursorID, err = r.Rope.GetVisID(visIx + len(text))
	if err != nil {
		return nil, cursorID, fmt.Errorf("Rope.GetVisID(%d): %w", visIx, err)
	}

	if !isPlainLine && !isHeader && !isImage {
		fop, err := r.Format(visIx, len(text), format)
		if err != nil {
			return nil, NoID, fmt.Errorf("Format(%d, 1, %v): %w", visIx, format, err)
		}
		mop = mop.Append(fop)
	} else {
		fop, err := r.Format(visIx, 1, format)
		if err != nil {
			return nil, NoID, fmt.Errorf("Format(%d, 1, %v): %w", visIx, format, err)
		}
		mop = mop.Append(fop)

		fop, err = r.Format(visIx+len(text), 1, FormatV3Line{})
		if err != nil {
			return nil, NoID, fmt.Errorf("Format(%d, 1, %v): %w", visIx+len(text), FormatV3Line{}, err)
		}
		mop = mop.Append(fop)
	}

	visID, err := r.Rope.GetVisID(visIx)
	if err != nil {
		return nil, NoID, fmt.Errorf("Rope.GetVisID(%d): %w", visIx, err)
	}

	curSpanFormat, err := r.GetCurSpanFormat(visID, visID)
	if err != nil {
		return nil, NoID, fmt.Errorf("GetCurSpanFormat(%v, %v): %w", visID, visID, err)
	}

	if !curSpanFormat.Equals(orgSpanFormat) {
		fop, err := r.Format(visIx, len(text), orgSpanFormat)
		if err != nil {
			return nil, NoID, fmt.Errorf("Format(%d, 1, %v): %w", visIx, orgSpanFormat, err)
		}

		mop = mop.Append(fop)
	}

	op := FlattenMop(mop)
	r.OpIndex.Put(op)
	ops = append(ops, op)
	return ops, cursorID, nil
}

func (r *Rogue) _linkSpan(visIx, length int) (op Op, err error) {
	if visIx <= 0 {
		return nil, nil
	}

	mop := MultiOp{}

	id, err := r.Rope.GetVisID(visIx)
	if err != nil {
		return nil, fmt.Errorf("Rope.GetVisID(%d): %w", visIx, err)
	}

	idNode, err := r.NOS.NoSticky.FindLeftSib(id)
	if err != nil {
		return nil, fmt.Errorf("FindLeftSib(%v): %w", id, err)
	}

	if idNode == nil {
		return nil, nil
	}

	ixNode, err := r.NOS.NoSticky.idToIxNosNode(idNode)
	if err != nil {
		return nil, fmt.Errorf("idToIxNosNode(%v): %w", idNode, err)
	}
	if ixNode == nil {
		return nil, nil
	}

	if ixNode.StartIx <= visIx && visIx <= ixNode.EndIx && ixNode.EndIx <= visIx+length {
		fop := FormatOp{
			ID:      r.NextID(1),
			StartID: idNode.StartID,
			EndID:   idNode.EndID,
			Format:  FormatV3Span{"a": ""},
		}

		_, err = r.MergeOp(fop)
		if err != nil {
			return nil, fmt.Errorf("FormatOp(%v): %w", fop, err)
		}
		mop = mop.Append(fop)

		op, err := r.Format(ixNode.StartIx, visIx-ixNode.StartIx, idNode.Format)
		if err != nil {
			return nil, fmt.Errorf("Format(%d, %d, %v): %w", ixNode.StartIx, visIx-ixNode.StartIx, idNode.Format, err)
		}
		mop = mop.Append(op)
	}

	return mop, nil
}

func (r *Rogue) RichDeleteLine(visIx int) (op Op, startIx int, err error) {
	if visIx <= 0 {
		return nil, visIx, nil
	}

	id, err := r.Rope.GetVisID(visIx - 1)
	if err != nil {
		return r.RichDelete(0, visIx)
	}

	lid, err := r.VisScanLeftOf(id, '\n')
	if err != nil {
		return nil, visIx, fmt.Errorf("VisScanRightOfFunc(%v): %w", id, err)
	}
	if lid == nil {
		return r.RichDelete(0, visIx)
	}

	svidx, _, err := r.Rope.GetIndex(*lid)
	if err != nil {
		return nil, visIx, fmt.Errorf("Rope.GetIndex(%v): %w", lid, err)
	}

	if visIx-svidx == 1 {
		// This is the start of the line we want to delete just the newline
		return r.RichDelete(svidx, 1)
	}

	// To follow the behavior of other editors, we do not delete the newline
	return r.RichDelete(svidx+1, visIx-svidx-1)
}

func (r *Rogue) RichDeleteWord(visIx int, forward bool) (op Op, startIx int, err error) {
	if forward {
		if visIx >= r.VisSize {
			return nil, visIx, nil
		}

		id, err := r.Rope.GetVisID(visIx + 1)
		if err != nil {
			return r.RichDelete(visIx, r.VisSize-visIx)
		}

		rid, err := r.VisScanRightOfFunc(id, func(v uint16) bool {
			return unicode.IsSpace(rune(v))
		})
		if err != nil {
			return nil, visIx, fmt.Errorf("VisScanRightOfFunc(%v): %w", id, err)
		}
		if rid == nil {
			return r.RichDelete(visIx, r.VisSize-visIx)
		}

		evidn, _, err := r.Rope.GetIndex(*rid)
		if err != nil {
			return nil, visIx, fmt.Errorf("Rope.GetIndex(%v): %w", rid, err)
		}

		return r.RichDelete(visIx, evidn-visIx)
	} else {
		if visIx <= 0 {
			return nil, visIx, nil
		}

		id, err := r.Rope.GetVisID(visIx - 1)
		if err != nil {
			return r.RichDelete(0, visIx)
		}

		lid, err := r.VisScanLeftOfFunc(id, func(v uint16) bool {
			return unicode.IsSpace(rune(v))
		})
		if err != nil {
			return nil, visIx, fmt.Errorf("VisScanRightOfFunc(%v): %w", id, err)
		}
		if lid == nil {
			return r.RichDelete(0, visIx)
		}

		svidx, _, err := r.Rope.GetIndex(*lid)
		if err != nil {
			return nil, visIx, fmt.Errorf("Rope.GetIndex(%v): %w", lid, err)
		}

		return r.RichDelete(svidx, visIx-svidx)
	}
}

func (r *Rogue) RichDelete(visIx, length int) (op Op, startIx int, err error) {
	if visIx < -1 || visIx+length > r.VisSize {
		return nil, visIx, fmt.Errorf("index: %d length: %d out of bounds for rogue size: %d", visIx, length, r.VisSize)
	}

	if length < 1 {
		return nil, visIx, fmt.Errorf("length: %d is less than 1", length)
	}

	// HANDLE DELETE AT BEGINNING OF DOC
	if length == 1 && visIx == -1 {
		curLineID, err := r.Rope.GetVisID(visIx + 1)
		if err != nil {
			return nil, visIx, fmt.Errorf("Rope.GetVisID(%d): %w", visIx+1, err)
		}

		format, err := r.GetCurLineFormat(curLineID, curLineID)
		if err != nil {
			return nil, visIx, fmt.Errorf("GetCurLineFormat(%v, %v): %w", curLineID, curLineID, err)
		}

		if _, ok := format.(FormatV3Line); !ok {
			fop, err := r.Format(visIx+1, 1, FormatV3Line{})
			if err != nil {
				return nil, visIx, fmt.Errorf("Format(%d, 1, %v): %w", visIx+1, FormatV3Line{}, err)
			}

			r.OpIndex.Put(fop)
			return fop, 0, nil
		}

		return nil, 0, nil
	}

	mop := MultiOp{}
	// HANDLE DELETE RIGHT SIDE OF LINK
	lop, err := r._linkSpan(visIx, length)
	if err != nil {
		return nil, visIx, fmt.Errorf("_linkSpan(%d): %w", visIx, err)
	}
	mop = mop.Append(lop)

	// HANDLE DELETE \n
	startChar := uint16(' ')
	if visIx >= 0 {
		startChar, err = r.GetChar(visIx)
		if err != nil {
			return nil, visIx, fmt.Errorf("GetChar(%d): %w", visIx, err)
		}
	}

	if length == 1 && startChar == '\n' {
		curLineID, err := r.Rope.GetVisID(visIx + 1)
		if err != nil {
			return nil, visIx, fmt.Errorf("Rope.GetVisID(%d): %w", visIx+1, err)
		}

		curLineFormat, err := r.GetCurLineFormat(curLineID, curLineID)
		if err != nil {
			return nil, visIx, fmt.Errorf("GetCurLineFormat(%v, %v): %w", curLineID, curLineID, err)
		}

		cur := getListMeta(curLineFormat)
		if cur.isList {
			var format FormatV3 = FormatV3Line{}
			if cur.indent > 0 {
				switch curLineFormat := curLineFormat.(type) {
				case FormatV3BulletList:
					format = FormatV3BulletList(curLineFormat - 1)
				case FormatV3OrderedList:
					format = FormatV3OrderedList(curLineFormat - 1)
				case FormatV3IndentedLine:
					format = FormatV3IndentedLine(curLineFormat - 1)
				default:
					return nil, visIx, stackerr.Errorf("Only implemented for bullet list, ordered lists, indented lines: %v", curLineFormat)
				}
			}

			fop, err := r.Format(visIx+1, 1, format)
			if err != nil {
				return nil, visIx, fmt.Errorf("Format(%d, 1, %v): %w", visIx+1, format, err)
			}

			r.OpIndex.Put(fop)
			return fop, visIx + 1, nil
		}

		prevIsEmpty, err := r.isEmptyLine(visIx)
		if err != nil {
			return nil, visIx, fmt.Errorf("isEmptyLine(%d): %w", visIx, err)
		}

		_, isHeader := curLineFormat.(FormatV3Header)
		if isHeader && prevIsEmpty {
			dop, err := r.Delete(visIx, 1)
			if err != nil {
				return nil, visIx, fmt.Errorf("Delete(%d, 1): %w", visIx, err)
			}
			mop = mop.Append(dop)

			fop, err := r.Format(visIx, 1, curLineFormat)
			if err != nil {
				return nil, visIx, fmt.Errorf("Format(%d, 1, %v): %w", visIx+1, FormatV3Line{}, err)
			}
			mop = mop.Append(fop)

			r.OpIndex.Put(mop)
			return mop, visIx, nil
		}
	}

	if IsLowSurrogate(startChar) {
		visIx--
	}

	endChar, err := r.GetChar(visIx + length - 1)
	if err != nil {
		return nil, visIx, fmt.Errorf("GetChar(%d): %w", visIx+length-1, err)
	}

	if IsHighSurrogate(endChar) {
		length++
	}

	startID, err := r.Rope.GetVisID(visIx)
	if err != nil {
		return nil, visIx, fmt.Errorf("Rope.GetVisID(%d): %w", visIx, err)
	}

	endID, err := r.Rope.GetVisID(visIx + length - 1)
	if err != nil {
		return nil, visIx, fmt.Errorf("Rope.GetVisID(%d): %w", visIx+length-1, err)
	}

	vis, err := r.Rope.GetBetween(startID, endID)
	if err != nil {
		return nil, visIx, fmt.Errorf("GetBetween(%v, %v): %w", startID, endID, err)
	}

	lineInd := findNewlineIndices(vis.Text)
	var format FormatV3
	if len(lineInd) > 0 {
		ix := lineInd[0]
		id := vis.IDs[ix]
		format, err = r.GetCurLineFormat(id, id)
		if err != nil {
			return nil, visIx, fmt.Errorf("GetCurLineFormat(%v, %v): %w", id, id, err)
		}
	}

	dop, err := r.Delete(visIx, length)
	if err != nil {
		return nil, visIx, fmt.Errorf("Delete(%d, %d): %w", visIx, length, err)
	}
	mop = mop.Append(dop)

	startIx = visIx
	if len(lineInd) > 0 {
		id, err := r.Rope.GetVisID(visIx)
		if err != nil {
			return nil, visIx, fmt.Errorf("Rope.GetVisID(%d): %w", visIx, err)
		}

		nid, err := r.VisScanRightOf(id, '\n')
		if err != nil {
			return nil, visIx, fmt.Errorf("VisScanRightOf(%v, '\\n'): %w", id, err)
		}

		if nid == nil {
			op, err := r.Insert(r.VisSize, "\n")
			if err != nil {
				return nil, visIx, fmt.Errorf("Insert(%d, '\\n'): %w", r.VisSize, err)
			}
			mop = mop.Append(op)
			nid = &op.ID
		}

		visIx, _, err = r.Rope.GetIndex(*nid)
		if err != nil {
			return nil, visIx, fmt.Errorf("Rope.GetIndex(%v): %w", *nid, err)
		}

		fop, err := r.Format(visIx, 1, format)
		if err != nil {
			return nil, visIx, fmt.Errorf("Format(%d, 1, %v): %w", visIx, format, err)
		}
		mop = mop.Append(fop)
	}

	r.OpIndex.Put(mop)

	return mop, startIx, nil
}

func (r *Rogue) _findPrevMatchingFormatIndents(startID ID, format FormatV3) (lessIndent, greaterIndent int, err error) {
	isBullet, isOrdered, isIndentedLine, curIndent := false, false, false, 0
	switch format := format.(type) {
	case FormatV3BulletList:
		isBullet = true
		curIndent = int(format)
	case FormatV3OrderedList:
		isOrdered = true
		curIndent = int(format)
	case FormatV3IndentedLine:
		isIndentedLine = true
		curIndent = int(format)
	default:
		return 0, 0, fmt.Errorf("Only implemented for bullet list, ordered lists, indented lines: %v", format)
	}

	id := startID
	lessIndent, greaterIndent = -1, -1
	for {
		id, err = r.VisLeftOf(id)
		if err != nil {
			if errors.As(err, &ErrorNoLeftVisSibling{}) {
				return lessIndent, greaterIndent, nil
			}
			return 0, 0, fmt.Errorf("VisLeftOf(%v): %w", startID, err)
		}

		lid, err := r.VisScanLeftOf(id, '\n')
		if err != nil {
			return 0, 0, fmt.Errorf("VisScanLeftOf(%v, %q): %w", id, '\n', err)
		}

		if lid == nil {
			return lessIndent, greaterIndent, nil
		}

		prevFormat, err := r.GetCurLineFormat(*lid, *lid)
		if err != nil {
			return 0, 0, fmt.Errorf("GetCurLineFormat(%v, %v): %w", lid, lid, err)
		}

		switch prevFormat := prevFormat.(type) {
		case FormatV3BulletList:
			if isBullet || isIndentedLine {
				prevIndent := int(prevFormat)
				if lessIndent == -1 && prevIndent < curIndent {
					lessIndent = prevIndent
					return lessIndent, greaterIndent, nil
				}

				if greaterIndent == -1 && prevIndent >= curIndent {
					greaterIndent = prevIndent
					if curIndent == 0 {
						return 0, greaterIndent, nil
					}
				}
			}
		case FormatV3OrderedList:
			if isOrdered || isIndentedLine {
				prevIndent := int(prevFormat)
				if lessIndent == -1 && prevIndent < curIndent {
					lessIndent = prevIndent
					return lessIndent, greaterIndent, nil
				}

				if greaterIndent == -1 && prevIndent >= curIndent {
					greaterIndent = prevIndent
					if curIndent == 0 {
						return 0, greaterIndent, nil
					}
				}
			}
		case FormatV3IndentedLine:
			continue
		default:
			return lessIndent, greaterIndent, nil
		}

		id = *lid
	}
}

// Get the list type of the prev list item with the same indent
func (r *Rogue) _findPrevMatchingIndentFormat(startID ID, format FormatV3) (prevIsBullet bool, prevIndent int, err error) {
	isBullet, indent := false, 0
	switch format := format.(type) {
	case FormatV3BulletList:
		isBullet = true
		indent = int(format)
	case FormatV3OrderedList:
		indent = int(format)
	default:
		return false, 0, fmt.Errorf("Only implemented for bullet and ordered lists")
	}

	id := startID

	for {
		id, err = r.VisLeftOf(id)
		if err != nil {
			if errors.As(err, &ErrorNoLeftVisSibling{}) {
				return isBullet, -1, nil
			}
			return false, 0, fmt.Errorf("VisLeftOf(%v): %w", startID, err)
		}

		lid, err := r.VisScanLeftOf(id, '\n')
		if err != nil {
			return false, 0, fmt.Errorf("VisScanLeftOf(%v, %q): %w", id, '\n', err)
		}

		if lid == nil {
			return isBullet, -1, nil
		}

		prevFormat, err := r.GetCurLineFormat(*lid, *lid)
		if err != nil {
			return false, 0, fmt.Errorf("GetCurLineFormat(%v, %v): %w", lid, lid, err)
		}

		switch prevFormat := prevFormat.(type) {
		case FormatV3BulletList:
			prevIndent := int(prevFormat)
			if prevIndent <= indent {
				return true, prevIndent, nil
			}
		case FormatV3OrderedList:
			prevIndent := int(prevFormat)
			if prevIndent <= indent {
				return false, prevIndent, nil
			}
		default:
			return isBullet, -1, nil
		}

		id = *lid
	}
}

var supportedLanguages = map[string]bool{
	"abap":         true,
	"actionscript": true,
	"ada":          true,
	"agda":         true,
	"algol":        true,
	"alice":        true,
	"asm":          true,
	"assembly":     true,
	"autohotkey":   true,
	"awk":          true,
	"bash":         true,
	"basic":        true,
	"batch":        true,
	"bc":           true,
	"brainfuck":    true,
	"c":            true,
	"c#":           true,
	"c++":          true,
	"cecil":        true,
	"ceylon":       true,
	"clojure":      true,
	"cobol":        true,
	"coffeescript": true,
	"common lisp":  true,
	"crystal":      true,
	"css":          true,
	"d":            true,
	"dart":         true,
	"delphi":       true,
	"dreamweaver":  true,
	"elixir":       true,
	"elm":          true,
	"erlang":       true,
	"f#":           true,
	"fortran":      true,
	"go":           true,
	"golang":       true,
	"groovy":       true,
	"haskell":      true,
	"haxe":         true,
	"html":         true,
	"idl":          true,
	"java":         true,
	"javascript":   true,
	"julia":        true,
	"kotlin":       true,
	"latex":        true,
	"lisp":         true,
	"logo":         true,
	"lua":          true,
	"matlab":       true,
	"mercury":      true,
	"nim":          true,
	"objective-c":  true,
	"ocaml":        true,
	"pascal":       true,
	"perl":         true,
	"php":          true,
	"pl/sql":       true,
	"postscript":   true,
	"prolog":       true,
	"python":       true,
	"r":            true,
	"racket":       true,
	"raku":         true,
	"ruby":         true,
	"rust":         true,
	"sas":          true,
	"scala":        true,
	"scheme":       true,
	"scratch":      true,
	"shell":        true,
	"smalltalk":    true,
	"sql":          true,
	"swift":        true,
	"tcl":          true,
	"typescript":   true,
	"vala":         true,
	"vb.net":       true,
	"visual basic": true,
	"webassembly":  true,
	"xml":          true,
	"yaml":         true,
	"zig":          true,
}

func _isCodeBlock(s string) *string {
	if !strings.HasPrefix(s, "```") {
		return nil
	}

	s = strings.TrimPrefix(s, "```")
	s = strings.TrimRight(s, " ") // trim trailing space

	if s == "" {
		return &s
	}

	language := strings.ToLower(s)

	if supportedLanguages[language] {
		return &language
	}

	return nil
}

func _parseMarkdownLink(s string) (int, int, int, string, error) {
	parenCount := 0
	lastParen := -1
	lastBracket := -1

	for i := len(s) - 1; i >= 0; i-- {
		if i > 0 && s[i-1] == '\\' {
			// Skip escaped characters
			i--
			continue
		}

		switch s[i] {
		case ')':
			if parenCount == 0 {
				lastParen = i
			}
			parenCount++
		case '(':
			parenCount--
			if parenCount == 0 && lastParen != -1 {
				// Found matching parentheses
				for j := i - 1; j >= 0; j-- {
					if j > 0 && s[j-1] == '\\' {
						// Skip escaped characters
						j--
						continue
					}
					if s[j] == ']' {
						lastBracket = j
						// Find opening bracket
						for k := j - 1; k >= 0; k-- {
							if k > 0 && s[k-1] == '\\' {
								// Skip escaped characters
								k--
								continue
							}
							if s[k] == '[' {
								url := s[lastBracket+2 : lastParen]
								k = Utf8ToUtf16Ix(s, k)
								lastBracket = Utf8ToUtf16Ix(s, lastBracket)
								lastParen = Utf8ToUtf16Ix(s, lastParen)
								return k, lastBracket, lastParen, url, nil
							}
						}
					}
				}
			}
		}
	}

	return -1, -1, -1, "", fmt.Errorf("no valid markdown link found")
}

func _spanMatchIndices(re *regexp.Regexp, lineText string, visIx, lineStartIx int) (startIx, endIx int) {
	allIndices := re.FindAllStringIndex(lineText, -1)

	startIx, endIx = -1, -1
	for _, indices := range allIndices {
		start, end := indices[0], indices[1]
		end = Utf8ToUtf16Ix(lineText, end-1)
		if end+lineStartIx == visIx {
			start = Utf8ToUtf16Ix(lineText, start)
			startIx = start + lineStartIx
			endIx = end + lineStartIx
			break
		}
	}

	return startIx, endIx
}

func (r *Rogue) IsSpanAllLists(startID, endID ID) (bool, error) {
	_, endIx, err := r.Rope.GetIndex(endID)
	if err != nil {
		return false, err
	}

	var nid *ID = &startID
	for {
		nid, err = r.VisScanRightOf(*nid, '\n')
		if err != nil {
			return false, err
		}

		if nid == nil {
			return false, nil // no more lines, so can't be all lists
		}

		format, err := r.GetCurLineFormat(*nid, *nid)
		if err != nil {
			return false, err
		}

		switch format.(type) {
		case FormatV3BulletList, FormatV3OrderedList, FormatV3IndentedLine:
		default:
			return false, nil // not all lists
		}

		_, nix, err := r.Rope.GetIndex(*nid)
		if err != nil {
			return false, err
		}

		if nix >= endIx {
			break // they are all lists
		}

		*nid, err = r.VisRightOf(*nid)
		if err != nil {
			if errors.As(err, &ErrorNoRightVisSibling{}) {
				return false, nil // no more lines, so can't be all lists
			}

			return false, err
		}
	}

	return true, nil
}

func (r *Rogue) IndentSpan(startID, endID ID, delta int) (op Op, err error) {
	_, endIx, err := r.Rope.GetIndex(endID)
	if err != nil {
		return nil, err
	}

	mop := MultiOp{}
	for {
		sid, err := r.VisScanRightOf(startID, '\n')
		if err != nil {
			return nil, err
		}

		if sid == nil {
			break
		}

		format, err := r.GetCurLineFormat(*sid, *sid)
		if err != nil {
			return nil, err
		}

		meta := getListMeta(format)

		if meta.isList {
			switch f := format.(type) {
			case FormatV3BulletList:
				format = FormatV3BulletList(max(0, int(f)+delta))
			case FormatV3OrderedList:
				format = FormatV3OrderedList(max(0, int(f)+delta))
			case FormatV3IndentedLine:
				format = FormatV3IndentedLine(max(0, int(f)+delta))
			}

			fop := FormatOp{
				ID:      r.NextID(1),
				StartID: *sid,
				EndID:   *sid,
				Format:  format,
			}
			mop = mop.Append(fop)
		}

		_, six, err := r.Rope.GetIndex(*sid)
		if err != nil {
			return nil, err
		}

		if six >= endIx {
			break
		}

		startID, err = r.VisRightOf(*sid)
		if err != nil {
			if errors.As(err, &ErrorNoRightVisSibling{}) {
				break
			}

			return nil, err
		}
	}

	op = FlattenMop(mop)
	_, err = r.MergeOp(op)
	if err != nil {
		return nil, err
	}

	return op, nil
}

// WIP: list tabbing can get a little wonky when deleting in the middle
// of a nested list such that you're first item in the list starts with
// a deep tab. Everything still works and renders, but sometimes you have
// to hit tab twice for a anything to happen on the screen. This is more an
// exception than the rule, but it's something to keep in mind.
//
// DOC: walk back through the list that the startID belongs to, until
// you hit the first item in the list, then return its indent.
// Return -1 if the startID is not in a list
func (r *Rogue) GetListStartOffset(startID ID) (int, error) {
	// this function isn't actually working yet
	startOffset := -1
	for {
		sid, err := r.VisScanLeftOf(startID, '\n')
		if err != nil {
			return -1, err
		}

		if sid == nil {
			return startOffset, nil
		}

		format, err := r.GetCurLineFormat(*sid, *sid)
		if err != nil {
			return 0, err
		}

		meta := getListMeta(format)
		if !meta.isList {
			break
		}

		startID, err = r.VisLeftOf(*sid)
		if err != nil {
			if errors.As(err, &ErrorNoLeftVisSibling{}) {
				return startOffset, nil
			}

			return -1, err
		}
	}

	return startOffset, nil
}
