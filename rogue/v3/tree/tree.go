package tree

import (
	"encoding/json"
	"fmt"

	heap "github.com/teamreviso/code/rogue/v3/heap"
)

type Node[K any, T any] struct {
	Value  T
	height int
	parent *Node[K, T]
	left   *Node[K, T]
	right  *Node[K, T]
}

type Tree[K any, T any] struct {
	Size     int
	root     *Node[K, T]
	KeyFunc  func(T) (K, error)
	LessThan func(K, K) bool
}

func NewTree[K any, T any](keyFunc func(T) (K, error), lessThan func(K, K) bool) *Tree[K, T] {
	return &Tree[K, T]{
		KeyFunc:  keyFunc,
		LessThan: lessThan,
	}
}

func (n *Node[K, T]) String() string {
	if n == nil {
		return ""
	}
	var leftStr, rightStr string
	if n.left != nil {
		leftStr = n.left.String()
	}
	if n.right != nil {
		rightStr = n.right.String()
	}
	return fmt.Sprintf("%#v (%s, %s)", n.Value, leftStr, rightStr)
}

func (t *Tree[K, T]) String() string {
	if t.root == nil {
		return "Tree is empty"
	}
	return fmt.Sprintf("Tree[%s]", t.root.String())
}

func (n *Node[K, T]) getHeight() int {
	if n == nil {
		return 0
	}
	return n.height
}

func (n *Node[K, T]) rightRotate() *Node[K, T] {
	x := n.left
	T2 := x.right

	x.right = n
	n.left = T2

	x.parent = n.parent
	n.parent = x
	if T2 != nil {
		T2.parent = n
	}

	n.height = 1 + max(n.left.getHeight(), n.right.getHeight())
	x.height = 1 + max(x.left.getHeight(), x.right.getHeight())

	return x
}

func (n *Node[K, T]) leftRotate() *Node[K, T] {
	y := n.right
	T2 := y.left

	y.left = n
	n.right = T2

	y.parent = n.parent
	n.parent = y
	if T2 != nil {
		T2.parent = n
	}

	n.height = 1 + max(n.left.getHeight(), n.right.getHeight())
	y.height = 1 + max(y.left.getHeight(), y.right.getHeight())

	return y
}

func (n *Node[K, T]) GetBalanceFactor() int {
	if n == nil {
		return 0
	}
	return n.left.getHeight() - n.right.getHeight()
}

func (n *Node[K, T]) rebalance() *Node[K, T] {
	n.height = 1 + max(n.left.getHeight(), n.right.getHeight())
	balance := n.GetBalanceFactor()

	if balance > 1 {
		if n.left.GetBalanceFactor() < 0 {
			n.left = n.left.leftRotate()
			if n.left != nil {
				n.left.parent = n
			}
		}
		return n.rightRotate()
	}
	if balance < -1 {
		if n.right.GetBalanceFactor() > 0 {
			n.right = n.right.rightRotate()
			if n.right != nil {
				n.right.parent = n
			}
		}
		return n.leftRotate()
	}
	return n
}

func (tree *Tree[K, T]) Put(value T) error {
	// get the key out here so we don't call it
	// on every recursion
	key, err := tree.KeyFunc(value)
	if err != nil {
		return err
	}

	n, err := tree.GetNode(key)
	if err != nil {
		return err
	}

	root, err := tree.root.put(key, value, tree.KeyFunc, tree.LessThan)
	if err != nil {
		return err
	}

	tree.root = root
	if n == nil {
		tree.Size++
	}

	return nil
}

