package v3

import (
	"errors"
	"fmt"
	"strings"

	"github.com/teamreviso/code/pkg/stackerr"
)

// LeafNode and InternalNode types
type RopeNode struct {
	VisWeight int
	TotWeight int
	Parent    *RopeNode
	Height    int
	Left      *RopeNode
	Right     *RopeNode
	Val       *FugueNode
}

func NewLeafNode(val *FugueNode) *RopeNode {
	totWeight := len(val.Text)
	visWeight := totWeight - val.Deleted()

	return &RopeNode{
		VisWeight: visWeight,
		TotWeight: totWeight,
		Parent:    nil,
		Height:    1,
		Left:      nil,
		Right:     nil,
		Val:       val,
	}
}

func NewInternalNode(visWeight, totWeight, height int, left, right *RopeNode, parent *RopeNode) *RopeNode {
	return &RopeNode{
		VisWeight: visWeight,
		TotWeight: totWeight,
		Parent:    parent,
		Height:    height,
		Left:      left,
		Right:     right,
		Val:       nil,
	}
}

type Indicies struct {
	Vis int
	Tot int
}

// Rope struct
type Rope struct {
	Root  *RopeNode
	Index RopeIndex
	Cache map[ID]Indicies
}

// Only use this when you're doing a lot of operations on the rope and the rope is NOT changing
// Always call in a sinlge function like
// r.SetCache()
// defer r.DelCache()
func (r *Rope) SetCache() {
	r.Cache = make(map[ID]Indicies)
}

func (r *Rope) DelCache() {
	r.Cache = nil
}

func NewRope(index RopeIndex) *Rope {
	return &Rope{
		Root:  nil,
		Index: index,
	}
}

func (r *Rope) Print() {
	PrintRopeNode(r.Root, 0)
}

func (node *RopeNode) IsLeaf() (bool, error) {
	if node.Val != nil {
		if node.Left != nil || node.Right != nil {
			return false, stackerr.New(fmt.Errorf("invalid leaf rope node"))
		}
		return true, nil
	}
	return false, nil
}

func (node *RopeNode) IsInternal() (bool, error) {
	if node.Val == nil {
		if node.Left == nil || node.Right == nil {
			return false, stackerr.New(fmt.Errorf("invalid internal rope node"))
		}
		return true, nil
	}
	return false, nil
}

func getNode(root *RopeNode, visIx int) (int, *RopeNode, error) {
	if root == nil || visIx < 0 {
		return -1, nil, stackerr.New(ErrorInvalidOffset{nil, visIx})
	}

	isLeaf, err := root.IsLeaf()
	if err != nil {
		return -1, nil, err
	}

	if isLeaf {
		if visIx >= len(root.Val.Text)-root.Val.Deleted() {
			return -1, nil, stackerr.New(ErrorInvalidOffset{&root.Val.ID, visIx})
		}

		return visIx, root, nil
	}

	if visIx < root.VisWeight {
		return getNode(root.Left, visIx)
	} else {
		return getNode(root.Right, visIx-root.VisWeight)
	}
}

func (r *Rope) GetNode(visIx int) (int, *RopeNode, error) {
	visOffset, node, err := getNode(r.Root, visIx)
	if err != nil {
		return -1, nil, err
	}

	return visOffset, node, err
}

func getTotNode(root *RopeNode, totIx int) (int, *RopeNode, error) {
	if root == nil || totIx < 0 {
		return -1, nil, stackerr.New(ErrorInvalidOffset{nil, totIx})
	}

	isLeaf, err := root.IsLeaf()
	if err != nil {
		return -1, nil, err
	}

	if isLeaf {
		if totIx >= len(root.Val.Text) {
			return -1, nil, stackerr.New(ErrorInvalidOffset{&root.Val.ID, totIx})
		}

		return totIx, root, nil
	}

	if totIx < root.TotWeight {
		return getTotNode(root.Left, totIx)
	} else {
		return getTotNode(root.Right, totIx-root.TotWeight)
	}
}

func (r *Rope) GetTotNode(totIx int) (int, *RopeNode, error) {
	totOffset, node, err := getTotNode(r.Root, totIx)
	if err != nil {
		return -1, nil, err
	}

	return totOffset, node, err
}

func (r *Rope) InsertWithIx(totIx int, val *FugueNode) (*RopeNode, error) {
	newNode, err := insertIntoRopeNode(r.Root, totIx, val)
	if err != nil {
		return nil, err
	}
	r.Index.Put(newNode)

	p := newNode
	for p.Parent != nil {
		p = p.Parent
		updateHeight(p)
		p, err = rebalance(p)
		if err != nil {
			return nil, err
		}
	}

	r.Root = p

	return newNode, nil
}

