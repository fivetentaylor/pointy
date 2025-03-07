package v3

import (
	"encoding/json"
	"errors"
	"fmt"
	"unicode"

	"github.com/fivetentaylor/pointy/pkg/stackerr"
	avl "github.com/fivetentaylor/pointy/rogue/v3/tree"
)

type NOSV2 struct {
	rope     *Rope
	Sticky   *NosTree // key is the cur ix of the StartID in the rope
	NoSticky *NosTree // key is the cur ix of the StartID in the rope
	Line     *NosTree // key is the cur ix of the StartID in the rope
}

type NOSV2Node struct {
	StartID ID
	EndID   ID
	Format  FormatV3
}

func NewNOSV2Node(op FormatOp) *NOSV2Node {
	return &NOSV2Node{
		StartID: op.StartID,
		EndID:   op.EndID,
		Format:  op.Format,
	}
}

func (n *NOSV2Node) String() string {
	return fmt.Sprintf("StartID: %v EndID: %v Format: %v", n.StartID, n.EndID, n.Format)
}

func (n NOSV2Node) AsJS() map[string]interface{} {
	return map[string]interface{}{
		"StartID": n.StartID.AsJS(),
		"EndID":   n.EndID.AsJS(),
		"Format":  n.Format.AsMap(),
	}
}

func (nos *NOSV2) ToSlice() []*NOSV2Node {
	out := make([]*NOSV2Node, 0)

	_ = nos.Sticky.Tree.Dft(func(value *NOSV2Node) error {
		out = append(out, value)
		return nil
	})

	_ = nos.NoSticky.Tree.Dft(func(value *NOSV2Node) error {
		out = append(out, value)
		return nil
	})

	_ = nos.Line.Tree.Dft(func(value *NOSV2Node) error {
		out = append(out, value)
		return nil
	})

	return out
}

func (nos *NOSV2) String() string {
	out := ""

	_ = nos.Sticky.Tree.Dft(func(value *NOSV2Node) error {
		out += fmt.Sprintf("    sticky: %s\n", value)
		return nil
	})

	_ = nos.NoSticky.Tree.Dft(func(value *NOSV2Node) error {
		out += fmt.Sprintf("not sticky: %s\n", value)
		return nil
	})

	_ = nos.Line.Tree.Dft(func(value *NOSV2Node) error {
		out += fmt.Sprintf("      line: %s\n", value)
		return nil
	})

	return out
}

func NewNOSV2(rope *Rope) *NOSV2 {
	return &NOSV2{
		rope:     rope,
		Sticky:   NewNosTree(rope, true),
		NoSticky: NewNosTree(rope, false),
		Line:     NewNosTree(rope, false),
	}
}

type NosTree struct {
	KeyFunc func(*NOSV2Node) (int, error)
	rope    *Rope
	Tree    *avl.Tree[int, *NOSV2Node]
	sticky  bool
}

func NewNosTree(rope *Rope, sticky bool) *NosTree {
	keyFunc := func(a *NOSV2Node) (int, error) {
		_, aIx, err := rope.GetIndex(a.StartID)
		if err != nil {
			return 0, err
		}

		return aIx, nil
	}

	return &NosTree{
		KeyFunc: keyFunc,
		rope:    rope,
		Tree:    avl.NewTree(keyFunc, lessThan),
		sticky:  sticky,
	}
}

