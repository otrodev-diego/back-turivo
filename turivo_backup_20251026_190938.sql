pg_dump: last built-in OID is 16383
pg_dump: reading extensions
pg_dump: identifying extension members
pg_dump: reading schemas
pg_dump: reading user-defined tables
pg_dump: reading user-defined functions
pg_dump: reading user-defined types
pg_dump: reading procedural languages
pg_dump: reading user-defined aggregate functions
pg_dump: reading user-defined operators
pg_dump: reading user-defined access methods
pg_dump: reading user-defined operator classes
pg_dump: reading user-defined operator families
pg_dump: reading user-defined text search parsers
pg_dump: reading user-defined text search templates
pg_dump: reading user-defined text search dictionaries
pg_dump: reading user-defined text search configurations
pg_dump: reading user-defined foreign-data wrappers
pg_dump: reading user-defined foreign servers
pg_dump: reading default privileges
pg_dump: reading user-defined collations
pg_dump: reading user-defined conversions
pg_dump: reading type casts
pg_dump: reading transforms
pg_dump: reading table inheritance information
pg_dump: reading event triggers
pg_dump: finding extension tables
pg_dump: finding inheritance relationships
pg_dump: reading column info for interesting tables
pg_dump: finding table default expressions
pg_dump: finding table check constraints
pg_dump: flagging inherited columns in subtables
pg_dump: reading partitioning data
pg_dump: reading indexes
pg_dump: flagging indexes in partitioned tables
pg_dump: reading extended statistics
pg_dump: reading constraints
pg_dump: reading triggers
pg_dump: reading rewrite rules
pg_dump: reading policies
pg_dump: reading row-level security policies
pg_dump: reading publications
pg_dump: reading publication membership of tables
pg_dump: reading publication membership of schemas
pg_dump: reading subscriptions
pg_dump: reading large objects
pg_dump: reading dependency data
pg_dump: saving encoding = UTF8
pg_dump: saving standard_conforming_strings = on
pg_dump: saving search_path = 
pg_dump: creating EXTENSION "uuid-ossp"
pg_dump: creating COMMENT "EXTENSION "uuid-ossp""
pg_dump: creating TYPE "public.background_check_status"
pg_dump: creating TYPE "public.company_sector"
pg_dump: creating TYPE "public.company_status"
pg_dump: creating TYPE "public.driver_status"
pg_dump: creating TYPE "public.language"
pg_dump: creating TYPE "public.license_class"
pg_dump: creating TYPE "public.payment_gateway"
pg_dump: creating TYPE "public.payment_status"
pg_dump: creating TYPE "public.request_status"
pg_dump: creating TYPE "public.reservation_status"
pg_dump: creating TYPE "public.user_role"
pg_dump: creating TYPE "public.user_status"
pg_dump: creating TYPE "public.vehicle_status"
pg_dump: creating TYPE "public.vehicle_type"
pg_dump: creating FUNCTION "public.update_updated_at_column()"
--
-- PostgreSQL database dump
--

\restrict D39Vw3DHTIZI3zqNf08wH1To83VA1fRb6so7TcHD8PYJAEXNACTVIukfiEpXEWH

-- Dumped from database version 15.14
-- Dumped by pg_dump version 16.10 (Homebrew)

-- Started on 2025-10-26 19:09:39 -03

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
-- TOC entry 2 (class 3079 OID 16396)
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- TOC entry 3712 (class 0 OID 0)
-- Dependencies: 2
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


--
-- TOC entry 889 (class 1247 OID 16464)
-- Name: background_check_status; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.background_check_status AS ENUM (
    'APPROVED',
    'PENDING',
    'REJECTED'
);


--
-- TOC entry 880 (class 1247 OID 16430)
-- Name: company_sector; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.company_sector AS ENUM (
    'HOTEL',
    'MINERIA',
    'TURISMO'
);


--
-- TOC entry 877 (class 1247 OID 16424)
-- Name: company_status; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.company_status AS ENUM (
    'ACTIVE',
    'SUSPENDED'
);


--
-- TOC entry 883 (class 1247 OID 16438)
-- Name: driver_status; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.driver_status AS ENUM (
    'ACTIVE',
    'INACTIVE'
);


--
-- TOC entry 907 (class 1247 OID 16516)
-- Name: language; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.language AS ENUM (
    'es',
    'en',
    'pt',
    'fr'
);


--
-- TOC entry 886 (class 1247 OID 16444)
-- Name: license_class; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.license_class AS ENUM (
    'A1',
    'A2',
    'A3',
    'A4',
    'A5',
    'B',
    'C',
    'D',
    'E'
);


--
-- TOC entry 901 (class 1247 OID 16504)
-- Name: payment_gateway; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.payment_gateway AS ENUM (
    'WEBPAY_PLUS'
);


--
-- TOC entry 904 (class 1247 OID 16508)
-- Name: payment_status; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.payment_status AS ENUM (
    'APPROVED',
    'REJECTED',
    'PENDING'
);


--
-- TOC entry 895 (class 1247 OID 16482)
-- Name: request_status; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.request_status AS ENUM (
    'PENDIENTE',
    'ASIGNADA',
    'EN_RUTA',
    'COMPLETADA',
    'CANCELADA'
);


--
-- TOC entry 898 (class 1247 OID 16494)
-- Name: reservation_status; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.reservation_status AS ENUM (
    'ACTIVA',
    'PROGRAMADA',
    'COMPLETADA',
    'CANCELADA'
);


--
-- TOC entry 871 (class 1247 OID 16408)
-- Name: user_role; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.user_role AS ENUM (
    'ADMIN',
    'USER',
    'DRIVER',
    'COMPANY'
);


--
-- TOC entry 874 (class 1247 OID 16418)
-- Name: user_status; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.user_status AS ENUM (
    'ACTIVE',
    'BLOCKED'
);


--
-- TOC entry 955 (class 1247 OID 16822)
-- Name: vehicle_status; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.vehicle_status AS ENUM (
    'AVAILABLE',
    'ASSIGNED',
    'MAINTENANCE',
    'INACTIVE'
);


--
-- TOC entry 892 (class 1247 OID 16472)
-- Name: vehicle_type; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.vehicle_type AS ENUM (
    'BUS',
    'VAN',
    'SEDAN',
    'SUV'
);


