package dynamo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashAndEncode(t *testing.T) {
	id0 := hashAndEncode([]byte("Hello World!"), 16)
	fmt.Printf("ID: %s\n", id0)

	id1 := hashAndEncode([]byte("Goodbye World!"), 16)
	fmt.Printf("ID: %s\n", id1)

	require.NotEqual(t, id0, id1)
}