func (n *Node[K, T]) put(key K, value T, keyFunc func(T) (K, error), lessThan func(K, K) bool) (*Node[K, T], error) {
	if n == nil {
		return &Node[K, T]{Value: value, height: 1}, nil
	}

	k, err := keyFunc(n.Value)
	if err != nil {
		return nil, err
	}

	if lessThan(key, k) {
		if n.left == nil {
			n.left = &Node[K, T]{Value: value, height: 1, parent: n}
		} else {
			left, err := n.left.put(key, value, keyFunc, lessThan)
			if err != nil {
				return nil, err
			}
			n.left = left
		}
	} else if lessThan(k, key) {
		if n.right == nil {
			n.right = &Node[K, T]{Value: value, height: 1, parent: n}
		} else {
			right, err := n.right.put(key, value, keyFunc, lessThan)
			if err != nil {
				return nil, err
			}
			n.right = right
		}
	} else {
		n.Value = value // overwrite the existing value
	}

	n.height = 1 + max(n.left.getHeight(), n.right.getHeight())

	return n.rebalance(), nil
}

func (tree *Tree[K, T]) Get(key K) (T, error) {
	var zero T

	n, err := tree.root.get(key, tree.KeyFunc, tree.LessThan)
	if err != nil {
		return zero, err
	}

	if n == nil {
		return zero, nil
	}

	return n.Value, nil
}

func (tree *Tree[K, T]) GetNode(key K) (*Node[K, T], error) {
	return tree.root.get(key, tree.KeyFunc, tree.LessThan)
}

func (n *Node[K, T]) get(key K, keyFunc func(T) (K, error), lessThan func(K, K) bool) (*Node[K, T], error) {
	if n == nil {
		return nil, nil
	}

	k, err := keyFunc(n.Value)
	if err != nil {
		return nil, err
	}

	if lessThan(key, k) {
		return n.left.get(key, keyFunc, lessThan)
	} else if lessThan(k, key) {
		return n.right.get(key, keyFunc, lessThan)
	} else {
		return n, nil // value found
	}
}

func (tree *Tree[K, T]) RemoveByKey(key K) error {
	n, err := tree.GetNode(key)
	if err != nil {
		return err
	}

	if n == nil {
		return nil
	}

	root, err := tree.root.remove(key, tree.KeyFunc, tree.LessThan)
	if err != nil {
		return err
	}

	tree.root = root
	tree.Size--

	return nil
}

func (tree *Tree[K, T]) Remove(value T) error {
	key, err := tree.KeyFunc(value)
	if err != nil {
		return err
	}

	return tree.RemoveByKey(key)
}

func (n *Node[K, T]) remove(key K, keyFunc func(T) (K, error), lessThan func(K, K) bool) (*Node[K, T], error) {
	if n == nil {
		return nil, nil
	}

	k, err := keyFunc(n.Value)
	if err != nil {
		return nil, err
	}

	if lessThan(key, k) {
		left, err := n.left.remove(key, keyFunc, lessThan)
		if err != nil {
			return nil, err
		}
		n.left = left
		if n.left != nil {
			n.left.parent = n
		}
	} else if lessThan(k, key) {
		right, err := n.right.remove(key, keyFunc, lessThan)
		if err != nil {
			return nil, err
		}
		n.right = right
		if n.right != nil {
			n.right.parent = n
		}
	} else {
		if n.left == nil {
			if n.right != nil {
				n.right.parent = n.parent
			}
			return n.right, nil
		} else if n.right == nil {
			if n.left != nil {
				n.left.parent = n.parent
			}
			return n.left, nil
		}

		// Node with two children
		successor := leftmost(n.right)
		successorKey, err := keyFunc(successor.Value) // get the key of the successor
		if err != nil {
			return nil, err
		}
		n.Value = successor.Value
		right, err := n.right.remove(successorKey, keyFunc, lessThan) // use the correct key for removal
		if err != nil {
			return nil, err
		}
		n.right = right
		if n.right != nil {
			n.right.parent = n
		}
	}

	n.height = 1 + max(n.left.getHeight(), n.right.getHeight())
	return n.rebalance(), nil
}

func leftmost[K any, T any](n *Node[K, T]) *Node[K, T] {
	current := n
	for current != nil && current.left != nil {
		current = current.left
	}
	return current
}

