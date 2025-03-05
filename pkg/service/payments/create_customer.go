package payments

import (
	"context"
	"fmt"

	"github.com/stripe/stripe-go/v81"

	"github.com/teamreviso/code/pkg/env"
)

func CreateCustomer(ctx context.Context, userID string) (string, error) {
	log := env.SLog(ctx)
	userTbl := env.Query(ctx).User

	user, err := userTbl.Where(userTbl.ID.Eq(userID)).First()
	if err != nil {
		return "", fmt.Errorf("[payments.CreateCustomer] error getting user: %s", err)
	}

	if user.StripeCustomerID != "" {
		return user.StripeCustomerID, nil
	}

	params := stripe.CustomerParams{
		Name:  stripe.String(user.Name),
		Email: stripe.String(user.Email),
		Params: stripe.Params{
			Metadata: map[string]string{
				"userID": userID,
			},
		},
	}

	sc := StripeClient()
	log.Info("👤 Creating stripe customer for user", "user", user.ID)

	customer, err := sc.Customers.New(&params)
	if err != nil {
		log.Error("[payments.CreateCustomer] error creating customer", "error", err)
		return "", fmt.Errorf("[payments.CreateCustomer] error creating customer: %s", err)
	}

	user.StripeCustomerID = customer.ID
	err = userTbl.Save(user)
	if err != nil {
		log.Error("[payments.CreateCustomer] error updating user", "error", err)
		return "", fmt.Errorf("[payments.CreateCustomer] error updating user: %s", err)
	}

	return user.StripeCustomerID, nil
}