func (nos *NosTree) mergeSpan(cur, next NOSV2Node) (out []*NOSV2Node, err error) {
	_, curStartIx, err := nos.rope.GetIndex(cur.StartID)
	if err != nil {
		return nil, err
	}

	_, curEndIx, err := nos.rope.GetIndex(cur.EndID)
	if err != nil {
		return nil, err
	}
	if nos.sticky {
		curEndIx = max(0, curEndIx-1)
	}

	_, nextStartIx, err := nos.rope.GetIndex(next.StartID)
	if err != nil {
		return nil, err
	}

	_, nextEndIx, err := nos.rope.GetIndex(next.EndID)
	if err != nil {
		return nil, err
	}
	if nos.sticky {
		nextEndIx = max(0, nextEndIx-1)
	}

	// new node fully to the right of the cur node
	if curEndIx < nextStartIx {
		return []*NOSV2Node{&cur, &next}, nil
	}

	// new node is fully to the left of the cur node
	if nextEndIx < curStartIx {
		return []*NOSV2Node{&next, &cur}, nil
	}

	// handle any leading span
	if curStartIx < nextStartIx {
		// existing span for interval tree
		endID := next.StartID
		if !nos.sticky {
			endID, err = nos.rope.TotLeftOf(next.StartID)
			if err != nil {
				return nil, err
			}
		}
		n := NOSV2Node{cur.StartID, endID, cur.Format}
		out = append(out, &n)
	}

	if nextStartIx < curStartIx {
		// new span for interval tree
		endID := cur.StartID
		if !nos.sticky {
			endID, err = nos.rope.TotLeftOf(cur.StartID)
			if err != nil {
				return nil, err
			}
		}
		n := NOSV2Node{next.StartID, endID, next.Format}
		out = append(out, &n)
	}

	// handle the overlapping section
	mStartID := cur.StartID
	if curStartIx < nextStartIx {
		mStartID = next.StartID
	}

	mEndID := cur.EndID
	if nextEndIx < curEndIx {
		mEndID = next.EndID
	}

	// new span
	merged := mergeFormats(cur.Format, next.Format).DropNull()
	n := NOSV2Node{mStartID, mEndID, merged}
	out = append(out, &n)

	// handle any trailing span
	if curEndIx < nextEndIx {
		startID := cur.EndID
		if !nos.sticky {
			startID, err = nos.rope.TotRightOf(cur.EndID)
			if err != nil {
				return nil, err
			}
		}
		n := NOSV2Node{startID, next.EndID, next.Format}
		out = append(out, &n)
	}

	if nextEndIx < curEndIx {
		startID := next.EndID
		if !nos.sticky {
			startID, err = nos.rope.TotRightOf(next.EndID)
			if err != nil {
				return nil, err
			}
		}
		n := NOSV2Node{startID, cur.EndID, cur.Format}
		out = append(out, &n)
	}

	return out, nil
}

// only called for sticky
func (nos *NosTree) shouldMerge(a, b *NOSV2Node) (bool, error) {
	if !a.Format.Equals(b.Format) {
		return false, nil
	}

	_, aEndIx, err := nos.rope.GetIndex(a.EndID)
	if err != nil {
		return false, err
	}

	if nos.sticky {
		aEndIx = max(0, aEndIx-1)
	}

	_, bStartIx, err := nos.rope.GetIndex(b.StartID)
	if err != nil {
		return false, err
	}

	return aEndIx+1 == bStartIx, nil
}

func (nos *NosTree) insertSpan(newNode *NOSV2Node) ([]*NOSV2Node, error) {
	_, newNodeStartIx, err := nos.rope.GetIndex(newNode.StartID)
	if err != nil {
		return nil, err
	}

	startNode, err := nos.Tree.FindLeftSibNode(newNodeStartIx)
	if err != nil {
		return nil, err
	}

	if startNode == nil {
		startNode, err = nos.Tree.FindRightSibNode(newNodeStartIx)
		if err != nil {
			return nil, err
		}
	}

	if startNode == nil {
		newNode.Format = newNode.Format.DropNull()
		if sf, ok := newNode.Format.(FormatV3Span); ok {
			delete(sf, "e")
			delete(sf, "en")
		}

		if !newNode.Format.Empty() {
			err := nos.Tree.Put(&NOSV2Node{
				StartID: newNode.StartID,
				EndID:   newNode.EndID,
				Format:  newNode.Format,
			})
			if err != nil {
				return nil, err
			}
		}

		return []*NOSV2Node{newNode}, nil
	}

	_, newNodeEndIx, err := nos.rope.GetIndex(newNode.EndID)
	if err != nil {
		return nil, err
	}

	if nos.sticky {
		newNodeEndIx = max(0, newNodeEndIx-1)
	}

	var toInsert []*NOSV2Node
	var toDelete []*NOSV2Node
	err = startNode.WalkRight(func(n *NOSV2Node) error {
		_, nStartIx, err := nos.rope.GetIndex(n.StartID)
		if err != nil {
			return err
		}

		if newNodeEndIx < nStartIx {
			return ErrorStopIteration{}
		}

		add, err := nos.mergeSpan(*n, *newNode)
		if err != nil {
			return err
		}

		newNode = add[len(add)-1]
		toInsert = append(toInsert, add[:len(add)-1]...)
		toDelete = append(toDelete, n)

		return nil
	})

	if err != nil {
		if !errors.As(err, &ErrorStopIteration{}) {
			return nil, err
		}
	}

	toInsert = append(toInsert, newNode)

	for _, n := range toDelete {
		err = nos.Tree.Remove(n)
		if err != nil {
			return nil, err
		}
	}

	// filter inserted to only reflect the spans that represent
	// new formats, which will save on storage in the interval tree
	toReturn := make([]*NOSV2Node, 0, len(toInsert))
	for _, n := range toInsert {
		_, nStartIx, err := nos.rope.GetIndex(n.StartID)
		if err != nil {
			return nil, err
		}

		_, nEndIx, err := nos.rope.GetIndex(n.EndID)
		if err != nil {
			return nil, err
		}

		if nos.sticky {
			nEndIx = max(0, nEndIx-1)
		}

		n.Format = n.Format.DropNull()
		if sf, ok := n.Format.(FormatV3Span); ok {
			delete(sf, "e")
			delete(sf, "en")
		}

		if newNodeStartIx <= nStartIx && nEndIx <= newNodeEndIx {
			toReturn = append(toReturn, &NOSV2Node{
				StartID: n.StartID,
				EndID:   n.EndID,
				Format:  n.Format,
			})
		}
	}

	// merge return nodes
	for i := len(toReturn) - 1; i > 0; i-- {
		a, b := toReturn[i-1], toReturn[i]

		if !a.Format.IsSpan() || !b.Format.IsSpan() {
			continue
		}

		yes, err := nos.shouldMerge(a, b)
		if err != nil {
			return nil, err
		}

		if yes {
			toReturn = DeleteAt(toReturn, i-1, 2)
			toReturn = InsertAt(toReturn, i-1, &NOSV2Node{
				StartID: a.StartID,
				EndID:   b.EndID,
				Format:  a.Format,
			})
		}
	}

	// merge nodes that are right next to each other
	// with the same format
	for i := len(toInsert) - 1; i > 0; i-- {
		a, b := toInsert[i-1], toInsert[i]

		if a.Format.Empty() || b.Format.Empty() {
			continue
		}

		// never merge line formats
		if !a.Format.IsSpan() || !b.Format.IsSpan() {
			continue
		}

		yes, err := nos.shouldMerge(a, b)
		if err != nil {
			return nil, err
		}

		if yes {
			toInsert = DeleteAt(toInsert, i-1, 2)
			toInsert = InsertAt(toInsert, i-1, &NOSV2Node{
				StartID: a.StartID,
				EndID:   b.EndID,
				Format:  a.Format,
			})
		}
	}

	for _, n := range toInsert {
		if n.Format.Empty() {
			continue
		}

		err = nos.Tree.Put(n)
		if err != nil {
			return nil, err
		}
	}

	return toReturn, nil
}

