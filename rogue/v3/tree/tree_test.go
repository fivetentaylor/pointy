package tree_test

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
	v3 "github.com/teamreviso/code/rogue/v3"
	"github.com/teamreviso/code/rogue/v3/tree"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

func id[T any](x T) (T, error) {
	return x, nil
}

func lessThan[T constraints.Ordered](a, b T) bool {
	return a < b
}

func TestAVLTree(t *testing.T) {
	tests := []struct {
		name   string
		inputs []int
		check  int
		exists bool
		delete int
		result []int
	}{
		{
			name:   "Insert multiple items",
			inputs: []int{33, 13, 53, 9, 21, 61, 8, 11},
			check:  21,
			exists: true,
			delete: 13,
			result: []int{8, 9, 11, 21, 33, 53, 61},
		},
		{
			name:   "Insert and delete root",
			inputs: []int{20, 4, 26, 3, 9, 15},
			check:  15,
			exists: true,
			delete: 20,
			result: []int{3, 4, 9, 15, 26},
		},
		{
			name:   "Left heavy operations",
			inputs: []int{30, 20, 40, 10, 25, 5, 15, 1},
			check:  5,
			exists: true,
			delete: 40,
			result: []int{1, 5, 10, 15, 20, 25, 30},
		},
		{
			name:   "Right heavy operations",
			inputs: []int{30, 20, 40, 35, 50, 45, 60},
			check:  45,
			exists: true,
			delete: 20,
			result: []int{30, 35, 40, 45, 50, 60},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			avl := tree.NewTree(id[int], lessThan)
			for _, input := range tt.inputs {
				err := avl.Put(input)
				require.NoError(t, err)

				avl.Print()
			}

			err := avl.Validate()
			require.NoError(t, err)

			gotNode, err := avl.Get(tt.check)
			require.NoError(t, err)
			if tt.exists {
				require.Equal(t, tt.check, gotNode)
			} else {
				require.Nil(t, gotNode)
			}

			avl.Print()
			err = avl.Remove(tt.delete)
			require.NoError(t, err)
			avl.Print()

			x := make([]int, 0, avl.Size)
			err = avl.Dft(func(n int) error {
				x = append(x, n)
				return nil
			})

			require.NoError(t, err)
			require.Equal(t, tt.result, x)

			if avl.Size > 2 {
				k, err := avl.FindRightSib(x[2])
				require.NoError(t, err)
				require.Equal(t, x[2], k)
				n, _ := avl.GetNode(k)

				y := []int{}
				err = n.WalkRight(func(n int) error {
					y = append(y, n)
					return nil
				})

				require.NoError(t, err)
				require.Equal(t, x[2:], y)
			}

			avl.Print()
		})
	}
}

func idp[T any](x *T) (T, error) {
	return *x, nil
}

func TestSibling(t *testing.T) {
	avl := tree.NewTree(idp[int], lessThan)
	for _, input := range []int{20, 10, 30, 5, 15, 25, 35} {
		i := input
		err := avl.Put(&i)
		require.NoError(t, err)

		avl.Print()
	}

	err := avl.Validate()
	require.NoError(t, err)

	n, err := avl.FindRightSib(11)
	require.NoError(t, err)
	require.Equal(t, 15, *n)

	n, err = avl.FindLeftSib(11)
	require.NoError(t, err)
	require.Equal(t, 10, *n)

	n, err = avl.FindRightSib(35)
	require.NoError(t, err)
	require.Equal(t, 35, *n)

	n, err = avl.FindLeftSib(5)
	require.NoError(t, err)
	require.Equal(t, 5, *n)

	n, err = avl.FindRightSib(36)
	require.NoError(t, err)
	require.Nil(t, n)

	n, err = avl.FindLeftSib(4)
	require.NoError(t, err)
	require.Nil(t, n)

	n, err = avl.FindLeftSib(40)
	require.NoError(t, err)
	require.Equal(t, 35, *n)

	n, err = avl.FindRightSib(0)
	require.NoError(t, err)
	require.Equal(t, 5, *n)

	n, err = avl.FindLeftSib(400)
	require.NoError(t, err)
	require.Equal(t, 35, *n)
}

func TestTreeMin(t *testing.T) {
	avlInt := tree.NewTree(id[int], lessThan)
	for _, input := range []int{20, 10, 30, 5, 15, 25, 35} {
		err := avlInt.Put(input)
		require.NoError(t, err)
	}
	require.Equal(t, 5, avlInt.Min())

	avlStr := tree.NewTree(id[string], lessThan)
	for _, input := range []string{"b", "c", "a"} {
		err := avlStr.Put(input)
		require.NoError(t, err)
	}

	require.Equal(t, "a", avlStr.Min())
}

