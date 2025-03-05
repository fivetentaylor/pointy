package dynamo

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestChannelCreate(t *testing.T) {
	db, err := NewDB()
	assert.NoError(t, err)

	channel := &Channel{
		DocID: uuid.New().String(),
		Type:  ChannelTypeDirect,
	}

	err = db.CreateChannel(channel)
	assert.NoError(t, err)
	assert.NotEmpty(t, channel.ChannelID)
	assert.NotNil(t, channel.UserIDs)

	rt, err := db.GetChannel(channel.DocID, channel.ChannelID)
	assert.NoError(t, err)
	assert.Equal(t, channel, rt)

	channels, err := db.GetDocumentChannels(channel.DocID)
	assert.NoError(t, err)
	assert.Len(t, channels, 1)
	assert.Equal(t, channel, channels[0])
}
