package ai_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/fivetentaylor/pointy/pkg/constants"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/models"
	"github.com/fivetentaylor/pointy/pkg/storage/dynamo"
)

type testConvoOption func(*dynamo.Message)

func createTestConvo(t *testing.T, ctx context.Context, options ...testConvoOption) (string, string, string, string) {
	docID := uuid.NewString()
	channelID := uuid.NewString()
	messageID := uuid.NewString()
	aiMsgID := uuid.NewString()
	userID := uuid.NewString()

	dydb := env.Dynamo(ctx)

	parentContainerID := fmt.Sprintf("%s%s", dynamo.ChannelPrefix, channelID)

	initialMessage := &dynamo.Message{
		DocID:          docID,
		ContainerID:    parentContainerID,
		ChannelID:      channelID,
		MessageID:      messageID,
		UserID:         userID,
		AuthorID:       userID[0:1],
		Content:        "hello world",
		LifecycleStage: dynamo.MessageLifecycleStageCompleted,
		Attachments: &models.AttachmentList{
			Attachments: []*models.Attachment{},
		},
	}

	for _, opt := range options {
		opt(initialMessage)
	}

	err := dydb.CreateMessage(initialMessage)
	if err != nil {
		t.Fatalf("CreateMessage() error = %v", err)
	}
	containerID := fmt.Sprintf("%s%s", dynamo.MsgPrefix, messageID)
	err = dydb.CreateMessage(&dynamo.Message{
		DocID:             docID,
		ParentContainerID: &parentContainerID,
		ContainerID:       containerID,
		ChannelID:         channelID,
		MessageID:         aiMsgID,
		UserID:            constants.RevisoUserID,
		AuthorID:          constants.RevisoAuthorID,
		Content:           "",
		LifecycleStage:    dynamo.MessageLifecycleStagePending,
		Attachments: &models.AttachmentList{
			Attachments: []*models.Attachment{},
		},
	})
	if err != nil {
		t.Fatalf("CreateMessage() error = %v", err)
	}

	return docID, channelID, messageID, aiMsgID
}
