package v3

import (
	"fmt"

	"github.com/teamreviso/code/pkg/stackerr"
)

func (r *Rogue) UndoDoc() (Op, error) {
	firstID, err := r.GetFirstTotID()
	if err != nil {
		return nil, err
	}

	lastID, err := r.GetLastTotID()
	if err != nil {
		return nil, err
	}

	return r.Undo(firstID, lastID)
}

func (r *Rogue) Undo(startID, endID ID) (Op, error) {
	if !r.CanUndo() {
		return nil, nil
	}

	if r.UndoState == nil || r.UndoState.address == nil {
		ca, err := r.GetFullAddress()
		if err != nil {
			return nil, fmt.Errorf("GetFullAddress(): %w", err)
		}

		// Just undo all authors for now
		// r.UndoAddress = ca.SelectAuthors(r.Author)
		if r.UndoState == nil {
			r.UndoState = &UndoState{address: ca}
		} else {
			r.UndoState.address = ca
			r.UndoState.redoStack = nil
		}
	}

	uop, rop, err := r._undoNext(startID, endID, r.UndoState.address)
	if err != nil {
		return nil, fmt.Errorf("_undoNext(%v, %v, %v): %w", startID, endID, r.UndoState.address, err)
	}

	_, err = r.MergeOp(uop)
	if err != nil {
		return nil, fmt.Errorf("MergeOp(%v): %w", uop, err)
	}

	r.UndoState.redoStack = append(r.UndoState.redoStack, rop)

	return uop, nil
}

func (r *Rogue) Redo() (Op, error) {
	mop := MultiOp{}

	if r.UndoState == nil {
		return mop, nil
	}

	if len(r.UndoState.redoStack) == 0 {
		return mop, nil
	}

	// clear the content address so we undo from the tip again
	r.UndoState.address = nil

	rs := r.UndoState.redoStack
	item := rs[len(rs)-1]
	r.UndoState.redoStack = rs[:len(rs)-1]

	var ops []Op
	if m, ok := item.(MultiOp); ok {
		ops = m.Mops
	} else {
		ops = []Op{item}
	}

	for i := len(ops) - 1; i >= 0; i-- {
		op := ops[i]
		switch op := op.(type) {
		case DeleteOp:
			op.ID = r.NextID(1)
			mop = mop.Append(op)
		case ShowOp:
			op.ID = r.NextID(1)
			mop = mop.Append(op)
		case FormatOp:
			op.ID = r.NextID(1)
			mop = mop.Append(op)
		case RewindOp:
			op = RewindOp{ID: r.NextID(1), Address: op.UndoAddress, UndoAddress: op.Address}
			mop = mop.Append(op)
		default:
			return mop, fmt.Errorf("not an undo op type: %T", item)
		}
	}

	op := FlattenMop(mop)
	_, err := r.MergeOp(op)
	if err != nil {
		return nil, fmt.Errorf("MergeOp(%v): %w", op, err)
	}

	return op, nil
}

// TODO: fix up rewind implementation!
func (r *Rogue) Rewind(startID, endID ID, address ContentAddress) (mop MultiOp, err error) {
	r.UndoState = nil

	_, startIx, err := r.Rope.GetIndex(startID)
	if err != nil {
		return mop, fmt.Errorf("GetIndex(%v): %w", startID, err)
	}

	_, endIx, err := r.Rope.GetIndex(endID)
	if err != nil {
		return mop, fmt.Errorf("GetIndex(%v): %w", endID, err)
	}

	_, caStartIx, err := r.Rope.GetIndex(address.StartID)
	if err != nil {
		return mop, fmt.Errorf("GetIndex(%v): %w", address.StartID, err)
	}

	_, caEndIx, err := r.Rope.GetIndex(address.EndID)
	if err != nil {
		return mop, fmt.Errorf("GetIndex(%v): %w", address.EndID, err)
	}

	if startIx < caStartIx {
		startID = address.StartID
	} else {
		address.StartID = startID
	}

	if caEndIx < endIx {
		endID = address.EndID
	} else {
		address.EndID = endID
	}

	undoAddress, err := r.GetAddress(startID, endID)
	if err != nil {
		return mop, fmt.Errorf("GetAddress(%v, %v): %w", startID, endID, err)
	}

	mop = mop.Append(RewindOp{
		ID:          r.NextID(1),
		Address:     address,
		UndoAddress: *undoAddress,
	})

	// TODO: implement actions
	_, err = r.MergeOp(mop)
	if err != nil {
		return mop, fmt.Errorf("MergeOp(%v): %w", mop, err)
	}

	return mop, nil
}

// TODO: add StartID and endID to RewindOp, remove from address
func (r *Rogue) rewindOp(op RewindOp) (err error) {
	node := r.RopeIndex.Get(op.Address.StartID)
	if node == nil {
		return stackerr.New(fmt.Errorf("node for id %v doesn't exist", op.Address.StartID))
	}

	startID, endID := op.Address.StartID, op.Address.EndID

	fops, err := r.rewindFormatTo(op.Address.StartID, op.Address.EndID, &op.Address, op.ID)
	if err != nil {
		return fmt.Errorf("rewindFormatToV2(%v, true): %w", op, err)
	}

	for _, fop := range fops {
		_, err := r.FormatOp(fop)
		if err != nil {
			return fmt.Errorf("FormatOp(%v): %w", fop, err)
		}
	}

	for {
		isStart := node.Val.ContainsID(startID)
		isEnd := node.Val.ContainsID(endID)

		fn := node.Val.Explode()
		for _, targetID := range fn.IDs {
			// continue if we haven't reached the startID
			if isStart && targetID.Seq < startID.Seq {
				continue
			}

			// BEGIN DO STUFF HERE
			marker := &Marker{ID: op.ID, IsDel: false}

			// Rewind before the creation of the character
			if !op.Address.Contains(targetID) {
				marker.IsDel = true
			} else {
				ch, ok := r.CharHistory[targetID]
				if ok {
					chNode := ch.MaxNode()
					for chNode != nil {
						if op.Address.Contains(chNode.Value.ID) {
							marker.IsDel = chNode.Value.IsDel
							break
						}
						chNode = chNode.StepLeft()
					}
				}
			}

			r.CharHistory.Add(targetID, marker)
			maxMarker := r.CharHistory.Max(targetID)
			err := r.MarkCharDel(targetID, maxMarker.IsDel)
			if err != nil {
				return fmt.Errorf("MarkCharDel(%v, %v): %w", targetID, maxMarker.IsDel, err)
			}
			// END DO STUFF HERE

			// exit if we reached the end id
			if isEnd && targetID.Seq >= endID.Seq {
				r.LamportClock = max(r.LamportClock, op.ID.Seq+1)
				return nil
			}
		}

		node, err = node.RightTotSibling()
		if err != nil {
			return fmt.Errorf("rightTotSibling(): %w", err)
		}
	}
}
