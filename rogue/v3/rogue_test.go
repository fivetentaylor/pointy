package v3

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"unicode/utf16"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRogueSizes(t *testing.T) {
	testCases := []struct {
		name            string
		fn              func(*Rogue)
		expectedVisSize int
		expectedTotSize int
	}{
		{
			name:            "blank doc",
			fn:              func(r *Rogue) {},
			expectedVisSize: 0,
			expectedTotSize: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := NewRogue("a")
			tc.fn(r)

			assert.Equal(t, tc.expectedVisSize, r.VisSize)
			assert.Equal(t, tc.expectedTotSize, r.TotSize)
		})
	}
}

func TestDelete(t *testing.T) {
	testCases := []struct {
		name         string
		fn           func(*Rogue)
		visIx        int
		len          int
		expectedHtml string
		expectedErr  error
		expectedOps  Op
	}{
		{
			name: "delete single character",
			fn: func(r *Rogue) {
				_, err := r.Insert(0, "abc 123 def 456")
				require.NoError(t, err)
			},
			visIx:        0,
			len:          1,
			expectedHtml: "<p><span>bc 123 def 456</span></p>",
		},
		{
			name: "delete multiple characters",
			fn: func(r *Rogue) {
				_, err := r.Insert(0, "abc 123 def 456")
				require.NoError(t, err)
			},
			visIx:        0,
			len:          4,
			expectedHtml: "<p><span>123 def 456</span></p>",
			expectedOps: DeleteOp{
				ID:         ID{"auth0", 18},
				TargetID:   ID{"auth0", 3},
				SpanLength: 4,
			},
		},
		{
			name: "delete multiple emojis characters",
			fn: func(r *Rogue) {
				_, err := r.Insert(0, "abc 🔥🤓🤘 def 456")
				require.NoError(t, err)
			},
			visIx:        4,
			len:          7,
			expectedHtml: "<p><span>abc def 456</span></p>",
			expectedOps: DeleteOp{
				ID:         ID{"auth0", 21},
				TargetID:   ID{"auth0", 7},
				SpanLength: 7,
			},
		},
		{
			name: "delete over multiple authors",
			fn: func(r *Rogue) {
				_, err := r.Insert(0, "abc def")
				require.NoError(t, err)
				r.Author = "auth1"
				_, err = r.Insert(3, " 123")
				require.NoError(t, err)
				_, err = r.Insert(11, " 456")
				require.NoError(t, err)
				r.Author = "auth0"
			},
			visIx:        0,
			len:          8,
			expectedHtml: "<p><span>def 456</span></p>",
			expectedOps: MultiOp{
				[]Op{
					DeleteOp{
						ID:         ID{"auth0", 18},
						TargetID:   ID{"auth0", 3},
						SpanLength: 3,
					},
					DeleteOp{
						ID:         ID{"auth0", 19},
						TargetID:   ID{"auth1", 10},
						SpanLength: 4,
					},
					DeleteOp{
						ID:         ID{"auth0", 20},
						TargetID:   ID{"auth0", 6},
						SpanLength: 1,
					},
				},
			},
		},
		{
			name: "delete over deletes",
			fn: func(r *Rogue) {
				_, err := r.Insert(0, "hello i'm deleted hi world")
				require.NoError(t, err)
				_, err = r.Delete(6, 11)
				require.NoError(t, err)
			},
			visIx:        6,
			len:          4,
			expectedHtml: "<p><span>hello world</span></p>",
			expectedOps: DeleteOp{
				ID:         ID{"auth0", 30},
				TargetID:   ID{"auth0", 20},
				SpanLength: 4,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := NewRogueForQuill("auth0")
			tc.fn(r)

			ops, err := r.Delete(tc.visIx, tc.len)
			if tc.expectedErr == nil {
				require.NoError(t, err)
			}
			if tc.expectedHtml != "" {
				html, err := r.GetHtml(RootID, LastID, false, false)
				require.NoError(t, err)
				require.Equal(t, tc.expectedHtml, html)
			}
			if tc.expectedOps != nil {
				require.Equal(t, tc.expectedOps, ops)
			}
		})
	}
}

