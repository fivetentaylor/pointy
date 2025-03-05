package v3

import (
	avl "github.com/teamreviso/code/rogue/v3/tree"
)

type RopeIndex map[string]*avl.Tree[int, *RopeNode]

func NewRopeIndex() RopeIndex {
	return make(RopeIndex)
}

func idSeq(node *RopeNode) (int, error) {
	return node.Val.ID.Seq, nil
}

func (index RopeIndex) Put(node *RopeNode) {
	author := node.Val.ID.Author
	nodes, exists := index[author]
	if !exists {
		nodes = avl.NewTree(idSeq, lessThan)
		index[author] = nodes
	}

	nodes.Put(node)
}

func (index RopeIndex) Get(id ID) *RopeNode {
	nodes, exists := index[id.Author]
	if !exists {
		return nil
	}

	n, _ := nodes.FindLeftSib(id.Seq)
	if n == nil {
		return nil
	}

	startSeq := n.Val.ID.Seq
	endSeq := startSeq + len(n.Val.Text)
	if id.Seq >= startSeq && id.Seq < endSeq {
		return n
	}

	return nil
}
