package dynamo

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/teamreviso/code/pkg/models"
)

func TestTimelineEvent_Create(t *testing.T) {
	db, err := NewDB()
	require.NoError(t, err)
	started := time.Now().UnixMicro()

	docId := uuid.NewString()
	userID := uuid.NewString()
	authorID := "0"

	event := &TimelineEvent{
		DocID:    docId,
		UserID:   userID,
		AuthorID: authorID,
		Event: &models.TimelineEventPayload{
			Payload: &models.TimelineEventPayload_Marker{
				Marker: &models.TimelineMarkerV1{
					Title: "First Version",
				},
			},
		},
	}

	err = db.CreateTimelineEvent(event)
	require.NoError(t, err)

	timeline, err := db.GetDocumentTimeline(docId)
	require.NoError(t, err)
	require.Len(t, timeline, 1)

	require.Equal(t, docId, timeline[0].DocID)
	require.Equal(t, userID, timeline[0].UserID)
	require.Equal(t, authorID, timeline[0].AuthorID)
	require.Equal(t, "", timeline[0].ReplyToID)
	require.Greater(t, timeline[0].CreatedAt, started)
}

func TestTimelineEvent_Create_reply(t *testing.T) {
	db, err := NewDB()
	require.NoError(t, err)
	started := time.Now().UnixMicro()

	docId := uuid.NewString()
	userID := uuid.NewString()
	userID2 := uuid.NewString()
	authorID := "0"
	authorID2 := "1"

	event := &TimelineEvent{
		DocID:    docId,
		UserID:   userID,
		AuthorID: authorID,
		Event: &models.TimelineEventPayload{
			Payload: &models.TimelineEventPayload_Message{
				Message: &models.TimelineMessageV1{
					Content: "I really like this part about Tacos",
				},
			},
		},
	}
	err = db.CreateTimelineEvent(event)
	require.NoError(t, err)

	replyEvent := &TimelineEvent{
		DocID:     docId,
		UserID:    userID2,
		AuthorID:  authorID2,
		ReplyToID: event.EventID,
		Event: &models.TimelineEventPayload{
			Payload: &models.TimelineEventPayload_Message{
				Message: &models.TimelineMessageV1{
					Content: "This is horrible, make it about burgers",
				},
			},
		},
	}
	err = db.CreateTimelineEvent(replyEvent)
	require.NoError(t, err)

	timeline, err := db.GetDocumentTimeline(docId)
	require.NoError(t, err)
	require.Len(t, timeline, 1)

	require.Equal(t, docId, timeline[0].DocID)
	require.Equal(t, userID, timeline[0].UserID)
	require.Equal(t, authorID, timeline[0].AuthorID)
	require.Equal(t, "", timeline[0].ReplyToID)
	require.Greater(t, timeline[0].CreatedAt, started)

	replyTimeline, err := db.GetDocumentTimelineReplies(docId, event.EventID)
	require.NoError(t, err)
	require.Len(t, replyTimeline, 1)

	require.Equal(t, docId, replyTimeline[0].DocID)
	require.Equal(t, userID2, replyTimeline[0].UserID)
	require.Equal(t, authorID2, replyTimeline[0].AuthorID)
	require.Equal(t, event.EventID, replyTimeline[0].ReplyToID)
	require.Greater(t, replyTimeline[0].CreatedAt, started)
}

func TestTimeline_GetLastUserUpdate(t *testing.T) {
	db, err := NewDB()
	require.NoError(t, err)

	docId := uuid.NewString()
	userID := uuid.NewString()
	authorID := "0"

	event := &TimelineEvent{
		DocID:    docId,
		UserID:   userID,
		AuthorID: authorID,
		Event: &models.TimelineEventPayload{
			Payload: &models.TimelineEventPayload_Update{
				Update: &models.TimelineDocumentUpdateV1{
					Title:   "First Version",
					Content: "I really like this part about Tacos",
					State:   models.UpdateState_COMPLETE_STATE,
				},
			},
		},
	}
	err = db.CreateTimelineEvent(event)
	require.NoError(t, err)

	event2 := &TimelineEvent{
		DocID:    docId,
		UserID:   userID,
		AuthorID: authorID,
		Event: &models.TimelineEventPayload{
			Payload: &models.TimelineEventPayload_Update{
				Update: &models.TimelineDocumentUpdateV1{
					Title:   "Second Version",
					Content: "I hate the part about Burritos",
					State:   models.UpdateState_COMPLETE_STATE,
				},
			},
		},
	}
	err = db.CreateTimelineEvent(event2)
	require.NoError(t, err)

	event3 := &TimelineEvent{
		DocID:    docId,
		UserID:   userID,
		AuthorID: authorID,
		Event: &models.TimelineEventPayload{
			Payload: &models.TimelineEventPayload_Update{
				Update: &models.TimelineDocumentUpdateV1{
					State: models.UpdateState_SUMMARIZING_STATE,
				},
			},
		},
	}
	err = db.CreateTimelineEvent(event3)
	require.NoError(t, err)

	returnedEvent, err := db.GetLastUserUpdate(docId, userID)
	require.NoError(t, err)
	require.Equal(t, event3.EventID, returnedEvent.EventID)
}
