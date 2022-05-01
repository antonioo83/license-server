--
-- PostgreSQL database dump
--

-- Dumped from database version 14.2
-- Dumped by pg_dump version 14.2

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
-- Name: ln_customers; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ln_customers (
                                     id integer NOT NULL,
                                     user_id integer NOT NULL,
                                     code character varying(64) DEFAULT ''::character varying NOT NULL,
                                     type character varying(20) DEFAULT 'service'::character varying NOT NULL,
                                     title character varying(100) DEFAULT ''::character varying NOT NULL,
                                     inn character varying(12) DEFAULT ''::character varying NOT NULL,
                                     description character varying(256) DEFAULT ''::character varying NOT NULL,
                                     created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
                                     updated_at timestamp without time zone,
                                     deleted_at timestamp without time zone
);


ALTER TABLE public.ln_customers OWNER TO postgres;

--
-- Name: ln_customers_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.ln_customers_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.ln_customers_id_seq OWNER TO postgres;

--
-- Name: ln_customers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.ln_customers_id_seq OWNED BY public.ln_customers.id;


--
-- Name: ln_licenses; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ln_licenses (
                                    id integer NOT NULL,
                                    customer_id integer NOT NULL,
                                    product_type character varying(20) DEFAULT ''::character varying NOT NULL,
                                    callback_url character varying(500) DEFAULT ''::character varying NOT NULL,
                                    count integer DEFAULT 1 NOT NULL,
                                    license_key character varying(500) DEFAULT ''::character varying NOT NULL,
                                    registration_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
                                    activation_at timestamp without time zone,
                                    expiration_at timestamp without time zone NOT NULL,
                                    deleted_at timestamp without time zone,
                                    duration integer NOT NULL,
                                    description character varying(500) DEFAULT ''::character varying NOT NULL,
                                    code character varying(64) NOT NULL,
                                    callback_attempts integer DEFAULT 0 NOT NULL,
                                    is_sent_callback smallint DEFAULT 0 NOT NULL
);


ALTER TABLE public.ln_licenses OWNER TO postgres;

--
-- Name: ln_licenses_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.ln_licenses_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.ln_licenses_id_seq OWNER TO postgres;

--
-- Name: ln_licenses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.ln_licenses_id_seq OWNED BY public.ln_licenses.id;


--
-- Name: ln_user_actions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ln_user_actions (
                                        id integer NOT NULL,
                                        action character varying(20) NOT NULL,
                                        description character varying(500) DEFAULT ''::character varying NOT NULL,
                                        created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
                                        updated_at timestamp without time zone,
                                        deleted_at timestamp without time zone
);


ALTER TABLE public.ln_user_actions OWNER TO postgres;

--
-- Name: ln_user_actions_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.ln_user_actions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.ln_user_actions_id_seq OWNER TO postgres;

--
-- Name: ln_user_actions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.ln_user_actions_id_seq OWNED BY public.ln_user_actions.id;


--
-- Name: ln_user_permissions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ln_user_permissions (
                                            id integer NOT NULL,
                                            user_id integer NOT NULL,
                                            action_id integer NOT NULL,
                                            product_type character varying(50) DEFAULT ''::character varying NOT NULL,
                                            created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
                                            updated_at timestamp without time zone,
                                            deleted_at timestamp without time zone
);


ALTER TABLE public.ln_user_permissions OWNER TO postgres;

--
-- Name: ln_users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ln_users (
                                 id integer NOT NULL,
                                 code character varying(64) DEFAULT ''::character varying NOT NULL,
                                 role character varying(20) DEFAULT 'service'::character varying NOT NULL,
                                 title character varying(100) DEFAULT ''::character varying NOT NULL,
                                 auth_token character varying(256) DEFAULT ''::character varying NOT NULL,
                                 description character varying(256) DEFAULT ''::character varying NOT NULL,
                                 created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
                                 updated_at timestamp without time zone,
                                 deleted_at timestamp without time zone
);


ALTER TABLE public.ln_users OWNER TO postgres;

--
-- Name: ln_users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.ln_users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.ln_users_id_seq OWNER TO postgres;

--
-- Name: ln_users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.ln_users_id_seq OWNED BY public.ln_users.id;


--
-- Name: user_permissions_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.user_permissions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.user_permissions_id_seq OWNER TO postgres;

--
-- Name: user_permissions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.user_permissions_id_seq OWNED BY public.ln_user_permissions.id;


--
-- Name: ln_customers id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ln_customers ALTER COLUMN id SET DEFAULT nextval('public.ln_customers_id_seq'::regclass);


--
-- Name: ln_licenses id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ln_licenses ALTER COLUMN id SET DEFAULT nextval('public.ln_licenses_id_seq'::regclass);


--
-- Name: ln_user_actions id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ln_user_actions ALTER COLUMN id SET DEFAULT nextval('public.ln_user_actions_id_seq'::regclass);


--
-- Name: ln_user_permissions id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ln_user_permissions ALTER COLUMN id SET DEFAULT nextval('public.user_permissions_id_seq'::regclass);


--
-- Name: ln_users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ln_users ALTER COLUMN id SET DEFAULT nextval('public.ln_users_id_seq'::regclass);


