package v3_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v3 "github.com/teamreviso/code/rogue/v3"
)

func TestRewind(t *testing.T) {
	tests := []struct {
		name          string
		fn            func(r *v3.Rogue, ca *v3.ContentAddress)
		startId       *v3.ID
		endId         *v3.ID
		expectedOp    v3.MultiOp
		expectedError error
		expectedHTML  string
	}{
		{
			name: "Basic",
			fn: func(r *v3.Rogue, ca *v3.ContentAddress) {
				_, err := r.Insert(0, "Hello World!")
				require.NoError(t, err)

				firstID, lastID, err := r.GetWrappingTotIDs()
				require.NoError(t, err)
				addr, err := r.GetAddress(firstID, lastID)
				require.NoError(t, err)

				*ca = *addr

				_, err = r.Delete(6, 6)
				require.NoError(t, err)

				_, err = r.Insert(6, "Friends!")
				require.NoError(t, err)
			},
			expectedHTML: "<p data-rid=\"q_1\"><span data-rid=\"1_3\">Hello World!</span></p>",
		},
		{
			name: "Selection",
			fn: func(r *v3.Rogue, ca *v3.ContentAddress) {
				_, err := r.Insert(0, "The brown fox jumps over the lazy dog.")
				require.NoError(t, err)

				firstID, lastID, err := r.GetWrappingTotIDs()
				require.NoError(t, err)
				addr, err := r.GetAddress(firstID, lastID)
				require.NoError(t, err)

				*ca = *addr

				_, err = r.Delete(4, 5)
				require.NoError(t, err)
				_, err = r.Insert(4, "white")
				require.NoError(t, err)
				_, err = r.Delete(29, 4)
				require.NoError(t, err)
				_, err = r.Insert(29, "sleeping")
				require.NoError(t, err)
			},
			startId:      &v3.ID{Author: "1", Seq: 6},
			endId:        &v3.ID{Author: "1", Seq: 12},
			expectedHTML: "<p data-rid=\"q_1\"><span data-rid=\"1_3\">The brown fox jumps over the sleeping dog.</span></p>",
		},
		{
			name: "Rewind after rewind",
			fn: func(r *v3.Rogue, ca *v3.ContentAddress) {
				_, err := r.Insert(0, "Hello World!")
				require.NoError(t, err)

				firstID, lastID, err := r.GetWrappingTotIDs()
				require.NoError(t, err)
				addr, err := r.GetAddress(firstID, lastID)
				require.NoError(t, err)

				*ca = *addr

				_, err = r.Delete(6, 6)
				require.NoError(t, err)

				_, err = r.Insert(6, "Friends!")
				require.NoError(t, err)

				_, err = r.Rewind(firstID, lastID, *addr)
				require.NoError(t, err)

				_, err = r.Delete(6, 6)
				require.NoError(t, err)

				_, err = r.Insert(6, "Enemies!")
				require.NoError(t, err)
			},
			expectedHTML: "<p data-rid=\"q_1\"><span data-rid=\"1_3\">Hello World!</span></p>",
		},
		{
			name: "Selection after rewind",
			fn: func(r *v3.Rogue, ca *v3.ContentAddress) {
				_, err := r.Insert(0, "The brown fox jumps over the lazy dog.")
				require.NoError(t, err)

				firstID, lastID, err := r.GetWrappingTotIDs()
				require.NoError(t, err)
				addr, err := r.GetAddress(firstID, lastID)
				require.NoError(t, err)

				*ca = *addr

				_, err = r.Delete(4, 5)
				require.NoError(t, err)
				_, err = r.Insert(4, "white")
				require.NoError(t, err)
				_, err = r.Delete(29, 4)
				require.NoError(t, err)
				_, err = r.Insert(29, "sleeping")
				require.NoError(t, err)

				_, err = r.Rewind(firstID, lastID, *addr)
				require.NoError(t, err)

				_, err = r.Delete(4, 5)
				require.NoError(t, err)
				_, err = r.Insert(4, "quick")
				require.NoError(t, err)

				fmt.Print(v3.ToJSON(r))
				html, err := r.GetHtmlDiff(firstID, lastID, ca, true, false)
				require.NoError(t, err)
				fmt.Print(html)
			},
			startId:      &v3.ID{Author: "1", Seq: 6},
			endId:        &v3.ID{Author: "1", Seq: 12},
			expectedHTML: "<p data-rid=\"q_1\"><span data-rid=\"1_3\">The brown fox jumps over the lazy dog.</span></p>",
		},
		{
			name: "Simple Hi, bye",
			fn: func(r *v3.Rogue, ca *v3.ContentAddress) {
				r.Author = "1"
				var err error
				_, err = r.Insert(0, "H")
				require.NoError(t, err)
				_, err = r.Insert(1, "i")

				first, last, err := r.GetWrappingIDs()
				require.NoError(t, err)

				addr, err := r.GetAddress(first, last)
				require.NoError(t, err)

				*ca = *addr

				v3.WithAuthor(r, "!1", func(ir *v3.Rogue) {
					_, err := ir.Delete(0, 2)
					require.NoError(t, err)
					_, err = ir.Insert(0, "Bye")
					require.NoError(t, err)
				})
			},
			expectedOp: v3.MultiOp{
				[]v3.Op{
					v3.RewindOp{ID: v3.ID{Author: "1", Seq: 9}, Address: v3.ContentAddress{StartID: v3.ID{Author: "1", Seq: 3}, EndID: v3.ID{Author: "q", Seq: 1}, MaxIDs: map[string]int{"1": 4, "q": 1}}, UndoAddress: v3.ContentAddress{StartID: v3.ID{Author: "1", Seq: 3}, EndID: v3.ID{Author: "q", Seq: 1}, MaxIDs: map[string]int{"!1": 8, "1": 4, "q": 1}}},
				},
			},
			expectedHTML: "<p data-rid=\"q_1\"><span data-rid=\"1_3\">Hi</span></p>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := v3.NewRogueForQuill("1")

			ca := &v3.ContentAddress{}
			tt.fn(r, ca)

			if tt.startId == nil {
				tt.startId = &v3.RootID
			}
			if tt.endId == nil {
				tt.endId = &v3.LastID
			}

			fmt.Printf("r.Rewind(%s, %s, %s)\n", v3.ToJSON(tt.startId), v3.ToJSON(tt.endId), v3.ToJSON(ca))
			ops, err := r.Rewind(*tt.startId, *tt.endId, *ca)
			if tt.expectedError != nil {
				require.Equal(t, tt.expectedError, err)
			} else {
				require.NoError(t, err)
			}
			if len(tt.expectedOp.Mops) > 0 {
				assert.Equal(t, tt.expectedOp, ops)
			}

			fmt.Printf("%s\n", v3.ToJSON(ops))

			_, err = r.MergeOp(ops)
			require.NoError(t, err, fmt.Sprintf("mop: %v", ops))

			html, err := r.GetHtml(v3.RootID, v3.LastID, true, false)
			require.NoError(t, err)

			require.Equal(t, tt.expectedHTML, html)
		})
	}
}

