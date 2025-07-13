-- Insert admin user (password: "admin123")
INSERT INTO users (email, password_hash, display_name, role, subscription_status, is_active, email_verified)
VALUES (
    'admin@example.com',
    '$2a$10$mVwEKKOa5M2YNVuZlowdXeJ3VnwPjmwkKQwcRLlFrzkqWRn2yTzU2', -- bcrypt hash for "admin123"
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
    '$2a$10$9Y9K8HcK2L.kDjZGQ8WJcewH6GN5cI4OYi7uQ2LpQpN4nHX2KJF2W', -- bcrypt hash for "user123"
    'Test User',
    'user',
    'active',
    true,
    true
) ON CONFLICT (email) DO NOTHING;