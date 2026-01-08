-- Create refresh_tokens table for token management and revocation
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(500) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    revoked_at TIMESTAMP NULL,
    replaced_by_token VARCHAR(500)
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user ON refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token ON refresh_tokens(token);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires ON refresh_tokens(expires_at);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_revoked ON refresh_tokens(revoked_at);

-- Add comments
COMMENT ON TABLE refresh_tokens IS 'Stores refresh tokens for secure token refresh and session management';
COMMENT ON COLUMN refresh_tokens.token IS 'Hashed refresh token value';
COMMENT ON COLUMN refresh_tokens.expires_at IS 'Token expiration timestamp (7-30 days based on role)';
COMMENT ON COLUMN refresh_tokens.revoked_at IS 'Timestamp when token was revoked (NULL = active)';
COMMENT ON COLUMN refresh_tokens.replaced_by_token IS 'New token that replaced this one (for token rotation)';
