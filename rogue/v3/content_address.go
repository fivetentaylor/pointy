package v3

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"

	"github.com/fivetentaylor/pointy/pkg/stackerr"
	avl "github.com/fivetentaylor/pointy/rogue/v3/tree"
	"golang.org/x/exp/slices"
)

func NewContentAddress() *ContentAddress {
	return &ContentAddress{
		MaxIDs: map[string]int{},
	}
}

func ParseContentAddress(s string) (*ContentAddress, error) {
	var ca ContentAddress
	err := json.Unmarshal([]byte(s), &ca)
	if err != nil {
		return nil, stackerr.Wrap(err)
	}

	return &ca, nil
}

func (ca *ContentAddress) MaxAuthorID(author string) ID {
	seq, ok := ca.MaxIDs[author]
	if !ok {
		return ID{author, -1}
	}

	return ID{author, seq}
}

func (ca *ContentAddress) DeepCopy() *ContentAddress {
	newCA := &ContentAddress{
		StartID: ca.StartID,
		EndID:   ca.EndID,
		MaxIDs:  make(map[string]int),
	}

	for key, value := range ca.MaxIDs {
		newCA.MaxIDs[key] = value
	}

	return newCA
}

func (ca *ContentAddress) AddID(id ID) {
	if maxSeq, ok := ca.MaxIDs[id.Author]; ok {
		if id.Seq > maxSeq {
			ca.MaxIDs[id.Author] = id.Seq
		}
	} else {
		ca.MaxIDs[id.Author] = id.Seq
	}
}

func (r *Rogue) GetEmptyAddress() (*ContentAddress, error) {
	startID, err := r.Rope.GetTotID(0)
	if err != nil {
		return nil, fmt.Errorf("GetTotID(0): %w", err)
	}

	endID, err := r.Rope.GetTotID(r.TotSize - 1)
	if err != nil {
		return nil, fmt.Errorf("GetTotID(%v): %w", r.TotSize-1, err)
	}

	return &ContentAddress{
		StartID: startID,
		EndID:   endID,
		MaxIDs:  map[string]int{},
	}, nil
}

func (r *Rogue) GetFullAddress() (*ContentAddress, error) {
	startID, err := r.Rope.GetTotID(0)
	if err != nil {
		return nil, err
	}

	endID, err := r.Rope.GetTotID(r.TotSize - 1)
	if err != nil {
		return nil, err
	}

	address := r.OpIndex.ContentAddress.DeepCopy()
	address.StartID = startID
	address.EndID = endID

	return address, nil
}

func (r *Rogue) GetAddress(startID, endID ID) (*ContentAddress, error) {
	address := ContentAddress{
		StartID: startID,
		EndID:   endID,
		MaxIDs:  map[string]int{},
	}

	tot, err := r.Rope.GetTotBetween(startID, endID)
	if err != nil {
		return nil, fmt.Errorf("GetTotBetween(%v, %v): %w", startID, endID, err)
	}

	for _, id := range tot.IDs {
		address.AddID(id)

		ch := r.CharHistory[id]
		if ch != nil {
			ch.Dft(func(m *Marker) error {
				address.AddID(m.ID)
				return nil
			})
		}
	}

	// GET FORMAT IDs
	_, startIx, err := r.Rope.GetIndex(startID)
	if err != nil {
		return nil, fmt.Errorf("GetIndex(%v): %w", startID, err)
	}

	_, endIx, err := r.Rope.GetIndex(endID)
	if err != nil {
		return nil, fmt.Errorf("GetIndex(%v): %w", endID, err)
	}

	stickyOps, err := r.Formats.Sticky.SearchOverlapping(startIx, endIx)
	if err != nil {
		return nil, fmt.Errorf("SearchOverlapping(%d, %d): %w", startIx, endIx, err)
	}

	for _, f := range stickyOps {
		address.AddID(f.ID)
	}

	noStickyOps, err := r.Formats.NoSticky.SearchOverlapping(startIx, endIx)
	if err != nil {
		return nil, fmt.Errorf("SearchOverlapping(%d, %d): %w", startIx, endIx, err)
	}

	for _, f := range noStickyOps {
		address.AddID(f.ID)
	}

	r.Formats.Lines.Tree.Slice(startIx, endIx, func(lh *LineHistory) error {
		lh.Formats.Dft(func(f FormatOp) error {
			address.AddID(f.ID)
			return nil
		})
		return nil
	})

	return &address, nil
}

