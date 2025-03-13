package email

import (
	"context"
	"fmt"
	"strings"

	"github.com/fivetentaylor/pointy/pkg/env"

	"github.com/charmbracelet/log"
	"github.com/fivetentaylor/pointy/pkg/service/email/templates"
)

// SendSharedToUserEmail sends an email to a user that already exists
func SendSharedToUserEmail(
	ctx context.Context,
	to,
	invitedByName,
	customMessage,
	documentId,
	documentTitle string,
) error {
	c := env.SES(ctx)
	from := fmt.Sprintf("Pointy <pointy@%s>", c.EmailDomain())
	subject := fmt.Sprintf("%s has shared a '%s' with you!", invitedByName, documentTitle)
	preheader := "A Pointy document has been shared with you"
	link := fmt.Sprintf("%s/drafts/%s", c.AppHost(), documentId)

	if c.Env() == "development" {
		log.Infof("DEV share link:\n\n %s\n\n for %s", link, to)
	}

	log.Infof("sending share doc email to %s", to)

	rctx := c.AttachHostValues(ctx)

	htmlBody := &strings.Builder{}

	err := templates.ShareDocHTML(preheader, documentTitle, invitedByName, customMessage, "View document", link).Render(rctx, htmlBody)
	if err != nil {
		return fmt.Errorf("failed to render share doc html: %w", err)
	}

	textBody := &strings.Builder{}
	err = templates.ShareDocText(documentTitle, invitedByName, customMessage, "View document", link).Render(rctx, textBody)
	if err != nil {
		return fmt.Errorf("failed to render share doc text: %w", err)
	}

	return c.EnqueueEmail(from, to, subject, textBody.String(), htmlBody.String())
}

// SendShareLinkEmail sends an email to a user that doesn't exist yet
func SendShareLinkEmail(
	ctx context.Context,
	to,
	invitedByName,
	customMessage,
	documentId,
	documentTitle,
	accessCode string,
) error {
	c := env.SES(ctx)
	from := fmt.Sprintf("Pointy <pointy@%s>", c.EmailDomain())
	subject := fmt.Sprintf("%s has shared a '%s' with you!", invitedByName, documentTitle)
	preheader := "A Pointy document has been shared with you"
	link := fmt.Sprintf("%s/invite/%s", c.AppHost(), accessCode)

	if c.Env() == "development" {
		log.Infof("DEV share link:\n\n %s\n\n for %s", link, to)
	}

	log.Infof("sending share link to %s", to)

	rctx := c.AttachHostValues(ctx)

	htmlBody := &strings.Builder{}
	err := templates.ShareDocHTML(preheader, documentTitle, invitedByName, customMessage, "Join Document", link).Render(rctx, htmlBody)
	if err != nil {
		return fmt.Errorf("failed to render share doc html: %w", err)
	}

	textBody := &strings.Builder{}
	err = templates.ShareDocText(documentTitle, invitedByName, customMessage, "Join Document", link).Render(rctx, textBody)
	if err != nil {
		return fmt.Errorf("failed to render share doc text: %w", err)
	}

	return c.EnqueueEmail(from, to, subject, textBody.String(), htmlBody.String())
}