func (nos *NosTree) insertLine(node *NOSV2Node) error {
	err := nos.Tree.Put(node) // just overwrite the current value
	if err != nil {
		return err
	}

	return nil
}

func (nos *NOSV2) Insert(op FormatOp) (ops []FormatOp, err error) {
	if op.Format == nil {
		return nil, nil
	}

	if op.Format.IsSpan() {
		f := op.Format.(FormatV3Span)
		sticky, noSticky := f.SplitSticky()

		if len(sticky) > 0 {
			newNode := NOSV2Node{
				StartID: op.StartID,
				EndID:   op.EndID,
				Format:  sticky,
			}

			inserted, err := nos.Sticky.insertSpan(&newNode)
			if err != nil {
				return nil, err
			}

			for _, n := range inserted {
				if sf, ok := n.Format.(FormatV3Span); ok {
					sf["e"] = "true"
				}

				ops = append(ops, FormatOp{
					ID:      op.ID,
					StartID: n.StartID,
					EndID:   n.EndID,
					Format:  n.Format,
				})
			}
		}

		if len(noSticky) > 0 {
			newNode := NOSV2Node{
				StartID: op.StartID,
				EndID:   op.EndID,
				Format:  noSticky,
			}

			inserted, err := nos.NoSticky.insertSpan(&newNode)
			if err != nil {
				return nil, err
			}

			for _, n := range inserted {
				if sf, ok := n.Format.(FormatV3Span); ok {
					sf["en"] = "true"
				}

				ops = append(ops, FormatOp{
					ID:      op.ID,
					StartID: n.StartID,
					EndID:   n.EndID,
					Format:  n.Format,
				})
			}
		}
	} else {
		newNode := NewNOSV2Node(op)
		err := nos.Line.insertLine(newNode)
		if err != nil {
			return nil, err
		}

		ops = append(ops, op)
	}

	return ops, nil
}

func (nos *NosTree) FindLeftVisSib(visIx int) (*NOSNode, error) {
	if visIx < 0 {
		return nil, nil
	}

	totIx, err := nos.rope.VisToTotIx(visIx)
	if err != nil {
		return nil, fmt.Errorf("VisToTotIx(%v): %w", visIx, err)
	}

	node, err := nos.Tree.FindLeftSibNode(totIx)
	if err != nil {
		return nil, fmt.Errorf("FindLeftSib(%v): %w", visIx, err)
	}

	for {
		if node == nil {
			return nil, nil
		}

		ixNode, err := nos.idToIxNosNode(node.Value)
		if err != nil {
			return nil, fmt.Errorf("idToIxNosNode(%v): %w", node.Value, err)
		}

		if ixNode != nil {
			return ixNode, nil
		}

		node = node.StepLeft()
	}
}

