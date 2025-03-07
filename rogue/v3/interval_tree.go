package v3

import (
	"fmt"

	"github.com/fivetentaylor/pointy/pkg/stackerr"
	"golang.org/x/exp/slices"
)

type IntervalNode struct {
	MaxID    ID
	MaxSeqID ID
	MinSeqID ID
	StartID  ID
	FormatOp FormatOp
	Left     *IntervalNode
	Right    *IntervalNode
	height   int
}

type IntervalTree struct {
	Rope     *Rope
	Root     *IntervalNode
	IsSticky bool
}

func (it *IntervalTree) KeyFunc(op FormatOp) (string, error) {
	_, startIx, err := it.Rope.GetIndex(op.StartID)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%016d%s", startIx, op.ID), nil
}

func NewIntervalTree(rope *Rope, isSticky bool) *IntervalTree {
	return &IntervalTree{
		Rope:     rope,
		Root:     nil,
		IsSticky: isSticky,
	}
}

func (it *IntervalTree) Insert(formatOp FormatOp) (err error) {
	sf, ok := formatOp.Format.Copy().(FormatV3Span)

	if !ok {
		return stackerr.New(fmt.Errorf("formatOp %v is not a FormatV3Span", formatOp))
	}

	sticky, noSticky := sf.SplitSticky()
	if it.IsSticky && len(noSticky) > 0 {
		return stackerr.New(fmt.Errorf("noSticky span found in sticky: %v", formatOp))
	} else if !it.IsSticky && len(sticky) > 0 {
		return stackerr.New(fmt.Errorf("sticky span found in noSticky: %v", formatOp))
	}

	delete(sf, "e")
	delete(sf, "en")

	it.Root, err = it._insert(it.Root, formatOp)
	if err != nil {
		return err
	}
	return nil
}

func (it *IntervalTree) _insert(node *IntervalNode, formatOp FormatOp) (*IntervalNode, error) {
	startID := formatOp.StartID
	endID := formatOp.EndID

	if node == nil {
		return &IntervalNode{
			StartID:  startID,
			MaxID:    endID,
			MaxSeqID: formatOp.ID,
			MinSeqID: formatOp.ID,
			FormatOp: formatOp,
		}, nil
	}

	nodeKey, err := it.KeyFunc(node.FormatOp)
	if err != nil {
		return nil, err
	}

	opKey, err := it.KeyFunc(formatOp)
	if err != nil {
		return nil, err
	}

	_, maxIx, err := it.Rope.GetIndex(node.MaxID)
	if err != nil {
		return nil, err
	}

	_, opEndIx, err := it.Rope.GetIndex(formatOp.EndID)
	if err != nil {
		return nil, err
	}

	if maxIx < opEndIx {
		node.MaxID = endID
	}

	if node.MaxSeqID.lessThan(formatOp.ID) {
		node.MaxSeqID = formatOp.ID
	}

	if formatOp.ID.lessThan(node.MinSeqID) {
		node.MinSeqID = formatOp.ID
	}

	if opKey == nodeKey {
		node.FormatOp = formatOp
	} else if opKey < nodeKey {
		node.Left, err = it._insert(node.Left, formatOp)
		if err != nil {
			return nil, err
		}
	} else {
		node.Right, err = it._insert(node.Right, formatOp)
		if err != nil {
			return nil, err
		}
	}

	// balance the tree
	node.height = 1 + max(node.Left.Height(), node.Right.Height())

	node, err = it.balance(opKey, node)
	if err != nil {
		return nil, err
	}

	return node, nil
}

func (it *IntervalTree) balance(opKey string, node *IntervalNode) (*IntervalNode, error) {
	balance := node.getBalance()

	if node.Left != nil {
		nodeKey, err := it.KeyFunc(node.Left.FormatOp)
		if err != nil {
			return nil, err
		}

		// Left Left Case
		if balance > 1 && opKey < nodeKey {
			return it.rotateRight(node)
		}

		// Left Right Case
		if balance > 1 && opKey > nodeKey {
			node.Left, err = it.rotateLeft(node.Left)
			if err != nil {
				return nil, err
			}
			return it.rotateRight(node)
		}
	}

	if node.Right != nil {
		nodeKey, err := it.KeyFunc(node.Right.FormatOp)
		if err != nil {
			return nil, err
		}

		// Right Right Case
		if balance < -1 && opKey > nodeKey {
			return it.rotateLeft(node)
		}

		// Right Left Case
		if balance < -1 && opKey < nodeKey {
			node.Right, err = it.rotateRight(node.Right)
			if err != nil {
				return nil, err
			}

			return it.rotateLeft(node)
		}
	}

	return node, nil
}

