package email

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/service/email/templates"
)

func SendWaitlistEmail(
	ctx context.Context,
	to string,
) error {
	c := env.SES(ctx)
	from := fmt.Sprintf("Reviso <reviso@%s>", c.EmailDomain())
	subject := fmt.Sprintf("You’ve been added to the Reviso waitlist")
	preheader := "You’ve been added to the Reviso waitlist"

	log.Infof("sending waitlist email to %s", to)

	rctx := c.AttachHostValues(ctx)

	htmlBody := &strings.Builder{}
	err := templates.WaitlistHTML(preheader, to).Render(rctx, htmlBody)
	if err != nil {
		return fmt.Errorf("failed to render magic link html: %w", err)
	}

	textBody := &strings.Builder{}
	err = templates.WaitlistText(to).Render(rctx, textBody)
	if err != nil {
		return fmt.Errorf("failed to render magic link text: %w", err)
	}

	return c.EnqueueEmail(from, to, subject, textBody.String(), htmlBody.String())
}