type IndexedNode struct {
	Index int
	Node  *FugueNode
}

func (r *Rogue) _trimFilterOp(caStartIx, caEndIx int, address ContentAddress, op Op) (ops []Op, err error) {
	switch op := op.(type) {
	case InsertOp:
		_, startIx, err := r.Rope.GetIndex(op.ID)
		if err != nil {
			return nil, fmt.Errorf("GetIndex(%v): %w", op.ID, err)
		}

		uText := StrToUint16(op.Text)
		endID := ID{op.ID.Author, op.ID.Seq + len(uText) - 1}
		_, endIx, err := r.Rope.GetIndex(endID)
		if err != nil {
			return nil, fmt.Errorf("GetIndex(%v): %w", endID, err)
		}

		if endIx < caStartIx || caEndIx < startIx {
			return nil, nil // continue since this isn't in the address
		}

		if startIx < caStartIx {
			for startIx < caStartIx {
				op.ID = ID{op.ID.Author, op.ID.Seq + 1}
				uText = uText[1:]
				_, startIx, err = r.Rope.GetIndex(op.ID)
				if err != nil {
					return nil, fmt.Errorf("GetIndex(%v): %w", op.ID, err)
				}
			}

			if op.Side == Right {
				op.Side = Root
			}
		}

		if caEndIx < endIx {
			for caEndIx < endIx {
				uText = uText[:len(uText)-1]
				endID = ID{endID.Author, endID.Seq - 1}
				_, endIx, err = r.Rope.GetIndex(endID)
				if err != nil {
					return nil, fmt.Errorf("GetIndex(%v): %w", endID, err)
				}
			}

			if op.Side == Left {
				op.Side = Root
			}
		}

		if r.ContainsID(op.ParentID) {
			_, parentIx, err := r.Rope.GetIndex(op.ParentID)
			if err != nil {
				return nil, fmt.Errorf("GetIndex(%v): %w", op.ParentID, err)
			}

			if parentIx < caStartIx || caEndIx < parentIx {
				op.Side = Root
			}
		}

		if len(uText) == 0 {
			return nil, nil
		}

		if op.Side == Root {
			op.ParentID = ID{"", startIx}
		}

		op.Text = Uint16ToStr(uText)
		ops = append(ops, op)
	case DeleteOp:
		_, startIx, err := r.Rope.GetIndex(op.TargetID)
		if err != nil {
			return nil, fmt.Errorf("GetIndex(%v): %w", op.TargetID, err)
		}
		if startIx > caEndIx {
			return nil, nil // delete is outside of the address
		}

		endID := ID{op.TargetID.Author, op.TargetID.Seq + op.SpanLength - 1}
		_, endIx, err := r.Rope.GetIndex(endID)
		if err != nil {
			return nil, fmt.Errorf("GetIndex(%v): %w", endID, err)
		}

		if endIx < caStartIx {
			return nil, nil // delete is outside of the address
		}

		if startIx < caStartIx {
			for i := 1; i < op.SpanLength; i++ {
				targetID := ID{Author: op.TargetID.Author, Seq: op.TargetID.Seq + i}
				_, startIx, err = r.Rope.GetIndex(targetID)
				if err != nil {
					return nil, fmt.Errorf("GetIndex(%v): %w", targetID, err)
				}

				if caStartIx <= startIx {
					op.SpanLength -= i
					op.TargetID = targetID
					break
				}
			}
		}

		if endIx > caEndIx {
			for i := op.SpanLength - 2; i >= 0; i-- {
				targetID := ID{Author: op.TargetID.Author, Seq: op.TargetID.Seq + i}
				_, endIx, err = r.Rope.GetIndex(targetID)
				if err != nil {
					return nil, fmt.Errorf("GetIndex(%v): %w", targetID, err)
				}

				if endIx <= caEndIx {
					op.SpanLength = i + 1
					break
				}
			}
		}

		if endIx < caStartIx || caEndIx < startIx {
			return nil, nil
		}

		if op.SpanLength == 0 {
			return nil, nil
		}

		ops = append(ops, op)
	case ShowOp:
		_, startIx, err := r.Rope.GetIndex(op.TargetID)
		if err != nil {
			return nil, fmt.Errorf("GetIndex(%v): %w", op.TargetID, err)
		}
		if startIx > caEndIx {
			return nil, nil // delete is outside of the address
		}

		endID := ID{op.TargetID.Author, op.TargetID.Seq + op.SpanLength - 1}
		_, endIx, err := r.Rope.GetIndex(endID)
		if err != nil {
			return nil, fmt.Errorf("GetIndex(%v): %w", endID, err)
		}

		if endIx < caStartIx {
			return nil, nil // delete is outside of the address
		}

		if startIx < caStartIx {
			for i := 1; i < op.SpanLength; i++ {
				targetID := ID{Author: op.TargetID.Author, Seq: op.TargetID.Seq + i}
				_, startIx, err = r.Rope.GetIndex(targetID)
				if err != nil {
					return nil, fmt.Errorf("GetIndex(%v): %w", targetID, err)
				}

				if caStartIx <= startIx {
					op.SpanLength -= i
					op.TargetID = targetID
					break
				}
			}
		}

		if endIx > caEndIx {
			for i := op.SpanLength - 2; i >= 0; i-- {
				targetID := ID{Author: op.TargetID.Author, Seq: op.TargetID.Seq + i}
				_, endIx, err = r.Rope.GetIndex(targetID)
				if err != nil {
					return nil, fmt.Errorf("GetIndex(%v): %w", targetID, err)
				}

				if endIx <= caEndIx {
					op.SpanLength = i + 1
					break
				}
			}
		}

		if endIx < caStartIx || caEndIx < startIx {
			return nil, nil
		}

		if op.SpanLength == 0 {
			return nil, nil
		}

		ops = append(ops, op)
	case FormatOp:
		_, fStartIx, err := r.Rope.GetIndex(op.StartID)
		if err != nil {
			return nil, fmt.Errorf("GetIndex(%v): %w", op.StartID, err)
		}

		_, fEndIx, err := r.Rope.GetIndex(op.EndID)
		if err != nil {
			return nil, fmt.Errorf("GetIndex(%v): %w", op.EndID, err)
		}

		// format doesn't overlap
		if fEndIx < caStartIx || caEndIx < fStartIx {
			return nil, nil
		}

		if fStartIx < caStartIx {
			op.StartID = address.StartID
		}

		if fEndIx > caEndIx {
			op.EndID = LastID
		}

		ops = append(ops, op)
	case MultiOp:
		newMop := MultiOp{}
		for _, op := range op.Mops {
			tops, err := r._trimFilterOp(caStartIx, caEndIx, address, op)
			if err != nil {
				return nil, err
			}
			newMop.Mops = append(newMop.Mops, tops...)
		}
		ops = append(ops, newMop)
	}

	return ops, nil
}