func (it *IntervalTree) _nodeSpan(node *IntervalNode) (int, int, error) {
	_, startIx, err := it.Rope.GetIndex(node.FormatOp.StartID)
	if err != nil {
		return 0, 0, err
	}
	_, endIx, err := it.Rope.GetIndex(node.FormatOp.EndID)
	if err != nil {
		return 0, 0, err
	}

	if it.IsSticky {
		endIx--
	}

	return startIx, endIx, nil
}

func (it *IntervalTree) _nodeOverlaps(node *IntervalNode, startIx int, endIx int) (bool, error) {
	nodeStartIx, nodeEndIx, err := it._nodeSpan(node)
	return startIx <= nodeEndIx && endIx >= nodeStartIx, err
}

func (it *IntervalTree) SearchOverlapping(startIx, endIx int) ([]FormatOp, error) {
	result := make([]FormatOp, 0)

	err := it._searchOverlapping(it.Root, startIx, endIx, &result)
	if err != nil {
		return nil, err
	}

	slices.SortFunc(result, func(a, b FormatOp) int {
		aSeq, aAuthor := a.ID.Seq, a.ID.Author
		bSeq, bAuthor := b.ID.Seq, b.ID.Author

		if aSeq < bSeq {
			return -1
		} else if aSeq > bSeq {
			return 1
		} else if aAuthor < bAuthor {
			return -1
		} else if aAuthor > bAuthor {
			return 1
		} else {
			return 0
		}
	})

	/*
	   // DEBUG
	   fmt.Printf("OVERLAPPING startIx: %d endIx: %d\n", startIx, endIx)
	   for _, op := range result {
	     it.Rope._printTotIxFormatOp(op)
	   }
	*/

	return result, nil
}

func (it *IntervalTree) _searchOverlapping(node *IntervalNode, startIx, endIx int, result *[]FormatOp) error {
	if node == nil {
		return nil
	}

	doesOverlap, err := it._nodeOverlaps(node, startIx, endIx)
	if err != nil {
		return err
	}

	if doesOverlap {
		*result = append(*result, node.FormatOp)
	}

	// recurse left if needed
	if node.Left != nil {
		_, leftMaxIx, err := it.Rope.GetIndex(node.Left.MaxID)
		if err != nil {
			return err
		}

		if startIx <= leftMaxIx {
			err = it._searchOverlapping(node.Left, startIx, endIx, result)
			if err != nil {
				return err
			}
		}
	}

	// recurse right if needed
	_, nodeStartIx, err := it.Rope.GetIndex(node.StartID)
	if err != nil {
		return err
	}

	if nodeStartIx <= endIx {
		err = it._searchOverlapping(node.Right, startIx, endIx, result)
		if err != nil {
			return err
		}
	}

	return nil
}

func (node *IntervalNode) Dft(callback func(*IntervalNode) error) error {
	if node == nil {
		return nil
	}

	node.Left.Dft(callback)
	err := callback(node)
	if err != nil {
		return err
	}
	node.Right.Dft(callback)

	return nil
}

func (root *IntervalNode) Bfs(callback func(*IntervalNode) error) error {
	if root == nil {
		return nil
	}

	queue := []*IntervalNode{root}

	for len(queue) > 0 {
		currentNode := queue[0]
		queue = queue[1:]

		err := callback(currentNode)
		if err != nil {
			return err
		}

		if currentNode.Left != nil {
			queue = append(queue, currentNode.Left)
		}

		if currentNode.Right != nil {
			queue = append(queue, currentNode.Right)
		}
	}

	return nil
}

func (node *IntervalNode) Height() int {
	if node == nil {
		return -1
	}
	return node.height
}

func (node *IntervalNode) getBalance() int {
	if node == nil {
		return 0
	}
	return node.Left.Height() - node.Right.Height()
}

func (it *IntervalTree) rotateRight(node *IntervalNode) (*IntervalNode, error) {
	if node == nil {
		return nil, nil
	}

	newRoot := node.Left
	if newRoot == nil {
		return node, nil
	}

	// Perform rotation
	node.Left = newRoot.Right
	newRoot.Right = node

	// Update heights
	node.height = max(node.Left.Height(), node.Right.Height()) + 1
	newRoot.height = max(newRoot.Left.Height(), newRoot.Right.Height()) + 1

	err := it.setMinMaxIDs(node)
	if err != nil {
		return nil, err
	}

	err = it.setMinMaxIDs(newRoot)
	if err != nil {
		return nil, err
	}

	// Return new root
	return newRoot, nil
}

