BEGIN;

-- Enable the citext extension if it's not already enabled
CREATE EXTENSION IF NOT EXISTS citext;

-- Create the users table
CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    email citext UNIQUE NOT NULL,
    password bytea NOT NULL,
    activated bool NOT NULL DEFAULT false,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 0
);

-- Create the tokens table
CREATE TABLE IF NOT EXISTS tokens (
    hash bytea PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    expiry timestamp with time zone NOT NULL,
    scope text NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp with time zone NOT NULL DEFAULT NOW()
);

-- Create the permissions table
CREATE TABLE IF NOT EXISTS permissions (
    id bigserial PRIMARY KEY,
    code text UNIQUE NOT NULL
);

-- Create the user <-> permissions table
CREATE TABLE IF NOT EXISTS user_permissions (
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    permission_id bigint NOT NULL REFERENCES permissions ON DELETE CASCADE,
    PRIMARY KEY (user_id, permission_id)
);

-- Add admin and superadmin roles
INSERT INTO permissions (code)
VALUES
    ('admin'),
    ('superadmin')
ON CONFLICT DO NOTHING;

COMMIT;