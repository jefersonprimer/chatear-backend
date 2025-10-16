-- Drop indexes

-- User deletions indexes
DROP INDEX IF EXISTS idx_user_deletions_status_date;
DROP INDEX IF EXISTS idx_user_deletions_scheduled_date;

-- Refresh tokens indexes
DROP INDEX IF EXISTS idx_refresh_tokens_user_revoked;
DROP INDEX IF EXISTS idx_refresh_token_token;

-- Magic links indexes
DROP INDEX IF EXISTS idx_magic_links_unused;
DROP INDEX IF EXISTS idx_magic_links_active;
DROP INDEX IF EXISTS idx_magic_links_user_type_active;
DROP INDEX IF EXISTS idx_magic_links_type_token;

-- Email sends indexes
DROP INDEX IF EXISTS idx_email_sends_user_type_time;
DROP INDEX IF EXISTS idx_email_sends_user_type_sent_at;

-- Action logs indexes
DROP INDEX IF EXISTS idx_action_logs_user_action;

-- Users table indexes
DROP INDEX IF EXISTS idx_users_deleted_at;
DROP INDEX IF EXISTS idx_users_is_email_verified;
DROP INDEX IF EXISTS idx_users_email_lower;