func (it *IntervalTree) rotateLeft(node *IntervalNode) (*IntervalNode, error) {
	if node == nil {
		return nil, nil
	}

	newRoot := node.Right
	if newRoot == nil {
		return node, nil
	}

	// Perform rotation
	node.Right = newRoot.Left
	newRoot.Left = node

	// Update heights
	node.height = max(node.Left.Height(), node.Right.Height()) + 1
	newRoot.height = max(newRoot.Left.Height(), newRoot.Right.Height()) + 1

	err := it.setMinMaxIDs(node)
	if err != nil {
		return nil, err
	}

	err = it.setMinMaxIDs(newRoot)
	if err != nil {
		return nil, err
	}

	// Return new root
	return newRoot, nil
}

func (node *IntervalNode) IsBalanced() (bool, int) {
	if node == nil {
		return true, -1
	}

	leftBalanced, leftHeight := node.Left.IsBalanced()
	if !leftBalanced {
		return false, 0
	}

	rightBalanced, rightHeight := node.Right.IsBalanced()
	if !rightBalanced {
		return false, 0
	}

	if abs(leftHeight-rightHeight) > 1 {
		return false, 0
	}

	height := max(leftHeight, rightHeight) + 1
	return true, height
}

func (node *IntervalNode) updateHeight() {
	if node == nil {
		return
	}
	node.Left.updateHeight()
	node.Right.updateHeight()
	node.height = 1 + max(node.Left.Height(), node.Right.Height())
}

func (it *IntervalTree) Print() {
	it.Root.Bfs(func(node *IntervalNode) error {
		fmt.Printf("%v\n", node.FormatOp)
		return nil
	})
}

func (it *IntervalTree) Validate() error {
	err := it.Root.Dft(func(node *IntervalNode) error {
		_, startIx, err := it.Rope.GetIndex(node.FormatOp.StartID)
		if err != nil {
			return fmt.Errorf("GetIndex(%v): %w", node.FormatOp.StartID, err)
		}

		_, endIx, err := it.Rope.GetIndex(node.FormatOp.EndID)
		if err != nil {
			return fmt.Errorf("GetIndex(%v): %w", node.FormatOp.EndID, err)
		}

		if it.IsSticky {
			endIx--
		}

		if endIx < startIx {
			return fmt.Errorf("endIx < startIx: %d < %d", endIx, startIx)
		}

		if f, ok := node.FormatOp.Format.(FormatV3Span); ok {
			sticky, noSticky := f.SplitSticky()

			if it.IsSticky && len(noSticky) > 0 {
				return fmt.Errorf("noSticky span found in sticky: %v", node.FormatOp)
			}

			if !it.IsSticky && len(sticky) > 0 {
				return fmt.Errorf("sticky span found in noSticky: %v", node.FormatOp)
			}
		} else {
			return fmt.Errorf("formatOp %v is not a FormatV3Span", node.FormatOp)
		}

		return nil
	})

	return err
}

func maxID(ids ...ID) ID {
	maxID := ids[0]
	for _, id := range ids {
		if maxID.lessThan(id) {
			maxID = id
		}
	}
	return maxID
}

func minID(ids ...ID) ID {
	minID := ids[0]
	for _, id := range ids {
		if id.lessThan(minID) {
			minID = id
		}
	}
	return minID
}

func (it *IntervalTree) setMinMaxIDs(node *IntervalNode) error {
	maxIDs, maxIxs, minSeqIDs, maxSeqIDs := make([]ID, 0, 3), make([]int, 0, 3), make([]ID, 0, 3), make([]ID, 0, 3)

	_, maxIx, err := it.Rope.GetIndex(node.FormatOp.EndID)
	if err != nil {
		return err
	}
	maxIxs = append(maxIxs, maxIx)
	maxIDs = append(maxIDs, node.FormatOp.EndID)

	minSeqIDs = append(minSeqIDs, node.FormatOp.ID)
	maxSeqIDs = append(maxSeqIDs, node.FormatOp.ID)

	if node.Left != nil {
		_, maxLeftIx, err := it.Rope.GetIndex(node.Left.MaxID)
		if err != nil {
			return err
		}
		maxIxs = append(maxIxs, maxLeftIx)
		maxIDs = append(maxIDs, node.Left.MaxID)

		minSeqIDs = append(minSeqIDs, node.Left.MinSeqID)
		maxSeqIDs = append(maxSeqIDs, node.Left.MaxSeqID)
	}

	if node.Right != nil {
		_, maxRightIx, err := it.Rope.GetIndex(node.Right.MaxID)
		if err != nil {
			return err
		}
		maxIxs = append(maxIxs, maxRightIx)
		maxIDs = append(maxIDs, node.Right.MaxID)

		minSeqIDs = append(minSeqIDs, node.Right.MinSeqID)
		maxSeqIDs = append(maxSeqIDs, node.Right.MaxSeqID)
	}

	ix := SliceMaxIx(maxIxs)

	node.MaxID = maxIDs[ix]
	node.MaxSeqID = maxID(maxSeqIDs...)
	node.MinSeqID = minID(minSeqIDs...)

	return nil
}

