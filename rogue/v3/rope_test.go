package v3

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRope_GetNode(t *testing.T) {
	r := NewRogueForQuill("auth0")

	for i, c := range "Hello World!" {
		op, err := r.Insert(i, string(c))
		if err != nil {
			t.Fatal(fmt.Errorf("op: %v, err: %v", op, err))
		}
	}

	if r.GetText() != "Hello World!\n" {
		t.Errorf("Expected 'Hello World!', got %q", r.GetText())
	}

	i, node, err := r.Rope.GetNode(11)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 11, i)
	assert.Equal(t, "!", string(rune(node.Val.Text[i])))

	r = NewRogueForQuill("auth0")
	op, err := r.Insert(0, "Hello World!")
	if err != nil {
		t.Fatal(fmt.Errorf("op: %v, err: %v", op, err))
	}

	if r.GetText() != "Hello World!\n" {
		t.Errorf("Expected 'Hello World!', got %q", r.GetText())
	}

	i, node, err = r.Rope.GetNode(11)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 11, i)
	assert.Equal(t, "!", string(rune(node.Val.Text[i])))
}

func TestRope_GetId(t *testing.T) {
	r := NewRogueForQuill("auth0")

	// ops := []Op{}
	for i, c := range "Hello World!" {
		op, err := r.Insert(i, string(c))
		if err != nil {
			t.Fatal(fmt.Errorf("op: %v, err: %v", op, err))
		}
		// ops = append(ops, op)
	}

	if r.GetText() != "Hello World!\n" {
		t.Errorf("Expected 'Hello World!', got %q", r.GetText())
	}

	id, err := r.Rope.GetVisID(0)
	assert.NoError(t, err)
	assert.Equal(t, 3, id.Seq) // it's 3 because 0 + (root node) + (end node) + (inserted \n for quill)

	id, err = r.Rope.GetVisID(11)
	assert.NoError(t, err)
	assert.Equal(t, 14, id.Seq) // it's 14 because 11 + (root node) + (end node) + (inserted \n for quill)
}

func TestRope_GetBetween(t *testing.T) {
	r := NewRogueForQuill("auth0")

	for i, c := range "Hello World!" {
		op, err := r.Insert(i, string(c))
		if err != nil {
			t.Fatal(fmt.Errorf("op: %v, err: %v", op, err))
		}
	}

	if r.GetText() != "Hello World!\n" {
		t.Errorf("Expected 'Hello World!', got %q", r.GetText())
	}

	start, err := r.Rope.GetVisID(0)
	assert.NoError(t, err)
	end, err := r.Rope.GetVisID(11)
	assert.NoError(t, err)

	text, err := r.Rope.GetBetween(start, end)
	assert.NoError(t, err)
	assert.Equal(t, "Hello World!", Uint16ToStr(text.Text))

	r = NewRogueForQuill("auth0")
	op, err := r.Insert(0, "Hello World!")
	if err != nil {
		t.Fatal(fmt.Errorf("op: %v, err: %v", op, err))
	}

	if r.GetText() != "Hello World!\n" {
		t.Errorf("Expected 'Hello World!', got %q", r.GetText())
	}

	start, err = r.Rope.GetVisID(0)
	assert.NoError(t, err)
	end, err = r.Rope.GetVisID(11)
	assert.NoError(t, err)

	text, err = r.Rope.GetBetween(start, end)
	assert.NoError(t, err)
	assert.Equal(t, "Hello World!", Uint16ToStr(text.Text))
}

func TestRope_GetBetween_starting_at_root(t *testing.T) {
	r := NewRogueForQuill("auth0")
	op, err := r.Insert(0, "Hello World!")
	if err != nil {
		t.Fatal(fmt.Errorf("op: %v, err: %v", op, err))
	}

	if r.GetText() != "Hello World!\n" {
		t.Errorf("Expected 'Hello World!', got %q", r.GetText())
	}

	end, err := r.Rope.GetVisID(11)
	assert.NoError(t, err)

	text, err := r.Rope.GetBetween(RootID, end)
	assert.NoError(t, err)
	assert.Equal(t, "Hello World!", Uint16ToStr(text.Text))
}

