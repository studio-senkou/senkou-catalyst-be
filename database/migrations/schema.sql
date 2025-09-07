SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: categories; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.categories (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    merchant_id character(16) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone
);


--
-- Name: categories_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.categories_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: categories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.categories_id_seq OWNED BY public.categories.id;


--
-- Name: merchants; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.merchants (
    id character(16) NOT NULL,
    name character varying(100) NOT NULL,
    username character varying(100) NOT NULL,
    owner_id integer,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone
);


--
-- Name: payment_transactions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.payment_transactions (
    id uuid NOT NULL,
    payment_type character varying(50) NOT NULL,
    payment_channel character varying(50) NOT NULL,
    fraud_status character varying(20),
    amount numeric(15,2) NOT NULL,
    currency character varying(10) NOT NULL,
    status character varying(20) NOT NULL,
    transaction_id character varying(100),
    transaction_time timestamp without time zone,
    signature_key character varying(255),
    expired_at timestamp without time zone,
    settled_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone
);


--
-- Name: product_metrics; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.product_metrics (
    id integer NOT NULL,
    product_id uuid NOT NULL,
    origin character varying(20) NOT NULL,
    ua_browser text,
    ua_os text,
    interaction text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone
);


--
-- Name: product_metrics_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.product_metrics_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: product_metrics_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.product_metrics_id_seq OWNED BY public.product_metrics.id;


--
-- Name: products; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.products (
    id uuid NOT NULL,
    merchant_id character(16) NOT NULL,
    category_id integer,
    title character varying(150) NOT NULL,
    price character varying(30) NOT NULL,
    description text,
    affiliate_url text NOT NULL,
    photos json,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone
);


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version character varying NOT NULL
);


--
-- Name: subscription_orders; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.subscription_orders (
    id uuid NOT NULL,
    user_id integer NOT NULL,
    subscription_id integer NOT NULL,
    payment_transaction_id uuid NOT NULL,
    amount numeric(10,2) NOT NULL,
    status character varying(50) DEFAULT 'pending'::character varying,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone
);


--
-- Name: subscription_plans; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.subscription_plans (
    id integer NOT NULL,
    sub_id integer NOT NULL,
    name character varying(100) NOT NULL,
    value text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone
);


--
-- Name: subscription_plans_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.subscription_plans_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: subscription_plans_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.subscription_plans_id_seq OWNED BY public.subscription_plans.id;


--
-- Name: subscriptions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.subscriptions (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    price numeric(10,2) NOT NULL,
    description text,
    duration smallint NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone
);


--
-- Name: subscriptions_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.subscriptions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: subscriptions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.subscriptions_id_seq OWNED BY public.subscriptions.id;


--
-- Name: user_subscriptions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_subscriptions (
    id integer NOT NULL,
    user_id integer NOT NULL,
    sub_id integer NOT NULL,
    started_at timestamp without time zone NOT NULL,
    expired_at timestamp without time zone NOT NULL,
    is_active boolean DEFAULT false,
    payment_status character varying(50) DEFAULT 'pending'::character varying,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone
);


--
-- Name: user_subscriptions_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.user_subscriptions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: user_subscriptions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.user_subscriptions_id_seq OWNED BY public.user_subscriptions.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    email character varying(100) NOT NULL,
    phone character varying(15),
    password character varying(255) NOT NULL,
    role character varying(50) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone
);


--
-- Name: users_has_token; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users_has_token (
    id integer NOT NULL,
    user_id integer NOT NULL,
    token character varying(255) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone
);


--
-- Name: users_has_token_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.users_has_token_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: users_has_token_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.users_has_token_id_seq OWNED BY public.users_has_token.id;


--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: categories id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.categories ALTER COLUMN id SET DEFAULT nextval('public.categories_id_seq'::regclass);


--
-- Name: product_metrics id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_metrics ALTER COLUMN id SET DEFAULT nextval('public.product_metrics_id_seq'::regclass);


--
-- Name: subscription_plans id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.subscription_plans ALTER COLUMN id SET DEFAULT nextval('public.subscription_plans_id_seq'::regclass);


--
-- Name: subscriptions id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.subscriptions ALTER COLUMN id SET DEFAULT nextval('public.subscriptions_id_seq'::regclass);