func (it *IntervalTree) MaxFormatAt(startIx, endIx int, address *ContentAddress) (*FormatOp, error) {
	var maxNode *IntervalNode
	var recurse func(node *IntervalNode) error

	recurse = func(node *IntervalNode) error {
		if node == nil {
			return nil
		}

		overlaps, err := it._nodeOverlaps(node, startIx, endIx)
		if err != nil {
			return err
		}

		if overlaps && (address == nil || address.Contains(node.FormatOp.ID)) {
			if maxNode == nil || maxNode.FormatOp.ID.lessThan(node.FormatOp.ID) {
				maxNode = node
			}
		}

		if node.Right != nil {
			_, nodeStartIx, err := it.Rope.GetIndex(node.StartID)
			if err != nil {
				return err
			}

			addressContainsTree := address == nil || address.Contains(node.Right.MinSeqID)
			treeGreaterThanCandidate := maxNode == nil || maxNode.FormatOp.ID.lessThan(node.Right.MaxSeqID)

			if nodeStartIx <= endIx && addressContainsTree && treeGreaterThanCandidate {
				err = recurse(node.Right)
				if err != nil {
					return err
				}
			}
		}

		if node.Left != nil {
			_, leftMaxIx, err := it.Rope.GetIndex(node.Left.MaxID)
			if err != nil {
				return err
			}

			addressContainsTree := address == nil || address.Contains(node.Left.MinSeqID)
			treeGreaterThanCandidate := maxNode == nil || maxNode.FormatOp.ID.lessThan(node.Left.MaxSeqID)

			if startIx <= leftMaxIx && addressContainsTree && treeGreaterThanCandidate {
				err = recurse(node.Left)
				if err != nil {
					return err
				}
			}
		}

		return nil
	}

	err := recurse(it.Root)
	if err != nil {
		return nil, err
	}

	if maxNode == nil {
		return nil, nil
	}

	return &maxNode.FormatOp, nil
}

func (it *IntervalTree) MaxFormatBefore(startIx, endIx int, targetID ID) (*FormatOp, error) {
	var maxNode *IntervalNode
	var recurse func(node *IntervalNode) error

	recurse = func(node *IntervalNode) error {
		if node == nil {
			return nil
		}

		overlaps, err := it._nodeOverlaps(node, startIx, endIx)
		if err != nil {
			return err
		}

		if overlaps && node.FormatOp.ID.lessThan(targetID) {
			if maxNode == nil || maxNode.FormatOp.ID.lessThan(node.FormatOp.ID) {
				maxNode = node
			}
		}

		if node.Right != nil {
			_, nodeStartIx, err := it.Rope.GetIndex(node.StartID)
			if err != nil {
				return err
			}

			opLessThanTarget := node.Right.MinSeqID.lessThan(targetID)
			treeGreaterThanCandidate := maxNode == nil || maxNode.FormatOp.ID.lessThan(node.Right.MaxSeqID)

			if nodeStartIx <= endIx && opLessThanTarget && treeGreaterThanCandidate {
				err = recurse(node.Right)
				if err != nil {
					return err
				}
			}
		}

		if node.Left != nil {
			_, leftMaxIx, err := it.Rope.GetIndex(node.Left.MaxID)
			if err != nil {
				return err
			}

			opLessThanTarget := node.Left.MinSeqID.lessThan(targetID)
			treeGreaterThanCandidate := maxNode == nil || maxNode.FormatOp.ID.lessThan(node.Left.MaxSeqID)

			if startIx <= leftMaxIx && opLessThanTarget && treeGreaterThanCandidate {
				err = recurse(node.Left)
				if err != nil {
					return err
				}
			}
		}

		return nil
	}

	err := recurse(it.Root)
	if err != nil {
		return nil, err
	}

	if maxNode == nil {
		return nil, nil
	}

	return &maxNode.FormatOp, nil
}
