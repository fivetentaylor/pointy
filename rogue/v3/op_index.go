package v3

import (
	"fmt"
	"strings"

	avl "github.com/fivetentaylor/pointy/rogue/v3/tree"
)

type OpIndex struct {
	AuthorOps      map[string]*avl.Tree[int, Op]
	ContentAddress *ContentAddress
}

func NewOpIndex() OpIndex {
	return OpIndex{
		AuthorOps:      make(map[string]*avl.Tree[int, Op]),
		ContentAddress: NewContentAddress(),
	}
}

func (oi OpIndex) Size() int {
	size := 0
	for _, tree := range oi.AuthorOps {
		size += tree.Size
	}
	return size
}

func (oi OpIndex) String() string {
	builder := strings.Builder{}

	for auth, tree := range oi.AuthorOps {
		builder.WriteString(fmt.Sprintf("Author: %s\n", auth))
		tree.Dft(func(op Op) error {

			builder.WriteString(fmt.Sprintf("  %s\n", op))
			return nil
		})
	}

	return builder.String()
}

func opIdSeq(op Op) (int, error) {
	return op.GetID().Seq, nil
}

func (index *OpIndex) Put(op Op) {
	author := op.GetID().Author
	tree, exists := index.AuthorOps[author]
	if !exists {
		tree = avl.NewTree(opIdSeq, lessThan)
		index.AuthorOps[author] = tree
	}

	// Clean up if these ops were already inserted previously
	if mop, ok := op.(MultiOp); ok {
		for _, op := range mop.Mops {
			tree.Remove(op)
		}
	}

	tree.Put(op)

	// Keep track of current content address
	maxID := MaxID(op)
	MapSetMax(index.ContentAddress.MaxIDs, maxID.Author, maxID.Seq)
}

func (index *OpIndex) Remove(op Op) {
	author := op.GetID().Author
	tree, exists := index.AuthorOps[author]
	if !exists {
		return
	}

	if mop, ok := op.(MultiOp); ok {
		for _, op := range mop.Mops {
			tree.Remove(op)
		}
	} else {
		tree.Remove(op)
	}
}

func (index *OpIndex) GetExact(id ID) Op {
	nodes, exists := index.AuthorOps[id.Author]
	if !exists {
		return nil
	}

	op, _ := nodes.FindLeftSib(id.Seq)
	if op == nil {
		return nil
	}

	if op.GetID() == id {
		return op
	}

	return nil
}

func (index *OpIndex) Get(id ID) Op {
	nodes, exists := index.AuthorOps[id.Author]
	if !exists {
		return nil
	}

	op, _ := nodes.FindLeftSib(id.Seq)
	if op == nil {
		return nil
	}

	if op.GetID() == id {
		return op
	}

	if id.Seq <= MaxID(op).Seq {
		return op
	}

	return nil
}

func (index *OpIndex) MaxSeq() int {
	max := 0
	for _, seq := range index.ContentAddress.MaxIDs {
		if seq > max {
			max = seq
		}
	}
	return max
}

func (index *OpIndex) nextSmallest(id ID) *ID {
	if id.Seq <= 0 {
		return nil
	}

	tree, exists := index.AuthorOps[id.Author]
	if !exists {
		return nil
	}

	op, _ := tree.FindLeftSib(id.Seq - 1)
	if op == nil {
		return nil
	}

	nextID := op.GetID()
	return &nextID
}
