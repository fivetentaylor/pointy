package v3

import (
	"errors"

	"github.com/fivetentaylor/pointy/pkg/stackerr"
	avl "github.com/fivetentaylor/pointy/rogue/v3/tree"
)

type ScrubState struct {
	FullAddress *ContentAddress // if full doc scrub

	StartID ID
	EndID   ID
	IdTree  *avl.Tree[ID, ID]
	CurNode *avl.Node[ID, ID]

	CurIx      int
	CurAddress *ContentAddress
}

func (r *Rogue) ScrubInit(startID, endID *ID) (int, error) {
	r.UndoState = nil

	if startID == nil && endID == nil {
		address, err := r.GetFullAddress()
		if err != nil {
			return -1, err
		}

		size := r.OpIndex.Size() - 1
		idTree := avl.NewTree(
			func(id ID) (ID, error) { return id, nil },
			lessThanRevID,
		)

		for _, tree := range r.OpIndex.AuthorOps {
			tree.Dft(func(op Op) error {
				idTree.Put(op.GetID())
				return nil
			})
		}

		firstID, err := r.GetFirstTotID()
		if err != nil {
			return -1, err
		}

		lastID, err := r.GetLastTotID()
		if err != nil {
			return -1, err
		}

		r.ScrubState = &ScrubState{
			FullAddress: address,
			StartID:     firstID,
			EndID:       lastID,
			CurIx:       size,
			CurAddress:  address,
			IdTree:      idTree,
			CurNode:     idTree.MaxNode(),
		}

		return size, nil
	}

	var err error
	if startID == nil {
		bid, err := r.GetFirstTotID()
		if err != nil {
			return -1, err
		}
		startID = &bid
	}

	if endID == nil {
		lid, err := r.GetLastTotID()
		if err != nil {
			return -1, err
		}
		endID = &lid
	}

	idTree, address, err := r._computeAddressState(*startID, *endID)
	if err != nil {
		return -1, err
	}

	r.ScrubState = &ScrubState{
		StartID: *startID,
		EndID:   *endID,
		IdTree:  idTree,
		CurNode: idTree.MaxNode(),

		CurAddress: address,
		CurIx:      idTree.Size - 1,
	}

	return idTree.Size - 1, nil
}

type ScrubStep struct {
	Span          *RenderSpan
	CursorStartID ID
	CursorEndID   ID
	Html          string
}

func (r *Rogue) ScrubTo(ix int) (*ScrubStep, error) {
	if r.ScrubState == nil {
		return nil, stackerr.Errorf("ScrubInit must be called before ScrubTo")
	}

	if ix == r.ScrubState.CurIx {
		return nil, nil
	}

	prevAddress := r.ScrubState.CurAddress.DeepCopy()

	if r.ScrubState.FullAddress != nil {
		if ix < 0 || ix >= r.OpIndex.Size() {
			return nil, stackerr.Errorf("Invalid index %d", ix)
		}

		totStartIx, totEndIx := r.TotSize-1, 0
		var startID, endID ID

		for ix != r.ScrubState.CurIx {
			var op Op

			if r.ScrubState.CurIx < ix {
				node := r.ScrubState.CurNode.StepRight()
				id := node.Value
				op = r.OpIndex.Get(id)
				r.ScrubState.CurAddress.MaxIDs[id.Author] = MaxID(op).Seq
				r.ScrubState.CurNode = node
				r.ScrubState.CurIx++
			} else if r.ScrubState.CurIx > ix {
				node := r.ScrubState.CurNode
				id := node.Value
				op = r.OpIndex.Get(id)
				r.ScrubState.CurAddress.MaxIDs[id.Author] = id.Seq - 1
				r.ScrubState.CurNode = node.StepLeft()
				r.ScrubState.CurIx--
			}

			opSpan, err := r.opSpan(op)
			if err != nil {
				return nil, err
			}

			if opSpan.TotStartIx < totStartIx {
				totStartIx = opSpan.TotStartIx
				startID = opSpan.StartID
			}

			if opSpan.TotEndIx > totEndIx {
				totEndIx = opSpan.TotEndIx
				endID = opSpan.EndID
			}
		}

		span, err := r.RenderSpanBetween(startID, endID, prevAddress, r.ScrubState.CurAddress)
		if err != nil {
			return nil, err
		}

		_, lastID, err := r.CursorAt(startID, endID, r.ScrubState.CurAddress)
		if err != nil {
			return nil, err
		}

		return &ScrubStep{
			Span:          span,
			CursorStartID: lastID,
			CursorEndID:   lastID,
			Html:          span.Html,
		}, nil

	} else if r.ScrubState.IdTree != nil {
		if ix < 0 || ix >= r.ScrubState.IdTree.Size {
			return nil, stackerr.Errorf("Invalid index %d", ix)
		}

		node := r.ScrubState.CurNode

		for ix != r.ScrubState.CurIx {
			if node == nil {
				return nil, stackerr.Errorf("Node should never be nil here")
			}

			if r.ScrubState.CurIx < ix {
				node = r.ScrubState.CurNode.StepRight()
				id := node.Value
				r.ScrubState.CurNode = node
				r.ScrubState.CurAddress.MaxIDs[id.Author] = id.Seq
				r.ScrubState.CurIx++
			} else if r.ScrubState.CurIx > ix {
				node = r.ScrubState.CurNode
				id := node.Value
				r.ScrubState.CurNode = node.StepLeft()
				r.ScrubState.CurAddress.MaxIDs[id.Author] = id.Seq - 1
				r.ScrubState.CurIx--
			}

		}

		span, err := r.RenderSpanBetween(r.ScrubState.StartID, r.ScrubState.EndID, prevAddress, r.ScrubState.CurAddress)
		if err != nil {
			return nil, err
		}

		firstID, lastID, err := r.CursorAt(r.ScrubState.StartID, r.ScrubState.EndID, r.ScrubState.CurAddress)
		if err != nil {
			return nil, err
		}

		return &ScrubStep{
			Span:          span,
			CursorStartID: firstID,
			CursorEndID:   lastID,
			Html:          span.Html,
		}, nil
	} else {
		return nil, stackerr.Errorf("Invalid ScrubState, must have FullAddress or IdTree")
	}
}

