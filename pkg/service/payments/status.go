package payments

import (
	"context"
	"errors"

	"github.com/stripe/stripe-go/v81"
	"github.com/teamreviso/code/pkg/env"
	"gorm.io/gorm"
)

type PaymentStatusResponse struct {
	Status string `json:"status"`
}

func Status(ctx context.Context, sessionID, userID string) (*PaymentStatusResponse, error) {
	log := env.SLog(ctx)
	response := &PaymentStatusResponse{
		Status: "unknown",
	}

	sc := StripeClient()
	stripeSession, err := sc.CheckoutSessions.Get(sessionID, nil)
	if err != nil {
		response.Status = "error"
		log.Error("error getting checkout session", "error", err)
		return response, err
	}

	if stripeSession == nil {
		response.Status = "error"
		log.Error("error getting checkout session (missing session)")
		return response, nil
	}

	userTbl := env.Query(ctx).User
	user, err := userTbl.Where(userTbl.ID.Eq(userID)).First()
	if err != nil {
		response.Status = "error"
		log.Error("error getting user", "error", err)
		return response, err
	}

	if user.StripeCustomerID != stripeSession.Customer.ID {
		response.Status = "error"
		log.Error("error getting user (customer mismatch)")
		return response, nil
	}

	userSubStatus, err := UserSubscriptionStatus(ctx, userID)
	if err != nil {
		response.Status = "error"
		log.Error("error getting user subscription status", "error", err)
		return response, err
	}

	status := "pending"
	if stripeSession.Status == stripe.CheckoutSessionStatusComplete && userSubStatus == "active" {
		status = "complete"
	}

	response.Status = string(status)
	return response, nil
}

func UserSubscriptionStatus(ctx context.Context, userID string) (string, error) {
	log := env.SLog(ctx)
	userSubTbl := env.Query(ctx).UserSubscription

	sub, err := userSubTbl.Where(userSubTbl.UserID.Eq(userID)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error("error getting user subscription", "error", err)
		return "", err
	}

	if sub == nil {
		return "inactive", nil
	}

	return sub.Status, nil
}
