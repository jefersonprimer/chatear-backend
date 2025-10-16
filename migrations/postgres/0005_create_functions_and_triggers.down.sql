-- Drop triggers and functions

-- Drop triggers
DROP TRIGGER IF EXISTS update_users_updated_at ON public.users;
DROP TRIGGER IF EXISTS trg_prevent_login_for_deleted_users ON public.user_logins;
DROP TRIGGER IF EXISTS trg_magic_link_is_active ON public.magic_links;

-- Drop functions
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP FUNCTION IF EXISTS update_magic_link_is_active();
DROP FUNCTION IF EXISTS prevent_login_for_deleted_users();
