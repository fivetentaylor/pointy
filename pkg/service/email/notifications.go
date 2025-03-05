package email

import (
	"context"
	"fmt"
	"strings"

	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/models"
	"github.com/teamreviso/code/pkg/service/email/templates"
	"github.com/teamreviso/code/pkg/service/rogue"
)

func SendTimelineCommentEmail(
	ctx context.Context,
	document *models.Document,
	message *models.TimelineMessageV1,
	fromUser *models.User,
	to *models.User,
	isMentioned bool,
) error {
	log := env.Log(ctx)
	c := env.SES(ctx)
	q := env.Query(ctx)
	from := fmt.Sprintf("Reviso <reviso@%s>", c.EmailDomain())

	subject := fmt.Sprintf("New comment on %s", document.Title)
	preheader := "New comment on Reviso"
	if isMentioned {
		subject = fmt.Sprintf("Someone mentioned you on %s", document.Title)
		preheader = "Someone mentioned you on Reviso"
	}

	log.Infof("sending mention email to %s", to.Email)

	docTbl := q.Document
	doc, err := docTbl.Where(docTbl.ID.Eq(document.ID)).First()
	if err != nil {
		log.Errorf("error getting document: %s", err)
		return err
	}

	rctx := c.AttachHostValues(ctx)

	var selection string
	if message.SelectionStartId != "" && message.SelectionEndId != "" {
		var err error
		selection, err = rogue.HTMLSelection(ctx, document.ID, message.SelectionStartId, message.SelectionEndId)
		if err != nil {
			log.Errorf("error getting selection: %s", err)
			return err
		}
	}

	data := &templates.TimelineMessageData{
		Document:    doc,
		FromUser:    fromUser,
		Message:     message,
		Selection:   selection,
		IsMentioned: isMentioned,
	}

	htmlBody := &strings.Builder{}
	err = templates.TimelineMessageHTML(preheader, data).Render(rctx, htmlBody)
	if err != nil {
		return fmt.Errorf("failed to render magic link html: %w", err)
	}

	textBody := &strings.Builder{}
	err = templates.TimelineMessageText(data).Render(rctx, textBody)
	if err != nil {
		return fmt.Errorf("failed to render magic link text: %w", err)
	}

	return c.EnqueueEmail(from, to.Email, subject, textBody.String(), htmlBody.String())
}

func SendTimelineMentionShareEmail(
	ctx context.Context,
	document *models.Document,
	message *models.TimelineMessageV1,
	fromUser *models.User,
	to *models.User,
	isMentioned bool,
) error {
	log := env.Log(ctx)
	c := env.SES(ctx)
	q := env.Query(ctx)
	from := fmt.Sprintf("Reviso <reviso@%s>", c.EmailDomain())

	subject := fmt.Sprintf("%s has shared %s with you!", fromUser.DisplayName, document.Title)
	preheader := "Someone shared a document with you on Reviso"

	log.Infof("sending mention share email to %s", to.Email)

	docTbl := q.Document
	doc, err := docTbl.Where(docTbl.ID.Eq(document.ID)).First()
	if err != nil {
		log.Errorf("error getting document: %s", err)
		return err
	}

	rctx := c.AttachHostValues(ctx)

	var selection string
	if message.SelectionStartId != "" && message.SelectionEndId != "" {
		var err error
		selection, err = rogue.HTMLSelection(ctx, document.ID, message.SelectionStartId, message.SelectionEndId)
		if err != nil {
			log.Errorf("error getting selection: %s", err)
			return err
		}
	}

	data := &templates.MentionShareData{
		Document:  doc,
		FromUser:  fromUser,
		Message:   message,
		Selection: selection,
	}

	htmlBody := &strings.Builder{}
	err = templates.MentionShareHTML(preheader, data).Render(rctx, htmlBody)
	if err != nil {
		return fmt.Errorf("failed to render mention share html: %w", err)
	}

	textBody := &strings.Builder{}
	err = templates.MentionShareText(data).Render(rctx, textBody)
	if err != nil {
		return fmt.Errorf("failed to render mention share text: %w", err)
	}

	return c.EnqueueEmail(from, to.Email, subject, textBody.String(), htmlBody.String())
}