func (r *Rope) Insert(val *FugueNode) (*RopeNode, error) {
	_, totIx, err := r.getInsertIx(val)
	if err != nil {
		return nil, err
	}

	return r.InsertWithIx(totIx, val)
}

func (r *Rope) GetVisID(visIx int) (ID, error) {
	if r.Root == nil {
		return ID{}, stackerr.New(fmt.Errorf("index %d out of bounds", visIx))
	}

	visOffset, node, err := r.GetNode(visIx)
	if err != nil {
		return ID{}, err
	}
	totOffset, err := node.Val.getTotOffset(visOffset)
	if err != nil {
		return ID{}, err
	}

	return ID{Author: node.Val.ID.Author, Seq: node.Val.ID.Seq + totOffset}, nil
}

func (r *Rope) GetTotID(totIx int) (ID, error) {
	if r.Root == nil {
		return ID{}, stackerr.New(fmt.Errorf("index %d out of bounds", totIx))
	}

	totOffset, node, err := r.GetTotNode(totIx)
	if err != nil {
		return ID{}, err
	}

	return ID{Author: node.Val.ID.Author, Seq: node.Val.ID.Seq + totOffset}, nil
}

// GetIndex returns the visible and total index of the node
func (r *Rope) GetIndex(id ID) (visIx, totIx int, err error) {
	if r.Cache != nil {
		if idx, ok := r.Cache[id]; ok {
			return idx.Vis, idx.Tot, nil
		}
	}

	node := r.Index.Get(id)
	if node == nil {
		return -1, -1, stackerr.New(fmt.Errorf("id: %+v not in rope", id))
	}

	totIx = id.Seq - node.Val.ID.Seq
	visIx, err = node.Val.getVisOffset(totIx)
	if err != nil {
		return -1, -1, err
	}

	n := node
	for n.Parent != nil {
		if n.Parent.Right == n {
			if visIx >= 0 {
				visIx += n.Parent.VisWeight
			}
			totIx += n.Parent.TotWeight
		}
		n = n.Parent
	}

	if r.Cache != nil {
		r.Cache[id] = Indicies{visIx, totIx}
	}

	return visIx, totIx, nil
}

func (r *Rope) GetVisWeight(id ID) (int, error) {
	rn := r.Index.Get(id)
	if rn == nil {
		return -1, stackerr.New(fmt.Errorf("id: %+v not in rope", id))
	}

	node := rn.Val
	visWeight := 0
	totOffset := id.Seq - node.ID.Seq
	for i := totOffset; i >= 0; i-- {
		if node.IsDeleted[i] == false {
			visWeight++
		}
	}

	n := rn
	for n.Parent != nil {
		if n.Parent.Right == n {
			visWeight += n.Parent.VisWeight
		}
		n = n.Parent
	}

	return visWeight, nil
}

func insertIntoRopeNode(root *RopeNode, ix int, val *FugueNode) (*RopeNode, error) {
	if root == nil {
		leaf := NewLeafNode(val)
		return leaf, nil
	}

	isLeaf, err := root.IsLeaf()
	if err != nil {
		return nil, err
	}
	if isLeaf {
		newNode := NewLeafNode(val)

		var left, right *RopeNode
		left, right = root, newNode
		if ix == 0 {
			left, right = right, left
		}

		newParent := NewInternalNode(left.VisWeight, left.TotWeight, 2, left, right, nil)
		_, err = swapChild(root, newParent)
		if err != nil {
			return nil, err
		}
		left.Parent, right.Parent = newParent, newParent

		return newNode, nil
	} else {
		weight := root.TotWeight
		if ix <= weight {
			root.VisWeight += len(val.Text) - val.Deleted()
			root.TotWeight += len(val.Text)
			return insertIntoRopeNode(root.Left, ix, val)
		} else {
			return insertIntoRopeNode(root.Right, ix-weight, val)
		}
	}
}

func swapChild(orgNode, newNode *RopeNode) (bool, error) {
	if orgNode == newNode {
		return true, nil
	}

	newNode.Parent = orgNode.Parent

	if orgNode.Parent != nil {
		if orgNode.Parent.Right == orgNode {
			orgNode.Parent.Right = newNode
		} else if orgNode.Parent.Left == orgNode {
			orgNode.Parent.Left = newNode
		} else {
			return false, stackerr.New(fmt.Errorf("node is not child of parent"))
		}
		return true, nil
	}

	return false, nil
}