func (nos *NosTree) FindRightVisSib(ix int) (*NOSNode, error) {
	totIx, err := nos.rope.VisToTotIx(ix)
	if err != nil {
		return nil, fmt.Errorf("VisToTotIx(%v): %w", ix, err)
	}

	node, err := nos.Tree.FindRightSibNode(totIx)
	if err != nil {
		return nil, fmt.Errorf("FindRightSib(%v): %w", ix, err)
	}

	for {
		if node == nil {
			return nil, nil
		}

		ixNode, err := nos.idToIxNosNode(node.Value)
		if err != nil {
			return nil, fmt.Errorf("idToIxNosNode(%v): %w", node.Value, err)
		}

		if ixNode != nil {
			return ixNode, nil
		}

		node = node.StepRight()
	}
}

func (nos *NosTree) FindLeftSib(id ID) (*NOSV2Node, error) {
	_, ix, err := nos.rope.GetIndex(id)
	if err != nil {
		return nil, fmt.Errorf("GetIndex(%v): %w", id, err)
	}

	node, err := nos.Tree.FindLeftSib(ix)
	if err != nil {
		return nil, fmt.Errorf("FindLeftSib(%v): %w", ix, err)
	}

	return node, nil
}

func (nos *NosTree) FindRightSib(id ID) (*NOSV2Node, error) {
	_, ix, err := nos.rope.GetIndex(id)
	if err != nil {
		return nil, fmt.Errorf("GetIndex(%v): %w", id, err)
	}

	node, err := nos.Tree.FindRightSib(ix)
	if err != nil {
		return nil, fmt.Errorf("FindRightSib(%v): %w", ix, err)
	}

	return node, nil
}

func (nos *NosTree) Contains(id ID, node *NOSV2Node) (bool, error) {
	_, startIx, err := nos.rope.GetIndex(node.StartID)
	if err != nil {
		return false, fmt.Errorf("GetIndex(%v): %w", node.StartID, err)
	}

	_, endIx, err := nos.rope.GetIndex(node.EndID)
	if err != nil {
		return false, fmt.Errorf("GetIndex(%v): %w", node.EndID, err)
	}

	_, ix, err := nos.rope.GetIndex(id)
	if err != nil {
		return false, fmt.Errorf("GetIndex(%v): %w", id, err)
	}

	return startIx <= ix && ix <= endIx, nil
}

func (nos *NosTree) getCurSpanFormat(startID, endID ID) (FormatV3Span, error) {
	startIx, _, err := nos.rope.GetIndex(startID)
	if err != nil {
		return nil, fmt.Errorf("GetIndex(%v): %w", startID, err)
	}

	if startIx == -1 {
		return nil, fmt.Errorf("startID %v is deleted", startID)
	}

	endIx, _, err := nos.rope.GetIndex(endID)
	if err != nil {
		return nil, fmt.Errorf("GetIndex(%v): %w", endID, err)
	}

	if endIx == -1 {
		return nil, fmt.Errorf("endID %v is deleted", endID)
	}

	if startIx == endIx {
		// if it's just the cursor we want to step back one position
		// since being on the left side of a sticky span should not
		// select its format
		startIx--
	}
	endIx-- // need special handling for non sticky later

	sib, err := nos.FindLeftVisSib(startIx)
	if err != nil {
		return nil, fmt.Errorf("nos.tree.FindLeftSib(%d): %w", startIx, err)
	}

	if sib == nil {
		return nil, nil
	}

	if startIx < sib.StartIx || sib.EndIx < startIx {
		// The span starts in an unformatted area
		return nil, nil
	}

	/*if startIx == endIx && startIx == sib.EndIx {
		// case where cursor is at the end of a span
		return nil, nil
	}*/

	format, ok := sib.Format.(FormatV3Span)
	if !ok {
		return nil, fmt.Errorf("expected FormatV3Span, got %T", sib.Format)
	}

	for {
		if sib.StartIx <= endIx && endIx <= sib.EndIx {
			return format, nil
		}

		nextSib, err := nos.FindRightVisSib(sib.EndIx + 1)
		if err != nil {
			return nil, fmt.Errorf("nos.FindRightSib(%v): %w", sib.EndIx, err)
		}

		if nextSib == nil {
			// The span ends in an unformatted area
			return nil, nil
		}

		if sib.EndIx < nextSib.StartIx-1 {
			// There is an unformated gap in the span
			return nil, nil
		}

		nextFormat, ok := nextSib.Format.(FormatV3Span)
		if !ok {
			return nil, fmt.Errorf("expected FormatV3Span, got %T", nextSib.Format)
		}

		format = IntersectMaps(format, nextFormat)
		if len(format) == 0 {
			// No shared formats across the span
			return nil, nil
		}

		sib = nextSib
	}
}

