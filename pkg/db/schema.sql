--
-- PostgreSQL database dump
--

-- Dumped from database version 16.3 (Debian 16.3-1.pgdg120+1)
-- Dumped by pg_dump version 16.4 (Homebrew)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;


--
-- Name: EXTENSION pgcrypto; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION pgcrypto IS 'cryptographic functions';


--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


--
-- Name: generate_access_link(integer); Type: FUNCTION; Schema: public; Owner: dev
--

CREATE FUNCTION public.generate_access_link(length integer) RETURNS character varying
    LANGUAGE plpgsql
    AS $$
DECLARE
  chars VARCHAR := 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
  result VARCHAR := '';
  i INT;
BEGIN
  FOR i IN 1..length LOOP
    result := result || substr(chars, floor(random() * length(chars) + 1)::INT, 1);
  END LOOP;
  RETURN result;
END;
$$;


ALTER FUNCTION public.generate_access_link(length integer) OWNER TO dev;

--
-- Name: generate_default_access_token(); Type: FUNCTION; Schema: public; Owner: dev
--

CREATE FUNCTION public.generate_default_access_token() RETURNS character varying
    LANGUAGE plpgsql
    AS $$
BEGIN
  RETURN generate_access_link(128); -- default length of 128
END;
$$;


ALTER FUNCTION public.generate_default_access_token() OWNER TO dev;

--
-- Name: generate_default_invite_link(); Type: FUNCTION; Schema: public; Owner: dev
--

CREATE FUNCTION public.generate_default_invite_link() RETURNS character varying
    LANGUAGE plpgsql
    AS $$
BEGIN
  RETURN generate_access_link(8); -- default length of 8
END;
$$;


ALTER FUNCTION public.generate_default_invite_link() OWNER TO dev;

--
-- Name: increment_author_id(uuid); Type: FUNCTION; Schema: public; Owner: dev
--

CREATE FUNCTION public.increment_author_id(doc_id uuid) RETURNS integer
    LANGUAGE plpgsql
    AS $$
DECLARE
    new_author_id INT;
BEGIN
    SELECT COALESCE(MAX(author_id), 0) + 1 INTO new_author_id
    FROM author_ids WHERE document_id = doc_id;

    RETURN new_author_id;
END;
$$;


ALTER FUNCTION public.increment_author_id(doc_id uuid) OWNER TO dev;

--
-- Name: set_default_root_parent_id(); Type: FUNCTION; Schema: public; Owner: dev
--

CREATE FUNCTION public.set_default_root_parent_id() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF NEW.root_parent_id IS NULL THEN
       NEW.root_parent_id := NEW.id;
    END IF;
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.set_default_root_parent_id() OWNER TO dev;

--
-- Name: update_updated_at_column(); Type: FUNCTION; Schema: public; Owner: dev
--

CREATE FUNCTION public.update_updated_at_column() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$;


ALTER FUNCTION public.update_updated_at_column() OWNER TO dev;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: author_ids; Type: TABLE; Schema: public; Owner: dev
--

