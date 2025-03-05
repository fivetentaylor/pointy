package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/webhook"
	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/models"
	"gorm.io/gorm"
)

var endpointSecret = os.Getenv("STRIPE_WEBHOOK_SECRET")

func (s *Server) StripeWebhook(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log := env.SLog(ctx)
	const MaxBodyBytes = int64(65536)
	req.Body = http.MaxBytesReader(w, req.Body, MaxBodyBytes)
	payload, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	event, err := webhook.ConstructEvent(payload, req.Header.Get("Stripe-Signature"),
		endpointSecret)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Info("👾 Stripe webhook", "event", event.Type)

	stripeWebhookEventsTbl := env.Query(ctx).StripeWebhookEvent

	record := &models.StripeWebhookEvent{
		EventID:    event.ID,
		EventType:  string(event.Type),
		Payload:    string(payload),
		ReceivedAt: time.Now(),
		Processed:  false,
	}
	err = stripeWebhookEventsTbl.Save(record)
	if err != nil {
		log.Error("error inserting webhook event", "error", err)
	}

	// Unmarshal the event data into an appropriate struct depending on its Type
	switch event.Type {
	case "customer.subscription.created", "customer.subscription.updated", "customer.subscription.deleted":
		err = handleSubscription(ctx, record, event)
		if err != nil {
			log.Error("error handling subscription", "error", err)
			fmt.Fprintf(os.Stderr, "Error handling subscription: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	// case "invoice.paid":
	// 	err = handleInvoice(ctx, record, event)
	// 	if err != nil {
	// 		log.Error("error handling invoice", "error", err)
	// 		fmt.Fprintf(os.Stderr, "Error handling invoice: %v\n", err)
	// 		w.WriteHeader(http.StatusBadRequest)
	// 		return
	// 	}
	default:
		log.Info("Unhandled webhook event type", "event", event.Type)
	}

	err = stripeWebhookEventsTbl.Save(record)
	if err != nil {
		log.Error("error inserting webhook event", "error", err)
	}

	w.WriteHeader(http.StatusOK)
}

func handleSubscription(ctx context.Context, record *models.StripeWebhookEvent, event stripe.Event) error {
	log := env.SLog(ctx)
	userTbl := env.Query(ctx).User
	userSubTbl := env.Query(ctx).UserSubscription
	subPlansTbl := env.Query(ctx).SubscriptionPlan

	var sub stripe.Subscription
	err := json.Unmarshal(event.Data.Raw, &sub)
	if err != nil {
		log.Error("Error parsing webhook JSON", "error", err)
		return fmt.Errorf("[server.handleSubscription] Error parsing webhook JSON: %v", err)
	}

	log.Info("🟢 Subscription", "subscription", sub)

	subscription, err := userSubTbl.Where(userSubTbl.StripeSubscriptionID.Eq(sub.ID)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("[server.handleSubscription] error getting user: %s", err)
	}

	if subscription == nil {
		subscription = &models.UserSubscription{}
	}

	subscription.StripeSubscriptionID = sub.ID
	subscription.Status = string(sub.Status)
	subscription.CurrentPeriodStart = time.Unix(sub.CurrentPeriodStart, 0)
	subscription.CurrentPeriodEnd = time.Unix(sub.CurrentPeriodEnd, 0)
	subscription.CanceledAt = time.Unix(sub.CanceledAt, 0)

	for _, item := range sub.Items.Data {
		subPlan, err := subPlansTbl.Where(subPlansTbl.StripePriceID.Eq(item.Price.ID)).First()
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("[server.handleSubscription] error getting subscription plan", "error", err)
		} else {
			subscription.SubscriptionPlanID = subPlan.ID
		}
	}

	var user *models.User
	if user, err = userTbl.Where(userTbl.StripeCustomerID.Eq(sub.Customer.ID)).First(); err != nil {
		log.Error("[server.handleSubscription] error getting user", "error", err)
	} else {
		subscription.UserID = user.ID
	}

	err = userSubTbl.Save(subscription)
	if err != nil {
		return fmt.Errorf("[server.handleSubscription] error saving subscription: %s", err)
	}

	record.Processed = true

	return nil
}
