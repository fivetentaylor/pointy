package client

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/charmbracelet/log"
	"github.com/jpoz/conveyor"
	"github.com/teamreviso/code/pkg/background/wire"
	"github.com/teamreviso/code/pkg/constants"
)

type SESInterface interface {
	Env() string
	AppHost() string
	WebHost() string
	EmailDomain() string
	SendRawEmail(from, to, subject, txtbody, htmlbody string) error
	EnqueueEmail(from, to, subject, txtbody, htmlbody string) error
	AttachHostValues(ctx context.Context) context.Context
}

type SES struct {
	appHost     string
	webHost     string
	emailDomain string

	env            string
	svc            *ses.SES
	bg             *conveyor.Client
	disabled       bool
	skipBackground bool
}

const DefaultEmailRegion = "us-east-1"

func NewSESFromEnv(bg *conveyor.Client) (*SES, error) {
	emailRegion := os.Getenv("EMAIL_REGION")
	if emailRegion == "" {
		emailRegion = DefaultEmailRegion
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(emailRegion),
	})
	if err != nil {
		return nil, err
	}

	// For development purposes, use credentials from ~/.aws/credentials
	// This allows email to be authenticated with AWS SES, yet all other AWS serveice
	// use the AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY in dev mode
	if os.Getenv("ENV") == "development" {
		log.Info("[development] using credentials from ~/.aws/credentials")
		creds := credentials.NewSharedCredentials("", "default")
		sess, err = session.NewSession(&aws.Config{
			Credentials: creds,
			Region:      aws.String(emailRegion),
		})
		if err != nil {
			return nil, err
		}
	}

	svc := ses.New(sess)

	emailDomain := os.Getenv("EMAIL_DOMAIN")
	if emailDomain == "" {
		emailDomain = "revi.so"
	}

	disableSendEmail := os.Getenv("DISABLE_EMAIL_SEND") == "true"
	skipBackground := os.Getenv("EMAIL_SEND_SKIP_BACKGROUND") == "true"

	log.Infof("[ses] email domain: %s, region: %s, disableSendEmail: %t, skipBackground: %t", emailDomain, emailRegion, disableSendEmail, skipBackground)

	return &SES{
		emailDomain:    emailDomain,
		env:            os.Getenv("ENV"),
		appHost:        os.Getenv("APP_HOST"),
		webHost:        os.Getenv("WEB_HOST"),
		svc:            svc,
		bg:             bg,
		disabled:       disableSendEmail,
		skipBackground: skipBackground,
	}, nil
}

func (c *SES) Env() string         { return c.env }
func (c *SES) AppHost() string     { return c.appHost }
func (c *SES) WebHost() string     { return c.webHost }
func (c *SES) EmailDomain() string { return c.emailDomain }

func (c *SES) EnqueueEmail(
	from,
	to,
	subject,
	txtbody,
	htmlbody string,
) error {
	if c.disabled {
		log.Warn("email sending is disabled. Would have sent:", "from", from, "to", to, "subject", subject, "txtbody", txtbody, "htmlbody", htmlbody)
		return nil
	}

	if c.Env() == "development" {
		log.Infof("📧 DEV email:\n\nfrom: %s\nto: %s\nsubject: %s\ntxtbody: %s\nhtmlbody: %s\n\n", from, to, subject, txtbody, htmlbody)
	}

	if c.skipBackground {
		return c.SendRawEmail(from, to, subject, txtbody, htmlbody)
	}

	_, err := c.bg.Enqueue(context.Background(), &wire.SendEmail{
		From:     from,
		To:       to,
		Subject:  subject,
		Txtbody:  txtbody,
		Htmlbody: htmlbody,
	})

	return err
}

func (c *SES) SendRawEmail(
	from,
	to,
	subject,
	txtbody,
	htmlbody string,
) error {
	params := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(to),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Data: aws.String(txtbody),
				},
				Html: &ses.Content{
					Data: aws.String(htmlbody),
				},
			},
			Subject: &ses.Content{
				Data: aws.String(subject),
			},
		},
		Source: aws.String(from),
	}

	log.Infof("sending %q email to %s", subject, to)
	output, err := c.svc.SendEmail(params)
	if err != nil {
		log.Errorf("failed to send email to %s: %s", to, err)
		return err
	}

	log.Infof("sent email: %s", aws.StringValue(output.MessageId))

	return nil
}

func (c *SES) AttachHostValues(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, constants.AppHostContextKey, c.appHost)
	return context.WithValue(ctx, constants.WebHostContextKey, c.webHost)
}
