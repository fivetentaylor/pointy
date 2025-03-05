package v3_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	v3 "github.com/teamreviso/code/rogue/v3"
	"github.com/teamreviso/code/rogue/v3/testcases"
)

func TestToOps(t *testing.T) {
	type toOpsTC struct {
		name     string
		fn       func(*v3.Rogue)
		expected []v3.Op
	}

	testCases := []toOpsTC{
		{
			name: "empty",
			fn:   func(*v3.Rogue) {},
			expected: []v3.Op{
				v3.InsertOp{ID: v3.RootID, Text: "x", Side: 0, ParentID: v3.ID{"", 0}},
				v3.InsertOp{ID: v3.ID{"q", 1}, Text: "\n", ParentID: v3.RootID, Side: v3.Right},
				v3.DeleteOp{ID: v3.ID{"q", 2}, TargetID: v3.RootID, SpanLength: 1},
			},
		},
		{
			name: "simple insert",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "hello world")
				require.NoError(t, err)
			},
			expected: []v3.Op{
				v3.InsertOp{ID: v3.RootID, Text: "x", Side: 0, ParentID: v3.ID{"", 0}},
				v3.InsertOp{ID: v3.ID{"q", 1}, Text: "\n", ParentID: v3.RootID, Side: v3.Right},
				v3.DeleteOp{ID: v3.ID{"q", 2}, TargetID: v3.RootID, SpanLength: 1},
				v3.InsertOp{ID: v3.ID{"auth0", 3}, Text: "hello world", ParentID: v3.ID{"q", 1}, Side: -1},
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
			expected: []v3.Op{
				v3.InsertOp{ID: v3.RootID, Text: "x", Side: 0, ParentID: v3.ID{"", 0}},
				v3.InsertOp{ID: v3.ID{"q", 1}, Text: "\n", ParentID: v3.RootID, Side: v3.Right},
				v3.DeleteOp{ID: v3.ID{"q", 2}, TargetID: v3.RootID, SpanLength: 1},
				v3.InsertOp{ID: v3.ID{"auth0", 3}, Text: "hello world", ParentID: v3.ID{"q", 1}, Side: -1},
				v3.FormatOp{
					ID:      v3.ID{"auth0", 14},
					StartID: v3.ID{"auth0", 9},
					EndID:   v3.ID{"q", 1},
					Format:  v3.FormatV3Span{"b": "true", "e": "true"},
				},
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
			expected: []v3.Op{
				v3.InsertOp{ID: v3.ID{Author: "root", Seq: 0}, Text: "x", ParentID: v3.ID{Author: "", Seq: 0}, Side: 0},
				v3.InsertOp{ID: v3.ID{Author: "q", Seq: 1}, Text: "\n", ParentID: v3.ID{Author: "root", Seq: 0}, Side: 1},
				v3.DeleteOp{ID: v3.ID{Author: "q", Seq: 2}, TargetID: v3.ID{Author: "root", Seq: 0}, SpanLength: 1},
				v3.InsertOp{ID: v3.ID{Author: "auth0", Seq: 3}, Text: "🥔🥒🥔", ParentID: v3.ID{Author: "q", Seq: 1}, Side: -1},
				v3.InsertOp{ID: v3.ID{Author: "auth0", Seq: 9}, Text: "🤡😿😝", ParentID: v3.ID{Author: "auth0", Seq: 5}, Side: -1},
			},
		},
	}

	realifeTestCases := testcases.LoadAll(t)
	for _, tc := range realifeTestCases {
		itc := tc
		testCases = append(testCases, toOpsTC{
			name: itc.Name,
			fn:   func(r *v3.Rogue) { r.Copy(itc.Doc) },
		})
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := v3.NewRogueForQuill("auth0")
			tc.fn(r)

			/*for i, root := range r.Roots {
				fmt.Printf("root[%d]:\n", i)
				root.Inspect()
				fmt.Println()
			}*/

			ops, err := r.ToOps()
			require.NoError(t, err)

			if tc.expected != nil {
				require.Equal(t, tc.expected, ops, "expected: %#v, got: %#v", tc.expected, ops)
			}

			r2 := v3.NewRogueForQuill("auth0")
			for _, op := range ops {
				switch fop := op.(type) {
				case v3.FormatOp:
					if _, ok := fop.Format.(v3.FormatV3Span); ok {
						_, startIx, err := r.Rope.GetIndex(fop.StartID)
						require.NoError(t, err)
						_, endIx, err := r.Rope.GetIndex(fop.EndID)
						require.NoError(t, err)

						fmt.Printf("id: %v, startIx: %d, endIx: %d, op: %v\n", op.GetID(), startIx, endIx, fop.Format)
					}
				case v3.InsertOp:
					fmt.Printf("id: %v, text: %q\n", op.GetID(), fop.Text)
				}
				// fmt.Printf("ser op[%d]: %#v\n", i, op)
				_, err := r2.MergeOp(op)
				require.NoError(t, err)
			}

			ops, err = r2.ToOps()
			require.NoError(t, err)

			r3 := v3.NewRogueForQuill("auth0")
			for _, op := range ops {
				_, err := r3.MergeOp(op)
				require.NoError(t, err)
			}

			html, err := r.GetHtml(v3.RootID, v3.LastID, true, false)
			require.NoError(t, err)
			html2, err := r2.GetHtml(v3.RootID, v3.LastID, true, false)
			require.NoError(t, err)
			html3, err := r3.GetHtml(v3.RootID, v3.LastID, true, false)
			require.NoError(t, err)

			if html != html2 {
				fmt.Println("html != html2 probably because of old deser discrepancy")
				require.Equal(t, html2, html3)
				fmt.Printf("html : %q\n", html)
				fmt.Printf("html2: %q\n", html2)
			}
		})
	}
}

