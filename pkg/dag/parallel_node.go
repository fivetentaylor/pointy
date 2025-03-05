package dag

import (
	"context"
	"sync"

	"github.com/teamreviso/code/pkg/env"
)

type ParallelNode struct {
	Nodes []Node
	After Node
}

func (n *ParallelNode) Run(ctx context.Context) (Node, error) {
	log := env.Log(ctx)

	waitGroup := &sync.WaitGroup{}
	dag := GetDag(ctx)

	for i, node := range n.Nodes {
		log.Info("running parallel node", "i", i, "node", node)
		n.RunNode(ctx, dag, waitGroup, node)
	}

	waitGroup.Wait()

	return n.After, nil
}

func (n *ParallelNode) RunNode(ctx context.Context, dag *Dag, waitGroup *sync.WaitGroup, node Node) {
	waitGroup.Add(1)
	go func(waitGroup *sync.WaitGroup, node Node) {
		defer waitGroup.Done()

		ctx = WithCurrentNode(ctx, node)
		next, err := dag.RunNode(ctx, node)
		if err != nil {
			env.Log(ctx).Error("[parallel] error running parallel node", "error", err)
			saveLogFile(ctx, "error.txt", err.Error())
		}

		if next != nil {
			if n.After != nil {
				env.Log(ctx).Info("❗️❗️❗️❗️ parallel node has multiple after nodes", "node", node, "after", n.After)
			} else {
				n.After = next
			}
		}
	}(waitGroup, node)
}