func TestMergeOp(t *testing.T) {
	testCases := []struct {
		name            string
		fn              func(*Rogue)
		Op              Op
		expectedHtml    string
		expectedErr     error
		expectedActions Actions
	}{
		{
			name: "insert hello",
			fn:   func(r *Rogue) {},
			Op: InsertOp{
				ID:       ID{"auth0", 3},
				Text:     "hello",
				ParentID: ID{"root", 0},
				Side:     Right,
			},
			expectedHtml: "<p><span>hello</span></p>",
			expectedActions: Actions{
				InsertAction{
					Index: 0,
					Text:  "hello",
				},
			},
		},
		{
			name: "delete he",
			fn: func(r *Rogue) {
				_, err := r.Insert(0, "hello")
				require.NoError(t, err)
			},
			Op: DeleteOp{
				ID:         ID{"larry", 7},
				TargetID:   ID{"auth0", 3},
				SpanLength: 2,
			},
			expectedHtml: "<p><span>llo</span></p>",
			expectedActions: Actions{
				DeleteAction{
					Index: 0,
					Count: 2,
				},
			},
		},
		{
			name: "delete in the middle of an insert",
			fn: func(r *Rogue) {
				_, err := r.Insert(0, "hello world")
				require.NoError(t, err)
			},
			Op: DeleteOp{
				ID:         ID{"larry", 7},
				TargetID:   ID{"auth0", 5},
				SpanLength: 7,
			},
			expectedHtml: "<p><span>held</span></p>",
			expectedActions: Actions{
				DeleteAction{
					Index: 2,
					Count: 7,
				},
			},
		},
		{
			name: "delete some deleted text",
			fn: func(r *Rogue) {
				_, err := r.Insert(0, "hello big world")
				require.NoError(t, err)
				_, err = r.Delete(5, 10)
				require.NoError(t, err)
			},
			Op: DeleteOp{
				ID:         ID{"larry", 7},
				TargetID:   ID{"auth0", 8},
				SpanLength: 4,
			},
			expectedHtml:    "<p><span>hello</span></p>",
			expectedActions: Actions{},
		},
		{
			name: "delete starting on some deleted text",
			fn: func(r *Rogue) {
				_, err := r.Insert(0, "hello big world")
				require.NoError(t, err)
				_, err = r.Delete(5, 4)
				require.NoError(t, err)
			},
			Op: DeleteOp{
				ID:         ID{"larry", 7},
				TargetID:   ID{"auth0", 8},
				SpanLength: 10,
			},
			expectedHtml: "<p><span>hello</span></p>",
			expectedActions: Actions{
				DeleteAction{
					Index: 5,
					Count: 6,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := NewRogueForQuill("auth0")
			tc.fn(r)

			html, err := r.GetHtml(RootID, LastID, false, false)
			require.NoError(t, err, "Failed to get html before merge")

			fmt.Println("html before merge", html)

			actions, err := r.MergeOp(tc.Op)
			if tc.expectedErr == nil {
				require.NoError(t, err)
			}

			for i, root := range r.Roots {
				fmt.Printf("root[%d]:\n", i)
				root.Inspect()
				fmt.Println()
			}

			if tc.expectedHtml != "" {
				html, err := r.GetHtml(RootID, LastID, false, false)
				require.NoError(t, err)
				fmt.Println("html after merge", html)
				require.Equal(t, tc.expectedHtml, html)
			}

			if tc.expectedActions != nil {
				require.Equal(t, tc.expectedActions, actions)
			}
		})
	}
}

func TestMergeOps(t *testing.T) {
	testCases := []struct {
		name            string
		fn              func(*Rogue)
		Ops             []Op
		expectedHtml    string
		expectedActions Actions
	}{
		{
			name: "insert hello",
			fn:   func(r *Rogue) {},
			Ops: []Op{
				InsertOp{ID: ID{Author: "root", Seq: 0}, Text: "x", ParentID: ID{Author: "root", Seq: 0}, Side: 0},
				DeleteOp{ID: ID{Author: "q", Seq: 2}, TargetID: ID{Author: "root", Seq: 0}, SpanLength: 1},
				InsertOp{ID: ID{Author: "q", Seq: 1}, Text: "\n", ParentID: ID{Author: "root", Seq: 0}, Side: 1},
				InsertOp{ID: ID{Author: "auth0", Seq: 3}, Text: "🥔", ParentID: ID{Author: "q", Seq: 1}, Side: -1},
				InsertOp{ID: ID{Author: "auth0", Seq: 5}, Text: "🥒🥔", ParentID: ID{Author: "auth0", Seq: 4}, Side: 1},
				InsertOp{ID: ID{Author: "auth0", Seq: 9}, Text: "🤡😿😝", ParentID: ID{Author: "auth0", Seq: 5}, Side: -1},
			},
			expectedHtml: "<p><span>🥔🤡😿😝🥒🥔</span></p>",
		},
		{
			name: "weird unicode issue",
			fn:   func(r *Rogue) {},
			Ops: []Op{
				InsertOp{ID: ID{Author: "q", Seq: 1}, Text: "\n", ParentID: ID{Author: "root", Seq: 0}, Side: 1},
				DeleteOp{ID: ID{Author: "q", Seq: 2}, TargetID: ID{Author: "root", Seq: 0}, SpanLength: 1},
				InsertOp{ID: ID{Author: "auth0", Seq: 3}, Text: "🚟", ParentID: ID{Author: "q", Seq: 1}, Side: -1},
				InsertOp{ID: ID{Author: "auth1", Seq: 7}, Text: "🎅è🕩", ParentID: ID{Author: "auth0", Seq: 3}, Side: -1},
				InsertOp{ID: ID{Author: "auth0", Seq: 17}, Text: "Ë¿þ🗴", ParentID: ID{Author: "auth1", Seq: 11}, Side: 1},
				InsertOp{ID: ID{Author: "auth1", Seq: 27}, Text: "🦳ç😳😂", ParentID: ID{Author: "auth1", Seq: 7}, Side: -1},
				InsertOp{ID: ID{Author: "auth0", Seq: 41}, Text: "l", ParentID: ID{Author: "auth1", Seq: 10}, Side: -1},
				InsertOp{ID: ID{Author: "auth1", Seq: 42}, Text: "😗Ѥ", ParentID: ID{Author: "auth1", Seq: 32}, Side: -1},
				InsertOp{ID: ID{Author: "auth0", Seq: 48}, Text: "🗘T", ParentID: ID{Author: "auth0", Seq: 4}, Side: 1},
				InsertOp{ID: ID{Author: "auth1", Seq: 53}, Text: "-🚔öЫ", ParentID: ID{Author: "auth0", Seq: 48}, Side: -1},
			},
			expectedHtml: "<p><span>🦳ç😳😗Ѥ😂🎅èl🕩Ë¿þ🗴🚟-🚔öЫ🗘T</span></p>",
			expectedActions: Actions{
				InsertAction{Index: 0, Text: "🚟"},
				InsertAction{Index: 0, Text: "🎅è🕩"},
				InsertAction{Index: 5, Text: "Ë¿þ🗴"},
				InsertAction{Index: 0, Text: "🦳ç😳😂"},
				InsertAction{Index: 10, Text: "l"},
				InsertAction{Index: 5, Text: "😗Ѥ"},
				InsertAction{Index: 23, Text: "🗘T"},
				InsertAction{Index: 23, Text: "-🚔öЫ"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := NewRogueForQuill("auth0")
			tc.fn(r)

			totActions := Actions{}
			for _, op := range tc.Ops {
				actions, err := r.MergeOp(op)
				require.NoError(t, err)
				totActions = append(totActions, actions...)
			}

			for i, root := range r.Roots {
				fmt.Printf("root[%d]:\n", i)
				root.Inspect()
				fmt.Println()
			}

			err := r.ValidateFugues()
			require.NoError(t, err)

			if tc.expectedHtml != "" {
				html, err := r.GetHtml(RootID, LastID, false, false)
				require.NoError(t, err)
				require.Equal(t, tc.expectedHtml, html)
			}

			if tc.expectedActions != nil {
				require.Equal(t, tc.expectedActions, totActions)
			}
		})
	}
}

func TestInsertDeleteExpectations(t *testing.T) {
	var ops []Op
	tree0 := NewRogueForQuill("auth0")
	tree1 := NewRogueForQuill("auth1")

	newOp, err := tree0.Insert(0, "abc")
	if err != nil {
		t.Fatalf("Failed to insert 'abc': %v", err)
	}
	ops = append(ops, newOp)
	if tree0.GetText() != "abc\n" {
		t.Errorf("Expected 'abc\\n', got '%s'", tree0.GetText())
	}

	newOp, err = tree0.Insert(3, " hello world")
	if err != nil {
		t.Fatalf("Failed to insert ' hello world': %v", err)
	}
	ops = append(ops, newOp)
	if tree0.GetText() != "abc hello world\n" {
		t.Errorf("Expected 'abc hello world\\n', got '%s'", tree0.GetText())
	}

	// ... continue with the rest of the operations and expectations

	for _, op := range ops {
		tree1.MergeOp(op)
	}

	if tree0.GetText() != tree1.GetText() {
		t.Errorf("Expected Tree0 and Tree1 to have the same text, got Tree0: '%s', Tree1: '%s'", tree0.GetText(), tree1.GetText())
	}
}

func compareFugueRopeWalkRight(t *testing.T, doc *Rogue) {
	t.Helper()

	_, rn, err := doc.Rope.GetTotNode(0)
	require.NoError(t, err)

	ropeNodes := []*RopeNode{}
	fugueNodes := []*FugueNode{}

	err = doc.Rope.WalkRight(rn.Val.ID, func(rn *RopeNode) error {
		ropeNodes = append(ropeNodes, rn)
		return nil
	})
	require.NoError(t, err)

	err = rn.Val.WalkRight(func(fn *FugueNode) error {
		fugueNodes = append(fugueNodes, fn)
		return nil
	})
	require.NoError(t, err)

	require.Equal(t, len(ropeNodes), len(fugueNodes))

	for i := 0; i < len(ropeNodes); i++ {
		require.Equal(t, ropeNodes[i].Val.ID, fugueNodes[i].ID)
		require.Equal(t, ropeNodes[i].Val.Text, fugueNodes[i].Text)
	}
}

func TestMultiRogueMultiplayer(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}

	r := rand.New(rand.NewSource(42))

	for i := 0; i < 5; i++ {
		var trees []*Rogue
		for j := 0; j < 5; j++ {
			trees = append(trees, NewRogueForQuill(fmt.Sprintf("auth%d", j)))
		}

		for k := 0; k < 2000; k++ {
			var allOps []Op

			for _, tree := range trees {
				var ops []Op

				// do a few ops in a row on each tree so things
				// have a chance to get jumbled some more
				for j := 0; j < 5; j++ {
					if r.Float64() > 0.75 && tree.VisSize > 0 {
						_, dop, err := tree.RandDelete(r, 10)
						require.NoError(t, err)

						ops = append(ops, dop)
					} else {
						_, op, err := tree.RandInsert(r, 10)
						require.NoError(t, err)

						ops = append(ops, op)
					}
				}

				// randomly interleave the ops to make this more realistic
				allOps = Interleave(r, allOps, ops)
			}

			for _, tree := range trees {
				for _, op := range allOps {
					tree.MergeOp(op)
				}
			}
		}

		// Assertions
		for _, tree := range trees {
			err := tree.Rope.Validate()
			require.NoError(t, err)

			err = tree.ValidateFugues()
			require.NoError(t, err)

			require.Equal(t, trees[0].VisSize, tree.VisSize)
			require.Equal(t, trees[0].TotSize, tree.TotSize)
		}

		_, totIx, err := trees[0].Rope.GetIndex(RootID)
		require.NoError(t, err)
		fmt.Printf("RootID totIx: %d\n", totIx)

		serRogue, err := json.Marshal(trees[0])
		require.NoError(t, err)

		var deRogue Rogue
		err = json.Unmarshal(serRogue, &deRogue)
		// fmt.Printf("serRogue: %q\n", string(serRogue)[:100])
		require.NoError(t, err)

		dTreeText := deRogue.GetText()

		for _, tree := range trees {
			tt := tree.GetText()
			if tt != dTreeText {
				t.Errorf("Expected all tree texts to be the same")
			}
		}

		compareFugueRopeWalkRight(t, trees[0])
	}
}

func TestMultiRogueMultiplayerActions(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}

	type RogueAction struct {
		r *Rogue
		s []uint16
	}

	r := rand.New(rand.NewSource(42))

	for i := 0; i < 5; i++ {
		var trees []*RogueAction
		for j := 0; j < 5; j++ {
			trees = append(trees, &RogueAction{
				r: NewRogueForQuill(fmt.Sprintf("auth%d", j)),
				s: StrToUint16("\n"),
			})
		}

		for k := 0; k < 1000; k++ {
			var allOps []Op

			for _, tree := range trees {
				var ops []Op

				// do a few ops in a row on each tree so things
				// have a chance to get jumbled some more
				for j := 0; j < 5; j++ {
					if r.Float64() > 0.75 && tree.r.VisSize > 0 {
						action, dop, err := tree.r.RandDelete(r, 10)
						require.NoError(t, err)

						ops = append(ops, dop)
						tree.s = DeleteAt(tree.s, action.Index, action.Count)
					} else {
						action, op, err := tree.r.RandInsert(r, 10)
						require.NoError(t, err)

						ops = append(ops, op)
						tree.s = InsertSliceAt(tree.s, action.Index, StrToUint16(action.Text))
					}
				}

				// randomly interleave the ops to make this more realistic
				allOps = Interleave(r, allOps, ops)
			}

			for _, tree := range trees {
				for _, op := range allOps {
					actions, err := tree.r.MergeOp(op)
					require.NoError(t, err)

					for _, action := range actions {
						switch action := action.(type) {
						case InsertAction:
							tree.s = InsertSliceAt(tree.s, action.Index, StrToUint16(action.Text))
						case DeleteAction:
							tree.s = DeleteAt(tree.s, action.Index, action.Count)
						}
					}
				}
			}
		}

		// Assertions
		t0 := trees[0]
		t0r := t0.r.GetText()
		t0s := Uint16ToStr(t0.s)
		require.Equal(t, t0s, t0r)
		for _, tree := range trees[1:] {
			require.Equal(t, t0r, tree.r.GetText())
			require.Equal(t, t0s, Uint16ToStr(tree.s))
		}
	}
}

