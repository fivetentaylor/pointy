package payments

import (
	"context"
	"fmt"
	"os"

	"github.com/stripe/stripe-go/v81"
	"github.com/teamreviso/code/pkg/env"
)

func Checkout(ctx context.Context, userID, subscriptionID string) (string, error) {
	log := env.SLog(ctx)
	customerID, err := CreateCustomer(ctx, userID)
	if err != nil {
		return "", err
	}

	subPlansTbl := env.Query(ctx).SubscriptionPlan
	p, err := subPlansTbl.Where(subPlansTbl.ID.Eq(subscriptionID)).First()
	if err != nil {
		return "", err
	}

	appHost := os.Getenv("APP_HOST")
	if appHost == "" {
		return "", fmt.Errorf("APP_HOST not configured")
	}

	successURL := fmt.Sprintf("%s/payments/success?session_id={CHECKOUT_SESSION_ID}", appHost)
	cancelURL := fmt.Sprintf("%s/payments/cancel", appHost)

	log.Info("🧾 Checkout", "plan", p, "userID", userID, "successURL", successURL, "cancelURL", cancelURL)

	params := &stripe.CheckoutSessionParams{
		Mode:     stripe.String("subscription"),
		Customer: stripe.String(customerID),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(p.StripePriceID),
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String(successURL),
		CancelURL:  stripe.String(cancelURL),
	}

	sc := StripeClient()
	s, err := sc.CheckoutSessions.New(params)
	if err != nil {
		return "", err
	}

	return s.URL, nil
}
