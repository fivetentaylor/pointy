package v3_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
	v3 "github.com/teamreviso/code/rogue/v3"
)

func TestGetAddress(t *testing.T) {
	r := v3.NewRogue("auth0")

	_, err := r.Insert(0, "hello world")
	require.NoError(t, err)

	_, err = r.Insert(5, " cruel")
	require.NoError(t, err)

	startID, err := r.Rope.GetTotID(0)
	require.NoError(t, err)

	endID, err := r.Rope.GetTotID(r.TotSize - 1)
	require.NoError(t, err)

	addr, err := r.GetAddress(startID, endID)
	require.NoError(t, err)

	fmt.Printf("Address: %v\n", addr)

	rb, err := json.Marshal(r)
	require.NoError(t, err)

	fmt.Printf("Rogue: %s\n", rb)
}

func TestGetAddressText(t *testing.T) {
	r := v3.NewRogueForQuill("auth0")

	_, err := r.Insert(0, "hello world")
	require.NoError(t, err)

	_, err = r.Insert(5, " cruel")
	require.NoError(t, err)

	startID, err := r.Rope.GetTotID(0)
	require.NoError(t, err)

	endID, err := r.Rope.GetTotID(r.TotSize - 1)
	require.NoError(t, err)

	addr, err := r.GetAddress(startID, endID)
	require.NoError(t, err)

	_, err = r.Delete(5, len(" cruel"))
	require.NoError(t, err)

	_, err = r.Insert(5, " friendly")
	require.NoError(t, err)

	addrRogue, err := r.GetAddressRogue(addr)
	require.NoError(t, err)

	require.Equal(t, "hello cruel world\n", addrRogue.GetText())
	require.Equal(t, "hello friendly world\n", r.GetText())
}

func TestGetAddressMarkdown(t *testing.T) {
	r := v3.NewRogueForQuill("auth0")

	_, err := r.Insert(0, "hello world")
	require.NoError(t, err)

	_, err = r.Insert(5, " cruel")
	require.NoError(t, err)

	_, err = r.Format(6, 5, v3.FormatV3Span{"b": "true"})
	require.NoError(t, err)

	startID, err := r.Rope.GetTotID(0)
	require.NoError(t, err)

	endID, err := r.Rope.GetTotID(r.TotSize - 1)
	require.NoError(t, err)

	addr, err := r.GetAddress(startID, endID)
	require.NoError(t, err)

	_, err = r.Format(6, 5, v3.FormatV3Span{"b": "null"})
	require.NoError(t, err)

	_, err = r.Delete(5, len(" cruel"))
	require.NoError(t, err)

	_, err = r.Insert(5, " friendly")
	require.NoError(t, err)

	addressRogue, err := r.GetAddressRogue(addr)
	require.NoError(t, err)

	fmt.Printf("Address Rogue: %s\n", addressRogue.GetText())

	startID, err = addressRogue.Rope.GetTotID(0)
	require.NoError(t, err)

	endID, err = addressRogue.Rope.GetTotID(addressRogue.TotSize - 1)
	require.NoError(t, err)

	text, err := addressRogue.GetMarkdownBeforeAfter(startID, endID)
	require.NoError(t, err)

	fmt.Printf("Address Markdown: %s\n", text)

	require.Equal(t, "hello **cruel** world\n\n", text)
	require.Equal(t, "hello friendly world\n", r.GetText())
}

type ContentAddressText struct {
	ContentAddress *v3.ContentAddress `json:"content_address"`
	Vis            *v3.FugueVis
	Tot            *v3.FugueTot
}