func TestMultiDeleteSplitNode(t *testing.T) {
	tree0 := NewRogueForQuill("auth0")
	tree1 := NewRogueForQuill("auth1")

	op0, err := tree0.Insert(0, "Hello World!")
	require.NoError(t, err)

	fmt.Printf("tree0.GetText(): %v\n", tree0.GetText())

	op1, err := tree0.Insert(5, " cruel")
	require.NoError(t, err)

	_, err = tree1.MergeOp(op0)
	require.NoError(t, err)
	require.Equal(t, "Hello World!\n", tree1.GetText())

	op2, err := tree1.Delete(0, 12)
	require.NoError(t, err)

	_, err = tree1.MergeOp(op1)
	require.NoError(t, err)

	fmt.Printf("op2: %v\n", op2)
	_, err = tree1.MergeOp(op2)
	require.NoError(t, err)
}

func TestMarshalUnmarshalComplexString(t *testing.T) {
	r := rand.New(rand.NewSource(42))
	originalText := RandomString(r, 100000)
	// randomly delete a bunch of characters in the originalText to create a new string

	var builder strings.Builder
	for i := 0; i < len(originalText); i++ {
		// randomly delete maybe creating broken unicode chars
		if r.Float64() > 0.5 {
			builder.WriteString(string(originalText[i]))
		}
	}

	newString := builder.String()
	fmt.Printf("len(newString): %d\n", len(newString))

	marshaled, err := json.Marshal(newString)
	assert.NoError(t, err)

	var unmarshaledText string
	err = json.Unmarshal(marshaled, &unmarshaledText)
	assert.NoError(t, err)

	assert.Equal(t, newString, unmarshaledText)
}

