package v3

import (
	"errors"
)

type RenderSpan struct {
	FirstBlockID ID
	LastBlockID  ID
	ToStartID    ID
	ToEndID      ID
	Html         string
}

// This should be called on an op after it has already been applied
// to the rogue. Used for partial rendering
func (r *Rogue) RenderOp(op Op) (renderSpan *RenderSpan, err error) {
	opSpan, err := r.opSpan(op)
	if err != nil {
		return nil, err
	}

	beforeAddr, err := r.AddressBeforeOp(op)
	if err != nil {
		return nil, err
	}

	afterAddr, err := r.AddressAfterOp(op)
	if err != nil {
		return nil, err
	}

	return r.RenderSpanBetween(opSpan.StartID, opSpan.EndID, beforeAddr, afterAddr)
}

func (r *Rogue) RenderSpan(startID, endID ID, address *ContentAddress) (renderSpan *RenderSpan, err error) {
	firstBlockStartID, _, err := r.GetBlockAt(startID, address)
	if err != nil {
		return nil, err
	}

	lastBlockStartID, lastBlockEndID, err := r.GetBlockAt(endID, address)
	if err != nil {
		return nil, err
	}

	startID, _, err = r.GetBlockAt(startID, nil)
	if err != nil {
		return nil, err
	}

	_, endID, err = r.GetBlockAt(endID, nil)
	if err != nil {
		return nil, err
	}

	_, startAtIx, err := r.Rope.GetIndex(firstBlockStartID)
	if err != nil {
		return nil, err
	}

	_, endAtIx, err := r.Rope.GetIndex(lastBlockEndID)
	if err != nil {
		return nil, err
	}

	_, startIx, err := r.Rope.GetIndex(startID)
	if err != nil {
		return nil, err
	}

	_, endIx, err := r.Rope.GetIndex(endID)
	if err != nil {
		return nil, err
	}

	if startAtIx < startIx {
		startID, _, err = r.GetBlockAt(firstBlockStartID, nil)
		if err != nil {
			return nil, err
		}
	} else if startIx < startAtIx {
		firstBlockStartID, _, err = r.GetBlockAt(startID, address)
		if err != nil {
			return nil, err
		}
	}

	if endIx < endAtIx {
		_, endID, err = r.GetBlockAt(lastBlockEndID, nil)
		if err != nil {
			return nil, err
		}
	} else if endAtIx < endIx {
		lastBlockStartID, lastBlockEndID, err = r.GetBlockAt(endID, address)
		if err != nil {
			return nil, err
		}
	}

	f, err := r.GetLineFormatAt(lastBlockEndID, lastBlockEndID, address)
	if err != nil {
		return nil, err
	}

	lastBlockID := lastBlockEndID
	if _, ok := f.(FormatV3CodeBlock); ok {
		lastBlockID = lastBlockStartID
	}

	html, err := r.GetHtml(startID, endID, true, true)
	if err != nil {
		return nil, err
	}

	// using the start id of the last block because it works for code blocks which
	// don't have an end id in the html currently
	return &RenderSpan{
		FirstBlockID: firstBlockStartID,
		LastBlockID:  lastBlockID,
		ToStartID:    startID,
		ToEndID:      endID,
		Html:         html,
	}, nil
}

