-- Insert admin user (password: "admin123")
INSERT INTO users (email, password_hash, display_name, role, subscription_status, is_active, email_verified)
VALUES (
    'admin@example.com',
    '$2a$10$mMoyaRT/vy75P/.QjLFbA.9cU5ozG6Is0D2X/0DuHzU5yJ2lRqJGu', -- bcrypt hash for "admin123"
    'Admin User',
    'admin',
    'active',
    true,
    true
) ON CONFLICT (email) DO NOTHING;

-- Insert test user (password: "user123")
INSERT INTO users (email, password_hash, display_name, role, subscription_status, is_active, email_verified)
VALUES (
    'user@example.com',
    '$2a$10$6wm5JXQaRirOdksBxMGzLO9ZkEeSNZ/3by70kCrFfOMI407y8DuHzU5yJ2lRqJGu', -- bcrypt hash for "user123"
    'Test User',
    'user',
    'active',
    true,
    true
) ON CONFLICT (email) DO NOTHING;