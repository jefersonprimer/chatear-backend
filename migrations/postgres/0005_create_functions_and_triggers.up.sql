-- Create functions and triggers based on docs/database.md

-- Function to prevent login for deleted users
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

-- Function to update magic link is_active status
CREATE OR REPLACE FUNCTION update_magic_link_is_active()
RETURNS TRIGGER AS $$
BEGIN
  NEW.is_active := NOT NEW.used AND NEW.expires_at > now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Function to update updated_at column
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers
CREATE TRIGGER trg_magic_link_is_active
  BEFORE INSERT OR UPDATE ON public.magic_links
  FOR EACH ROW
  EXECUTE FUNCTION update_magic_link_is_active();

CREATE TRIGGER trg_prevent_login_for_deleted_users
  BEFORE INSERT ON public.user_logins
  FOR EACH ROW
  EXECUTE FUNCTION prevent_login_for_deleted_users();

CREATE TRIGGER update_users_updated_at
  BEFORE UPDATE ON public.users
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();