func (r *Rogue) RenderSpanBetween(startID, endID ID, beforeAddr, afterAddr *ContentAddress) (renderSpan *RenderSpan, err error) {
	firstBlockStartID, _, err := r.GetBlockAt(startID, beforeAddr)
	if err != nil {
		return nil, err
	}

	lastBlockStartID, lastBlockEndID, err := r.GetBlockAt(endID, beforeAddr)
	if err != nil {
		return nil, err
	}

	startID, _, err = r.GetBlockAt(startID, afterAddr)
	if err != nil {
		return nil, err
	}

	_, endID, err = r.GetBlockAt(endID, afterAddr)
	if err != nil {
		return nil, err
	}

	_, startAtIx, err := r.Rope.GetIndex(firstBlockStartID)
	if err != nil {
		return nil, err
	}

	_, endAtIx, err := r.Rope.GetIndex(lastBlockEndID)
	if err != nil {
		return nil, err
	}

	_, startIx, err := r.Rope.GetIndex(startID)
	if err != nil {
		return nil, err
	}

	_, endIx, err := r.Rope.GetIndex(endID)
	if err != nil {
		return nil, err
	}

	if startAtIx < startIx {
		startID, _, err = r.GetBlockAt(firstBlockStartID, afterAddr)
		if err != nil {
			return nil, err
		}
	} else if startIx < startAtIx {
		firstBlockStartID, _, err = r.GetBlockAt(startID, beforeAddr)
		if err != nil {
			return nil, err
		}
	}

	if endIx < endAtIx {
		_, endID, err = r.GetBlockAt(lastBlockEndID, afterAddr)
		if err != nil {
			return nil, err
		}
	} else if endAtIx < endIx {
		lastBlockStartID, lastBlockEndID, err = r.GetBlockAt(endID, beforeAddr)
		if err != nil {
			return nil, err
		}
	}

	f, err := r.GetLineFormatAt(lastBlockEndID, lastBlockEndID, beforeAddr)
	if err != nil {
		return nil, err
	}

	lastBlockID := lastBlockEndID
	if _, ok := f.(FormatV3CodeBlock); ok {
		lastBlockID = lastBlockStartID
	}

	// DEBUG
	/*
		firstID, err := r.GetFirstTotID()
		if err != nil {
			return nil, err
		}

		lastID, err := r.GetLastTotID()
		if err != nil {
			return nil, err
		}

		beforeHtml, err := r.GetHtmlAt(firstID, lastID, beforeAddr, true, true)
		if err != nil {
			return nil, err
		}

		afterHtml, err := r.GetHtmlAt(firstID, lastID, afterAddr, true, true)
		if err != nil {
			return nil, err
		}

		fmt.Printf("beforeHtml: %q\n", beforeHtml)
		fmt.Printf("afterHtml: %q\n", afterHtml)
	*/
	// DEBUG

	html, err := r.GetHtmlAt(startID, endID, afterAddr, true, true)
	if err != nil {
		return nil, err
	}

	// using the start id of the last block because it works for code blocks which
	// don't have an end id in the html currently
	return &RenderSpan{
		FirstBlockID: firstBlockStartID,
		LastBlockID:  lastBlockID,
		ToStartID:    startID,
		ToEndID:      endID,
		Html:         html,
	}, nil
}

func isBlock(f FormatV3) bool {
	switch f.(type) {
	case FormatV3Line, FormatV3Header, FormatV3Rule, FormatV3Image:
		return false
	}

	return true
}

func _lineFormatEquals(f, g FormatV3) bool {
	fType, gType := "x", "y"

	switch f := f.(type) {
	case FormatV3BulletList, FormatV3OrderedList, FormatV3IndentedLine:
		fType = "list"
	case FormatV3CodeBlock:
		fType = string(f)
	case FormatV3BlockQuote:
		fType = "quote"
	}

	switch g := g.(type) {
	case FormatV3BulletList, FormatV3OrderedList, FormatV3IndentedLine:
		gType = "list"
	case FormatV3CodeBlock:
		gType = string(g)
	case FormatV3BlockQuote:
		gType = "quote"
	}

	return fType == gType
}