--
-- Name: user_subscriptions id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_subscriptions ALTER COLUMN id SET DEFAULT nextval('public.user_subscriptions_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: users_has_token id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users_has_token ALTER COLUMN id SET DEFAULT nextval('public.users_has_token_id_seq'::regclass);


--
-- Name: categories categories_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (id);


--
-- Name: merchants merchants_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.merchants
    ADD CONSTRAINT merchants_pkey PRIMARY KEY (id);


--
-- Name: merchants merchants_username_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.merchants
    ADD CONSTRAINT merchants_username_key UNIQUE (username);


--
-- Name: payment_transactions payment_transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payment_transactions
    ADD CONSTRAINT payment_transactions_pkey PRIMARY KEY (id);


--
-- Name: payment_transactions payment_transactions_transaction_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payment_transactions
    ADD CONSTRAINT payment_transactions_transaction_id_key UNIQUE (transaction_id);


--
-- Name: product_metrics product_metrics_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_metrics
    ADD CONSTRAINT product_metrics_pkey PRIMARY KEY (id);


--
-- Name: products products_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: subscription_orders subscription_orders_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.subscription_orders
    ADD CONSTRAINT subscription_orders_pkey PRIMARY KEY (id);


--
-- Name: subscription_plans subscription_plans_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.subscription_plans
    ADD CONSTRAINT subscription_plans_pkey PRIMARY KEY (id);


--
-- Name: subscriptions subscriptions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.subscriptions
    ADD CONSTRAINT subscriptions_pkey PRIMARY KEY (id);


--
-- Name: user_subscriptions user_subscriptions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_subscriptions
    ADD CONSTRAINT user_subscriptions_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users_has_token users_has_token_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users_has_token
    ADD CONSTRAINT users_has_token_pkey PRIMARY KEY (id);


--
-- Name: users users_phone_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_phone_key UNIQUE (phone);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_categories_merchant_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_categories_merchant_id ON public.categories USING btree (merchant_id);


--
-- Name: idx_merchants_owner_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_merchants_owner_id ON public.merchants USING btree (owner_id);


--
-- Name: idx_product_metrics_product_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_product_metrics_product_id ON public.product_metrics USING btree (product_id);


--
-- Name: idx_products_category_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_products_category_id ON public.products USING btree (category_id);


--
-- Name: idx_products_merchant_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_products_merchant_id ON public.products USING btree (merchant_id);


--
-- Name: idx_subscription_orders_payment_transaction_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_subscription_orders_payment_transaction_id ON public.subscription_orders USING btree (payment_transaction_id);


--
-- Name: idx_subscription_orders_subscription_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_subscription_orders_subscription_id ON public.subscription_orders USING btree (subscription_id);


--
-- Name: idx_subscription_orders_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_subscription_orders_user_id ON public.subscription_orders USING btree (user_id);


--
-- Name: idx_subscription_plans_sub_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_subscription_plans_sub_id ON public.subscription_plans USING btree (sub_id);


--
-- Name: idx_user_subscriptions_sub_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_subscriptions_sub_id ON public.user_subscriptions USING btree (sub_id);


--
-- Name: idx_user_subscriptions_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_subscriptions_user_id ON public.user_subscriptions USING btree (user_id);


--
-- Name: categories fk_categories_merchant; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT fk_categories_merchant FOREIGN KEY (merchant_id) REFERENCES public.merchants(id) ON DELETE CASCADE;


--
-- Name: merchants fk_merchant_owner; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.merchants
    ADD CONSTRAINT fk_merchant_owner FOREIGN KEY (owner_id) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: product_metrics fk_product; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_metrics
    ADD CONSTRAINT fk_product FOREIGN KEY (product_id) REFERENCES public.products(id) ON DELETE CASCADE;


--
-- Name: products fk_products_category; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT fk_products_category FOREIGN KEY (category_id) REFERENCES public.categories(id) ON DELETE SET NULL;


--
-- Name: products fk_products_merchant; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT fk_products_merchant FOREIGN KEY (merchant_id) REFERENCES public.merchants(id) ON DELETE CASCADE;


--
-- Name: subscription_orders fk_subscription_orders_payment_transaction; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.subscription_orders
    ADD CONSTRAINT fk_subscription_orders_payment_transaction FOREIGN KEY (payment_transaction_id) REFERENCES public.payment_transactions(id) ON DELETE CASCADE;


--
-- Name: subscription_orders fk_subscription_orders_subscription; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.subscription_orders
    ADD CONSTRAINT fk_subscription_orders_subscription FOREIGN KEY (subscription_id) REFERENCES public.subscriptions(id) ON DELETE CASCADE;


--
-- Name: subscription_orders fk_subscription_orders_user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.subscription_orders
    ADD CONSTRAINT fk_subscription_orders_user FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: subscription_plans fk_subscription_plans_subscription; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.subscription_plans
    ADD CONSTRAINT fk_subscription_plans_subscription FOREIGN KEY (sub_id) REFERENCES public.subscriptions(id) ON DELETE CASCADE;


--
-- Name: user_subscriptions fk_user_subscriptions_subscription; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_subscriptions
    ADD CONSTRAINT fk_user_subscriptions_subscription FOREIGN KEY (sub_id) REFERENCES public.subscriptions(id) ON DELETE CASCADE;


--
-- Name: user_subscriptions fk_user_subscriptions_user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_subscriptions
    ADD CONSTRAINT fk_user_subscriptions_user FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: users_has_token fk_users_has_token; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users_has_token
    ADD CONSTRAINT fk_users_has_token FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--


--
-- Dbmate schema migrations
--

INSERT INTO public.schema_migrations (version) VALUES
    ('20250906180807'),
    ('20250906181125'),
    ('20250906184723'),
    ('20250906184949'),
    ('20250906190008'),
    ('20250906190502'),
    ('20250906190730'),
    ('20250906190914'),
    ('20250906191222'),
    ('20250906191449'),
    ('20250906191821');
