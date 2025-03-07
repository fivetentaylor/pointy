package timeline

import (
	"context"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/models"
	"github.com/fivetentaylor/pointy/pkg/query"
	"github.com/fivetentaylor/pointy/pkg/service/notifications"
	"github.com/fivetentaylor/pointy/pkg/service/pubsub"
	"github.com/fivetentaylor/pointy/pkg/storage/dynamo"
	"gorm.io/gorm"
)

func CreateTimelineEvent(ctx context.Context, event *dynamo.TimelineEvent) error {
	log := env.SLog(ctx)
	log.Info("creating timeline event", "eventID", event.EventID, "docID", event.DocID)

	mentionShareUsers, err := ProcessMentions(ctx, event)

	if err != nil {
		log.Error("error processing mentions", "error", err)
		return err
	}

	err = env.Dynamo(ctx).CreateTimelineEvent(event)
	if err != nil {
		log.Error("error creating timeline event", "error", err)
		return err
	}

	eventToPublish := event
	eventType := EventTypeInsert

	if event.ReplyToID != "" {
		// Publish an update event for the parent comment
		eventToPublish, err = env.Dynamo(ctx).GetTimelineEvent(event.DocID, event.ReplyToID)
		eventType = EventTypeUpdate
		if err != nil {
			log.Error("error fetching parent event", "error", err)
			return err
		}
	}

	err = PublishTimelineEvent(ctx, eventToPublish, eventType)
	if err != nil {
		log.Error("error publishing timeline event", "error", err)
		return err
	}

	err = notifications.EnqueueNewComment(ctx, event.DocID, event.EventID, mentionShareUsers)
	if err != nil {
		log.Error("error sending new comment email", "error", err)
		return err
	}

	err = shareMentionedUsers(ctx, event, mentionShareUsers)
	if err != nil {
		log.Error("error sending mention notifications", "error", err)
		return err
	}

	log.Info("timeline event created", "eventID", event.EventID, "docID", event.DocID, "event", "timeline_event_created")
	return nil
}

func UpdateTimelineEvent(ctx context.Context, event *dynamo.TimelineEvent) error {
	log := env.Log(ctx)
	log.Debug("updating timeline event", "eventID", event.EventID, "docID", event.DocID)

	mentionShareUsers, err := ProcessMentions(ctx, event)
	if err != nil {
		log.Errorf("error processing mentions: %s", err)
		return err
	}

	err = env.Dynamo(ctx).UpdateTimelineEvent(event)
	if err != nil {
		log.Errorf("error updating timeline event: %s", err)
		return err
	}

	err = PublishTimelineEvent(ctx, event, EventTypeUpdate)
	if err != nil {
		log.Errorf("error publishing timeline event: %s", err)
		return err
	}

	err = notifications.EnqueueNewComment(ctx, event.DocID, event.EventID, mentionShareUsers)
	if err != nil {
		log.Errorf("error sending new comment email: %s", err)
		return err
	}

	err = shareMentionedUsers(ctx, event, mentionShareUsers)
	if err != nil {
		log.Errorf("error sending mention notifications: %s", err)
		return err
	}

	return nil
}

func DeleteTimelineEvent(ctx context.Context, event *dynamo.TimelineEvent) error {
	log := env.Log(ctx)
	log.Debug("updating timeline event", "eventID", event.EventID, "docID", event.DocID)

	err := env.Dynamo(ctx).DeleteTimelineEvent(event.DocID, event.EventID)
	if err != nil {
		log.Errorf("error updating timeline event: %s", err)
		return err
	}

	eventToPublish := event
	eventType := EventTypeDelete

	err = PublishTimelineEvent(ctx, eventToPublish, eventType)
	if err != nil {
		log.Errorf("error publishing timeline event: %s", err)
		return err
	}

	return nil
}

type ProcessMentionsOptions struct {
	ExistingIDs []string
}

type MentionStatus int

const (
	MentionStatusNew MentionStatus = iota
	MentionStatusExisting
)