func (r *Rogue) GetOldRogue(address *ContentAddress) (*Rogue, error) {
	newRogue := NewRogue(r.Author)

	vals := MapValues(r.OpIndex.AuthorOps)
	err := avl.Merge(vals, func(op Op) error {
		if !address.Contains(op.GetID()) {
			return ErrorStopIteration{}
		}

		_, err := newRogue.MergeOp(op)
		return err
	})

	if err != nil {
		if !errors.As(err, &ErrorStopIteration{}) {
			return nil, err
		}
	}

	return newRogue, nil
}

func (r *Rogue) GetAddressRogue(address *ContentAddress) (*Rogue, error) {
	// GET NODES FROM ORIGINAL ROGUE
	_, caStartIx, err := r.Rope.GetIndex(address.StartID)
	if err != nil {
		return nil, fmt.Errorf("GetIndex(%v): %w", address.StartID, err)
	}

	_, caEndIx, err := r.Rope.GetIndex(address.EndID)
	if err != nil {
		return nil, fmt.Errorf("GetIndex(%v): %w", address.EndID, err)
	}

	ops := make([]Op, 0, (caEndIx-caStartIx)*2)

	for author, tree := range r.OpIndex.AuthorOps {
		if _, ok := address.MaxIDs[author]; !ok {
			continue
		}

		err := tree.Dft(func(op Op) error {
			iid := op.GetID()

			maxIDSeq := address.MaxIDs[iid.Author]
			if maxIDSeq < iid.Seq {
				return ErrorStopIteration{} // can stop since we know any more are after address
			}

			trimmedOps, err := r._trimFilterOp(caStartIx, caEndIx, *address, op)
			if err != nil {
				return fmt.Errorf("_trimFilterOp(%d, %d, %v, %v): %w", caStartIx, caEndIx, *address, op, err)
			}

			ops = append(ops, trimmedOps...)
			return nil
		})

		if err != nil {
			if errors.As(err, &ErrorStopIteration{}) {
				continue
			}

			return nil, fmt.Errorf("index.Dft(): %w", err)
		}
	}

	/*
		maxSeq := address.MaxID().Seq

		// HANDLE FIRST NODE BOUNDARY
		startNode := r.Index.GetRopeNode(address.StartID)
		if startNode != nil {
			lc := startNode.Val.LeftChildren
			if len(lc) > 0 {
				leftSib := lc[len(lc)-1].rightmost()
				id := ID{leftSib.ID.Author, leftSib.ID.Seq + len(leftSib.Text) - 1}
				text := leftSib.Text[len(leftSib.Text)-1:]

				insertOp := InsertOp{
					ID:       id,
					Text:     Uint16ToStr(text),
					ParentID: startNode.Val.ID,
					Side:     Left,
				}
				fmt.Printf("insertOp: %v\n", insertOp)
				ops = append(ops, insertOp)

				seq := max(maxSeq+1, id.Seq+1)
				deleteOp := DeleteOp{
					ID:         ID{Author: "t", Seq: seq},
					TargetID:   id,
					SpanLength: 1,
				}
				ops = append(ops, deleteOp)
			}
		}

		// HANDLE LAST NODE BOUNDARY
		endNode := r.Index.GetRopeNode(address.EndID)
		if endNode != nil {
			rc := endNode.Val.RightChildren
			if len(rc) > 0 {
				rightSib := rc[0].leftmost()
				text := rightSib.Text[:1]

				if IsHighSurrogate(text[0]) {
					text = rightSib.Text[:2]
				}

				insertOp := InsertOp{
					ID:       rightSib.ID,
					Text:     Uint16ToStr(text),
					ParentID: endNode.Val.ID,
					Side:     Right,
				}
				fmt.Printf("insertOp: %v\n", insertOp)
				ops = append(ops, insertOp)

				seq := max(maxSeq+2, rightSib.ID.Seq+1)
				deleteOp := DeleteOp{
					ID:         ID{Author: "j", Seq: seq},
					TargetID:   rightSib.ID,
					SpanLength: 1,
				}
				ops = append(ops, deleteOp)
			}
		}
	*/

	// SORT NODES BY SEQ FOR BREADTH FIRST INSERTION
	slices.SortFunc(ops, func(a, b Op) int {
		return a.GetID().Seq - b.GetID().Seq
	})

	// Populate new rogue
	addressRogue := NewRogue(r.Author)
	for _, op := range ops {
		_, err := addressRogue.MergeOp(op)
		if err != nil {
			return nil, fmt.Errorf("MergeOp(%v): %w", op, err)
		}
	}

	// ADD FORMAT OPS
	/*formatOps, err := r.Formats.SearchOverlapping(address.StartID, address.EndID)
	if err != nil {
		return nil, fmt.Errorf("SearchOverlapping(%v, %v): %w", address.StartID, address.EndID, err)
	}

	for _, f := range formatOps {
		if maxSeq, ok := address.MaxIDs[f.ID.Author]; ok {
			if f.ID.Seq <= maxSeq {
				_, fStartIx, err := r.Rope.GetIndex(f.StartID)
				if err != nil {
					return nil, fmt.Errorf("GetIndex(%v): %w", f.StartID, err)
				}

				if fStartIx < caStartIx {
					f.StartID = address.StartID
				}

				_, fEndIx, err := r.Rope.GetIndex(f.EndID)
				if err != nil {
					return nil, fmt.Errorf("GetIndex(%v): %w", f.EndID, err)
				}

				if fEndIx > caEndIx {
					f.EndID = LastID
				}

				_, err = addressRogue.MergeOp(f)
				if err != nil {
					return nil, fmt.Errorf("Insert(%v): %w", f, err)
				}
			}
		}
	}*/

	return addressRogue, nil
}

