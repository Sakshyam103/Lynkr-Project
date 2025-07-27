-- Migration: 002_add_geolocation
-- Created: 2023-08-15
-- Description: Add geolocation data to attendances table

-- Up Migration
-- Add latitude and longitude columns to attendances table
ALTER TABLE attendances ADD COLUMN latitude REAL;
ALTER TABLE attendances ADD COLUMN longitude REAL;

-- Create index on latitude and longitude for potential spatial queries
CREATE INDEX IF NOT EXISTS idx_attendances_location ON attendances(latitude, longitude);

-- Record this migration
INSERT INTO migrations (name) VALUES ('002_add_geolocation');

-- Down Migration
-- DROP INDEX IF EXISTS idx_attendances_location;
-- ALTER TABLE attendances DROP COLUMN longitude;
-- ALTER TABLE attendances DROP COLUMN latitude;