-- migrate:up
-- Initialize the database with necessary extensions and settings

-- Enable UUID extension (needed for UUID generation)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Set the default encoding to UTF8 for client connections
SET client_encoding = 'UTF8';

-- Ensure standard conforming strings (to prevent SQL injection)
SET standard_conforming_strings = on;

-- Disable row-level security globally (can be adjusted per table if needed later)
SET row_security = off;

-- Optionally, set the default timezone (adjust as needed)
SET timezone = 'UTC';

-- migrate:down
-- Drop the UUID extension and revert settings if needed
DROP EXTENSION IF EXISTS "uuid-ossp";