-- Add the citext extension to support case-insensitive text types.
-- The `citext` extension provides a case-insensitive text type for PostgreSQL.
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users
(
    id         BIGSERIAL PRIMARY KEY,
    username   CITEXT                   NOT NULL UNIQUE,
    email      CITEXT                   NOT NULL UNIQUE,
    password   TEXT                     NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    version    BIGINT                   NOT NULL DEFAULT 1
);