func TestRogueInsertDelete(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}

	r := rand.New(rand.NewSource(1))

	// Initialize the FugueNode and a slice for comparison
	tree := NewRogueForQuill("auth0")
	slice := StrToUint16("\n")
	var err error

	for i := 0; i < 10000; i++ {
		if r.Float64() > 0.75 && tree.VisSize > 0 {
			action, _, err := tree.RandDelete(r, 10)
			require.NoError(t, err)
			slice = DeleteAt(slice, action.Index, action.Count)
		} else {
			action, _, err := tree.RandInsert(r, 10)
			require.NoError(t, err)
			slice = InsertSliceAt(slice, action.Index, StrToUint16(action.Text))
		}

		// Compare the results
		if Uint16ToStr(slice) != tree.GetText() {
			fmt.Printf("SLICE LEN: %d\n", len(Uint16ToStr(slice)))
			fmt.Printf("ROGUE LEN: %d\n", len(tree.GetText()))

			x, y := utf16.Decode(slice), []rune(tree.GetText())
			fmt.Printf("SLICE: %v\n", x)
			fmt.Printf("ROGUE: %v\n", y)

			fmt.Printf("SLICE: %v\n", Uint16ToStr(slice))
			fmt.Printf("ROGUE: %v\n", tree.GetText())

			t.Fatalf("Mismatch after operation %d", i)
		}
	}

	// Validate the node's structure if applicable
	err = tree.Rope.Validate()
	require.NoError(t, err)

	err = tree.ValidateFugues()
	require.NoError(t, err)

	serTree, err := json.Marshal(tree)
	assert.NoError(t, err)

	var tree2 Rogue
	err = json.Unmarshal(serTree, &tree2)
	assert.NoError(t, err)

	assert.Equal(t, Uint16ToStr(slice), tree2.GetText())

	// Test visLeftOf and visRightOf
	ix := tree.VisSize
	id, err := tree.Rope.GetVisID(ix - 1)
	require.NoError(t, err)
	for {
		newIx, _, err := tree.Rope.GetIndex(id)
		require.NoError(t, err)

		assert.Equal(t, ix-1, newIx)

		ix = newIx

		c, err := tree.GetChar(ix)
		require.NoError(t, err)
		assert.Equal(t, slice[ix], c)

		c2, err := tree2.GetChar(ix)
		require.NoError(t, err)
		assert.Equal(t, slice[ix], c2)

		id, err = tree.VisLeftOf(id)
		if err != nil {
			if errors.As(err, &ErrorNoLeftVisSibling{}) {
				break
			}
			require.NoError(t, err)
		}
	}

	ix = -1
	id, err = tree.Rope.GetVisID(0)
	require.NoError(t, err)
	for {
		newIx, _, err := tree.Rope.GetIndex(id)
		require.NoError(t, err)

		assert.Equal(t, ix+1, newIx)

		ix = newIx

		c, err := tree.GetChar(ix)
		require.NoError(t, err)
		assert.Equal(t, slice[ix], c)

		c2, err := tree2.GetChar(ix)
		require.NoError(t, err)
		assert.Equal(t, slice[ix], c2)

		id, err = tree.VisRightOf(id)
		if err != nil {
			if errors.As(err, &ErrorNoRightVisSibling{}) {
				break
			}
			require.NoError(t, err)
		}
	}
}

func TestMultiRogueSerialize(t *testing.T) {
	rogue0 := NewRogueForQuill("auth0")
	var ops []Op

	var op Op
	op, err := rogue0.Insert(0, "Hello World!")
	require.NoError(t, err)
	ops = append(ops, op)
	fop, err := rogue0.Format(0, 2, FormatV3Span{"bold": "true"})
	require.NoError(t, err)
	ops = append(ops, fop)
	dop, err := rogue0.Delete(0, 1)
	require.NoError(t, err)
	ops = append(ops, dop)

	sRogue0, err := json.Marshal(rogue0)
	require.NoError(t, err)

	rogue1 := NewRogueForQuill("auth0")
	for _, op := range ops {
		action, err := rogue1.MergeOp(op)
		require.NoError(t, err)
		actionJSON, err := json.Marshal(action)
		require.NoError(t, err)
		t.Log(string(actionJSON))
	}

	sRogue1, err := json.Marshal(rogue1)
	require.NoError(t, err)

	require.Equal(t, string(sRogue0), string(sRogue1))

	var dRogue0 Rogue
	err = json.Unmarshal(sRogue0, &dRogue0)
	require.NoError(t, err)

	t.Log(dRogue0.GetText())
	require.Equal(t, rogue0.GetText(), dRogue0.GetText())
	require.Equal(t, rogue0.VisSize, dRogue0.VisSize)
}

func TestRogueWithRealWorldDoc(t *testing.T) {
	r := NewRogueForQuill("auth0")

	ops := []Op{
		InsertOp{ID: ID{"891d-38", 3}, Text: "Kyla", ParentID: ID{"q", 1}, Side: Left},
		FormatOp{ID: ID{"891d-38", 7}, StartID: ID{"891d-38", 3}, EndID: ID{"891d-38", 6}, Format: FormatV3Span{"bold": "true", "link": "http://google.com"}},
		InsertOp{ID: ID{"891d-38", 8}, Text: "\u00a0(born January\u00a05, 1981) is a Filipino singer-songwriter. She is the recipient of\u00a0", ParentID: ID{"891d-38", 6}, Side: Right},
		InsertOp{ID: ID{"891d-38", 89}, Text: " numerous accolades", ParentID: ID{"891d-38", 88}, Side: Left},
		FormatOp{ID: ID{"891d-38", 108}, StartID: ID{"891d-38", 90}, EndID: ID{"891d-38", 107}, Format: FormatV3Span{"bold": "true", "link": "http://google.com"}},
		InsertOp{ID: ID{"891d-38", 109}, Text: ",", ParentID: ID{"891d-38", 107}, Side: Right},
	}

	for _, op := range ops {
		_, err := r.MergeOp(op)
		if err != nil {
			t.Fatalf("Failed to merge op: %v", err)
		}
	}

	assert.Equal(t, "Kyla\u00a0(born January\u00a05, 1981) is a Filipino singer-songwriter. She is the recipient of numerous accolades,\u00a0\n", r.GetText())
}

func TestRogueWithUnicode(t *testing.T) {
	r := NewRogueForQuill("auth0")

	ops := []Op{
		InsertOp{ID: ID{"891d-38", 3}, Text: "Hello", ParentID: ID{"q", 1}, Side: Left},
		InsertOp{ID: ID{"891d-38", 8}, Text: "\u00a0World\u00a0", ParentID: ID{"891d-38", 6}, Side: Right},
		InsertOp{ID: ID{"891d-38", 15}, Text: "you're great", ParentID: ID{"891d-38", 14}, Side: Right},
	}

	for _, op := range ops {
		_, err := r.MergeOp(op)
		if err != nil {
			t.Fatalf("Failed to merge op: %v", err)
		}
	}

	assert.Equal(t, "Hello\u00a0World\u00a0you're great\n", r.GetText())
}

func TestGetIdBug(t *testing.T) {
	rogue := NewRogueForQuill("auth0")
	_, err := rogue.Insert(0, "Hello World!")
	assert.NoError(t, err)

	for i := 0; i < rogue.VisSize; i++ {
		_, err := rogue.Rope.GetVisID(i)
		assert.NoError(t, err)
	}

	_, err = rogue.Rope.GetVisID(-1)
	assert.Error(t, err)

	_, err = rogue.Rope.GetVisID(13)
	assert.Error(t, err)
}

func TestRogueSerialization(t *testing.T) {
	rogue := NewRogueForQuill("auth0")
	rogue.Insert(0, "Hello World!")
	rogue.Format(0, 2, FormatV3Span{"bold": "true"})
	rogue.Delete(0, 1)

	serRogue, err := json.Marshal(rogue)
	if err != nil {
		t.Fatalf("Failed to serialize rogue: %v", err)
	}
	fmt.Println(string(serRogue))

	var deRogue Rogue
	err = json.Unmarshal(serRogue, &deRogue)
	if err != nil {
		t.Fatalf("Failed to deserialize rogue: %v", err)
	}

	t.Log(deRogue.GetText())
}

