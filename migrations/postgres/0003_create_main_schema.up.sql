-- Create main schema based on docs/database.md

-- Create users table
CREATE TABLE IF NOT EXISTS public.users (
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

-- Create action_logs table
CREATE TABLE IF NOT EXISTS public.action_logs (
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  user_id uuid,
  action text NOT NULL,
  created_at timestamp without time zone DEFAULT now(),
  meta jsonb,
  CONSTRAINT action_logs_pkey PRIMARY KEY (id),
  CONSTRAINT action_logs_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id)
);

-- Create deletion_capacity table
CREATE TABLE IF NOT EXISTS public.deletion_capacity (
  day date NOT NULL,
  count integer NOT NULL DEFAULT 0,
  max_limit integer NOT NULL DEFAULT 10,
  updated_at timestamp without time zone DEFAULT now(),
  CONSTRAINT deletion_capacity_pkey PRIMARY KEY (day)
);

-- Create email_sends table
CREATE TABLE IF NOT EXISTS public.email_sends (
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  user_id uuid,
  type text NOT NULL CHECK (type = ANY (ARRAY['verification'::text, 'password_reset'::text])),
  sent_at timestamp without time zone DEFAULT now(),
  CONSTRAINT email_sends_pkey PRIMARY KEY (id),
  CONSTRAINT email_sends_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id)
);

-- Create magic_links table
CREATE TABLE IF NOT EXISTS public.magic_links (
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

-- Create user_deletion_cycles table
CREATE TABLE IF NOT EXISTS public.user_deletion_cycles (
  user_id uuid NOT NULL,
  cycles integer NOT NULL DEFAULT 0,
  last_cycle_at timestamp without time zone,
  CONSTRAINT user_deletion_cycles_pkey PRIMARY KEY (user_id),
  CONSTRAINT user_deletion_cycles_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id)
);

-- Create user_deletions table
CREATE TABLE IF NOT EXISTS public.user_deletions (
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

-- Create user_logins table
CREATE TABLE IF NOT EXISTS public.user_logins (
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  user_id uuid,
  ip_address text,
  user_agent text,
  created_at timestamp without time zone DEFAULT now(),
  success boolean DEFAULT true,
  CONSTRAINT user_logins_pkey PRIMARY KEY (id),
  CONSTRAINT user_logins_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id)
);
