-- Create devices table
CREATE TABLE IF NOT EXISTS devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    ip_address VARCHAR(45) NOT NULL,
    location VARCHAR(200),
    status VARCHAR(20) DEFAULT 'offline',
    version VARCHAR(50),
    last_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index for faster queries
CREATE INDEX idx_devices_status ON devices(status);
CREATE INDEX idx_devices_name ON devices(name);

-- Insert sample data
INSERT INTO devices (name, ip_address, location, status, last_seen) VALUES
    ('Router-BDG-01', '192.168.100.11', 'Bandung', 'online', NOW()),
    ('Router-JKT-01', '192.168.100.12', 'Jakarta', 'online', NOW()),
    ('Router-SBY-01', '192.168.100.13', 'Surabaya', 'offline', NOW() - INTERVAL '1 hour')
ON CONFLICT (name) DO NOTHING;