func TestRealDoc(t *testing.T) {
	// convert the ops below into rogue ops
	/*
	 */
	ops := []Op{
		InsertOp{ID{"8cd8-au", 3}, "H", ID{"q", 1}, Left},
		InsertOp{ID{"8cd8-au", 4}, "e", ID{"8cd8-au", 3}, Right},
		InsertOp{ID{"8cd8-au", 5}, "l", ID{"8cd8-au", 4}, Right},
		InsertOp{ID{"8cd8-au", 6}, "l", ID{"8cd8-au", 5}, Right},
		InsertOp{ID{"8cd8-au", 7}, "o", ID{"8cd8-au", 6}, Right},
		InsertOp{ID{"8cd8-au", 8}, "W", ID{"8cd8-au", 7}, Right},
		InsertOp{ID{"8cd8-au", 9}, "o", ID{"8cd8-au", 8}, Right},
		InsertOp{ID{"8cd8-au", 10}, "r", ID{"8cd8-au", 9}, Right},
		InsertOp{ID{"8cd8-au", 11}, "l", ID{"8cd8-au", 10}, Right},
		InsertOp{ID{"8cd8-au", 12}, "d", ID{"8cd8-au", 11}, Right},
		InsertOp{ID{"8cd8-au", 13}, "!", ID{"8cd8-au", 12}, Right},
		InsertOp{ID{"8cd8-ue", 14}, " ", ID{"8cd8-au", 13}, Right}, // unicode nbsp
		InsertOp{ID{"8cd8-ue", 15}, " h", ID{"8cd8-ue", 14}, Left},
		DeleteOp{ID{"8cd8-ue", 17}, ID{"8cd8-ue", 14}, 1},
		InsertOp{ID{"8cd8-ue", 18}, "o", ID{"8cd8-ue", 16}, Right},
	}

	rogue := NewRogueForQuill("auth0")
	for _, op := range ops {
		fmt.Printf("OP: %+v\n", op)

		actions, err := rogue.MergeOp(op)
		assert.NoError(t, err)

		fmt.Printf("ACTIONS: %+v\n", actions)
		fmt.Printf("ACTIONS: %T\n", actions)
	}

	fmt.Println(rogue.GetText())
}

func TestTotLeftRight(t *testing.T) {
	rogue := NewRogueForQuill("auth0")
	rogue.Insert(0, "Hello World!")

	leftID, err := rogue.Rope.GetVisID(0)
	assert.NoError(t, err)

	rightID, err := rogue.Rope.GetVisID(11)
	assert.NoError(t, err)

	startID, err := rogue.TotLeftOf(leftID)
	assert.NoError(t, err)

	endID, err := rogue.TotRightOf(rightID)
	assert.NoError(t, err)

	assert.Equal(t, "root_0", startID.String())
	assert.Equal(t, "q_1", endID.String())
}

func TestFormat(t *testing.T) {
	rogue := NewRogueForQuill("auth0")
	sentence := "Hello World!"
	var op Op
	var err error
	var ops []Op

	op, err = rogue.Insert(0, sentence)
	assert.NoError(t, err)
	ops = append(ops, op)

	op, err = rogue.Insert(0, sentence)
	assert.NoError(t, err)
	ops = append(ops, op)

	fop, err := rogue.Format(0, 5, FormatV3Span{"u": "true"})
	assert.NoError(t, err)
	ops = append(ops, fop)

	fop, err = rogue.Format(2, 1, FormatV3Span{"u": "null"})
	assert.NoError(t, err)
	ops = append(ops, fop)

	for i := 0; i < len(sentence); i++ {
		dop, err := rogue.Delete(0, 1)
		assert.NoError(t, err)
		ops = append(ops, dop)
	}

	fmt.Printf("TEXT: %s\n", rogue.GetText())

	fop, err = rogue.Format(2, 8, FormatV3Span{"b": "true"})
	assert.NoError(t, err)
	ops = append(ops, fop)

	fop, err = rogue.Format(3, 3, FormatV3Span{"b": "null"})
	assert.NoError(t, err)
	ops = append(ops, fop)

	rogue2 := NewRogueForQuill("auth0")

	var action []Action
	var actions []Action
	for _, op := range ops {
		action, err = rogue2.MergeOp(op)
		PrintJson("ACTION", action)
		require.NoError(t, err)
		actions = append(actions, action...)
	}
	PrintJson("ACTIONS", actions)
	PrintJson("ROGUE", rogue2)
}

func TestRopeVisSibling(t *testing.T) {
	rogue := NewRogueForQuill("auth0")
	sentence := "Hello World!"
	var err error

	_, err = rogue.Insert(0, sentence)
	assert.NoError(t, err)

	_, err = rogue.Insert(0, sentence)
	assert.NoError(t, err)

	_, err = rogue.Insert(len(sentence), sentence)
	assert.NoError(t, err)

	// Delete all the characters in left node
	_, err = rogue.Delete(0, len(sentence))
	assert.NoError(t, err)
	fmt.Printf("text0: %q\n", rogue.GetText())

	// Delete all the characters in the right node
	// len(sentence)+1 to delete implicit newline
	_, err = rogue.Delete(len(sentence), len(sentence)+1)
	assert.NoError(t, err)
	fmt.Printf("text1: %q\n", rogue.GetText())

	rogueText := rogue.GetText()
	// fmt.Println(rogueText)
	require.Equal(t, sentence, rogueText)

	// Check left side
	visOffset, ropeNode, err := rogue.Rope.GetNode(0)
	assert.NoError(t, err)
	assert.Equal(t, 0, visOffset)

	leftTot, err := ropeNode.LeftTotSibling()
	assert.NoError(t, err)
	fmt.Printf("leftTot: %+v\n", leftTot.Val.ID)
	fmt.Printf("leftText: %s\n", Uint16ToStr(leftTot.Val.Text))

	leftVis, err := ropeNode.LeftVisSibling()
	assert.Error(t, err)
	fmt.Printf("%+v\n", err)
	fmt.Printf("leftVis: %+v\n", leftVis)

	// Check right side
	visOffset, ropeNode, err = rogue.Rope.GetNode(len(sentence) - 1)
	assert.NoError(t, err)
	assert.Equal(t, len(sentence)-1, visOffset)

	rightTot, err := ropeNode.RightTotSibling()
	assert.NoError(t, err)
	fmt.Printf("rightTot: %+v\n", rightTot.Val.ID)
	fmt.Printf("rightText: %s\n", Uint16ToStr(rightTot.Val.Text))
	fmt.Printf("rightVisText: %s\n", string(rightTot.Val.VisibleText()))

	rightVis, err := ropeNode.RightVisSibling()
	assert.Error(t, err)
	fmt.Printf("%+v\n", err)
	fmt.Printf("rightVis: %+v\n", rightVis)
}

func TestEmoji(t *testing.T) {
	rogue := NewRogueForQuill("auth0")
	sentence := "Hello World! 🌎🙂"
	var err error

	_, err = rogue.Insert(0, sentence)
	assert.NoError(t, err)
}

