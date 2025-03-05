package payments

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/stripe/stripe-go/v81"
	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/models"
	"gorm.io/gorm"
)

func SyncPlans(ctx context.Context) error {
	log := env.SLog(ctx)
	sc := StripeClient()

	subPlansTbl := env.Query(ctx).SubscriptionPlan

	iter := sc.Prices.List(&stripe.PriceListParams{})

	for iter.Next() {
		price := iter.Price()
		if price == nil {
			continue
		}

		log.Info("🧾 Syncing price", "price", price.Product.Name, "priceID", price.ID, "price", price)

		plan, err := subPlansTbl.Where(subPlansTbl.StripePriceID.Eq(price.ID)).First()
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("error getting subscription plan", "error", err)
			return fmt.Errorf("error getting subscription plan: %s", err)
		}

		product, err := sc.Products.Get(price.Product.ID, nil)
		if err != nil {
			log.Error("error getting product", "error", err)
			return fmt.Errorf("error getting product: %s", err)
		}

		if plan == nil {
			log.Info("🧾 Creating subscription plan", "price", price.Product.Name, "priceID", price.ID)
			plan = &models.SubscriptionPlan{}
		} else {
			log.Info("🧾 Updating subscription plan", "price", price.Product.Name, "priceID", price.ID)
		}

		plan.Name = product.Name
		plan.StripePriceID = price.ID
		plan.Currency = string(price.Currency)
		plan.PriceCents = int32(price.UnitAmount)
		plan.Interval = string(price.Recurring.Interval)
		plan.UpdatedAt = time.Now()

		if price.Active {
			plan.Status = "active"
		} else {
			plan.Status = "inactive"
		}

		err = subPlansTbl.Save(plan)
		if err != nil {
			log.Error("error saving subscription plan", "error", err)
			return fmt.Errorf("error saving subscription plan: %s", err)
		}
	}

	return iter.Err()
}
