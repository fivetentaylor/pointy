package v3

import (
	"encoding/json"
	"fmt"

	"github.com/teamreviso/code/pkg/stackerr"
)

type FugueNode struct {
	ID            ID
	Text          []uint16
	Side          Side
	Parent        *FugueNode
	IsDeleted     []bool
	LeftChildren  []*FugueNode
	RightChildren []*FugueNode
}

func NewFugueNode(id ID, text []uint16, side Side, parent *FugueNode) *FugueNode {
	node := &FugueNode{
		ID:            id,
		Text:          text,
		Side:          side,
		Parent:        parent,
		IsDeleted:     make([]bool, len(text)),
		LeftChildren:  []*FugueNode{},
		RightChildren: []*FugueNode{},
	}
	if len(text) != len(node.IsDeleted) {
		panic(fmt.Sprintf("text and IsDeleted length mismatch: %d != %d", len(text), len(node.IsDeleted)))
	}
	if parent == nil {
		node.Parent = node
	}
	return node
}

func (node FugueNode) String() string {
	if node.Parent == nil {
		return fmt.Sprintf("FugueNode{%s: %q, Parent: nil, Side: %s, Length: %d}", node.ID, Uint16ToStr(node.Text), node.Side, len(node.Text))
	}
	return fmt.Sprintf("FugueNode{%s: %q, Parent: %s, Side: %s, Length: %d}", node.ID, Uint16ToStr(node.Text), adjustedParentId(&node), node.Side, len(node.Text))
}

func (node *FugueNode) ContainsID(id ID) bool {
	if node.ID.Author != id.Author {
		return false
	}

	return node.ID.Seq <= id.Seq && id.Seq < node.ID.Seq+len(node.Text)
}

func (node *FugueNode) VisibleText() string {
	fv := node.Explode().Visible()
	return Uint16ToStr(fv.Text)
}

func (node *FugueNode) Bfs(cb func(node *FugueNode) error) error {
	if node == nil {
		return fmt.Errorf("Bfs() nil root")
	}

	queue := []*FugueNode{node}

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		if node == nil {
			return fmt.Errorf("Bfs() child is nil")
		}

		err := cb(node)
		if err != nil {
			return err
		}

		if node.LeftChildren != nil {
			queue = append(queue, node.LeftChildren...)
		}

		if node.RightChildren != nil {
			queue = append(queue, node.RightChildren...)
		}
	}

	return nil
}

// Flatten returns a list of all nodes in the tree, in BFS order
// Only used in admin debugging
func (node *FugueNode) Flatten() []*FugueNode {
	nodes := []*FugueNode{}

	node.Bfs(func(n *FugueNode) error {
		nodes = append(nodes, n)
		return nil
	})

	return nodes
}

type TextTuple struct {
	ID   ID
	Text string
}

// FlatTextNodes returns a list of all Text
func (node *FugueNode) TextTuples() []TextTuple {
	tups := []TextTuple{}
	for i := range node.Text {
		tup := TextTuple{
			ID:   ID{Author: node.ID.Author, Seq: node.ID.Seq + i},
			Text: Uint16ToStr(node.Text[i : i+1]),
		}
		tups = append(tups, tup)
	}
	return tups
}

func (node *FugueNode) getVisOffset(totOffset int) (int, error) {
	if totOffset < 0 || totOffset >= len(node.Text) {
		return -1, stackerr.New(fmt.Errorf("index out of bounds: %d for node: %v", totOffset, node.ID))
	}

	if node.IsDeleted[totOffset] == true {
		return -1, nil
	}

	visOffset := 0
	for i := 0; i < totOffset; i++ {
		if node.IsDeleted[i] == false {
			visOffset++
		}
	}

	return visOffset, nil
}

func (node *FugueNode) getTotOffset(visOffset int) (int, error) {
	if visOffset < 0 {
		return -1, stackerr.New(fmt.Errorf("index out of bounds: %d for node: %v", visOffset, node.ID))
	}

	ix := visOffset
	totOffset := 0

	for _, ts := range node.IsDeleted {
		if ts == false {
			if ix == 0 {
				break
			}

			ix--
		}

		totOffset++
	}

	if ix > 0 {
		return -1, stackerr.New(fmt.Errorf("index out of bounds: %d for node: %v", visOffset, node.ID))
	}

	return totOffset, nil
}

func (node *FugueNode) leftmost() *FugueNode {
	for len(node.LeftChildren) > 0 {
		node = node.LeftChildren[0]
	}
	return node
}

func (node *FugueNode) rightmost() *FugueNode {
	for len(node.RightChildren) > 0 {
		node = node.RightChildren[len(node.RightChildren)-1]
	}
	return node
}

