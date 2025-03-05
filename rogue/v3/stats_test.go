package v3_test

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
	v3 "github.com/teamreviso/code/rogue/v3"
)

func Test_OpStats(t *testing.T) {
	type toOpsTC struct {
		name     string
		fn       func(*v3.Rogue)
		expected v3.OpStats
	}

	testCases := []toOpsTC{
		{
			name: "empty",
			fn:   func(*v3.Rogue) {},
			expected: v3.OpStats{
				Inserts:              []int{0},
				Deletes:              []int{0},
				InsertsByPrefix:      map[string][]int{},
				DeletesByPrefix:      map[string][]int{},
				CurrentCharsByPrefix: map[string]int{},
				Segments:             1,
			},
		},
		{
			name: "simple insert",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "hello world")
				require.NoError(t, err)
			},
			expected: v3.OpStats{
				Inserts: []int{11},
				Deletes: []int{0},
				InsertsByPrefix: map[string][]int{
					"": {11},
				},
				CurrentCharsByPrefix: map[string]int{
					"": 11,
				},
				DeletesByPrefix: map[string][]int{},
				Segments:        1,
			},
		},
		{
			name: "simple insert with format",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "hello world")
				require.NoError(t, err)
				_, err = r.Format(6, 5, v3.FormatV3Span{"b": "true"})
				require.NoError(t, err)
			},
			expected: v3.OpStats{
				Inserts: []int{11},
				Deletes: []int{0},
				InsertsByPrefix: map[string][]int{
					"": {11},
				},
				CurrentCharsByPrefix: map[string]int{
					"": 11,
				},
				DeletesByPrefix: map[string][]int{},
				Segments:        1,
			},
		},
		{
			name: "insert emojis in emojis",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "🥔🥒🥔")
				require.NoError(t, err)
				_, err = r.Insert(2, "🤡😿😝")
				require.NoError(t, err)
			},
			expected: v3.OpStats{
				Inserts: []int{24},
				Deletes: []int{0},
				InsertsByPrefix: map[string][]int{
					"": {24},
				},
				DeletesByPrefix: map[string][]int{},
				CurrentCharsByPrefix: map[string]int{
					"": 12,
				},
				Segments: 1,
			},
		},
		{
			name: "insert and deletes",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "hello world")
				require.NoError(t, err)
				_, err = r.Delete(3, 8)
				require.NoError(t, err)
			},
			expected: v3.OpStats{
				Inserts: []int{11},
				Deletes: []int{8},
				InsertsByPrefix: map[string][]int{
					"": {11},
				},
				DeletesByPrefix: map[string][]int{
					"": {8},
				},
				CurrentCharsByPrefix: map[string]int{
					"": 3,
				},
				Segments: 1,
			},
		},
		{
			name: "insert and deletes with different author prefixes",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "hello world")
				require.NoError(t, err)
				_, err = r.Delete(3, 8)
				require.NoError(t, err)

				r.Author = "!0"
				_, err = r.Insert(4, " i am hungry")
				require.NoError(t, err)
				_, err = r.Delete(3, 8)
				require.NoError(t, err)

				r.Author = "#0"
				_, err = r.Insert(4, " nother prefix! omg")
				require.NoError(t, err)
				_, err = r.Delete(10, 5)
				require.NoError(t, err)
			},
			expected: v3.OpStats{
				Inserts: []int{42},
				Deletes: []int{21},
				InsertsByPrefix: map[string][]int{
					"":  {11},
					"!": {12},
					"#": {19},
				},
				DeletesByPrefix: map[string][]int{
					"":  {8},
					"!": {8},
					"#": {5},
				},
				CurrentCharsByPrefix: map[string]int{
					"":  3,
					"!": 5,
					"#": 14,
				},
				Segments: 1,
			},
		},
		{
			name: "insert and deletes over multiple segments",
			fn: func(r *v3.Rogue) {
				rn := rand.New(rand.NewSource(1))
				for i := 0; i < v3.OpsPerStatSegment*2; i++ {
					_, _, err := r.RandInsert(rn, 10)
					require.NoError(t, err)
					_, _, err = r.RandDelete(rn, 10)
					require.NoError(t, err)
				}
			},
			expected: v3.OpStats{Inserts: []int{165, 163, 153, 161, 143, 150, 148, 162, 149, 155, 146, 157, 131, 170, 139, 138, 168, 155, 160, 141, 150, 85}, Deletes: []int{67, 87, 84, 69, 89, 105, 82, 88, 34, 106, 88, 71, 117, 91, 84, 118, 71, 85, 123, 83, 82, 46}, InsertsByPrefix: map[string][]int{"": {165, 163, 153, 161, 143, 150, 148, 162, 149, 155, 146, 157, 131, 170, 139, 138, 168, 155, 160, 141, 150, 85}}, DeletesByPrefix: map[string][]int{"": {67, 87, 84, 69, 89, 105, 82, 88, 34, 106, 88, 71, 117, 91, 84, 118, 71, 85, 123, 83, 82, 46}}, CurrentCharsByPrefix: map[string]int{"": 365}, Segments: 22},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := v3.NewRogueForQuill("0")
			tc.fn(r)

			ops, err := r.OpStats()
			require.NoError(t, err)

			require.Equal(t, tc.expected, *ops, "\nexpected: %#v\n\ngot: %#v", tc.expected, ops)
		})
	}
}
