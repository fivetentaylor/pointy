package v3

import (
	"fmt"
	"strings"
)

func (r *Rogue) PasteHtml(idx int, html string) (MultiOp, error) {
	mop := MultiOp{}

	if idx < 0 {
		return mop, fmt.Errorf("index out of bounds: %d", idx)
	}
	if html == "" {
		return mop, nil
	}

	text, spans, err := ParseHtml(html)
	if err != nil {
		return mop, fmt.Errorf("ParseHtml(%q): failed to parse HTML: %w", html, err)
	}

	mop.Mops = make([]Op, 0, len(spans)+1)

	op, err := r.Insert(idx, text)
	if err != nil {
		return mop, err
	}

	mop = mop.Append(op)

	for _, span := range spans {
		// This is currently happening with image tags that have no content,
		// will have to figure out how to handle these eventually
		if span.EndIndex <= span.StartIndex {
			continue
		}

		fop, err := r.Format(idx+span.StartIndex, (span.EndIndex - span.StartIndex), span.Format)
		if err != nil {
			return mop, fmt.Errorf("Format(%d, %d, %v): %w", idx+span.StartIndex, idx+span.EndIndex, span.Format, err)
		}
		mop = mop.Append(fop)
	}

	return mop, nil
}

type PasteItem struct {
	Kind string
	Mime string
	Data string
}

func (r *Rogue) Paste(visIx, selLen int, orgSpanFormat FormatV3Span, items []PasteItem) (ops []Op, cursorID ID, err error) {
	mop := MultiOp{}
	plaintext := ""
	spans := []TextSpan{}

	if _checkForMimeType(items, "text/_notion") {
		// fmt.Println("FROM NOTION")
	} else if _checkForMimeType(items, "application/x-vnd.google-docs-document-slice-clip+wrapped") {
		// fmt.Println("FROM GOOGLE DOCS")
	} else if _fromApple(items) {
		// fmt.Println("FROM APPLE")
	}

	if _checkForMimeType(items, "text/html") {
		// fmt.Println("FROM HTML")
		data := _getMimeType(items, "text/html")
		plaintext, spans, err = ParseHtml(data)
		if err != nil {
			return nil, NoID, fmt.Errorf("ParseHtml(%q): %w", data, err)
		}
	} else if _checkForMimeType(items, "text/plain") {
		plaintext = _getMimeType(items, "text/plain")
		return r.RichInsert(visIx, selLen, orgSpanFormat, plaintext)
	} else {
		// fmt.Println("FROM UNKNOWN")
	}

	if plaintext == "" {
		cursorID, err = r.Rope.GetVisID(visIx + selLen)
		if err != nil {
			return nil, NoID, err
		}

		return nil, cursorID, nil
	}

	if selLen > 0 {
		dop, err := r.Delete(visIx, selLen)
		if err != nil {
			return nil, NoID, fmt.Errorf("Delete(%d, %d): %w", visIx, selLen, err)
		}
		mop = mop.Append(dop)
	}

	iop, err := r.Insert(visIx, plaintext)
	if err != nil {
		return nil, NoID, fmt.Errorf("Insert(%d, %q): %w", visIx, plaintext, err)
	}
	mop = mop.Append(iop)

	for _, span := range spans {
		// This is currently happening with image tags that have no content,
		// will have to figure out how to handle these eventually
		if span.EndIndex <= span.StartIndex {
			continue
		}

		fop, err := r.Format(visIx+span.StartIndex, (span.EndIndex - span.StartIndex), span.Format)
		if err != nil {
			return nil, NoID, fmt.Errorf("Format(%d, %d, %v): %w", visIx+span.StartIndex, visIx+span.EndIndex, span.Format, err)
		}
		mop = mop.Append(fop)
	}
	r.OpIndex.Put(mop)

	// Set cursor at end of the pasted text
	cursorIx := visIx + len(StrToUint16(plaintext))
	if len(plaintext) > 0 && plaintext[len(plaintext)-1] == '\n' {
		cursorIx--
	}

	cursorID, err = r.Rope.GetVisID(cursorIx)
	if err != nil {
		return nil, NoID, fmt.Errorf("GetVisID(%d): %w", cursorIx, err)
	}

	return []Op{mop}, cursorID, nil
}

func _checkForMimeType(items []PasteItem, mimePrefix string) bool {
	for _, item := range items {
		if item.Kind == "string" && strings.HasPrefix(item.Mime, mimePrefix) {
			return true
		}
	}

	return false
}

func _fromApple(items []PasteItem) bool {
	for _, item := range items {
		if item.Kind == "string" && item.Mime == "text/rtf" {
			if strings.Contains(item.Data, "cocoartf2761") {
				return true
			}
		}
	}

	return false
}

func _getMimeType(items []PasteItem, mimePrefix string) string {
	for _, item := range items {
		if item.Kind == "string" && strings.HasPrefix(item.Mime, mimePrefix) {
			return item.Data
		}
	}

	return ""
}
