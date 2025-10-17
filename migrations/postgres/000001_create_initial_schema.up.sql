-- Tables
CREATE TABLE public.action_logs (
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  user_id uuid,
  action text NOT NULL,
  created_at timestamp without time zone DEFAULT now(),
  meta jsonb,
  CONSTRAINT action_logs_pkey PRIMARY KEY (id),
  CONSTRAINT action_logs_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id)
);

CREATE TABLE public.deletion_capacity (
  day date NOT NULL,
  count integer NOT NULL DEFAULT 0,
  max_limit integer NOT NULL DEFAULT 10,
  updated_at timestamp without time zone DEFAULT now(),
  CONSTRAINT deletion_capacity_pkey PRIMARY KEY (day)
);

CREATE TABLE public.email_sends (
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  user_id uuid,
  type text NOT NULL CHECK (type = ANY (ARRAY['verification'::text, 'password_reset'::text])),
  sent_at timestamp without time zone DEFAULT now(),
  CONSTRAINT email_sends_pkey PRIMARY KEY (id),
  CONSTRAINT email_sends_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id)
);

CREATE TABLE public.magic_links (
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  user_id uuid,
  token text NOT NULL UNIQUE,
  expires_at timestamp without time zone NOT NULL,
  used boolean DEFAULT false,
  created_at timestamp without time zone DEFAULT now(),
  type text NOT NULL DEFAULT 'email_verification'::text CHECK (type = ANY (ARRAY['email_verification'::text, 'password_reset'::text])),
  used_at timestamp without time zone,
  is_active boolean DEFAULT true,
  CONSTRAINT magic_links_pkey PRIMARY KEY (id),
  CONSTRAINT magic_links_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id)
);

CREATE TABLE public.notifications (
  id uuid NOT NULL,
  type character varying NOT NULL,
  recipient character varying NOT NULL,
  subject character varying NOT NULL,
  body text NOT NULL,
  sent_at timestamp with time zone NOT NULL,
  created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT notifications_pkey PRIMARY KEY (id)
);

CREATE TABLE public.refresh_tokens (
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  user_id uuid,
  token text NOT NULL UNIQUE,
  expires_at timestamp without time zone NOT NULL,
  created_at timestamp without time zone DEFAULT now(),
  revoked boolean DEFAULT false,
  ip_address text,
  user_agent text,
  CONSTRAINT refresh_tokens_pkey PRIMARY KEY (id),
  CONSTRAINT refresh_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id)
);

CREATE TABLE public.user_deletion_cycles (
  user_id uuid NOT NULL,
  cycles integer NOT NULL DEFAULT 0,
  last_cycle_at timestamp without time zone,
  CONSTRAINT user_deletion_cycles_pkey PRIMARY KEY (user_id),
  CONSTRAINT user_deletion_cycles_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id)
);

CREATE TABLE public.user_deletions (
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL,
  scheduled_date date NOT NULL,
  executed boolean DEFAULT false,
  created_at timestamp without time zone DEFAULT now(),
  status text DEFAULT 'queued'::text CHECK (status = ANY (ARRAY['queued'::text, 'scheduled'::text, 'executed'::text, 'cancelled'::text])),
  token text UNIQUE,
  token_expires_at timestamp without time zone,
  recovery_token text UNIQUE,
  recovery_token_expires_at timestamp without time zone,
  CONSTRAINT user_deletions_pkey PRIMARY KEY (id),
  CONSTRAINT user_deletions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id)
);

CREATE TABLE public.user_logins (
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  user_id uuid,
  ip_address text,
  user_agent text,
  created_at timestamp without time zone DEFAULT now(),
  success boolean DEFAULT true,
  CONSTRAINT user_logins_pkey PRIMARY KEY (id),
  CONSTRAINT user_logins_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id)
);

CREATE TABLE public.users (
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  name text NOT NULL,
  email text NOT NULL UNIQUE,
  password_hash text NOT NULL,
  created_at timestamp without time zone DEFAULT now(),
  updated_at timestamp without time zone DEFAULT now(),
  is_email_verified boolean NOT NULL DEFAULT false,
  deleted_at timestamp without time zone,
  avatar_url text,
  deletion_due_at timestamp without time zone,
  last_login_at timestamp without time zone,
  is_deleted boolean DEFAULT false,
  CONSTRAINT users_pkey PRIMARY KEY (id)
);

-- Functions
CREATE OR REPLACE FUNCTION prevent_login_for_deleted_users()
RETURNS TRIGGER AS $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM public.users u WHERE u.id = NEW.user_id AND u.is_deleted = true
  ) THEN
    RAISE EXCEPTION 'User is deleted and cannot log in';
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_magic_link_is_active()
RETURNS TRIGGER AS $$
BEGIN
  NEW.is_active := NOT NEW.used AND NEW.expires_at > now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Triggers
CREATE TRIGGER trg_magic_link_is_active
BEFORE INSERT OR UPDATE ON public.magic_links
FOR EACH ROW EXECUTE FUNCTION update_magic_link_is_active();

CREATE TRIGGER trg_prevent_login_for_deleted_users
BEFORE INSERT ON public.user_logins
FOR EACH ROW EXECUTE FUNCTION prevent_login_for_deleted_users();

CREATE TRIGGER update_users_updated_at
BEFORE UPDATE ON public.users
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