func TestWalkRightFromIDAtAddress(t *testing.T) {
	r := v3.NewRogueForQuill("1")
	_, err := r.Insert(0, "Hello World!")
	require.NoError(t, err)

	addr0, err := r.GetFullAddress()
	require.NoError(t, err)

	_, err = r.Delete(6, 5)
	require.NoError(t, err)

	_, err = r.Insert(6, "Friends")
	require.NoError(t, err)

	addr1, err := r.GetFullAddress()
	require.NoError(t, err)

	startID, err := r.GetFirstID()
	require.NoError(t, err)

	result := []uint16{}
	for v, err := range r.WalkRightFromAt(startID, addr0) {
		require.NoError(t, err)
		result = append(result, v.Char)
	}
	require.Equal(t, "Hello World!\n", v3.Uint16ToStr(result))

	result = []uint16{}
	for v, err := range r.WalkRightFromAt(startID, addr1) {
		require.NoError(t, err)
		result = append(result, v.Char)
	}
	require.Equal(t, "Hello Friends!\n", v3.Uint16ToStr(result))
}

func TestIDFromIDAndOffset(t *testing.T) {
	r := v3.NewRogueForQuill("1")
	_, err := r.Insert(0, "Hello World!")
	require.NoError(t, err)

	addr0, err := r.GetFullAddress()
	require.NoError(t, err)

	_, err = r.Delete(6, 5)
	require.NoError(t, err)

	_, err = r.Insert(6, "Friends")
	require.NoError(t, err)

	addr1, err := r.GetFullAddress()
	require.NoError(t, err)

	startID, err := r.GetFirstID()
	require.NoError(t, err)

	id, err := r.IDFromIDAndOffset(startID, 6, addr0)
	require.NoError(t, err)
	require.Equal(t, v3.ID{"1", 9}, id)

	id, err = r.IDFromIDAndOffset(startID, 6, addr1)
	require.NoError(t, err)
	require.Equal(t, v3.ID{"1", 16}, id)

	id, err = r.IDFromIDAndOffset(startID, 0, nil)
	require.NoError(t, err)
	require.Equal(t, v3.ID{"1", 3}, id)

	id, err = r.IDFromIDAndOffset(startID, 6, nil)
	require.NoError(t, err)
	require.Equal(t, v3.ID{"1", 16}, id)
}

