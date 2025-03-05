package notifications

import (
	"context"
	"fmt"

	"github.com/teamreviso/code/pkg/background/wire"
	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/models"
	"github.com/teamreviso/code/pkg/service/email"
	"github.com/teamreviso/code/pkg/utils"
)

func EnqueueNewComment(ctx context.Context, docID, eventID string, excludeUserIds []string) error {
	log := env.Log(ctx)
	bg := env.Background(ctx)
	msg := &wire.NotifyNewTimelineComment{
		DocId:          docID,
		EventId:        eventID,
		ExcludeUserIds: excludeUserIds,
	}
	_, err := bg.Enqueue(ctx, msg)
	if err != nil {
		log.Errorf("error enqueuing timeline mention: %s", err)
		return err
	}

	return nil
}

func SendNewComment(ctx context.Context, docID, eventID string, excludeUserIds []string) error {
	log := env.Log(ctx)
	dydb := env.Dynamo(ctx)
	event, err := env.Dynamo(ctx).GetTimelineEvent(docID, eventID)
	if err != nil {
		env.Log(ctx).Errorf("error getting timeline event: %s", err)
		return err
	}

	log.Info("📣 Sending new comment notifications", "event", event, "docID", docID)

	var message *models.TimelineMessageV1
	switch p := event.Event.Payload.(type) {
	case *models.TimelineEventPayload_Message:
		message = p.Message
	default:
		log.Warnf("unknown payload type: %T", p)
		return nil
	}

	usertbl := env.Query(ctx).User
	doctbl := env.Query(ctx).Document
	docAccessTbl := env.Query(ctx).DocumentAccess

	// find all owners, not including the user
	ownersAccess, err := docAccessTbl.
		Where(docAccessTbl.DocumentID.Eq(docID)).
		Where(docAccessTbl.AccessLevel.In("owner", "write")).
		Where(docAccessTbl.UserID.Neq(event.UserID)).
		Find()
	if err != nil {
		log.Error("error finding owners access", "err", err)
		return err
	}

	fromUser, err := usertbl.Where(usertbl.ID.Eq(event.UserID)).First()
	if err != nil {
		log.Errorf("error getting user for mention email: %s", err)
		return err
	}

	doc, err := doctbl.Where(doctbl.ID.Eq(docID)).First()
	if err != nil {
		log.Errorf("error getting document for mention email: %s", err)
		return err
	}

	var errs []error

	for _, owner := range ownersAccess {
		docPref, err := dydb.GetDocNotificationPreference(owner.UserID, docID)
		if err != nil {
			log.Error("Failed to get notification preference", "err", err)
			errs = append(errs, err)
			continue
		}

		log.Info("📣 checking", "owner", owner.UserID, "mentioned", message.MentionedUserIds, "pref", docPref)

		if utils.Contains(message.MentionedUserIds, owner.UserID) && !utils.Contains(excludeUserIds, owner.UserID) {
			if docPref.Preference.EnableMentionNotifications {
				toUser, err := usertbl.Where(usertbl.ID.Eq(owner.UserID)).First()
				if err != nil {
					log.Errorf("error getting user for mention email: %s", err)
					return err
				}
				err = email.SendTimelineCommentEmail(ctx, doc, message, fromUser, toUser, true)
				if err != nil {
					log.Errorf("error sending mention email: %s", err)
					return err
				}
				continue
			}
		}

		if !docPref.Preference.EnableAllCommentNotifications {
			continue
		}

		toUser, err := usertbl.Where(usertbl.ID.Eq(owner.UserID)).First()
		if err != nil {
			log.Errorf("error getting user for mention email: %s", err)
			return err
		}

		err = email.SendTimelineCommentEmail(ctx, doc, message, fromUser, toUser, false)
		if err != nil {
			log.Errorf("error sending mention email: %s", err)
			return err
		}
	}

	if len(errs) > 0 {
		log.Errorf("error sending timeline mention emails: %s", errs)
	}

	return nil
}

func EnqueueNewMentionShare(ctx context.Context, docID, eventID, recipientID string) error {
	log := env.Log(ctx)
	bg := env.Background(ctx)
	msg := &wire.NotifyNewMentionShare{
		DocId:       docID,
		EventId:     eventID,
		RecipientId: recipientID,
	}
	_, err := bg.Enqueue(ctx, msg)
	if err != nil {
		log.Errorf("error enqueuing timeline mention: %s", err)
		return err
	}

	return nil
}

