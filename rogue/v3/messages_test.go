package v3

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalFormat(t *testing.T) {
	bytes := []byte(`[[2,["8cd8-qg",130],["8cd8-mq",72],["8cd8-mq",77],{"underline":true}]]`)

	var ops []FormatOp
	err := json.Unmarshal(bytes, &ops)
	assert.NoError(t, err)
}

func TestMarshalUnmarshalFormatOpV3(t *testing.T) {
	cases := []struct {
		name         string
		inputJSON    string
		expectedJSON string
	}{
		{
			name:         "Simple Ordered List",
			inputJSON:    `[2,["2*",745],["2*",729],["2*",729],{"ol":3}]`,
			expectedJSON: `[2,["2*",745],["2*",729],["2*",729],{"ol":"3"}]`,
		},
		{
			name:         "Simple Bullet List",
			inputJSON:    `[2,["2*",745],["2*",729],["2*",729],{"ul":"5"}]`,
			expectedJSON: `[2,["2*",745],["2*",729],["2*",729],{"ul":"5"}]`,
		},
		{
			name:         "Null list",
			inputJSON:    `[2,["2*",745],["2*",729],["2*",729],{"ul":"null"}]`,
			expectedJSON: `[2,["2*",745],["2*",729],["2*",729],{}]`,
		},
		{
			name:         "Empty list and null header",
			inputJSON:    `[2,["2*",745],["2*",729],["2*",729],{"ul":"","h":"null"}]`,
			expectedJSON: `[2,["2*",745],["2*",729],["2*",729],{}]`,
		},
		{
			name:         "Old style Bullet List",
			inputJSON:    `[2,["2*",745],["2*",729],["2*",729],{"list":"bullet"}]`,
			expectedJSON: `[2,["2*",745],["2*",729],["2*",729],{"ul":"0"}]`,
		},
		{
			name:         "Old style Bullet List With Indent",
			inputJSON:    `[2,["2*",745],["2*",729],["2*",729],{"list":"bullet","indent":2}]`,
			expectedJSON: `[2,["2*",745],["2*",729],["2*",729],{"ul":"2"}]`,
		},
		{
			name:         "Old style Bullet List With String Indent",
			inputJSON:    `[2,["2*",745],["2*",729],["2*",729],{"list":"bullet","indent":"2"}]`,
			expectedJSON: `[2,["2*",745],["2*",729],["2*",729],{"ul":"2"}]`,
		},
		{
			name:         "Old style Bullet List With String Indent",
			inputJSON:    `[2,["2*",745],["2*",729],["2*",729],{"list":"bullet","indent":"2"}]`,
			expectedJSON: `[2,["2*",745],["2*",729],["2*",729],{"ul":"2"}]`,
		},
		{
			name:         "New Header1",
			inputJSON:    `[2,["2*",745],["2*",729],["2*",729],{"h":1}]`,
			expectedJSON: `[2,["2*",745],["2*",729],["2*",729],{"h":"1"}]`,
		},
		{
			name:         "Old Header1",
			inputJSON:    `[2,["2*",745],["2*",729],["2*",729],{"header":1}]`,
			expectedJSON: `[2,["2*",745],["2*",729],["2*",729],{"h":"1"}]`,
		},
		{
			name:         "Old Header string 1",
			inputJSON:    `[2,["2*",745],["2*",729],["2*",729],{"header":"1"}]`,
			expectedJSON: `[2,["2*",745],["2*",729],["2*",729],{"h":"1"}]`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var op FormatOp
			err := json.Unmarshal([]byte(tc.inputJSON), &op)
			require.NoError(t, err)

			outputJSON, err := json.Marshal(op)
			require.NoError(t, err)

			require.Equal(t, tc.expectedJSON, string(outputJSON))
		})
	}
}

func TestRewindMarshal(t *testing.T) {
	op := RewindOp{
		ID: ID{"1", 10},
		Address: ContentAddress{
			StartID: ID{"1", 0},
			EndID:   ID{"1", 9},
			MaxIDs:  map[string]int{"1": 10},
		},
	}

	bytes, err := json.Marshal(op)
	require.NoError(t, err)
	fmt.Println(string(bytes))

	var op2 RewindOp
	err = json.Unmarshal(bytes, &op2)
	require.NoError(t, err)
}

func TestMarshalUnmarshalBadOps(t *testing.T) {
	cases := []struct {
		name         string
		inputJSON    string
		expectedJSON string
	}{
		{
			name:         "only bad",
			inputJSON:    `{"version":"v0","ops":[[12],[5]]}`,
			expectedJSON: `{"version":"v0","ops":[]}`,
		},
		{
			// drop only the bad ops
			name:         "good and bad",
			inputJSON:    `{"version":"v0","ops":[[0,["0",0],"Hello World!",["",0],0],[12],[5]]}`,
			expectedJSON: `{"version":"v0","ops":[[0,["0",0],"Hello World!",["",0],0]]}`,
		},
		{
			// drop only the bad part of the MultiOp
			name:         "bad multi op",
			inputJSON:    `{"version":"v0","ops":[[6,["0",0],[[15],[0,["0",0],"Hello World!",["",0],0],[1,["0",12],["0",0],5]]]]}`,
			expectedJSON: `{"version":"v0","ops":[[6,["0",0],[[0,["0",0],"Hello World!",["",0],0],[1,["0",12],["0",0],5]]]]}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var r *Rogue
			err := json.Unmarshal([]byte(tc.inputJSON), &r)
			require.NoError(t, err)

			outputJSON, err := json.Marshal(r)
			require.NoError(t, err)

			require.Equal(t, tc.expectedJSON, string(outputJSON))
		})
	}
}

func TestMaxID(t *testing.T) {
	cases := []struct {
		name  string
		op    Op
		maxID ID
	}{
		{
			name: "DeleteOp",
			op: DeleteOp{
				ID:         ID{"0", 11},
				TargetID:   ID{"0", 5},
				SpanLength: 5,
			},
			maxID: ID{"0", 11},
		},
		{
			name: "FormatOp",
			op: FormatOp{
				ID:      ID{"0", 11},
				StartID: ID{"0", 5},
				EndID:   ID{"0", 10},
			},
			maxID: ID{"0", 11},
		},
		{
			name: "InsertOp",
			op: InsertOp{
				ID:   ID{"0", 0},
				Text: "Hello World!",
			},
			maxID: ID{"0", 11},
		},
		{
			name: "MultiOp",
			op: MultiOp{
				Mops: []Op{
					InsertOp{
						ID:   ID{"0", 0},
						Text: "Hello World!",
					},
					DeleteOp{
						ID:         ID{"0", 12},
						TargetID:   ID{"0", 5},
						SpanLength: 5,
					},
					InsertOp{
						ID:   ID{"0", 13},
						Text: "Hello World!",
					},
				},
			},
			maxID: ID{"0", 24},
		},
		{
			name: "MultiOp out of order",
			op: MultiOp{
				Mops: []Op{
					InsertOp{
						ID:   ID{"0", 13},
						Text: "Hello World!",
					},
					DeleteOp{
						ID:         ID{"0", 12},
						TargetID:   ID{"0", 5},
						SpanLength: 5,
					},
					InsertOp{
						ID:   ID{"0", 0},
						Text: "Hello World!",
					},
				},
			},
			maxID: ID{"0", 24},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.maxID, MaxID(tc.op))
		})
	}

}
