package v3

import (
	"errors"
	"fmt"

	rbt "github.com/emirpasic/gods/trees/redblacktree"
	aSet "github.com/teamreviso/code/rogue/v3/set"
	avl "github.com/teamreviso/code/rogue/v3/tree"
)

type NOS struct {
	tree *avl.Tree[int, *NOSNode]
}

type NOSNode struct {
	StartIx int
	EndIx   int
	Format  FormatV3
}

func (n NOSNode) String() string {
	return fmt.Sprintf("StartIx: %d EndIx: %d Format: %v", n.StartIx, n.EndIx, n.Format)
}

func (nos *NOS) String() string {
	out := ""
	nos.tree.Dft(func(node *NOSNode) error {
		out += fmt.Sprintf("%s\n", node)
		return nil
	})

	return out
}

func (nos *NOS) AsSlice() []*NOSNode {
	return nos.tree.AsSlice()
}

func NewNOS() *NOS {
	keyFunc := func(n *NOSNode) (int, error) {
		return n.StartIx, nil
	}

	return &NOS{
		tree: avl.NewTree(keyFunc, lessThan),
	}
}

// StringToLineNOS returns the start and end indices of
// each line in the string as a NOS tree. This makes it easy
// to find the start or end of a line for a given index
func StringToLineNOS(s []uint16) *NOS {
	nos := NewNOS()

	prevIx := 0
	for i, c := range s {
		if c == '\n' {
			nos.tree.Put(&NOSNode{prevIx, i, FormatV3Line{}})
			prevIx = i + 1
		}
	}

	if prevIx < len(s) {
		nos.tree.Put(&NOSNode{prevIx, len(s), FormatV3Line{}})
	}

	return nos
}

func mergeFormats(curFmt, newFmt FormatV3) FormatV3 {
	if !newFmt.IsSpan() {
		return newFmt
	}

	if !curFmt.IsSpan() {
		return curFmt
	}

	curSticky, curNoSticky := curFmt.(FormatV3Span).SplitSticky()
	newSticky, newNoSticky := newFmt.(FormatV3Span).SplitSticky()

	var merged FormatV3Span

	// "e" for exclusive
	if _, ok := newSticky["e"]; ok {
		merged = newSticky
	} else {
		merged = MergeMaps(curSticky, newSticky)
	}

	// "en" for no exclusive nonsticky
	if _, ok := newNoSticky["en"]; ok {
		merged = MergeMaps(merged, newNoSticky)
	} else {
		merged = MergeMaps(merged, MergeMaps(curNoSticky, newNoSticky))
	}

	return merged
}

func (nos NOSNode) IsLineFormat() bool {
	if !nos.Format.IsSpan() {
		return nos.StartIx == nos.EndIx
	}

	return false
}

func mergeNodes(cur, new NOSNode) (toInsert, toDelete []NOSNode, remaining *NOSNode, err error) {
	if cur.EndIx < new.StartIx {
		new.Format = new.Format.DropNull()
		if !new.Format.Empty() {
			return nil, nil, &new, nil
		} else {
			return nil, nil, nil, nil
		}
	}

	if new.EndIx < cur.StartIx {
		new.Format = new.Format.DropNull()
		if !new.Format.Empty() {
			return []NOSNode{new}, nil, nil, nil
		} else {
			return nil, nil, nil, nil
		}
	}

	toDelete = []NOSNode{cur}
	toInsert = []NOSNode{}
	if cur.StartIx < new.StartIx {
		f := cur.Format.DropNull()
		if !f.Empty() {
			n := NOSNode{cur.StartIx, new.StartIx - 1, f}
			toInsert = append(toInsert, n)
		}
	}

	if new.StartIx < cur.StartIx {
		f := new.Format.DropNull()
		if !f.Empty() {
			n := NOSNode{new.StartIx, cur.StartIx - 1, f}
			toInsert = append(toInsert, n)
		}
	}

	// handle the overlapping section
	mStart := max(cur.StartIx, new.StartIx)
	mEnd := min(cur.EndIx, new.EndIx)
	// call drop null after the merge so deleted formats get cleaned up
	merged := mergeFormats(cur.Format, new.Format).DropNull()
	if !merged.Empty() {
		n := NOSNode{mStart, mEnd, merged}
		toInsert = append(toInsert, n)
	}

	if cur.EndIx < new.EndIx {
		f := new.Format.DropNull()
		if !f.Empty() {
			n := NOSNode{cur.EndIx + 1, new.EndIx, f}
			remaining = &n
		}
	}

	if new.EndIx < cur.EndIx {
		f := cur.Format.DropNull()
		if !f.Empty() {
			n := NOSNode{new.EndIx + 1, cur.EndIx, f}
			toInsert = append(toInsert, n)
		}
	}

	// merge adjacent nodes with the same format
	for i := len(toInsert) - 1; i > 0; i-- {
		a, b := toInsert[i-1], toInsert[i]

		b.Format = b.Format.DropNull()

		if b.Format.Empty() {
			toInsert = DeleteAt(toInsert, i, 1)
			continue
		}

		a.Format = a.Format.DropNull()

		yes := a.EndIx == b.StartIx-1 && a.Format.Equals(b.Format)

		if yes {
			toInsert = DeleteAt(toInsert, i-1, 2)
			toInsert = InsertAt(toInsert, i-1, NOSNode{
				StartIx: a.StartIx,
				EndIx:   b.EndIx,
				Format:  a.Format,
			})
		}
	}

	return toInsert, toDelete, remaining, nil
}

