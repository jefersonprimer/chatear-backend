-- WARNING: This schema is for context only and is not meant to be run.
-- Table order and constraints may not be valid for execution.

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


### functions
prevent_login_for_deleted_users
BEGIN
  IF EXISTS (
    SELECT 1 FROM public.users u WHERE u.id = NEW.user_id AND u.is_deleted = true
  ) THEN
    RAISE EXCEPTION 'User is deleted and cannot log in';
  END IF;
  RETURN NEW;
END;

update_magic_link_is_active
BEGIN
  NEW.is_active := NOT NEW.used AND NEW.expires_at > now();
  RETURN NEW;
END;

update_updated_at_column
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;


### triggers
[
  {
    "table_name": "buckets",
    "trigger_name": "enforce_bucket_name_length_trigger",
    "action_timing": "BEFORE",
    "event": "INSERT",
    "action_statement": "EXECUTE FUNCTION storage.enforce_bucket_name_length()"
  },
  {
    "table_name": "buckets",
    "trigger_name": "enforce_bucket_name_length_trigger",
    "action_timing": "BEFORE",
    "event": "UPDATE",
    "action_statement": "EXECUTE FUNCTION storage.enforce_bucket_name_length()"
  },
  {
    "table_name": "magic_links",
    "trigger_name": "trg_magic_link_is_active",
    "action_timing": "BEFORE",
    "event": "UPDATE",
    "action_statement": "EXECUTE FUNCTION update_magic_link_is_active()"
  },
  {
    "table_name": "magic_links",
    "trigger_name": "trg_magic_link_is_active",
    "action_timing": "BEFORE",
    "event": "INSERT",
    "action_statement": "EXECUTE FUNCTION update_magic_link_is_active()"
  },
  {
    "table_name": "objects",
    "trigger_name": "objects_delete_delete_prefix",
    "action_timing": "AFTER",
    "event": "DELETE",
    "action_statement": "EXECUTE FUNCTION storage.delete_prefix_hierarchy_trigger()"
  },
  {
    "table_name": "objects",
    "trigger_name": "objects_insert_create_prefix",
    "action_timing": "BEFORE",
    "event": "INSERT",
    "action_statement": "EXECUTE FUNCTION storage.objects_insert_prefix_trigger()"
  },
  {
    "table_name": "objects",
    "trigger_name": "objects_update_create_prefix",
    "action_timing": "BEFORE",
    "event": "UPDATE",
    "action_statement": "EXECUTE FUNCTION storage.objects_update_prefix_trigger()"
  },
  {
    "table_name": "objects",
    "trigger_name": "update_objects_updated_at",
    "action_timing": "BEFORE",
    "event": "UPDATE",
    "action_statement": "EXECUTE FUNCTION storage.update_updated_at_column()"
  },
  {
    "table_name": "prefixes",
    "trigger_name": "prefixes_create_hierarchy",
    "action_timing": "BEFORE",
    "event": "INSERT",
    "action_statement": "EXECUTE FUNCTION storage.prefixes_insert_trigger()"
  },
  {
    "table_name": "prefixes",
    "trigger_name": "prefixes_delete_hierarchy",
    "action_timing": "AFTER",
    "event": "DELETE",
    "action_statement": "EXECUTE FUNCTION storage.delete_prefix_hierarchy_trigger()"
  },
  {
    "table_name": "subscription",
    "trigger_name": "tr_check_filters",
    "action_timing": "BEFORE",
    "event": "INSERT",
    "action_statement": "EXECUTE FUNCTION realtime.subscription_check_filters()"
  },
  {
    "table_name": "subscription",
    "trigger_name": "tr_check_filters",
    "action_timing": "BEFORE",
    "event": "UPDATE",
    "action_statement": "EXECUTE FUNCTION realtime.subscription_check_filters()"
  },
  {
    "table_name": "user_logins",
    "trigger_name": "trg_prevent_login_for_deleted_users",
    "action_timing": "BEFORE",
    "event": "INSERT",
    "action_statement": "EXECUTE FUNCTION prevent_login_for_deleted_users()"
  },
  {
    "table_name": "users",
    "trigger_name": "update_users_updated_at",
    "action_timing": "BEFORE",
    "event": "UPDATE",
    "action_statement": "EXECUTE FUNCTION update_updated_at_column()"
  }
]

