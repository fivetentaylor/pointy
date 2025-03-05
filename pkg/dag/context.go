package dag

import (
	"context"
)

type dagKey struct{}
type dagStateKey struct{}
type currentNodeKey struct{}
type runIdKey struct{}

func WithDag(ctx context.Context, d *Dag) context.Context {
	return context.WithValue(ctx, dagKey{}, d)
}

func WithCurrentNode(ctx context.Context, node Node) context.Context {
	return context.WithValue(ctx, currentNodeKey{}, node)
}

func ClearCurrentNode(ctx context.Context) context.Context {
	return context.WithValue(ctx, currentNodeKey{}, nil)
}

func WithDagState(ctx context.Context, state *State) context.Context {
	return context.WithValue(ctx, dagStateKey{}, state)
}

func GetDag(ctx context.Context) *Dag {
	return ctx.Value(dagKey{}).(*Dag)
}

func GetCurrentNode(ctx context.Context) Node {
	if node, ok := ctx.Value(currentNodeKey{}).(Node); ok {
		return node
	}
	return nil
}

func GetState(ctx context.Context) *State {
	return ctx.Value(dagStateKey{}).(*State)
}

func WithRunID(ctx context.Context, runID string) context.Context {
	return context.WithValue(ctx, runIdKey{}, runID)
}

func GetRunID(ctx context.Context) string {
	runID, ok := ctx.Value(runIdKey{}).(string)
	if !ok {
		return ""
	}
	return runID
}
