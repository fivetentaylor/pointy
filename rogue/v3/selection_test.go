package v3_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	v3 "github.com/fivetentaylor/pointy/rogue/v3"
)

func Test_GetSelection(t *testing.T) {
	tests := []struct {
		name              string
		fn                func(*v3.Rogue)
		startId           v3.ID
		endId             v3.ID
		expectedSelection *v3.Selection
	}{
		{
			name:    "empty doc",
			fn:      func(r *v3.Rogue) {},
			startId: v3.ID{"root", 0},
			endId:   v3.ID{"q", 1},
			expectedSelection: &v3.Selection{
				StartSpanID:     v3.ID{"q", 1},
				StartSpanOffset: 0,
				EndSpanID:       v3.ID{"q", 1},
				EndSpanOffset:   0,
			},
		},
		{
			name:    "simple",
			fn:      func(r *v3.Rogue) { r.Insert(0, "hello") },
			startId: v3.ID{"root", 0},
			endId:   v3.ID{"q", 1},
			expectedSelection: &v3.Selection{
				StartSpanID:     v3.ID{"0", 3},
				StartSpanOffset: 0,
				EndSpanID:       v3.ID{"0", 3},
				EndSpanOffset:   5,
			},
		},
		{
			name: "with returns",
			fn: func(r *v3.Rogue) {
				r.Insert(0, "hello\nhow\nare\nyou\ndoing")
			},
			startId: v3.ID{"0", 5},
			endId:   v3.ID{"0", 22},
			expectedSelection: &v3.Selection{
				StartSpanID:     v3.ID{"0", 3},
				StartSpanOffset: 2,
				EndSpanID:       v3.ID{"0", 21},
				EndSpanOffset:   1,
			},
		},
		{
			name: "with quotes",
			fn: func(r *v3.Rogue) {
				r.Insert(0, "\n\"hello\"")
			},
			startId: v3.ID{"0", 5},
			endId:   v3.ID{"0", 10},
			expectedSelection: &v3.Selection{
				StartSpanID:     v3.ID{"0", 4},
				StartSpanOffset: 1,
				EndSpanID:       v3.ID{"0", 10},
				EndSpanOffset:   0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := v3.NewRogueForQuill("0")
			tt.fn(r)

			selection, err := r.GetSelection(tt.startId, tt.endId, nil, true)
			require.NoError(t, err)
			require.Equal(t, tt.expectedSelection, selection)
		})
	}
}
