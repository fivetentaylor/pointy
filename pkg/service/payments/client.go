package payments

import (
	"os"

	"github.com/stripe/stripe-go/v81/client"
)

func StripeClient() *client.API {
	sc := &client.API{}
	sc.Init(os.Getenv("STRIPE_API_KEY"), nil)
	return sc
}
