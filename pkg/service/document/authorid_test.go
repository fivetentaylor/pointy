package document_test

import (
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/fivetentaylor/pointy/pkg/service/document"
	"github.com/fivetentaylor/pointy/pkg/testutils"
)

func TestNewAuthorID(t *testing.T) {
	t.Parallel()
	testutils.EnsureStorage()

	ctx := testutils.TestContext()

	docID := uuid.New().String()
	userID := uuid.New().String()

	authorID, err := document.NewAuthorID(ctx, docID, userID)
	require.NoError(t, err)
	assert.Equal(t, "1", authorID)

	authorID2, err := document.NewAuthorID(ctx, docID, userID)
	require.NoError(t, err)
	assert.Equal(t, "2", authorID2)

	for i := 0; i < 98; i++ {
		_, err := document.NewAuthorID(ctx, docID, userID)
		require.NoError(t, err)
	}

	authorID3, err := document.NewAuthorID(ctx, docID, userID)
	require.NoError(t, err)

	newAuthorID, err := strconv.ParseInt(authorID3, 16, 64)

	assert.Equal(t, int64(101), newAuthorID)
}

func TestValidateAuthorID(t *testing.T) {
	t.Parallel()
	testutils.EnsureStorage()

	ctx := testutils.TestContext()
	docID := uuid.New().String()
	userID := uuid.New().String()

	authorID, err := document.NewAuthorID(ctx, docID, userID)
	require.NoError(t, err)
	require.NotEmpty(t, authorID)

	valid, err := document.ValidateAuthorID(ctx, authorID, docID, userID)
	require.NoError(t, err)
	require.True(t, valid)

	// should not be valid for another user
	otherUserID := uuid.New().String()
	valid, err = document.ValidateAuthorID(ctx, authorID, docID, otherUserID)
	require.NoError(t, err)
	require.False(t, valid)

	// should not be valid for the same user but different document
	otherDocID := uuid.New().String()
	valid, err = document.ValidateAuthorID(ctx, authorID, otherDocID, userID)
	require.NoError(t, err)
	require.False(t, valid)

	// invalid id should error
	badAuthorID := "root"
	valid, err = document.ValidateAuthorID(ctx, badAuthorID, docID, userID)
	require.Error(t, err)

	badAuthorID = "q"
	valid, err = document.ValidateAuthorID(ctx, badAuthorID, docID, userID)
	require.Error(t, err)
}