func (ca *ContentAddress) MaxID() *ID {
	var maxID *ID

	for author, seq := range ca.MaxIDs {
		if maxID == nil || seq > maxID.Seq || (seq == maxID.Seq && author > maxID.Author) {
			maxID = &ID{author, seq}
		}
	}

	return maxID
}

func (ca *ContentAddress) MinID() *ID {
	var minID *ID

	for author, seq := range ca.MaxIDs {
		if minID == nil || seq < minID.Seq || (seq == minID.Seq && author < minID.Author) {
			minID = &ID{author, seq}
		}
	}

	return minID
}

func idCompFunc(reverse bool) func(ID, ID) int {
	return func(a, b ID) int {
		x := 0
		if a.Seq < b.Seq {
			x = -1
		} else if a.Seq > b.Seq {
			x = 1
		} else if a.Author < b.Author {
			x = -1
		} else if a.Author > b.Author {
			x = 1
		} else {
			x = 0
		}

		if reverse {
			return -x
		} else {
			return x
		}
	}
}

func (ca *ContentAddress) SortedIDs(reverse bool) []ID {
	ids := make([]ID, 0, len(ca.MaxIDs))

	for author, seq := range ca.MaxIDs {
		ids = append(ids, ID{author, seq})
	}

	slices.SortFunc(ids, idCompFunc(reverse))

	return ids
}