// GetCurSpanFormat returns the current format of a span which
// must be contiguous within the span
func (r *Rogue) GetCurSpanFormat(startID, endID ID) (FormatV3Span, error) {
	sticky, err := r.NOS.Sticky.getCurSpanFormat(startID, endID)
	if err != nil {
		return nil, fmt.Errorf("sticky.getCurrentSpanFormat(%v, %v): %w", startID, endID, err)
	}

	noSticky, err := r.NOS.NoSticky.getCurSpanFormat(startID, endID)
	if err != nil {
		return nil, fmt.Errorf("noSticky.getCurrentSpanFormat(%v, %v): %w", startID, endID, err)
	}

	f := MergeMaps(sticky, noSticky)
	f["e"] = "true"
	f["en"] = "true"
	return FormatV3Span(f), nil
}

func (r *Rogue) GetLineFormatAt(startID, endID ID, address *ContentAddress) (FormatV3, error) {
	_, endIx, err := r.Rope.GetIndex(endID)
	if err != nil {
		return nil, err
	}

	var curFormat FormatV3
	for item, err := range r.WalkRightFromAt(startID, address) {
		if err != nil {
			return nil, err
		}

		if item.Char != '\n' {
			continue
		}

		if item.Char == '\n' {
			fop, err := r.Formats.Lines.FormatAt(item.ID, address)
			if err != nil {
				return nil, err
			}

			if fop == nil {
				return FormatV3Line{}, nil
			}

			if curFormat == nil {
				curFormat = fop.Format
			} else if !curFormat.Equals(fop.Format) {
				return FormatV3Line{}, nil
			}

			_, ix, err := r.Rope.GetIndex(item.ID)
			if err != nil {
				return nil, err
			}

			if ix >= endIx {
				break
			}
		}
	}

	return curFormat, nil
}

func (r *Rogue) GetCurLineFormat(startID, endID ID) (FormatV3, error) {
	var format FormatV3

	_, endIx, err := r.Rope.GetIndex(endID)
	if err != nil {
		return nil, err
	}

	nextID := &startID
	for {
		nextID, err = r.VisScanRightOf(*nextID, '\n')
		if err != nil {
			return nil, err
		}

		if nextID == nil {
			return FormatV3Line{}, nil
		}

		_, curIx, err := r.Rope.GetIndex(*nextID)
		if err != nil {
			return nil, err
		}

		node, err := r.NOS.Line.Tree.Get(curIx)
		if err != nil {
			return nil, err
		}

		if node == nil {
			return FormatV3Line{}, nil
		}

		if format == nil {
			format = node.Format
		} else if format != node.Format {
			// can stop here because there are different line formats
			// within the span
			return FormatV3Line{}, nil
		}

		if curIx >= endIx {
			return format, nil
		}

		*nextID, err = r.TotRightOf(*nextID)
		if err != nil {
			if errors.As(err, &ErrorNoRightTotSibling{}) {
				return format, nil
			}
			return nil, err
		}
	}
}

func (r *Rogue) GetPrevLineFormat(id ID) (FormatV3, error) {
	id, err := r.Rope.NearestVisLeftOf(id)
	if err != nil {
		return nil, err
	}

	c, err := r.GetCharByID(id)
	if err != nil {
		return nil, err
	}

	if c != '\n' {
		lid, err := r.VisScanLeftOf(id, '\n')
		if err != nil {
			return nil, err
		}

		if lid == nil {
			return FormatV3Line{}, nil // no previous line
		}

		id = *lid
	}

	format, err := r.GetCurLineFormat(id, id)
	if err != nil {
		return nil, err
	}

	return format, nil
}

func (nos *NosTree) slice(startID, endID ID, callback func(*NOSV2Node) error) error {
	_, startIx, err := nos.rope.GetIndex(startID)
	if err != nil {
		return fmt.Errorf("GetIndex(%v): %w", startID, err)
	}

	startNode, err := nos.Tree.FindLeftSibNode(startIx)
	if err != nil {
		return fmt.Errorf("FindLeftSib(%v): %w", startIx, err)
	}

	if startNode == nil {
		startNode, err = nos.Tree.FindRightSibNode(startIx)
		if err != nil {
			return fmt.Errorf("FindRightSib(%v): %w", startIx, err)
		}
	}

	_, endIx, err := nos.rope.GetIndex(endID)
	if err != nil {
		return fmt.Errorf("GetIndex(%v): %w", endID, err)
	}

	err = startNode.WalkRight(func(node *NOSV2Node) error {
		_, nStartIx, err := nos.rope.GetIndex(node.StartID)
		if err != nil {
			return fmt.Errorf("GetIndex(%v): %w", node.StartID, err)
		}

		_, nEndIx, err := nos.rope.GetIndex(node.EndID)
		if err != nil {
			return fmt.Errorf("GetIndex(%v): %w", node.EndID, err)
		}

		if nEndIx < startIx {
			return nil // continue
		}

		if endIx < nStartIx {
			return ErrorStopIteration{} // break
		}

		containsStart := nStartIx <= startIx && startIx <= nEndIx
		containsEnd := nStartIx <= endIx && endIx <= nEndIx

		sid, eid := node.StartID, node.EndID
		if containsStart {
			sid = startID
		}

		if containsEnd {
			eid = endID
		}

		err = callback(&NOSV2Node{
			StartID: sid,
			EndID:   eid,
			Format:  node.Format,
		})

		if err != nil {
			return err
		}

		if containsEnd {
			return ErrorStopIteration{}
		}

		return nil
	})

	if err != nil {
		if errors.As(err, &ErrorStopIteration{}) {
			return nil
		}

		return fmt.Errorf("WalkRight: %w", err)
	}

	return nil
}