--
-- TOC entry 244 (class 1255 OID 16751)
-- Name: update_updated_at_column(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.update_updated_at_column() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = NOW();
  pg_dump: creating FUNCTION "public.update_vehicle_status_on_assignment()"
pg_dump: creating TABLE "public.companies"
pg_dump: creating TABLE "public.driver_availability"
pg_dump: creating TABLE "public.driver_background_checks"
pg_dump: creating TABLE "public.driver_feedback"
pg_dump: creating TABLE "public.driver_licenses"
pg_dump: creating TABLE "public.drivers"
pg_dump: creating TABLE "public.hotels"
  RETURN NEW;
END;
$$;


--
-- TOC entry 245 (class 1255 OID 16838)
-- Name: update_vehicle_status_on_assignment(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.update_vehicle_status_on_assignment() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF NEW.driver_id IS NOT NULL AND (OLD.driver_id IS NULL OR OLD.driver_id != NEW.driver_id) THEN
        NEW.status = 'ASSIGNED';
    ELSIF NEW.driver_id IS NULL AND OLD.driver_id IS NOT NULL THEN
        NEW.status = 'AVAILABLE';
    END IF;
    RETURN NEW;
END;
$$;


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 217 (class 1259 OID 16539)
-- Name: companies; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.companies (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying(255) NOT NULL,
    rut character varying(50) NOT NULL,
    contact_email character varying(255) NOT NULL,
    status public.company_status DEFAULT 'ACTIVE'::public.company_status NOT NULL,
    sector public.company_sector NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- TOC entry 224 (class 1259 OID 16628)
-- Name: driver_availability; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.driver_availability (
    driver_id character varying(20) NOT NULL,
    regions jsonb DEFAULT '[]'::jsonb NOT NULL,
    days jsonb DEFAULT '[]'::jsonb NOT NULL,
    time_ranges jsonb DEFAULT '[]'::jsonb NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- TOC entry 221 (class 1259 OID 16586)
-- Name: driver_background_checks; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.driver_background_checks (
    driver_id character varying(20) NOT NULL,
    status public.background_check_status DEFAULT 'PENDING'::public.background_check_status NOT NULL,
    file_url character varying(500),
    checked_at timestamp with time zone
);


--
-- TOC entry 231 (class 1259 OID 24657)
-- Name: driver_feedback; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.driver_feedback (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    driver_id character varying(20) NOT NULL,
    reservation_id character varying(20) NOT NULL,
    rating numeric(2,1) NOT NULL,
    comment text,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT driver_feedback_rating_check CHECK (((rating >= 1.0) AND (rating <= 5.0)))
);


--
-- TOC entry 220 (class 1259 OID 16574)
-- Name: driver_licenses; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.driver_licenses (
    driver_id character varying(20) NOT NULL,
    number character varying(100) NOT NULL,
    class public.license_class NOT NULL,
    issued_at date,
    expires_at date,
    file_url character varying(500)
);


--
-- TOC entry 219 (class 1259 OID 16562)
-- Name: drivers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.drivers (
    id character varying(20) NOT NULL,
    first_name character varying(255) NOT NULL,
    last_name character varying(255) NOT NULL,
    rut_or_dni character varying(50) NOT NULL,
    birth_date date,
    phone character varying(50),
    email character varying(255),
    photo_url character varying(500),
    status public.driver_status DEFAULT 'ACTIVE'::public.driver_status NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    user_id uuid,
    company_id uuid,
    vehicle_id uuid
);


--
-- TOC entry 218 (class 1259 OID 16552)
-- Name: hotels; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.hotels (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying(255) NOT NULL,
    city character varying(255) NOT NULL,
    contact_email character varying(255) NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);
pg_dump: creating TABLE "public.migrations"
pg_dump: creating SEQUENCE "public.migrations_id_seq"
pg_dump: creating SEQUENCE OWNED BY "public.migrations_id_seq"
pg_dump: creating TABLE "public.payments"
pg_dump: creating TABLE "public.refresh_tokens"
pg_dump: creating TABLE "public.registration_tokens"
pg_dump: creating COMMENT "public.COLUMN registration_tokens.company_profile"
pg_dump: creating TABLE "public.requests"
pg_dump: creating TABLE "public.reservation_timeline"
pg_dump: creating TABLE "public.reservations"


--
-- TOC entry 233 (class 1259 OID 24728)
-- Name: migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.migrations (
    id integer NOT NULL,
    migration character varying(255) NOT NULL,
    batch integer NOT NULL
);


--
-- TOC entry 232 (class 1259 OID 24727)
-- Name: migrations_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.migrations_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 3713 (class 0 OID 0)
-- Dependencies: 232
-- Name: migrations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.migrations_id_seq OWNED BY public.migrations.id;


--
-- TOC entry 228 (class 1259 OID 16703)
-- Name: payments; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.payments (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    reservation_id character varying(20) NOT NULL,
    gateway public.payment_gateway NOT NULL,
    amount numeric(12,2) NOT NULL,
    currency character varying(3) DEFAULT 'CLP'::character varying NOT NULL,
    status public.payment_status DEFAULT 'PENDING'::public.payment_status NOT NULL,
    transaction_ref character varying(255),
    payload jsonb,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- TOC entry 229 (class 1259 OID 16764)
-- Name: refresh_tokens; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.refresh_tokens (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    token text NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- TOC entry 230 (class 1259 OID 16787)
-- Name: registration_tokens; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.registration_tokens (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    token character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    org_id uuid,
    role public.user_role NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    used boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT now(),
    company_profile character varying(50)
);


--
-- TOC entry 3714 (class 0 OID 0)
-- Dependencies: 230
-- Name: COLUMN registration_tokens.company_profile; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.registration_tokens.company_profile IS 'Profile type for COMPANY role registrations. Valid values: COMPANY_ADMIN, COMPANY_USER';


--
-- TOC entry 225 (class 1259 OID 16644)
-- Name: requests; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.requests (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    hotel_id uuid,
    company_id uuid,
    fecha timestamp with time zone NOT NULL,
    origin jsonb NOT NULL,
    destination jsonb NOT NULL,
    pax integer NOT NULL,
    vehicle_type public.vehicle_type NOT NULL,
    language public.language,
    status public.request_status DEFAULT 'PENDIENTE'::public.request_status NOT NULL,
    assigned_driver_id character varying(20),
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT check_org_id CHECK (((hotel_id IS NOT NULL) OR (company_id IS NOT NULL))),
    CONSTRAINT requests_pax_check CHECK ((pax > 0))
);


--
-- TOC entry 227 (class 1259 OID 16688)
-- Name: reservation_timeline; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.reservation_timeline (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    reservation_id character varying(20) NOT NULL,
    title character varying(255) NOT NULL,
    description text NOT NULL,
    at timestamp with time zone NOT NULL,
    variant character varying(50) DEFAULT 'default'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- TOC entry 226 (class 1259 OID 16672)
-- Name: reservations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.reservations (
    id character varying(20) NOT NULL,pg_dump: creating TABLE "public.schema_migrations"
pg_dump: creating TABLE "public.users"
pg_dump: creating COMMENT "public.COLUMN users.company_profile"
pg_dump: creating TABLE "public.vehicle_photos"
pg_dump: creating TABLE "public.vehicles"
pg_dump: creating DEFAULT "public.migrations id"
pg_dump: processing data for table "public.companies"
pg_dump: dumping contents of table "public.companies"

    user_id uuid,
    org_id uuid,
    pickup character varying(500) NOT NULL,
    destination character varying(500) NOT NULL,
    datetime timestamp with time zone NOT NULL,
    passengers integer NOT NULL,
    status public.reservation_status DEFAULT 'ACTIVA'::public.reservation_status NOT NULL,
    amount numeric(12,2),
    notes text,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    assigned_driver_id character varying(20),
    distance_km numeric(10,2),
    arrived_on_time boolean DEFAULT true,
    CONSTRAINT reservations_passengers_check CHECK ((passengers > 0))
);


--
-- TOC entry 215 (class 1259 OID 16389)
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


--
-- TOC entry 216 (class 1259 OID 16525)
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    password_hash character varying(255) NOT NULL,
    role public.user_role DEFAULT 'USER'::public.user_role NOT NULL,
    status public.user_status DEFAULT 'ACTIVE'::public.user_status NOT NULL,
    org_id uuid,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    company_profile character varying(50)
);


--
-- TOC entry 3715 (class 0 OID 0)
-- Dependencies: 216
-- Name: COLUMN users.company_profile; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.users.company_profile IS 'Profile type for COMPANY role users. Valid values: COMPANY_ADMIN (full access), COMPANY_USER (limited access)';


--
-- TOC entry 223 (class 1259 OID 16614)
-- Name: vehicle_photos; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.vehicle_photos (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    vehicle_id uuid NOT NULL,
    url character varying(500) NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- TOC entry 222 (class 1259 OID 16599)
-- Name: vehicles; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.vehicles (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    driver_id character varying(20),
    type public.vehicle_type NOT NULL,
    brand character varying(100) NOT NULL,
    model character varying(100) NOT NULL,
    year integer,
    plate character varying(20),
    vin character varying(100),
    color character varying(50),
    insurance_policy character varying(100),
    insurance_expires_at date,
    inspection_expires_at date,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    status public.vehicle_status DEFAULT 'AVAILABLE'::public.vehicle_status NOT NULL,
    capacity integer,
    CONSTRAINT vehicles_capacity_check CHECK (((capacity > 0) AND (capacity <= 60)))
);


--
-- TOC entry 3432 (class 2604 OID 24731)
-- Name: migrations id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.migrations ALTER COLUMN id SET DEFAULT nextval('public.migrations_id_seq'::regclass);


--
-- TOC entry 3690 (class 0 OID 16539)
-- Dependencies: 217
-- Data for Name: companies; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.companies (id, name, rut, contact_email, status, sector, created_at, updated_at) FROM stdin;
dd28931b-c0d1-4c8b-a1b5-85ada9390657	Turismo Andes	76.234.567-8	info@turismoandes.cl	ACTIVE	TURISMO	2025-10-02 20:16:41.204462+00	2025-10-02 20:16:41.204462+00
6c2b7ad0-2629-4b1b-9e3e-0d90bbf1a212	Minera del Norte	96.345.678-9	contacto@mineranorte.cl	ACTIVE	MINERIA	2025-10-02 20:16:41.204462+00	2025-10-02 20:16:41.204462+00
17a55804-4644-46ea-a28f-5c4efc8fbf67	Turivo ok	76.123.456-7	contacto@turivo.com	ACTIVE	TURISMO	2025-10-02 20:16:41.204462+00	2025-10-11 01:27:02.539805+00
\.


--
-- TOC entry 3697 (class 0 OID 16628)
-- Dependencipg_dump: processing data for table "public.driver_availability"
pg_dump: dumping contents of table "public.driver_availability"
pg_dump: processing data for table "public.driver_background_checks"
pg_dump: dumping contents of table "public.driver_background_checks"
pg_dump: processing data for table "public.driver_feedback"
pg_dump: dumping contents of table "public.driver_feedback"
pg_dump: processing data for table "public.driver_licenses"
pg_dump: dumping contents of table "public.driver_licenses"
pg_dump: processing data for table "public.drivers"
pg_dump: dumping contents of table "public.drivers"
pg_dump: processing data for table "public.hotels"
pg_dump: dumping contents of table "public.hotels"
pg_dump: processing data for table "public.migrations"
pg_dump: dumping contents of table "public.migrations"
pg_dump: processing data for table "public.payments"
pg_dump: dumping contents of table "public.payments"
pg_dump: processing data for table "public.refresh_tokens"
pg_dump: dumping contents of table "public.refresh_tokens"
es: 224
-- Data for Name: driver_availability; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.driver_availability (driver_id, regions, days, time_ranges, updated_at) FROM stdin;
\.


--
-- TOC entry 3694 (class 0 OID 16586)
-- Dependencies: 221
-- Data for Name: driver_background_checks; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.driver_background_checks (driver_id, status, file_url, checked_at) FROM stdin;
\.


--
-- TOC entry 3704 (class 0 OID 24657)
-- Dependencies: 231
-- Data for Name: driver_feedback; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.driver_feedback (id, driver_id, reservation_id, rating, comment, created_at, updated_at) FROM stdin;
\.


--
-- TOC entry 3693 (class 0 OID 16574)
-- Dependencies: 220
-- Data for Name: driver_licenses; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.driver_licenses (driver_id, number, class, issued_at, expires_at, file_url) FROM stdin;
\.


--
-- TOC entry 3692 (class 0 OID 16562)
-- Dependencies: 219
-- Data for Name: drivers; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.drivers (id, first_name, last_name, rut_or_dni, birth_date, phone, email, photo_url, status, created_at, updated_at, user_id, company_id, vehicle_id) FROM stdin;
DRV001	Juan	Pérez	12345678-9	\N	+56912345678	juan.perez@turivo.com	\N	ACTIVE	2025-10-24 03:40:22.231718+00	2025-10-24 03:40:22.231718+00	\N	\N	\N
DRV002	María	González	98765432-1	\N	+56987654321	maria.gonzalez@turivo.com	\N	ACTIVE	2025-10-24 03:40:22.231718+00	2025-10-24 03:40:22.231718+00	\N	\N	\N
DRV-1761363000000	Nuevo	Conductor	987654321	\N	+56987654321	nuevo@conductor.com	\N	ACTIVE	2025-10-25 03:35:03.008869+00	2025-10-25 03:35:03.008869+00	\N	\N	\N
DRV003	Carlos	Rodríguez	11223344-5	\N	+56911223344	carlos.rodriguez@turivo.com	\N	ACTIVE	2025-10-24 03:40:22.231718+00	2025-10-26 18:43:07.941679+00	\N	\N	\N
DRV-1761368724331	Diego	Jara	12.345.678-91111	\N	+56912345678	diego@otrodev.cl	\N	ACTIVE	2025-10-25 05:05:24.433734+00	2025-10-26 20:50:06.537006+00	\N	\N	\N
\.


--
-- TOC entry 3691 (class 0 OID 16552)
-- Dependencies: 218
-- Data for Name: hotels; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.hotels (id, name, city, contact_email, created_at, updated_at) FROM stdin;
56e9a543-468d-4cf5-8d06-87e91c812a30	Hotel Miramar	Valparaíso	reservas@hotelmiramar.cl	2025-10-02 20:16:41.204462+00	2025-10-02 20:16:41.204462+00
32ebf7c0-b2e7-4fe1-b73a-5a10a76f58c1	Hotel Andes	Santiago	contacto@hotelandes.cl	2025-10-02 20:16:41.204462+00	2025-10-02 20:16:41.204462+00
284ed600-0f5d-49a5-a844-d1ee2ae25e09	Hotel Patagonia	Puerto Montt	info@hotelpatagonia.cl	2025-10-02 20:16:41.204462+00	2025-10-02 20:16:41.204462+00
\.


--
-- TOC entry 3706 (class 0 OID 24728)
-- Dependencies: 233
-- Data for Name: migrations; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.migrations (id, migration, batch) FROM stdin;
\.


--
-- TOC entry 3701 (class 0 OID 16703)
-- Dependencies: 228
-- Data for Name: payments; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.payments (id, reservation_id, gateway, amount, currency, status, transaction_ref, payload, created_at) FROM stdin;
5ac4db94-90fa-4e23-aaaa-b02f5b278cce	RSV-1003	WEBPAY_PLUS	250000.00	CLP	APPROVED	WP_1759436201	{"vci": "TSY", "status": "AUTHORIZED"}	2025-10-02 20:16:41.204462+00
\.


--
-- TOC entry 3702 (class 0 OID 16764)
-- Dependencies: 229
-- Data for Name: refresh_tokens; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.refresh_tokens (id, user_id, token, expires_at, created_at) FROM stdin;
888e8078-862d-4ccf-8bf6-4d823f68982b	4999dcea-9508-4452-9c5a-bbaf201808ff	73efdb056953d0d3f68f65e5629d2c75c76deb85638babfc2b3e1316c1dbe2a1	2025-10-09 20:39:12.341746+00	2025-10-02 20:39:12.343234+00
a5958f78-3219-4d77-a1d5-5b39c3b968bb	4999dcea-9508-4452-9c5a-bbaf201808ff	954ca769e3efb2ce4d2cebb92a964abaa2bfd0410aa43fc46945fd3c43d317f8	2025-10-10 00:55:41.273781+00	2025-10-03 00:55:41.273592+00
e04a31d7-767e-4775-bd32-46fee4b58a56	4999dcea-9508-4452-9c5a-bbaf201808ff	c483780262cb5bb3d3755061ae6243701d7da620d310f2c7f05eea156cde7219	2025-10-10 01:34:18.742623+00	2025-10-03 01:34:18.744166+00
74dcbef5-b8ee-4b31-a773-6f5ed9ccb0d9	4999dcea-9508-4452-9c5a-bbaf201808ff	10fb34d30f0959a96090e83391cfb294b5330ffe5ba9dea1ad7c00d9d294215e	2025-10-10 01:51:06.194806+00	2025-10-03 01:51:06.195774+00
1f5eef2e-9adf-4628-a26a-3a0df74e1973	4999dcea-9508-4452-9c5a-bbaf201808ff	cf5e9a8dfb47ba77f2d04172b93ccdc62d56a26a2e3d52f706e9dbbce30f807c	2025-10-10 02:06:29.472506+00	2025-10-03 02:06:29.473491+00
6c9ca5be-c52d-4515-8b3b-86b8bb146e5b	4999dcea-9508-4452-9c5a-bbaf201808ff	1963ea4fae8dd2bda464148f04484b794d5c820181cc757d88470d46df1fc6f2	2025-10-10 02:23:03.973171+00	2025-10-03 02:23:03.972545+00
2e969c41-057c-49cf-b885-69748817bd02	4999dcea-9508-4452-9c5a-bbaf201808ff	59edb83abde8dc5bb4237eed4afcb2761b1089164fd117fd01f19252e97ba27a	2025-10-10 02:54:00.271295+00	2025-10-03 02:54:00.270417+00
abde2bd0-64eb-4e1b-be58-c2f219d1f646	4999dcea-9508-4452-9c5a-bbaf201808ff	67096be0a53a690d7522c5d9468c5aa0bbeb9cacd289bd4a749abd5530d94884	2025-10-10 03:16:34.428915+00	2025-10-03 03:16:34.429511+00
c99bf2ac-359b-44de-8bf9-4c821570f43d	4999dcea-9508-4452-9c5a-bbaf201808ff	b0fcecf877b0983eb897eb062de3d57561208e28e82862e449a4ad675056404c	2025-10-18 01:14:22.289027+00	2025-10-11 01:14:22.289694+00
83655d44-c443-40b5-977f-4f0376a26034	4999dcea-9508-4452-9c5a-bbaf201808ff	7f4ca630aab927a0eb9e92b7652e60209822096f999667852628c03f924b9d12	2025-10-18 01:29:31.034107+00	2025-10-11 01:29:31.035496+00
e3761b2d-edbe-483d-ad9c-d9c0a92373c6	4999dcea-9508-4452-9c5a-bbaf201808ff	515c1474b713687019d073cd467d62e03b9bde61c4123a546c5af6c5d41993b1	2025-10-18 02:56:02.36145+00	2025-10-11 02:56:02.361933+00
e8830b1f-eaf1-4840-8fe8-7a44bad9f82e	4999dcea-9508-4452-9c5a-bbaf201808ff	9e806743f645b14c884cc85b2ffcee58e96cf6118ac39c1356042d1b9afb8a5e	2025-10-31 03:25:37.131321+00	2025-10-24 03:25:37.133097+00
c9d17b2d-b306-492d-9950-24cd5d03aadb	4999dcea-9508-4452-9c5a-bbaf201808ff	f830af8d65d802fa2c35c45e3f6a595ae18e086116c0e6fb2f285fa298fbed2f	2025-10-31 03:32:31.952758+00	2025-10-24 03:32:31.95384+00
a25f774e-17df-46cd-b108-7195ed4a8872	4999dcea-9508-4452-9c5a-bbaf201808ff	cb1da0ece8eed7c3be4a69bcb445f81bcd1c154e51d6bf898be9baf4d831be51	2025-10-31 03:44:09.798653+00	2025-10-24 03:44:09.798692+00
e30b13c8-8d46-44fe-8bf5-e2d7a2a0b90f	4999dcea-9508-4452-9c5a-bbaf201808ff	0b08d4498d942841db001ada8759952102c4d670b6c88f20a2506fafd9f62f0a	2025-10-31 03:45:36.798512+00	2025-10-24 03:45:36.801595+00
0bb47bdd-bd64-4b1b-b096-a6c610cc2e2b	4999dcea-9508-4452-9c5a-bbaf201808ff	752b02292a9700d29fe412cdf5536e6f6fac3e8ceffdbd5cec767cde7e73791f	2025-10-31 03:46:00.146589+00	2025-10-24 03:46:00.146229+00
367b1d33-60aa-4c02-94b4-858ac1fb3cbc	4999dcea-9508-4452-9c5a-bbaf201808ff	8dbc0b7e71c43a95582c66c1662110852f29fdf3b55e331d4450f5fa5420c6b0	2025-10-31 04:22:05.444792+00	2025-10-24 04:22:05.444696+00
4c97f03a-6924-47fa-a3d5-46b0cede59ed	4999dcea-9508-4452-9c5a-bbaf201808ff	b0868c4dbe242decbe6bfc84b398388e07b5bbc6b2f73a320c51f792faf55f08	2025-10-31 04:25:52.315701+00	2025-10-24 04:25:52.316148+00
7da31644-2699-4f78-b313-99715745d829	4999dcea-9508-4452-9c5a-bbaf201808ff	061c5712db1f6d7843b2094c29bcd18972d2198d18af6e581694ce25f1837144	2025-10-31 04:27:32.371557+00	2025-10-24 04:27:32.372256+00
75e54a69-d790-4a97-8781-f08da9e92406	4999dcea-9508-4452-9c5a-bbaf201808ff	e312a282b5e3de9fc90097bcdcf48c0a3ddccac7f8834bd4a1b1c677f01fc896	2025-10-31 04:28:04.584157+00	2025-10-24 04:28:04.584837+00
d9b8d858-7961-42fb-b86a-7ddfa7e61fd8	4999dcea-9508-4452-9c5a-bbaf201808ff	67d899afe4aebda55c935bd717f35def1e6387d4a44972c7c05a0bd0618e6e10	2025-10-31 04:54:49.168712+00	2025-10-24 04:54:49.168614+00
1a343603-5886-4be9-ab38-94502958215b	4999dcea-9508-4452-9c5a-bbaf201808ff	021d1f2664a72fd034c35332d9ed3e342572bfe2e0dc66b58a55e4269c673f03	2025-10-31 04:58:24.198399+00	2025-10-24 04:58:24.198416+00
35009ba1-ccd7-4210-8d07-a683c61837c5	4999dcea-9508-4452-9c5a-bbaf201808ff	101a6ef2bee949e58142f693d1e41cc89e5a78a17ec14faa27af30fd44b49646	2025-10-31 04:59:49.455177+00	2025-10-24 04:59:49.455934+00
0f9e97db-7c5f-49c0-9cf8-8e5f4b85149e	4999dcea-9508-4452-9c5a-bbaf201808ff	e8a6e5e13ef84c0a9f299b3738fee79cbcb8e194dc552342c6c0f1b9b210d750	2025-10-31 05:04:36.328221+00	2025-10-24 05:04:36.329828+00
69f1bd0e-70da-438e-ab36-5aabd3418388	4999dcea-9508-4452-9c5a-bbaf201808ff	57f333bfb33be43bc93609b4edd7c4973846d5157523c90d11786d6e655bef6a	2025-10-31 05:14:16.004283+00	2025-10-24 05:14:16.004572+00
088d98a6-10df-4e63-b0c7-7f65963d4515	4999dcea-9508-4452-9c5a-bbaf201808ff	f2edc5324d31943a5cab73878bdcca87dc07ee40e70a12822e4e47f086786f5c	2025-10-31 05:22:47.791604+00	2025-10-24 05:22:47.792412+00
6d06110c-c211-48f3-af1f-66baa4c21c43	4999dcea-9508-4452-9c5a-bbaf201808ff	c58b9020904d32899151162d25ff3d5879c20e0c442b3efcd6b99b941734f704	2025-10-31 05:33:05.091928+00	2025-10-24 05:33:05.092219+00
f59c641a-3276-4b47-ad67-50cc2ec057e6	4999dcea-9508-4452-9c5a-bbaf201808ff	f943ae57b0876b0b3224c564a08ea3f36c6e4e43ac6a317678e9816e704093ee	2025-10-31 05:33:38.253574+00	2025-10-24 05:33:38.254991+00
98efd6f0-8359-4182-bbef-bc639e773cac	4999dcea-9508-4452-9c5a-bbaf201808ff	033d880ffcd73eefb04a29ff45bc771b233d412aeb79e6c45f8d4bb94daa0956	2025-10-31 05:33:43.395305+00	2025-10-24 05:33:43.395127+00
5ac64b8a-668b-4968-8094-6b6b68697275	4999dcea-9508-4452-9c5a-bbaf201808ff	623a997050ac82de8d5347bb9ff843a6e05f5b1161dab0f784b9092ede42acb3	2025-10-31 05:33:48.403774+00	2025-10-24 05:33:48.4033+00
0667142f-17bd-42a7-a18a-66d9e6514a89	4999dcea-9508-4452-9c5a-bbaf201808ff	9138d9863daaf0c812fa3e31262d39d3c5925b2ad3cfef83b1e0284c73e9a6ee	2025-10-31 05:34:33.789711+00	2025-10-24 05:34:33.789416+00
896a0077-758c-46af-a4e0-54719d092c99	4999dcea-9508-4452-9c5a-bbaf201808ff	552206dc245067739684e6c6a8a699c23ad9e74981dea8b6ee1d1ab712721fa6	2025-10-31 05:34:40.176171+00	2025-10-24 05:34:40.176159+00
08f63e14-be4e-4ce6-bd50-3c219beb6ff5	4999dcea-9508-4452-9c5a-bbaf201808ff	0190d50b4f20e3e21e1aef352c69b746413a91eacaaf0b0b43d73972b1000e34	2025-10-31 05:35:32.524746+00	2025-10-24 05:35:32.527567+00
232dc6a6-1405-42cc-b2f8-b77ce50c15ad	4999dcea-9508-4452-9c5a-bbaf201808ff	57fea6e05976ffcd8cb0ff896ea9178fad48e51b7052d4e96a1026b90c4585bc	2025-10-31 05:37:41.003622+00	2025-10-24 05:37:41.003933+00
51fbdc78-5903-4803-a192-320f644220b4	4999dcea-9508-4452-9c5a-bbaf201808ff	abbf6e8af012879854bd21d23449e77ecf5a3b717d9ba47686772721b169a4f5	2025-11-01 03:30:09.012485+00	2025-10-25 03:30:09.007083+00
185b4c50-00ad-4f6f-90a4-c9c6b1927750	4999dcea-9508-4452-9c5a-bbaf201808ff	cf1c3394dd80d781ebef54cbb63ae8788ece0ba90125ba045f93315053017df7	2025-11-01 03:34:36.599343+00	2025-10-25 03:34:36.600053+00
53850baf-0a11-43bd-a382-aa9481d06d9e	4999dcea-9508-4452-9c5a-bbaf201808ff	b4fb845db7a9d3ec29e683c2b582d0aee123b862724fc4e9ff69418f62fd852e	2025-11-01 03:34:57.704962+00	2025-10-25 03:34:57.70512+00
d6c39937-27ca-4abc-88e9-2053999285c2	4999dcea-9508-4452-9c5a-bbaf201808ff	1a0e0aaaf7a126478448cc426ac58be1f60f3483788d0e79e3cf4edc7fe4622f	2025-11-01 03:36:07.233542+00	2025-10-25 03:36:07.234212+00
3692691b-30a7-42d2-a7fc-912235aa411d	4999dcea-9508-4452-9c5a-bbaf201808ff	2309d442f9fc8eaeef3f5ec80d2d8ce7f9470e3792b1b81b8301dac68383fa9b	2025-11-01 03:37:14.699635+00	2025-10-25 03:37:14.694967+00
143e809f-e4ce-4c3a-abb5-a0e6c84b9aec	4999dcea-9508-4452-9c5a-bbaf201808ff	66b152f977f8c05a56d48afe6ec083677811397987d2a40a899bf3d4234ea111	2025-11-01 03:39:51.945478+00	2025-10-25 03:39:51.942771+00
41aa91fd-8f10-4384-8abd-0c9201061040	4999dcea-9508-4452-9c5a-bbaf201808ff	7cec45b25d20258be0db3203a945e378731bba4d2d99f836e425b931dfedffa0	2025-11-01 03:40:26.766971+00	2025-10-25 03:40:26.721185+00
403a7fb1-5694-4a66-9c5a-fb7de248f427	4999dcea-9508-4452-9c5a-bbaf201808ff	2604ee7ff7ab86727bec21f4fcaab32aee54fe0e16c6e149082d46fb30da9699	2025-11-01 03:46:12.758788+00	2025-10-25 03:46:12.759412+00
1d7b9e88-0653-42dd-9085-2bffab83462e	4999dcea-9508-4452-9c5a-bbaf201808ff	0bb9bb94546ae3bc509b24e4b16f13238743bbba6398614e1e4a7e364d05bda1	2025-11-01 04:07:54.315305+00	2025-10-25 04:07:54.315346+00
f43d8cf3-3a4d-4d0c-8890-7b19a55bd6ed	4999dcea-9508-4452-9c5a-bbaf201808ff	ce836c9067a1aba589393536da52162671f1a4a11e1604dd1d676e7c3d45bd67	2025-11-01 04:13:39.638964+00	2025-10-25 04:13:39.638655+00
0d17cba3-b500-4a6d-88c4-cd1e235a52b0	4999dcea-9508-4452-9c5a-bbaf201808ff	6f698077c6f0e250dec3fabfedb90202b3a26eb0aa9b75523398186cf01d811d	2025-11-01 04:21:52.990573+00	2025-10-25 04:21:52.348719+00
4de0824c-4727-4cd6-a604-ded341ae74d6	4999dcea-9508-4452-9c5a-bbaf201808ff	2ed5591d7d3f0f8c49ffd47a25749aab57c5ce82a90c43df47bd359d879f7c33	2025-11-01 04:27:56.410611+00	2025-10-25 04:27:56.411889+00
85cc1170-f180-4483-8623-caf9b773c2d2	4999dcea-9508-4452-9c5a-bbaf201808ff	7089be6936f61b09e5c80c8cfca06227ab1436b5e28a20d1841ebe3bc43eee39	2025-11-01 04:57:33.450451+00	2025-10-25 04:57:33.452231+00
2851e9e5-7e3b-4eeb-91af-52caa237923e	4999dcea-9508-4452-9c5a-bbaf201808ff	1c89873f4da5ca019365bcaffcc7fd95735ee6c630defd9391d7d7f9894f1b35	2025-11-01 04:58:32.392527+00	2025-10-25 04:58:32.392484+00
6cebf078-4416-4ba1-9ebe-a4f959608af8	4999dcea-9508-4452-9c5a-bbaf201808ff	d1458004219f831aa8b96e3a4a0b78c9ad098fc52a57365f230af088a5e10ba3	2025-11-01 05:04:41.449215+00	2025-10-25 05:04:41.451342+00
401b7bd3-bedc-4f1b-954e-02faa6c1d926	bfd6036f-b1c6-420c-b782-c1e7e918a27b	426b3e76a1ce5179bbc2ecc5fc2035530383f606730cd00c6932d70bd085b9bc	2025-11-01 19:55:13.994144+00	2025-10-25 19:55:13.997388+00
5558f128-7c25-43b2-b8a8-aafbbe5df132	bfd6036f-b1c6-420c-b782-c1e7e918a27b	61be51d50cfb225a177904dcaeaeb59960be37862e4c45b67a6efefebeba8e44	2025-11-01 19:58:39.710206+00	2025-10-25 19:58:39.710702+00
09e9bf81-aaf5-4a79-a865-f13bf028bee7	4999dcea-9508-4452-9c5a-bbaf201808ff	5996215b2155bdc231c60454b1fc2bead88756b11542ac0cd2792f0906b9ad05	2025-11-01 19:58:58.194354+00	2025-10-25 19:58:58.194595+00
236e725b-8588-4449-bdda-2417e58ca037	bfd6036f-b1c6-420c-b782-c1e7e918a27b	0ef4c97f6302665da652d60e96be3a8c23f842c10b0823ad09492bae176d3ed1	2025-11-01 20:07:00.007737+00	2025-10-25 20:07:00.007233+00
eee4bc9b-1528-4f67-a119-c916c4067214	bfd6036f-b1c6-420c-b782-c1e7e918a27b	dddf996ea52f252d3998e1c0b6eb325b553d82b8606330d605c591372dae3a98	2025-11-01 20:13:10.39713+00	2025-10-25 20:13:10.397912+00
40926d30-3d48-4540-b236-7e6983b683c9	bfd6036f-b1c6-420c-b782-c1e7e918a27b	2bfe9b8965b1c42bd6d0e7e548b002a283e572d7a88e2e2d3b85ef8070cce746	2025-11-01 20:16:07.855236+00	2025-10-25 20:16:07.855138+00
141bc70e-59dd-41cb-a411-6298e80701e0	bfd6036f-b1c6-420c-b782-c1e7e918a27b	466694f19e5d764019e72d0eba1756b5c7e8db16a7f63696e2ea56669283839b	2025-11-01 20:16:41.86369+00	2025-10-25 20:16:41.863546+00
08114ec6-c1fc-498c-b5a5-3f60b70f9ee7	bfd6036f-b1c6-420c-b782-c1e7e918a27b	de310e4a96280632d3cfe7e0bef00e2ce47eeb9d7e172f639c20164968650199	2025-11-01 20:17:01.775507+00	2025-10-25 20:17:01.775219+00
3e8ab57d-0197-4b9a-8123-d3e6fc9148cf	bfd6036f-b1c6-420c-b782-c1e7e918a27b	e033bd1a02bd072f576e65705fdc025ff9cec8d9d5ac64c2f5f58782a9826a4d	2025-11-01 20:17:21.499794+00	2025-10-25 20:17:21.499754+00
6e0dd1e5-89ff-4037-b924-4a79d3ed4f96	bfd6036f-b1c6-420c-b782-c1e7e918a27b	df607adbbe460aa5078e4367ce36216df9aa8f775f72d4467bb4510ce36f13c0	2025-11-01 20:18:08.827075+00	2025-10-25 20:18:08.826562+00
c340bb12-d694-422d-953c-50e3e9148d79	bfd6036f-b1c6-420c-b782-c1e7e918a27b	4cffa77fa3070c1eddeac7ad9f15c43260bc2a4f6f02ba3a14c1d455081ee3e5	2025-11-01 20:22:30.348431+00	2025-10-25 20:22:30.348677+00
9f917152-4b5a-422f-a810-9c4a6cec2304	bfd6036f-b1c6-420c-b782-c1e7e918a27b	9d1a131aed1e236b34a1b9b54e736cfd8ded908e73c6b8a0146385a85f282821	2025-11-01 20:32:04.135913+00	2025-10-25 20:32:04.138552+00
b8b103e2-f4cc-4bff-b3f6-beab9fdbc673	bfd6036f-b1c6-420c-b782-c1e7e918a27b	e3e8703e12edadbaade14265e7e602d591dec041254d0d1abbda410f77a192af	2025-11-01 20:32:14.667481+00	2025-10-25 20:32:14.6681+00
e81554b8-6f28-4449-ae96-78d794bcf70b	bfd6036f-b1c6-420c-b782-c1e7e918a27b	e978991beb9219438e913c1730411df556506f40eee5e0dededa7e1dee9e6c1d	2025-11-01 20:34:45.820434+00	2025-10-25 20:34:45.822178+00
de8f1085-e8a7-4da9-9107-d4848a1f62dc	bfd6036f-b1c6-420c-b782-c1e7e918a27b	e32835a42d3bdaf9507a150e1bd665ac08b3563826a1fce95e328d5facf99e33	2025-11-01 20:35:09.072889+00	2025-10-25 20:35:09.073161+00
259d6592-c0be-497e-af5c-b5abd41ac121	bfd6036f-b1c6-420c-b782-c1e7e918a27b	2872a7apg_dump: processing data for table "public.registration_tokens"
pg_dump: dumping contents of table "public.registration_tokens"
pg_dump: processing data for table "public.requests"
85c16dae9a86f579e0b09b26917c82b57a56707c5ec85fdedc9ad69a7	2025-11-01 20:35:09.199814+00	2025-10-25 20:35:09.19979+00
a5fcc2bd-fd48-4fd5-9a6f-dfb7c895bf26	4999dcea-9508-4452-9c5a-bbaf201808ff	c9b4f26262fcd249752e4cd890215cd3324a64ce7d078f69839e92757830c0ed	2025-11-02 17:09:08.778377+00	2025-10-26 17:09:08.781219+00
\.


--
-- TOC entry 3703 (class 0 OID 16787)
-- Dependencies: 230
-- Data for Name: registration_tokens; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.registration_tokens (id, token, email, org_id, role, expires_at, used, created_at, company_profile) FROM stdin;
0bf4d9fb-4fbb-4ea3-9366-522d61c4986d	59f04d08ecc64fcee1e292e3249de47c76aab077a10518d812a265ee93e02460	diego@otrodev.cl	17a55804-4644-46ea-a28f-5c4efc8fbf67	COMPANY	2025-10-04 02:37:48+00	t	2025-10-03 02:37:48+00	\N
14e3891b-2331-4032-ae77-12cc2aff04b9	bc3d95d0e034b31f2aa81e17a567270c22927cdd95ba653061ff0e631bbcf452	djaramontenegro@gmail.com	17a55804-4644-46ea-a28f-5c4efc8fbf67	COMPANY	2025-10-04 04:39:11+00	t	2025-10-03 04:39:11+00	COMPANY_USER
56d51055-b519-4a55-9b32-9fc09426cbd1	c9dc9b7f42b5878d35502df01e5936d9243d13409fb202036b8beaa743e0568e	diego@otrodev.cl	\N	USER	2025-10-12 01:30:41+00	t	2025-10-11 01:30:41+00	\N
8db6516b-d189-41d9-800c-1cef6e848c39	49d411b9-b001-4240-a8a2-515afc4b69d4	final@test.com	\N	DRIVER	2025-10-26 03:37:20+00	f	2025-10-25 03:37:20+00	\N
3be8a9a7-30f5-4e54-9e94-71e5fce0e87f	33aa4ec3-7e3d-491d-b01b-f5c5053a7da6	token@test.com	\N	DRIVER	2025-10-26 03:39:57+00	f	2025-10-25 03:39:57+00	\N
1867db91-496d-4d07-a3ff-386d7f118a65	b8c6ea0f-a723-4de6-b0db-468f5f34e11c	token@debug.com	\N	DRIVER	2025-10-26 03:40:33+00	f	2025-10-25 03:40:33+00	\N
8f85f89b-5c38-4fc3-8dc0-aab36788988f	a3539cb2-37de-4f7f-996b-00ea9c246d6c	test@complete.com	\N	DRIVER	2025-10-26 03:46:20+00	t	2025-10-25 03:46:20+00	\N
8d56bc65-6957-4916-947a-d55085d0e06a	d0482ff2-7a59-4db8-8cdc-bb4e10a83365	djaramontenegro@gmail.com	\N	DRIVER	2025-10-26 03:43:13+00	t	2025-10-25 03:43:13+00	\N
3e4e5688-5e31-432e-9ffe-da23c5f5f7b1	22806931-8e2f-473f-9189-b22cbce44d05	djaramontenegro@gmail.com	\N	DRIVER	2025-10-26 04:08:06+00	f	2025-10-25 04:08:06+00	\N
61daaecc-d906-4440-a920-b350494718e2	8f55fdd3-123e-4ce8-8bf2-51f1b98b0123	djaramontenegro@gmail.com	\N	DRIVER	2025-10-26 04:15:37+00	f	2025-10-25 04:15:37+00	\N
b7de5919-455b-4257-997c-593d04b49188	328e6ab9-692b-427b-aea6-623941d249d8	djaramontenegro@gmail.com	\N	DRIVER	2025-10-26 04:23:06+00	f	2025-10-25 04:23:06+00	\N
264a9590-61d6-4300-812c-29b179eff9aa	848e23fb-263b-4e48-b53f-8ca0fb825829	diego@otrodev.cl	\N	DRIVER	2025-10-26 04:53:20+00	t	2025-10-25 04:53:20+00	\N
08903b19-26e5-4e1b-ad47-ed783db3020b	bb170c20-d444-499a-bfb4-182fe07768b2	diego@otrodev.cl	\N	DRIVER	2025-10-26 05:05:24+00	t	2025-10-25 05:05:24+00	\N
3dedeb2d-c465-4b25-b812-5c894651d92d	27bd14d2a2f14066c26d9d9c14d0437f97f974682919e92fc63506f39d2201c1	djaramontenegro@gmail.com	\N	USER	2025-10-26 13:00:58+00	t	2025-10-25 13:00:58+00	\N
747f5926-58e5-47d3-ad1b-800f974cf353	85a811a1b1f30cf19cfe12508dfe031c6666e1be49990056795bbeceb0a08d0d	djaramontenegro@gmail.com	\N	USER	2025-10-26 13:09:21+00	t	2025-10-25 13:09:21+00	\N
8cd7e50d-7df9-474d-bede-a62243e49993	017aee037cccfd3335eebe34fcd5571e853648db556e1428f2668550c8291094	test1761423441@example.com	\N	USER	2025-10-26 20:17:21+00	f	2025-10-25 20:17:21+00	\N
de8a29eb-d20d-4352-b325-61618f39360d	8b732c63365084df78702a27d6a809e790ed2e451e1179566073cdd99b41d7cc	test.pago1761423488@example.com	\N	USER	2025-10-26 20:18:08+00	f	2025-10-25 20:18:08+00	\N
b4efa648-df23-464b-87c4-eff3eafd6374	eaac2eaded16b1417586ef7c192778b3b0ca80e47e1caf5ff4bcf35d9c062dde	test.pago1761423750@example.com	\N	USER	2025-10-26 20:22:30+00	f	2025-10-25 20:22:30+00	\N
574e27e4-257a-4601-a535-9a1d3a72754c	3d12277e299262724d99f40bfd1a21cb06d9541076c0e1e6dd4f2a3ed3879aae	test.reserva.pendiente1761424318@example.com	\N	USER	2025-10-26 20:31:58+00	f	2025-10-25 20:31:58+00	\N
\.


--
-- TOC entry 3698 (class 0 OID 16644)
-- Dependencies: 225
-- Data for Name: requests; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.requests (id, hopg_dump: dumping contents of table "public.requests"
pg_dump: processing data for table "public.reservation_timeline"
pg_dump: dumping contents of table "public.reservation_timeline"
pg_dump: processing data for table "public.reservations"
pg_dump: dumping contents of table "public.reservations"
tel_id, company_id, fecha, origin, destination, pax, vehicle_type, language, status, assigned_driver_id, created_at, updated_at) FROM stdin;
\.


--
-- TOC entry 3700 (class 0 OID 16688)
-- Dependencies: 227
-- Data for Name: reservation_timeline; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.reservation_timeline (id, reservation_id, title, description, at, variant, created_at) FROM stdin;
c9309088-40c7-4d71-b399-ac2133734e33	RSV-1001	Reserva creada	La reserva ha sido creada exitosamente	2025-10-01 20:16:41.226742+00	success	2025-10-02 20:16:41.204462+00
f26bfcee-a467-417e-8e10-51aa76766fb4	RSV-1002	Reserva creada	La reserva ha sido creada exitosamente	2025-10-01 20:16:41.228468+00	success	2025-10-02 20:16:41.204462+00
732af19f-674d-4f15-8d05-7a4d6d5f4b65	RSV-1002	Reserva programada	La reserva ha sido programada	2025-10-02 08:16:41.228469+00	info	2025-10-02 20:16:41.204462+00
8be9278d-5fd1-4333-a68f-344da0ea3cba	RSV-1003	Reserva creada	La reserva ha sido creada exitosamente	2025-10-01 20:16:41.231913+00	success	2025-10-02 20:16:41.204462+00
b5bf4775-bca2-4088-8e74-c89f5e94c2e2	RSV-1003	Reserva programada	La reserva ha sido programada	2025-10-02 08:16:41.231914+00	info	2025-10-02 20:16:41.204462+00
247b6c2b-aa85-4b94-8b03-e2c1dda1d17c	RSV-1003	Servicio completado	El servicio ha sido completado exitosamente	2025-10-02 19:16:41.231914+00	success	2025-10-02 20:16:41.204462+00
c3ee4a4c-d588-4514-8ec5-a4932714905f	RSV-24509	Reserva creada	La reserva ha sido creada exitosamente	2025-10-25 20:35:09.213818+00	success	2025-10-25 20:35:09.214069+00
\.


--
-- TOC entry 3699 (class 0 OID 16672)
-- Dependencies: 226
-- Data for Name: reservations; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.reservations (id, user_id, org_id, pickup, destination, datetime, passengers, status, amount, notes, created_at, updated_at, assigned_driver_id, distance_km, arrived_on_time) FROM stdin;
RSV-1003	\N	\N	Oficina Central	Mina El Teniente	2025-10-01 20:16:41.22488+00	8	COMPLETADA	250000.00	\N	2025-10-02 20:16:41.204462+00	2025-10-02 20:16:41.204462+00	\N	\N	t
RSV-1001	4323f59c-d1f4-4cff-a693-fb210a9d68a0	\N	Hotel Miramar, Valparaíso	Aeropuerto SCL, Santiago	2025-10-09 20:16:41.22488+00	2	COMPLETADA	120000.00	\N	2025-10-02 20:16:41.204462+00	2025-10-11 01:29:50.88271+00	\N	\N	t
RSV-1002	4323f59c-d1f4-4cff-a693-fb210a9d68a0	\N	Hotel Andes, Santiago	Mall Costanera Center	2025-10-03 20:16:41.22488+00	4	COMPLETADA	80000.00	\N	2025-10-02 20:16:41.204462+00	2025-10-11 01:30:00.16094+00	\N	\N	t
RES001	4999dcea-9508-4452-9c5a-bbaf201808ff	\N	Aeropuerto SCL	Hotel Plaza	2024-01-15 10:00:00+00	2	COMPLETADA	\N	\N	2025-10-24 04:48:27.822058+00	2025-10-24 04:48:27.822058+00	DRV003	15.50	t
RES002	4999dcea-9508-4452-9c5a-bbaf201808ff	\N	Hotel Plaza	Centro Comercial	2024-01-16 14:30:00+00	4	COMPLETADA	\N	\N	2025-10-24 04:48:27.822058+00	2025-10-24 04:48:27.822058+00	DRV003	8.20	t
RES003	4999dcea-9508-4452-9c5a-bbaf201808ff	\N	Centro	Aeropuerto SCL	2024-01-17 16:00:00+00	1	COMPLETADA	\N	\N	2025-10-24 04:48:27.822058+00	2025-10-24 04:48:27.822058+00	DRV003	12.80	f
RES004	4999dcea-9508-4452-9c5a-bbaf201808ff	\N	Hotel	Mall	2024-01-18 09:00:00+00	3	CANCELADA	\N	\N	2025-10-24 04:48:27.822058+00	2025-10-24 04:48:27.822058+00	DRV003	5.50	t
RES005	4999dcea-9508-4452-9c5a-bbaf201808ff	\N	Aeropuerto SCL	Hotel	2024-01-19 11:30:00+00	2	COMPLETADA	\N	\N	2025-10-24 04:48:27.822058+00	2025-10-24 04:48:27.822058+00	DRV003	18.30	t
RES006	4999dcea-9508-4452-9c5a-bbaf201808ff	\N	Centro	Aeropuerto SCL	2024-01-20 08:00:00+00	1	COMPLETADA	\N	\N	2025-10-24 05:04:12.172471+00	2025-10-24 05:04:12.172471+00	DRV001	25.00	t
RES007	4999dcea-9508-4452-9c5a-bbaf201808ff	\N	Hotel	Mall	2024-01-21 15:30:00+00	2	COMPLETADA	\N	\N	2025-10-24 05:04:12.172471+00	2025-10-24 05:04:12.172471+00	DRV001	12.50	t
RES008	4999dcea-9508-4452-9c5a-bbaf201808ff	\N	Aeropuerto SCL	Hotel	2024-01-22 10:15:00+00	3	CANCELADA	\N	\N	2025-10-24 05:04:12.172471+00	2025-10-24 05:04:12.172471+00	DRV001	18.00	t
RES009	4999dcea-9508-4452-9c5a-bbaf201808ff	\N	Centro	Aeropuerto SCL	2024-01-23 09:00:00+00	2	pg_dump: processing data for table "public.schema_migrations"
pg_dump: dumping contents of table "public.schema_migrations"
pg_dump: processing data for table "public.users"
pg_dump: dumping contents of table "public.users"
COMPLETADA	\N	\N	2025-10-24 05:04:12.172471+00	2025-10-24 05:04:12.172471+00	DRV002	22.00	f
RES010	4999dcea-9508-4452-9c5a-bbaf201808ff	\N	Hotel	Centro	2024-01-24 14:00:00+00	1	COMPLETADA	\N	\N	2025-10-24 05:04:12.172471+00	2025-10-24 05:04:12.172471+00	DRV002	8.50	t
RSV-60893	\N	\N	prueba 1 11, Santiago, Región Metropolitana	prueba 2 222, La Florida, Región Metropolitana	2025-10-17 03:07:00+00	1	COMPLETADA	80000.00	\N	2025-10-03 03:08:13.615634+00	2025-10-24 05:24:04.683948+00	\N	\N	t
RSV-49188	\N	\N	C. la Cueca 1415, Quillota, Valparaíso, Chile	Av Costanera 4247, Renca, Región Metropolitana, Chile	2025-10-25 02:19:00+00	1	COMPLETADA	80000.00	\N	2025-10-11 02:19:48.351063+00	2025-10-25 01:23:34.307263+00	\N	\N	t
RSV-51289	\N	\N	Zañartu 2132, 7780299 Ñuñoa, Región Metropolitana, Chile	Aeropuerto SCL, Pudahuel, Región Metropolitana, Chile	2025-10-31 02:53:00+00	1	COMPLETADA	80000.00	\N	2025-10-11 02:54:49.877111+00	2025-10-25 01:23:34.307263+00	\N	\N	t
RSV-97960	\N	\N	C. la Cueca 1415, Quillota, Valparaíso, Chile	Av Costanera 4247, Renca, Región Metropolitana, Chile	2025-11-13 13:12:00+00	1	COMPLETADA	25000.00	\N	2025-10-25 13:12:40.76775+00	2025-10-25 19:59:10.210241+00	DRV-1761368724331	15.50	t
RSV-23195	bfd6036f-b1c6-420c-b782-c1e7e918a27b	\N	Test Origin	Test Destination	2025-10-26 09:00:00+00	1	ACTIVA	120000.00	\N	2025-10-25 20:13:16.014498+00	2025-10-25 20:13:16.014498+00	\N	\N	t
RSV-24509	bfd6036f-b1c6-420c-b782-c1e7e918a27b	\N	Santiago Centro, Chile	Aeropuerto Arturo Merino Benítez, Santiago, Chile	2025-10-26 20:35:08+00	3	ACTIVA	155000.00	Reserva desde prueba automática - Usuario pendiente de registro	2025-10-25 20:35:09.211727+00	2025-10-25 20:35:09.211727+00	\N	18.50	t
\.


--
-- TOC entry 3688 (class 0 OID 16389)
-- Dependencies: 215
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.schema_migrations (version, dirty) FROM stdin;
10	f
\.


--
-- TOC entry 3689 (class 0 OID 16525)
-- Dependencies: 216
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.users (id, name, email, password_hash, role, status, org_id, created_at, updated_at, company_profile) FROM stdin;
4999dcea-9508-4452-9c5a-bbaf201808ff	Admin Sistema	admin@turivo.com	$2a$10$9pjpOYuBh0O/loXHGwBOg.l6hrmZHtsae/2UqhD23ff4O5nTwgymC	ADMIN	ACTIVE	\N	2025-10-02 20:16:41.204462+00	2025-10-02 20:16:41.204462+00	\N
4323f59c-d1f4-4cff-a693-fb210a9d68a0	Cliente Demo	cliente@demo.com	$2a$10$9pjpOYuBh0O/loXHGwBOg.l6hrmZHtsae/2UqhD23ff4O5nTwgymC	USER	ACTIVE	\N	2025-10-02 20:16:41.204462+00	2025-10-02 20:16:41.204462+00	\N
a39fe240-a9f9-4947-8055-2aabdb14c29f	Juan Pérez	juan@turivo.com	$2a$10$9pjpOYuBh0O/loXHGwBOg.l6hrmZHtsae/2UqhD23ff4O5nTwgymC	COMPANY	ACTIVE	17a55804-4644-46ea-a28f-5c4efc8fbf67	2025-10-02 20:16:41.204462+00	2025-10-03 03:55:11.032901+00	COMPANY_ADMIN
c36f9105-1bec-4334-bc0a-8c41b07f0285	María González	maria@turismoandes.cl	$2a$10$9pjpOYuBh0O/loXHGwBOg.l6hrmZHtsae/2UqhD23ff4O5nTwgymC	COMPANY	ACTIVE	dd28931b-c0d1-4c8b-a1b5-85ada9390657	2025-10-02 20:16:41.204462+00	2025-10-03 03:55:11.032901+00	COMPANY_ADMIN
e2b48993-e03f-4fc4-a134-9fac83d240ac	Test Driver	test@driver.com	temp_hash	DRIVER	ACTIVE	\N	2025-10-25 03:34:42.594163+00	2025-10-25 03:34:42.594163+00	\N
b8ecbbed-4b60-47dd-ba15-fd7bd6ef0ce9	Nuevo Conductor	nuevo@conductor.com	temp_hash	DRIVER	ACTIVE	\N	2025-10-25 03:35:03.004314+00	2025-10-25 03:35:03.004314+00	\N
2516493d-9919-4174-afe1-3a479b624bc8	Test Token	test@token.com	temp_hash	DRIVER	ACTIVE	\N	2025-10-25 03:36:31.413637+00	2025-10-25 03:36:31.413637+00	\N
0114d672-5c03-4780-b19a-a278fb7f0238	Final Test	final@test.com	temp_hash	DRIVER	ACTIVE	\N	2025-10-25 03:37:20.880708+00	2025-10-25 03:37:20.880708+00	\N
e155ddb8-d7c6-447b-80d0-66bc43914198	Token Test	token@test.com	temp_hash	DRIVER	ACTIVE	\N	2025-10-25 03:39:57.329038+00	2025-10-25 03:39:57.329038+00	\N
dfc0ded4-5493-4045-ac4f-794fa5a08e53	Token Debug	token@debug.com	temp_hash	DRIVER	ACTIVE	\N	2025-10-25 03:40:33.685947+00	2025-10-25 03:40:33.685947+00	\N
2c72496b-bdc5-4dff-86e6-3b75e912905pg_dump: processing data for table "public.vehicle_photos"
pg_dump: dumping contents of table "public.vehicle_photos"
pg_dump: processing data for table "public.vehicles"
pg_dump: dumping contents of table "public.vehicles"
pg_dump: executing SEQUENCE SET migrations_id_seq
pg_dump: creating CONSTRAINT "public.companies companies_pkey"
pg_dump: creating CONSTRAINT "public.companies companies_rut_key"
pg_dump: creating CONSTRAINT "public.driver_availability driver_availability_pkey"
pg_dump: creating CONSTRAINT "public.driver_background_checks driver_background_checks_pkey"
pg_dump: creating CONSTRAINT "public.driver_feedback driver_feedback_pkey"
pg_dump: creating CONSTRAINT "public.driver_licenses driver_licenses_pkey"
pg_dump: creating CONSTRAINT "public.drivers drivers_pkey"
pg_dump: creating CONSTRAINT "public.drivers drivers_rut_or_dni_key"
f	Test Complete	test@complete.com	$2a$10$kVbYqnnJQr4L4BSAT7RaKeQS1PBH7CcLd4IiEySAsYF.LN4FZKuMy	DRIVER	ACTIVE	\N	2025-10-25 03:46:20.550087+00	2025-10-25 03:47:43.536132+00	\N
0bb4896a-8d3b-4273-9c68-aa7fc572279c	Juan Pérez	juan.perez@test.com	temp_hash	DRIVER	ACTIVE	\N	2025-10-25 04:57:43.502503+00	2025-10-25 04:57:43.502503+00	\N
bfd6036f-b1c6-420c-b782-c1e7e918a27b	Sistema Turivo	system@turivo.com	$2y$10$6PeOWiC0L.SbItNa3E5YVuZmaOUya76CBB3wa7yYFeg9tAOmEHlPm	ADMIN	ACTIVE	\N	2025-10-25 18:53:16.669996+00	2025-10-25 19:55:10.505571+00	\N
328d2373-40e3-40ba-b42e-f629d8efb445	Usuario Test Registro	usuario.registro@test.com	$2a$10$ALwTC1f23uARtlPgS8O26.qRUxkm9.WGXrrUhnW9fTY60KV9baiN.	USER	ACTIVE	\N	2025-10-25 20:35:09.130544+00	2025-10-25 20:35:09.130544+00	\N
\.


--
-- TOC entry 3696 (class 0 OID 16614)
-- Dependencies: 223
-- Data for Name: vehicle_photos; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.vehicle_photos (id, vehicle_id, url, created_at) FROM stdin;
\.


--
-- TOC entry 3695 (class 0 OID 16599)
-- Dependencies: 222
-- Data for Name: vehicles; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.vehicles (id, driver_id, type, brand, model, year, plate, vin, color, insurance_policy, insurance_expires_at, inspection_expires_at, created_at, updated_at, status, capacity) FROM stdin;
550e8400-e29b-41d4-a716-446655440002	\N	VAN	Ford	Transit	2022	DEF-456	VIN987654321	Azul	\N	\N	\N	2025-10-24 03:40:37.208684+00	2025-10-24 03:40:37.208684+00	AVAILABLE	12
550e8400-e29b-41d4-a716-446655440003	\N	SEDAN	Toyota	Camry	2023	GHI-789	VIN456789123	Negro	\N	\N	\N	2025-10-24 03:40:37.208684+00	2025-10-24 03:40:37.208684+00	AVAILABLE	4
550e8400-e29b-41d4-a716-446655440004	\N	SUV	Nissan	Pathfinder	2023	JKL-012	VIN111222333	Gris	\N	\N	\N	2025-10-24 03:40:37.208684+00	2025-10-24 03:40:37.208684+00	AVAILABLE	7
550e8400-e29b-41d4-a716-446655440001	DRV-1761368724331	BUS	Mercedes	Sprinter	2023	ABC-123	VIN123456789	Blanco	\N	\N	\N	2025-10-24 03:40:37.208684+00	2025-10-26 18:31:50.344484+00	ASSIGNED	20
\.


--
-- TOC entry 3716 (class 0 OID 0)
-- Dependencies: 232
-- Name: migrations_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.migrations_id_seq', 1, false);


--
-- TOC entry 3449 (class 2606 OID 16549)
-- Name: companies companies_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.companies
    ADD CONSTRAINT companies_pkey PRIMARY KEY (id);


--
-- TOC entry 3451 (class 2606 OID 16551)
-- Name: companies companies_rut_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.companies
    ADD CONSTRAINT companies_rut_key UNIQUE (rut);


--
-- TOC entry 3478 (class 2606 OID 16638)
-- Name: driver_availability driver_availability_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.driver_availability
    ADD CONSTRAINT driver_availability_pkey PRIMARY KEY (driver_id);


--
-- TOC entry 3467 (class 2606 OID 16593)
-- Name: driver_background_checks driver_background_checks_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.driver_background_checks
    ADD CONSTRAINT driver_background_checks_pkey PRIMARY KEY (driver_id);


--
-- TOC entry 3511 (class 2606 OID 24667)
-- Name: driver_feedback driver_feedback_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.driver_feedback
    ADD CONSTRAINT driver_feedback_pkey PRIMARY KEY (id);


--
-- TOC entry 3465 (class 2606 OID 16580)
-- Name: driver_licenses driver_licenses_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.driver_licenses
    ADD CONSTRAINT driver_licenses_pkey PRIMARY KEY (driver_id);


--
-- TOC entry 3456 (class 2606 OID 16571)
-- Name: drivers drivers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.drivers
    ADD CONSTRAINT drivers_pkey PRIMARY KEY (id);


--
-- TOC entry 3458 (class 2606 OID 16573)
-- Name: drivers drivers_rut_or_dni_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.drivers
    ADD CONSTRAINT drivers_rpg_dump: creating CONSTRAINT "public.hotels hotels_pkey"
pg_dump: creating CONSTRAINT "public.migrations migrations_pkey"
pg_dump: creating CONSTRAINT "public.payments payments_pkey"
pg_dump: creating CONSTRAINT "public.refresh_tokens refresh_tokens_pkey"
pg_dump: creating CONSTRAINT "public.refresh_tokens refresh_tokens_token_key"
pg_dump: creating CONSTRAINT "public.registration_tokens registration_tokens_pkey"
pg_dump: creating CONSTRAINT "public.registration_tokens registration_tokens_token_key"
pg_dump: creating CONSTRAINT "public.requests requests_pkey"
pg_dump: creating CONSTRAINT "public.reservation_timeline reservation_timeline_pkey"
pg_dump: creating CONSTRAINT "public.reservations reservations_pkey"
pg_dump: creating CONSTRAINT "public.schema_migrations schema_migrations_pkey"
pg_dump: creating CONSTRAINT "public.vehicles unique_driver_vehicle"
pg_dump: creating CONSTRAINT "public.users users_email_key"
pg_dump: creating CONSTRAINT "public.users users_pkey"
pg_dump: creating CONSTRAINT "public.vehicle_photos vehicle_photos_pkey"
pg_dump: creating CONSTRAINT "public.vehicles vehicles_pkey"
pg_dump: creating INDEX "public.idx_companies_rut"
pg_dump: creating INDEX "public.idx_driver_feedback_created_at"
pg_dump: creating INDEX "public.idx_driver_feedback_driver_id"
ut_or_dni_key UNIQUE (rut_or_dni);


--
-- TOC entry 3454 (class 2606 OID 16561)
-- Name: hotels hotels_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.hotels
    ADD CONSTRAINT hotels_pkey PRIMARY KEY (id);


--
-- TOC entry 3516 (class 2606 OID 24733)
-- Name: migrations migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.migrations
    ADD CONSTRAINT migrations_pkey PRIMARY KEY (id);


--
-- TOC entry 3494 (class 2606 OID 16713)
-- Name: payments payments_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_pkey PRIMARY KEY (id);


--
-- TOC entry 3499 (class 2606 OID 16772)
-- Name: refresh_tokens refresh_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_pkey PRIMARY KEY (id);


--
-- TOC entry 3501 (class 2606 OID 16774)
-- Name: refresh_tokens refresh_tokens_token_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_token_key UNIQUE (token);


--
-- TOC entry 3507 (class 2606 OID 16796)
-- Name: registration_tokens registration_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.registration_tokens
    ADD CONSTRAINT registration_tokens_pkey PRIMARY KEY (id);


--
-- TOC entry 3509 (class 2606 OID 16798)
-- Name: registration_tokens registration_tokens_token_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.registration_tokens
    ADD CONSTRAINT registration_tokens_token_key UNIQUE (token);


--
-- TOC entry 3483 (class 2606 OID 16656)
-- Name: requests requests_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.requests
    ADD CONSTRAINT requests_pkey PRIMARY KEY (id);


--
-- TOC entry 3490 (class 2606 OID 16697)
-- Name: reservation_timeline reservation_timeline_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.reservation_timeline
    ADD CONSTRAINT reservation_timeline_pkey PRIMARY KEY (id);


--
-- TOC entry 3488 (class 2606 OID 16682)
-- Name: reservations reservations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.reservations
    ADD CONSTRAINT reservations_pkey PRIMARY KEY (id);


--
-- TOC entry 3439 (class 2606 OID 16393)
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- TOC entry 3472 (class 2606 OID 16834)
-- Name: vehicles unique_driver_vehicle; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.vehicles
    ADD CONSTRAINT unique_driver_vehicle UNIQUE (driver_id);


--
-- TOC entry 3445 (class 2606 OID 16538)
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- TOC entry 3447 (class 2606 OID 16536)
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- TOC entry 3476 (class 2606 OID 16622)
-- Name: vehicle_photos vehicle_photos_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.vehicle_photos
    ADD CONSTRAINT vehicle_photos_pkey PRIMARY KEY (id);


--
-- TOC entry 3474 (class 2606 OID 16608)
-- Name: vehicles vehicles_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.vehicles
    ADD CONSTRAINT vehicles_pkey PRIMARY KEY (id);


--
-- TOC entry 3452 (class 1259 OID 16742)
-- Name: idx_companies_rut; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_companies_rut ON public.companies USING btree (rut);


--
-- TOC entry 3512 (class 1259 OID 24680)
-- Name: idx_driver_feedback_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_driver_feedback_created_at ON public.driver_feedback USING btree (created_at);


--
--pg_dump: creating INDEX "public.idx_driver_feedback_reservation_id"
pg_dump: creating INDEX "public.idx_drivers_company_id"
pg_dump: creating INDEX "public.idx_drivers_status"
pg_dump: creating INDEX "public.idx_drivers_unique_vehicle"
pg_dump: creating INDEX "public.idx_drivers_user_id"
pg_dump: creating INDEX "public.idx_drivers_vehicle_id"
pg_dump: creating INDEX "public.idx_payments_reservation_id"
pg_dump: creating INDEX "public.idx_payments_status"
pg_dump: creating INDEX "public.idx_refresh_tokens_expires_at"
pg_dump: creating INDEX "public.idx_refresh_tokens_token"
pg_dump: creating INDEX "public.idx_refresh_tokens_user_id"
pg_dump: creating INDEX "public.idx_registration_tokens_email"
pg_dump: creating INDEX "public.idx_registration_tokens_expires_at"
pg_dump: creating INDEX "public.idx_registration_tokens_token"
pg_dump: creating INDEX "public.idx_registration_tokens_used"
pg_dump: creating INDEX "public.idx_requests_assigned_driver"
pg_dump: creating INDEX "public.idx_requests_fecha"
pg_dump: creating INDEX "public.idx_requests_status"
pg_dump: creating INDEX "public.idx_reservations_assigned_driver"
 TOC entry 3513 (class 1259 OID 24678)
-- Name: idx_driver_feedback_driver_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_driver_feedback_driver_id ON public.driver_feedback USING btree (driver_id);


--
-- TOC entry 3514 (class 1259 OID 24679)
-- Name: idx_driver_feedback_reservation_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_driver_feedback_reservation_id ON public.driver_feedback USING btree (reservation_id);


--
-- TOC entry 3459 (class 1259 OID 24712)
-- Name: idx_drivers_company_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_drivers_company_id ON public.drivers USING btree (company_id);


--
-- TOC entry 3460 (class 1259 OID 16743)
-- Name: idx_drivers_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_drivers_status ON public.drivers USING btree (status);


--
-- TOC entry 3461 (class 1259 OID 24719)
-- Name: idx_drivers_unique_vehicle; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_drivers_unique_vehicle ON public.drivers USING btree (vehicle_id) WHERE (vehicle_id IS NOT NULL);


--
-- TOC entry 3462 (class 1259 OID 24706)
-- Name: idx_drivers_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_drivers_user_id ON public.drivers USING btree (user_id);


--
-- TOC entry 3463 (class 1259 OID 24718)
-- Name: idx_drivers_vehicle_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_drivers_vehicle_id ON public.drivers USING btree (vehicle_id);


--
-- TOC entry 3491 (class 1259 OID 16750)
-- Name: idx_payments_reservation_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_payments_reservation_id ON public.payments USING btree (reservation_id);


--
-- TOC entry 3492 (class 1259 OID 16749)
-- Name: idx_payments_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_payments_status ON public.payments USING btree (status);


--
-- TOC entry 3495 (class 1259 OID 16782)
-- Name: idx_refresh_tokens_expires_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_refresh_tokens_expires_at ON public.refresh_tokens USING btree (expires_at);


--
-- TOC entry 3496 (class 1259 OID 16781)
-- Name: idx_refresh_tokens_token; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_refresh_tokens_token ON public.refresh_tokens USING btree (token);


--
-- TOC entry 3497 (class 1259 OID 16780)
-- Name: idx_refresh_tokens_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_refresh_tokens_user_id ON public.refresh_tokens USING btree (user_id);


--
-- TOC entry 3502 (class 1259 OID 16805)
-- Name: idx_registration_tokens_email; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_registration_tokens_email ON public.registration_tokens USING btree (email);


--
-- TOC entry 3503 (class 1259 OID 16806)
-- Name: idx_registration_tokens_expires_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_registration_tokens_expires_at ON public.registration_tokens USING btree (expires_at);


--
-- TOC entry 3504 (class 1259 OID 16804)
-- Name: idx_registration_tokens_token; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_registration_tokens_token ON public.registration_tokens USING btree (token);


--
-- TOC entry 3505 (class 1259 OID 16807)
-- Name: idx_registration_tokens_used; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_registration_tokens_used ON public.registration_tokens USING btree (used);


--
-- TOC entry 3479 (class 1259 OID 16746)
-- Name: idx_requests_assigned_driver; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_requests_assigned_driver ON public.requests USING btree (assigned_driver_id);


--
-- TOC entry 3480 (class 1259 OID 16745)
-- Name: idx_requests_fecha; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_requests_fecha ON public.requests USING btree (fecha);


--
-- TOC entry 3481 (class 1259 OID 16744)
-- Name: idx_requests_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_requests_status ON public.requests USING btree (status);


--
-- TOC entry 3484 (class 1259 OID 16817)
-- Name: idx_repg_dump: creating INDEX "public.idx_reservations_datetime"
pg_dump: creating INDEX "public.idx_reservations_status"
pg_dump: creating INDEX "public.idx_users_company_profile"
pg_dump: creating INDEX "public.idx_users_email"
pg_dump: creating INDEX "public.idx_users_org_id"
pg_dump: creating INDEX "public.idx_users_role"
pg_dump: creating INDEX "public.idx_vehicles_driver_id"
pg_dump: creating INDEX "public.idx_vehicles_status"
pg_dump: creating INDEX "public.idx_vehicles_type"
pg_dump: creating TRIGGER "public.vehicles trigger_update_vehicle_status"
pg_dump: creating TRIGGER "public.companies update_companies_updated_at"
pg_dump: creating TRIGGER "public.driver_availability update_driver_availability_updated_at"
pg_dump: creating TRIGGER "public.drivers update_drivers_updated_at"
pg_dump: creating TRIGGER "public.hotels update_hotels_updated_at"
pg_dump: creating TRIGGER "public.requests update_requests_updated_at"
pg_dump: creating TRIGGER "public.reservations update_reservations_updated_at"
pg_dump: creating TRIGGER "public.users update_users_updated_at"
servations_assigned_driver; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_reservations_assigned_driver ON public.reservations USING btree (assigned_driver_id);


--
-- TOC entry 3485 (class 1259 OID 16748)
-- Name: idx_reservations_datetime; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_reservations_datetime ON public.reservations USING btree (datetime);


--
-- TOC entry 3486 (class 1259 OID 16747)
-- Name: idx_reservations_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_reservations_status ON public.reservations USING btree (status);


--
-- TOC entry 3440 (class 1259 OID 16840)
-- Name: idx_users_company_profile; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_users_company_profile ON public.users USING btree (company_profile) WHERE (company_profile IS NOT NULL);


--
-- TOC entry 3441 (class 1259 OID 16739)
-- Name: idx_users_email; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_users_email ON public.users USING btree (email);


--
-- TOC entry 3442 (class 1259 OID 16741)
-- Name: idx_users_org_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_users_org_id ON public.users USING btree (org_id);


--
-- TOC entry 3443 (class 1259 OID 16740)
-- Name: idx_users_role; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_users_role ON public.users USING btree (role);


--
-- TOC entry 3468 (class 1259 OID 16836)
-- Name: idx_vehicles_driver_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_vehicles_driver_id ON public.vehicles USING btree (driver_id);


--
-- TOC entry 3469 (class 1259 OID 16835)
-- Name: idx_vehicles_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_vehicles_status ON public.vehicles USING btree (status);


--
-- TOC entry 3470 (class 1259 OID 16837)
-- Name: idx_vehicles_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_vehicles_type ON public.vehicles USING btree (type);


--
-- TOC entry 3541 (class 2620 OID 16839)
-- Name: vehicles trigger_update_vehicle_status; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trigger_update_vehicle_status BEFORE UPDATE ON public.vehicles FOR EACH ROW EXECUTE FUNCTION public.update_vehicle_status_on_assignment();


--
-- TOC entry 3538 (class 2620 OID 16753)
-- Name: companies update_companies_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_companies_updated_at BEFORE UPDATE ON public.companies FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 3543 (class 2620 OID 16759)
-- Name: driver_availability update_driver_availability_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_driver_availability_updated_at BEFORE UPDATE ON public.driver_availability FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 3540 (class 2620 OID 16755)
-- Name: drivers update_drivers_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_drivers_updated_at BEFORE UPDATE ON public.drivers FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 3539 (class 2620 OID 16754)
-- Name: hotels update_hotels_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_hotels_updated_at BEFORE UPDATE ON public.hotels FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 3544 (class 2620 OID 16757)
-- Name: requests update_requests_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_requests_updated_at BEFORE UPDATE ON public.requests FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 3545 (class 2620 OID 16758)
-- Name: reservations update_reservations_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_reservations_updated_at BEFORE UPDATE ON public.reservations FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 3537 (class 2620 OID 16752)
-- Name: users update_users_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGEpg_dump: creating TRIGGER "public.vehicles update_vehicles_updated_at"
pg_dump: creating FK CONSTRAINT "public.driver_availability driver_availability_driver_id_fkey"
pg_dump: creating FK CONSTRAINT "public.driver_background_checks driver_background_checks_driver_id_fkey"
pg_dump: creating FK CONSTRAINT "public.driver_feedback driver_feedback_driver_id_fkey"
pg_dump: creating FK CONSTRAINT "public.driver_feedback driver_feedback_reservation_id_fkey"
pg_dump: creating FK CONSTRAINT "public.driver_licenses driver_licenses_driver_id_fkey"
pg_dump: creating FK CONSTRAINT "public.drivers drivers_company_id_fkey"
pg_dump: creating FK CONSTRAINT "public.drivers drivers_user_id_fkey"
pg_dump: creating FK CONSTRAINT "public.drivers drivers_vehicle_id_fkey"
pg_dump: creating FK CONSTRAINT "public.users fk_users_org_id"
pg_dump: creating FK CONSTRAINT "public.payments payments_reservation_id_fkey"
pg_dump: creating FK CONSTRAINT "public.refresh_tokens refresh_tokens_user_id_fkey"
pg_dump: creating FK CONSTRAINT "public.registration_tokens registration_tokens_org_id_fkey"
R update_users_updated_at BEFORE UPDATE ON public.users FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 3542 (class 2620 OID 16756)
-- Name: vehicles update_vehicles_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_vehicles_updated_at BEFORE UPDATE ON public.vehicles FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- TOC entry 3525 (class 2606 OID 16639)
-- Name: driver_availability driver_availability_driver_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.driver_availability
    ADD CONSTRAINT driver_availability_driver_id_fkey FOREIGN KEY (driver_id) REFERENCES public.drivers(id) ON DELETE CASCADE;


--
-- TOC entry 3522 (class 2606 OID 16594)
-- Name: driver_background_checks driver_background_checks_driver_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.driver_background_checks
    ADD CONSTRAINT driver_background_checks_driver_id_fkey FOREIGN KEY (driver_id) REFERENCES public.drivers(id) ON DELETE CASCADE;


--
-- TOC entry 3535 (class 2606 OID 24668)
-- Name: driver_feedback driver_feedback_driver_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.driver_feedback
    ADD CONSTRAINT driver_feedback_driver_id_fkey FOREIGN KEY (driver_id) REFERENCES public.drivers(id) ON DELETE CASCADE;


--
-- TOC entry 3536 (class 2606 OID 24673)
-- Name: driver_feedback driver_feedback_reservation_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.driver_feedback
    ADD CONSTRAINT driver_feedback_reservation_id_fkey FOREIGN KEY (reservation_id) REFERENCES public.reservations(id) ON DELETE CASCADE;


--
-- TOC entry 3521 (class 2606 OID 16581)
-- Name: driver_licenses driver_licenses_driver_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.driver_licenses
    ADD CONSTRAINT driver_licenses_driver_id_fkey FOREIGN KEY (driver_id) REFERENCES public.drivers(id) ON DELETE CASCADE;


--
-- TOC entry 3518 (class 2606 OID 24707)
-- Name: drivers drivers_company_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.drivers
    ADD CONSTRAINT drivers_company_id_fkey FOREIGN KEY (company_id) REFERENCES public.companies(id) ON DELETE SET NULL;


--
-- TOC entry 3519 (class 2606 OID 24701)
-- Name: drivers drivers_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.drivers
    ADD CONSTRAINT drivers_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- TOC entry 3520 (class 2606 OID 24713)
-- Name: drivers drivers_vehicle_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.drivers
    ADD CONSTRAINT drivers_vehicle_id_fkey FOREIGN KEY (vehicle_id) REFERENCES public.vehicles(id) ON DELETE SET NULL;


--
-- TOC entry 3517 (class 2606 OID 16734)
-- Name: users fk_users_org_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT fk_users_org_id FOREIGN KEY (org_id) REFERENCES public.companies(id) ON DELETE SET NULL;


--
-- TOC entry 3532 (class 2606 OID 16714)
-- Name: payments payments_reservation_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_reservation_id_fkey FOREIGN KEY (reservation_id) REFERENCES public.reservations(id) ON DELETE CASCADE;


--
-- TOC entry 3533 (class 2606 OID 16775)
-- Name: refresh_tokens refresh_tokens_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- TOC entry 3534 (class 2606 OID 16799)
-- Name: registration_tokens registration_tokens_org_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.registration_tokens
    ADD CONSTRAINT registration_tokens_org_id_fkey FOREIGN KEY (org_id) REFERENCES public.companies(id) ON DELEpg_dump: creating FK CONSTRAINT "public.requests requests_assigned_driver_id_fkey"
pg_dump: creating FK CONSTRAINT "public.requests requests_company_id_fkey"
pg_dump: creating FK CONSTRAINT "public.requests requests_hotel_id_fkey"
pg_dump: creating FK CONSTRAINT "public.reservation_timeline reservation_timeline_reservation_id_fkey"
pg_dump: creating FK CONSTRAINT "public.reservations reservations_assigned_driver_id_fkey"
pg_dump: creating FK CONSTRAINT "public.reservations reservations_user_id_fkey"
pg_dump: creating FK CONSTRAINT "public.vehicle_photos vehicle_photos_vehicle_id_fkey"
pg_dump: creating FK CONSTRAINT "public.vehicles vehicles_driver_id_fkey"
TE SET NULL;


--
-- TOC entry 3526 (class 2606 OID 16667)
-- Name: requests requests_assigned_driver_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.requests
    ADD CONSTRAINT requests_assigned_driver_id_fkey FOREIGN KEY (assigned_driver_id) REFERENCES public.drivers(id) ON DELETE SET NULL;


--
-- TOC entry 3527 (class 2606 OID 16662)
-- Name: requests requests_company_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.requests
    ADD CONSTRAINT requests_company_id_fkey FOREIGN KEY (company_id) REFERENCES public.companies(id) ON DELETE SET NULL;


--
-- TOC entry 3528 (class 2606 OID 16657)
-- Name: requests requests_hotel_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.requests
    ADD CONSTRAINT requests_hotel_id_fkey FOREIGN KEY (hotel_id) REFERENCES public.hotels(id) ON DELETE SET NULL;


--
-- TOC entry 3531 (class 2606 OID 16698)
-- Name: reservation_timeline reservation_timeline_reservation_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.reservation_timeline
    ADD CONSTRAINT reservation_timeline_reservation_id_fkey FOREIGN KEY (reservation_id) REFERENCES public.reservations(id) ON DELETE CASCADE;


--
-- TOC entry 3529 (class 2606 OID 16812)
-- Name: reservations reservations_assigned_driver_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.reservations
    ADD CONSTRAINT reservations_assigned_driver_id_fkey FOREIGN KEY (assigned_driver_id) REFERENCES public.drivers(id) ON DELETE SET NULL;


--
-- TOC entry 3530 (class 2606 OID 16683)
-- Name: reservations reservations_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.reservations
    ADD CONSTRAINT reservations_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- TOC entry 3524 (class 2606 OID 16623)
-- Name: vehicle_photos vehicle_photos_vehicle_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.vehicle_photos
    ADD CONSTRAINT vehicle_photos_vehicle_id_fkey FOREIGN KEY (vehicle_id) REFERENCES public.vehicles(id) ON DELETE CASCADE;


--
-- TOC entry 3523 (class 2606 OID 16609)
-- Name: vehicles vehicles_driver_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.vehicles
    ADD CONSTRAINT vehicles_driver_id_fkey FOREIGN KEY (driver_id) REFERENCES public.drivers(id) ON DELETE SET NULL;


-- Completed on 2025-10-26 19:09:39 -03

--
-- PostgreSQL database dump complete
--

\unrestrict D39Vw3DHTIZI3zqNf08wH1To83VA1fRb6so7TcHD8PYJAEXNACTVIukfiEpXEWH

