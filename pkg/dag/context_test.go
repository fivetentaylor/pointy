package dag_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/teamreviso/code/pkg/dag"
)

func Test_WithRunId(t *testing.T) {
	ctx := dag.WithRunID(context.Background(), "test")
	runId := dag.GetRunID(ctx)
	require.Equal(t, "test", runId)
}

func Test_WithDagState(t *testing.T) {
	ctx := dag.WithDagState(context.Background(), &dag.State{})
	state := dag.GetState(ctx)
	require.NotNil(t, state)
}
