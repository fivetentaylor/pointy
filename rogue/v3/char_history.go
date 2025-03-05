package v3

import (
	"fmt"

	avl "github.com/teamreviso/code/rogue/v3/tree"
)

type Marker struct {
	ID    ID
	IsDel bool
}

type CharHistory map[ID]*avl.Tree[ID, *Marker]

func NewCharHistory() CharHistory {
	return make(CharHistory)
}

func markerKey(marker *Marker) (ID, error) {
	return marker.ID, nil
}

func (ch CharHistory) Add(targetID ID, m *Marker) {
	tree, ok := ch[targetID]

	if !ok {
		tree = avl.NewTree(markerKey, lessThanRevID)
		ch[targetID] = tree
	}

	tree.Put(m)
}

func (ch CharHistory) Max(targetID ID) Marker {
	tree, ok := ch[targetID]

	if !ok {
		return Marker{ID: targetID, IsDel: false}
	}

	m := tree.Max()
	if m == nil {
		return Marker{ID: targetID, IsDel: false}
	}
	return *m
}

func (r *Rogue) MarkCharDel(id ID, isDel bool) error {
	rn := r.RopeIndex.Get(id)
	if rn == nil {
		return fmt.Errorf("MarkCharVis: could not find node for id %v", id)
	}

	node := rn.Val
	ix := id.Seq - node.ID.Seq
	curDel := node.IsDeleted[ix]
	if curDel != isDel {
		if isDel {
			r.VisSize--
		} else {
			r.VisSize++
		}

		node.IsDeleted[ix] = isDel
		rn.updateWeight()
	}

	return nil
}

func (ch CharHistory) Print(tartetID ID) {
	tree, ok := ch[tartetID]

	if !ok {
		return
	}

	out := make([]Marker, 0, tree.Size)
	tree.Dft(func(m *Marker) error {
		out = append(out, *m)
		return nil
	})
	fmt.Printf("%v\n", out)
}