--
-- Data for Name: ln_customers; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: ln_licenses; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: ln_user_actions; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.ln_user_actions (id, action, description, created_at, updated_at, deleted_at) VALUES (3, 'delete', '', '2022-04-19 10:21:07.563603', '1000-01-01 00:00:00', '1000-01-01 00:00:00');
INSERT INTO public.ln_user_actions (id, action, description, created_at, updated_at, deleted_at) VALUES (2, 'update', '', '2022-04-19 10:21:07.563603', '1000-01-01 00:00:00', '1000-01-01 00:00:00');
INSERT INTO public.ln_user_actions (id, action, description, created_at, updated_at, deleted_at) VALUES (1, 'create', '', '2022-04-19 10:21:07.563603', '1000-01-01 00:00:00', '1000-01-01 00:00:00');
INSERT INTO public.ln_user_actions (id, action, description, created_at, updated_at, deleted_at) VALUES (4, 'get', '', '2022-04-19 10:21:07.563603', '1000-01-01 00:00:00', '1000-01-01 00:00:00');


--
-- Data for Name: ln_user_permissions; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.ln_user_permissions (id, user_id, action_id, product_type, created_at, updated_at, deleted_at) VALUES (78, 44, 1, 'courier', '2022-04-21 16:48:37.353374', NULL, NULL);
INSERT INTO public.ln_user_permissions (id, user_id, action_id, product_type, created_at, updated_at, deleted_at) VALUES (79, 44, 2, 'courier', '2022-04-21 16:48:37.353374', NULL, NULL);
INSERT INTO public.ln_user_permissions (id, user_id, action_id, product_type, created_at, updated_at, deleted_at) VALUES (80, 44, 3, 'courier', '2022-04-21 16:48:37.353374', NULL, NULL);
INSERT INTO public.ln_user_permissions (id, user_id, action_id, product_type, created_at, updated_at, deleted_at) VALUES (81, 44, 4, 'courier', '2022-04-21 16:48:37.353374', NULL, NULL);


--
-- Data for Name: ln_users; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.ln_users (id, code, role, title, auth_token, description, created_at, updated_at, deleted_at) VALUES (44, '909fdfea-fff1-4dd5-97f7-612cf9840b82', 'admin', 'Admin', '54d1ba805e2a4891aeac9299b618945e', '', '2022-04-21 16:48:37.351955', NULL, NULL);


--
-- Name: ln_customers_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.ln_customers_id_seq', 25, true);


--
-- Name: ln_licenses_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.ln_licenses_id_seq', 32, true);


--
-- Name: ln_user_actions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.ln_user_actions_id_seq', 4, true);


--
-- Name: ln_users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.ln_users_id_seq', 47, true);


--
-- Name: user_permissions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.user_permissions_id_seq', 108, true);


--
-- Name: ln_customers ln_customers_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ln_customers
    ADD CONSTRAINT ln_customers_pk PRIMARY KEY (id);


--
-- Name: ln_customers ln_customers_pk_2; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ln_customers
    ADD CONSTRAINT ln_customers_pk_2 UNIQUE (user_id, code);


--
-- Name: ln_licenses ln_licenses_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ln_licenses
    ADD CONSTRAINT ln_licenses_pk PRIMARY KEY (id);


--
-- Name: ln_user_actions ln_user_actions_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ln_user_actions
    ADD CONSTRAINT ln_user_actions_pk PRIMARY KEY (id);


--
-- Name: ln_user_permissions ln_user_permissions_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ln_user_permissions
    ADD CONSTRAINT ln_user_permissions_pk PRIMARY KEY (id);


--
-- Name: ln_users ln_users_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ln_users
    ADD CONSTRAINT ln_users_pk PRIMARY KEY (id);


--
-- Name: ln_licenses_code_customer_id_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX ln_licenses_code_customer_id_uindex ON public.ln_licenses USING btree (code, customer_id);


--
-- Name: ln_licenses_license_key_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX ln_licenses_license_key_uindex ON public.ln_licenses USING btree (license_key);


--
-- Name: ln_user_actions_action_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX ln_user_actions_action_uindex ON public.ln_user_actions USING btree (action);


--
-- Name: ln_user_permissions_user_id, action_id_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX "ln_user_permissions_user_id, action_id_uindex" ON public.ln_user_permissions USING btree (user_id, action_id);


--
-- Name: ln_users_auth_token_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX ln_users_auth_token_uindex ON public.ln_users USING btree (auth_token);


--
-- Name: ln_users_code_deleted_at_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX ln_users_code_deleted_at_uindex ON public.ln_users USING btree (code, deleted_at);


--
-- Name: ln_customers ln_customers_ln_users_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ln_customers
    ADD CONSTRAINT ln_customers_ln_users_id_fk FOREIGN KEY (user_id) REFERENCES public.ln_users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: ln_licenses ln_licenses_ln_customers_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ln_licenses
    ADD CONSTRAINT ln_licenses_ln_customers_id_fk FOREIGN KEY (customer_id) REFERENCES public.ln_customers(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: ln_user_permissions ln_user_permissions_ln_user_actions_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ln_user_permissions
    ADD CONSTRAINT ln_user_permissions_ln_user_actions_id_fk FOREIGN KEY (action_id) REFERENCES public.ln_user_actions(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: ln_user_permissions ln_user_permissions_ln_users_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ln_user_permissions
    ADD CONSTRAINT ln_user_permissions_ln_users_id_fk FOREIGN KEY (user_id) REFERENCES public.ln_users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

