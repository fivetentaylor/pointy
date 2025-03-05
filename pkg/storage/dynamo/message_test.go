package dynamo

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/teamreviso/code/pkg/models"
)

func TestCreateMessage(t *testing.T) {
	db, err := NewDB()
	assert.NoError(t, err)

	docID := uuid.New().String()
	chanID := uuid.New().String()

	msg := &Message{
		DocID:          docID,
		ContainerID:    "chan#" + chanID,
		ChannelID:      chanID,
		MessageID:      uuid.New().String(),
		CreatedAt:      time.Now().Unix(),
		Content:        "hello world",
		UserID:         uuid.New().String(),
		AuthorID:       "1",
		LifecycleStage: MessageLifecycleStageCompleted,
		Attachments: &models.AttachmentList{
			Attachments: []*models.Attachment{},
		},
	}

	err = db.CreateMessage(msg)
	assert.NoError(t, err)

	rm, err := db.GetMessage(msg.ContainerID, msg.MessageID)
	assert.NoError(t, err)

	assert.Equal(t, msg.MessageID, rm.MessageID)
}

func TestCreateMessageReply(t *testing.T) {
	db, err := NewDB()
	assert.NoError(t, err)

	docID := uuid.New().String()
	chanID := uuid.New().String()

	msg := &Message{
		DocID:          docID,
		ContainerID:    "chan#" + chanID,
		ChannelID:      chanID,
		MessageID:      uuid.New().String(),
		CreatedAt:      time.Now().Unix(),
		Content:        "first message",
		UserID:         uuid.New().String(),
		AuthorID:       "1",
		LifecycleStage: MessageLifecycleStageCompleted,
		Attachments: &models.AttachmentList{
			Attachments: []*models.Attachment{},
		},
	}
	err = db.CreateMessage(msg)
	assert.NoError(t, err)

	firstMsg, err := db.GetMessage(msg.ContainerID, msg.MessageID)
	assert.NoError(t, err)

	replyMsg := &Message{
		DocID:             docID,
		ParentContainerID: &firstMsg.ContainerID,
		ContainerID:       firstMsg.GetContainerID(),
		ChannelID:         chanID,
		MessageID:         uuid.New().String(),
		CreatedAt:         time.Now().Unix(),
		Content:           "",
		UserID:            uuid.New().String(),
		AuthorID:          "!1",
		LifecycleStage:    MessageLifecycleStagePending,
		Attachments: &models.AttachmentList{
			Attachments: []*models.Attachment{},
		},
	}
	err = db.CreateMessage(replyMsg)
	assert.NoError(t, err)

	firstMsg, err = db.GetMessage(msg.ContainerID, msg.MessageID)
	assert.NoError(t, err)

	reply2Msg := &Message{
		DocID:             docID,
		ParentContainerID: &firstMsg.ContainerID,
		ContainerID:       firstMsg.GetContainerID(),
		ChannelID:         chanID,
		MessageID:         uuid.New().String(),
		CreatedAt:         time.Now().Unix(),
		Content:           "thanks!",
		UserID:            firstMsg.UserID,
		AuthorID:          "!1",
		LifecycleStage:    MessageLifecycleStageCompleted,
		Attachments: &models.AttachmentList{
			Attachments: []*models.Attachment{},
		},
	}
	err = db.CreateMessage(reply2Msg)
	assert.NoError(t, err)

	firstMsg, err = db.GetMessage(msg.ContainerID, msg.MessageID)
	assert.NoError(t, err)
}

func TestUpdatedMessage(t *testing.T) {
	db, err := NewDB()
	assert.NoError(t, err)
	docID := uuid.New().String()
	chanID := uuid.New().String()

	msg := &Message{
		DocID:          docID,
		ContainerID:    "chan#" + chanID,
		ChannelID:      chanID,
		MessageID:      uuid.New().String(),
		CreatedAt:      time.Now().Unix(),
		Content:        "hello world",
		AuthorID:       "1",
		UserID:         uuid.New().String(),
		LifecycleStage: MessageLifecycleStageCompleted,
		Attachments: &models.AttachmentList{
			Attachments: []*models.Attachment{},
		},
	}

	err = db.CreateMessage(msg)
	assert.NoError(t, err)

	msg.Content = "new content"
	err = db.UpdateMessage(msg)
	assert.NoError(t, err)

	msg, err = db.GetMessage(msg.ContainerID, msg.MessageID)
	assert.NoError(t, err)
	assert.Equal(t, "new content", msg.Content)

	msg.Content = "new content 2"
}
