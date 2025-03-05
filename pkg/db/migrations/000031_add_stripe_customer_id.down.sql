-- Drop indexes
DROP INDEX IF EXISTS public.idx_stripe_webhook_events_event_id;
DROP INDEX IF EXISTS public.idx_payment_history_user_id;
DROP INDEX IF EXISTS public.idx_user_subscriptions_plan_id;
DROP INDEX IF EXISTS public.idx_user_subscriptions_user_id;

-- Drop the stripe_webhook_events table
DROP TABLE IF EXISTS public.stripe_webhook_events;

-- Remove added columns from the users table
ALTER TABLE public.users
DROP COLUMN IF EXISTS stripe_customer_id;

-- Drop the payment_history table
DROP TABLE IF EXISTS public.payment_history;

-- Drop the user_subscriptions table
DROP TABLE IF EXISTS public.user_subscriptions;

-- Drop the subscription_plans table
DROP TABLE IF EXISTS public.subscription_plans;
