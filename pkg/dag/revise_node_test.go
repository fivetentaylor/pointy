package dag_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/fivetentaylor/pointy/pkg/constants"
	"github.com/fivetentaylor/pointy/pkg/dag"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/models"
	"github.com/fivetentaylor/pointy/pkg/storage/dynamo"
	"github.com/fivetentaylor/pointy/pkg/testutils"
	v3 "github.com/fivetentaylor/pointy/rogue/v3"
)

type ReviseNodeTestSuite struct {
	suite.Suite

	ctx      context.Context
	docId    string
	threadId string
	userId   string

	state *dag.State

	inMsg  *dynamo.Message
	outMsg *dynamo.Message
	thrd   *dynamo.Thread
	prompt *models.Prompt
	rogue  *v3.Rogue
}

func (s *ReviseNodeTestSuite) SetupTest() {
	testutils.EnsureStorage()
	s.ctx = testutils.TestContext()

	docID := uuid.New().String()
	threadID := uuid.New().String()
	userID := uuid.New().String()

	s.rogue, _ = testutils.CreateTestDocument(s.T(), s.ctx, docID, "Hello world!")

	ca, err := s.rogue.GetFullAddress()
	s.Require().NoError(err)
	cabts, err := json.Marshal(ca)
	s.Require().NoError(err)

	s.inMsg = &dynamo.Message{
		DocID:          docID,
		ContainerID:    fmt.Sprintf("%s%s", dynamo.AiThreadPrefix, threadID),
		ChannelID:      threadID,
		MessageID:      uuid.New().String(),
		CreatedAt:      time.Now().Unix(),
		AuthorID:       "0",
		UserID:         userID,
		Content:        "message content",
		LifecycleStage: dynamo.MessageLifecycleStageCompleted,
		Attachments: &models.AttachmentList{
			Attachments: []*models.Attachment{},
		},
		MessageMetadata: &models.MessageMetadata{
			AllowDraftEdits: true,
			ContentAddress:  string(cabts),
		},
	}

	s.outMsg = &dynamo.Message{
		DocID:          docID,
		ContainerID:    fmt.Sprintf("%s%s", dynamo.AiThreadPrefix, threadID),
		ChannelID:      threadID,
		MessageID:      uuid.New().String(),
		CreatedAt:      time.Now().Unix(),
		AuthorID:       "!0",
		UserID:         constants.RevisoUserID,
		Content:        "",
		LifecycleStage: dynamo.MessageLifecycleStagePending,
		Attachments: &models.AttachmentList{
			Attachments: []*models.Attachment{},
		},
		MessageMetadata: &models.MessageMetadata{
			AllowDraftEdits: true,
		},
	}

	s.thrd = &dynamo.Thread{
		DocID:    docID,
		ThreadID: threadID,
		UserID:   userID,
	}

	s.prompt = &models.Prompt{
		PromptName:  constants.PromptDraftsConvo,
		Provider:    "test",
		ContentJSON: "{}",
	}

	err = env.Query(s.ctx).Prompt.Save(s.prompt)
	s.Require().NoError(err)

	err = env.Dynamo(s.ctx).CreateMessage(s.inMsg)
	s.Require().NoError(err)

	err = env.Dynamo(s.ctx).CreateMessage(s.outMsg)
	s.Require().NoError(err)

	err = env.Dynamo(s.ctx).CreateThread(s.thrd)
	s.Require().NoError(err)
}

func (s *ReviseNodeTestSuite) TestGenerate() {
	adapter := NewTestLLMAdapterWithResponse("<conclusion>I love bananas</conclusion>")
	node := dag.ReviseNode{}
	node.SetAdapter(adapter)

	d := dag.New("test", &node)
	err := d.Run(s.ctx, map[string]any{
		"docId":           s.inMsg.DocID,
		"inputMessageId":  s.inMsg.MessageID,
		"outputMessageId": s.outMsg.MessageID,
		"threadId":        s.inMsg.ChannelID,
		"authorId":        s.inMsg.AuthorID,
		"promptName":      "testPrompt",
	})
	s.Require().NoError(err)

	msg, err := env.Dynamo(s.ctx).GetAiThreadMessage(
		s.outMsg.ChannelID,
		s.outMsg.MessageID,
	)
	s.Require().NoError(err)

	s.Equal(dynamo.MessageLifecycleStageRevised, msg.LifecycleStage)

	var conclusion *models.Attachment_Content
	for _, a := range msg.Attachments.Attachments {
		if v, ok := a.Value.(*models.Attachment_Content); ok {
			conclusion = v
		}
	}
	s.NotNil(conclusion)
	s.Equal("I love bananas", conclusion.Content.Text)
	s.Equal("conclusion", conclusion.Content.Role)
}

func TestGenerateNode(t *testing.T) {
	suite.Run(t, new(ReviseNodeTestSuite))
}
