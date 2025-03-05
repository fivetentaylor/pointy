package dynamo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateContentAddress(t *testing.T) {
	db, err := NewDB()
	require.NoError(t, err)

	addressID, err := db.CreateContentAddress("docID", []byte("Hello Address!"))
	require.NoError(t, err)

	fmt.Printf("AddressID: %s\n", addressID)

	bytes, err := db.GetContentAddress("docID", addressID)
	require.NoError(t, err)

	addressIDs, _, err := db.GetContentAddressIDs("docID", nil)
	require.NoError(t, err)
	require.Equal(t, []string{addressID}, addressIDs)

	require.Equal(t, []byte("Hello Address!"), bytes)

}
