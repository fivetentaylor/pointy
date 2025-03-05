package v3

import avl "github.com/teamreviso/code/rogue/v3/tree"

type FailedOps struct {
	*avl.Tree[ID, Op]
}

func NewFailedOps() *FailedOps {
	tree := avl.NewTree(opGetID, lessThanRevID)
	return &FailedOps{tree}
}

func (fo *FailedOps) Put(op Op) error {
	if op, ok := op.(MultiOp); ok {
		for _, op := range op.Mops {
			fo.Tree.Remove(op)
		}
	}

	return fo.Tree.Put(op)
}

func (fo *FailedOps) Size() int {
	return fo.Tree.Size
}