func updateHeight(node *RopeNode) {
	if node == nil {
		return
	}

	leftHeight, rightHeight := 0, 0

	if node.Left != nil {
		leftHeight = node.Left.Height
	}
	if node.Right != nil {
		rightHeight = node.Right.Height
	}

	node.Height = 1 + max(leftHeight, rightHeight)
}

func rebalance(node *RopeNode) (*RopeNode, error) {
	if node == nil {
		return node, nil // Can only rebalance internal nodes
	}
	isLeaf, err := node.IsLeaf()
	if err != nil {
		return nil, err
	}
	if isLeaf {
		return node, nil
	}

	balance, err := getBalanceFactor(node)
	if err != nil {
		return nil, err
	}

	if balance > 1 {
		isInternal, err := node.Left.IsInternal()
		if err != nil {
			return nil, err
		}
		balanceFactor, err := getBalanceFactor(node.Left)
		if err != nil {
			return nil, err
		}
		if balanceFactor < 0 && isInternal {
			node.Left, err = rotateLeft(node.Left)
			if err != nil {
				return nil, err
			}
		}
		return rotateRight(node)
	}

	if balance < -1 {
		isInternal, err := node.Right.IsInternal()
		if err != nil {
			return nil, err
		}
		balanceFactor, err := getBalanceFactor(node.Right)
		if err != nil {
			return nil, err
		}
		if balanceFactor > 0 && isInternal {
			node.Right, err = rotateRight(node.Right)
			if err != nil {
				return nil, err
			}
		}
		return rotateLeft(node)
	}

	return node, nil // No rotation needed
}

func getBalanceFactor(node *RopeNode) (int, error) {
	if node == nil {
		return 0, nil
	}

	leftHeight, rightHeight := 0, 0

	isInternal, err := node.IsInternal()
	if err != nil {
		return 0, err
	}
	if isInternal {
		if node.Left != nil {
			leftHeight = node.Left.Height
		}
		if node.Right != nil {
			rightHeight = node.Right.Height
		}
	}

	return leftHeight - rightHeight, nil
}

func rotateLeft(root *RopeNode) (*RopeNode, error) {
	if root.Right == nil {
		return root, nil
	}

	isLeaf, err := root.Right.IsLeaf()
	if err != nil {
		return nil, err
	}
	if isLeaf {
		return root, nil
	}

	newRoot := root.Right
	root.Right = newRoot.Left
	if root.Right != nil {
		root.Right.Parent = root
	}
	newRoot.Left = root

	_, err = swapChild(root, newRoot)
	if err != nil {
		return nil, err
	}
	root.Parent = newRoot

	newRoot.VisWeight += root.VisWeight
	newRoot.TotWeight += root.TotWeight
	updateHeight(root)
	updateHeight(newRoot)

	return newRoot, nil
}

func rotateRight(root *RopeNode) (*RopeNode, error) {
	if root.Left == nil {
		return root, nil
	}

	isLeaf, err := root.Left.IsLeaf()
	if err != nil {
		return nil, err
	}

	if isLeaf {
		return root, nil
	}

	newRoot := root.Left
	root.Left = newRoot.Right
	if root.Left != nil {
		root.Left.Parent = root
	}
	newRoot.Right = root

	_, err = swapChild(root, newRoot)
	if err != nil {
		return nil, err
	}
	root.Parent = newRoot

	root.VisWeight -= newRoot.VisWeight
	root.TotWeight -= newRoot.TotWeight
	updateHeight(root)
	updateHeight(newRoot)

	return newRoot, nil
}

func (node *RopeNode) Side() Side {
	if node.Parent == nil {
		return Root
	}
	if node.Parent.Left == node {
		return Left
	}
	return Right
}

func (node *RopeNode) RightVisWeight() (int, error) {
	isLeaf, err := node.IsLeaf()
	if err != nil {
		return 0, err
	}

	if isLeaf {
		return node.VisWeight, nil
	}

	n := node.Right
	rightVisWeight := 0
	for {
		rightVisWeight += n.VisWeight

		isLeaf, err := n.IsLeaf()
		if err != nil {
			return 0, err
		}

		if isLeaf {
			return rightVisWeight, nil
		}

		n = n.Right
	}
}

func (node *RopeNode) rightmostLeaf() (*RopeNode, error) {
	isLeaf, err := node.IsLeaf()
	if err != nil {
		return nil, stackerr.New(fmt.Errorf("rightmostLeaf failed to check if node is leaf: %v", err))
	}
	if isLeaf {
		return node, nil
	}
	return node.Right.rightmostLeaf()
}

