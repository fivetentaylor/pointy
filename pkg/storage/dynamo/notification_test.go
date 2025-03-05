package dynamo_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/teamreviso/code/pkg/models"
	"github.com/teamreviso/code/pkg/storage/dynamo"
)

func TestNotification_Create(t *testing.T) {
	db, err := dynamo.NewDB()
	require.NoError(t, err)

	userId := uuid.New().String()

	n := dynamo.Notification{
		ID:      uuid.New().String(),
		UserID:  userId,
		DocID:   uuid.New().String(),
		Payload: &models.NotificationPayload{},
	}

	err = db.UpsertNotification(&n)
	require.NoError(t, err)

	notifs, _, err := db.GetNotificationsForUser(userId, false, dynamo.PaginationParams{
		Limit: 10,
	})
	require.NoError(t, err)
	require.Equal(t, 1, len(notifs))

	notif := notifs[0]
	assert.Equal(t, n.UserID, notif.UserID)
	assert.Equal(t, n.Payload, notif.Payload)
	assert.NotEmpty(t, notif.ID, "notification id should not be empty")
	assert.NotEmpty(t, notif.CreatedAt, "created at should not be empty")
}

func TestNotification_GetNotificationsForUser(t *testing.T) {
	db, err := dynamo.NewDB()
	require.NoError(t, err)

	userId := uuid.New().String()
	for i := 0; i < 10; i++ {
		n := dynamo.Notification{
			ID:      uuid.New().String(),
			UserID:  userId,
			DocID:   uuid.New().String(),
			Payload: &models.NotificationPayload{},
		}
		err = db.UpsertNotification(&n)
		require.NoError(t, err)
	}

	notifs, next, err := db.GetNotificationsForUser(userId, false, dynamo.PaginationParams{
		Limit: 3,
	})
	require.NoError(t, err)

	assert.Equal(t, 3, len(notifs))
	lastNotif := notifs[len(notifs)-1]

	notifs, next, err = db.GetNotificationsForUser(userId, false, dynamo.PaginationParams{
		Limit:             3,
		ExclusiveStartKey: next,
	})
	require.NoError(t, err)

	assert.Equal(t, 3, len(notifs))
	firstNotif := notifs[0]
	assert.Greater(t, lastNotif.CreatedAt, firstNotif.CreatedAt)
}

func TestNotification_GetUnreadNotificationsForUser(t *testing.T) {
	db, err := dynamo.NewDB()
	require.NoError(t, err)

	userId := uuid.New().String()
	for i := 0; i < 10; i++ {
		n := dynamo.Notification{
			ID:      uuid.New().String(),
			UserID:  userId,
			DocID:   uuid.New().String(),
			Payload: &models.NotificationPayload{},
		}
		err = db.UpsertNotification(&n)
		require.NoError(t, err)
	}

	count, err := db.GetNotificationCountForUser(userId, false)
	require.NoError(t, err)
	assert.Equal(t, int64(10), count)
}

func TestNotification_ReadNotification(t *testing.T) {
	db, err := dynamo.NewDB()
	require.NoError(t, err)

	userID := uuid.New().String()
	docID := uuid.New().String()
	n := &dynamo.Notification{
		UserID: userID,
		DocID:  docID,
		ID:     uuid.New().String(),
	}
	err = db.UpsertNotification(n)
	require.NoError(t, err)

	count, err := db.GetNotificationCountForUser(userID, false)
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)

	err = db.MarkNotifications(userID, docID, []string{n.ID}, true)
	require.NoError(t, err)

	count, err = db.GetNotificationCountForUser(userID, false)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)

	// duplicate read
	err = db.MarkNotifications(userID, docID, []string{n.ID}, true)
	require.NoError(t, err)

	notf, err := db.GetNotification(userID, n.ID)
	require.NoError(t, err)

	assert.Equal(t, notf.Read, true)
	assert.Equal(t, notf.UserID, userID)
	assert.NotEmpty(t, notf.CreatedAt)
}
