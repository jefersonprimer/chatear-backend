-- Drop Triggers
DROP TRIGGER IF EXISTS trg_magic_link_is_active ON public.magic_links;
DROP TRIGGER IF EXISTS trg_prevent_login_for_deleted_users ON public.user_logins;
DROP TRIGGER IF EXISTS update_users_updated_at ON public.users;

-- Drop Functions
DROP FUNCTION IF EXISTS prevent_login_for_deleted_users();
DROP FUNCTION IF EXISTS update_magic_link_is_active();
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop Tables (in reverse order of dependency)
DROP TABLE IF EXISTS public.user_logins;
DROP TABLE IF EXISTS public.user_deletions;
DROP TABLE IF EXISTS public.user_deletion_cycles;
DROP TABLE IF EXISTS public.refresh_tokens;
DROP TABLE IF EXISTS public.magic_links;
DROP TABLE IF EXISTS public.email_sends;
DROP TABLE IF EXISTS public.action_logs;
DROP TABLE IF EXISTS public.notifications;
DROP TABLE IF EXISTS public.deletion_capacity;
DROP TABLE IF EXISTS public.users;
