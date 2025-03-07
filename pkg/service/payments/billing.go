package payments

import (
	"context"
	"fmt"
	"os"

	"github.com/stripe/stripe-go/v81"
	"github.com/fivetentaylor/pointy/pkg/env"
)

func BillingPortalSession(ctx context.Context, userID string) (string, error) {
	log := env.SLog(ctx)
	userTbl := env.Query(ctx).User

	user, err := userTbl.Where(userTbl.ID.Eq(userID)).First()
	if err != nil {
		return "", fmt.Errorf("[payments.BillingPortalSession] error getting user: %s", err)
	}

	if user.StripeCustomerID == "" {
		return "", fmt.Errorf("[payments.BillingPortalSession] user has no stripe customer id")
	}

	sc := StripeClient()

	appHost := os.Getenv("APP_HOST")
	if appHost == "" {
		return "", fmt.Errorf("APP_HOST not configured")
	}

	returnURL := fmt.Sprintf("%s/drafts", appHost)

	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(user.StripeCustomerID),
		ReturnURL: stripe.String(returnURL),
	}

	s, err := sc.BillingPortalSessions.New(params)
	if err != nil {
		log.Error("[payments.BillingPortalSession] error creating billing portal session", "error", err)
		return "", fmt.Errorf("[payments.BillingPortalSession] error creating billing portal session: %s", err)
	}

	return s.URL, nil
}
