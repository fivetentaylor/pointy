package admin

import (
	"net/http"

	"github.com/fivetentaylor/pointy/pkg/admin/templates"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/service/payments"
)

func SubscriptionPlans(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	subscriptionPlansTbl := env.Query(ctx).SubscriptionPlan
	subscriptionPlans, err := subscriptionPlansTbl.Find()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templates.SubscriptionPlans(subscriptionPlans).Render(ctx, w)
}

func SyncSubscriptionPlans(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := payments.SyncPlans(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/payment/subscription/plans", http.StatusSeeOther)
}