func TestAllUnicodeMarshal(t *testing.T) {
	t.Skip("Skipping test")
	ranges := [][]int{
		// U+0000 - U+007F	Basic Latin	128
		{0x0000, 0x007F},
		// U+0080 - U+00FF	Latin-1 Supplement	128
		{0x0080, 0x00FF},
		// U+0100 - U+017F	Latin Extended-A	128
		{0x0100, 0x017F},
		// U+0180 - U+024F	Latin Extended-B	208
		{0x0180, 0x024F},
		// U+0250 - U+02AF	IPA Extensions	96
		{0x0250, 0x02AF},
		// U+02B0 - U+02FF	Spacing Modifier Letters	80
		{0x02B0, 0x02FF},
		// U+0300 - U+036F	Combining Diacritical Marks	112
		{0x0300, 0x036F},
		// U+0370 - U+03FF	Greek and Coptic	135
		{0x0370, 0x03FF},
		// U+0400 - U+04FF	Cyrillic	256
		{0x0400, 0x04FF},
		// U+0500 - U+052F	Cyrillic Supplement	48
		{0x0500, 0x052F},
		// U+0530 - U+058F	Armenian	91
		{0x0530, 0x058F},
		// U+0590 - U+05FF	Hebrew	88
		{0x0590, 0x05FF},
		// U+0600 - U+06FF	Arabic	255
		{0x0600, 0x06FF},
		// U+0700 - U+074F	Syriac	77
		{0x0700, 0x074F},
		// U+0750 - U+077F	Arabic Supplement	48
		{0x0750, 0x077F},
		// U+0780 - U+07BF	Thaana	50
		{0x0780, 0x07BF},
		// U+07C0 - U+07FF	NKo	62
		{0x07C0, 0x07FF},
		// U+0800 - U+083F	Samaritan	61
		{0x0800, 0x083F},
		// U+0840 - U+085F	Mandaic	29
		{0x0840, 0x085F},
		// U+0860 - U+086F	Syriac Supplement	11
		{0x0860, 0x086F},
		// U+08A0 - U+08FF	Arabic Extended-A	84
		{0x08A0, 0x08FF},
		// U+0900 - U+097F	Devanagari	128
		{0x0900, 0x097F},
		// U+0980 - U+09FF	Bengali	96
		{0x0980, 0x09FF},
		// U+0A00 - U+0A7F	Gurmukhi	80
		{0x0A00, 0x0A7F},
		// U+0A80 - U+0AFF	Gujarati	91
		{0x0A80, 0x0AFF},
		// U+0B00 - U+0B7F	Oriya	91
		{0x0B00, 0x0B7F},
		// U+0B80 - U+0BFF	Tamil	72
		{0x0B80, 0x0BFF},
		// U+0C00 - U+0C7F	Telugu	98
		{0x0C00, 0x0C7F},
		// U+0C80 - U+0CFF	Kannada	89
		{0x0C80, 0x0CFF},
		// U+0D00 - U+0D7F	Malayalam	118
		{0x0D00, 0x0D7F},
		// U+0D80 - U+0DFF	Sinhala	91
		{0x0D80, 0x0DFF},
		// U+0E00 - U+0E7F	Thai	87
		{0x0E00, 0x0E7F},
		// U+0E80 - U+0EFF	Lao	82
		{0x0E80, 0x0EFF},
		// U+0F00 - U+0FFF	Tibetan	211
		{0x0F00, 0x0FFF},
		// U+1000 - U+109F	Myanmar	160
		{0x1000, 0x109F},
		// U+10A0 - U+10FF	Georgian	88
		{0x10A0, 0x10FF},
		// U+1100 - U+11FF	Hangul Jamo	256
		{0x1100, 0x11FF},
		// U+1200 - U+137F	Ethiopic	358
		{0x1200, 0x137F},
		// U+1380 - U+139F	Ethiopic Supplement	26
		{0x1380, 0x139F},
		// U+13A0 - U+13FF	Cherokee	92
		{0x13A0, 0x13FF},
		// U+1400 - U+167F	Unified Canadian Aboriginal Syllabics	640
		{0x1400, 0x167F},
		// U+1680 - U+169F	Ogham	29
		{0x1680, 0x169F},
		// U+16A0 - U+16FF	Runic	89
		{0x16A0, 0x16FF},
		// U+1700 - U+171F	Tagalog	20
		{0x1700, 0x171F},
		// U+1720 - U+173F	Hanunoo	23
		{0x1720, 0x173F},
		// U+1740 - U+175F	Buhid	20
		{0x1740, 0x175F},
		// U+1760 - U+177F	Tagbanwa	18
		{0x1760, 0x177F},
		// U+1780 - U+17FF	Khmer	114
		{0x1780, 0x17FF},
		// U+1800 - U+18AF	Mongolian	157
		{0x1800, 0x18AF},
		// U+18B0 - U+18FF	Unified Canadian Aboriginal Syllabics Extended	70
		{0x18B0, 0x18FF},
		// U+1900 - U+194F	Limbu	68
		{0x1900, 0x194F},
		// U+1950 - U+197F	Tai Le	35
		{0x1950, 0x197F},
		// U+1980 - U+19DF	New Tai Lue	83
		{0x1980, 0x19DF},
		// U+19E0 - U+19FF	Khmer Symbols	32
		{0x19E0, 0x19FF},
		// U+1A00 - U+1A1F	Buginese	30
		{0x1A00, 0x1A1F},
		// U+1A20 - U+1AAF	Tai Tham	127
		{0x1A20, 0x1AAF},
		// U+1AB0 - U+1AFF	Combining Diacritical Marks Extended	17
		{0x1AB0, 0x1AFF},
		// U+1B00 - U+1B7F	Balinese	121
		{0x1B00, 0x1B7F},
		// U+1B80 - U+1BBF	Sundanese	64
		{0x1B80, 0x1BBF},
		// U+1BC0 - U+1BFF	Batak	56
		{0x1BC0, 0x1BFF},
		// U+1C00 - U+1C4F	Lepcha	74
		{0x1C00, 0x1C4F},
		// U+1C50 - U+1C7F	Ol Chiki	48
		{0x1C50, 0x1C7F},
		// U+1C80 - U+1C8F	Cyrillic Extended-C	9
		{0x1C80, 0x1C8F},
		// U+1C90 - U+1CBF	Georgian Extended	46
		{0x1C90, 0x1CBF},
		// U+1CC0 - U+1CCF	Sundanese Supplement	8
		{0x1CC0, 0x1CCF},
		// U+1CD0 - U+1CFF	Vedic Extensions	43
		{0x1CD0, 0x1CFF},
		// U+1D00 - U+1D7F	Phonetic Extensions	128
		{0x1D00, 0x1D7F},
		// U+1D80 - U+1DBF	Phonetic Extensions Supplement	64
		{0x1D80, 0x1DBF},
		// U+1DC0 - U+1DFF	Combining Diacritical Marks Supplement	63
		{0x1DC0, 0x1DFF},
		// U+1E00 - U+1EFF	Latin Extended Additional	256
		{0x1E00, 0x1EFF},
		// U+1F00 - U+1FFF	Greek Extended	233
		{0x1F00, 0x1FFF},
		// U+2000 - U+206F	General Punctuation	111
		{0x2000, 0x206F},
		// U+2070 - U+209F	Superscripts and Subscripts	42
		{0x2070, 0x209F},
		// U+20A0 - U+20CF	Currency Symbols	32
		{0x20A0, 0x20CF},
		// U+20D0 - U+20FF	Combining Diacritical Marks for Symbols	33
		{0x20D0, 0x20FF},
		// U+2100 - U+214F	Letterlike Symbols	80
		{0x2100, 0x214F},
		// U+2150 - U+218F	Number Forms	60
		{0x2150, 0x218F},
		// U+2190 - U+21FF	Arrows	112
		{0x2190, 0x21FF},
		// U+2200 - U+22FF	Mathematical Operators	256
		{0x2200, 0x22FF},
		// U+2300 - U+23FF	Miscellaneous Technical	256
		{0x2300, 0x23FF},
		// U+2400 - U+243F	Control Pictures	39
		{0x2400, 0x243F},
		// U+2440 - U+245F	Optical Character Recognition	11
		{0x2440, 0x245F},
		// U+2460 - U+24FF	Enclosed Alphanumerics	160
		{0x2460, 0x24FF},
		// U+2500 - U+257F	Box Drawing	128
		{0x2500, 0x257F},
		// U+2580 - U+259F	Block Elements	32
		{0x2580, 0x259F},
		// U+25A0 - U+25FF	Geometric Shapes	96
		{0x25A0, 0x25FF},
		// U+2600 - U+26FF	Miscellaneous Symbols	256
		{0x2600, 0x26FF},
		// U+2700 - U+27BF	Dingbats	192
		{0x2700, 0x27BF},
		// U+27C0 - U+27EF	Miscellaneous Mathematical Symbols-A	48
		{0x27C0, 0x27EF},
		// U+27F0 - U+27FF	Supplemental Arrows-A	16
		{0x27F0, 0x27FF},
		// U+2800 - U+28FF	Braille Patterns	256
		{0x2800, 0x28FF},
		// U+2900 - U+297F	Supplemental Arrows-B	128
		{0x2900, 0x297F},
		// U+2980 - U+29FF	Miscellaneous Mathematical Symbols-B	128
		{0x2980, 0x29FF},
		// U+2A00 - U+2AFF	Supplemental Mathematical Operators	256
		{0x2A00, 0x2AFF},
		// U+2B00 - U+2BFF	Miscellaneous Symbols and Arrows	253
		{0x2B00, 0x2BFF},
		// U+2C00 - U+2C5F	Glagolitic	94
		{0x2C00, 0x2C5F},
		// U+2C60 - U+2C7F	Latin Extended-C	32
		{0x2C60, 0x2C7F},
		// U+2C80 - U+2CFF	Coptic	123
		{0x2C80, 0x2CFF},
		// U+2D00 - U+2D2F	Georgian Supplement	40
		{0x2D00, 0x2D2F},
		// U+2D30 - U+2D7F	Tifinagh	59
		{0x2D30, 0x2D7F},
		// U+2D80 - U+2DDF	Ethiopic Extended	79
		{0x2D80, 0x2DDF},
		// U+2DE0 - U+2DFF	Cyrillic Extended-A	32
		{0x2DE0, 0x2DFF},
		// U+2E00 - U+2E7F	Supplemental Punctuation	83
		{0x2E00, 0x2E7F},
		// U+2E80 - U+2EFF	CJK Radicals Supplement	115
		{0x2E80, 0x2EFF},
		// U+2F00 - U+2FDF	Kangxi Radicals	214
		{0x2F00, 0x2FDF},
		// U+2FF0 - U+2FFF	Ideographic Description Characters	12
		{0x2FF0, 0x2FFF},
		// U+3000 - U+303F	CJK Symbols and Punctuation	64
		{0x3000, 0x303F},
		// U+3040 - U+309F	Hiragana	93
		{0x3040, 0x309F},
		// U+30A0 - U+30FF	Katakana	96
		{0x30A0, 0x30FF},
		// U+3100 - U+312F	Bopomofo	43
		{0x3100, 0x312F},
		// U+3130 - U+318F	Hangul Compatibility Jamo	94
		{0x3130, 0x318F},
		// U+3190 - U+319F	Kanbun	16
		{0x3190, 0x319F},
		// U+31A0 - U+31BF	Bopomofo Extended	32
		{0x31A0, 0x31BF},
		// U+31C0 - U+31EF	CJK Strokes	36
		{0x31C0, 0x31EF},
		// U+31F0 - U+31FF	Katakana Phonetic Extensions	16
		{0x31F0, 0x31FF},
		// U+3200 - U+32FF	Enclosed CJK Letters and Months	255
		{0x3200, 0x32FF},
		// U+3300 - U+33FF	CJK Compatibility	256
		{0x3300, 0x33FF},
		// U+3400 - U+4DBF	CJK Unified Ideographs Extension A	6,592
		{0x3400, 0x4DBF},
		// U+4DC0 - U+4DFF	Yijing Hexagram Symbols	64
		{0x4DC0, 0x4DFF},
		// U+4E00 - U+9FFF	CJK Unified Ideographs	20,989
		{0x4E00, 0x9FFF},
		// U+A000 - U+A48F	Yi Syllables	1,165
		{0xA000, 0xA48F},
		// U+A490 - U+A4CF	Yi Radicals	55
		{0xA490, 0xA4CF},
		// U+A4D0 - U+A4FF	Lisu	48
		{0xA4D0, 0xA4FF},
		// U+A500 - U+A63F	Vai	300
		{0xA500, 0xA63F},
		// U+A640 - U+A69F	Cyrillic Extended-B	96
		{0xA640, 0xA69F},
		// U+A6A0 - U+A6FF	Bamum	88
		{0xA6A0, 0xA6FF},
		// U+A700 - U+A71F	Modifier Tone Letters	32
		{0xA700, 0xA71F},
		// U+A720 - U+A7FF	Latin Extended-D	180
		{0xA720, 0xA7FF},
		// U+A800 - U+A82F	Syloti Nagri	45
		{0xA800, 0xA82F},
		// U+A830 - U+A83F	Common Indic Number Forms	10
		{0xA830, 0xA83F},
		// U+A840 - U+A87F	Phags-pa	56
		{0xA840, 0xA87F},
		// U+A880 - U+A8DF	Saurashtra	82
		{0xA880, 0xA8DF},
		// U+A8E0 - U+A8FF	Devanagari Extended	32
		{0xA8E0, 0xA8FF},
		// U+A900 - U+A92F	Kayah Li	48
		{0xA900, 0xA92F},
		// U+A930 - U+A95F	Rejang	37
		{0xA930, 0xA95F},
		// U+A960 - U+A97F	Hangul Jamo Extended-A	29
		{0xA960, 0xA97F},
		// U+A980 - U+A9DF	Javanese	91
		{0xA980, 0xA9DF},
		// U+A9E0 - U+A9FF	Myanmar Extended-B	31
		{0xA9E0, 0xA9FF},
		// U+AA00 - U+AA5F	Cham	83
		{0xAA00, 0xAA5F},
		// U+AA60 - U+AA7F	Myanmar Extended-A	32
		{0xAA60, 0xAA7F},
		// U+AA80 - U+AADF	Tai Viet	72
		{0xAA80, 0xAADF},
		// U+AAE0 - U+AAFF	Meetei Mayek Extensions	23
		{0xAAE0, 0xAAFF},
		// U+AB00 - U+AB2F	Ethiopic Extended-A	32
		{0xAB00, 0xAB2F},
		// U+AB30 - U+AB6F	Latin Extended-E	60
		{0xAB30, 0xAB6F},
		// U+AB70 - U+ABBF	Cherokee Supplement	80
		{0xAB70, 0xABBF},
		// U+ABC0 - U+ABFF	Meetei Mayek	56
		{0xABC0, 0xABFF},
		// U+AC00 - U+D7AF	Hangul Syllables	11,172
		{0xAC00, 0xD7AF},
		// U+D7B0 - U+D7FF	Hangul Jamo Extended-B	72
		{0xD7B0, 0xD7FF},
		// U+D800 - U+DB7F	High Surrogates	0
		{0xD800, 0xDB7F},
		// U+DB80 - U+DBFF	High Private Use Surrogates	0
		{0xDB80, 0xDBFF},
		// U+DC00 - U+DFFF	Low Surrogates	0
		{0xDC00, 0xDFFF},
		// U+E000 - U+F8FF	Private Use Area	0
		{0xE000, 0xF8FF},
		// U+F900 - U+FAFF	CJK Compatibility Ideographs	472
		{0xF900, 0xFAFF},
		// U+FB00 - U+FB4F	Alphabetic Presentation Forms	58
		{0xFB00, 0xFB4F},
		// U+FB50 - U+FDFF	Arabic Presentation Forms-A	611
		{0xFB50, 0xFDFF},
		// U+FE00 - U+FE0F	Variation Selectors	16
		{0xFE00, 0xFE0F},
		// U+FE10 - U+FE1F	Vertical Forms	10
		{0xFE10, 0xFE1F},
		// U+FE20 - U+FE2F	Combining Half Marks	16
		{0xFE20, 0xFE2F},
		// U+FE30 - U+FE4F	CJK Compatibility Forms	32
		{0xFE30, 0xFE4F},
		// U+FE50 - U+FE6F	Small Form Variants	26
		{0xFE50, 0xFE6F},
		// U+FE70 - U+FEFF	Arabic Presentation Forms-B	141
		{0xFE70, 0xFEFF},
		// U+FF00 - U+FFEF	Halfwidth and Fullwidth Forms	225
		{0xFF00, 0xFFEF},
		// U+FFF0 - U+FFFF	Specials	5
		{0xFFF0, 0xFFFF},
	}

	for _, r := range ranges {
		root := NewFugueNode(ID{"root", 0}, []uint16{0}, Root, nil)
		lastNode := root
		depth := 0
		for i := r[0]; i <= r[1]; i++ {
			node := NewFugueNode(ID{"auth0", i}, []uint16{uint16(i)}, Left, lastNode)
			lastNode.LeftChildren = append(lastNode.LeftChildren, node)
			lastNode = node
			depth++

			if depth%100 == 0 {
				t.Run(fmt.Sprintf("U+%x - U+%x", uint16(r[0]), uint16(i)), func(t *testing.T) {
					fmt.Printf("U+%x - U+%x\n DEPTH %d", uint16(r[0]), uint16(r[1]), depth)
					_, err := root.MarshalJSON()
					assert.NoError(t, err)
				})
				root = NewFugueNode(ID{"root", 0}, []uint16{0}, Root, nil)
				lastNode = root
				depth = 0
			}
		}

	}
}

