package v3

import (
	"fmt"

	avl "github.com/teamreviso/code/rogue/v3/tree"
)

type LineHistory struct {
	TargetID ID
	Formats  *avl.Tree[ID, FormatOp]
}

func NewLineHistory(targetID ID) *LineHistory {
	keyFunc := func(op FormatOp) (ID, error) {
		return op.ID, nil
	}

	return &LineHistory{
		TargetID: targetID,
		Formats:  avl.NewTree(keyFunc, lessThanRevID),
	}
}

type Lines struct {
	Rope *Rope
	Tree *avl.Tree[int, *LineHistory]
}

func NewLines(rope *Rope) *Lines {
	keyFunc := func(lh *LineHistory) (int, error) {
		_, ix, err := rope.GetIndex(lh.TargetID)
		if err != nil {
			return 0, err
		}

		return ix, nil
	}

	return &Lines{
		Rope: rope,
		Tree: avl.NewTree(keyFunc, lessThan),
	}
}

func (l *Lines) Insert(op FormatOp) error {
	if _, ok := op.Format.(FormatV3Span); ok {
		return fmt.Errorf("formatOp %v is not a line format", op)
	}

	_, lineIx, err := l.Rope.GetIndex(op.StartID)
	if err != nil {
		return fmt.Errorf("GetIndex(%v): %w", op.StartID, err)
	}

	lh, err := l.Tree.Get(lineIx)
	if err != nil {
		return fmt.Errorf("Get(%v): %w", lineIx, err)
	}

	if lh == nil {
		lh = NewLineHistory(op.StartID)
		l.Tree.Put(lh)
	}

	lh.Formats.Put(op)
	return nil
}

func (l *Lines) Validate() error {
	err := l.Tree.Dft(func(lh *LineHistory) error {
		err := lh.Formats.Dft(func(op FormatOp) error {
			if _, ok := op.Format.(FormatV3Span); ok {
				return fmt.Errorf("formatOp %v is not a line format", op)
			}

			if op.StartID != lh.TargetID {
				return fmt.Errorf("formatOp %v does not belong to line %v", op, lh.TargetID)
			}

			if op.StartID != op.EndID {
				return fmt.Errorf("line formatOp %v is has different start and end ID", op)
			}

			return nil
		})
		return err
	})
	return err
}

func (l *Lines) FormatAt(id ID, address *ContentAddress) (*FormatOp, error) {
	_, lineIx, err := l.Rope.GetIndex(id)
	if err != nil {
		return nil, err
	}

	lh, err := l.Tree.Get(lineIx)
	if err != nil {
		return nil, err
	}

	if lh == nil {
		return nil, nil
	}

	if address == nil {
		fop := lh.Formats.Max()
		return &fop, nil
	}

	_, aStartIx, err := l.Rope.GetIndex(address.StartID)
	if err != nil {
		return nil, err
	}

	_, aEndIx, err := l.Rope.GetIndex(address.EndID)
	if err != nil {
		return nil, err
	}

	if lineIx < aStartIx || aEndIx < lineIx {
		fop := lh.Formats.Max()
		return &fop, nil
	}

	n, err := lh.Formats.FindLeftSibNode(*address.MaxID())
	if err != nil {
		return nil, err
	}

	for {
		if n == nil {
			return nil, nil
		}

		if address.Contains(n.Value.ID) {
			return &n.Value, nil
		}

		n = n.StepLeft()
	}
}
