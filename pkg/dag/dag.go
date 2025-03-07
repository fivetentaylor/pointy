package dag

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/stackerr"
)

type Node interface {
	Run(ctx context.Context) (Node, error)
}

type Dag struct {
	Root Node

	Uuid       string
	ParentId   string // usually the uuid of the document
	Name       string
	OnError    func(ctx context.Context, node Node, err error)
	OnComplete func(ctx context.Context, dag *Dag)

	step  atomic.Uint64
	state *State
	mutex sync.Mutex
}

func New(name string, root Node) *Dag {
	return &Dag{
		Name:     name,
		ParentId: "unknown",
		Root:     root,
		OnError:  func(ctx context.Context, node Node, err error) {},
	}
}

func (d *Dag) Run(ctx context.Context, values map[string]any) error {
	log := env.SLog(ctx)
	started := time.Now()
	defer func() {
		log.Info("[dag] done", "uuid", d.Uuid, "dag", d.Name, "duration", time.Since(started).Seconds(), "event", "dag_complete")
		if r := recover(); r != nil {
			log.Error("[dag] panic", "PANIC", r, "stack", string(debug.Stack()), "event", "dag_panic")
			d.OnError(ctx, nil, stackerr.Errorf("%v", r))
		}
	}()

	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.Uuid = uuid.NewString()

	d.state = NewState(values)

	log.Info("[dag] starting", "uuid", d.Uuid, "dag", d.Name, "state", d.state, "event", "dag_start")

	ctx = WithDag(ctx, d)
	ctx = WithDagState(ctx, d.state)
	ctx = WithRunID(ctx, d.Uuid)

	currentNode := d.Root
	d.step.Store(1)
	var err error
	for currentNode != nil {
		currentNode, err = d.RunNode(ctx, currentNode)
		if err != nil {
			log.Error("[dag] error", "error", err, "event", "dag_error")
			go saveLogFile(ctx, "final_error.txt", err.Error())
			log.Error("[dag] RETURNING ERROR", "error", err, "event", "dag_error")
			return err
		}
	}

	ctx = ClearCurrentNode(ctx)

	defer func() {
		go d.logState(ctx, 99)
	}()

	if d.OnComplete != nil {
		d.OnComplete(ctx, d)
	}

	return nil
}

func (d *Dag) RunNode(ctx context.Context, node Node) (Node, error) {
	log := env.SLog(ctx)
	step := d.step.Load()
	nodeType := GetNodeType(node)
	log.Info("[dag] running node", "step", step, "node", nodeType, "dag", d.Name)

	ctx = WithCurrentNode(ctx, node)
	d.logState(ctx, step)

	nextNode, err := node.Run(ctx)
	log.Info("[dag] node done", "step", step, "node", nodeType, "dag", d.Name, "err", err)

	if err != nil {
		defer d.OnError(ctx, node, err)
		saveLogFile(ctx, "error.txt", err.Error())
		return nil, stackerr.Wrap(err)
	}

	d.step.Add(1)
	return nextNode, nil
}

// This should only be used in tests, if you need to inspect the state during a dag run
// use the ctx instead
func (d *Dag) State() *State {
	return d.state
}

func (d *Dag) logState(ctx context.Context, step uint64) {
	log := env.SLog(ctx)
	filtered := d.state.Filtered()
	log.Info("[dag] state", "step", step, "state", filtered)

	go saveLogFile(ctx, fmt.Sprintf("state_%02d.json", step), filtered)
}