func (node *RopeNode) rightmostVisLeaf() (*RopeNode, error) {
	n := node
	for {
		isLeaf, err := n.IsLeaf()
		if err != nil {
			return nil, err
		}

		if isLeaf {
			if n.VisWeight == 0 {
				return nil, stackerr.New(fmt.Errorf("no visible right sibling"))
			}

			return n, nil
		}

		rightWeight, err := n.RightVisWeight()
		if err != nil {
			return nil, err
		}

		if rightWeight > 0 {
			n = n.Right
		} else {
			n = n.Left
		}
	}
}

func (node *RopeNode) leftmostLeaf() (*RopeNode, error) {
	isLeaf, err := node.IsLeaf()
	if err != nil {
		return nil, fmt.Errorf("IsLeaf(): %w", err)
	}
	if isLeaf {
		return node, nil
	}
	return node.Left.leftmostLeaf()
}

func (node *RopeNode) leftmostVisLeaf() (*RopeNode, error) {
	n := node
	for {
		isLeaf, err := n.IsLeaf()
		if err != nil {
			return nil, err
		}

		if isLeaf {
			if n.VisWeight == 0 {
				return nil, stackerr.New(fmt.Errorf("no visible left sibling"))
			}

			return n, nil
		}

		if n.VisWeight > 0 {
			n = n.Left
		} else {
			n = n.Right
		}
	}
}

func (node *RopeNode) LeftTotSibling() (*RopeNode, error) {
	n := node
	for n.Parent != nil {
		if n.Side() == Right {
			return n.Parent.Left.rightmostLeaf()
		}
		n = n.Parent
	}

	return nil, stackerr.New(ErrorNoLeftTotSibling{node.Val.ID})
}

func (node *RopeNode) RightTotSibling() (*RopeNode, error) {
	n := node
	for n.Parent != nil {
		if n.Side() == Left {
			return n.Parent.Right.leftmostLeaf()
		}
		n = n.Parent
	}

	return nil, stackerr.New(ErrorNoRightTotSibling{node.Val.ID})
}

func (node *RopeNode) LeftVisSibling() (*RopeNode, error) {
	n := node
	for n.Parent != nil {
		if n.Side() == Right && n.Parent.VisWeight > 0 {
			return n.Parent.Left.rightmostVisLeaf()
		}
		n = n.Parent
	}

	return nil, stackerr.New(ErrorNoLeftVisSibling{node.Val.ID})
}

func (node *RopeNode) RightVisSibling() (*RopeNode, error) {
	n := node
	for n.Parent != nil {
		visWeight, err := n.Parent.RightVisWeight()
		if err != nil {
			return nil, err
		}

		if n.Side() == Left && visWeight > 0 {
			return n.Parent.Right.leftmostVisLeaf()
		}
		n = n.Parent
	}

	return nil, stackerr.New(ErrorNoRightVisSibling{node.Val.ID})
}

func (node *RopeNode) validateWeights() (int, error) {
	if node == nil {
		return 0, nil
	}

	isLeaf, err := node.IsLeaf()
	if err != nil {
		return 0, err
	}
	if isLeaf {
		return node.VisWeight, nil
	}

	leftWeight, err := node.Left.validateWeights()
	if err != nil {
		return 0, err
	}
	rightWeight, err := node.Right.validateWeights()
	if err != nil {
		return 0, err
	}

	if node.VisWeight != leftWeight {
		return 0, stackerr.New(fmt.Errorf("invalid node found: %+v", node))
	}

	return node.VisWeight + rightWeight, nil
}

func (node *RopeNode) validateHeights() (bool, error) {
	if node == nil {
		return true, nil
	}

	isLeaf, err := node.IsLeaf()
	if err != nil {
		return false, err
	}
	if isLeaf {
		return node.Height == 1, err
	}

	leftHeight, rightHeight := 0, 0
	if node.Left != nil {
		leftHeight = node.Left.Height
	}
	if node.Right != nil {
		rightHeight = node.Right.Height
	}

	expectedHeight := max(leftHeight, rightHeight) + 1
	if node.Height != expectedHeight {
		return false, nil
	}

	leftHeightValid, err := node.Left.validateHeights()
	if err != nil {
		return false, err
	}
	rightHeightValid, err := node.Right.validateHeights()
	if err != nil {
		return false, err
	}

	return leftHeightValid && rightHeightValid, nil
}

