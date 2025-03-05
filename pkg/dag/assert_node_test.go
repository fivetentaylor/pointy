package dag

import "context"

type AssertNode struct {
	Ran bool
}

func (n *AssertNode) Run(ctx context.Context) (Node, error) {
	n.Ran = true
	return nil, nil
}