func TestRope_GetBetweenTable(t *testing.T) {
	type testCase struct {
		name        string
		setup       func(*testing.T) *Rogue
		startID     ID
		endID       ID
		expectedVis string
		expectedTot string
	}

	tests := []testCase{
		{
			name: "hello world",
			setup: func(t *testing.T) *Rogue {
				r := NewRogueForQuill("auth0")
				_, err := r.Insert(0, "Hello World!")
				require.NoError(t, err)
				return r
			},
			startID:     RootID,
			endID:       LastID,
			expectedVis: "Hello World!\n",
			expectedTot: "xHello World!\n",
		},
		{
			name: "hello cruel world",
			setup: func(t *testing.T) *Rogue {
				r := NewRogueForQuill("auth0")
				_, err := r.Insert(0, "Hello World!")
				require.NoError(t, err)

				_, err = r.Insert(5, " cruel")
				require.NoError(t, err)

				return r
			},
			startID:     ID{"auth0", 7},
			endID:       ID{"auth0", 9},
			expectedVis: "o cruel W",
			expectedTot: "o cruel W",
		},
		{
			name: "hello deleted world",
			setup: func(t *testing.T) *Rogue {
				r := NewRogueForQuill("auth0")
				_, err := r.Insert(0, "Hello World!")
				require.NoError(t, err)

				_, err = r.Insert(5, " cruel")
				require.NoError(t, err)

				for i := 0; i < len(" cruel"); i++ {
					_, err = r.Delete(5, 1)
					require.NoError(t, err)
				}

				return r
			},
			startID:     ID{"auth0", 7},
			endID:       ID{"auth0", 9},
			expectedVis: "o W",
			expectedTot: "o cruel W",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := tc.setup(t)

			vis, err := r.Rope.GetBetween(tc.startID, tc.endID)
			require.NoError(t, err)

			tot, err := r.Rope.GetTotBetween(tc.startID, tc.endID)
			require.NoError(t, err)

			require.Equal(t, tc.expectedVis, Uint16ToStr(vis.Text))
			require.Equal(t, vis.Text, tot.Visible().Text)
			require.Equal(t, tc.expectedTot, Uint16ToStr(tot.Text))
		})
	}
}

func TestRope_LeftTotSibling(t *testing.T) {
	r := NewRogueForQuill("auth0")
	op, err := r.Insert(0, "Hello")
	if err != nil {
		t.Fatal(fmt.Errorf("op: %v, err: %v", op, err))
	}

	r.Rope.Print()

	endId, err := r.Rope.GetVisID(5)
	assert.NoError(t, err)

	end := r.Rope.Index.Get(endId)
	left, err := end.LeftTotSibling()
	assert.NoError(t, err)
	if left == nil {
		t.Fatal("left is nil")
	}
	assert.Equal(t, "auth0_3", left.Val.ID.String())
}

