package dynamo_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/fivetentaylor/pointy/pkg/storage/dynamo"
)

func TestSentReceipt_FindOrCreateSentReceipt(t *testing.T) {
	db, err := dynamo.NewDB()
	require.NoError(t, err)

	userId := uuid.New().String()

	sr, err := db.CreateSentReceipt(userId, "messageID", "containerID")
	require.NoError(t, err)

	srOut, err := db.GetSentReceipt(userId, "messageID")
	require.NoError(t, err)

	assert.Equal(t, sr, srOut)
}

func TestSentReceipt_GetUnsentSentReceiptsForUser(t *testing.T) {
	db, err := dynamo.NewDB()
	require.NoError(t, err)
	userId := uuid.New().String()

	for i := 0; i < 5; i++ {
		msgId := uuid.New().String()
		containerId := uuid.New().String()
		sr, err := db.CreateSentReceipt(userId, msgId, containerId)
		require.NoError(t, err)
		fmt.Println(sr)
		if i%2 == 0 {
			sr, err = db.UpdateSentReceiptSent(userId, msgId, true)
			require.NoError(t, err)
			fmt.Println("UPDATED", sr)
		}
	}

	unsents, err := db.GetUnsentSentReceiptsForUser(userId)
	require.NoError(t, err)
	require.Equal(t, 2, len(unsents))
}

func TestSentReceipt_LastSendReceipt(t *testing.T) {
	db, err := dynamo.NewDB()
	require.NoError(t, err)
	userId := uuid.New().String()

	for i := 0; i < 5; i++ {
		msgId := uuid.New().String()
		containerId := uuid.New().String()
		sr, err := db.CreateSentReceipt(userId, msgId, containerId)
		require.NoError(t, err)
		fmt.Println(sr)
		if i%2 == 0 {
			sr, err = db.UpdateSentReceiptSent(userId, msgId, true)
			require.NoError(t, err)
			fmt.Println("UPDATED", sr)
		}
	}

	sr, err := db.GetLastSentReceipt(userId)
	require.NoError(t, err)

	assert.NotEqual(t, 0, sr.SentAt)
}
