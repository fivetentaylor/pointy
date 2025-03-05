package v3_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	v3 "github.com/teamreviso/code/rogue/v3"
)

func TestRenderSpan(t *testing.T) {
	tests := []struct {
		name    string
		fn      func(r *v3.Rogue) v3.Op
		expSpan *v3.RenderSpan
	}{
		{
			name: "basic render span",
			fn: func(r *v3.Rogue) v3.Op {
				_, err := r.Insert(0, "a\nb\nc\nd\ne")
				require.NoError(t, err)

				op, err := r.Format(4, 3, v3.FormatV3BulletList(0))
				require.NoError(t, err)

				return op
			},
			expSpan: &v3.RenderSpan{
				FirstBlockID: v3.ID{"0", 7},
				LastBlockID:  v3.ID{"0", 10},
				ToStartID:    v3.ID{"0", 7},
				ToEndID:      v3.ID{"0", 10},
				Html:         "<ul><li data-rid=\"0_8\"><span data-rid=\"0_7\">c</span></li><li data-rid=\"0_10\"><span data-rid=\"0_9\">d</span></li></ul>",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := v3.NewRogueForQuill("0")
			op := tt.fn(r)

			span, err := r.RenderOp(op)
			require.NoError(t, err)

			require.Equal(t, tt.expSpan, span)
		})
	}
}

func TestGetBlockAt(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(r *v3.Rogue) *v3.ContentAddress
		getAtID    v3.ID
		expStartID v3.ID
		expEndID   v3.ID
	}{
		{
			name: "get code block",
			fn: func(r *v3.Rogue) *v3.ContentAddress {
				_, err := r.Insert(0, "a\nb\nc\nd\ne")
				require.NoError(t, err)

				_, err = r.Format(2, 5, v3.FormatV3CodeBlock("python"))
				require.NoError(t, err)

				return nil
			},
			getAtID:    v3.ID{"0", 7},
			expStartID: v3.ID{"0", 5},
			expEndID:   v3.ID{"0", 10},
		},
		{
			name: "get code block at end",
			fn: func(r *v3.Rogue) *v3.ContentAddress {
				_, err := r.Insert(0, "a\nb\nc\nd\ne")
				require.NoError(t, err)

				_, err = r.Format(2, 5, v3.FormatV3CodeBlock("python"))
				require.NoError(t, err)

				ca, err := r.GetFullAddress()
				require.NoError(t, err)

				return ca
			},
			getAtID:    v3.ID{"0", 10},
			expStartID: v3.ID{"0", 5},
			expEndID:   v3.ID{"0", 10},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := v3.NewRogueForQuill("0")

			ca := tt.fn(r)

			startID, endID, err := r.GetBlockAt(tt.getAtID, ca)
			require.NoError(t, err)

			firstID, err := r.GetFirstTotID()
			require.NoError(t, err)

			lastID, err := r.GetLastTotID()
			require.NoError(t, err)

			html, err := r.GetHtmlAt(firstID, lastID, ca, true, true)
			require.NoError(t, err)

			fmt.Printf("ca  : %v\n", ca)
			fmt.Printf("html: %q\n", html)

			require.Equal(t, tt.expStartID, startID)
			require.Equal(t, tt.expEndID, endID)
		})
	}
}