func validateParentPointers(node, parent *RopeNode) (bool, error) {
	if node == nil {
		return true, nil
	}

	if node.Parent != parent {
		return false, nil
	}

	isLeaf, err := node.IsLeaf()
	if err != nil {
		return false, fmt.Errorf("validateParentPointers failed to check if node is leaf: %v", err)
	}
	if !isLeaf {
		leftValid, err := validateParentPointers(node.Left, node)
		if err != nil {
			return false, err
		}
		rightValid, err := validateParentPointers(node.Right, node)
		if err != nil {
			return false, err
		}
		return leftValid && rightValid, nil
	}

	return true, nil
}

func (r *Rope) Validate() error {
	_, err := r.Root.validateWeights() // errors when encountering bad weight
	if err != nil {
		return fmt.Errorf("invalid weights: %w", err)
	}

	valid, err := r.Root.validateHeights() // errors when encountering bad height
	if err != nil {
		return fmt.Errorf("invalid heights: %w", err)
	}
	if !valid {
		return fmt.Errorf("invalid heights")
	}

	valid, err = validateParentPointers(r.Root, nil)
	if err != nil {
		return fmt.Errorf("invalid parent pointers: %w", err)
	}
	if !valid {
		return fmt.Errorf("invalid parent pointers")
	}

	return nil
}

func idToString(id ID) string {
	return fmt.Sprintf("%s_%d", id.Author, id.Seq)
}

func PrintRopeNode(node *RopeNode, depth int) {
	if node == nil {
		return
	}

	indent := strings.Repeat(" ", depth)
	isLeaf, err := node.IsLeaf()
	if err != nil {
		fmt.Printf("%sError: %v\n", indent, err)
		return
	}
	if isLeaf {
		// Assuming node has fields like Text and Id or a reference to another struct with these fields
		fmt.Printf("%s(%p) %q %d %d %d %s\n", indent, node, Uint16ToStr(node.Val.Text), node.VisWeight, node.TotWeight, node.Height, idToString(node.Val.ID))
	} else {
		fmt.Printf("%s(%p) %d %d %d\n", indent, node, node.VisWeight, node.TotWeight, node.Height)
		fmt.Printf("%sLeft:\n", indent)
		PrintRopeNode(node.Left, depth+1)
		fmt.Printf("%sRight:\n", indent)
		PrintRopeNode(node.Right, depth+1)
	}
}

func (node *RopeNode) updateWeight() {
	totLength := len(node.Val.Text)
	visLength := totLength - node.Val.Deleted()

	totDiff := totLength - node.TotWeight
	visDiff := visLength - node.VisWeight

	// nothing to update
	if totDiff == 0 && visDiff == 0 {
		return
	}

	node.TotWeight = totLength
	node.VisWeight = visLength
	n := node

	for n.Parent != nil {
		if n.Parent.Left == n {
			n.Parent.VisWeight += visDiff
			n.Parent.TotWeight += totDiff
		}
		n = n.Parent
	}
}

func (node *RopeNode) Dft(callback func(*RopeNode) error) error {
	if node == nil {
		return nil
	}

	node.Left.Dft(callback)

	ok, err := node.IsLeaf()
	if err != nil {
		return err
	}

	if ok {
		err := callback(node)
		if err != nil {
			return err
		}
	}

	node.Right.Dft(callback)
	return nil
}

func (node *RopeNode) DftVis(callback func(*RopeNode) error) error {
	if node == nil {
		return nil
	}

	// only visit nodes with visible text
	if node.VisWeight > 0 {
		err := node.Left.DftVis(callback)
		if err != nil {
			return err
		}

		ok, err := node.IsLeaf()
		if err != nil {
			return err
		}

		if ok {
			err := callback(node)
			if err != nil {
				return err
			}
		}
	}

	return node.Right.DftVis(callback)
}

func (r *Rope) TotRightOf(id ID) (ID, error) {
	ropeNode := r.Index.Get(id)
	if ropeNode == nil {
		return ID{}, stackerr.New(fmt.Errorf("node with ID %+v doesn't exist", id))
	}

	rightID := ropeNode.Val.RightMostID()
	if id == rightID {
		rn, err := ropeNode.RightTotSibling()
		if err != nil {
			return ID{}, err
		}

		leftID := rn.Val.LeftMostID()
		return leftID, nil
	}

	return ID{Author: id.Author, Seq: id.Seq + 1}, nil

}

