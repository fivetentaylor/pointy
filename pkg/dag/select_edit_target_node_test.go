package dag_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/teamreviso/code/pkg/constants"
	"github.com/teamreviso/code/pkg/dag"
	"github.com/teamreviso/code/pkg/dag/mocks"
	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/models"
	"github.com/teamreviso/code/pkg/storage/dynamo"
	"github.com/teamreviso/code/pkg/testutils"
	v3 "github.com/teamreviso/code/rogue/v3"
)

type SelectFocusNodeTestSuite struct {
	suite.Suite

	ctx      context.Context
	docId    string
	threadId string
	userId   string
	docText  string
	rogueDoc *v3.Rogue

	state *dag.State

	msg    *dynamo.Message
	outMsg *dynamo.Message
	prompt *models.Prompt

	ctrl *gomock.Controller
	llm  *mocks.MockModel
}

func (s *SelectFocusNodeTestSuite) SetupTest() {
	testutils.EnsureStorage()
	s.ctx = testutils.TestContext()
	s.msg = &dynamo.Message{
		MessageID: uuid.NewString(),
	}

	s.ctrl = gomock.NewController(s.T())

	docID := uuid.New().String()
	threadID := uuid.New().String()
	userID := uuid.New().String()

	s.docText = `Two roads diverged in a yellow wood,
And sorry I could not travel both
And be one traveler, long I stood
And looked down one as far as I could
To where it bent in the undergrowth;

Then took the other, as just as fair,
And having perhaps the better claim,
Because it was grassy and wanted wear;
Though as for that the passing there
Had worn them really about the same,

And both that morning equally lay
In leaves no step had trodden black.
Oh, I kept the first for another day!
Yet knowing how way leads on to way,
I doubted if I should ever come back.

I shall be telling this with a sigh
Somewhere ages and ages hence:
Two roads diverged in a wood, and I—
I took the one less traveled by,
And that has made all the difference.`

	s.rogueDoc, _ = testutils.CreateTestDocument(s.T(), s.ctx, docID, s.docText)

	contentAddress, err := s.rogueDoc.GetFullAddress()
	s.Require().NoError(err)
	bts, err := json.Marshal(contentAddress)
	s.Require().NoError(err)

	s.msg = &dynamo.Message{
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
			ContentAddress:  string(bts),
		},
	}

	s.outMsg = &dynamo.Message{
		DocID:          docID,
		ContainerID:    fmt.Sprintf("%s%s", dynamo.AiThreadPrefix, threadID),
		ChannelID:      threadID,
		MessageID:      uuid.New().String(),
		CreatedAt:      time.Now().Unix(),
		AuthorID:       constants.RevisoAuthorID,
		UserID:         constants.RevisoUserID,
		LifecycleStage: dynamo.MessageLifecycleStageCompleted,
		Attachments: &models.AttachmentList{
			Attachments: []*models.Attachment{},
		},
	}

	err = env.Dynamo(s.ctx).CreateMessage(s.msg)
	s.Require().NoError(err)

	err = env.Dynamo(s.ctx).CreateMessage(s.outMsg)
	s.Require().NoError(err)
}

func (s *SelectFocusNodeTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *SelectFocusNodeTestSuite) TestSelectFocusNode_FullDocument() {
	node := dag.SelectEditTargetNode{}
	llmResponse := `<response>
<full_document></full_document>
</response>`

	adapter := NewTestLLMAdapterWithResponse(llmResponse)
	node.SetAdapter(adapter)

	d := dag.New("test select focus node", &node)
	err := d.Run(s.ctx, map[string]any{
		"docId":           s.msg.DocID,
		"inputMessageId":  s.msg.MessageID,
		"outputMessageId": s.outMsg.MessageID,
		"threadId":        s.msg.ChannelID,
		"authorId":        s.msg.AuthorID,
	})
	s.Require().NoError(err)

	state := d.State()
	s.Require().NotNil(state)
	s.Require().Nil(state.Get("editTargets")) // Nil means the full document is selected
}

func (s *SelectFocusNodeTestSuite) TestSelectFocusNode_Paragraph() {
	node := dag.SelectEditTargetNode{}

	llmResponse := `<response>
<relevant_section>
<p data-rid="0_145" />
</relevant_section>
</response>`

	adapter := NewTestLLMAdapterWithResponse(llmResponse)
	node.SetAdapter(adapter)

	d := dag.New("test select focus node", &node)
	err := d.Run(s.ctx, map[string]any{
		"docId":           s.msg.DocID,
		"inputMessageId":  s.msg.MessageID,
		"outputMessageId": s.outMsg.MessageID,
		"threadId":        s.msg.ChannelID,
		"authorId":        s.msg.AuthorID,
	})
	s.Require().NoError(err)

	state := d.State()
	s.Require().NotNil(state)

	html, err := s.rogueDoc.GetFullHtml(true, false)
	s.Require().NoError(err)
	fmt.Printf("doc html: %q\n", html)

	ets := state.Get("editTargets").([]dag.EditTarget)
	s.Require().Len(ets, 1)
	et := ets[0]

	s.Require().Equal(v3.ID{Author: "0", Seq: 107}, et.BeforeID)
	s.Require().Equal(v3.ID{Author: "0", Seq: 145}, et.AfterID)
}

func (s *SelectFocusNodeTestSuite) TestSelectFocusNode_large_chunk_as_relevant() {
	node := dag.SelectEditTargetNode{}

	llmResponse := `
<response>
<relevant_section>
  <p data-rid="0_39" /><span data-rid="0_594" />
</relevant_section>
</response>`
	adapter := NewTestLLMAdapterWithResponse(llmResponse)
	node.SetAdapter(adapter)

	d := dag.New("test select focus node", &node)
	err := d.Run(s.ctx, map[string]any{
		"docId":           s.msg.DocID,
		"inputMessageId":  s.msg.MessageID,
		"outputMessageId": s.outMsg.MessageID,
		"threadId":        s.msg.ChannelID,
		"authorId":        s.msg.AuthorID,
	})
	s.Require().NoError(err)

	state := d.State()
	s.Require().NotNil(state)

	/*
		html, err := s.rogueDoc.GetFullHtml(true)
		s.Require().NoError(err)
		fmt.Printf("doc html: %q\n", html)
	*/

	ets := state.Get("editTargets").([]dag.EditTarget)
	s.Require().Len(ets, 1)
	et := ets[0]

	s.Require().Equal(v3.ID{Author: "root", Seq: 0}, et.BeforeID)
	s.Require().Equal(v3.ID{Author: "0", Seq: 624}, et.AfterID)

	/*
		html, err = s.rogueDoc.GetHtml(et.BeforeID, et.AfterID, true)
		s.Require().NoError(err)
		fmt.Printf("selected html: %q\n", html)
	*/
}

func TestSelectFocusNode(t *testing.T) {
	suite.Run(t, new(SelectFocusNodeTestSuite))
}
