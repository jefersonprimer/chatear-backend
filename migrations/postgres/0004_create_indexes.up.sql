-- Create indexes based on docs/database.md

-- Users table indexes
CREATE INDEX IF NOT EXISTS idx_users_email_lower ON public.users USING btree (lower(email));
CREATE INDEX IF NOT EXISTS idx_users_is_email_verified ON public.users USING btree (is_email_verified);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON public.users USING btree (deleted_at);

-- Action logs indexes
CREATE INDEX IF NOT EXISTS idx_action_logs_user_action ON public.action_logs USING btree (user_id, action, created_at DESC);

-- Email sends indexes
CREATE INDEX IF NOT EXISTS idx_email_sends_user_type_sent_at ON public.email_sends USING btree (user_id, type, sent_at);
CREATE INDEX IF NOT EXISTS idx_email_sends_user_type_time ON public.email_sends USING btree (user_id, type, sent_at DESC);

-- Magic links indexes
CREATE INDEX IF NOT EXISTS idx_magic_links_type_token ON public.magic_links USING btree (type, token);
CREATE INDEX IF NOT EXISTS idx_magic_links_user_type_active ON public.magic_links USING btree (user_id, type, is_active);
CREATE INDEX IF NOT EXISTS idx_magic_links_active ON public.magic_links USING btree (token) WHERE (is_active = true);
CREATE INDEX IF NOT EXISTS idx_magic_links_unused ON public.magic_links USING btree (token) WHERE (used = false);

-- Refresh tokens indexes
CREATE INDEX IF NOT EXISTS idx_refresh_token_token ON public.refresh_tokens USING btree (token);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_revoked ON public.refresh_tokens USING btree (user_id, revoked);

-- User deletions indexes
CREATE INDEX IF NOT EXISTS idx_user_deletions_scheduled_date ON public.user_deletions USING btree (scheduled_date);
CREATE INDEX IF NOT EXISTS idx_user_deletions_status_date ON public.user_deletions USING btree (status, scheduled_date);
