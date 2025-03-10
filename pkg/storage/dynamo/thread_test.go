package dynamo

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/fivetentaylor/pointy/pkg/models"
)

func TestThreadCreate(t *testing.T) {
	db, err := NewDB()
	assert.NoError(t, err)

	thread := &Thread{
		DocID:     uuid.New().String(),
		ThreadID:  uuid.New().String(),
		UserID:    uuid.New().String(),
		Title:     "whatup",
		UpdatedAt: time.Now().Unix(),
	}

	err = db.CreateThread(thread)
	assert.NoError(t, err)

	rt, err := db.GetThreadForUser(thread.DocID, thread.ThreadID, thread.UserID)
	assert.NoError(t, err)
	assert.Equal(t, thread, rt)
}

func TestGetThreadsForDoc(t *testing.T) {
	db, err := NewDB()
	assert.NoError(t, err)

	docID := uuid.New().String()
	userID1 := uuid.New().String()
	userID2 := uuid.New().String()

	thread := &Thread{
		DocID:     docID,
		UserID:    userID1,
		Title:     "New Chat user 1",
		UpdatedAt: time.Now().Unix(),
	}
	err = db.CreateThread(thread)
	require.NoError(t, err)

	thread2 := &Thread{
		DocID:     docID,
		UserID:    userID2,
		Title:     "New Chat user 2",
		UpdatedAt: time.Now().Unix(),
	}
	err = db.CreateThread(thread2)
	require.NoError(t, err)

	threads, err := db.GetThreadsForDocForUser(docID, userID1)
	assert.NoError(t, err)
	require.Equal(t, 1, len(threads))
	assert.Equal(t, thread.ThreadID, threads[0].ThreadID)
	assert.Equal(t, thread.Title, threads[0].Title)
}

func TestMessagesForThread(t *testing.T) {
	db, err := NewDB()
	assert.NoError(t, err)

	docID := uuid.New().String()
	userID1 := uuid.New().String()

	thread := &Thread{
		DocID:     docID,
		UserID:    userID1,
		Title:     "New Chat user 1",
		UpdatedAt: time.Now().Unix(),
	}
	err = db.CreateThread(thread)
	require.NoError(t, err)

	msg := &Message{
		DocID:          docID,
		ContainerID:    AiThreadPrefix + thread.ThreadID,
		ChannelID:      uuid.New().String(),
		MessageID:      uuid.New().String(),
		CreatedAt:      time.Now().Unix(),
		Content:        "hello world",
		UserID:         userID1,
		AuthorID:       "1",
		LifecycleStage: MessageLifecycleStageCompleted,
		Attachments: &models.AttachmentList{
			Attachments: []*models.Attachment{},
		},
	}
	err = db.CreateMessage(msg)
	require.NoError(t, err)

	msgs, err := db.GetMessagesForThread(thread.ThreadID)
	assert.NoError(t, err)
	require.Equal(t, 1, len(msgs))
	assert.Equal(t, msg.MessageID, msgs[0].MessageID)
}