func (ca *ContentAddress) IsMin() bool {
	if len(ca.MaxIDs) == 0 {
		return true
	}

	for _, seq := range ca.MaxIDs {
		if seq != -1 {
			return false
		}
	}

	return true
}

func (ca *ContentAddress) IsMax() bool {
	if len(ca.MaxIDs) == 0 {
		return true
	}

	for _, seq := range ca.MaxIDs {
		if seq != math.MaxInt {
			return false
		}
	}

	return true
}

func (r *Rogue) ValidAddress(address ContentAddress) bool {
	if !r.ContainsID(address.StartID) {
		return false
	}

	if !r.ContainsID(address.EndID) {
		return false
	}

	for author, maxSeq := range address.MaxIDs {
		if !r.ContainsID(ID{Author: author, Seq: maxSeq}) {
			return false
		}
	}

	return true
}

func (ca *ContentAddress) SelectAuthors(authors ...string) *ContentAddress {
	newCA := ca.DeepCopy()
	newMaxIDs := map[string]int{}

	for _, author := range authors {
		if maxSeq, ok := ca.MaxIDs[author]; ok {
			newMaxIDs[author] = maxSeq
		}
	}

	newCA.MaxIDs = newMaxIDs
	return newCA
}

func (ca *ContentAddress) Contains(id ID) bool {
	if maxSeq, ok := ca.MaxIDs[id.Author]; ok {
		return id.Seq <= maxSeq
	}

	return false
}

func (r *Rogue) Compact(address *ContentAddress) (*Rogue, error) {
	rogue := NewRogue(r.Author)

	firstID, err := r.GetFirstTotID()
	if err != nil {
		return nil, err
	}

	lastID, err := r.GetLastTotID()
	if err != nil {
		return nil, err
	}

	vis, span, line, err := r.ToIndexNos(firstID, lastID, address, false)
	if err != nil {
		return nil, err
	}

	sText := Uint16ToStr(vis.Text)
	_, err = rogue.Insert(0, sText)
	if err != nil {
		return nil, err
	}

	err = span.tree.Dft(func(node *NOSNode) error {
		if node.EndIx >= len(vis.Text) {
			return nil
		}

		_, err := rogue.Format(node.StartIx, node.EndIx-node.StartIx+1, node.Format)
		return err
	})
	if err != nil {
		return nil, err
	}

	err = line.tree.Dft(func(node *NOSNode) error {
		_, err := rogue.Format(node.StartIx, 1, node.Format)
		return err
	})
	if err != nil {
		return nil, err
	}

	return rogue, nil
}