func TestDeeplyNestedMarshal(t *testing.T) {
	t.Skip("Still fails with json: error calling MarshalJSON for type *v3.FugueNode: invalid character '[' exceeded max depth")
	root := NewFugueNode(ID{"root", 0}, []uint16{0}, Root, nil)
	lastNode := root

	for i := 0; i < 6_000; i++ {
		node := NewFugueNode(ID{"auth0", i}, []uint16{uint16(1)}, Left, lastNode)
		lastNode.LeftChildren = append(lastNode.LeftChildren, node)
		lastNode = node
	}

	_, err := root.MarshalJSON()
	assert.NoError(t, err)
}

func BenchmarkMultiRogueMultiplayer(b *testing.B) {
	r := rand.New(rand.NewSource(42))
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		var trees []*Rogue
		for j := 0; j < 5; j++ {
			trees = append(trees, NewRogueForQuill(fmt.Sprintf("auth%d", j)))
		}

		for i := 0; i < 5; i++ {
			var allOps []Op

			for _, tree := range trees {
				var ops []Op
				for j := 0; j < 5; j++ {
					if r.Float64() > 0.75 && tree.VisSize > 0 {
						_, dop, err := tree.RandDelete(r, 10)
						if err != nil {
							b.Fatal(err)
						}
						ops = append(ops, dop)
					} else {
						_, op, err := tree.RandInsert(r, 10)
						if err != nil {
							b.Fatal(err)
						}
						ops = append(ops, op)
					}
				}
				allOps = Interleave(r, allOps, ops)
			}

			for _, tree := range trees {
				for _, op := range allOps {
					tree.MergeOp(op)
				}
			}
		}
	}
}

