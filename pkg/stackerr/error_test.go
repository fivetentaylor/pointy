package stackerr_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/fivetentaylor/pointy/pkg/stackerr"
)

func TestNew(t *testing.T) {
	t.Run("raw error", func(t *testing.T) {
		err := stackerr.New(fmt.Errorf("some error"))
		assert.Contains(t, err.Error(), stackerr.Header)
	})

	t.Run("with wrapped error", func(t *testing.T) {
		err := stackerr.New(fmt.Errorf("some error"))
		err1 := stackerr.New(err)

		assert.Equal(t, err.Error(), err1.Error())
	})

	t.Run("with nil", func(t *testing.T) {
		err := stackerr.New(nil)
		assert.Nil(t, err)
	})
}
