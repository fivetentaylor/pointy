package voice

import "github.com/fivetentaylor/pointy/pkg/dag"

func UpdateDocumentDag() *dag.Dag {
	postNode := &dag.TitleThreadNode{}

	d := dag.New("voice_update", &dag.DirectSelectEditTargetNode{
		Next: &dag.DirectReviseNode{
			Next: postNode,
		},
	})

	return d
}