func TestIDFromIDAndOffsetTable(t *testing.T) {
	type testCase struct {
		startID v3.ID
		offset  int
		expID   v3.ID
	}

	tests := []struct {
		name  string
		fn    func(r *v3.Rogue, ca *v3.ContentAddress)
		cases []testCase
	}{
		{
			name: "basic",
			fn: func(r *v3.Rogue, ca *v3.ContentAddress) {
				_, err := r.Insert(0, "Hello World!")
				require.NoError(t, err)
			},
			cases: []testCase{
				{
					startID: v3.ID{"1", 3},
					offset:  3,
					expID:   v3.ID{"1", 6},
				},
				{
					startID: v3.ID{"1", 3},
					offset:  11,
					expID:   v3.ID{"1", 14},
				},
			},
		},
		{
			name: "quotes",
			fn: func(r *v3.Rogue, ca *v3.ContentAddress) {
				_, err := r.Insert(0, "\n\"hello\"")
				require.NoError(t, err)
			},
			cases: []testCase{
				{
					startID: v3.ID{"1", 3},
					offset:  0,
					expID:   v3.ID{"1", 3},
				},
				{
					startID: v3.ID{"1", 4},
					offset:  0,
					expID:   v3.ID{"1", 4},
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%s_%d", tt.name, i), func(t *testing.T) {
			r := v3.NewRogueForQuill("1")
			tt.fn(r, nil)

			for _, c := range tt.cases {
				id, err := r.IDFromIDAndOffset(c.startID, c.offset, nil)
				require.NoError(t, err)
				require.Equal(t, c.expID, id)
			}
		})
	}
}

func TestIDsToEnclosingSpan(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(r *v3.Rogue, ca *v3.ContentAddress)
		selStartID v3.ID
		selEndID   v3.ID
		expStartID v3.ID
		expEndID   v3.ID
	}{
		{
			name: "Basic",
			fn: func(r *v3.Rogue, ca *v3.ContentAddress) {
				_, err := r.Insert(0, "Hello World!")
				require.NoError(t, err)
			},
			selStartID: v3.ID{"1", 3},
			selEndID:   v3.ID{"q", 1},
			expStartID: v3.ID{"1", 3},
			expEndID:   v3.ID{"q", 1},
		},
		{
			name: "Max Selection",
			fn: func(r *v3.Rogue, ca *v3.ContentAddress) {
				_, err := r.Insert(0, "Hello World!")
				require.NoError(t, err)
			},
			selStartID: v3.ID{"root", 0},
			selEndID:   v3.ID{"q", 1},
			expStartID: v3.ID{"root", 0},
			expEndID:   v3.ID{"q", 1},
		},
		{
			name: "Multi Line Selection Middle to Middle",
			fn: func(r *v3.Rogue, ca *v3.ContentAddress) {
				_, err := r.Insert(0, "First line\nSecond line\nThird line")
				require.NoError(t, err)
			},
			selStartID: v3.ID{"1", 5},
			selEndID:   v3.ID{"1", 16},
			expStartID: v3.ID{"1", 3},
			expEndID:   v3.ID{"1", 25},
		},
		{
			name: "Selection Within Single Line",
			fn: func(r *v3.Rogue, ca *v3.ContentAddress) {
				_, err := r.Insert(0, "First line\nSecond line\nThird line")
				require.NoError(t, err)
			},
			selStartID: v3.ID{"1", 4},
			selEndID:   v3.ID{"1", 6},
			expStartID: v3.ID{"1", 3},
			expEndID:   v3.ID{"1", 13},
		},
		{
			name: "Selection Across Three Lines",
			fn: func(r *v3.Rogue, ca *v3.ContentAddress) {
				_, err := r.Insert(0, "First line\nSecond line\nThird line")
				require.NoError(t, err)
			},
			selStartID: v3.ID{"1", 8},
			selEndID:   v3.ID{"1", 26},
			expStartID: v3.ID{"1", 3},
			expEndID:   v3.ID{"q", 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := v3.NewRogueForQuill("1")
			tt.fn(r, nil)
			startID, endID, err := r.IDsToEnclosingSpan([]v3.ID{tt.selStartID, tt.selEndID}, nil)
			require.NoError(t, err)
			require.Equal(t, tt.expStartID, startID)
			require.Equal(t, tt.expEndID, endID)
		})
	}
}

