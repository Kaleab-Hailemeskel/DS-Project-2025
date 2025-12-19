-- Optional migrations for user-service PoC
-- Add account state and last_login timestamp
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT TRUE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_login TIMESTAMP;

-- Optional DB-based token blacklist (recommend Redis with TTL instead of DB writes)
CREATE TABLE IF NOT EXISTS token_blacklist (
    jti TEXT PRIMARY KEY,
    expires_at TIMESTAMP NOT NULL
);
