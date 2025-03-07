package v3_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	v3 "github.com/fivetentaylor/pointy/rogue/v3"
)

func TestFailedOps(t *testing.T) {
	r := v3.NewRogue("0")

	r.Insert(0, "Hello, world!")
	r.Delete(0, 5)

	badOps := []v3.Op{
		v3.InsertOp{
			ID:       v3.ID{"1", 10},
			Text:     "I'm a bad boy",
			ParentID: v3.ID{"1", 8},
			Side:     v3.Left,
		},
		v3.MultiOp{[]v3.Op{
			v3.InsertOp{
				ID:       v3.ID{"1", 12},
				Text:     "So bad",
				ParentID: v3.ID{"1", 9},
				Side:     v3.Right,
			},
			v3.DeleteOp{
				ID:         v3.ID{"1", 15},
				TargetID:   v3.ID{"1", 8},
				SpanLength: 5,
			},
			v3.FormatOp{
				ID:      v3.ID{"1", 18},
				StartID: v3.ID{"1", 8},
				EndID:   v3.ID{"1", 10},
				Format:  v3.FormatV3Span{"b": "true"},
			},
		}},
	}

	for _, op := range badOps {
		_, err := r.MergeOp(op)
		require.Error(t, err)
	}

	failedOps := r.FailedOps.AsSlice()
	require.Equal(t, badOps, failedOps)

	allOps, err := r.ToOps()
	require.NoError(t, err)
	require.Equal(t, 2+len(badOps), len(allOps))
}
