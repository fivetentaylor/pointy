package v3_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
	v3 "github.com/teamreviso/code/rogue/v3"
)

func TestRewindTo(t *testing.T) {
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