func TestTreeMax(t *testing.T) {
	avlInt := tree.NewTree(id[int], lessThan)
	for _, input := range []int{20, 10, 30, 5, 15, 25, 35} {
		err := avlInt.Put(input)
		require.NoError(t, err)
	}
	require.Equal(t, 35, avlInt.Max())

	avlStr := tree.NewTree(id[string], lessThan)
	for _, input := range []string{"b", "c", "a"} {
		err := avlStr.Put(input)
		require.NoError(t, err)
	}

	require.Equal(t, "c", avlStr.Max())
}

func TestFuzzy(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	avl := tree.NewTree(id[string], lessThan)

	for i := 0; i < 10000; i++ {
		s := v3.RandomString(rng, 5)
		err := avl.Put(s)
		require.NoError(t, err)
	}

	x := make([]string, 0, avl.Size)
	err := avl.Dft(func(n string) error {
		x = append(x, n)
		return nil
	})
	require.NoError(t, err)
	require.True(t, slices.IsSorted(x))

	lastKey := x[len(x)-1]
	n, _ := avl.FindLeftSibNode(lastKey)
	for i := len(x) - 1; i >= 0; i-- {
		require.Equal(t, x[i], n.Value)
		n = n.StepLeft()
	}

	err = avl.Validate()
	require.NoError(t, err)

	for i := 0; i < 10000; i++ {
		ix := rng.Intn(len(x)) // rand ix
		k := x[ix]
		x = v3.DeleteAt(x, ix, 1)
		err := avl.Remove(k)
		require.NoError(t, err)

		require.Equal(t, avl.Size, len(x))
		require.Equal(t, x, avl.AsSlice())
	}
}

// TestMerge tests the Merge function with multiple scenarios.
func TestMerge(t *testing.T) {
	tests := []struct {
		name     string
		trees    [][]int
		expected []int
	}{
		{
			name: "simple merge without errors",
			trees: [][]int{
				{1, 2, 3},
				{4, 5, 6},
			},
			expected: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name: "merge with overlapping elements",
			trees: [][]int{
				{1, 3, 5},
				{2, 3, 4},
			},
			expected: []int{1, 2, 3, 3, 4, 5},
		},
		{
			name: "merge with different sized trees",
			trees: [][]int{
				{1},
				{2, 3, 4, 5},
			},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name: "merge with an empty tree",
			trees: [][]int{
				{1, 2, 3},
				{},
			},
			expected: []int{1, 2, 3},
		},
		{
			name: "merge with all empty trees",
			trees: [][]int{
				{},
				{},
			},
			expected: []int{},
		},
		{
			name: "merge single tree",
			trees: [][]int{
				{1, 2, 3, 4, 5, 6},
			},
			expected: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name: "merge three trees",
			trees: [][]int{
				{1, 3, 5},
				{2, 4, 6},
				{0, 7, 8},
			},
			expected: []int{0, 1, 2, 3, 4, 5, 6, 7, 8},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := []int{}
			callback := func(n int) error {
				out = append(out, n)
				return nil
			}

			trees := make([]*tree.Tree[int, int], 0, len(tt.trees))
			for _, ints := range tt.trees {
				tree := tree.NewTree(id[int], lessThan)
				trees = append(trees, tree)

				for _, n := range ints {
					err := tree.Put(n)
					if err != nil {
						return
					}
				}
			}

			err := tree.Merge(trees, callback)
			require.NoError(t, err)
			require.Equal(t, tt.expected, out)
		})
	}
}

type ID struct {
	Author string
	Seq    int
}

func keyFunc(id ID) (int, error) {
	return id.Seq, nil
}

func TestOverwrite(t *testing.T) {
	avl := tree.NewTree(keyFunc, lessThan)
	a := ID{"a", 1}
	b := ID{"b", 1}
	c := ID{"c", 2}

	err := avl.Put(a)
	require.NoError(t, err)

	a_, err := avl.Get(a.Seq)
	require.NoError(t, err)
	require.Equal(t, a, a_)

	err = avl.Put(b)
	require.NoError(t, err)

	b_, err := avl.Get(b.Seq)
	require.NoError(t, err)
	require.Equal(t, b, b_)

	err = avl.Put(c)
	require.NoError(t, err)

	c_, err := avl.Get(c.Seq)
	require.NoError(t, err)
	require.Equal(t, c, c_)
}