func rightmost[K any, T any](n *Node[K, T]) *Node[K, T] {
	current := n
	for current != nil && current.right != nil {
		current = current.right
	}
	return current
}

func (tree *Tree[K, T]) Print() {
	tree.root.Dft(func(value T) error {
		fmt.Printf("%v | ", value)
		return nil
	})
	fmt.Println()
}

func (tree *Tree[K, T]) Dft(callback func(T) error) error {
	return tree.root.Dft(callback)
}

func (n *Node[K, T]) Dft(callback func(T) error) error {
	if n == nil {
		return nil
	}

	err := n.left.Dft(callback)
	if err != nil {
		return err
	}

	err = callback(n.Value)
	if err != nil {
		return err
	}

	err = n.right.Dft(callback)
	if err != nil {
		return err
	}

	return nil
}

func (tree *Tree[K, T]) ReverseDft(callback func(T) error) error {
	return tree.root.ReverseDft(callback)
}

func (n *Node[K, T]) ReverseDft(callback func(T) error) error {
	if n == nil {
		return nil
	}

	err := n.right.ReverseDft(callback)
	if err != nil {
		return err
	}

	err = callback(n.Value)
	if err != nil {
		return err
	}

	err = n.left.ReverseDft(callback)
	if err != nil {
		return err
	}

	return nil
}

func (tree *Tree[K, T]) AsSlice() []T {
	slice := make([]T, 0, tree.Size)
	tree.Dft(func(value T) error {
		slice = append(slice, value)
		return nil
	})

	return slice
}

func (tree *Tree[K, T]) FindRightSib(key K) (T, error) {
	var zero T

	rightSib, err := tree.FindRightSibNode(key)
	if err != nil {
		return zero, err
	}

	if rightSib == nil {
		return zero, nil
	}

	return rightSib.Value, nil
}

func (tree *Tree[K, T]) FindRightSibNode(key K) (*Node[K, T], error) {
	n := tree.root
	var rightNeighbor *Node[K, T] = nil

	for n != nil {
		k, err := tree.KeyFunc(n.Value)
		if err != nil {
			return nil, err
		}

		if tree.LessThan(k, key) {
			n = n.right
		} else if tree.LessThan(key, k) {
			rightNeighbor = n
			n = n.left
		} else {
			return n, nil
		}
	}

	return rightNeighbor, nil
}

func (tree *Tree[K, T]) FindLeftSib(key K) (T, error) {
	var zero T

	leftSib, err := tree.FindLeftSibNode(key)
	if err != nil {
		return zero, err
	}

	if leftSib == nil {
		return zero, nil
	}

	return leftSib.Value, nil
}

func (tree *Tree[K, T]) FindLeftSibNode(key K) (*Node[K, T], error) {
	n := tree.root
	var leftNeighbor *Node[K, T] = nil

	for n != nil {
		k, err := tree.KeyFunc(n.Value)
		if err != nil {
			return nil, err
		}

		if tree.LessThan(k, key) {
			leftNeighbor = n
			n = n.right
		} else if tree.LessThan(key, k) {
			n = n.left
		} else {
			return n, nil
		}
	}

	return leftNeighbor, nil
}

func (startNode *Node[K, T]) WalkRight(callback func(T) error) error {
	if startNode == nil {
		return nil
	}

	err := callback(startNode.Value)
	if err != nil {
		return err
	}

	if startNode.right != nil {
		err := startNode.right.Dft(callback)
		if err != nil {
			return err
		}
	}

	node := startNode
	for node.parent != nil {
		if node.parent.left == node {
			err := callback(node.parent.Value)
			if err != nil {
				return err
			}

			err = node.parent.right.Dft(callback)
			if err != nil {
				return err
			}
		}

		node = node.parent
	}

	return nil
}

