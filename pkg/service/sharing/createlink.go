package sharing

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/models"
	"github.com/fivetentaylor/pointy/pkg/query"
)

func CreateShareLink(
	ctx context.Context,
	document *models.Document,
	invitedBy *models.User,
	inviteeEmail,
	customMessage string,
) (*models.SharedDocumentLink, error) {
	sharedLink := &models.SharedDocumentLink{
		InviterID:    invitedBy.ID,
		InviteeEmail: inviteeEmail,
		DocumentID:   document.ID,
	}

	err := env.Query(ctx).Transaction(func(tx *query.Query) error {
		err := tx.SharedDocumentLink.Create(sharedLink)
		if err != nil {
			log.Infof("failed to create share link: %v", err)
			return err
		}

		waitlistUser, err := tx.WaitlistUser.Where(tx.WaitlistUser.Email.Eq(inviteeEmail)).First()
		if err != nil && err.Error() != "record not found" {
			log.Infof("failed to check for waitlist user for invite: %v", err)
			return err
		}
		if waitlistUser != nil {
			return nil
		}

		err = tx.WaitlistUser.Create(&models.WaitlistUser{
			Email:       inviteeEmail,
			AllowAccess: true,
		})

		if err != nil {
			log.Infof("failed to create waitlist user for invite: %v", err)
			return err
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create share link: %w", err)
	}

	return sharedLink, nil
}
