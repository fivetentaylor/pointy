package sharing

import (
	"context"
	"fmt"

	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/models"
	"github.com/fivetentaylor/pointy/pkg/service/email"
	"github.com/fivetentaylor/pointy/pkg/service/pubsub"
	"github.com/fivetentaylor/pointy/pkg/service/timeline"
	"github.com/fivetentaylor/pointy/pkg/storage/dynamo"
)

func ShareDoc(
	ctx context.Context,
	document *models.Document,
	invitedBy *models.User,
	inviteeEmails []string,
	message string,
) ([]*models.SharedDocumentLink, error) {
	log := env.Log(ctx)
	userTbl := env.Query(ctx).User

	var userIdentifiers []string
	var links []*models.SharedDocumentLink

	for _, inviteeEmail := range inviteeEmails {
		invitedUser, _ := userTbl.Where(userTbl.Email.Eq(inviteeEmail)).First()
		if invitedUser != nil {
			log.Infof("user already exists: %s", inviteeEmail)
			err := shareDocToUser(ctx, document, invitedBy, invitedUser, message)
			if err != nil {
				log.Errorf("error sharing doc to user: %s", err)
				return nil, fmt.Errorf("sorry, we could not share your document")
			}
			userIdentifiers = append(userIdentifiers, invitedUser.Name)
			continue
		}

		log.Infof("user doesn't exist: %s", inviteeEmail)
		link, err := shareDocToEmail(ctx, document, invitedBy, inviteeEmail, message)
		if err != nil {
			log.Errorf("error sharing doc to user: %s", err)
			return nil, fmt.Errorf("sorry, we could not share your document")
		}

		links = append(links, link)
		userIdentifiers = append(userIdentifiers, inviteeEmail)
	}

	event := &dynamo.TimelineEvent{
		DocID:  document.ID,
		UserID: invitedBy.ID,
		Event: &models.TimelineEventPayload{
			Payload: &models.TimelineEventPayload_AccessChange{
				AccessChange: &models.TimelineAccessChangeV1{
					Action:          models.TimelineAccessChangeAction_INVITE_ACTION,
					UserIdentifiers: userIdentifiers,
				},
			},
		},
	}

	err := timeline.CreateTimelineEvent(ctx, event)
	if err != nil {
		log.Errorf("error creating timeline event: %s", err)
	}

	return links, nil
}

func shareDocToUser(
	ctx context.Context,
	document *models.Document,
	invitedBy *models.User,
	invitedUser *models.User,
	message string,
) error {
	log := env.Log(ctx)

	role := "write"
	if invitedUser.Educator {
		role = "admin"
	}

	_, err := JoinDoc(ctx, document.ID, invitedUser.ID, role)
	if err != nil {
		log.Errorf("error adding user to document: %s", err)
		return fmt.Errorf("sorry, we could not share your document")
	}

	// update metadata for users who already have the document
	err = pubsub.PublishDocument(ctx, document.ID)
	if err != nil {
		log.Errorf("error publishing document: %s", err)
	}

	// let newly shared users know about the document
	err = pubsub.PublishNewDocument(ctx, invitedUser.ID, document.ID)
	if err != nil {
		log.Errorf("error publishing document: %s", err)
	}

	err = email.SendSharedToUserEmail(
		ctx,
		invitedUser.Email,
		invitedBy.Name,
		message,
		document.ID,
		document.Title,
	)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func shareDocToEmail(
	ctx context.Context,
	document *models.Document,
	invitedBy *models.User,
	inviteeEmail string,
	message string,
) (*models.SharedDocumentLink, error) {
	log := env.Log(ctx)
	shareLink, err := CreateShareLink(ctx, document, invitedBy, inviteeEmail, message)
	if err != nil {
		log.Errorf("error creating share link: %s", err)
		return nil, fmt.Errorf("sorry, we could not share your document")
	}

	err = email.SendShareLinkEmail(
		ctx,
		inviteeEmail,
		invitedBy.Name,
		message,
		document.ID,
		document.Title,
		shareLink.InviteLink,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to send email: %w", err)
	}

	return shareLink, nil
}
