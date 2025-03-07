package email

import (
	"context"
	"fmt"
	"strings"

	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/models"
	"github.com/fivetentaylor/pointy/pkg/service/email/templates"
)

func SendFirstOpen(
	ctx context.Context,
	to string,
	user *models.User,
	doc *models.Document,
) error {
	c := env.SES(ctx)
	from := fmt.Sprintf("Reviso <reviso@%s>", c.EmailDomain())
	subject := fmt.Sprintf("%s opened %s", user.Name, doc.Title)
	preheader := "A Reviso document has been shared with you"

	rctx := c.AttachHostValues(ctx)

	htmlBody := &strings.Builder{}
	err := templates.FirstOpenHTML(preheader, user, doc).Render(rctx, htmlBody)
	if err != nil {
		return fmt.Errorf("failed to render share doc html: %w", err)
	}

	textBody := &strings.Builder{}
	err = templates.FirstOpenText(user, doc).Render(rctx, textBody)
	if err != nil {
		return fmt.Errorf("failed to render share doc text: %w", err)
	}

	return c.EnqueueEmail(from, to, subject, textBody.String(), htmlBody.String())
}
