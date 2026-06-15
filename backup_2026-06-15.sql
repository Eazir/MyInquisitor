--
-- PostgreSQL database dump
--

\restrict BeMnSNKvfKY8R0MiIoNqdPr1a5cAuPRafnaIC5JOeqU4OuwcZENbUk1CnOGSq6D

-- Dumped from database version 16.14
-- Dumped by pg_dump version 16.14

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

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: categories; Type: TABLE; Schema: public; Owner: myinquisitor_app
--

CREATE TABLE public.categories (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    name text NOT NULL,
    type text NOT NULL,
    icon text,
    color text,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT categories_type_check CHECK ((type = ANY (ARRAY['income'::text, 'expense'::text, 'debt'::text])))
);


ALTER TABLE public.categories OWNER TO myinquisitor_app;

--
-- Name: debt_monthly_status; Type: TABLE; Schema: public; Owner: myinquisitor_app
--

CREATE TABLE public.debt_monthly_status (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    debt_id uuid NOT NULL,
    month date NOT NULL,
    installment_num integer NOT NULL,
    total_installments integer NOT NULL,
    amount_due numeric(14,2) NOT NULL,
    interest_amount numeric(14,2) DEFAULT 0,
    principal_amount numeric(14,2) DEFAULT 0,
    amount_paid numeric(14,2) DEFAULT 0,
    paid boolean DEFAULT false NOT NULL,
    paid_at timestamp with time zone,
    notes text,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.debt_monthly_status OWNER TO myinquisitor_app;

--
-- Name: debts; Type: TABLE; Schema: public; Owner: myinquisitor_app
--

CREATE TABLE public.debts (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    category_id uuid,
    name text NOT NULL,
    description text,
    total_amount numeric(14,2) NOT NULL,
    remaining_amount numeric(14,2) NOT NULL,
    interest_rate numeric(5,2) DEFAULT 0,
    total_installments integer DEFAULT 1 NOT NULL,
    current_installment integer DEFAULT 1 NOT NULL,
    status text DEFAULT 'active'::text NOT NULL,
    start_date date NOT NULL,
    end_date date,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    due_day integer,
    CONSTRAINT debts_status_check CHECK ((status = ANY (ARRAY['active'::text, 'paid'::text, 'paused'::text])))
);


ALTER TABLE public.debts OWNER TO myinquisitor_app;

--
-- Name: expense_monthly_status; Type: TABLE; Schema: public; Owner: myinquisitor_app
--

CREATE TABLE public.expense_monthly_status (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    expense_id uuid NOT NULL,
    month date NOT NULL,
    paid boolean DEFAULT false NOT NULL,
    paid_at timestamp with time zone,
    amount_paid numeric(14,2),
    notes text,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.expense_monthly_status OWNER TO myinquisitor_app;

--
-- Name: invite_tokens; Type: TABLE; Schema: public; Owner: myinquisitor_app
--

CREATE TABLE public.invite_tokens (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    token text NOT NULL,
    created_by uuid NOT NULL,
    used boolean DEFAULT false NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.invite_tokens OWNER TO myinquisitor_app;

--
-- Name: monthly_summary; Type: TABLE; Schema: public; Owner: myinquisitor_app
--

CREATE TABLE public.monthly_summary (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    month date NOT NULL,
    total_income numeric(14,2) DEFAULT 0 NOT NULL,
    income_breakdown jsonb,
    total_expenses numeric(14,2) DEFAULT 0 NOT NULL,
    expense_breakdown jsonb,
    total_debt_payments numeric(14,2) DEFAULT 0 NOT NULL,
    debt_breakdown jsonb,
    total_obligations numeric(14,2) DEFAULT 0 NOT NULL,
    net_balance numeric(14,2) DEFAULT 0 NOT NULL,
    projected_income numeric(14,2),
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.monthly_summary OWNER TO myinquisitor_app;

--
-- Name: recurring_expenses; Type: TABLE; Schema: public; Owner: myinquisitor_app
--

CREATE TABLE public.recurring_expenses (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    category_id uuid,
    name text NOT NULL,
    description text,
    amount numeric(14,2) NOT NULL,
    frequency text NOT NULL,
    due_day integer,
    status text DEFAULT 'active'::text NOT NULL,
    start_date date NOT NULL,
    end_date date,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    billing_month integer,
    CONSTRAINT recurring_expenses_billing_month_check CHECK (((billing_month >= 1) AND (billing_month <= 12))),
    CONSTRAINT recurring_expenses_frequency_check CHECK ((frequency = ANY (ARRAY['monthly'::text, 'yearly'::text, 'weekly'::text, 'biweekly'::text, 'once'::text]))),
    CONSTRAINT recurring_expenses_status_check CHECK ((status = ANY (ARRAY['active'::text, 'cancelled'::text])))
);


ALTER TABLE public.recurring_expenses OWNER TO myinquisitor_app;

--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: myinquisitor_app
--

CREATE TABLE public.schema_migrations (
    filename text NOT NULL,
    applied_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO myinquisitor_app;

--
-- Name: transactions; Type: TABLE; Schema: public; Owner: myinquisitor_app
--

CREATE TABLE public.transactions (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    category_id uuid,
    type text NOT NULL,
    amount numeric(14,2) NOT NULL,
    description text,
    source text,
    reference_date date DEFAULT CURRENT_DATE NOT NULL,
    is_recurring boolean DEFAULT false NOT NULL,
    recurring_expense_id uuid,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT transactions_type_check CHECK ((type = ANY (ARRAY['income'::text, 'expense'::text, 'transfer'::text])))
);


ALTER TABLE public.transactions OWNER TO myinquisitor_app;

--
-- Name: users; Type: TABLE; Schema: public; Owner: myinquisitor_app
--

CREATE TABLE public.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    email text NOT NULL,
    email_hash text NOT NULL,
    password_hash text NOT NULL,
    full_name text NOT NULL,
    phone text,
    role text DEFAULT 'user'::text NOT NULL,
    active boolean DEFAULT true NOT NULL,
    super_admin boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT users_role_check CHECK ((role = ANY (ARRAY['user'::text, 'super_admin'::text])))
);


ALTER TABLE public.users OWNER TO myinquisitor_app;

--
-- Data for Name: categories; Type: TABLE DATA; Schema: public; Owner: myinquisitor_app
--

COPY public.categories (id, user_id, name, type, icon, color, created_at) FROM stdin;
\.


--
-- Data for Name: debt_monthly_status; Type: TABLE DATA; Schema: public; Owner: myinquisitor_app
--

COPY public.debt_monthly_status (id, debt_id, month, installment_num, total_installments, amount_due, interest_amount, principal_amount, amount_paid, paid, paid_at, notes, created_at, updated_at) FROM stdin;
8ba75546-4962-4aa0-82f2-d3b1abcc49aa	504c5c0d-2def-42cd-8b58-a71f7ba5a2e9	2026-07-01	2	19	105236.84	26289.47	78947.37	0.00	f	\N	\N	2026-06-02 00:16:12.329189+00	2026-06-02 00:16:12.329189+00
0d58a6b9-24c7-451e-b8d7-76e5684adc08	504c5c0d-2def-42cd-8b58-a71f7ba5a2e9	2026-08-01	3	19	103776.32	24828.95	78947.37	0.00	f	\N	\N	2026-06-02 00:16:12.334243+00	2026-06-02 00:16:12.334243+00
58e507c0-e594-4f64-a9d5-3503bec67c72	504c5c0d-2def-42cd-8b58-a71f7ba5a2e9	2026-09-01	4	19	102315.79	23368.42	78947.37	0.00	f	\N	\N	2026-06-02 00:16:12.339106+00	2026-06-02 00:16:12.339106+00
7cf19cb6-7851-4e82-b5d2-54a991c7e532	504c5c0d-2def-42cd-8b58-a71f7ba5a2e9	2026-10-01	5	19	100855.26	21907.89	78947.37	0.00	f	\N	\N	2026-06-02 00:16:12.344085+00	2026-06-02 00:16:12.344085+00
49d3cf4f-8e21-4c6c-b7b4-5db834f4e269	504c5c0d-2def-42cd-8b58-a71f7ba5a2e9	2026-11-01	6	19	99394.74	20447.37	78947.37	0.00	f	\N	\N	2026-06-02 00:16:12.349062+00	2026-06-02 00:16:12.349062+00
36c0d079-251f-48ac-a7e3-e6cb2544557f	504c5c0d-2def-42cd-8b58-a71f7ba5a2e9	2026-12-01	7	19	97934.21	18986.84	78947.37	0.00	f	\N	\N	2026-06-02 00:16:12.354014+00	2026-06-02 00:16:12.354014+00
353dcf53-dfdc-4574-b38c-29248827c88a	504c5c0d-2def-42cd-8b58-a71f7ba5a2e9	2027-01-01	8	19	96473.69	17526.32	78947.37	0.00	f	\N	\N	2026-06-02 00:16:12.358753+00	2026-06-02 00:16:12.358753+00
90c7934f-8d1b-4d94-8b19-dc6cb8d9250b	504c5c0d-2def-42cd-8b58-a71f7ba5a2e9	2027-02-01	9	19	95013.16	16065.79	78947.37	0.00	f	\N	\N	2026-06-02 00:16:12.363249+00	2026-06-02 00:16:12.363249+00
efe14311-5d5c-46bd-a4b8-2db1456cd0af	504c5c0d-2def-42cd-8b58-a71f7ba5a2e9	2027-03-01	10	19	93552.63	14605.26	78947.37	0.00	f	\N	\N	2026-06-02 00:16:12.368045+00	2026-06-02 00:16:12.368045+00
9518e6e6-4c23-403c-bd62-278176f9913e	504c5c0d-2def-42cd-8b58-a71f7ba5a2e9	2027-04-01	11	19	92092.11	13144.74	78947.37	0.00	f	\N	\N	2026-06-02 00:16:12.372738+00	2026-06-02 00:16:12.372738+00
c38dee92-1428-436e-bb73-62132490588e	504c5c0d-2def-42cd-8b58-a71f7ba5a2e9	2027-05-01	12	19	90631.58	11684.21	78947.37	0.00	f	\N	\N	2026-06-02 00:16:12.377442+00	2026-06-02 00:16:12.377442+00
c1283373-7286-4b44-8771-ff0e679a4940	504c5c0d-2def-42cd-8b58-a71f7ba5a2e9	2027-06-01	13	19	89171.05	10223.68	78947.37	0.00	f	\N	\N	2026-06-02 00:16:12.381988+00	2026-06-02 00:16:12.381988+00
bc5c62a0-1f06-4cd5-be06-70d016c312fe	504c5c0d-2def-42cd-8b58-a71f7ba5a2e9	2027-07-01	14	19	87710.53	8763.16	78947.37	0.00	f	\N	\N	2026-06-02 00:16:12.386643+00	2026-06-02 00:16:12.386643+00
9a9e1a6c-c4e5-4890-a946-05bc52de59b7	504c5c0d-2def-42cd-8b58-a71f7ba5a2e9	2027-08-01	15	19	86250.00	7302.63	78947.37	0.00	f	\N	\N	2026-06-02 00:16:12.391308+00	2026-06-02 00:16:12.391308+00
8f3bda18-4602-4c34-9a86-f9fd828a1c5e	504c5c0d-2def-42cd-8b58-a71f7ba5a2e9	2027-09-01	16	19	84789.48	5842.11	78947.37	0.00	f	\N	\N	2026-06-02 00:16:12.396069+00	2026-06-02 00:16:12.396069+00
fde17bd5-7ab0-485b-9c7c-62bf7ad55ffa	504c5c0d-2def-42cd-8b58-a71f7ba5a2e9	2027-10-01	17	19	83328.95	4381.58	78947.37	0.00	f	\N	\N	2026-06-02 00:16:12.400734+00	2026-06-02 00:16:12.400734+00
a9a878ad-706e-4509-abdb-80777dcd3124	504c5c0d-2def-42cd-8b58-a71f7ba5a2e9	2027-11-01	18	19	81868.42	2921.05	78947.37	0.00	f	\N	\N	2026-06-02 00:16:12.405352+00	2026-06-02 00:16:12.405352+00
575f1115-1c1f-4872-84f8-37c8a3a8ceeb	504c5c0d-2def-42cd-8b58-a71f7ba5a2e9	2027-12-01	19	19	80407.90	1460.53	78947.37	0.00	f	\N	\N	2026-06-02 00:16:12.410236+00	2026-06-02 00:16:12.410236+00
f81ab824-a45d-4713-8de9-cf538d1695c5	504c5c0d-2def-42cd-8b58-a71f7ba5a2e9	2026-06-01	1	19	106697.37	27750.00	78947.37	106697.37	t	\N	\N	2026-06-02 00:16:12.324166+00	2026-06-02 00:16:36.521038+00
6503ace1-3ff4-43e2-8e8f-4ba51f388d2a	ce14bd06-6dde-4f21-a390-11ad71115684	2026-07-01	1	4	94525.76	6611.18	87914.58	0.00	f	\N	\N	2026-06-12 21:55:49.559911+00	2026-06-12 21:55:49.559911+00
65586b6e-31d4-40a1-a495-b1eab1384674	ce14bd06-6dde-4f21-a390-11ad71115684	2026-08-01	2	4	92872.96	4958.38	87914.58	0.00	f	\N	\N	2026-06-12 21:55:49.571018+00	2026-06-12 21:55:49.571018+00
4a19795c-bf2e-46f1-ae0c-e579e1a680a1	ce14bd06-6dde-4f21-a390-11ad71115684	2026-09-01	3	4	91220.17	3305.59	87914.58	0.00	f	\N	\N	2026-06-12 21:55:49.576768+00	2026-06-12 21:55:49.576768+00
0351766a-2b8c-4ba4-a4bb-23e73b0f1a8b	ce14bd06-6dde-4f21-a390-11ad71115684	2026-10-01	4	4	89567.37	1652.79	87914.58	0.00	f	\N	\N	2026-06-12 21:55:49.582752+00	2026-06-12 21:55:49.582752+00
37bb16bd-52d5-4379-91f9-6f3ec817d19e	017ed55d-a4e0-4623-87e2-2e1f4010e7ce	2026-06-01	1	18	139199.88	35195.19	104004.69	0.00	f	\N	\N	2026-06-12 22:06:41.244064+00	2026-06-12 22:06:41.244064+00
6276735f-3ebd-45bd-a4f2-6312fe9f4ce6	017ed55d-a4e0-4623-87e2-2e1f4010e7ce	2026-07-01	2	18	137244.59	33239.90	104004.69	0.00	f	\N	\N	2026-06-12 22:06:41.250151+00	2026-06-12 22:06:41.250151+00
0f1a7086-d501-453d-8dde-57c9662d0e7a	017ed55d-a4e0-4623-87e2-2e1f4010e7ce	2026-08-01	3	18	135289.30	31284.61	104004.69	0.00	f	\N	\N	2026-06-12 22:06:41.255832+00	2026-06-12 22:06:41.255832+00
6ed9cab8-356d-40cd-8411-962738ec733e	017ed55d-a4e0-4623-87e2-2e1f4010e7ce	2026-09-01	4	18	133334.01	29329.32	104004.69	0.00	f	\N	\N	2026-06-12 22:06:41.262158+00	2026-06-12 22:06:41.262158+00
e91c66f0-3c77-4ff2-a0be-4db4ff7d420b	017ed55d-a4e0-4623-87e2-2e1f4010e7ce	2026-10-01	5	18	131378.72	27374.03	104004.69	0.00	f	\N	\N	2026-06-12 22:06:41.267886+00	2026-06-12 22:06:41.267886+00
77ca0c93-ecf6-4c46-b849-1d12db80c7b1	017ed55d-a4e0-4623-87e2-2e1f4010e7ce	2026-11-01	6	18	129423.44	25418.75	104004.69	0.00	f	\N	\N	2026-06-12 22:06:41.273828+00	2026-06-12 22:06:41.273828+00
693a222c-9914-4ad4-984f-c92978c784f0	017ed55d-a4e0-4623-87e2-2e1f4010e7ce	2026-12-01	7	18	127468.15	23463.46	104004.69	0.00	f	\N	\N	2026-06-12 22:06:41.279896+00	2026-06-12 22:06:41.279896+00
b283dba7-5fc1-4ae7-815e-e78563668149	017ed55d-a4e0-4623-87e2-2e1f4010e7ce	2027-01-01	8	18	125512.86	21508.17	104004.69	0.00	f	\N	\N	2026-06-12 22:06:41.285885+00	2026-06-12 22:06:41.285885+00
59ab076a-9078-4382-9d43-4f1fd38bb92c	017ed55d-a4e0-4623-87e2-2e1f4010e7ce	2027-02-01	9	18	123557.57	19552.88	104004.69	0.00	f	\N	\N	2026-06-12 22:06:41.292209+00	2026-06-12 22:06:41.292209+00
2b4c6c84-8f13-4cd2-a4cd-d58699c7135d	017ed55d-a4e0-4623-87e2-2e1f4010e7ce	2027-03-01	10	18	121602.28	17597.59	104004.69	0.00	f	\N	\N	2026-06-12 22:06:41.298126+00	2026-06-12 22:06:41.298126+00
c7a059b0-0bcf-4490-a68c-7e731755b8ab	017ed55d-a4e0-4623-87e2-2e1f4010e7ce	2027-04-01	11	18	119647.00	15642.31	104004.69	0.00	f	\N	\N	2026-06-12 22:06:41.304342+00	2026-06-12 22:06:41.304342+00
ee50c93e-f6a8-4eff-8f75-fc82d9a03d01	017ed55d-a4e0-4623-87e2-2e1f4010e7ce	2027-05-01	12	18	117691.71	13687.02	104004.69	0.00	f	\N	\N	2026-06-12 22:06:41.310569+00	2026-06-12 22:06:41.310569+00
14dfb16d-5771-4bff-997c-c645d7b1cbd6	017ed55d-a4e0-4623-87e2-2e1f4010e7ce	2027-06-01	13	18	115736.42	11731.73	104004.69	0.00	f	\N	\N	2026-06-12 22:06:41.316818+00	2026-06-12 22:06:41.316818+00
58cad300-99d4-4a38-8a2c-9b284f45d7f3	017ed55d-a4e0-4623-87e2-2e1f4010e7ce	2027-07-01	14	18	113781.13	9776.44	104004.69	0.00	f	\N	\N	2026-06-12 22:06:41.322724+00	2026-06-12 22:06:41.322724+00
c45057c7-d985-419d-b8d1-d55eb9ed0696	017ed55d-a4e0-4623-87e2-2e1f4010e7ce	2027-08-01	15	18	111825.84	7821.15	104004.69	0.00	f	\N	\N	2026-06-12 22:06:41.328634+00	2026-06-12 22:06:41.328634+00
b4234019-7264-4e10-9b7c-57ffaccd6a7a	017ed55d-a4e0-4623-87e2-2e1f4010e7ce	2027-09-01	16	18	109870.55	5865.86	104004.69	0.00	f	\N	\N	2026-06-12 22:06:41.334123+00	2026-06-12 22:06:41.334123+00
36c6f763-c891-4316-b2cf-233890338daf	017ed55d-a4e0-4623-87e2-2e1f4010e7ce	2027-10-01	17	18	107915.27	3910.58	104004.69	0.00	f	\N	\N	2026-06-12 22:06:41.339295+00	2026-06-12 22:06:41.339295+00
9d5dfe81-c47c-4331-8c0f-3e6026ba5249	017ed55d-a4e0-4623-87e2-2e1f4010e7ce	2027-11-01	18	18	105959.98	1955.29	104004.69	0.00	f	\N	\N	2026-06-12 22:06:41.344859+00	2026-06-12 22:06:41.344859+00
\.


--
-- Data for Name: debts; Type: TABLE DATA; Schema: public; Owner: myinquisitor_app
--

COPY public.debts (id, user_id, category_id, name, description, total_amount, remaining_amount, interest_rate, total_installments, current_installment, status, start_date, end_date, created_at, updated_at, due_day) FROM stdin;
504c5c0d-2def-42cd-8b58-a71f7ba5a2e9	7fabe4d1-f6c7-48e3-a38c-78ffc64d52e1	\N	PC NU	\N	1500000.00	1393302.63	1.85	19	1	active	2026-06-01	\N	2026-06-02 00:16:12.313265+00	2026-06-02 00:16:36.53074+00	20
ce14bd06-6dde-4f21-a390-11ad71115684	b96400c5-92d3-461d-98a1-360a320e09d6	\N	Rappi Card	\N	351658.30	351658.30	1.88	4	0	active	2026-07-01	\N	2026-06-12 21:55:49.514509+00	2026-06-12 21:55:49.514509+00	10
017ed55d-a4e0-4623-87e2-2e1f4010e7ce	b96400c5-92d3-461d-98a1-360a320e09d6	\N	PC NU	\N	1872084.40	1872084.40	1.88	18	0	active	2026-06-01	\N	2026-06-12 22:06:41.200409+00	2026-06-12 22:06:41.200409+00	28
\.


--
-- Data for Name: expense_monthly_status; Type: TABLE DATA; Schema: public; Owner: myinquisitor_app
--

COPY public.expense_monthly_status (id, expense_id, month, paid, paid_at, amount_paid, notes, created_at) FROM stdin;
146cd69b-b7bb-4d00-afcb-f06216d50d02	014f306f-5e3c-4732-b2c5-472b6ac10f38	2026-06-01	t	\N	\N	\N	2026-06-01 23:03:33.127482+00
128fd98d-3bc4-4478-a169-e57be40a01d2	dd521805-fb4b-4f67-b989-ed0b0d0ecdd3	2026-06-01	t	\N	\N	\N	2026-06-01 23:03:34.198264+00
\.


--
-- Data for Name: invite_tokens; Type: TABLE DATA; Schema: public; Owner: myinquisitor_app
--

COPY public.invite_tokens (id, token, created_by, used, expires_at, created_at) FROM stdin;
96773ec4-1d38-45e8-8a43-6c280bb03daf	b926814b11fc3a41716d6c6ae499ff80ce06c5b032d267074cb08139d287e8ff	24ff88fb-3b05-4b37-88d9-546c7d4a797f	t	2026-05-31 15:39:55.565952+00	2026-05-28 15:39:55.566276+00
\.


--
-- Data for Name: monthly_summary; Type: TABLE DATA; Schema: public; Owner: myinquisitor_app
--

COPY public.monthly_summary (id, user_id, month, total_income, income_breakdown, total_expenses, expense_breakdown, total_debt_payments, debt_breakdown, total_obligations, net_balance, projected_income, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: recurring_expenses; Type: TABLE DATA; Schema: public; Owner: myinquisitor_app
--

COPY public.recurring_expenses (id, user_id, category_id, name, description, amount, frequency, due_day, status, start_date, end_date, created_at, updated_at, billing_month) FROM stdin;
dd521805-fb4b-4f67-b989-ed0b0d0ecdd3	7fabe4d1-f6c7-48e3-a38c-78ffc64d52e1	\N	comida de los gatos	\N	50000.00	monthly	30	active	2026-05-28	\N	2026-05-28 16:16:42.925037+00	2026-05-28 16:16:42.925037+00	\N
014f306f-5e3c-4732-b2c5-472b6ac10f38	7fabe4d1-f6c7-48e3-a38c-78ffc64d52e1	\N	Datos Movistar	\N	25000.00	monthly	10	active	2026-05-28	\N	2026-05-28 16:17:03.067077+00	2026-05-28 16:17:03.067077+00	\N
42c06be0-560d-4e50-9492-bb5f7ad7b2b0	b96400c5-92d3-461d-98a1-360a320e09d6	\N	Tidal	\N	7500.00	monthly	20	active	2026-06-01	\N	2026-06-12 21:40:28.790997+00	2026-06-12 21:40:28.790997+00	\N
b7315608-6ed2-4676-93f2-40a5a5c6e094	b96400c5-92d3-461d-98a1-360a320e09d6	\N	comida de los gatos	\N	70000.00	monthly	\N	active	2026-06-01	\N	2026-06-12 21:44:30.905566+00	2026-06-12 21:44:30.905566+00	\N
bdfca0f9-5865-4211-8f07-09eca33fae2f	b96400c5-92d3-461d-98a1-360a320e09d6	\N	Datos Movistar	\N	25000.00	monthly	10	active	2026-06-01	\N	2026-06-12 21:45:07.361771+00	2026-06-12 21:45:07.361771+00	\N
e024235a-9dcc-419a-a107-ae97cc15f24f	b96400c5-92d3-461d-98a1-360a320e09d6	\N	ahorro productos personales	\N	30000.00	monthly	\N	active	2026-06-01	\N	2026-06-12 21:46:06.582876+00	2026-06-12 21:46:06.582876+00	\N
857b8347-796e-4df5-98d9-a8e5f98a2b53	b96400c5-92d3-461d-98a1-360a320e09d6	\N	gasolina	\N	25000.00	monthly	\N	active	2026-06-01	\N	2026-06-12 21:46:33.60109+00	2026-06-12 21:46:33.60109+00	\N
53951d92-b8bd-4666-9ad5-5d89b00d560a	b96400c5-92d3-461d-98a1-360a320e09d6	\N	ahorro gatos	\N	10000.00	monthly	\N	active	2026-06-01	\N	2026-06-12 21:47:41.742293+00	2026-06-12 21:47:41.742293+00	\N
60f4539e-145c-4c0c-8443-deca2b4d6d9a	b96400c5-92d3-461d-98a1-360a320e09d6	\N	premios gatos	\N	20000.00	monthly	\N	active	2026-06-01	\N	2026-06-12 21:48:04.895068+00	2026-06-12 21:48:04.895068+00	\N
a0031b4a-e5b9-4b9c-8226-dd5ccfbf6513	b96400c5-92d3-461d-98a1-360a320e09d6	\N	ahorro personal	\N	50000.00	monthly	\N	active	2026-06-01	\N	2026-06-12 21:49:56.003022+00	2026-06-12 21:49:56.003022+00	\N
a2175e75-4570-46fd-8fcd-04fbd1a1d963	b96400c5-92d3-461d-98a1-360a320e09d6	\N	cuota de manejo Nu	\N	12000.00	monthly	28	active	2026-06-01	\N	2026-06-12 22:07:20.963088+00	2026-06-12 22:07:20.963088+00	\N
23238b95-b934-46e3-96e5-b0b708576c82	b96400c5-92d3-461d-98a1-360a320e09d6	\N	saldo a universidad	\N	200000.00	once	\N	cancelled	2026-06-12	2026-06-12	2026-06-12 22:36:03.134065+00	2026-06-12 22:36:03.134065+00	\N
2f979638-fbfa-45f1-acd2-b72e2db0ba0e	b96400c5-92d3-461d-98a1-360a320e09d6	\N	saldo a universidad	\N	200000.00	once	\N	cancelled	2026-06-01	2026-06-01	2026-06-12 22:45:16.317253+00	2026-06-12 22:45:16.317253+00	\N
b0afabd2-620f-4fdd-8c8e-293fb7683412	b96400c5-92d3-461d-98a1-360a320e09d6	\N	saldo a universidad	\N	150000.00	once	30	active	2026-06-01	2026-06-01	2026-06-12 22:47:14.957643+00	2026-06-15 14:10:56.674444+00	\N
\.


--
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: myinquisitor_app
--

COPY public.schema_migrations (filename, applied_at) FROM stdin;
001_init.up.sql	2026-05-28 15:38:37.666846+00
002_invite_tokens.up.sql	2026-05-28 15:38:37.704879+00
003_debt_due_day.up.sql	2026-05-28 15:38:37.714651+00
004_paused_status.up.sql	2026-06-01 22:52:48.448825+00
002_expense_updates.up.sql	2026-06-12 22:35:26.634924+00
\.


--
-- Data for Name: transactions; Type: TABLE DATA; Schema: public; Owner: myinquisitor_app
--

COPY public.transactions (id, user_id, category_id, type, amount, description, source, reference_date, is_recurring, recurring_expense_id, created_at) FROM stdin;
638bac69-68f0-4558-99c5-c06453cec948	7fabe4d1-f6c7-48e3-a38c-78ffc64d52e1	\N	expense	50000.00	Gasto: prueba - 2026-05-01	\N	2026-05-28	f	\N	2026-05-28 16:17:35.151112+00
4da3bd56-b0bf-42c9-9a4d-4d5b6bf1d1b2	7fabe4d1-f6c7-48e3-a38c-78ffc64d52e1	\N	expense	80545.03	Pago de deuda: PC NU - Cuota 19/19	\N	2026-06-01	f	\N	2026-06-01 17:02:13.485691+00
0cac0584-fec8-4125-b532-55250569bff3	7fabe4d1-f6c7-48e3-a38c-78ffc64d52e1	\N	income	80545.03	\N	\N	2026-06-01	f	\N	2026-06-01 23:03:17.155727+00
1c0785ed-5b0e-48d6-a99d-9e40541f2e5d	7fabe4d1-f6c7-48e3-a38c-78ffc64d52e1	\N	expense	25000.00	Gasto: Datos Movistar - 2026-06-01	\N	2026-06-01	f	\N	2026-06-01 23:03:33.172054+00
53405721-0adb-46fb-8d20-96c6fbc1e078	7fabe4d1-f6c7-48e3-a38c-78ffc64d52e1	\N	expense	7500.00	Gasto: Tidal - 2026-06-01	\N	2026-06-01	f	\N	2026-06-01 23:03:33.735837+00
02b56299-1ec7-423f-bccc-d1dd78f73792	7fabe4d1-f6c7-48e3-a38c-78ffc64d52e1	\N	expense	50000.00	Gasto: comida de los gatos - 2026-06-01	\N	2026-06-01	f	\N	2026-06-01 23:03:34.207489+00
a619c489-6a43-49fe-a69e-e00afe948e66	7fabe4d1-f6c7-48e3-a38c-78ffc64d52e1	\N	expense	81259.87	Pago de deuda: nu - Cuota 1/19	\N	2026-06-01	f	\N	2026-06-01 23:11:26.751581+00
8e96f1d3-0312-4623-b792-cf4b3a68622a	7fabe4d1-f6c7-48e3-a38c-78ffc64d52e1	\N	expense	81138.16	Pago de deuda: nu - Cuota 2/19	\N	2026-06-01	f	\N	2026-06-01 23:11:28.845444+00
d4af45ce-ef9e-48d8-8585-5f327f0cc090	7fabe4d1-f6c7-48e3-a38c-78ffc64d52e1	\N	expense	106697.37	Pago de deuda: PC NU - Cuota 1/19	\N	2026-06-01	f	\N	2026-06-02 00:16:36.535774+00
18280555-2bc9-4126-882c-2f64caf9178c	7fabe4d1-f6c7-48e3-a38c-78ffc64d52e1	\N	expense	100000.00	Gasto: Kevin 11 - 2026-06-01	\N	2026-06-12	f	\N	2026-06-12 22:50:11.191702+00
57d38b92-03d9-456c-8af7-5ba2a00b9d1f	7fabe4d1-f6c7-48e3-a38c-78ffc64d52e1	\N	expense	100000.00	Gasto: pruebaaa - 2026-06-01	\N	2026-06-12	f	\N	2026-06-12 22:57:45.62113+00
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: myinquisitor_app
--

COPY public.users (id, email, email_hash, password_hash, full_name, phone, role, active, super_admin, created_at, updated_at) FROM stdin;
24ff88fb-3b05-4b37-88d9-546c7d4a797f	f861b833f2c9b3c7bfa736cefc378ab24d5b45d89f7b6df5fd94d51e9b24d7b5a5b0bb18e1d359a62bf8dfaf6c9690c02e28	8531e4c91f8217750dfb95f9fdbed99ff2e7966e6f16145b3a2d6ef006e2b037	$2a$12$.agA0nE1/QeDCT0Aq992lOezHmci2WuB5P.UvBHpbVUPhDqJRNdOy	2d474a5c94ee23d75b8df6eeee1b3e5755e3d03307b5c910454c51e34fcfee75c866e891307f21	\N	super_admin	t	t	2026-05-28 15:39:43.905343+00	2026-05-28 15:39:43.905343+00
dcf17f39-076f-470b-b0b3-f34788a37e04	974faac9d0bf49a74d20eac45a20e263f5134178591d67e9d5419b2673997c0ca56a3e00d9bb989dc34dd795cad19375e8b1088a	0009bbb9ff9132187c83d77bedac565241b9d3c774d3a90468b156d6b15b51d4	$2a$12$VjB.0QduULSA9Ni5gdredevouf8CL.d46QpuqzE9HxF29V2/zdChu	889caee981577296a22f2f1b77d9ddeda7abf5ff01ce0b06511ef362af2e6aabe85f3ee2fade0e095f2e61ece87c248fb3526e6a510e7eb9d745d0	\N	user	t	f	2026-05-28 15:42:20.190116+00	2026-05-28 15:42:20.190116+00
7fabe4d1-f6c7-48e3-a38c-78ffc64d52e1	18f0d73926b9f363ed54cdc58cecd85668c905905b81e4fd1fa79b677c4757db714141ed7ad7efe9ea	f660ab912ec121d1b1e928a0bb4bc61b15f5ad44d5efdc4e1c92a25e99b8e44a	$2a$12$TFd5XUxP0icijyO95oSXRepRtQOD7.ZP2.FB1zzej18JFWTpTVWYS	994b55c8ea75f202dc2af0258f908e0c2fbdf604e4ebc408781867be6c40b397d4ed7687f2c1	\N	user	t	f	2026-05-28 16:15:38.317903+00	2026-06-12 21:34:30.712288+00
b96400c5-92d3-461d-98a1-360a320e09d6	0c9f4ad3620cd546eb488f70e905ed63fdb7f9087008cbe948f4654ee6502d7a0bb658	53445089386ad8e5d02a7c7f02674052d5a45a9acac7b57e4f4342075ecf2787	$2a$12$FqUAwzVh8iFBvjY2QkGn1u2Wh2/xNpe5GfBMUbR9j2B1ZOfnD6hfG	ffb638a792d458c127534dd698dff74b217112126e4429fb548fce40852d56f873	\N	super_admin	t	f	2026-06-12 21:34:59.312239+00	2026-06-12 21:35:16.776718+00
\.


--
-- Name: categories categories_pkey; Type: CONSTRAINT; Schema: public; Owner: myinquisitor_app
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (id);


--
-- Name: debt_monthly_status debt_monthly_status_debt_id_month_key; Type: CONSTRAINT; Schema: public; Owner: myinquisitor_app
--

ALTER TABLE ONLY public.debt_monthly_status
    ADD CONSTRAINT debt_monthly_status_debt_id_month_key UNIQUE (debt_id, month);


--
-- Name: debt_monthly_status debt_monthly_status_pkey; Type: CONSTRAINT; Schema: public; Owner: myinquisitor_app
--

ALTER TABLE ONLY public.debt_monthly_status
    ADD CONSTRAINT debt_monthly_status_pkey PRIMARY KEY (id);


--
-- Name: debts debts_pkey; Type: CONSTRAINT; Schema: public; Owner: myinquisitor_app
--

ALTER TABLE ONLY public.debts
    ADD CONSTRAINT debts_pkey PRIMARY KEY (id);


--
-- Name: expense_monthly_status expense_monthly_status_expense_id_month_key; Type: CONSTRAINT; Schema: public; Owner: myinquisitor_app
--

ALTER TABLE ONLY public.expense_monthly_status
    ADD CONSTRAINT expense_monthly_status_expense_id_month_key UNIQUE (expense_id, month);


--
-- Name: expense_monthly_status expense_monthly_status_pkey; Type: CONSTRAINT; Schema: public; Owner: myinquisitor_app
--

ALTER TABLE ONLY public.expense_monthly_status
    ADD CONSTRAINT expense_monthly_status_pkey PRIMARY KEY (id);


--
-- Name: invite_tokens invite_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: myinquisitor_app
--

ALTER TABLE ONLY public.invite_tokens
    ADD CONSTRAINT invite_tokens_pkey PRIMARY KEY (id);


--
-- Name: invite_tokens invite_tokens_token_key; Type: CONSTRAINT; Schema: public; Owner: myinquisitor_app
--

ALTER TABLE ONLY public.invite_tokens
    ADD CONSTRAINT invite_tokens_token_key UNIQUE (token);


--
-- Name: monthly_summary monthly_summary_pkey; Type: CONSTRAINT; Schema: public; Owner: myinquisitor_app
--

ALTER TABLE ONLY public.monthly_summary
    ADD CONSTRAINT monthly_summary_pkey PRIMARY KEY (id);


--
-- Name: monthly_summary monthly_summary_user_id_month_key; Type: CONSTRAINT; Schema: public; Owner: myinquisitor_app
--

ALTER TABLE ONLY public.monthly_summary
    ADD CONSTRAINT monthly_summary_user_id_month_key UNIQUE (user_id, month);


--
-- Name: recurring_expenses recurring_expenses_pkey; Type: CONSTRAINT; Schema: public; Owner: myinquisitor_app
--

ALTER TABLE ONLY public.recurring_expenses
    ADD CONSTRAINT recurring_expenses_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: myinquisitor_app
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (filename);


--
-- Name: transactions transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: myinquisitor_app
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_pkey PRIMARY KEY (id);


--
-- Name: users users_email_hash_key; Type: CONSTRAINT; Schema: public; Owner: myinquisitor_app
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_hash_key UNIQUE (email_hash);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: myinquisitor_app
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_categories_user_id; Type: INDEX; Schema: public; Owner: myinquisitor_app
--

CREATE INDEX idx_categories_user_id ON public.categories USING btree (user_id);


--
-- Name: idx_debt_monthly_status_debt_id; Type: INDEX; Schema: public; Owner: myinquisitor_app
--

CREATE INDEX idx_debt_monthly_status_debt_id ON public.debt_monthly_status USING btree (debt_id);


--
-- Name: idx_debt_monthly_status_month; Type: INDEX; Schema: public; Owner: myinquisitor_app
--

CREATE INDEX idx_debt_monthly_status_month ON public.debt_monthly_status USING btree (month);


--
-- Name: idx_debts_status; Type: INDEX; Schema: public; Owner: myinquisitor_app
--

CREATE INDEX idx_debts_status ON public.debts USING btree (status);


--
-- Name: idx_debts_user_id; Type: INDEX; Schema: public; Owner: myinquisitor_app
--

CREATE INDEX idx_debts_user_id ON public.debts USING btree (user_id);


--
-- Name: idx_expense_monthly_status_expense_id; Type: INDEX; Schema: public; Owner: myinquisitor_app
--

CREATE INDEX idx_expense_monthly_status_expense_id ON public.expense_monthly_status USING btree (expense_id);


--
-- Name: idx_invite_tokens_created_by; Type: INDEX; Schema: public; Owner: myinquisitor_app
--

CREATE INDEX idx_invite_tokens_created_by ON public.invite_tokens USING btree (created_by);


--
-- Name: idx_invite_tokens_token; Type: INDEX; Schema: public; Owner: myinquisitor_app
--

CREATE INDEX idx_invite_tokens_token ON public.invite_tokens USING btree (token);


--
-- Name: idx_monthly_summary_month; Type: INDEX; Schema: public; Owner: myinquisitor_app
--

CREATE INDEX idx_monthly_summary_month ON public.monthly_summary USING btree (month);


--
-- Name: idx_monthly_summary_user_id; Type: INDEX; Schema: public; Owner: myinquisitor_app
--

CREATE INDEX idx_monthly_summary_user_id ON public.monthly_summary USING btree (user_id);


--
-- Name: idx_recurring_expenses_user_id; Type: INDEX; Schema: public; Owner: myinquisitor_app
--

CREATE INDEX idx_recurring_expenses_user_id ON public.recurring_expenses USING btree (user_id);


--
-- Name: idx_transactions_reference_date; Type: INDEX; Schema: public; Owner: myinquisitor_app
--

CREATE INDEX idx_transactions_reference_date ON public.transactions USING btree (reference_date);


--
-- Name: idx_transactions_type; Type: INDEX; Schema: public; Owner: myinquisitor_app
--

CREATE INDEX idx_transactions_type ON public.transactions USING btree (type);


--
-- Name: idx_transactions_user_id; Type: INDEX; Schema: public; Owner: myinquisitor_app
--

CREATE INDEX idx_transactions_user_id ON public.transactions USING btree (user_id);


--
-- Name: categories categories_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: myinquisitor_app
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: debt_monthly_status debt_monthly_status_debt_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: myinquisitor_app
--

ALTER TABLE ONLY public.debt_monthly_status
    ADD CONSTRAINT debt_monthly_status_debt_id_fkey FOREIGN KEY (debt_id) REFERENCES public.debts(id) ON DELETE CASCADE;


--
-- Name: debts debts_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: myinquisitor_app
--

ALTER TABLE ONLY public.debts
    ADD CONSTRAINT debts_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.categories(id);


--
-- Name: debts debts_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: myinquisitor_app
--

ALTER TABLE ONLY public.debts
    ADD CONSTRAINT debts_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: expense_monthly_status expense_monthly_status_expense_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: myinquisitor_app
--

ALTER TABLE ONLY public.expense_monthly_status
    ADD CONSTRAINT expense_monthly_status_expense_id_fkey FOREIGN KEY (expense_id) REFERENCES public.recurring_expenses(id) ON DELETE CASCADE;


--
-- Name: invite_tokens invite_tokens_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: myinquisitor_app
--

ALTER TABLE ONLY public.invite_tokens
    ADD CONSTRAINT invite_tokens_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: monthly_summary monthly_summary_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: myinquisitor_app
--

ALTER TABLE ONLY public.monthly_summary
    ADD CONSTRAINT monthly_summary_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: recurring_expenses recurring_expenses_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: myinquisitor_app
--

ALTER TABLE ONLY public.recurring_expenses
    ADD CONSTRAINT recurring_expenses_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.categories(id);


--
-- Name: recurring_expenses recurring_expenses_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: myinquisitor_app
--

ALTER TABLE ONLY public.recurring_expenses
    ADD CONSTRAINT recurring_expenses_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: transactions transactions_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: myinquisitor_app
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.categories(id);


--
-- Name: transactions transactions_recurring_expense_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: myinquisitor_app
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_recurring_expense_id_fkey FOREIGN KEY (recurring_expense_id) REFERENCES public.recurring_expenses(id);


--
-- Name: transactions transactions_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: myinquisitor_app
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: pg_database_owner
--

GRANT USAGE ON SCHEMA public TO myinquisitor_app;


--
-- Name: DEFAULT PRIVILEGES FOR TABLES; Type: DEFAULT ACL; Schema: public; Owner: myinquisitor_app
--

ALTER DEFAULT PRIVILEGES FOR ROLE myinquisitor_app IN SCHEMA public GRANT SELECT,INSERT,DELETE,UPDATE ON TABLES TO myinquisitor_app;


--
-- PostgreSQL database dump complete
--

\unrestrict BeMnSNKvfKY8R0MiIoNqdPr1a5cAuPRafnaIC5JOeqU4OuwcZENbUk1CnOGSq6D