func isBadID(id ID) bool {
	badIDs := []ID{
		{Author: "0", Seq: 609},
		{Author: "0", Seq: 612},
		{Author: "0", Seq: 656},
	}

	for _, bad := range badIDs {
		if id == bad {
			return true
		}
	}

	return false
}

func (nos *NosTree) _idNosToIxNos(diff *FugueDiff, node *NOSV2Node) (*NOSNode, error) {
	_, startTotIx, err := nos.rope.GetIndex(node.StartID)
	if err != nil {
		return nil, fmt.Errorf("GetIndex(%v): %w", node.StartID, err)
	}

	startOffset := bisectLeft(diff.TotIxs, startTotIx, Identity)
	endOffset := startOffset

	if startOffset == len(diff.TotIxs) {
		return nil, nil
	}

	endTotIx := startTotIx
	if node.StartID != node.EndID {
		_, endTotIx, err = nos.rope.GetIndex(node.EndID)
		if err != nil {
			return nil, fmt.Errorf("GetIndex(%v): %w", node.EndID, err)
		}
	} else if diff.TotIxs[startOffset] != startTotIx {
		// line format that's been deleted or empty sticky span
		return nil, nil
	}

	if nos.sticky {
		endTotIx--
	}

	if endTotIx < startTotIx {
		return nil, nil
	}

	endOffset = bisectLeft(diff.TotIxs, endTotIx, Identity)
	if (endOffset == len(diff.TotIxs) || diff.TotIxs[endOffset] != endTotIx) && endOffset > 0 {
		endOffset--
	}

	if startOffset <= endOffset {
		return &NOSNode{
			StartIx: startOffset,
			EndIx:   endOffset,
			Format:  node.Format,
		}, nil
	}

	return nil, nil
}

func (nos *NosTree) idToIxNosNode(node *NOSV2Node) (*NOSNode, error) {
	var err error
	startID, endID := node.StartID, node.EndID

	if startID == endID && !node.Format.IsSpan() {
		startIx, _, err := nos.rope.GetIndex(startID)
		if err != nil {
			return nil, fmt.Errorf("GetIndex(%v): %w", startID, err)
		}

		if startIx == -1 {
			return nil, nil
		}

		nosNode := &NOSNode{
			StartIx: startIx,
			EndIx:   startIx,
			Format:  node.Format,
		}

		return nosNode, nil
	}

	startIx, _, err := nos.rope.GetIndex(startID)
	if err != nil {
		return nil, fmt.Errorf("GetIndex(%v): %w", startID, err)
	}

	if startIx == -1 {
		startID, err = nos.rope.VisRightOf(startID)
		if err != nil {
			if errors.As(err, &ErrorNoRightVisSibling{}) {
				return nil, nil
			}
			return nil, fmt.Errorf("VisRightOf(%v): %w", node.StartID, err)
		}

		startIx, _, err = nos.rope.GetIndex(startID)
		if err != nil {
			return nil, fmt.Errorf("GetIndex(%v): %w", startID, err)
		}
	}

	if nos.sticky {
		endID, err = nos.rope.TotLeftOf(endID)
		if err != nil {
			if errors.As(err, &ErrorNoLeftTotSibling{}) {
				return nil, nil
			}
			return nil, fmt.Errorf("TotLeftOf(%v): %w", node.EndID, err)
		}
	}

	endIx, _, err := nos.rope.GetIndex(endID)
	if err != nil {
		return nil, fmt.Errorf("GetIndex(%v): %w", endID, err)
	}

	if endIx == -1 {
		endID, err = nos.rope.VisLeftOf(endID)
		if err != nil {
			if errors.As(err, &ErrorNoLeftVisSibling{}) {
				return nil, nil
			}
			return nil, fmt.Errorf("VisLeftOf(%v): %w", node.EndID, err)
		}

		endIx, _, err = nos.rope.GetIndex(endID)
		if err != nil {
			return nil, fmt.Errorf("GetIndex(%v): %w", endID, err)
		}
	}

	if startIx > endIx {
		return nil, nil
	}

	nosNode := &NOSNode{
		StartIx: startIx,
		EndIx:   endIx,
		Format:  node.Format,
	}

	return nosNode, nil
}

