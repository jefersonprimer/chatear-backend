-- Drop main schema tables in reverse order

DROP TABLE IF EXISTS public.user_logins;
DROP TABLE IF EXISTS public.user_deletions;
DROP TABLE IF EXISTS public.user_deletion_cycles;
DROP TABLE IF EXISTS public.magic_links;
DROP TABLE IF EXISTS public.email_sends;
DROP TABLE IF EXISTS public.deletion_capacity;
DROP TABLE IF EXISTS public.action_logs;
DROP TABLE IF EXISTS public.users;