// Insert a new node into a nos tree, merging formats
// to represent the latest state of formatting
func (nos *NOS) Insert(newNode NOSNode) error {
	startNode, err := nos.tree.FindLeftSibNode(newNode.StartIx)
	if err != nil {
		return err
	}

	if startNode == nil {
		startNode, err = nos.tree.FindRightSibNode(newNode.StartIx)
		if err != nil {
			return err
		}
	}

	if startNode == nil {
		nos.tree.Put(&newNode)
		return nil
	}

	var toInsert []NOSNode
	var toDelete []NOSNode
	// var newNodes []NOSNode
	var remaining *NOSNode
	err = startNode.WalkRight(func(n *NOSNode) error {
		var insert, delete []NOSNode
		insert, delete, remaining, err = mergeNodes(*n, newNode)
		if err != nil {
			return fmt.Errorf("mergeNodes(%v, %v): %w", n, newNode, err)
		}

		toInsert = append(toInsert, insert...)
		toDelete = append(toDelete, delete...)

		if remaining != nil {
			newNode = *remaining
		} else {
			return ErrorStopIteration{}
		}

		return nil
	})

	if err != nil {
		if !errors.As(err, &ErrorStopIteration{}) {
			return err
		}
	}

	for _, n := range toDelete {
		nos.tree.RemoveByKey(n.StartIx)
	}

	for _, n := range toInsert {
		x := n
		nos.tree.Put(&x)
	}

	if remaining != nil {
		nos.tree.Put(remaining)
	}

	return nil
}

func (nos *NOS) between(startIx, endIx int, callback func(*NOSNode) error) error {
	// right sib is inclusive if the index matches perfectly
	startNode, err := nos.tree.FindRightSibNode(startIx)
	if err != nil {
		return fmt.Errorf("findRightSib(%v): %w", startIx, err)
	}

	err = startNode.WalkRight(func(n *NOSNode) error {
		if n.StartIx > endIx {
			return ErrorStopIteration{}
		}

		if n.EndIx >= startIx {
			return callback(n)
		}

		return nil
	})

	if err != nil {
		if errors.As(err, &ErrorStopIteration{}) {
			return nil
		}

		return err
	}

	return nil
}

type NOSActions struct {
	Tree *rbt.Tree
}

func NewNOSActions() *NOSActions {
	return &NOSActions{
		Tree: rbt.NewWithIntComparator(),
	}
}

func (na *NOSActions) Insert(ix int, s string) {
	list, ok := na.Tree.Get(ix)
	if ok {
		// na.Tree.Put(ix, InsertAt(list.([]string), 0, s))
		na.Tree.Put(ix, append(list.([]string), s))
	} else {
		na.Tree.Put(ix, []string{s})
	}
}