func (sib *FugueNode) insertRight(node *FugueNode) {
	if len(sib.RightChildren) == 0 {
		node.Parent = sib
		sib.RightChildren = append(sib.RightChildren, node)
		node.Side = Right
	} else {
		rightSibling := sib.RightChildren[0].leftmost()
		node.Parent = rightSibling
		rightSibling.LeftChildren = append(rightSibling.LeftChildren, node)
		node.Side = Left
	}
}

func (sib *FugueNode) insertLeft(node *FugueNode) {
	if len(sib.LeftChildren) == 0 {
		node.Parent = sib
		sib.LeftChildren = append(sib.LeftChildren, node)
		node.Side = Left
	} else {
		lc := sib.LeftChildren
		leftSibling := lc[len(lc)-1].rightmost()
		node.Parent = leftSibling
		leftSibling.RightChildren = append(leftSibling.RightChildren, node)
		node.Side = Right
	}
}

func compFugueNodes(a, b *FugueNode) int {
	if a.ID == b.ID {
		return 0
	}

	if a.ID.Author == b.ID.Author {
		if a.ID.Seq < b.ID.Seq {
			return -1
		}
	} else if a.ID.Author < b.ID.Author {
		return -1
	}

	return 1
}

func (node *FugueNode) insertChild(side Side, child *FugueNode) {
	if side == Left {
		_, c := sortedInsert(node.LeftChildren, child, compFugueNodes)
		node.LeftChildren = c
	} else {
		_, c := sortedInsert(node.RightChildren, child, compFugueNodes)
		node.RightChildren = c
	}

}

func (node *FugueNode) LeftMostVisID() (ID, error) {
	for i := 0; i < len(node.Text); i++ {
		if node.IsDeleted[i] == false {
			return ID{
				Author: node.ID.Author,
				Seq:    node.ID.Seq + i,
			}, nil
		}
	}

	return ID{}, stackerr.New(fmt.Errorf("no visible text in node %v", node.ID))
}

func (node *FugueNode) RightMostVisID() (ID, error) {
	for i := len(node.Text) - 1; i >= 0; i-- {
		if node.IsDeleted[i] == false {
			return ID{
				Author: node.ID.Author,
				Seq:    node.ID.Seq + i,
			}, nil
		}
	}

	return ID{}, stackerr.New(fmt.Errorf("no visible text in node %v", node.ID))
}

func (node *FugueNode) LeftMostID() ID {
	return node.ID
}

func (node *FugueNode) RightMostID() ID {
	return ID{
		Author: node.ID.Author,
		Seq:    node.ID.Seq + (len(node.Text) - 1),
	}
}

func (node *FugueNode) getChildIndex() (int, error) {
	if node.Side == Root {
		return -1, nil
	}

	var children []*FugueNode
	if node.Side == Right {
		children = node.Parent.RightChildren
	} else {
		children = node.Parent.LeftChildren
	}

	for ix, child := range children {
		if child.ID == node.ID {
			return ix, nil
		}
	}

	return -1, stackerr.New(fmt.Errorf("node %+v is not a child of parent %+v", node.ID, node.Parent.ID))
}

func adjustedParentId(node *FugueNode) ID {
	if node.Parent == nil {
		return NoID
	}

	if node.Side == Left {
		return node.Parent.ID
	} else {
		seq := node.Parent.ID.Seq + max(0, len(node.Parent.Text)-1)
		return ID{Author: node.Parent.ID.Author, Seq: seq}
	}
}

// TODO: remove once all rogues use op based serialization
func (node *FugueNode) MarshalJSON() ([]byte, error) {
	serializedNode := []interface{}{
		node.ID,
		Uint16ToStr(node.Text),
		node.LeftChildren,
		node.RightChildren,
	}

	return json.Marshal(serializedNode)
}

func (node *FugueNode) Inspect() {
	fmt.Println("node id: ", node.ID)
	fmt.Printf("\ttext: %q bytes: %#v\n", Uint16ToStr(node.Text), node.Text)
	fmt.Printf("\tleft: %s\n", node.LeftChildren)
	fmt.Printf("\tright: %s\n", node.RightChildren)
	for _, left := range node.LeftChildren {
		left.Inspect()
	}
	for _, right := range node.RightChildren {
		right.Inspect()
	}
}

func (node *FugueNode) dft(callback func(*FugueNode) error) error {
	for _, child := range node.LeftChildren {
		child.dft(callback)
	}

	err := callback(node)
	if err != nil {
		return err
	}

	for _, child := range node.RightChildren {
		child.dft(callback)
	}

	return nil
}

