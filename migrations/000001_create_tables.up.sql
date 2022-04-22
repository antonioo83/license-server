--
-- PostgreSQL database dump
--

-- Dumped from database version 14.2
-- Dumped by pg_dump version 14.2

-- Started on 2022-04-22 15:16:28

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
-- TOC entry 216 (class 1259 OID 25001)
-- Name: ln_customers; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ln_customers (
                                     id integer NOT NULL,
                                     user_id integer NOT NULL,
                                     code character varying(32) DEFAULT ''::character varying NOT NULL,
                                     type character varying(20) DEFAULT 'service'::character varying NOT NULL,
                                     title character varying(100) DEFAULT ''::character varying NOT NULL,
                                     inn character varying(12) DEFAULT ''::character varying NOT NULL,
                                     description character varying(32) DEFAULT ''::character varying NOT NULL,
                                     created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
                                     updated_at timestamp without time zone,
                                     deleted_at timestamp without time zone
);


ALTER TABLE public.ln_customers OWNER TO postgres;

--
-- TOC entry 215 (class 1259 OID 25000)
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
-- TOC entry 3382 (class 0 OID 0)
-- Dependencies: 215
-- Name: ln_customers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.ln_customers_id_seq OWNED BY public.ln_customers.id;


--
-- TOC entry 217 (class 1259 OID 25013)
-- Name: ln_licenses; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ln_licenses (
                                    id integer NOT NULL,
                                    customer_id integer NOT NULL,
                                    product_type character varying(20) DEFAULT ''::character varying NOT NULL,
                                    callback_url character varying(500) DEFAULT ''::character varying NOT NULL,
                                    count integer DEFAULT 1 NOT NULL,
                                    license_key character varying(500) DEFAULT ''::character varying NOT NULL,
                                    registration_at timestamp without time zone NOT NULL,
                                    activation_at timestamp without time zone,
                                    expiration_at timestamp without time zone NOT NULL,
                                    deleted_at timestamp without time zone,
                                    duration integer NOT NULL,
                                    description character varying(500) DEFAULT ''::character varying NOT NULL
);


ALTER TABLE public.ln_licenses OWNER TO postgres;

--
-- TOC entry 218 (class 1259 OID 25023)
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
-- TOC entry 3383 (class 0 OID 0)
-- Dependencies: 218
-- Name: ln_licenses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.ln_licenses_id_seq OWNED BY public.ln_licenses.id;


--
-- TOC entry 212 (class 1259 OID 24983)
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
-- TOC entry 211 (class 1259 OID 24982)
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
-- TOC entry 3384 (class 0 OID 0)
-- Dependencies: 211
-- Name: ln_user_actions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.ln_user_actions_id_seq OWNED BY public.ln_user_actions.id;


--
-- TOC entry 214 (class 1259 OID 24994)
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
-- TOC entry 209 (class 1259 OID 24966)
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
-- TOC entry 210 (class 1259 OID 24969)
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
-- TOC entry 3385 (class 0 OID 0)
-- Dependencies: 210
-- Name: ln_users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.ln_users_id_seq OWNED BY public.ln_users.id;


--
-- TOC entry 213 (class 1259 OID 24993)
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
-- TOC entry 3386 (class 0 OID 0)
-- Dependencies: 213
-- Name: user_permissions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.user_permissions_id_seq OWNED BY public.ln_user_permissions.id;


--
-- TOC entry 3197 (class 2604 OID 25004)
-- Name: ln_customers id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ln_customers ALTER COLUMN id SET DEFAULT nextval('public.ln_customers_id_seq'::regclass);


--
-- TOC entry 3209 (class 2604 OID 25024)
-- Name: ln_licenses id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ln_licenses ALTER COLUMN id SET DEFAULT nextval('public.ln_licenses_id_seq'::regclass);


--
-- TOC entry 3191 (class 2604 OID 24986)
-- Name: ln_user_actions id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ln_user_actions ALTER COLUMN id SET DEFAULT nextval('public.ln_user_actions_id_seq'::regclass);


--
-- TOC entry 3194 (class 2604 OID 24997)
-- Name: ln_user_permissions id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ln_user_permissions ALTER COLUMN id SET DEFAULT nextval('public.user_permissions_id_seq'::regclass);


