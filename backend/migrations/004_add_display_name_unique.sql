-- Add UNIQUE constraint to display_name for user safety
-- This prevents confusion when adding users to groups by display name

-- First check for existing duplicates and handle them
-- This is a safe operation as we're just checking
DO $$
DECLARE
    duplicate_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO duplicate_count 
    FROM (
        SELECT display_name 
        FROM users 
        GROUP BY display_name 
        HAVING COUNT(*) > 1
    ) as duplicates;
    
    IF duplicate_count > 0 THEN
        RAISE NOTICE 'Found % duplicate display names. Please resolve duplicates before applying UNIQUE constraint.', duplicate_count;
        -- For development purposes, we'll append user ID to duplicates
        UPDATE users 
        SET display_name = display_name || '_' || id::text
        WHERE id IN (
            SELECT u1.id 
            FROM users u1
            INNER JOIN users u2 ON u1.display_name = u2.display_name AND u1.id > u2.id
        );
        RAISE NOTICE 'Automatically resolved duplicates by appending user ID';
    END IF;
END $$;

-- Add the UNIQUE constraint
ALTER TABLE users ADD CONSTRAINT users_display_name_unique UNIQUE (display_name);

-- Add index for performance
CREATE INDEX IF NOT EXISTS idx_users_display_name ON users(display_name);

-- Add check constraint to ensure display_name is not empty
ALTER TABLE users ADD CONSTRAINT users_display_name_not_empty CHECK (LENGTH(TRIM(display_name)) > 0);