func (r *Rope) TotLeftOf(id ID) (ID, error) {
	ropeNode := r.Index.Get(id)
	if ropeNode == nil {
		return ID{}, stackerr.New(fmt.Errorf("node with ID %+v doesn't exist", id))
	}

	leftID := ropeNode.Val.LeftMostID()
	if id == leftID {
		rn, err := ropeNode.LeftTotSibling()
		if err != nil {
			return ID{}, err
		}

		rightID := rn.Val.RightMostID()
		return rightID, nil
	}

	return ID{Author: id.Author, Seq: id.Seq - 1}, nil
}

func (r *Rope) VisRightOf(id ID) (ID, error) {
	ropeNode := r.Index.Get(id)
	if ropeNode == nil {
		return ID{}, stackerr.New(fmt.Errorf("node with ID %+v doesn't exist", id))
	}

	n := ropeNode.Val
	totOffset := id.Seq - n.ID.Seq

	for i := totOffset + 1; i < len(n.Text); i++ {
		if n.IsDeleted[i] == false {
			return ID{Author: id.Author, Seq: n.ID.Seq + i}, nil
		}
	}

	rSib, err := ropeNode.RightVisSibling()
	if err != nil {
		return ID{}, err
	}

	leftVisID, err := rSib.Val.LeftMostVisID()
	if err != nil {
		return ID{}, err
	}

	return leftVisID, nil
}

func (r *Rope) VisLeftOf(id ID) (ID, error) {
	ropeNode := r.Index.Get(id)
	if ropeNode == nil {
		return ID{}, stackerr.New(fmt.Errorf("node with ID %+v doesn't exist", id))
	}

	n := ropeNode.Val
	totOffset := id.Seq - n.ID.Seq

	for i := totOffset - 1; i >= 0; i-- {
		if n.IsDeleted[i] == false {
			return ID{Author: id.Author, Seq: n.ID.Seq + i}, nil
		}
	}

	lSib, err := ropeNode.LeftVisSibling()
	if err != nil {
		return ID{}, err
	}

	rightVisID, err := lSib.Val.RightMostVisID()
	if err != nil {
		return ID{}, err
	}

	return rightVisID, nil
}

func (r *Rope) VisLeftOfConstrainedByAddress(id ID, ca *ContentAddress) (ID, error) {
	leftId := id
	var err error
	for {
		leftId, err = r.VisLeftOf(leftId)
		if err != nil {
			return ID{}, fmt.Errorf("VisLeftOf(%v): %w", id, err)
		}

		if ca.Contains(leftId) {
			return leftId, nil
		}
	}
}

func (r *Rope) IsDeleted(id ID) (bool, error) {
	ropeNode := r.Index.Get(id)
	if ropeNode == nil {
		return false, stackerr.New(fmt.Errorf("node with ID %+v doesn't exist", id))
	}

	node := ropeNode.Val
	totOffset := id.Seq - node.ID.Seq
	return node.IsDeleted[totOffset], nil
}

func (r *Rope) GetBetween(startID, endID ID) (*FugueVis, error) {
	isDel, err := r.IsDeleted(startID)
	if err != nil {
		return nil, err
	}

	if isDel {
		startID, err = r.VisRightOf(startID)
		if err != nil {
			if errors.As(err, &ErrorNoRightVisSibling{}) {
				return &FugueVis{
					IDs:  []ID{},
					Text: []uint16{},
				}, nil
			}
			return nil, err
		}
	}

	isDel, err = r.IsDeleted(endID)
	if err != nil {
		return nil, fmt.Errorf("IsDeleted(%v): %w", endID, err)
	}

	if isDel {
		endID, err = r.VisLeftOf(endID)
		if err != nil {
			if errors.As(err, &ErrorNoLeftVisSibling{}) {
				return &FugueVis{
					IDs:  []ID{},
					Text: []uint16{},
				}, nil
			}
			return nil, fmt.Errorf("VisLeftOf(%v): %w", endID, err)
		}
	}

	startIx, _, err := r.GetIndex(startID)
	if err != nil {
		return nil, fmt.Errorf("GetIndex(%v): %w", startID, err)
	}

	endIx, _, err := r.GetIndex(endID)
	if err != nil {
		return nil, fmt.Errorf("GetIndex(%v): %w", endID, err)
	}

	if startIx > endIx {
		return &FugueVis{
			IDs:  []ID{},
			Text: []uint16{},
		}, nil
	}

	out := FugueVis{
		IDs:  make([]ID, 0, endIx-startIx+1),
		Text: make([]uint16, 0, endIx-startIx+1),
	}

	node := r.Index.Get(startID)
	if node == nil {
		return nil, fmt.Errorf("node %v doesn't exist", startID)
	}

	for {
		isStart := node.Val.ContainsID(startID)
		isEnd := node.Val.ContainsID(endID)

		fn := node.Val.Explode()
		for i, id := range fn.IDs {
			if isStart && id.Seq < startID.Seq {
				continue
			}

			if fn.IsDeleted[i] == false {
				out.IDs = append(out.IDs, id)
				out.Text = append(out.Text, fn.Text[i])
			}

			if isEnd && id.Seq >= endID.Seq {
				return &out, nil
			}

		}

		node, err = node.RightVisSibling()
		if err != nil {
			return nil, fmt.Errorf("rightTotSibling(): %w", err)
		}
	}
}