func TestWalkLeftFromAtAddress(t *testing.T) {
	tests := []struct {
		name    string
		fn      func(r *v3.Rogue, ca *v3.ContentAddress)
		exp     string
		startID v3.ID
	}{
		{
			name: "basic",
			fn: func(r *v3.Rogue, ca *v3.ContentAddress) {
				_, err := r.Insert(0, "Hello World!")
				require.NoError(t, err)
			},
			exp:     "Hello World!\n",
			startID: v3.ID{"q", 1},
		},
		{
			name: "deleted text",
			fn: func(r *v3.Rogue, ca *v3.ContentAddress) {
				_, err := r.Insert(0, "Hello World!")
				require.NoError(t, err)

				_, err = r.Delete(9, 1)
				require.NoError(t, err)

				_, err = r.Delete(4, 3)
				require.NoError(t, err)

				_, err = r.Delete(0, 3)
				require.NoError(t, err)
			},
			exp:     "lord!\n",
			startID: v3.ID{"q", 1},
		},
	}

	for _, tt := range tests {
		r := v3.NewRogueForQuill("1")
		tt.fn(r, nil)

		s := []uint16{}
		for item, err := range r.WalkLeftFromAt(tt.startID, nil) {
			require.NoError(t, err)
			s = append(s, item.Char)
		}

		v3.Reverse(s)
		require.Equal(t, tt.exp, v3.Uint16ToStr(s))
	}
}

func TestWalkRightFromAtAddress(t *testing.T) {
	tests := []struct {
		name    string
		fn      func(r *v3.Rogue, ca *v3.ContentAddress)
		exp     string
		startID v3.ID
	}{
		{
			name: "last id",
			fn: func(r *v3.Rogue, ca *v3.ContentAddress) {
				_, err := r.Insert(0, "Hello World!")
				require.NoError(t, err)
			},
			exp:     "\n",
			startID: v3.ID{"q", 1},
		},
		{
			name: "last id deleted",
			fn: func(r *v3.Rogue, ca *v3.ContentAddress) {
				_, err := r.Insert(0, "Hello World!")
				require.NoError(t, err)

				_, err = r.Delete(12, 1)
				require.NoError(t, err)
			},
			exp:     "",
			startID: v3.ID{"q", 1},
		},
		{
			name: "deleted text",
			fn: func(r *v3.Rogue, ca *v3.ContentAddress) {
				_, err := r.Insert(0, "Hello World!")
				require.NoError(t, err)

				_, err = r.Delete(9, 1)
				require.NoError(t, err)

				_, err = r.Delete(4, 3)
				require.NoError(t, err)

				_, err = r.Delete(0, 3)
				require.NoError(t, err)
			},
			exp:     "lord!\n",
			startID: v3.ID{"1", 3},
		},
	}

	for _, tt := range tests {
		r := v3.NewRogueForQuill("1")
		tt.fn(r, nil)

		s := []uint16{}
		for item, err := range r.WalkRightFromAt(tt.startID, nil) {
			require.NoError(t, err)
			s = append(s, item.Char)
		}

		require.Equal(t, tt.exp, v3.Uint16ToStr(s))
	}
}

func TestGetLine(t *testing.T) {
	tests := []struct {
		name    string
		fn      func(r *v3.Rogue, ca *v3.ContentAddress)
		id      v3.ID
		startID v3.ID
		endID   v3.ID
		offset  int
	}{
		{
			name: "only newlines",
			fn: func(r *v3.Rogue, ca *v3.ContentAddress) {
				_, err := r.Insert(0, "\n")
				require.NoError(t, err)
			},
			id:      v3.ID{"q", 1},
			startID: v3.ID{"q", 1},
			endID:   v3.ID{"q", 1},
			offset:  0,
		},
		{
			name: "basic",
			fn: func(r *v3.Rogue, ca *v3.ContentAddress) {
				_, err := r.Insert(0, "Hello World!")
				require.NoError(t, err)
			},
			id:      v3.ID{"1", 3},
			startID: v3.ID{"1", 3},
			endID:   v3.ID{"q", 1},
			offset:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := v3.NewRogueForQuill("1")
			tt.fn(r, nil)

			startID, endID, offset, err := r.GetLineAt(tt.id, nil)
			require.NoError(t, err)
			require.Equal(t, tt.startID, startID)
			require.Equal(t, tt.endID, endID)
			require.Equal(t, tt.offset, offset)
		})
	}
}
