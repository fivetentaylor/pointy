CREATE TABLE public.subscription_plans (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying(255) NOT NULL,
    price_cents integer NOT NULL,
    currency character varying(3) DEFAULT 'USD' NOT NULL,
    interval character varying(20) NOT NULL,
    status character varying(20) DEFAULT 'active',
    stripe_price_id character varying(255) NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE public.user_subscriptions (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    subscription_plan_id uuid NOT NULL,
    stripe_subscription_id character varying(255) NOT NULL,
    status character varying(20) NOT NULL,
    current_period_start timestamp with time zone NOT NULL,
    current_period_end timestamp with time zone NOT NULL,
    canceled_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT fk_user_subscription_user FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_subscription_plan FOREIGN KEY (subscription_plan_id) REFERENCES public.subscription_plans(id)
);

CREATE TABLE public.payment_history (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL, 
    stripe_payment_intent_id character varying(255) NOT NULL,
    amount_cents integer NOT NULL,
    currency character varying(3) DEFAULT 'USD' NOT NULL,
    status character varying(20) NOT NULL, 
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT fk_payment_history_user FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE
);

ALTER TABLE public.users
ADD COLUMN stripe_customer_id character varying(255);

CREATE TABLE public.stripe_webhook_events (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    event_id character varying(255) NOT NULL,
    event_type character varying(255) NOT NULL,
    payload jsonb NOT NULL,
    received_at timestamp with time zone DEFAULT now() NOT NULL,
    processed boolean DEFAULT false NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT unique_event_id UNIQUE (event_id)
);

CREATE INDEX idx_user_subscriptions_user_id ON public.user_subscriptions (user_id);
CREATE INDEX idx_user_subscriptions_plan_id ON public.user_subscriptions (subscription_plan_id);
CREATE INDEX idx_payment_history_user_id ON public.payment_history (user_id);
CREATE INDEX idx_stripe_webhook_events_event_id ON public.stripe_webhook_events (event_id);