func (r *Rope) GetTotBetween(startID, endID ID) (*FugueTot, error) {
	_, startIx, err := r.GetIndex(startID)
	if err != nil {
		return nil, err
	}

	_, endIx, err := r.GetIndex(endID)
	if err != nil {
		return nil, err
	}

	if startIx > endIx {
		return &FugueTot{
			IDs:       []ID{},
			IsDeleted: []bool{},
			Text:      []uint16{},
			Deleted:   0,
		}, nil
	}

	out := FugueTot{
		IDs:       make([]ID, 0, endIx-startIx),
		IsDeleted: make([]bool, 0, endIx-startIx),
		Text:      make([]uint16, 0, endIx-startIx),
		Deleted:   0,
	}

	node := r.Index.Get(startID)
	if node == nil {
		return nil, stackerr.New(fmt.Errorf("node %v doesn't exist", startID))
	}

	for {
		isStart := node.Val.ContainsID(startID)
		isEnd := node.Val.ContainsID(endID)

		fn := node.Val.Explode()
		for i, id := range fn.IDs {
			if isStart && id.Seq < startID.Seq {
				continue
			}

			out.IDs = append(out.IDs, id)
			out.IsDeleted = append(out.IsDeleted, fn.IsDeleted[i])
			out.Text = append(out.Text, fn.Text[i])
			if fn.IsDeleted[i] == true {
				out.Deleted++
			}

			if isEnd && id.Seq >= endID.Seq {
				return &out, nil
			}
		}

		node, err = node.RightTotSibling()
		if err != nil {
			return nil, err
		}
	}
}

func (r *Rope) WalkRight(startID ID, callback func(*RopeNode) error) error {
	startNode := r.Index.Get(startID)
	if startNode == nil {
		return fmt.Errorf("node %v doesn't exist", startID)
	}

	node := startNode
	for node != nil {
		err := callback(node)
		if err != nil {
			if errors.As(err, &ErrorStopIteration{}) {
				return nil
			}

			return err
		}

		node, err = node.RightTotSibling()
		if err != nil {
			if errors.As(err, &ErrorNoRightTotSibling{}) {
				return nil
			}

			return fmt.Errorf("RightTotSibling(): %w", err)
		}
	}

	return nil
}

func (r *Rope) getInsertIx(node *FugueNode) (int, int, error) {
	if r.Root == nil {
		return 0, 0, nil
	}

	if len(node.RightChildren) > 0 {
		rightSib := node.RightChildren[0].leftmost()
		return r.GetIndex(rightSib.ID)
	}

	lc := node.LeftChildren
	if len(lc) > 0 {
		leftSib := lc[len(lc)-1].rightmost()
		visIx, totIx, err := r.GetIndex(leftSib.ID)
		if err != nil {
			return 0, 0, err
		}
		totWeight := len(leftSib.Text)
		visWeight := totWeight - leftSib.Deleted()
		return visIx + visWeight, totIx + totWeight, nil
	}

	childIx, err := node.getChildIndex()
	if err != nil {
		return 0, 0, err
	}

	if node.Side == Right {
		var leftSib *FugueNode
		if childIx == 0 {
			leftSib = node.Parent
		} else {
			leftSib = node.Parent.RightChildren[childIx-1].rightmost()
		}
		visIx, totIx, err := r.GetIndex(leftSib.ID)
		if err != nil {
			return 0, 0, err
		}
		totWeight := len(leftSib.Text)
		visWeight := totWeight - leftSib.Deleted()
		return visIx + visWeight, totIx + totWeight, nil
	}

	if node.Side == Left {
		lc := node.Parent.LeftChildren
		var rightSib *FugueNode
		if childIx == len(lc)-1 {
			rightSib = node.Parent
		} else {
			rightSib = node.Parent.LeftChildren[childIx+1].leftmost()
		}
		visIx, totIx, err := r.GetIndex(rightSib.ID)
		if err != nil {
			return 0, 0, err
		}
		return visIx, totIx, nil
	}

	return 0, 0, stackerr.New(fmt.Errorf("couldn't find sibling for node: %+v", node.ID))
}

