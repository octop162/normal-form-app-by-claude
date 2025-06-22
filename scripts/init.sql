-- Initialize Database for Normal Form App
-- This script sets up the initial database structure

-- Create extension for UUID generation if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create initial database (this might already exist from environment)
-- Database creation is handled by Docker environment variables

-- Set timezone
SET timezone = 'Asia/Tokyo';

-- Create a simple health check table for initial testing
CREATE TABLE IF NOT EXISTS health_check (
    id SERIAL PRIMARY KEY,
    status VARCHAR(20) NOT NULL DEFAULT 'ok',
    checked_at TIMESTAMP DEFAULT NOW()
);

-- Insert initial health check record
INSERT INTO health_check (status) VALUES ('ok') ON CONFLICT DO NOTHING;

-- Grant necessary permissions
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO postgres;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO postgres;