// used for compaction to mutate the tree in place
func (node *FugueNode) pft(callback func(*FugueNode) error) error {
	for i := len(node.RightChildren) - 1; i >= 0; i-- {
		child := node.RightChildren[i]
		child.dft(callback)
	}

	err := callback(node)
	if err != nil {
		return err
	}

	for i := len(node.LeftChildren) - 1; i >= 0; i-- {
		child := node.LeftChildren[i]
		child.dft(callback)
	}

	return nil
}

type FugueTot struct {
	IDs       []ID
	IsDeleted []bool
	Text      []uint16
	Deleted   int
}

type FugueVis struct {
	IDs  []ID
	Text []uint16
}

type FugueRune struct {
	ID    ID
	Rune  uint16
	IsDel bool
}

func (node *FugueNode) Deleted() int {
	deleted := 0
	for _, isDel := range node.IsDeleted {
		if isDel {
			deleted++
		}
	}
	return deleted
}

func (node *FugueNode) Explode() FugueTot {
	ft := FugueTot{
		IDs:       make([]ID, len(node.Text)),
		IsDeleted: node.IsDeleted,
		Text:      node.Text,
		Deleted:   node.Deleted(),
	}

	for i := range node.Text {
		ft.IDs[i] = ID{node.ID.Author, node.ID.Seq + i}
	}

	return ft
}

func (ft FugueTot) Visible() FugueVis {
	visLength := len(ft.Text) - ft.Deleted
	if visLength == 0 {
		return FugueVis{}
	}

	fv := FugueVis{
		IDs:  make([]ID, 0, visLength),
		Text: make([]uint16, 0, visLength),
	}

	for i, char := range ft.Text {
		if ft.IsDeleted[i] == false {
			fv.IDs = append(fv.IDs, ft.IDs[i])
			fv.Text = append(fv.Text, char)
		}
	}

	return fv
}

func (node *FugueNode) VisRunes() []FugueRune {
	ft := node.Explode()
	fv := ft.Visible()

	frs := make([]FugueRune, len(fv.Text))
	for i, char := range fv.Text {
		frs[i] = FugueRune{fv.IDs[i], char, false}
	}

	return frs
}

func (node *FugueNode) Runes() []FugueRune {
	ft := node.Explode()

	frs := make([]FugueRune, len(ft.Text))
	for i, char := range ft.Text {
		frs[i] = FugueRune{ft.IDs[i], char, ft.IsDeleted[i]}
	}

	return frs
}

func (node *FugueNode) WalkRight(callback func(*FugueNode) error) error {
	if node == nil {
		return fmt.Errorf("nil node")
	}

	prevSide := Left
	n := node
	for {
		if prevSide == Left {
			err := callback(n)
			if err != nil {
				return err
			}

			for _, child := range n.RightChildren {
				err := child.dft(callback)
				if err != nil {
					return err
				}
			}
		}

		if n.Side == Root {
			return nil
		}

		ix, err := n.getChildIndex()
		if err != nil {
			return err
		}

		children := n.Parent.RightChildren
		if n.Side == Left {
			children = n.Parent.LeftChildren
		}

		for ix++; ix < len(children); ix++ {
			err := children[ix].dft(callback)
			if err != nil {
				return err
			}
		}

		prevSide = n.Side
		n = n.Parent
	}
}

// ValidateFugue validates the Fugue tree
//   - Checks that all children have the correct parent
func (node *FugueNode) ValidateParentSide() error {
	if node == nil {
		return fmt.Errorf("nil root")
	}

	// Validate all children have the correct parent
	queue := []*FugueNode{node}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		if node.LeftChildren != nil {
			for _, lnode := range node.LeftChildren {
				if lnode.Parent == nil {
					return fmt.Errorf("child %+v has no parent, it should be %+v", lnode.ID, node.ID)
				}
				if lnode.Parent != node {
					return fmt.Errorf("child %+v has parent %+v, it should be %+v", lnode.ID, node.Parent.ID, node.ID)
				}
				if lnode.Side != Left {
					return fmt.Errorf("child %+v has side %+v, it should be %+v", lnode.ID, lnode.Side, Left)
				}
			}

			queue = append(queue, node.LeftChildren...)
		}

		if node.RightChildren != nil {
			for _, rnode := range node.RightChildren {
				if rnode.Parent == nil {
					return fmt.Errorf("child %+v has no parent, it should be %+v", rnode.ID, node.ID)
				}
				if rnode.Parent != node {
					return fmt.Errorf("child %+v has parent %+v, it should be %+v", rnode.ID, node.Parent.ID, node.ID)
				}
				if rnode.Side != Right {
					return fmt.Errorf("child %+v has side %+v, it should be %+v", rnode.ID, rnode.Side, Right)
				}
			}
			queue = append(queue, node.RightChildren...)
		}
	}

	return nil
}
