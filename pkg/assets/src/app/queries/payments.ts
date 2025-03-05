import { gql } from "@/__generated__";

export const SubscriptionPlans = gql(`
  query SubscriptionPlans {
    subscriptionPlans {
      id
      name
      priceCents
      currency
      interval
    }
  }
`);

export const Checkout = gql(`
  mutation Checkout($id: ID!) {
    checkoutSubscriptionPlan(id: $id) {
      url
    }
  }
`);

export const BillingPortalSession = gql(`
  mutation BillingPortalSession {
    billingPortalSession {
      url
    }
  }
`);