func (na *NOSActions) Merge(b *NOSActions) {
	it := b.Tree.Iterator()
	for it.Next() {
		ix := it.Key().(int)
		sym := it.Value().([]string)

		if c, ok := na.Tree.Get(ix); ok {
			na.Tree.Put(ix, append(c.([]string), sym...))
		} else {
			na.Tree.Put(ix, sym)
		}
	}
}

type PendingActions struct {
	Tree *rbt.Tree
}

func NewPendingActions() *PendingActions {
	return &PendingActions{
		Tree: rbt.NewWithIntComparator(),
	}
}

func (p *PendingActions) Add(ix int, s string) {
	set, ok := p.Tree.Get(ix)
	if ok {
		set.(aSet.Set[string]).Add(s)
	} else {
		p.Tree.Put(ix, aSet.NewSet(s))
	}
}

func (p *PendingActions) Remove(ix int, s string) {
	set, ok := p.Tree.Get(ix)
	if ok {
		set.(aSet.Set[string]).Pop(s)
	}
}

func (p *PendingActions) Symbols() (out aSet.Set[string]) {
	out = aSet.Set[string]{}
	it := p.Tree.Iterator()
	for it.Next() {
		syms := it.Value().(aSet.Set[string])
		out = out.Or(syms)
	}
	return out
}

type pendAction struct {
	ix int
	s  string
}

func (nos *NOS) IterPairs(callback func(prev, next *NOSNode) error) error {
	var prev *NOSNode = nil
	var next *NOSNode = nil

	nos.tree.Dft(func(n *NOSNode) error {
		next = n

		err := callback(prev, next)
		if err != nil {
			if errors.As(err, &ErrorStopIteration{}) {
				return nil
			}
			return err
		}

		prev = next

		return nil
	})

	// handle the last node
	err := callback(prev, nil)
	if err != nil {
		if errors.As(err, &ErrorStopIteration{}) {
			return nil
		}
		return err
	}

	return nil
}

func (nos *NOS) betweenPairs(startIx, endIx int, callback func(prev, next *NOSNode) error) error {
	var prev *NOSNode = nil
	var next *NOSNode = nil

	// right sib is inclusive if the index matches perfectly
	startNode, err := nos.tree.FindRightSibNode(startIx)
	if err != nil {
		return fmt.Errorf("findRightSib(%v): %w", startIx, err)
	}

	err = startNode.WalkRight(func(n *NOSNode) error {
		if n.StartIx > endIx {
			return ErrorStopIteration{}
		}

		next = n

		err := callback(prev, next)
		if err != nil {
			return err
		}

		prev = next

		return nil
	})

	if err != nil {
		if !errors.As(err, &ErrorStopIteration{}) {
			return fmt.Errorf("walkRight(%v): %w", startNode, err)
		}
	}

	// handle the last node
	err = callback(prev, nil)
	if err != nil {
		if errors.As(err, &ErrorStopIteration{}) {
			return nil
		}
		return err
	}

	return nil
}

func (nos *NOS) shouldMerge(a, b *NOSNode) bool {
	if !a.Format.Equals(b.Format) {
		return false
	}

	return a.EndIx+1 == b.StartIx
}

func (nos *NOS) mergeNeighbors() *NOS {
	toInsert := make([]*NOSNode, 0, nos.tree.Size)

	nos.tree.Dft(func(n *NOSNode) error {
		f := n.Format.DropNull()
		if sf, ok := f.(FormatV3Span); ok {
			delete(sf, "e")
			delete(sf, "en")
		}

		if !f.Empty() {
			toInsert = append(toInsert, &NOSNode{
				StartIx: n.StartIx,
				EndIx:   n.EndIx,
				Format:  f,
			})
		}

		return nil
	})

	for i := len(toInsert) - 1; i > 0; i-- {
		a, b := toInsert[i-1], toInsert[i]

		if !a.Format.IsSpan() || !b.Format.IsSpan() {
			continue
		}

		if nos.shouldMerge(a, b) {
			toInsert = DeleteAt(toInsert, i-1, 2)
			toInsert = InsertAt(toInsert, i-1, &NOSNode{
				StartIx: a.StartIx,
				EndIx:   b.EndIx,
				Format:  a.Format,
			})
		}
	}

	out := NewNOS()
	for _, n := range toInsert {
		out.tree.Put(n)
	}

	return out
}
