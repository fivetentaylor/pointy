package dag_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/teamreviso/code/pkg/dag"
	"github.com/teamreviso/code/pkg/testutils"
)

type testNode struct {
	called bool
	error  error
}

func (n *testNode) Run(ctx context.Context) (dag.Node, error) {
	n.called = true
	return nil, n.error
}

func TestDag_Run(t *testing.T) {
	testCtx := testutils.TestContext()
	testutils.EnsureStorage()

	t.Run("Run will run all nodes in the dag", func(t *testing.T) {
		n := &testNode{called: false}
		d := dag.New("test", n)
		err := d.Run(testCtx, nil)
		require.NoError(t, err)
		require.True(t, n.called)
	})

	t.Run("Run will call OnError if there is an error", func(t *testing.T) {
		onErrorCaleld := false
		n := &testNode{called: false, error: fmt.Errorf("error")}
		d := dag.New("test", n)
		d.OnError = func(ctx context.Context, node dag.Node, err error) {
			require.Equal(t, n, node)
			require.Equal(t, fmt.Errorf("error"), err)
			onErrorCaleld = true
		}
		err := d.Run(testCtx, nil)
		require.Error(t, err)
		require.True(t, n.called)
		require.True(t, onErrorCaleld)
	})

	t.Run("Run will not call OnError if there is no error", func(t *testing.T) {
		onErrorCaleld := false
		n := &testNode{called: false, error: nil}
		d := dag.New("test", n)
		d.OnError = func(ctx context.Context, node dag.Node, err error) {
			require.Equal(t, n, node)
			require.Equal(t, fmt.Errorf("error"), err)
			onErrorCaleld = true
		}
		err := d.Run(testCtx, nil)
		require.NoError(t, err)
		require.True(t, n.called)
		require.False(t, onErrorCaleld)
	})
}
