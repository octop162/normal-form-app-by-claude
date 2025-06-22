-- Create user_options table for user selected options
CREATE TABLE user_options (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    option_type VARCHAR(10) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_user_options_user_id ON user_options(user_id);
CREATE INDEX idx_user_options_option_type ON user_options(option_type);
CREATE INDEX idx_user_options_created_at ON user_options(created_at);

-- Add constraints
ALTER TABLE user_options ADD CONSTRAINT chk_user_options_option_type 
    CHECK (option_type IN ('AA', 'BB', 'AB'));

-- Add unique constraint to prevent duplicate option selections
ALTER TABLE user_options ADD CONSTRAINT uk_user_options_user_option 
    UNIQUE (user_id, option_type);

-- Add comments
COMMENT ON TABLE user_options IS 'User selected options for membership plans';
COMMENT ON COLUMN user_options.id IS 'Primary key';
COMMENT ON COLUMN user_options.user_id IS 'Foreign key to users table';
COMMENT ON COLUMN user_options.option_type IS 'Option type: AA (A plan only), BB (B plan only), AB (both plans)';
COMMENT ON COLUMN user_options.created_at IS 'Record creation timestamp';