func (r *Rogue) _computeAddressState(startID, endID ID) (*avl.Tree[ID, ID], *ContentAddress, error) {
	address := &ContentAddress{
		MaxIDs:  make(map[string]int),
		StartID: startID,
		EndID:   endID,
	}

	idTree := avl.NewTree(
		func(id ID) (ID, error) { return id, nil },
		lessThanRevID,
	)

	var err error
	ops, err := r.Formats.SearchOverlapping(startID, endID)
	if err != nil {
		return nil, nil, err
	}

	for _, op := range ops {
		id := op.GetID()

		if err := idTree.Put(id); err != nil {
			return nil, nil, err
		}

		address.AddID(id)
	}

	for item, err := range r.WalkRightFromTot(startID) {
		if err != nil {
			return nil, nil, err
		}

		if err := idTree.Put(item.ID); err != nil {
			return nil, nil, err
		}

		address.AddID(item.ID)

		if ch, ok := r.CharHistory[item.ID]; ok {
			err = ch.Dft(func(m *Marker) error {
				if err := idTree.Put(m.ID); err != nil {
					return err
				}

				address.AddID(m.ID)

				return nil
			})
			if err != nil {
				return nil, nil, err
			}
		}

		if item.ID == endID {
			break
		}
	}

	return idTree, address, nil
}

func (r *Rogue) CursorAt(startID, endID ID, address *ContentAddress) (firstID, lastID ID, err error) {
	_, endIx, err := r.Rope.GetIndex(endID)
	if err != nil {
		return NoID, NoID, err
	}

	firstID, lastID = startID, endID
	for item, err := range r.WalkRightFromAt(startID, address) {
		if err != nil {
			for item, err := range r.WalkLeftFromAt(startID, address) {
				if err != nil {
					return NoID, NoID, err
				}

				return item.ID, item.ID, nil
			}
		}

		firstID = item.ID
		break
	}

	_, firstIx, err := r.Rope.GetIndex(firstID)
	if err != nil {
		return NoID, NoID, err
	}

	if endIx < firstIx {
		return firstID, firstID, nil // Cursort at the first right char
	}

	// There's something visible between startID and endID
	for item, err := range r.WalkLeftFromAt(endID, address) {
		if err != nil {
			return NoID, NoID, err
		}

		lastID = item.ID
		break
	}

	lid, err := r.RightOfAt(lastID, address)
	if err != nil {
		if !errors.As(err, &ErrorNoRightSiblingAt{}) {
			return NoID, NoID, err
		}

		return firstID, lastID, nil
	}

	return firstID, lid, nil
}