--
-- TOC entry 3184 (class 2604 OID 24970)
-- Name: ln_users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ln_users ALTER COLUMN id SET DEFAULT nextval('public.ln_users_id_seq'::regclass);

--
-- TOC entry 3387 (class 0 OID 0)
-- Dependencies: 215
-- Name: ln_customers_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.ln_customers_id_seq', 1, false);


--
-- TOC entry 3388 (class 0 OID 0)
-- Dependencies: 218
-- Name: ln_licenses_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.ln_licenses_id_seq', 1, false);


--
-- TOC entry 3389 (class 0 OID 0)
-- Dependencies: 211
-- Name: ln_user_actions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.ln_user_actions_id_seq', 4, true);


--
-- TOC entry 3390 (class 0 OID 0)
-- Dependencies: 210
-- Name: ln_users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.ln_users_id_seq', 44, true);


--
-- TOC entry 3391 (class 0 OID 0)
-- Dependencies: 213
-- Name: user_permissions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.user_permissions_id_seq', 81, true);


--
-- TOC entry 3221 (class 2606 OID 25012)
-- Name: ln_customers ln_customers_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ln_customers
    ADD CONSTRAINT ln_customers_pk PRIMARY KEY (id);


--
-- TOC entry 3223 (class 2606 OID 25026)
-- Name: ln_licenses ln_licenses_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ln_licenses
    ADD CONSTRAINT ln_licenses_pk PRIMARY KEY (id);


--
-- TOC entry 3216 (class 2606 OID 24990)
-- Name: ln_user_actions ln_user_actions_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ln_user_actions
    ADD CONSTRAINT ln_user_actions_pk PRIMARY KEY (id);


--
-- TOC entry 3218 (class 2606 OID 25028)
-- Name: ln_user_permissions ln_user_permissions_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ln_user_permissions
    ADD CONSTRAINT ln_user_permissions_pk PRIMARY KEY (id);


--
-- TOC entry 3213 (class 2606 OID 24992)
-- Name: ln_users ln_users_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ln_users
    ADD CONSTRAINT ln_users_pk PRIMARY KEY (id);


--
-- TOC entry 3214 (class 1259 OID 25060)
-- Name: ln_user_actions_action_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX ln_user_actions_action_uindex ON public.ln_user_actions USING btree (action);


--
-- TOC entry 3219 (class 1259 OID 25059)
-- Name: ln_user_permissions_user_id, action_id_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX "ln_user_permissions_user_id, action_id_uindex" ON public.ln_user_permissions USING btree (user_id, action_id);


--
-- TOC entry 3210 (class 1259 OID 25058)
-- Name: ln_users_auth_token_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX ln_users_auth_token_uindex ON public.ln_users USING btree (auth_token);


--
-- TOC entry 3211 (class 1259 OID 25062)
-- Name: ln_users_code_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX ln_users_code_uindex ON public.ln_users USING btree (code);


--
-- TOC entry 3226 (class 2606 OID 25034)
-- Name: ln_customers ln_customers_ln_users_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ln_customers
    ADD CONSTRAINT ln_customers_ln_users_id_fk FOREIGN KEY (user_id) REFERENCES public.ln_users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- TOC entry 3227 (class 2606 OID 25029)
-- Name: ln_licenses ln_licenses_ln_customers_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ln_licenses
    ADD CONSTRAINT ln_licenses_ln_customers_id_fk FOREIGN KEY (customer_id) REFERENCES public.ln_customers(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- TOC entry 3225 (class 2606 OID 25049)
-- Name: ln_user_permissions ln_user_permissions_ln_user_actions_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ln_user_permissions
    ADD CONSTRAINT ln_user_permissions_ln_user_actions_id_fk FOREIGN KEY (action_id) REFERENCES public.ln_user_actions(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- TOC entry 3224 (class 2606 OID 25044)
-- Name: ln_user_permissions ln_user_permissions_ln_users_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ln_user_permissions
    ADD CONSTRAINT ln_user_permissions_ln_users_id_fk FOREIGN KEY (user_id) REFERENCES public.ln_users(id) ON UPDATE CASCADE ON DELETE CASCADE;


-- Completed on 2022-04-22 15:16:28

--
-- PostgreSQL database dump complete
--