func (r *Rogue) GetBlockAt(id ID, address *ContentAddress) (startID, endID ID, err error) {
	curFormat, err := r.GetLineFormatAt(id, id, address)
	if err != nil {
		return NoID, NoID, err
	}

	startID, endID, _, err = r.GetLineAt(id, address)
	if err != nil {
		return NoID, NoID, err
	}

	if !isBlock(curFormat) {
		return startID, endID, nil
	}

	lid, err := r.LeftOfAt(startID, address)
	if err != nil {
		if !errors.As(err, &ErrorNoLeftSiblingAt{}) {
			return NoID, NoID, err
		}

		lid = startID
	}

	for item, err := range r.WalkLeftFromAt(lid, address) {
		if err != nil {
			return NoID, NoID, err
		}

		if item.Char == '\n' {
			format, err := r.GetLineFormatAt(item.ID, item.ID, address)
			if err != nil {
				return NoID, NoID, err
			}

			if !_lineFormatEquals(curFormat, format) {
				break
			}
		}

		startID = item.ID
	}

	rid, err := r.RightOfAt(endID, address)
	if err != nil {
		if !errors.As(err, &ErrorNoRightSiblingAt{}) {
			return NoID, NoID, err
		}

		rid = endID
	}

	for item, err := range r.WalkRightFromAt(rid, address) {
		if err != nil {
			return NoID, NoID, err
		}

		if item.Char == '\n' {
			format, err := r.GetLineFormatAt(item.ID, item.ID, address)
			if err != nil {
				return NoID, NoID, err
			}

			if !_lineFormatEquals(curFormat, format) {
				break
			}

			endID = item.ID
		}
	}

	return startID, endID, nil
}

type OpSpan struct {
	StartID    ID
	EndID      ID
	TotStartIx int
	TotEndIx   int
}

func (r *Rogue) opSpan(op Op) (opSpan OpSpan, err error) {
	startIx, endIx := r.TotSize-1, 0

	ops := []Op{op}
	if mop, ok := op.(MultiOp); ok {
		ops = mop.Mops
	}

	for _, op := range ops {
		sid, eid := NoID, NoID
		switch op := op.(type) {
		case InsertOp:
			sid = op.ID
			eid = MaxID(op)
		case DeleteOp:
			sid = op.TargetID
			eid = ID{Author: op.TargetID.Author, Seq: op.TargetID.Seq + op.SpanLength - 1}
		case FormatOp:
			sid = op.StartID
			eid = op.EndID
		case ShowOp:
			sid = op.TargetID
			eid = ID{Author: op.TargetID.Author, Seq: op.TargetID.Seq + op.SpanLength - 1}
		case RewindOp:
			sid = op.Address.StartID
			eid = op.Address.EndID
		default:
			continue
		}

		_, six, err := r.Rope.GetIndex(sid)
		if err != nil {
			return opSpan, err
		}

		if six < startIx {
			startIx = six
		}

		_, eix, err := r.Rope.GetIndex(eid)
		if err != nil {
			return opSpan, err
		}

		if eix > endIx {
			endIx = eix
		}
	}

	if endIx < startIx {
		startID, err := r.GetFirstID()
		if err != nil {
			return opSpan, err
		}

		endID, err := r.GetLastID()
		if err != nil {
			return opSpan, err
		}

		// not sure how this would happen, but just render whole doc then
		return OpSpan{
			StartID:    startID,
			EndID:      endID,
			TotStartIx: 0,
			TotEndIx:   r.TotSize - 1,
		}, nil
	}

	startID, err := r.Rope.GetTotID(startIx)
	if err != nil {
		return opSpan, err
	}

	endID, err := r.Rope.GetTotID(endIx)
	if err != nil {
		return opSpan, err
	}

	return OpSpan{
		StartID:    startID,
		EndID:      endID,
		TotStartIx: startIx,
		TotEndIx:   endIx,
	}, nil
}

// Op should always come from a single author
func (r *Rogue) AddressBeforeOp(op Op) (*ContentAddress, error) {
	id := op.GetID()
	id = ID{Author: id.Author, Seq: id.Seq - 1}

	address, err := r.GetFullAddress()
	if err != nil {
		return nil, err
	}

	address.MaxIDs[id.Author] = id.Seq
	return address, nil
}

func (r *Rogue) AddressAfterOp(op Op) (*ContentAddress, error) {
	id := MaxID(op)

	address, err := r.GetFullAddress()
	if err != nil {
		return nil, err
	}

	address.MaxIDs[id.Author] = id.Seq
	return address, nil
}