func (r *Rope) ContainsVisibleBetween(startID, endID ID) (bool, error) {
	startVisWeight, err := r.GetVisWeight(startID)
	if err != nil {
		return false, err
	}

	endVisWeight, err := r.GetVisWeight(endID)
	if err != nil {
		return false, err
	}

	return startVisWeight < endVisWeight-1, nil
}

func (r *Rope) ContainsVisible(startID, endID ID) (bool, error) {
	visStartIx, totStartIx, err := r.GetIndex(startID)
	if err != nil {
		return false, err
	}

	visEndIx, totEndIx, err := r.GetIndex(endID)
	if err != nil {
		return false, err
	}

	if totEndIx < totStartIx {
		return false, nil
	}

	if visStartIx > -1 || visEndIx > -1 {
		return true, nil
	}

	containsVis, err := r.ContainsVisibleBetween(startID, endID)
	if err != nil {
		return false, err
	}

	return containsVis, nil
}

func (r *Rope) VisToTotIx(visIx int) (int, error) {
	id, err := r.GetVisID(visIx)
	if err != nil {
		return -1, err
	}

	_, totIx, err := r.GetIndex(id)
	if err != nil {
		return -1, err
	}

	return totIx, nil
}

func (r *Rope) VisIxLeft(id ID) (int, error) {
	visIx, _, err := r.GetIndex(id)
	if err != nil {
		return -1, err
	}

	if visIx < 0 {
		id, err = r.VisLeftOf(id)
		if err != nil {
			if errors.As(err, &ErrorNoLeftVisSibling{}) {
				return 0, nil
			}

			return -1, fmt.Errorf("VisLeftOf(%v): %w", id, err)
		}

		visIx, _, err = r.GetIndex(id)
		if err != nil {
			return -1, fmt.Errorf("GetIndex(%v): %w", id, err)
		}
	}

	return visIx, nil
}

func (r *Rope) VisIxLeftConstrainedByAddress(id ID, ca *ContentAddress) (int, error) {
	visIx, _, err := r.GetIndex(id)
	if err != nil {
		return -1, fmt.Errorf("GetIndex(%v): %w", id, err)
	}

	if ca.Contains(id) && visIx >= 0 {
		return visIx, nil
	}

	for {
		id, err = r.VisLeftOf(id)
		if err != nil {
			if errors.As(err, &ErrorNoLeftVisSibling{}) {
				return 0, nil
			}

			return -1, fmt.Errorf("VisLeftOf(%v): %w", id, err)
		}

		if !ca.Contains(id) {
			continue
		}

		visIx, _, err = r.GetIndex(id)
		if err != nil {
			return -1, fmt.Errorf("GetIndex(%v): %w", id, err)
		}
		return visIx, nil
	}
}

func (r *Rope) VisIxRight(id ID) (int, error) {
	visIx, _, err := r.GetIndex(id)
	if err != nil {
		return -1, fmt.Errorf("GetIndex(%v): %w", id, err)
	}

	if visIx < 0 {
		id, err = r.VisRightOf(id)
		if err != nil {
			if errors.As(err, &ErrorNoRightVisSibling{}) {
				ix := max(0, r.visSize()-1)
				return ix, nil
			}

			return -1, fmt.Errorf("VisRightOf(%v): %w", id, err)
		}

		visIx, _, err = r.GetIndex(id)
		if err != nil {
			return -1, fmt.Errorf("GetIndex(%v): %w", id, err)
		}
	}

	return visIx, nil
}

func (r *Rope) visSize() int {
	if r.Root == nil {
		return 0
	}

	size := 0
	rn := r.Root
	for rn != nil {
		size += rn.VisWeight
		rn = rn.Right
	}

	return size
}

func (r *Rope) NearestVisRightOf(id ID) (ID, error) {
	isDel, err := r.IsDeleted(id)
	if err != nil {
		return NoID, err
	}

	if !isDel {
		return id, nil
	}

	visID, err := r.VisRightOf(id)
	if err != nil {
		return NoID, err
	}

	return visID, nil
}

func (r *Rope) NearestVisLeftOf(id ID) (ID, error) {
	isDel, err := r.IsDeleted(id)
	if err != nil {
		return NoID, err
	}

	if !isDel {
		return id, nil
	}

	visID, err := r.VisLeftOf(id)
	if err != nil {
		return NoID, err
	}

	return visID, nil
}