func TestContentAddressFuzz(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}

	rng := rand.New(rand.NewSource(42))
	cas := []ContentAddressText{}

	var trees []*v3.Rogue
	for j := 0; j < 5; j++ {
		trees = append(trees, v3.NewRogueForQuill(fmt.Sprintf("a%d", j)))
	}

	doc := trees[0]

	for k := 0; k < 1000; k++ {
		var allOps []v3.Op

		for _, tree := range trees {
			var ops []v3.Op

			// do a few ops in a row on each tree so things
			// have a chance to get jumbled some more
			for j := 0; j < 5; j++ {
				if rng.Float64() > 0.75 && tree.VisSize > 0 {
					_, dop, err := tree.RandDelete(rng, 10)
					require.NoError(t, err)

					ops = append(ops, dop)
				} else {
					_, op, err := tree.RandInsert(rng, 10)
					require.NoError(t, err)

					ops = append(ops, op)
				}
			}

			// randomly interleave the ops to make this more realistic
			allOps = v3.Interleave(rng, allOps, ops)
		}

		for _, tree := range trees {
			for _, op := range allOps {
				tree.MergeOp(op)
			}
		}

		// take a content address 10% of the time
		if rng.Float64() < 0.1 {
			startIx := max(rng.Intn(doc.VisSize)-10, 0)
			endIx := startIx + rng.Intn(doc.VisSize-startIx)

			startChar, err := doc.GetChar(startIx)
			require.NoError(t, err)

			if v3.IsLowSurrogate(startChar) {
				startIx--
			}

			endChar, err := doc.GetChar(endIx)
			require.NoError(t, err)

			if v3.IsHighSurrogate(endChar) {
				endIx++
			}

			startID, err := doc.Rope.GetVisID(startIx)
			require.NoError(t, err)

			endID, err := doc.Rope.GetVisID(endIx)
			require.NoError(t, err)

			ca, err := doc.GetAddress(startID, endID)
			require.NoError(t, err)

			vis, err := doc.Rope.GetBetween(startID, endID)
			require.NoError(t, err)

			tot, err := doc.Rope.GetTotBetween(startID, endID)
			require.NoError(t, err)

			cas = append(cas, ContentAddressText{
				ContentAddress: ca,
				Vis:            vis,
				Tot:            tot,
			})
		}
	}

	require.Greater(t, trees[0].VisSize, 0)
	require.Greater(t, len(cas), 0)
	fmt.Printf("len(cas): %d\n", len(cas))

	for i, ca := range cas {
		fmt.Printf("ITERATION: %d\n", i)
		visText := v3.Uint16ToStr(ca.Vis.Text)

		ar, err := doc.GetAddressRogue(ca.ContentAddress)
		require.NoError(t, err)

		fmt.Printf("ar.RootSeqs: %v\n", ar.RootSeqs)
		fmt.Printf("ar.Roots: %v\n", ar.Roots)

		totText := v3.Uint16ToStr(ca.Tot.Text)
		require.Equal(t, totText, ar.GetTotText())
		require.Equal(t, visText, ar.GetText())

		startID, endID := ca.ContentAddress.StartID, ca.ContentAddress.EndID
		filterText, err := doc.Filter(startID, endID, ca.ContentAddress)
		require.NoError(t, err)
		require.Equal(t, visText, v3.Uint16ToStr(filterText.Text))

		walkLeft := []uint16{}
		for v, err := range doc.WalkLeftFromAt(endID, ca.ContentAddress) {
			require.NoError(t, err)
			walkLeft = append(walkLeft, v.Char)

			if v.ID == startID {
				break
			}
		}
		v3.Reverse(walkLeft)
		require.Equal(t, visText, v3.Uint16ToStr(walkLeft))

		walkRight := []uint16{}
		for v, err := range doc.WalkRightFromAt(startID, ca.ContentAddress) {
			require.NoError(t, err)
			walkRight = append(walkRight, v.Char)
			if v.ID == endID {
				break
			}
		}
		require.Equal(t, visText, v3.Uint16ToStr(walkRight))

		// TODO: Add this back if we ever to compact again
		/*car, err := ar.Compact()
		require.NoError(t, err)

		require.Equal(t, ar.GetText(), car.GetText())*/
	}
}

func TestParseOldAddress(t *testing.T) {
	oldAddress := `{"startID":["!1",16635],"endID":["q",1],"maxIDs":[{"key":"!1","value":16650},{"key":"1","value":22314},{"key":"q","value":1}]}`

	ca := v3.ContentAddress{}
	err := json.Unmarshal([]byte(oldAddress), &ca)
	require.NoError(t, err)
}

