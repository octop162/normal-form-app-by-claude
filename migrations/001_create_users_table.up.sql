-- Create users table for member registration
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    last_name VARCHAR(15) NOT NULL,
    first_name VARCHAR(15) NOT NULL,
    last_name_kana VARCHAR(15) NOT NULL,
    first_name_kana VARCHAR(15) NOT NULL,
    phone1 VARCHAR(5) NOT NULL,
    phone2 VARCHAR(4) NOT NULL,
    phone3 VARCHAR(4) NOT NULL,
    postal_code1 CHAR(3) NOT NULL,
    postal_code2 CHAR(4) NOT NULL,
    prefecture VARCHAR(10) NOT NULL,
    city VARCHAR(50) NOT NULL,
    town VARCHAR(50),
    chome VARCHAR(10),
    banchi VARCHAR(10) NOT NULL,
    go VARCHAR(10),
    building VARCHAR(100),
    room VARCHAR(20),
    email VARCHAR(256) NOT NULL UNIQUE,
    plan_type VARCHAR(10) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_plan_type ON users(plan_type);
CREATE INDEX idx_users_created_at ON users(created_at);
CREATE INDEX idx_users_postal_code ON users(postal_code1, postal_code2);

-- Add constraints
ALTER TABLE users ADD CONSTRAINT chk_users_plan_type 
    CHECK (plan_type IN ('A', 'B'));

-- Add comments
COMMENT ON TABLE users IS 'User registration information';
COMMENT ON COLUMN users.id IS 'Primary key';
COMMENT ON COLUMN users.last_name IS 'User last name (max 15 chars)';
COMMENT ON COLUMN users.first_name IS 'User first name (max 15 chars)';
COMMENT ON COLUMN users.last_name_kana IS 'User last name in katakana (max 15 chars)';
COMMENT ON COLUMN users.first_name_kana IS 'User first name in katakana (max 15 chars)';
COMMENT ON COLUMN users.phone1 IS 'Phone number part 1 (area code)';
COMMENT ON COLUMN users.phone2 IS 'Phone number part 2 (exchange)';
COMMENT ON COLUMN users.phone3 IS 'Phone number part 3 (number)';
COMMENT ON COLUMN users.postal_code1 IS 'Postal code first 3 digits';
COMMENT ON COLUMN users.postal_code2 IS 'Postal code last 4 digits';
COMMENT ON COLUMN users.prefecture IS 'Prefecture name';
COMMENT ON COLUMN users.city IS 'City name';
COMMENT ON COLUMN users.town IS 'Town name (optional)';
COMMENT ON COLUMN users.chome IS 'Chome (optional)';
COMMENT ON COLUMN users.banchi IS 'Banchi (house number)';
COMMENT ON COLUMN users.go IS 'Go (optional)';
COMMENT ON COLUMN users.building IS 'Building name (optional)';
COMMENT ON COLUMN users.room IS 'Room number (optional)';
COMMENT ON COLUMN users.email IS 'Email address (unique)';
COMMENT ON COLUMN users.plan_type IS 'Plan type: A or B';
COMMENT ON COLUMN users.created_at IS 'Record creation timestamp';
COMMENT ON COLUMN users.updated_at IS 'Record update timestamp';