func TestSerializable_AsJS(t *testing.T) {
	type tcs struct {
		name string
		fn   func(*v3.Rogue)
	}

	testCases := []tcs{
		{
			name: "empty",
			fn:   func(*v3.Rogue) {},
		},
		{
			name: "simple insert",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "hello world")
				require.NoError(t, err)
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
		},
		{
			name: "overlapping formats",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "hello world")
				require.NoError(t, err)
				_, err = r.Format(6, 5, v3.FormatV3Span{"b": "true"})
				require.NoError(t, err)
				_, err = r.Format(6, 5, v3.FormatV3Span{"b": "null"})
				require.NoError(t, err)
				_, err = r.Format(6, 5, v3.FormatV3Span{"i": "true"})
				require.NoError(t, err)
			},
		},
	}

	realifeTestCases := testcases.LoadAll(t)
	for _, tc := range realifeTestCases {
		itc := tc
		testCases = append(testCases, tcs{
			name: tc.Name,
			fn:   func(r *v3.Rogue) { r.Copy(itc.Doc) },
		})
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := v3.NewRogueForQuill("auth0")
			tc.fn(r)

			serializable, err := r.Serializable()
			require.NoError(t, err)

			require.True(t, len(serializable.Ops) > 0)
			require.Equal(t, "v0", *serializable.Version)
		})
	}
}

func TestSerializable(t *testing.T) {
	type toOpsTC struct {
		name string
		fn   func(*v3.Rogue)
	}

	testCases := []toOpsTC{
		{
			name: "empty",
			fn:   func(*v3.Rogue) {},
		},
		{
			name: "simple insert",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "hello world")
				require.NoError(t, err)
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
		},
		{
			name: "insert and deltes with emojis then insert",
			fn: func(r *v3.Rogue) {
				_, err := r.Insert(0, "🥸 👩🏽‍🚀😗")
				require.NoError(t, err)
				_, err = r.Delete(2, 1)
				require.NoError(t, err)
				_, err = r.Insert(2, "🤡😿😝")
				require.NoError(t, err)
			},
		},
	}

	realifeTestCases := testcases.LoadAll(t)
	for _, tc := range realifeTestCases {
		itc := tc
		testCases = append(testCases, toOpsTC{
			name: itc.Name,
			fn:   func(r *v3.Rogue) { r.Copy(itc.Doc) },
		})
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := v3.NewRogueForQuill("auth0")
			tc.fn(r)

			html, err := r.GetHtml(v3.RootID, v3.LastID, true, false)
			require.NoError(t, err)

			fmt.Printf("HTML before %q\n", html)

			snapshot, err := r.NewSnapshotOp()
			require.NoError(t, err)

			var r2 v3.Rogue
			_, err = r2.MergeOp(snapshot)
			require.NoError(t, err)

			html2, err := r2.GetHtml(v3.RootID, v3.LastID, true, false)
			require.NoError(t, err)

			var r3 v3.Rogue
			snapshot2, err := r2.NewSnapshotOp()
			require.NoError(t, err)
			_, err = r3.MergeOp(snapshot2)
			require.NoError(t, err)

			html3, err := r3.GetHtml(v3.RootID, v3.LastID, true, false)
			require.NoError(t, err)

			// fmt.Println("HTML after", html2)

			if html != html2 {
				fmt.Println("html != html2 probably because of old deser discrepancy")
				require.Equal(t, html2, html3)
			}
		})
	}
}

func TestFormatV5Migration(t *testing.T) {
	snapshot := `{"ops":[[0,["root",0],"x",["",0],0],[0,["q",1],"\n",["root",0],1],[1,["q",2],["root",0],1],[6,["1",3],[[0,["1",3],"Untitled",["q",1],-1],[2,["1",11],["1",3],["q",1],{"e":"true"}]]],[2,["1",12],["q",1],["q",1],{"h":"1"}],[6,["1",13],[[1,["1",13],["1",3],8]]],[2,["1",14],["q",1],["q",1],{}],[6,["1",15],[[0,["1",15],"H",["1",10],1],[2,["1",16],["1",15],["q",1],{"e":"true"}]]],[0,["1",17],"e",["1",15],1],[0,["1",18],"l",["1",17],1],[0,["1",19],"l",["1",18],1],[0,["1",20],"o",["1",19],1],[0,["1",21]," ",["1",20],1],[0,["1",22],"W",["1",21],1],[0,["1",23],"o",["1",22],1],[0,["1",24],"r",["1",23],1],[0,["1",25],"l",["1",24],1],[0,["1",26],"d",["1",25],1],[0,["1",27],"!",["1",26],1],[2,["1",28],["1",22],["1",27],{"u":"true"}],[2,["1",29],["1",15],["1",21],{"b":"true"}],[2,["1",30],["1",18],["1",25],{"s":"true"}],[2,["1",31],["1",15],["q",1],{"i":"true"}]]}`

	r := v3.NewRogueForQuill("0")
	err := json.Unmarshal([]byte(snapshot), &r)
	require.NoError(t, err)

	expectedHtml := `<p data-rid="q_1"><strong><em data-rid="1_15">He</em></strong><strong><em><s data-rid="1_18">llo</s></em></strong><em><s data-rid="1_21"> </s></em><em><s><u data-rid="1_22">Wor</u></s></em><em><u data-rid="1_25">ld</u></em><em data-rid="1_27">!</em></p>`
	html, err := r.GetHtml(v3.RootID, v3.LastID, true, true)
	require.NoError(t, err)

	require.Equal(t, expectedHtml, html)
}