### indexes
[
  {
    "schemaname": "public",
    "tablename": "action_logs",
    "indexname": "action_logs_pkey",
    "indexdef": "CREATE UNIQUE INDEX action_logs_pkey ON public.action_logs USING btree (id)"
  },
  {
    "schemaname": "public",
    "tablename": "action_logs",
    "indexname": "idx_action_logs_user_action",
    "indexdef": "CREATE INDEX idx_action_logs_user_action ON public.action_logs USING btree (user_id, action, created_at DESC)"
  },
  {
    "schemaname": "public",
    "tablename": "deletion_capacity",
    "indexname": "deletion_capacity_pkey",
    "indexdef": "CREATE UNIQUE INDEX deletion_capacity_pkey ON public.deletion_capacity USING btree (day)"
  },
  {
    "schemaname": "public",
    "tablename": "email_sends",
    "indexname": "email_sends_pkey",
    "indexdef": "CREATE UNIQUE INDEX email_sends_pkey ON public.email_sends USING btree (id)"
  },
  {
    "schemaname": "public",
    "tablename": "email_sends",
    "indexname": "idx_email_sends_user_type_sent_at",
    "indexdef": "CREATE INDEX idx_email_sends_user_type_sent_at ON public.email_sends USING btree (user_id, type, sent_at)"
  },
  {
    "schemaname": "public",
    "tablename": "email_sends",
    "indexname": "idx_email_sends_user_type_time",
    "indexdef": "CREATE INDEX idx_email_sends_user_type_time ON public.email_sends USING btree (user_id, type, sent_at DESC)"
  },
  {
    "schemaname": "public",
    "tablename": "magic_links",
    "indexname": "idx_magic_links_active",
    "indexdef": "CREATE INDEX idx_magic_links_active ON public.magic_links USING btree (token) WHERE (is_active = true)"
  },
  {
    "schemaname": "public",
    "tablename": "magic_links",
    "indexname": "idx_magic_links_type_token",
    "indexdef": "CREATE INDEX idx_magic_links_type_token ON public.magic_links USING btree (type, token)"
  },
  {
    "schemaname": "public",
    "tablename": "magic_links",
    "indexname": "idx_magic_links_unused",
    "indexdef": "CREATE INDEX idx_magic_links_unused ON public.magic_links USING btree (token) WHERE (used = false)"
  },
  {
    "schemaname": "public",
    "tablename": "magic_links",
    "indexname": "idx_magic_links_user_type_active",
    "indexdef": "CREATE INDEX idx_magic_links_user_type_active ON public.magic_links USING btree (user_id, type, is_active)"
  },
  {
    "schemaname": "public",
    "tablename": "magic_links",
    "indexname": "magic_links_pkey",
    "indexdef": "CREATE UNIQUE INDEX magic_links_pkey ON public.magic_links USING btree (id)"
  },
  {
    "schemaname": "public",
    "tablename": "magic_links",
    "indexname": "magic_links_token_key",
    "indexdef": "CREATE UNIQUE INDEX magic_links_token_key ON public.magic_links USING btree (token)"
  },
  {
    "schemaname": "public",
    "tablename": "refresh_tokens",
    "indexname": "idx_refresh_token_token",
    "indexdef": "CREATE INDEX idx_refresh_token_token ON public.refresh_tokens USING btree (token)"
  },
  {
    "schemaname": "public",
    "tablename": "refresh_tokens",
    "indexname": "idx_refresh_tokens_user_revoked",
    "indexdef": "CREATE INDEX idx_refresh_tokens_user_revoked ON public.refresh_tokens USING btree (user_id, revoked)"
  },
  {
    "schemaname": "public",
    "tablename": "refresh_tokens",
    "indexname": "refresh_tokens_pkey",
    "indexdef": "CREATE UNIQUE INDEX refresh_tokens_pkey ON public.refresh_tokens USING btree (id)"
  },
  {
    "schemaname": "public",
    "tablename": "refresh_tokens",
    "indexname": "refresh_tokens_token_unique",
    "indexdef": "CREATE UNIQUE INDEX refresh_tokens_token_unique ON public.refresh_tokens USING btree (token)"
  },
  {
    "schemaname": "public",
    "tablename": "user_deletion_cycles",
    "indexname": "user_deletion_cycles_pkey",
    "indexdef": "CREATE UNIQUE INDEX user_deletion_cycles_pkey ON public.user_deletion_cycles USING btree (user_id)"
  },
  {
    "schemaname": "public",
    "tablename": "user_deletions",
    "indexname": "idx_user_deletions_scheduled_date",
    "indexdef": "CREATE INDEX idx_user_deletions_scheduled_date ON public.user_deletions USING btree (scheduled_date)"
  },
  {
    "schemaname": "public",
    "tablename": "user_deletions",
    "indexname": "idx_user_deletions_status_date",
    "indexdef": "CREATE INDEX idx_user_deletions_status_date ON public.user_deletions USING btree (status, scheduled_date)"
  },
  {
    "schemaname": "public",
    "tablename": "user_deletions",
    "indexname": "user_deletions_pkey",
    "indexdef": "CREATE UNIQUE INDEX user_deletions_pkey ON public.user_deletions USING btree (id)"
  },
  {
    "schemaname": "public",
    "tablename": "user_deletions",
    "indexname": "user_deletions_recovery_token_key",
    "indexdef": "CREATE UNIQUE INDEX user_deletions_recovery_token_key ON public.user_deletions USING btree (recovery_token)"
  },
  {
    "schemaname": "public",
    "tablename": "user_deletions",
    "indexname": "user_deletions_token_key",
    "indexdef": "CREATE UNIQUE INDEX user_deletions_token_key ON public.user_deletions USING btree (token)"
  },
  {
    "schemaname": "public",
    "tablename": "user_logins",
    "indexname": "user_logins_pkey",
    "indexdef": "CREATE UNIQUE INDEX user_logins_pkey ON public.user_logins USING btree (id)"
  },
  {
    "schemaname": "public",
    "tablename": "users",
    "indexname": "idx_users_deleted_at",
    "indexdef": "CREATE INDEX idx_users_deleted_at ON public.users USING btree (deleted_at)"
  },
  {
    "schemaname": "public",
    "tablename": "users",
    "indexname": "idx_users_email_lower",
    "indexdef": "CREATE INDEX idx_users_email_lower ON public.users USING btree (lower(email))"
  },
  {
    "schemaname": "public",
    "tablename": "users",
    "indexname": "idx_users_is_email_verified",
    "indexdef": "CREATE INDEX idx_users_is_email_verified ON public.users USING btree (is_email_verified)"
  },
  {
    "schemaname": "public",
    "tablename": "users",
    "indexname": "users_email_key",
    "indexdef": "CREATE UNIQUE INDEX users_email_key ON public.users USING btree (email)"
  },
  {
    "schemaname": "public",
    "tablename": "users",
    "indexname": "users_pkey",
    "indexdef": "CREATE UNIQUE INDEX users_pkey ON public.users USING btree (id)"
  }
]
