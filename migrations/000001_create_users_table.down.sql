BEGIN;

-- The order of dropping tables is important due to constraints!

-- Drop the user_permissions table
DROP TABLE IF EXISTS user_permissions;

-- Drop the permissions table
DROP TABLE IF EXISTS permissions;

-- Drop the tokens table
DROP TABLE IF EXISTS tokens;

-- Drop the users table
DROP TABLE IF EXISTS users;

COMMIT;