func TestFormatSpans(t *testing.T) {
	testCases := []struct {
		name    string
		fn      func(*Rogue)
		visIx   int
		length  int
		format  FormatV3
		expHtml string
	}{
		{
			name: "bold span",
			fn: func(r *Rogue) {
				_, err := r.Insert(0, "Hello World!")
				require.NoError(t, err)
			},
			visIx:   0,
			length:  5,
			format:  FormatV3Span{"b": "true"},
			expHtml: `<p><strong>Hello</strong><span> World!</span></p>`,
		},
		{
			name: "unbold single char",
			fn: func(r *Rogue) {
				_, err := r.Insert(0, "Hello World!")
				require.NoError(t, err)

				_, err = r.Format(0, 5, FormatV3Span{"b": "true"})
				require.NoError(t, err)
			},
			visIx:   1,
			length:  1,
			format:  FormatV3Span{"b": ""},
			expHtml: `<p><strong>H</strong><span>e</span><strong>llo</strong><span> World!</span></p>`,
		},
		{
			name: "multiple lines to list",
			fn: func(r *Rogue) {
				_, err := r.Insert(0, "Hello\nWorld!\nHello\nWorld!")
				require.NoError(t, err)
			},
			visIx:   0,
			length:  20,
			format:  FormatV3BulletList(0),
			expHtml: `<ul><li><span>Hello</span></li><li><span>World!</span></li><li><span>Hello</span></li><li><span>World!</span></li></ul>`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := NewRogue("0")
			tc.fn(r)

			_, err := r.Format(tc.visIx, tc.length, tc.format)
			require.NoError(t, err)

			firstID, err := r.GetFirstTotID()
			require.NoError(t, err)

			lastID, err := r.GetLastTotID()
			require.NoError(t, err)

			html, err := r.GetHtml(firstID, lastID, false, false)
			require.NoError(t, err)

			require.Equal(t, tc.expHtml, html)
		})
	}
}
