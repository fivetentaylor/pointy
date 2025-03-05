package dag_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/teamreviso/code/pkg/constants"
	"github.com/teamreviso/code/pkg/dag"
	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/models"
	"github.com/teamreviso/code/pkg/storage/dynamo"
	"github.com/teamreviso/code/pkg/testutils"
	v3 "github.com/teamreviso/code/rogue/v3"
)

type AskNodeTestSuite struct {
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

func (s *AskNodeTestSuite) SetupTest() {
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

	err = env.Dynamo(s.ctx).CreateMessage(s.inMsg)
	s.Require().NoError(err)

	err = env.Dynamo(s.ctx).CreateMessage(s.outMsg)
	s.Require().NoError(err)

	err = env.Dynamo(s.ctx).CreateThread(s.thrd)
	s.Require().NoError(err)
}

func (s *AskNodeTestSuite) Test_AskWithMessage() {
	adapter := NewTestLLMAdapterWithResponse("<message>I hate bananas</message>")
	node := dag.AskNode{}
	node.SetAdapter(adapter)

	d := dag.New("test", &node)
	err := d.Run(s.ctx, map[string]any{
		"docId":           s.inMsg.DocID,
		"inputMessageId":  s.inMsg.MessageID,
		"outputMessageId": s.outMsg.MessageID,
		"threadId":        s.inMsg.ChannelID,
		"authorId":        s.inMsg.AuthorID,
	})
	s.Require().NoError(err)

	msg, err := env.Dynamo(s.ctx).GetAiThreadMessage(
		s.outMsg.ChannelID,
		s.outMsg.MessageID,
	)
	s.Require().NoError(err)

	var conclusion *models.Attachment_Content
	for _, a := range msg.Attachments.Attachments {
		if v, ok := a.Value.(*models.Attachment_Content); ok {
			conclusion = v
		}
	}
	s.NotNil(conclusion)
	s.Equal("I hate bananas", conclusion.Content.Text)
	s.Equal("answer", conclusion.Content.Role)
}

func (s *AskNodeTestSuite) Test_AskWithMessageContent() {
	adapter := NewTestLLMAdapterWithResponse("<message><content>I hate bananas</content></message>")
	node := dag.AskNode{}
	node.SetAdapter(adapter)

	d := dag.New("test", &node)
	err := d.Run(s.ctx, map[string]any{
		"docId":           s.inMsg.DocID,
		"inputMessageId":  s.inMsg.MessageID,
		"outputMessageId": s.outMsg.MessageID,
		"threadId":        s.inMsg.ChannelID,
		"authorId":        s.inMsg.AuthorID,
	})
	s.Require().NoError(err)

	msg, err := env.Dynamo(s.ctx).GetAiThreadMessage(
		s.outMsg.ChannelID,
		s.outMsg.MessageID,
	)
	s.Require().NoError(err)

	var conclusion *models.Attachment_Content
	for _, a := range msg.Attachments.Attachments {
		if v, ok := a.Value.(*models.Attachment_Content); ok {
			conclusion = v
		}
	}
	s.NotNil(conclusion)
	s.Equal("I hate bananas", conclusion.Content.Text)
	s.Equal("answer", conclusion.Content.Role)
}

func TestAskNode(t *testing.T) {
	suite.Run(t, new(AskNodeTestSuite))
}