func TestRopeGetIndex(t *testing.T) {
	t.Parallel()
	type testCase struct {
		name        string
		initialDoc  func(r *Rogue)
		id          ID
		expectedVis int
		expectedTot int
		expectedErr error
	}

	testCases := []testCase{
		{
			name: "Root!",
			initialDoc: func(r *Rogue) {
				_, err := r.Insert(0, "Hello World!")
				require.NoError(t, err)
			},
			id:          ID{Author: "root", Seq: 0},
			expectedVis: -1,
			expectedTot: 0,
		},
		{
			name: "Q2",
			initialDoc: func(r *Rogue) {
				_, err := r.Insert(0, "Hello World!")
				require.NoError(t, err)
			},
			id:          ID{Author: "q", Seq: 1},
			expectedVis: 12,
			expectedTot: 13,
		},
		{
			name: "Deleted char",
			initialDoc: func(r *Rogue) {
				_, err := r.Insert(0, "Hello World!")
				require.NoError(t, err)
				_, err = r.Delete(4, 1)
				require.NoError(t, err)
			},
			id:          ID{Author: "auth0", Seq: 7},
			expectedVis: -1,
			expectedTot: 5,
		},
	}

	// Execute each test case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := NewRogueForQuill("auth0")
			tc.initialDoc(r)
			vis, tot, err := r.Rope.GetIndex(tc.id)
			assert.Equal(t, tc.expectedVis, vis)
			assert.Equal(t, tc.expectedTot, tot)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestContainsVisibleBetween(t *testing.T) {
	r := NewRogue("0")

	_, err := r.Insert(0, "Hello cruel World!")
	require.NoError(t, err)

	_, err = r.Delete(6, 6)
	require.NoError(t, err)

	isVis, err := r.Rope.ContainsVisibleBetween(ID{"0", 6}, ID{"0", 11})
	require.NoError(t, err)
	require.False(t, isVis)

	isVis, err = r.Rope.ContainsVisibleBetween(ID{"0", 0}, ID{"0", 1})
	require.NoError(t, err)
	require.False(t, isVis)

	isVis, err = r.Rope.ContainsVisibleBetween(ID{"0", 5}, ID{"0", 12})
	require.NoError(t, err)
	require.False(t, isVis)

	isVis, err = r.Rope.ContainsVisibleBetween(ID{"0", 0}, ID{"0", 2})
	require.NoError(t, err)
	require.True(t, isVis)

	isVis, err = r.Rope.ContainsVisibleBetween(ID{"0", 4}, ID{"0", 12})
	require.NoError(t, err)
	require.True(t, isVis)

	isVis, err = r.Rope.ContainsVisibleBetween(ID{"0", 5}, ID{"0", 13})
	require.NoError(t, err)
	require.True(t, isVis)

	// IDs are the same
	isVis, err = r.Rope.ContainsVisibleBetween(ID{"0", 5}, ID{"0", 5})
	require.NoError(t, err)
	require.False(t, isVis)

	// Ids out of order
	isVis, err = r.Rope.ContainsVisibleBetween(ID{"0", 10}, ID{"0", 5})
	require.NoError(t, err)
	require.False(t, isVis)
}

func TestContainsVisible(t *testing.T) {
	r := NewRogue("0")

	_, err := r.Insert(0, "Hello cruel World!")
	require.NoError(t, err)

	_, err = r.Delete(6, 6)
	require.NoError(t, err)

	isVis, err := r.Rope.ContainsVisible(ID{"0", 6}, ID{"0", 11})
	require.NoError(t, err)
	require.False(t, isVis)

	isVis, err = r.Rope.ContainsVisible(ID{"0", 0}, ID{"0", 1})
	require.NoError(t, err)
	require.True(t, isVis)

	isVis, err = r.Rope.ContainsVisible(ID{"0", 6}, ID{"0", 12})
	require.NoError(t, err)
	require.True(t, isVis)

	isVis, err = r.Rope.ContainsVisible(ID{"0", 5}, ID{"0", 11})
	require.NoError(t, err)
	require.True(t, isVis)

	isVis, err = r.Rope.ContainsVisible(ID{"0", 0}, ID{"0", 2})
	require.NoError(t, err)
	require.True(t, isVis)

	isVis, err = r.Rope.ContainsVisible(ID{"0", 4}, ID{"0", 12})
	require.NoError(t, err)
	require.True(t, isVis)

	isVis, err = r.Rope.ContainsVisible(ID{"0", 5}, ID{"0", 13})
	require.NoError(t, err)
	require.True(t, isVis)

	// IDs are the same
	isVis, err = r.Rope.ContainsVisible(ID{"0", 5}, ID{"0", 5})
	require.NoError(t, err)
	require.True(t, isVis)

	// Ids out of order
	isVis, err = r.Rope.ContainsVisible(ID{"0", 10}, ID{"0", 5})
	require.NoError(t, err)
	require.False(t, isVis)
}

func TestGetTotID(t *testing.T) {
	r := NewRogueForQuill("0")

	_, err := r.Insert(0, "Hello World!")
	require.NoError(t, err)

	id, err := r.Rope.GetTotID(13)
	require.NoError(t, err)

	require.Equal(t, id, LastID)
}
