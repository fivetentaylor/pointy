package email

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/service/email/templates"
)

func SendMagicLinkEmail(
	ctx context.Context,
	to,
	accessToken,
	next string,
) error {
	c := env.SES(ctx)
	from := fmt.Sprintf("Pointy <pointy@%s>", c.EmailDomain())
	subject := fmt.Sprintf("One time access link")
	preheader := "Your one time access link"

	accessLink := fmt.Sprintf("%s/access/%s", c.AppHost(), accessToken)
	if next != "" {
		accessLink = fmt.Sprintf("%s/access/%s?next=%s", c.AppHost(), accessToken, next)
	}

	if c.WebHost() == "https://www.dev.pointy.ai" {
		log.Infof("DEV access link:\n\n %s\n\n for %s", accessLink, to)
	}

	log.Infof("sending access link to %s", to)

	rctx := c.AttachHostValues(ctx)

	htmlBody := &strings.Builder{}
	err := templates.MagicLinkHTML(preheader, to, accessLink).Render(rctx, htmlBody)
	if err != nil {
		return fmt.Errorf("failed to render magic link html: %w", err)
	}

	textBody := &strings.Builder{}
	err = templates.MagicLinkText(to, accessLink).Render(rctx, textBody)
	if err != nil {
		return fmt.Errorf("failed to render magic link text: %w", err)
	}

	return c.EnqueueEmail(from, to, subject, textBody.String(), htmlBody.String())
}
