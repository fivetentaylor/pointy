package v3_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
	v3 "github.com/fivetentaylor/pointy/rogue/v3"
)

func _printHtml(tree *v3.Rogue) {
	firstID, err := tree.GetFirstTotID()
	if err != nil {
		panic(err)
	}

	lastID, err := tree.GetLastTotID()
	if err != nil {
		panic(err)
	}

	tt, err := tree.GetHtml(firstID, lastID, true, false)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%q\n", tt)
}

func TestFuzzyUndo(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}

	r := rand.New(rand.NewSource(42))

	for i := 0; i < 5; i++ {
		var trees []*v3.Rogue
		for j := 0; j < 5; j++ {
			trees = append(trees, v3.NewRogueForQuill(fmt.Sprintf("auth%d", j)))
		}

		for k := 0; k < 100; k++ {
			var allOps []v3.Op

			for _, tree := range trees {
				var ops []v3.Op

				// do a few ops in a row on each tree so things
				// have a chance to get jumbled some more
				for j := 0; j < 5; j++ {
					prob := r.Float64()

					if prob < 0.1 && tree.VisSize > 0 {
						_, dop, err := tree.RandDelete(r, 10)
						require.NoError(t, err)

						ops = append(ops, dop)
					} else if prob < 0.2 && tree.VisSize > 2 {
						/*_, fop, err := tree.RandFormat(r, 10)
						require.NoError(t, err)

						ops = append(ops, fop)*/
					} else if prob < 0.3 && tree.VisSize > 1 {
						// random undo
						firstID, err := tree.GetFirstTotID()
						require.NoError(t, err)

						lastID, err := tree.GetLastTotID()
						require.NoError(t, err)

						uop, err := tree.Undo(firstID, lastID)
						require.NoError(t, err)

						ops = append(ops, uop)
					} else if prob <= 1.0 {
						_, op, err := tree.RandInsert(r, 10)
						require.NoError(t, err)

						ops = append(ops, op)
					}
				}

				// randomly interleave the ops to make this more realistic
				allOps = v3.Interleave(r, allOps, ops)
			}

			/*for _, op := range allOps {
				fmt.Printf("op: %v\n", op)
			}*/

			for _, tree := range trees {
				for _, op := range allOps {
					tree.MergeOp(op)
				}
			}
		}

		/*for _, tree := range trees {
			_printHtml(tree)
		}*/

		firstID, err := trees[0].GetFirstTotID()
		require.NoError(t, err)

		lastID, err := trees[0].GetLastTotID()
		require.NoError(t, err)

		tot0, err := trees[0].Rope.GetTotBetween(firstID, lastID)
		require.NoError(t, err)

		firstID, err = trees[1].GetFirstTotID()
		require.NoError(t, err)

		lastID, err = trees[1].GetLastTotID()
		require.NoError(t, err)

		tot1, err := trees[1].Rope.GetTotBetween(firstID, lastID)
		require.NoError(t, err)

		// DEBUG
		for i := 0; i < min(len(tot0.Text), len(tot1.Text)); i++ {
			if tot0.IsDeleted[i] != tot1.IsDeleted[i] {
				ch0 := trees[0].CharHistory[tot0.IDs[i]]
				fmt.Println("ch0")
				ch0.Dft(func(m *v3.Marker) error {
					fmt.Printf("%v\n", m)
					return nil
				})
				fmt.Println()

				ch1 := trees[1].CharHistory[tot1.IDs[i]]
				fmt.Println("ch1")
				ch1.Dft(func(m *v3.Marker) error {
					fmt.Printf("%v\n", m)
					return nil
				})

				require.Fail(t, "tot0 and tot1 differ")
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

		_, totIx, err := trees[0].Rope.GetIndex(v3.RootID)
		require.NoError(t, err)
		fmt.Printf("RootID totIx: %d\n", totIx)

		serRogue, err := json.Marshal(trees[0])
		require.NoError(t, err)

		var deRogue v3.Rogue
		err = json.Unmarshal(serRogue, &deRogue)
		// fmt.Printf("serRogue: %q\n", string(serRogue)[:100])
		require.NoError(t, err)

		firstID, err = deRogue.GetFirstTotID()
		require.NoError(t, err)

		lastID, err = deRogue.GetLastTotID()
		require.NoError(t, err)

		dTreeText, err := deRogue.GetHtml(firstID, lastID, true, false)
		require.NoError(t, err)

		for _, tree := range trees {
			firstID, err := deRogue.GetFirstTotID()
			require.NoError(t, err)

			lastID, err := deRogue.GetLastTotID()
			require.NoError(t, err)

			tt, err := tree.GetHtml(firstID, lastID, true, false)
			require.NoError(t, err)

			require.Equal(t, dTreeText, tt)
		}
	}
}

func TestUndoRedoToEmpty(t *testing.T) {
	t.Parallel()

	r := v3.NewRogueForQuill("0")

	_, err := r.Insert(0, "abc")
	require.NoError(t, err)

	_, err = r.Format(0, 3, v3.FormatV3Span{"b": "true"})
	require.NoError(t, err)

	firstID, err := r.GetFirstTotID()
	require.NoError(t, err)

	lastID, err := r.GetLastTotID()
	require.NoError(t, err)

	_, err = r.Undo(firstID, lastID)
	require.NoError(t, err)
	require.Equal(t, "<p><span>abc</span></p>", r.DisplayAllHtml(false, false))

	_, err = r.Redo()
	require.NoError(t, err)
	require.Equal(t, "<p><strong>abc</strong></p>", r.DisplayAllHtml(false, false))

	_, err = r.Undo(firstID, lastID)
	require.NoError(t, err)
	require.Equal(t, "<p><span>abc</span></p>", r.DisplayAllHtml(false, false))

	err = r.Formats.Validate()
	require.NoError(t, err)
}

func TestUndoLink(t *testing.T) {
	t.Parallel()

	r := v3.NewRogueForQuill("0")

	_, err := r.Insert(0, "google.com")
	require.NoError(t, err)

	_, err = r.Format(0, 10, v3.FormatV3Span{"a": "https://www.google.com"})
	require.NoError(t, err)

	require.Equal(t, "<p><a href=\"https://www.google.com\">google.com</a></p>", r.DisplayAllHtml(false, false))

	firstID, err := r.GetFirstTotID()
	require.NoError(t, err)

	lastID, err := r.GetLastTotID()
	require.NoError(t, err)

	_, err = r.Undo(firstID, lastID)
	require.NoError(t, err)

	require.Equal(t, "<p><span>google.com</span></p>", r.DisplayAllHtml(false, false))

	err = r.Formats.Validate()
	require.NoError(t, err)
}

func TestUndoStickyAndNoSticky(t *testing.T) {
	t.Parallel()

	r := v3.NewRogueForQuill("0")

	_, err := r.Insert(0, "google.com")
	require.NoError(t, err)

	_, err = r.Format(0, 10, v3.FormatV3Span{"a": "https://www.google.com"})
	require.NoError(t, err)

	require.Equal(t, "<p><a href=\"https://www.google.com\">google.com</a></p>", r.DisplayAllHtml(false, false))

	_, err = r.Format(0, 6, v3.FormatV3Span{"s": "true"})
	require.NoError(t, err)

	require.Equal(t, "<p><a href=\"https://www.google.com\"><s>google</s></a><a href=\"https://www.google.com\">.com</a></p>", r.DisplayAllHtml(false, false))

	fmt.Println("UNDO")
	_, err = r.UndoDoc()
	require.NoError(t, err)
	require.Equal(t, "<p><a href=\"https://www.google.com\">google.com</a></p>", r.DisplayAllHtml(false, false))

	fmt.Println("UNDO")
	_, err = r.UndoDoc()
	require.NoError(t, err)
	require.Equal(t, "<p><span>google.com</span></p>", r.DisplayAllHtml(false, false))

	fmt.Println("REDO")
	_, err = r.Redo()
	require.NoError(t, err)
	require.Equal(t, "<p><a href=\"https://www.google.com\">google.com</a></p>", r.DisplayAllHtml(false, false))

	fmt.Println("REDO")
	_, err = r.Redo()
	require.NoError(t, err)
	require.Equal(t, "<p><a href=\"https://www.google.com\"><s>google</s></a><a href=\"https://www.google.com\">.com</a></p>", r.DisplayAllHtml(false, false))

	fmt.Println(r.OpIndex.String())

	fmt.Println("UNDO")
	_, err = r.UndoDoc()
	require.NoError(t, err)
	require.Equal(t, "<p><a href=\"https://www.google.com\">google.com</a></p>", r.DisplayAllHtml(false, false))

	fmt.Println("NoSticky")
	r.Formats.NoSticky.Print()

	fmt.Printf("\nSticky\n")
	r.Formats.Sticky.Print()

	fmt.Println("UNDO")
	_, err = r.UndoDoc()
	require.NoError(t, err)
	require.Equal(t, "<p><span>google.com</span></p>", r.DisplayAllHtml(false, false))

	err = r.Formats.Validate()
	require.NoError(t, err)
}