func TestFastContentAddressFuzz(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}

	rng := rand.New(rand.NewSource(42))
	cas := []ContentAddressText{}

	var trees []*v3.Rogue
	for j := 0; j < 5; j++ {
		trees = append(trees, v3.NewRogueForQuill(fmt.Sprintf("a%d", j)))
	}

	doc := trees[0]

	for k := 0; k < 1000; k++ {
		var allOps []v3.Op

		for _, tree := range trees {
			var ops []v3.Op

			// do a few ops in a row on each tree so things
			// have a chance to get jumbled some more
			for j := 0; j < 5; j++ {
				if rng.Float64() > 0.75 && tree.VisSize > 0 {
					_, dop, err := tree.RandDelete(rng, 10)
					require.NoError(t, err)

					ops = append(ops, dop)
				} else {
					_, op, err := tree.RandInsert(rng, 10)
					require.NoError(t, err)

					ops = append(ops, op)
				}
			}

			// randomly interleave the ops to make this more realistic
			allOps = v3.Interleave(rng, allOps, ops)
		}

		for _, tree := range trees {
			for _, op := range allOps {
				tree.MergeOp(op)
			}
		}

		// take a content address 10% of the time
		if rng.Float64() < 0.1 {
			startID, err := doc.GetFirstTotID()
			require.NoError(t, err)

			endID, err := doc.GetLastTotID()
			require.NoError(t, err)

			sca, err := doc.GetAddress(startID, endID)
			require.NoError(t, err)

			ca, err := doc.GetFullAddress()
			require.NoError(t, err)

			// Make sure the fast address and the old address are equal
			require.Equal(t, ca, sca)

			vis, err := doc.Rope.GetBetween(startID, endID)
			require.NoError(t, err)

			tot, err := doc.Rope.GetTotBetween(startID, endID)
			require.NoError(t, err)

			cas = append(cas, ContentAddressText{
				ContentAddress: ca,
				Vis:            vis,
				Tot:            tot,
			})
		}
	}

	require.Greater(t, trees[0].VisSize, 0)
	require.Greater(t, len(cas), 0)
	fmt.Printf("len(cas): %d\n", len(cas))

	for i, ca := range cas {
		fmt.Printf("ITERATION: %d\n", i)
		visText := v3.Uint16ToStr(ca.Vis.Text)
		filterText, err := doc.Filter(ca.ContentAddress.StartID, ca.ContentAddress.EndID, ca.ContentAddress)
		require.NoError(t, err)
		require.Equal(t, visText, v3.Uint16ToStr(filterText.Text))

		// TODO: Add this back if we ever to compact again
		/*car, err := ar.Compact()
		require.NoError(t, err)

		require.Equal(t, ar.GetText(), car.GetText())*/
	}
}

func TestContentAddressFuzzSinglePlayer(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}

	type ContentAddressHtml struct {
		ContentAddress *v3.ContentAddress
		Html           string
		Plaintext      string
	}

	rng := rand.New(rand.NewSource(42))
	cas := []ContentAddressHtml{}

	doc := v3.NewRogueForQuill("0")

	deletes, inserts, undos, formats := 0, 0, 0, 0
	for k := 0; k < 2000; k++ {
		prob := rng.Float64()

		if prob < 0.1 && doc.VisSize > 0 {
			deletes++
			_, _, err := doc.RandDelete(rng, 10)
			require.NoError(t, err)
		} else if prob < 0.3 && doc.VisSize > 2 {
			formats++
			_, _, err := doc.RandFormat(rng, 10)
			require.NoError(t, err)
		} else if prob < 0.4 && doc.VisSize > 1 {
			undos++
			_, err := doc.UndoDoc()
			require.NoError(t, err)
		} else if prob <= 1.0 {
			inserts++
			_, _, err := doc.RandInsert(rng, 10)
			require.NoError(t, err)
		}

		prob = rng.Float64()

		// take a content address 10% of the time
		if rng.Float64() < 0.1 {
			ca, err := doc.GetFullAddress()
			require.NoError(t, err)

			startID, err := doc.GetFirstTotID()
			require.NoError(t, err)

			endID, err := doc.GetLastTotID()
			require.NoError(t, err)

			html, err := doc.GetHtml(startID, endID, true, false)
			require.NoError(t, err)

			cas = append(cas, ContentAddressHtml{
				ContentAddress: ca,
				Html:           html,
				Plaintext:      doc.GetText(),
			})
		}
	}

	fmt.Printf("deletes: %d, inserts: %d, undos: %d, formats: %d\n", deletes, inserts, undos, formats)

	firstID, err := doc.GetFirstTotID()
	require.NoError(t, err)

	lastID, err := doc.GetLastTotID()
	require.NoError(t, err)

	for _, ca := range cas {
		plaintext, err := doc.Filter(firstID, lastID, ca.ContentAddress)
		require.NoError(t, err)

		require.Equal(t, ca.Plaintext, v3.Uint16ToStr(plaintext.Text))

		html, err := doc.GetHtmlAt(firstID, lastID, ca.ContentAddress, true, false)
		require.NoError(t, err)

		require.Equal(t, ca.Html, html)

		r2, err := doc.Compact(ca.ContentAddress)
		require.NoError(t, err)

		r2StartID, err := r2.GetFirstTotID()
		require.NoError(t, err)

		r2EndID, err := r2.GetLastTotID()
		require.NoError(t, err)

		r2Html, err := r2.GetHtml(r2StartID, r2EndID, false, false)
		require.NoError(t, err)

		htmlNoID, err := doc.GetHtmlAt(firstID, lastID, ca.ContentAddress, false, false)
		require.NoError(t, err)

		require.Equal(t, htmlNoID, r2Html)
	}
}
