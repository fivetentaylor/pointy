package dynamo_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/fivetentaylor/pointy/pkg/storage/dynamo"
)

func TestReadReceipt_SK(t *testing.T) {
	unread := &dynamo.ReadReceipt{
		DocID:     "docID",
		ChannelID: "channelID",
		ContainerID: fmt.Sprintf(
			"%s%s", dynamo.ChannelPrefix, "channelID",
		),
		MessageID: "messageID",
		Read:      true,
		CreatedAt: 12345678,
	}

	assert.Equal(t, "rr#true#docID#channelID#chan#channelID", unread.SK1())
}

func TestReadReceipt_HydrateFromSK(t *testing.T) {
	unread := &dynamo.ReadReceipt{}
	err := unread.HydrateFromSK1("unrd#true#docID#channelID#chan#channelID")
	require.NoError(t, err)

	assert.Equal(t, "docID", unread.DocID)
	assert.Equal(t, "channelID", unread.ChannelID)
	assert.Equal(t, "chan#channelID", unread.ContainerID)
	assert.Equal(t, true, unread.Read)
}

func TestReadReceipt_FindOrCreateReadReceipt(t *testing.T) {
	db, err := dynamo.NewDB()
	require.NoError(t, err)

	userId := uuid.New().String()
	docId := uuid.New().String()

	rr, err := db.FindOrCreateReadReceipt(userId, docId, "channelID", "chan#channelID", "messageID", false)
	require.NoError(t, err)

	assert.NotEmpty(t, rr.CreatedAt)

	rrOut, err := db.GetReadReceipt(userId, "messageID")
	require.NoError(t, err)

	assert.Equal(t, rr, rrOut)
}

func TestReadReceipt_UserReadingMessageBeforeReadReceiptIsCreated(t *testing.T) {
	db, err := dynamo.NewDB()
	require.NoError(t, err)
	userId := uuid.New().String()

	for i := 0; i < 2; i++ {
		docId := uuid.New().String()
		for j := 0; j < 4; j++ {
			channelId := uuid.New().String()
			containerId := fmt.Sprintf("%s%s", dynamo.ChannelPrefix, channelId)

			for k := 0; k < 5; k++ {
				messageId := uuid.New().String()
				_, err := db.FindOrCreateReadReceipt(userId, docId, channelId, containerId, messageId, false)
				require.NoError(t, err)

				unreads, err := db.GetUnreadCountForContainer(userId, docId, channelId, containerId)
				require.NoError(t, err)
				require.Equal(t, int64(1+k), unreads)
			}

			unreads, err := db.GetUnreadCountForChannel(userId, docId, channelId)
			require.NoError(t, err)
			require.Equal(t, int64(5), unreads)
		}

		unreads, err := db.GetUnreadCountForDocument(userId, docId)
		require.NoError(t, err)
		require.Equal(t, int64(20), unreads)
	}

	unreads, err := db.GetUnreadCountForUser(userId)
	require.NoError(t, err)
	require.Equal(t, int64(40), unreads)
}

func TestReadReceipt_MarkAsRead(t *testing.T) {
	db, err := dynamo.NewDB()
	require.NoError(t, err)

	userId := uuid.New().String()
	docId := uuid.New().String()
	channelId := uuid.New().String()
	containerId := fmt.Sprintf("%s%s", dynamo.ChannelPrefix, channelId)
	messageId := uuid.New().String()

	_, err = db.MarkReadReceiptRead(userId, docId, channelId, containerId, messageId)
	require.NoError(t, err)

	_, err = db.MarkReadReceiptRead(userId, docId, channelId, containerId, messageId)
	require.NoError(t, err)

	unreads, err := db.GetUnreadCountForContainer(userId, docId, channelId, containerId)
	require.NoError(t, err)
	require.Equal(t, int64(0), unreads)
}

func TestReadReceipt_MarkAsUnread(t *testing.T) {
	db, err := dynamo.NewDB()
	require.NoError(t, err)

	userId := uuid.New().String()
	docId := uuid.New().String()
	channelId := uuid.New().String()
	containerId := fmt.Sprintf("%s%s", dynamo.ChannelPrefix, channelId)
	messageId := uuid.New().String()

	_, err = db.MarkReadReceiptUnread(userId, docId, channelId, containerId, messageId, false)
	require.NoError(t, err)

	_, err = db.MarkReadReceiptUnread(userId, docId, channelId, containerId, messageId, false)
	require.NoError(t, err)

	unreads, err := db.GetUnreadCountForContainer(userId, docId, channelId, containerId)
	require.NoError(t, err)
	require.Equal(t, int64(1), unreads)
}