func SendNewMentionShare(ctx context.Context, docID, eventID, recipientID string) error {
	log := env.Log(ctx)
	if docID == "" || eventID == "" || recipientID == "" {
		log.Warnf("invalid arguments: docID or eventID or recipientID is empty")
		return nil
	}

	event, err := env.Dynamo(ctx).GetTimelineEvent(docID, eventID)
	if event == nil {
		err = fmt.Errorf("timeline event not found: docID: %s, eventID: %s", docID, eventID)
		return err
	}

	if err != nil {
		env.Log(ctx).Errorf("error getting timeline event: %s", err)
		return err
	}

	log.Info("📣 Sending new mention share notifications", "event", event, "docID", docID)

	var message *models.TimelineMessageV1
	switch p := event.Event.Payload.(type) {
	case *models.TimelineEventPayload_Message:
		message = p.Message
	default:
		log.Warnf("unknown payload type: %T", p)
		return nil
	}

	usertbl := env.Query(ctx).User
	doctbl := env.Query(ctx).Document

	fromUser, err := usertbl.Where(usertbl.ID.Eq(event.UserID)).First()
	if err != nil {
		log.Errorf("error getting user for mention email: %s", err)
		return err
	}

	doc, err := doctbl.Where(doctbl.ID.Eq(docID)).First()
	if err != nil {
		log.Errorf("error getting document for mention email: %s", err)
		return err
	}

	toUser, err := usertbl.Where(usertbl.ID.Eq(recipientID)).First()
	if err != nil {
		log.Errorf("error getting user for mention email: %s", err)
		return err
	}

	err = email.SendTimelineMentionShareEmail(ctx, doc, message, fromUser, toUser, false)
	if err != nil {
		log.Errorf("error sending mention email: %s", err)
		return err
	}

	return nil
}

func EnqueueFirstOpen(ctx context.Context, docID, readerID string) error {
	log := env.Log(ctx)

	bg := env.Background(ctx)
	msg := &wire.NotifyFirstOpen{
		DocId:    docID,
		ReaderId: readerID,
	}
	_, err := bg.Enqueue(ctx, msg)
	if err != nil {
		log.Errorf("error enqueuing timeline mention: %s", err)
		return err
	}

	return nil
}

func SendFirstOpen(ctx context.Context, docID, readerID string) error {
	log := env.Log(ctx)
	dydb := env.Dynamo(ctx)
	q := env.Query(ctx)

	docTbl := q.Document
	docAccessTbl := q.DocumentAccess
	userTlb := q.User

	// find all owners, not including the user
	ownersAccess, err := docAccessTbl.
		Where(docAccessTbl.DocumentID.Eq(docID)).
		Where(docAccessTbl.AccessLevel.Eq("owner")).
		Where(docAccessTbl.UserID.Neq(readerID)).
		Find()
	if err != nil {
		log.Error("error finding owners access", "err", err)
		return err
	}

	doc, err := docTbl.Where(docTbl.ID.Eq(docID)).First()
	if err != nil {
		log.Errorf("error finding document: %s", err)
		return err
	}

	reader, err := userTlb.Where(userTlb.ID.Eq(readerID)).First()
	if err != nil {
		log.Errorf("error finding user (reader): %s", err)
		return err
	}

	for _, owner := range ownersAccess {
		docPref, err := dydb.GetDocNotificationPreference(owner.UserID, docID)
		if err != nil {
			log.Error("Failed to get notification preference", "err", err)
			return err
		}

		if !docPref.Preference.EnableFirstOpenNotifications {
			log.Info("user preference disable first open notifications", "userID", owner.UserID)
			continue
		}

		toUser, err := userTlb.Where(userTlb.ID.Eq(owner.UserID)).First()
		if err != nil {
			log.Errorf("error finding user: %s", err)
			return err
		}

		err = email.SendFirstOpen(ctx, toUser.Email, reader, doc)
		if err != nil {
			log.Errorf("error sending first open email: %s", err)
			return err
		}
	}

	return nil
}