func (r *Rogue) ToIndexNos(startID, endID ID, address *ContentAddress, smartQuote bool) (diff *FugueDiff, span, line *NOS, err error) {
	r.Rope.SetCache()
	defer r.Rope.DelCache()

	diff, err = r.Filter(startID, endID, address)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Filter(%v, %v): %w", startID, endID, err)
	}

	if diff == nil || len(diff.Text) == 0 {
		return nil, nil, nil, nil
	}

	line = StringToLineNOS(diff.Text)

	nos := r.NOS

	// build filtered nos if necessary
	if address != nil {
		nos = NewNOSV2(r.Rope)

		ops, err := r.rewindFormatTo(startID, endID, address, NoID)
		if err != nil {
			return nil, nil, nil, err
		}

		for _, op := range ops {
			_, err := nos.Insert(op)
			if err != nil {
				return nil, nil, nil, err
			}
		}
	}

	// BEGIN: Get format from closing line if the selection stops short of the newline
	_, startIx, err := r.Rope.GetIndex(diff.IDs[0])
	if err != nil {
		return nil, nil, nil, err
	}

	startLineID, endLineID, _, err := r.GetLineAt(endID, address)
	if err != nil {
		return nil, nil, nil, err
	}

	_, startLineIx, err := r.Rope.GetIndex(startLineID)
	if err != nil {
		return nil, nil, nil, err
	}

	// if the selection starts at the beginning of the last line or before it
	if startIx < startLineIx {
		_, endLineTotIx, err := r.Rope.GetIndex(endLineID)
		if err != nil {
			return nil, nil, nil, err
		}

		lineNode, err := nos.Line.Tree.Get(endLineTotIx)
		if err != nil {
			return nil, nil, nil, err
		}

		endLine, err := line.tree.FindLeftSib(len(diff.Text))
		if err != nil {
			return nil, nil, nil, err
		}

		if endLine != nil && lineNode != nil && lineNode.Format != nil {
			endLine.Format = lineNode.Format.Copy()
		}
	}
	// END: Get format from closing line if the selection stops short of the newline

	// Get any real line formats from line nos id and transfer to index nos
	err = nos.Line.slice(startID, endID, func(node *NOSV2Node) error {
		ixNode, err := nos.Line._idNosToIxNos(diff, node)
		if err != nil {
			return fmt.Errorf("nos._idNosToIxNos(%v, %v): %w", diff, node, err)
		}

		if ixNode == nil {
			return nil
		}

		n, err := line.tree.FindLeftSib(ixNode.EndIx)
		if err != nil {
			return fmt.Errorf("Get(%v): %w", ixNode.EndIx, err)
		}

		if n != nil && n.EndIx == ixNode.EndIx {
			n.Format = ixNode.Format
		}

		return nil
	})

	if err != nil {
		return nil, nil, nil, fmt.Errorf("line.slice(%v, %v): %w", startID, endID, err)
	}

	span = NewNOS()

	// seed span with line formats to split up spans
	err = line.tree.Dft(func(node *NOSNode) error {
		newNode := NOSNode{
			StartIx: node.EndIx,
			EndIx:   node.EndIx,
			Format:  node.Format,
		}
		return span.tree.Put(&newNode)
	})

	if err != nil {
		return nil, nil, nil, fmt.Errorf("line.tree.Dft: %w", err)
	}

	err = nos.Sticky.slice(startID, endID, func(node *NOSV2Node) error {
		ixNode, err := nos.Sticky._idNosToIxNos(diff, node)
		if err != nil {
			return fmt.Errorf("nos._idNosToIxNos(%v, %v): %w", diff, node, err)
		}

		if ixNode == nil {
			return nil
		}

		err = span.Insert(*ixNode)
		if err != nil {
			return fmt.Errorf("Insert(%v): %w", ixNode, err)
		}

		return nil
	})

	if err != nil {
		return nil, nil, nil, fmt.Errorf("sticky.slice(%v, %v): %w", startID, endID, err)
	}

	err = nos.NoSticky.slice(startID, endID, func(node *NOSV2Node) error {
		ixNode, err := nos.NoSticky._idNosToIxNos(diff, node)
		if err != nil {
			return fmt.Errorf("nos._idNosToIxNos(%v, %v): %w", diff, node, err)
		}

		if ixNode == nil {
			return nil
		}

		err = span.Insert(*ixNode)
		if err != nil {
			return fmt.Errorf("Insert(%v): %w", ixNode, err)
		}

		return nil
	})

	if err != nil {
		return nil, nil, nil, fmt.Errorf("noSticky.slice(%v, %v): %w", startID, endID, err)
	}

	// smart quote handling
	if smartQuote {
		for i, c := range diff.Text {
			if c == '"' || c == '\'' {
				// unicode whitespace is all in the bmp so this works
				isLeftSpace, isRightSpace := false, false

				if i == 0 || unicode.IsSpace(rune(diff.Text[i-1])) {
					isLeftSpace = true
				}

				if i == len(diff.Text)-1 || unicode.IsSpace(rune(diff.Text[i+1])) {
					isRightSpace = true
				}

				if isLeftSpace != isRightSpace {
					if isLeftSpace {
						span.Insert(NOSNode{
							StartIx: i,
							EndIx:   i,
							Format:  FormatV3Span{"ql": "true"},
						})
					} else {
						span.Insert(NOSNode{
							StartIx: i,
							EndIx:   i,
							Format:  FormatV3Span{"qr": "true"},
						})
					}
				}
			}
		}
	}

	// DEBUG
	/*span.tree.Dft(func(node *NOSNode) error {
		fmt.Printf("span node: %v\n", node)
		return nil
	})*/

	span = span.mergeNeighbors()

	return diff, span, line, nil
}