func (startNode *Node[K, T]) StepLeft() *Node[K, T] {
	if startNode == nil {
		return nil
	}

	if startNode.left != nil {
		return rightmost(startNode.left)
	}

	node := startNode
	for node.parent != nil {
		if node.parent.right == node {
			return node.parent
		}
		node = node.parent
	}

	return nil
}

func (startNode *Node[K, T]) StepRight() *Node[K, T] {
	if startNode == nil {
		return nil
	}

	if startNode.right != nil {
		return leftmost(startNode.right)
	}

	node := startNode
	for node.parent != nil {
		if node.parent.left == node {
			return node.parent
		}
		node = node.parent
	}

	return nil
}

func (root *Node[K, T]) validateParents() error {
	if root.left != nil {
		if root.left.parent != root {
			return fmt.Errorf("parent of left child is not the root")
		}
		err := root.left.validateParents()
		if err != nil {
			return err
		}
	}

	if root.right != nil {
		if root.right.parent != root {
			return fmt.Errorf("parent of right child is not the root")
		}
		err := root.right.validateParents()
		if err != nil {
			return err
		}
	}

	return nil
}

func (tree Tree[K, T]) Validate() error {
	return tree.root.validateParents()
}

func (tree Tree[K, T]) Max() T {
	var zero T
	if tree.root == nil {
		return zero
	}

	return rightmost(tree.root).Value
}

func (tree Tree[K, T]) Min() T {
	var zero T
	if tree.root == nil {
		return zero
	}

	return leftmost(tree.root).Value
}

func (tree Tree[K, T]) MaxNode() *Node[K, T] {
	return rightmost(tree.root)
}

func (tree Tree[K, T]) MinNode() *Node[K, T] {
	return leftmost(tree.root)
}

// function to iterate over multiple trees in order
// note that they must all have the same key type
func Merge[K any, T any](trees []*Tree[K, T], callback func(T) error) error {
	minHeap := heap.NewMinHeap(func(node *Node[K, T]) K {
		k, _ := trees[0].KeyFunc(node.Value)
		return k
	}, trees[0].LessThan)

	for _, tree := range trees {
		if tree.Size > 0 {
			minHeap.Push(tree.MinNode())
		}
	}

	for minHeap.Size > 0 {
		node, ok := minHeap.Pop()
		if !ok {
			return fmt.Errorf("failed to pop from min heap")
		}

		err := callback(node.Value)
		if err != nil {
			return err
		}

		node = node.StepRight()
		if node != nil {
			minHeap.Push(node)
		}
	}

	return nil
}

type treeKeyFunc[K any, T any] func(T) (K, error)

func (tree *Tree[K, T]) Slice(startKey, endKey K, callback func(T) error) error {
	n, err := tree.FindRightSibNode(startKey)
	if err != nil {
		return err
	}

	for n != nil {
		k, err := tree.KeyFunc(n.Value)
		if err != nil {
			return err
		}

		if tree.LessThan(endKey, k) {
			break
		}

		err = callback(n.Value)
		if err != nil {
			return err
		}

		n = n.StepRight()
	}

	return nil
}

func (tree Tree[K, T]) MarshalJSON() ([]byte, error) {
	out := make([]T, 0, tree.Size)
	tree.Dft(func(value T) error {
		out = append(out, value)
		return nil
	})

	return json.Marshal(out)
}

func (tree *Tree[K, T]) UnmarshalJSON(b []byte) error {
	var out []T
	err := json.Unmarshal(b, &out)
	if err != nil {
		return err
	}

	for _, value := range out {
		err := tree.Put(value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (tree *Tree[K, T]) GetAtIx(ix int) (T, error) {
	var zero T
	if ix < 0 || ix >= tree.Size {
		return zero, fmt.Errorf("index out of bounds")
	}

	// TODO: make this faster, tree can optionally be a rope if needed
	var out T
	err := tree.Dft(func(value T) error {
		if ix == 0 {
			out = value
			return fmt.Errorf("stop")
		}
		ix--
		return nil
	})

	return out, err
}
