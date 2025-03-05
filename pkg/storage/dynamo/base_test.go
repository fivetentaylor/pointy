package dynamo

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBase(t *testing.T) {
	dydb, err := NewDB()
	if err != nil {
		t.Fatal(err)
	}

	docId := uuid.New().String()
	channel := &Channel{
		DocID: docId,
		Type:  ChannelTypeDirect,
		UserIDs: []string{
			uuid.New().String(),
		},
	}

	// Create channel for messages
	err = dydb.CreateChannel(channel)
	assert.NoError(t, err)

	// Check to make sure channel was created
	docChannels, err := dydb.GetDocumentChannels(docId)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(docChannels))
	assert.Equal(t, channel, docChannels[0])

	// Create first message in the channel
	msg := &Message{
		DocID:            docId,
		ContainerID:      fmt.Sprintf("%s%s", ChannelPrefix, channel.ChannelID),
		ChannelID:        channel.ChannelID,
		UserID:           uuid.New().String(),
		AuthorID:         "1",
		MentionedUserIds: []string{},
	}
	err = dydb.CreateMessage(msg)
	assert.NoError(t, err)
	rmsg, err := dydb.GetMessage(msg.ContainerID, msg.MessageID)
	assert.NoError(t, err)
	assert.Equal(t, *msg, *rmsg)

	// Create second message in the channel
	msg2 := &Message{
		DocID:       docId,
		ContainerID: fmt.Sprintf("%s%s", ChannelPrefix, channel.ChannelID),
		ChannelID:   channel.ChannelID,
		UserID:      uuid.New().String(),
		AuthorID:    "1",
	}
	err = dydb.CreateMessage(msg2)
	assert.NoError(t, err)

	// Create a threaded reply to the first message
	threadMsg := &Message{
		DocID:             docId,
		ChannelID:         channel.ChannelID,
		ContainerID:       fmt.Sprintf("%s%s", MsgPrefix, msg.MessageID),
		UserID:            uuid.New().String(),
		AuthorID:          "1",
		ParentContainerID: &msg.ContainerID,
	}
	err = dydb.CreateMessage(threadMsg)
	assert.NoError(t, err)
	updatedMsg, err := dydb.GetMessage(msg.ContainerID, msg.MessageID)
	assert.NoError(t, err)
	// Get the thread of messages
	msgs, err := dydb.GetThreadMessages(msg.ContainerID, msg.MessageID)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(msgs))

	// Get just the channel messages
	msgs, err = dydb.GetChannelMessages(channel.ChannelID)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(msgs))
	assert.Equal(t, updatedMsg, msgs[0])
}