func (r *Rogue) GetNearestVisID(id ID) (ID, error) {
	if r.VisSize == 0 {
		return NoID, stackerr.New(fmt.Errorf("no visible text"))
	}

	isDel, err := r.Rope.IsDeleted(id)
	if err != nil {
		return NoID, err
	}

	outID := id
	if isDel {
		outID, err = r.Rope.VisLeftOf(id)
		if err != nil {
			if !errors.As(err, &ErrorNoLeftVisSibling{}) {
				return NoID, err
			} else {
				outID, err = r.Rope.VisRightOf(id)
				if err != nil {
					return NoID, err
				}
			}
		}
	}

	return outID, nil
}

func (r *Rogue) _enclosingCodeBlock(id ID, address *ContentAddress) (startID ID, offset int, err error) {
	offset, startID = 0, id
	var format FormatV3

	format, err = r.GetLineFormatAt(id, id, address)
	if err != nil {
		return NoID, -1, fmt.Errorf("GetLineFormatAt(%v): %w", id, err)
	}

	if _, ok := format.(FormatV3CodeBlock); !ok {
		addrBytes, err := json.Marshal(address)
		if err != nil {
			return NoID, -1, stackerr.New(err)
		}

		return NoID, -1, stackerr.New(ErrorNotCodeblock{
			ID:      id,
			Address: string(addrBytes),
		})
	}

	isFirst := true
	for item, err := range r.WalkLeftFromAt(id, address) {
		if err != nil {
			return NoID, -1, err
		}

		if isFirst {
			isFirst = false
			continue
		}

		if item.Char == '\n' {
			prevFormat, err := r.GetLineFormatAt(item.ID, item.ID, address)
			if err != nil {
				return NoID, -1, err
			}

			if !format.Equals(prevFormat) {
				break
			}
		}

		startID = item.ID
		offset++
	}

	return startID, offset, nil
}

func (r *Rogue) EnclosingSpanID(id ID, address *ContentAddress, smartQuote bool) (startID ID, offset int, err error) {
	// HANDLE CODE BLOCK
	startID, offset, err = r._enclosingCodeBlock(id, address)
	if err == nil {
		return startID, offset, nil
	} else {
		if !errors.As(err, &ErrorNotCodeblock{}) {
			return NoID, -1, err
		}
	}

	firstID, lastID, offset, err := r.GetLineAt(id, address)
	if err != nil {
		return NoID, -1, err
	}

	if firstID == lastID {
		return firstID, 0, nil
	}

	diff, span, _, err := r.ToIndexNos(firstID, lastID, address, smartQuote)
	if err != nil {
		return NoID, -1, err
	}

	item := span.tree.Max()
	if !item.Format.IsSpan() {
		// drop the line format from the spans
		err = span.tree.RemoveByKey(item.EndIx)
		if err != nil {
			return NoID, -1, err
		}
	}

	n, err := span.tree.FindLeftSib(offset)
	if err != nil {
		return NoID, -1, err
	}

	if n == nil {
		return firstID, offset, nil
	}

	if n.StartIx <= offset && offset <= n.EndIx+1 {
		startID = diff.IDs[n.StartIx]
		offset = offset - n.StartIx
	} else {
		startID = diff.IDs[n.EndIx+1]
		offset = offset - n.EndIx - 1
	}

	return startID, offset, nil
}
