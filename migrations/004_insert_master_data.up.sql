-- Insert master data for testing and initial setup

-- Create options_master table for option management
CREATE TABLE options_master (
    id SERIAL PRIMARY KEY,
    option_type VARCHAR(10) NOT NULL UNIQUE,
    option_name VARCHAR(100) NOT NULL,
    description TEXT,
    plan_compatibility VARCHAR(10) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Insert option master data
INSERT INTO options_master (option_type, option_name, description, plan_compatibility) VALUES
('AA', 'AAオプション', 'Aプラン専用のオプションサービス', 'A'),
('BB', 'BBオプション', 'Bプラン専用のオプションサービス', 'B'),
('AB', 'ABオプション', 'A・B両プラン共通のオプションサービス', 'AB');

-- Create prefectures_master table for address validation
CREATE TABLE prefectures_master (
    id SERIAL PRIMARY KEY,
    prefecture_code CHAR(2) NOT NULL UNIQUE,
    prefecture_name VARCHAR(10) NOT NULL UNIQUE,
    region VARCHAR(20) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Insert prefecture master data (main prefectures for demonstration)
INSERT INTO prefectures_master (prefecture_code, prefecture_name, region) VALUES
('01', '北海道', '北海道'),
('13', '東京都', '関東'),
('14', '神奈川県', '関東'),
('23', '愛知県', '中部'),
('27', '大阪府', '関西'),
('28', '兵庫県', '関西'),
('40', '福岡県', '九州');

-- Create indexes
CREATE INDEX idx_options_master_option_type ON options_master(option_type);
CREATE INDEX idx_options_master_plan_compatibility ON options_master(plan_compatibility);
CREATE INDEX idx_prefectures_master_code ON prefectures_master(prefecture_code);
CREATE INDEX idx_prefectures_master_name ON prefectures_master(prefecture_name);

-- Add constraints
ALTER TABLE options_master ADD CONSTRAINT chk_options_master_plan_compatibility 
    CHECK (plan_compatibility IN ('A', 'B', 'AB'));

-- Add comments
COMMENT ON TABLE options_master IS 'Master data for available options';
COMMENT ON TABLE prefectures_master IS 'Master data for Japanese prefectures';
COMMENT ON COLUMN options_master.option_type IS 'Option type identifier';
COMMENT ON COLUMN options_master.option_name IS 'Display name for the option';
COMMENT ON COLUMN options_master.plan_compatibility IS 'Which plans can use this option';
COMMENT ON COLUMN prefectures_master.prefecture_code IS 'JIS prefecture code';
COMMENT ON COLUMN prefectures_master.prefecture_name IS 'Prefecture name';
COMMENT ON COLUMN prefectures_master.region IS 'Geographic region';