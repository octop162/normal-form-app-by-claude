-- Create user_sessions table for temporary data storage
CREATE TABLE user_sessions (
    id VARCHAR(255) PRIMARY KEY,
    user_data JSONB NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_user_sessions_expires_at ON user_sessions(expires_at);
CREATE INDEX idx_user_sessions_created_at ON user_sessions(created_at);

-- Add comments
COMMENT ON TABLE user_sessions IS 'Temporary user session data for form persistence';
COMMENT ON COLUMN user_sessions.id IS 'Session ID (UUID or similar)';
COMMENT ON COLUMN user_sessions.user_data IS 'JSON data containing form information';
COMMENT ON COLUMN user_sessions.expires_at IS 'Session expiration timestamp';
COMMENT ON COLUMN user_sessions.created_at IS 'Record creation timestamp';
COMMENT ON COLUMN user_sessions.updated_at IS 'Record update timestamp';