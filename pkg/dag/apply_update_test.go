package dag_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/teamreviso/code/pkg/dag"
	"github.com/teamreviso/code/pkg/testutils"
	v3 "github.com/teamreviso/code/rogue/v3"
)

func TestApplyUpdate_replace_entire_document(t *testing.T) {
	ctx := testutils.TestContext()
	id := uuid.NewString()
	doc := v3.NewRogueForQuill("0")

	state := dag.NewState(map[string]any{})
	ctx = dag.WithDagState(ctx, state)

	_, err := doc.Insert(0, essay)
	require.NoError(t, err)

	revision, err := dag.ApplyUpdate(ctx, dag.ApplyInput{
		DocId:    id,
		AuthorId: "0",
		Update:   "replacement of entire document",
		Document: doc,
	})

	require.NoError(t, err)
	require.Equal(t, "replacement of entire document", revision.Updated)
	require.Equal(t, "replacement of entire document\n", doc.GetText())
}

func TestApplyUpdate_edit_target(t *testing.T) {
	ctx := testutils.TestContext()
	id := uuid.NewString()
	doc := v3.NewRogueForQuill("0")

	state := dag.NewState(map[string]any{})
	ctx = dag.WithDagState(ctx, state)

	_, err := doc.Insert(0, essay)
	require.NoError(t, err)

	beforeID, err := doc.TotLeftOf(v3.ID{Author: "0", Seq: 3})
	require.NoError(t, err)
	afterID := v3.ID{Author: "0", Seq: 445}

	revision, err := dag.ApplyUpdate(ctx, dag.ApplyInput{
		DocId:    id,
		AuthorId: "0",
		Update:   "replacement first paragraph",
		Document: doc,
		EditTarget: &dag.EditTarget{
			ID:       id,
			BeforeID: beforeID,
			AfterID:  afterID,
			Markdown: `Lifelong learning is essential in today’s rapidly changing world, where new technologies and evolving industries constantly reshape the job market. Engaging in continuous education allows individuals to stay relevant and competitive, ensuring they possess the necessary skills to adapt to new challenges. Moreover, lifelong learning fosters personal growth and intellectual fulfillment, contributing to a more satisfying and meaningful life.`,
		},
	})

	require.NoError(t, err)
	require.Equal(t, "replacement first paragraph", revision.Updated)

	mkdown, err := doc.GetFullMarkdown()
	require.NoError(t, err)
	require.Equal(t, "replacement first paragraph\n\n\n\nOne significant advantage of lifelong learning is the ability to enhance career prospects and job security. As industries advance, employees who continuously update their knowledge and skills are more likely to secure promotions and navigate career transitions successfully. Additionally, employers value workers who demonstrate a commitment to personal development, often leading to increased opportunities and professional recognition. Beyond the workplace, acquiring new skills can open doors to diverse hobbies and interests, enriching one’s personal life.\n\n\n\nFurthermore, lifelong learning promotes critical thinking and problem-solving abilities, which are crucial in both professional and personal settings. By engaging with new ideas and perspectives, individuals become more adaptable and resilient in the face of change. This ongoing intellectual stimulation also contributes to mental well-being, reducing the risk of cognitive decline as one ages. In essence, lifelong learning not only supports economic and career advancement but also enhances overall quality of life, making it a vital pursuit for individuals of all ages.\n\n", mkdown)
}