CREATE TABLE public.author_ids (
    id integer NOT NULL,
    author_id integer NOT NULL,
    document_id uuid NOT NULL,
    user_id uuid NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.author_ids OWNER TO dev;

--
-- Name: author_ids_id_seq; Type: SEQUENCE; Schema: public; Owner: dev
--

CREATE SEQUENCE public.author_ids_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.author_ids_id_seq OWNER TO dev;

--
-- Name: author_ids_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: dev
--

ALTER SEQUENCE public.author_ids_id_seq OWNED BY public.author_ids.id;


--
-- Name: comments; Type: TABLE; Schema: public; Owner: dev
--

CREATE TABLE public.comments (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    document_id uuid NOT NULL,
    thread_id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,
    user_id uuid NOT NULL,
    body text NOT NULL,
    selection text,
    selection_start text,
    selection_finish text,
    notes text
);


ALTER TABLE public.comments OWNER TO dev;

--
-- Name: default_documents; Type: TABLE; Schema: public; Owner: dev
--

CREATE TABLE public.default_documents (
    name character varying(255) NOT NULL,
    doc_id character varying(255) NOT NULL
);


ALTER TABLE public.default_documents OWNER TO dev;

--
-- Name: document_access; Type: TABLE; Schema: public; Owner: dev
--

CREATE TABLE public.document_access (
    id integer NOT NULL,
    document_id uuid NOT NULL,
    user_id uuid NOT NULL,
    access_level text NOT NULL,
    last_accessed_at timestamp with time zone,
    CONSTRAINT document_access_access_level_check CHECK ((access_level = ANY (ARRAY['comment'::text, 'write'::text, 'owner'::text, 'admin'::text])))
);


ALTER TABLE public.document_access OWNER TO dev;

--
-- Name: document_access_id_seq; Type: SEQUENCE; Schema: public; Owner: dev
--

CREATE SEQUENCE public.document_access_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.document_access_id_seq OWNER TO dev;

--
-- Name: document_access_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: dev
--

ALTER SEQUENCE public.document_access_id_seq OWNED BY public.document_access.id;


--
-- Name: document_attachments; Type: TABLE; Schema: public; Owner: dev
--

CREATE TABLE public.document_attachments (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    document_id uuid NOT NULL,
    s3_id uuid NOT NULL,
    filename text NOT NULL,
    content_type text NOT NULL,
    size bigint NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.document_attachments OWNER TO dev;

--
-- Name: document_versions; Type: TABLE; Schema: public; Owner: dev
--

CREATE TABLE public.document_versions (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    document_id uuid NOT NULL,
    name text NOT NULL,
    content_address text NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_by uuid NOT NULL,
    updated_by uuid NOT NULL
);


ALTER TABLE public.document_versions OWNER TO dev;

--
-- Name: documents; Type: TABLE; Schema: public; Owner: dev
--

CREATE TABLE public.documents (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    title character varying(255) NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    parent_id uuid,
    deleted_at timestamp with time zone,
    is_public boolean DEFAULT true NOT NULL,
    rogue_version character varying(255),
    root_parent_id uuid NOT NULL,
    parent_address text,
    is_folder boolean DEFAULT false,
    folder_id uuid
);


ALTER TABLE public.documents OWNER TO dev;

--
-- Name: one_time_access_tokens; Type: TABLE; Schema: public; Owner: dev
--

CREATE TABLE public.one_time_access_tokens (
    id integer NOT NULL,
    user_id uuid NOT NULL,
    token character varying(128) DEFAULT public.generate_default_access_token() NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    is_used boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.one_time_access_tokens OWNER TO dev;

--
-- Name: one_time_access_tokens_id_seq; Type: SEQUENCE; Schema: public; Owner: dev
--

CREATE SEQUENCE public.one_time_access_tokens_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.one_time_access_tokens_id_seq OWNER TO dev;

--
-- Name: one_time_access_tokens_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: dev
--

ALTER SEQUENCE public.one_time_access_tokens_id_seq OWNED BY public.one_time_access_tokens.id;


--
-- Name: payment_history; Type: TABLE; Schema: public; Owner: dev
--

CREATE TABLE public.payment_history (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    stripe_payment_intent_id character varying(255) NOT NULL,
    amount_cents integer NOT NULL,
    currency character varying(3) DEFAULT 'USD'::character varying NOT NULL,
    status character varying(20) NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.payment_history OWNER TO dev;

--
-- Name: prompts; Type: TABLE; Schema: public; Owner: dev
--

CREATE TABLE public.prompts (
    id integer NOT NULL,
    prompt_name character varying(255) NOT NULL,
    system_content text,
    version text DEFAULT ''::text NOT NULL,
    provider character varying(255) NOT NULL,
    model_name character varying(255) NOT NULL,
    temperature double precision DEFAULT 1.0,
    max_tokens integer DEFAULT 0,
    top_p double precision DEFAULT 1.0,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    content_json json
);


ALTER TABLE public.prompts OWNER TO dev;

--
-- Name: prompts_id_seq; Type: SEQUENCE; Schema: public; Owner: dev
--

CREATE SEQUENCE public.prompts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.prompts_id_seq OWNER TO dev;

--
-- Name: prompts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: dev
--

ALTER SEQUENCE public.prompts_id_seq OWNED BY public.prompts.id;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: dev
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO dev;

--
-- Name: shared_document_links; Type: TABLE; Schema: public; Owner: dev
--

CREATE TABLE public.shared_document_links (
    id integer NOT NULL,
    document_id uuid NOT NULL,
    inviter_id uuid NOT NULL,
    invitee_email character varying(255) NOT NULL,
    invite_link character varying(8) DEFAULT public.generate_default_invite_link() NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    is_active boolean DEFAULT true
);


ALTER TABLE public.shared_document_links OWNER TO dev;

--
-- Name: shared_document_links_id_seq; Type: SEQUENCE; Schema: public; Owner: dev
--

CREATE SEQUENCE public.shared_document_links_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.shared_document_links_id_seq OWNER TO dev;

--
-- Name: shared_document_links_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: dev
--

ALTER SEQUENCE public.shared_document_links_id_seq OWNED BY public.shared_document_links.id;


--
-- Name: stripe_webhook_events; Type: TABLE; Schema: public; Owner: dev
--

CREATE TABLE public.stripe_webhook_events (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    event_id character varying(255) NOT NULL,
    event_type character varying(255) NOT NULL,
    payload jsonb NOT NULL,
    received_at timestamp with time zone DEFAULT now() NOT NULL,
    processed boolean DEFAULT false NOT NULL
);


ALTER TABLE public.stripe_webhook_events OWNER TO dev;

--
-- Name: subscription_plans; Type: TABLE; Schema: public; Owner: dev
--

CREATE TABLE public.subscription_plans (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying(255) NOT NULL,
    price_cents integer NOT NULL,
    currency character varying(3) DEFAULT 'USD'::character varying NOT NULL,
    "interval" character varying(20) NOT NULL,
    status character varying(20) DEFAULT 'active'::character varying,
    stripe_price_id character varying(255) NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.subscription_plans OWNER TO dev;

--
-- Name: user_subscriptions; Type: TABLE; Schema: public; Owner: dev
--

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
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.user_subscriptions OWNER TO dev;

--
-- Name: users; Type: TABLE; Schema: public; Owner: dev
--

CREATE TABLE public.users (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    password_hash character varying(60),
    provider text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    picture text,
    admin boolean DEFAULT false NOT NULL,
    display_name character varying(255) NOT NULL,
    educator boolean DEFAULT false NOT NULL,
    stripe_customer_id character varying(255)
);


ALTER TABLE public.users OWNER TO dev;

--
-- Name: waitlist_users; Type: TABLE; Schema: public; Owner: dev
--

CREATE TABLE public.waitlist_users (
    email character varying(255) NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    allow_access boolean DEFAULT false NOT NULL
);


ALTER TABLE public.waitlist_users OWNER TO dev;

--
-- Name: author_ids id; Type: DEFAULT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.author_ids ALTER COLUMN id SET DEFAULT nextval('public.author_ids_id_seq'::regclass);


--
-- Name: document_access id; Type: DEFAULT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.document_access ALTER COLUMN id SET DEFAULT nextval('public.document_access_id_seq'::regclass);


--
-- Name: one_time_access_tokens id; Type: DEFAULT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.one_time_access_tokens ALTER COLUMN id SET DEFAULT nextval('public.one_time_access_tokens_id_seq'::regclass);


--
-- Name: prompts id; Type: DEFAULT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.prompts ALTER COLUMN id SET DEFAULT nextval('public.prompts_id_seq'::regclass);


--
-- Name: shared_document_links id; Type: DEFAULT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.shared_document_links ALTER COLUMN id SET DEFAULT nextval('public.shared_document_links_id_seq'::regclass);


--
-- Name: author_ids author_ids_author_id_document_id_key; Type: CONSTRAINT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.author_ids
    ADD CONSTRAINT author_ids_author_id_document_id_key UNIQUE (author_id, document_id);


--
-- Name: author_ids author_ids_pkey; Type: CONSTRAINT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.author_ids
    ADD CONSTRAINT author_ids_pkey PRIMARY KEY (id);


--
-- Name: comments comments_pkey; Type: CONSTRAINT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.comments
    ADD CONSTRAINT comments_pkey PRIMARY KEY (id);


--
-- Name: default_documents default_documents_pkey; Type: CONSTRAINT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.default_documents
    ADD CONSTRAINT default_documents_pkey PRIMARY KEY (name);


--
-- Name: document_access document_access_document_id_user_id_key; Type: CONSTRAINT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.document_access
    ADD CONSTRAINT document_access_document_id_user_id_key UNIQUE (document_id, user_id);


--
-- Name: document_access document_access_pkey; Type: CONSTRAINT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.document_access
    ADD CONSTRAINT document_access_pkey PRIMARY KEY (id);


--
-- Name: document_attachments document_attachments_pkey; Type: CONSTRAINT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.document_attachments
    ADD CONSTRAINT document_attachments_pkey PRIMARY KEY (id);


--
-- Name: document_versions document_versions_pkey; Type: CONSTRAINT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.document_versions
    ADD CONSTRAINT document_versions_pkey PRIMARY KEY (id);


--
-- Name: documents documents_pkey; Type: CONSTRAINT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.documents
    ADD CONSTRAINT documents_pkey PRIMARY KEY (id);


--
-- Name: one_time_access_tokens one_time_access_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.one_time_access_tokens
    ADD CONSTRAINT one_time_access_tokens_pkey PRIMARY KEY (id);


--
-- Name: one_time_access_tokens one_time_access_tokens_token_key; Type: CONSTRAINT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.one_time_access_tokens
    ADD CONSTRAINT one_time_access_tokens_token_key UNIQUE (token);


--
-- Name: payment_history payment_history_pkey; Type: CONSTRAINT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.payment_history
    ADD CONSTRAINT payment_history_pkey PRIMARY KEY (id);


--
-- Name: prompts prompts_pkey; Type: CONSTRAINT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.prompts
    ADD CONSTRAINT prompts_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: shared_document_links shared_document_links_invite_link_key; Type: CONSTRAINT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.shared_document_links
    ADD CONSTRAINT shared_document_links_invite_link_key UNIQUE (invite_link);


--
-- Name: shared_document_links shared_document_links_pkey; Type: CONSTRAINT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.shared_document_links
    ADD CONSTRAINT shared_document_links_pkey PRIMARY KEY (id);


--
-- Name: stripe_webhook_events stripe_webhook_events_pkey; Type: CONSTRAINT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.stripe_webhook_events
    ADD CONSTRAINT stripe_webhook_events_pkey PRIMARY KEY (id);


--
-- Name: subscription_plans subscription_plans_pkey; Type: CONSTRAINT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.subscription_plans
    ADD CONSTRAINT subscription_plans_pkey PRIMARY KEY (id);


--
-- Name: users unique_email; Type: CONSTRAINT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT unique_email UNIQUE (email);


--
-- Name: stripe_webhook_events unique_event_id; Type: CONSTRAINT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.stripe_webhook_events
    ADD CONSTRAINT unique_event_id UNIQUE (event_id);


--
-- Name: user_subscriptions user_subscriptions_pkey; Type: CONSTRAINT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.user_subscriptions
    ADD CONSTRAINT user_subscriptions_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: waitlist_users waitlist_users_email_key; Type: CONSTRAINT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.waitlist_users
    ADD CONSTRAINT waitlist_users_email_key UNIQUE (email);


--
-- Name: idx_author_document; Type: INDEX; Schema: public; Owner: dev
--

CREATE INDEX idx_author_document ON public.author_ids USING btree (author_id, document_id);


--
-- Name: idx_comments_document_id; Type: INDEX; Schema: public; Owner: dev
--

CREATE INDEX idx_comments_document_id ON public.comments USING btree (document_id);


--
-- Name: idx_comments_user_thread; Type: INDEX; Schema: public; Owner: dev
--

CREATE INDEX idx_comments_user_thread ON public.comments USING btree (user_id, thread_id);


--
-- Name: idx_document_attachments_document_id; Type: INDEX; Schema: public; Owner: dev
--

CREATE INDEX idx_document_attachments_document_id ON public.document_attachments USING btree (document_id);


--
-- Name: idx_document_attachments_user_id; Type: INDEX; Schema: public; Owner: dev
--

CREATE INDEX idx_document_attachments_user_id ON public.document_attachments USING btree (user_id);


--
-- Name: idx_document_versions_document_id; Type: INDEX; Schema: public; Owner: dev
--

CREATE INDEX idx_document_versions_document_id ON public.document_versions USING btree (document_id);


--
-- Name: idx_documents_folder_id; Type: INDEX; Schema: public; Owner: dev
--

CREATE INDEX idx_documents_folder_id ON public.documents USING btree (folder_id);


--
-- Name: idx_documents_is_folder; Type: INDEX; Schema: public; Owner: dev
--

CREATE INDEX idx_documents_is_folder ON public.documents USING btree (is_folder);


--
-- Name: idx_documents_parent_folder; Type: INDEX; Schema: public; Owner: dev
--

CREATE INDEX idx_documents_parent_folder ON public.documents USING btree (parent_id) WHERE (is_folder = true);


--
-- Name: idx_invite_link; Type: INDEX; Schema: public; Owner: dev
--

CREATE UNIQUE INDEX idx_invite_link ON public.shared_document_links USING btree (invite_link);


--
-- Name: idx_payment_history_user_id; Type: INDEX; Schema: public; Owner: dev
--

CREATE INDEX idx_payment_history_user_id ON public.payment_history USING btree (user_id);


--
-- Name: idx_prompts_prompt_name; Type: INDEX; Schema: public; Owner: dev
--

CREATE INDEX idx_prompts_prompt_name ON public.prompts USING btree (prompt_name);


--
-- Name: idx_stripe_webhook_events_event_id; Type: INDEX; Schema: public; Owner: dev
--

CREATE INDEX idx_stripe_webhook_events_event_id ON public.stripe_webhook_events USING btree (event_id);


--
-- Name: idx_user_subscriptions_plan_id; Type: INDEX; Schema: public; Owner: dev
--

CREATE INDEX idx_user_subscriptions_plan_id ON public.user_subscriptions USING btree (subscription_plan_id);


--
-- Name: idx_user_subscriptions_user_id; Type: INDEX; Schema: public; Owner: dev
--

CREATE INDEX idx_user_subscriptions_user_id ON public.user_subscriptions USING btree (user_id);


--
-- Name: documents set_root_parent_id_trigger; Type: TRIGGER; Schema: public; Owner: dev
--

CREATE TRIGGER set_root_parent_id_trigger BEFORE INSERT ON public.documents FOR EACH ROW EXECUTE FUNCTION public.set_default_root_parent_id();


--
-- Name: document_versions update_document_versions_updated_at; Type: TRIGGER; Schema: public; Owner: dev
--

CREATE TRIGGER update_document_versions_updated_at BEFORE UPDATE ON public.document_versions FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: prompts update_prompts_updated_at; Type: TRIGGER; Schema: public; Owner: dev
--

CREATE TRIGGER update_prompts_updated_at BEFORE UPDATE ON public.prompts FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: comments comments_document_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.comments
    ADD CONSTRAINT comments_document_id_fkey FOREIGN KEY (document_id) REFERENCES public.documents(id);


--
-- Name: comments comments_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.comments
    ADD CONSTRAINT comments_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: document_access document_access_document_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.document_access
    ADD CONSTRAINT document_access_document_id_fkey FOREIGN KEY (document_id) REFERENCES public.documents(id);


--
-- Name: document_access document_access_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.document_access
    ADD CONSTRAINT document_access_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: document_attachments document_attachments_document_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.document_attachments
    ADD CONSTRAINT document_attachments_document_id_fkey FOREIGN KEY (document_id) REFERENCES public.documents(id) ON DELETE CASCADE;


--
-- Name: document_attachments document_attachments_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.document_attachments
    ADD CONSTRAINT document_attachments_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: document_versions document_versions_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.document_versions
    ADD CONSTRAINT document_versions_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id);


--
-- Name: document_versions document_versions_document_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.document_versions
    ADD CONSTRAINT document_versions_document_id_fkey FOREIGN KEY (document_id) REFERENCES public.documents(id) ON DELETE CASCADE;


--
-- Name: document_versions document_versions_updated_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.document_versions
    ADD CONSTRAINT document_versions_updated_by_fkey FOREIGN KEY (updated_by) REFERENCES public.users(id);


--
-- Name: documents documents_folder_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.documents
    ADD CONSTRAINT documents_folder_id_fkey FOREIGN KEY (folder_id) REFERENCES public.documents(id);


--
-- Name: documents fk_parent_id; Type: FK CONSTRAINT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.documents
    ADD CONSTRAINT fk_parent_id FOREIGN KEY (parent_id) REFERENCES public.documents(id);


--
-- Name: payment_history fk_payment_history_user; Type: FK CONSTRAINT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.payment_history
    ADD CONSTRAINT fk_payment_history_user FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: documents fk_root_parent_id; Type: FK CONSTRAINT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.documents
    ADD CONSTRAINT fk_root_parent_id FOREIGN KEY (root_parent_id) REFERENCES public.documents(id);


--
-- Name: user_subscriptions fk_user_subscription_plan; Type: FK CONSTRAINT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.user_subscriptions
    ADD CONSTRAINT fk_user_subscription_plan FOREIGN KEY (subscription_plan_id) REFERENCES public.subscription_plans(id);


--
-- Name: user_subscriptions fk_user_subscription_user; Type: FK CONSTRAINT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.user_subscriptions
    ADD CONSTRAINT fk_user_subscription_user FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: one_time_access_tokens one_time_access_tokens_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.one_time_access_tokens
    ADD CONSTRAINT one_time_access_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: shared_document_links shared_document_links_document_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.shared_document_links
    ADD CONSTRAINT shared_document_links_document_id_fkey FOREIGN KEY (document_id) REFERENCES public.documents(id);


--
-- Name: shared_document_links shared_document_links_inviter_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dev
--

ALTER TABLE ONLY public.shared_document_links
    ADD CONSTRAINT shared_document_links_inviter_id_fkey FOREIGN KEY (inviter_id) REFERENCES public.users(id);


--
-- PostgreSQL database dump complete
--

