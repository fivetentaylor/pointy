package dag_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/fivetentaylor/pointy/pkg/dag"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/models"
	"github.com/fivetentaylor/pointy/pkg/rogue"
	"github.com/fivetentaylor/pointy/pkg/testutils"
	v3 "github.com/fivetentaylor/pointy/rogue/v3"
)

type SummarizeUpdateNodeTestSuite struct {
	suite.Suite

	ctx context.Context

	docId    string
	threadId string
	userId   string

	state *dag.State

	doc      *v3.Rogue
	docStore *rogue.DocStore

	prompt *models.Prompt
}

func (s *SummarizeUpdateNodeTestSuite) SetupTest() {
	testutils.EnsureStorage()
	ctx := testutils.TestContext()

	s.docId = uuid.NewString()
	s.threadId = uuid.NewString()
	s.userId = uuid.NewString()

	s.state = dag.NewState(map[string]any{
		"docId":    s.docId,
		"threadId": s.threadId,
		"userId":   s.userId,
	})

	s.ctx = dag.WithDagState(ctx, s.state)

	s.doc, s.docStore = testutils.CreateTestDocument(s.T(), s.ctx, s.docId, "Hello I am a doc")

	authorId := &models.AuthorID{
		AuthorID:   0,
		DocumentID: s.docId,
		UserID:     s.userId,
	}
	err := env.Query(s.ctx).AuthorID.Save(authorId)
	s.Require().NoError(err)
}

func (s *SummarizeUpdateNodeTestSuite) TestDocDifference_emptyDoc() {
	newDocID := uuid.NewString()
	s.doc, s.docStore = testutils.CreateTestDocument(s.T(), s.ctx, newDocID, "")
	authorId := &models.AuthorID{
		AuthorID:   0,
		DocumentID: newDocID,
		UserID:     s.userId,
	}
	err := env.Query(s.ctx).AuthorID.Save(authorId)
	s.Require().NoError(err)

	_, err = s.doc.Insert(0, "draft")
	s.Require().NoError(err)

	s.docStore.SaveDocToS3(s.ctx, newDocID, 2, s.doc)

	startingAddress, err := s.doc.GetEmptyAddress()
	s.Require().NoError(err)
	endingAddress, err := s.doc.GetFullAddress()
	s.Require().NoError(err)

	diff, err := dag.SummarizeDocDifference(s.ctx, newDocID, s.userId, startingAddress, endingAddress)
	s.Require().NoError(err)

	s.Equal("Before:\n```\n\n\n\n```\n\nAfter:\n```\ndraft\n\n\n```", diff)
}

func (s *SummarizeUpdateNodeTestSuite) TestDocDifference_singleAuthor() {
	startingAddress, err := s.doc.GetFullAddress()
	s.Require().NoError(err)

	_, err = s.doc.Delete(13, 3)
	s.Require().NoError(err)
	_, err = s.doc.Insert(13, "draft")
	s.Require().NoError(err)

	s.docStore.SaveDocToS3(s.ctx, s.docId, 2, s.doc)

	endingAddress, err := s.doc.GetFullAddress()
	s.Require().NoError(err)

	diff, err := dag.SummarizeDocDifference(s.ctx, s.docId, s.userId, startingAddress, endingAddress)
	s.Require().NoError(err)

	s.Equal(
		"Before:\n```\nHello I am a doc\n\n\n```\n\nAfter:\n```\nHello I am a draft\n\n\n```",
		diff,
	)
}

func (s *SummarizeUpdateNodeTestSuite) TestDocDifference_MultiAuthor() {
	otherUserID := uuid.NewString()
	otherUserAuthorId := int32(1)

	authorId := &models.AuthorID{
		AuthorID:   otherUserAuthorId,
		DocumentID: s.docId,
		UserID:     otherUserID,
	}
	err := env.Query(s.ctx).AuthorID.Save(authorId)
	s.Require().NoError(err)

	startingAddress, err := s.doc.GetFullAddress()
	s.Require().NoError(err)

	s.doc.Author = fmt.Sprintf("%d", otherUserAuthorId)
	_, err = s.doc.Delete(13, 3)
	s.Require().NoError(err)
	_, err = s.doc.Insert(13, "draft")
	s.Require().NoError(err)

	// back to original author
	s.doc.Author = fmt.Sprintf("%d", 0)
	_, err = s.doc.Delete(0, 5)
	s.Require().NoError(err)
	_, err = s.doc.Insert(0, "Goodbye")
	s.Require().NoError(err)

	s.docStore.SaveDocToS3(s.ctx, s.docId, 2, s.doc)

	endingAddress, err := s.doc.GetFullAddress()
	s.Require().NoError(err)

	diff, err := dag.SummarizeDocDifference(s.ctx, s.docId, s.userId, startingAddress, endingAddress)
	s.Require().NoError(err)

	// Notice that author 1 changes are included in the before and after summaries, so their changes will not be summarized
	s.Equal(
		"Before:\n```\nHello I am a draft\n\n\n```\n\nAfter:\n```\nGoodbye I am a draft\n\n\n```",
		diff,
	)
}

func TestSummarizeUpdateNode(t *testing.T) {
	suite.Run(t, new(SummarizeUpdateNodeTestSuite))
}