func ProcessMentions(ctx context.Context, event *dynamo.TimelineEvent) ([]string, error) {
	log := env.Log(ctx)
	mentionedUsers := make(map[string]MentionStatus)
	sliceMentionedUsers := make([]string, 0)

	switch message := event.Event.Payload.(type) {
	case *models.TimelineEventPayload_Message:
		var existingIDs []string

		// if it's an existing event, we need to get the existing mentioned user IDs
		if event.EventID != "" {
			existingEvent, err := env.Dynamo(ctx).GetTimelineEvent(event.DocID, event.EventID)
			if err != nil {
				log.Errorf("error fetching existing event: %s", err)
				return nil, err
			}
			existingIDs = existingEvent.Event.Payload.(*models.TimelineEventPayload_Message).Message.MentionedUserIds
			log.Info("existing mentioned user IDs", "existingIDs", existingIDs)
		}

		// get the mentions from the message content
		mentions := ExtractMentions(ctx, message.Message.Content, existingIDs)

		if len(mentions) == 0 {
			log.Info("no new mentions")
			return nil, nil
		}

		for _, mention := range mentions {
			userID := mention.ID

			// first, check if it's actually a new user
			if userID == "new-user-id" {
				userTbl := env.Query(ctx).User
				user, err := userTbl.Where(userTbl.Email.Eq(mention.Name)).First()
				if err != nil {
					if err == gorm.ErrRecordNotFound {
						user, err = createNewUser(ctx, mention)
						if err != nil {
							log.Errorf("error creating new user: %s", err)
							return nil, err
						}
						userID = user.ID
						message.Message.Content = UpdateBase64Mentions(message.Message.Content, user, UpdateBase64MentionsOptions{SpecificMatch: mention.Base64Match})
					} else {
						log.Errorf("error checking for user: %s", err)
						return nil, err
					}
				} else {
					// User exists, use the found user
					userID = user.ID
					message.Message.Content = UpdateBase64Mentions(message.Message.Content, user, UpdateBase64MentionsOptions{SpecificMatch: mention.Base64Match})
				}
			}

			docAccessTbl := env.Query(ctx).DocumentAccess

			hasAccess, err := docAccessTbl.Where(docAccessTbl.DocumentID.Eq(event.DocID), docAccessTbl.UserID.Eq(userID)).Count()
			if err != nil {
				log.Errorf("error checking document access: %s", err)
				return nil, err
			}

			if hasAccess < 1 {
				mentionedUsers[userID] = MentionStatusNew
			} else {
				mentionedUsers[userID] = MentionStatusExisting
			}
		}

		sliceUserIds := make([]string, 0, len(mentionedUsers))
		for userID := range mentionedUsers {
			sliceUserIds = append(sliceUserIds, userID)
			if mentionedUsers[userID] == MentionStatusNew {
				sliceMentionedUsers = append(sliceMentionedUsers, userID)
			}
		}

		// finally, set the unique mentioned user IDs on the event
		event.Event.Payload.(*models.TimelineEventPayload_Message).Message.MentionedUserIds = sliceUserIds
	}

	return sliceMentionedUsers, nil
}

func createNewUser(ctx context.Context, mention Mention) (*models.User, error) {
	userID := uuid.NewString()
	user := &models.User{
		ID:          userID,
		Email:       mention.Name,
		Name:        mention.Name,
		DisplayName: mention.Name,
		Provider:    "manual",
	}
	env.Query(ctx).Transaction(func(tx *query.Query) error {
		err := tx.User.Create(user)
		if err != nil {
			log.Errorf("error creating user: %s", err)
			return err
		}

		waitlistUser, err := tx.WaitlistUser.Where(tx.WaitlistUser.Email.Eq(mention.Name)).First()
		if err != nil && err.Error() != "record not found" {
			log.Infof("failed to check for waitlist user for invite: %v", err)
			return err
		}
		if waitlistUser != nil {
			return nil
		}

		err = tx.WaitlistUser.Create(&models.WaitlistUser{
			Email:       mention.Name,
			AllowAccess: true,
		})

		if err != nil {
			log.Infof("failed to create waitlist user for invite: %v", err)
			return err
		}

		return nil
	})

	return user, nil
}

func addMentionedUser(ctx context.Context, event *dynamo.TimelineEvent, userID string) error {
	log := env.Log(ctx)
	docAccessTbl := env.Query(ctx).DocumentAccess

	err := docAccessTbl.Create(&models.DocumentAccess{
		DocumentID:  event.DocID,
		UserID:      userID,
		AccessLevel: "write",
	})

	if err != nil {
		log.Error("error adding document access", "err", err)
		return err
	}

	userTbl := env.Query(ctx).User
	user, err := userTbl.Where(userTbl.ID.Eq(userID)).First()
	if err != nil {
		log.Errorf("error fetching user: %s", err)
		return err
	}

	err = pubsub.PublishDocument(ctx, event.DocID)
	if err != nil {
		log.Errorf("error publishing document: %s", err)
	}

	shareEvent := &dynamo.TimelineEvent{
		DocID:  event.DocID,
		UserID: event.UserID,
		Event: &models.TimelineEventPayload{
			Payload: &models.TimelineEventPayload_AccessChange{
				AccessChange: &models.TimelineAccessChangeV1{
					Action:          models.TimelineAccessChangeAction_INVITE_ACTION,
					UserIdentifiers: []string{user.Email},
				},
			},
		},
	}

	err = CreateTimelineEvent(ctx, shareEvent)
	if err != nil {
		log.Errorf("error creating timeline event: %s", err)
		return err
	}

	return nil
}

func shareMentionedUsers(ctx context.Context, event *dynamo.TimelineEvent, mentionedUsers []string) error {
	log := env.Log(ctx)

	if len(mentionedUsers) == 0 {
		log.Info("no mentioned users")
		return nil
	}

	// iterate through mentionedUsers
	for _, userID := range mentionedUsers {
		err := addMentionedUser(ctx, event, userID)
		if err != nil {
			log.Errorf("error adding mentioned user: %s", err)
			return err
		}

		log.Info("📣 enqueueing new mention share", "docID", event.DocID, "eventID", event.EventID, "recipientID", userID)

		err = notifications.EnqueueNewMentionShare(ctx, event.DocID, event.EventID, userID)
		if err != nil {
			log.Errorf("error sending mention emails: %s", err)
			return err
		}

	}

	return nil
}
